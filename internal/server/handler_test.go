//go:build ignore
// +build ignore

package server

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/javanhut/CarrionLSP/internal/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

// testStream implements jsonrpc2.ObjectStream for testing
type testStream struct {
	in  chan json.RawMessage
	out chan json.RawMessage
}

func newTestStream() *testStream {
	return &testStream{
		in:  make(chan json.RawMessage, 10),
		out: make(chan json.RawMessage, 10),
	}
}

func (s *testStream) ReadObject(v interface{}) error {
	msg, ok := <-s.in
	if !ok {
		return io.EOF
	}
	return json.Unmarshal(msg, v)
}

func (s *testStream) WriteObject(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	s.out <- data
	return nil
}

func (s *testStream) Close() error {
	close(s.in)
	close(s.out)
	return nil
}

// testHandler captures method calls for verification
type testHandler struct {
	replies       []interface{}
	errors        []*jsonrpc2.Error
	notifications []string
}

func (h *testHandler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// This is used to capture responses from the server
}

// createTestConn creates a test connection for handler testing
func createTestConn(t *testing.T, handler *Handler) (*jsonrpc2.Conn, *testStream) {
	stream := newTestStream()
	conn := jsonrpc2.NewConn(
		context.Background(),
		stream,
		handler,
	)
	return conn, stream
}

// sendRequest sends a request through the test stream
func sendRequest(t *testing.T, stream *testStream, method string, params interface{}) {
	req := &jsonrpc2.Request{
		Method: method,
		ID:     jsonrpc2.ID{Num: 1},
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			t.Fatal(err)
		}
		raw := json.RawMessage(data)
		req.Params = &raw
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	stream.in <- reqData
}

// getResponse reads a response from the test stream
func getResponse(t *testing.T, stream *testStream) *jsonrpc2.Response {
	select {
	case msg := <-stream.out:
		var resp jsonrpc2.Response
		if err := json.Unmarshal(msg, &resp); err != nil {
			t.Fatal(err)
		}
		return &resp
	default:
		return nil
	}
}

func TestHandler_NewHandler(t *testing.T) {
	handler := NewHandler()

	if handler == nil {
		t.Fatal("Expected handler to be created")
	}

	if handler.analyzer == nil {
		t.Error("Expected analyzer to be initialized")
	}

	if handler.workspaces == nil {
		t.Error("Expected workspaces map to be initialized")
	}

	if handler.initialized {
		t.Error("Expected handler to start uninitialized")
	}
}

func TestHandler_Initialize(t *testing.T) {
	handler := NewHandler()
	conn, stream := createTestConn(t, handler)
	defer conn.Close()

	params := protocol.InitializeParams{
		Capabilities: &protocol.ClientCapabilities{},
	}

	// Start handling in background
	go func() {
		<-conn.DisconnectNotify()
	}()

	// Send initialize request
	sendRequest(t, stream, "initialize", params)

	// Get response
	resp := getResponse(t, stream)
	if resp == nil {
		t.Fatal("Expected response")
	}

	if resp.Error != nil {
		t.Errorf("Expected no error, got %v", resp.Error)
	}

	// Check that result is InitializeResult
	var result protocol.InitializeResult
	if err := json.Unmarshal(*resp.Result, &result); err != nil {
		t.Fatal(err)
	}

	if result.Capabilities.TextDocumentSync == nil {
		t.Error("Expected TextDocumentSync capability to be set")
		if result.Capabilities.CompletionProvider == nil {
			t.Error("Expected CompletionProvider capability to be set")
		}
		if result.ServerInfo == nil {
			t.Error("Expected ServerInfo to be set")
		}
		if result.ServerInfo.Name != "Carrion Language Server" {
			t.Errorf("Expected server name to be 'Carrion Language Server', got %s", result.ServerInfo.Name)
		}
	} else {
		t.Error("Expected reply to be InitializeResult")
	}
}

func TestHandler_Initialize_WithWorkspace(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	rootURI := "file:///test/workspace"
	params := protocol.InitializeParams{
		RootURI:      &rootURI,
		Capabilities: &protocol.ClientCapabilities{},
	}

	req := newMockRequest("initialize", params)
	handler.handleInitialize(ctx, conn, req)

	// Check that workspace was created
	if len(handler.workspaces) != 1 {
		t.Errorf("Expected 1 workspace, got %d", len(handler.workspaces))
	}

	if _, exists := handler.workspaces["/test/workspace"]; !exists {
		t.Error("Expected workspace to be created at /test/workspace")
	}
}

func TestHandler_Initialized(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	req := newMockRequest("initialized", nil)
	handler.handleInitialized(ctx, conn, req)

	if !handler.initialized {
		t.Error("Expected handler to be marked as initialized")
	}
}

func TestHandler_DidOpen(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell greet(): return \"Hello\"",
		},
	}

	req := newMockRequest("textDocument/didOpen", params)
	handler.handleDidOpen(ctx, conn, req)

	// Check that document was analyzed
	doc := handler.analyzer.GetDocument("file:///test.crl")
	if doc == nil {
		t.Error("Expected document to be analyzed and stored")
	}

	// Check that diagnostics were published
	if len(conn.notifications) == 0 {
		t.Error("Expected diagnostics notification to be sent")
	} else {
		notification := conn.notifications[0]
		if notification.method != "textDocument/publishDiagnostics" {
			t.Errorf("Expected publishDiagnostics notification, got %s", notification.method)
		}
	}
}

func TestHandler_DidOpen_NonCarrionFile(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	params := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.txt",
			LanguageID: "text",
			Version:    1,
			Text:       "This is not Carrion code",
		},
	}

	req := newMockRequest("textDocument/didOpen", params)
	handler.handleDidOpen(ctx, conn, req)

	// Check that document was not analyzed
	doc := handler.analyzer.GetDocument("file:///test.txt")
	if doc != nil {
		t.Error("Expected non-Carrion file to not be analyzed")
	}
}

func TestHandler_DidChange(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// First open the document
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell greet(): return \"Hello\"",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear notifications from open
	conn.notifications = nil

	// Now change the document
	changeParams := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			URI:     "file:///test.crl",
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Text: "spell greet(name): return \"Hello, \" + name",
			},
		},
	}

	changeReq := newMockRequest("textDocument/didChange", changeParams)
	handler.handleDidChange(ctx, conn, changeReq)

	// Check that document was re-analyzed
	doc := handler.analyzer.GetDocument("file:///test.crl")
	if doc == nil {
		t.Error("Expected document to exist after change")
	}

	// Check that diagnostics were published
	if len(conn.notifications) == 0 {
		t.Error("Expected diagnostics notification to be sent after change")
	}
}

func TestHandler_DidSave(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	text := "spell greet(): return \"Hello\""
	params := protocol.DidSaveTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///test.crl",
		},
		Text: &text,
	}

	req := newMockRequest("textDocument/didSave", params)
	handler.handleDidSave(ctx, conn, req)

	// Check that document was analyzed
	doc := handler.analyzer.GetDocument("file:///test.crl")
	if doc == nil {
		t.Error("Expected document to be analyzed on save")
	}
}

func TestHandler_DidClose(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// First open the document
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell greet(): return \"Hello\"",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Verify document exists
	if doc := handler.analyzer.GetDocument("file:///test.crl"); doc == nil {
		t.Error("Expected document to exist before close")
	}

	// Now close the document
	closeParams := protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///test.crl",
		},
	}

	closeReq := newMockRequest("textDocument/didClose", closeParams)
	handler.handleDidClose(ctx, conn, closeReq)

	// Check that document was removed
	if doc := handler.analyzer.GetDocument("file:///test.crl"); doc != nil {
		t.Error("Expected document to be removed after close")
	}
}

func TestHandler_Completion(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// First open a document
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "sp",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear previous replies
	conn.replies = nil

	// Request completion
	completionParams := protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.crl",
			},
			Position: protocol.Position{Line: 0, Character: 2},
		},
	}

	completionReq := newMockRequest("textDocument/completion", completionParams)
	handler.handleCompletion(ctx, conn, completionReq)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 completion reply, got %d", len(conn.replies))
	}

	if result, ok := conn.replies[0].(protocol.CompletionList); ok {
		if len(result.Items) == 0 {
			t.Error("Expected completion items to be returned")
		}

		// Should have "spell" keyword
		foundSpell := false
		for _, item := range result.Items {
			if item.Label == "spell" {
				foundSpell = true
				break
			}
		}
		if !foundSpell {
			t.Error("Expected 'spell' in completion items")
		}
	} else {
		t.Error("Expected reply to be CompletionList")
	}
}

func TestHandler_Hover(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// Open a document with a function
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "print(\"Hello\")",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear previous replies
	conn.replies = nil

	// Request hover over "print"
	hoverParams := protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.crl",
			},
			Position: protocol.Position{Line: 0, Character: 2},
		},
	}

	hoverReq := newMockRequest("textDocument/hover", hoverParams)
	handler.handleHover(ctx, conn, hoverReq)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 hover reply, got %d", len(conn.replies))
	}

	// Check if hover information is returned
	if hover, ok := conn.replies[0].(*protocol.Hover); ok && hover != nil {
		if hover.Contents == nil {
			t.Error("Expected hover contents to be present")
		}
	}
}

func TestHandler_Definition(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// Open a document with a function definition
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell greet(): return \"Hello\"\nresult = greet()",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear previous replies
	conn.replies = nil

	// Request definition of "greet" on second line
	definitionParams := protocol.DefinitionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.crl",
			},
			Position: protocol.Position{Line: 1, Character: 9},
		},
	}

	definitionReq := newMockRequest("textDocument/definition", definitionParams)
	handler.handleDefinition(ctx, conn, definitionReq)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 definition reply, got %d", len(conn.replies))
	}
}

func TestHandler_DocumentSymbol(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// Open a document with grimoires and spells
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text: `
grim Person:
    init(name):
        self.name = name
    
    spell greet():
        return "Hello"

spell standalone():
    return "Standalone function"
`,
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear previous replies
	conn.replies = nil

	// Request document symbols
	symbolParams := protocol.DocumentSymbolParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///test.crl",
		},
	}

	symbolReq := newMockRequest("textDocument/documentSymbol", symbolParams)
	handler.handleDocumentSymbol(ctx, conn, symbolReq)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 symbol reply, got %d", len(conn.replies))
	}

	if symbols, ok := conn.replies[0].([]protocol.DocumentSymbol); ok {
		if len(symbols) == 0 {
			t.Error("Expected document symbols to be returned")
		}

		// Should have Person grimoire and standalone function
		foundPerson := false
		foundStandalone := false

		for _, symbol := range symbols {
			if symbol.Name == "Person" && symbol.Kind == protocol.SymbolKindClass {
				foundPerson = true
			}
			if symbol.Name == "standalone" && symbol.Kind == protocol.SymbolKindFunction {
				foundStandalone = true
			}
		}

		if !foundPerson {
			t.Error("Expected Person grimoire in document symbols")
		}
		if !foundStandalone {
			t.Error("Expected standalone function in document symbols")
		}
	} else {
		t.Error("Expected reply to be []DocumentSymbol")
	}
}

func TestHandler_SemanticTokens(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// Open a document
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell test(): return 42",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.handleDidOpen(ctx, conn, openReq)

	// Clear previous replies
	conn.replies = nil

	// Request semantic tokens
	tokensParams := protocol.SemanticTokensParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///test.crl",
		},
	}

	tokensReq := newMockRequest("textDocument/semanticTokens/full", tokensParams)
	handler.handleSemanticTokens(ctx, conn, tokensReq)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 semantic tokens reply, got %d", len(conn.replies))
	}

	if tokens, ok := conn.replies[0].(*protocol.SemanticTokens); ok && tokens != nil {
		if len(tokens.Data) == 0 {
			t.Error("Expected semantic token data to be returned")
		}
	}
}

func TestHandler_InvalidMethod(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	req := newMockRequest("invalidMethod", nil)
	handler.Handle(ctx, conn, req)

	if len(conn.errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(conn.errors))
	}

	if conn.errors[0].Code != jsonrpc2.CodeMethodNotFound {
		t.Errorf("Expected method not found error, got %d", conn.errors[0].Code)
	}
}

func TestHandler_InvalidParams(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// Send initialize with invalid JSON
	req := &jsonrpc2.Request{
		Method: "initialize",
		Params: (*json.RawMessage)(&[]byte(`{"invalid": json}`)[0:15]), // Truncated invalid JSON
		ID:     jsonrpc2.ID{Num: 1},
	}

	handler.handleInitialize(ctx, conn, req)

	if len(conn.errors) != 1 {
		t.Errorf("Expected 1 error for invalid params, got %d", len(conn.errors))
	}

	if conn.errors[0].Code != jsonrpc2.CodeInvalidParams {
		t.Errorf("Expected invalid params error, got %d", conn.errors[0].Code)
	}
}

func TestHandler_Shutdown(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	req := newMockRequest("shutdown", nil)
	handler.handleShutdown(ctx, conn, req)

	if len(conn.replies) != 1 {
		t.Errorf("Expected 1 shutdown reply, got %d", len(conn.replies))
	}

	if conn.replies[0] != nil {
		t.Error("Expected shutdown reply to be null")
	}
}

func TestHandler_Exit(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	req := newMockRequest("exit", nil)
	handler.handleExit(ctx, conn, req)

	// Exit should not send any replies
	if len(conn.replies) != 0 {
		t.Errorf("Expected 0 replies for exit, got %d", len(conn.replies))
	}
}

func TestHandler_DidChangeConfiguration(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	req := newMockRequest("workspace/didChangeConfiguration", nil)
	handler.handleDidChangeConfiguration(ctx, conn, req)

	// Should handle without error
	if len(conn.errors) != 0 {
		t.Errorf("Expected 0 errors for configuration change, got %d", len(conn.errors))
	}
}

// Integration test for full LSP workflow
func TestHandler_FullWorkflow(t *testing.T) {
	handler := NewHandler()
	conn := &mockConn{}
	ctx := context.Background()

	// 1. Initialize
	initParams := protocol.InitializeParams{
		Capabilities: &protocol.ClientCapabilities{},
	}
	initReq := newMockRequest("initialize", initParams)
	handler.Handle(ctx, conn, initReq)

	// 2. Initialized
	initializedReq := newMockRequest("initialized", nil)
	handler.Handle(ctx, conn, initializedReq)

	// 3. Open document
	openParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.crl",
			LanguageID: "carrion",
			Version:    1,
			Text:       "spell greet(name): return \"Hello, \" + name",
		},
	}
	openReq := newMockRequest("textDocument/didOpen", openParams)
	handler.Handle(ctx, conn, openReq)

	// Clear replies/notifications
	conn.replies = nil
	conn.notifications = nil

	// 4. Request completion
	completionParams := protocol.CompletionParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.crl",
			},
			Position: protocol.Position{Line: 0, Character: 30},
		},
	}
	completionReq := newMockRequest("textDocument/completion", completionParams)
	handler.Handle(ctx, conn, completionReq)

	// 5. Request hover
	hoverParams := protocol.HoverParams{
		TextDocumentPositionParams: protocol.TextDocumentPositionParams{
			TextDocument: protocol.TextDocumentIdentifier{
				URI: "file:///test.crl",
			},
			Position: protocol.Position{Line: 0, Character: 6},
		},
	}
	hoverReq := newMockRequest("textDocument/hover", hoverParams)
	handler.Handle(ctx, conn, hoverReq)

	// 6. Shutdown
	shutdownReq := newMockRequest("shutdown", nil)
	handler.Handle(ctx, conn, shutdownReq)

	// 7. Exit
	exitReq := newMockRequest("exit", nil)
	handler.Handle(ctx, conn, exitReq)

	// Verify we got appropriate responses
	if len(conn.replies) < 3 { // completion, hover, shutdown
		t.Errorf("Expected at least 3 replies in full workflow, got %d", len(conn.replies))
	}

	if len(conn.errors) > 0 {
		t.Errorf("Expected no errors in full workflow, got %d", len(conn.errors))
	}
}
