package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
	"github.com/sourcegraph/jsonrpc2"
)

type Handler struct {
	analyzer    *analyzer.Analyzer
	initialized bool
	clientCaps  *protocol.ClientCapabilities
	workspaces  map[string]*analyzer.Workspace
}

func NewHandler() *Handler {
	return &Handler{
		analyzer:   analyzer.New(),
		workspaces: make(map[string]*analyzer.Workspace),
	}
}

func (h *Handler) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	switch req.Method {
	case "initialize":
		h.handleInitialize(ctx, conn, req)
	case "initialized":
		h.handleInitialized(ctx, conn, req)
	case "textDocument/didOpen":
		h.handleDidOpen(ctx, conn, req)
	case "textDocument/didChange":
		h.handleDidChange(ctx, conn, req)
	case "textDocument/didSave":
		h.handleDidSave(ctx, conn, req)
	case "textDocument/didClose":
		h.handleDidClose(ctx, conn, req)
	case "textDocument/completion":
		h.handleCompletion(ctx, conn, req)
	case "textDocument/hover":
		h.handleHover(ctx, conn, req)
	case "textDocument/definition":
		h.handleDefinition(ctx, conn, req)
	case "textDocument/references":
		h.handleReferences(ctx, conn, req)
	case "textDocument/documentSymbol":
		h.handleDocumentSymbol(ctx, conn, req)
	case "textDocument/semanticTokens/full":
		h.handleSemanticTokens(ctx, conn, req)
	case "textDocument/formatting":
		h.handleFormatting(ctx, conn, req)
	case "workspace/didChangeConfiguration":
		h.handleDidChangeConfiguration(ctx, conn, req)
	case "shutdown":
		h.handleShutdown(ctx, conn, req)
	case "exit":
		h.handleExit(ctx, conn, req)
	default:
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		})
	}
}

func (h *Handler) handleInitialize(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	h.clientCaps = params.Capabilities

	// Initialize workspace if provided
	if params.RootURI != nil {
		workspacePath := strings.TrimPrefix(*params.RootURI, "file://")
		h.workspaces[workspacePath] = analyzer.NewWorkspace(workspacePath)
	}

	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindIncremental,
				Save: &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", "(", " "},
				ResolveProvider:   false,
			},
			HoverProvider:          true,
			DefinitionProvider:     true,
			ReferencesProvider:     true,
			DocumentSymbolProvider: true,
			SemanticTokensProvider: &protocol.SemanticTokensOptions{
				Legend: protocol.SemanticTokensLegend{
					TokenTypes: []string{
						"keyword", "string", "number", "operator", "variable",
						"function", "class", "parameter", "property", "comment",
					},
					TokenModifiers: []string{"definition", "readonly", "static", "deprecated"},
				},
				Full: true,
			},
			DocumentFormattingProvider: true,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "Carrion Language Server",
			Version: func() *string { s := "0.1.0"; return &s }(),
		},
	}

	conn.Reply(ctx, req.ID, result)
}

func (h *Handler) handleInitialized(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	h.initialized = true
	log.Println("Carrion LSP server initialized")
}

func (h *Handler) handleDidOpen(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		log.Printf("Error unmarshaling didOpen params: %v", err)
		return
	}

	// Only handle .crl files
	if !strings.HasSuffix(params.TextDocument.URI, ".crl") {
		return
	}

	// Parse the document and update analysis
	h.analyzeDocument(ctx, conn, params.TextDocument.URI, params.TextDocument.Text)
}

func (h *Handler) handleDidChange(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		log.Printf("Error unmarshaling didChange params: %v", err)
		return
	}

	// Apply incremental changes
	for _, change := range params.ContentChanges {
		if change.Range == nil {
			// Full document change
			h.analyzeDocument(ctx, conn, params.TextDocument.URI, change.Text)
		} else {
			// Incremental change - for now, treat as full document update
			// TODO: Implement proper incremental updates
			h.analyzeDocument(ctx, conn, params.TextDocument.URI, change.Text)
		}
	}
}

func (h *Handler) handleDidSave(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		log.Printf("Error unmarshaling didSave params: %v", err)
		return
	}

	// Re-analyze the document on save
	if params.Text != nil {
		h.analyzeDocument(ctx, conn, params.TextDocument.URI, *params.Text)
	}
}

func (h *Handler) handleDidClose(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		log.Printf("Error unmarshaling didClose params: %v", err)
		return
	}

	// Remove document from analysis
	h.analyzer.RemoveDocument(params.TextDocument.URI)
}

func (h *Handler) analyzeDocument(ctx context.Context, conn *jsonrpc2.Conn, uri, content string) {
	// Dynamic parsing using TheCarrionLanguage parser
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	errors := p.Errors()
	var diagnostics []protocol.Diagnostic

	for _, err := range errors {
		diagnostics = append(diagnostics, protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{Line: 0, Character: 0}, // TODO: Extract actual position
				End:   protocol.Position{Line: 0, Character: 0},
			},
			Severity: protocol.DiagnosticSeverityError,
			Message:  err,
			Source:   "carrion-lsp",
		})
	}

	// Update analyzer with parsed AST
	h.analyzer.UpdateDocument(uri, content, program)

	// Send diagnostics to client
	conn.Notify(ctx, "textDocument/publishDiagnostics", protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	})
}

func (h *Handler) handleCompletion(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.CompletionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	completions := h.analyzer.GetCompletions(params.TextDocument.URI, params.Position)

	result := protocol.CompletionList{
		IsIncomplete: false,
		Items:        completions,
	}

	conn.Reply(ctx, req.ID, result)
}

func (h *Handler) handleHover(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	hover := h.analyzer.GetHover(params.TextDocument.URI, params.Position)
	conn.Reply(ctx, req.ID, hover)
}

func (h *Handler) handleDefinition(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DefinitionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	locations := h.analyzer.GetDefinition(params.TextDocument.URI, params.Position)
	conn.Reply(ctx, req.ID, locations)
}

func (h *Handler) handleReferences(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.ReferenceParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	locations := h.analyzer.GetReferences(params.TextDocument.URI, params.Position, params.Context.IncludeDeclaration)
	conn.Reply(ctx, req.ID, locations)
}

func (h *Handler) handleDocumentSymbol(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DocumentSymbolParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	symbols := h.analyzer.GetDocumentSymbols(params.TextDocument.URI)
	conn.Reply(ctx, req.ID, symbols)
}

func (h *Handler) handleSemanticTokens(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.SemanticTokensParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	tokens := h.analyzer.GetSemanticTokens(params.TextDocument.URI)
	conn.Reply(ctx, req.ID, tokens)
}

func (h *Handler) handleFormatting(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	var params protocol.DocumentFormattingParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		conn.ReplyWithError(ctx, req.ID, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeInvalidParams,
			Message: err.Error(),
		})
		return
	}

	edits := h.analyzer.FormatDocument(params.TextDocument.URI, params.Options)
	conn.Reply(ctx, req.ID, edits)
}

func (h *Handler) handleDidChangeConfiguration(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	// Handle configuration changes - could reload parser if needed
	log.Println("Configuration changed - could trigger parser reload")
}

func (h *Handler) handleShutdown(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	conn.Reply(ctx, req.ID, nil)
}

func (h *Handler) handleExit(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) {
	log.Println("LSP server exiting")
	// Graceful shutdown
}
