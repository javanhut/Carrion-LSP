package protocol

import (
	"encoding/json"
	"testing"
)

func TestPosition(t *testing.T) {
	pos := Position{Line: 10, Character: 5}

	if pos.Line != 10 {
		t.Errorf("Expected line 10, got %d", pos.Line)
	}

	if pos.Character != 5 {
		t.Errorf("Expected character 5, got %d", pos.Character)
	}
}

func TestRange(t *testing.T) {
	start := Position{Line: 0, Character: 0}
	end := Position{Line: 1, Character: 10}

	r := Range{Start: start, End: end}

	if r.Start.Line != 0 || r.Start.Character != 0 {
		t.Error("Expected start position to be (0, 0)")
	}

	if r.End.Line != 1 || r.End.Character != 10 {
		t.Error("Expected end position to be (1, 10)")
	}
}

func TestLocation(t *testing.T) {
	location := Location{
		URI: "file:///test.crl",
		Range: Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 0, Character: 5},
		},
	}

	if location.URI != "file:///test.crl" {
		t.Errorf("Expected URI to be 'file:///test.crl', got %s", location.URI)
	}
}

func TestTextDocumentItem(t *testing.T) {
	doc := TextDocumentItem{
		URI:        "file:///test.crl",
		LanguageID: "carrion",
		Version:    1,
		Text:       "spell test(): return 42",
	}

	if doc.URI != "file:///test.crl" {
		t.Errorf("Expected URI to be 'file:///test.crl', got %s", doc.URI)
	}

	if doc.LanguageID != "carrion" {
		t.Errorf("Expected language ID to be 'carrion', got %s", doc.LanguageID)
	}

	if doc.Version != 1 {
		t.Errorf("Expected version to be 1, got %d", doc.Version)
	}
}

func TestCompletionItem(t *testing.T) {
	item := CompletionItem{
		Label:            "spell",
		Kind:             CompletionItemKindKeyword,
		Detail:           "spell keyword",
		InsertText:       "spell ${1:name}(${2:params}):\\n\\t${3:body}",
		InsertTextFormat: InsertTextFormatSnippet,
	}

	if item.Label != "spell" {
		t.Errorf("Expected label to be 'spell', got %s", item.Label)
	}

	if item.Kind != CompletionItemKindKeyword {
		t.Errorf("Expected kind to be keyword, got %d", item.Kind)
	}

	if item.InsertTextFormat != InsertTextFormatSnippet {
		t.Errorf("Expected insert text format to be snippet, got %d", item.InsertTextFormat)
	}
}

func TestCompletionItemKindConstants(t *testing.T) {
	tests := []struct {
		kind     CompletionItemKind
		expected int
	}{
		{CompletionItemKindText, 1},
		{CompletionItemKindMethod, 2},
		{CompletionItemKindFunction, 3},
		{CompletionItemKindConstructor, 4},
		{CompletionItemKindField, 5},
		{CompletionItemKindVariable, 6},
		{CompletionItemKindClass, 7},
		{CompletionItemKindKeyword, 14},
	}

	for _, test := range tests {
		if int(test.kind) != test.expected {
			t.Errorf("Expected %d, got %d for completion kind", test.expected, int(test.kind))
		}
	}
}

func TestSymbolKindConstants(t *testing.T) {
	tests := []struct {
		kind     SymbolKind
		expected int
	}{
		{SymbolKindFile, 1},
		{SymbolKindModule, 2},
		{SymbolKindNamespace, 3},
		{SymbolKindPackage, 4},
		{SymbolKindClass, 5},
		{SymbolKindMethod, 6},
		{SymbolKindFunction, 12},
		{SymbolKindVariable, 13},
	}

	for _, test := range tests {
		if int(test.kind) != test.expected {
			t.Errorf("Expected %d, got %d for symbol kind", test.expected, int(test.kind))
		}
	}
}

func TestDiagnosticSeverityConstants(t *testing.T) {
	tests := []struct {
		severity DiagnosticSeverity
		expected int
	}{
		{DiagnosticSeverityError, 1},
		{DiagnosticSeverityWarning, 2},
		{DiagnosticSeverityInformation, 3},
		{DiagnosticSeverityHint, 4},
	}

	for _, test := range tests {
		if int(test.severity) != test.expected {
			t.Errorf("Expected %d, got %d for diagnostic severity", test.expected, int(test.severity))
		}
	}
}

func TestTextDocumentSyncKindConstants(t *testing.T) {
	tests := []struct {
		kind     TextDocumentSyncKind
		expected int
	}{
		{TextDocumentSyncKindNone, 0},
		{TextDocumentSyncKindFull, 1},
		{TextDocumentSyncKindIncremental, 2},
	}

	for _, test := range tests {
		if int(test.kind) != test.expected {
			t.Errorf("Expected %d, got %d for text document sync kind", test.expected, int(test.kind))
		}
	}
}

func TestDiagnostic(t *testing.T) {
	diagnostic := Diagnostic{
		Range: Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 0, Character: 5},
		},
		Severity: DiagnosticSeverityError,
		Message:  "Syntax error",
		Source:   "carrion-lsp",
	}

	if diagnostic.Severity != DiagnosticSeverityError {
		t.Errorf("Expected error severity, got %d", diagnostic.Severity)
	}

	if diagnostic.Message != "Syntax error" {
		t.Errorf("Expected message 'Syntax error', got %s", diagnostic.Message)
	}

	if diagnostic.Source != "carrion-lsp" {
		t.Errorf("Expected source 'carrion-lsp', got %s", diagnostic.Source)
	}
}

func TestDocumentSymbol(t *testing.T) {
	symbol := DocumentSymbol{
		Name: "Person",
		Kind: SymbolKindClass,
		Range: Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 10, Character: 0},
		},
		SelectionRange: Range{
			Start: Position{Line: 0, Character: 5},
			End:   Position{Line: 0, Character: 11},
		},
		Children: []DocumentSymbol{
			{
				Name: "init",
				Kind: SymbolKindMethod,
				Range: Range{
					Start: Position{Line: 1, Character: 4},
					End:   Position{Line: 3, Character: 0},
				},
				SelectionRange: Range{
					Start: Position{Line: 1, Character: 4},
					End:   Position{Line: 1, Character: 8},
				},
			},
		},
	}

	if symbol.Name != "Person" {
		t.Errorf("Expected name 'Person', got %s", symbol.Name)
	}

	if symbol.Kind != SymbolKindClass {
		t.Errorf("Expected class kind, got %d", symbol.Kind)
	}

	if len(symbol.Children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(symbol.Children))
	}

	if symbol.Children[0].Name != "init" {
		t.Errorf("Expected child name 'init', got %s", symbol.Children[0].Name)
	}
}

func TestSemanticTokens(t *testing.T) {
	tokens := SemanticTokens{
		ResultID: "test-result",
		Data:     []int{0, 0, 5, 0, 0, 0, 6, 4, 1, 0},
	}

	if tokens.ResultID != "test-result" {
		t.Errorf("Expected result ID 'test-result', got %s", tokens.ResultID)
	}

	if len(tokens.Data) != 10 {
		t.Errorf("Expected 10 data items, got %d", len(tokens.Data))
	}

	// Data should be in groups of 5
	if len(tokens.Data)%5 != 0 {
		t.Error("Expected semantic token data to be in groups of 5")
	}
}

func TestFormattingOptions(t *testing.T) {
	options := FormattingOptions{
		TabSize:                4,
		InsertSpaces:           true,
		TrimTrailingWhitespace: true,
		InsertFinalNewline:     true,
		TrimFinalNewlines:      false,
	}

	if options.TabSize != 4 {
		t.Errorf("Expected tab size 4, got %d", options.TabSize)
	}

	if !options.InsertSpaces {
		t.Error("Expected insert spaces to be true")
	}

	if !options.TrimTrailingWhitespace {
		t.Error("Expected trim trailing whitespace to be true")
	}
}

func TestTextEdit(t *testing.T) {
	edit := TextEdit{
		Range: Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 0, Character: 5},
		},
		NewText: "spell",
	}

	if edit.NewText != "spell" {
		t.Errorf("Expected new text 'spell', got %s", edit.NewText)
	}
}

func TestHover(t *testing.T) {
	hover := Hover{
		Contents: "**spell**: Function definition keyword",
		Range: &Range{
			Start: Position{Line: 0, Character: 0},
			End:   Position{Line: 0, Character: 5},
		},
	}

	if hover.Contents != "**spell**: Function definition keyword" {
		t.Errorf("Expected specific contents, got %v", hover.Contents)
	}

	if hover.Range == nil {
		t.Error("Expected range to be set")
	}
}

// Test JSON serialization/deserialization
func TestInitializeParamsSerialization(t *testing.T) {
	params := InitializeParams{
		ProcessID: func() *int { i := 1234; return &i }(),
		ClientInfo: &ClientInfo{
			Name:    "Test Client",
			Version: func() *string { s := "1.0.0"; return &s }(),
		},
		RootURI: func() *string { s := "file:///test"; return &s }(),
		Capabilities: &ClientCapabilities{
			TextDocument: &TextDocumentClientCapabilities{},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeParams: %v", err)
	}

	var unmarshaled InitializeParams
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal InitializeParams: %v", err)
	}

	if *unmarshaled.ProcessID != 1234 {
		t.Errorf("Expected process ID 1234, got %d", *unmarshaled.ProcessID)
	}

	if unmarshaled.ClientInfo.Name != "Test Client" {
		t.Errorf("Expected client name 'Test Client', got %s", unmarshaled.ClientInfo.Name)
	}
}

func TestInitializeResultSerialization(t *testing.T) {
	result := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: &TextDocumentSyncOptions{
				OpenClose: true,
				Change:    TextDocumentSyncKindIncremental,
			},
			CompletionProvider: &CompletionOptions{
				TriggerCharacters: []string{".", "("},
				ResolveProvider:   false,
			},
			HoverProvider:      true,
			DefinitionProvider: true,
		},
		ServerInfo: &ServerInfo{
			Name:    "Carrion Language Server",
			Version: func() *string { s := "0.1.0"; return &s }(),
		},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal InitializeResult: %v", err)
	}

	var unmarshaled InitializeResult
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal InitializeResult: %v", err)
	}

	if unmarshaled.ServerInfo.Name != "Carrion Language Server" {
		t.Errorf("Expected server name 'Carrion Language Server', got %s", unmarshaled.ServerInfo.Name)
	}
}

func TestCompletionListSerialization(t *testing.T) {
	list := CompletionList{
		IsIncomplete: false,
		Items: []CompletionItem{
			{
				Label:            "spell",
				Kind:             CompletionItemKindKeyword,
				Detail:           "Function definition keyword",
				InsertText:       "spell ${1:name}(${2:params}):\\n\\t${3:body}",
				InsertTextFormat: InsertTextFormatSnippet,
			},
			{
				Label:  "print",
				Kind:   CompletionItemKindFunction,
				Detail: "Print function",
			},
		},
	}

	data, err := json.Marshal(list)
	if err != nil {
		t.Fatalf("Failed to marshal CompletionList: %v", err)
	}

	var unmarshaled CompletionList
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal CompletionList: %v", err)
	}

	if len(unmarshaled.Items) != 2 {
		t.Errorf("Expected 2 completion items, got %d", len(unmarshaled.Items))
	}

	if unmarshaled.Items[0].Label != "spell" {
		t.Errorf("Expected first item label 'spell', got %s", unmarshaled.Items[0].Label)
	}
}

func TestPublishDiagnosticsParamsSerialization(t *testing.T) {
	params := PublishDiagnosticsParams{
		URI: "file:///test.crl",
		Diagnostics: []Diagnostic{
			{
				Range: Range{
					Start: Position{Line: 0, Character: 0},
					End:   Position{Line: 0, Character: 5},
				},
				Severity: DiagnosticSeverityError,
				Message:  "Syntax error",
				Source:   "carrion-lsp",
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal PublishDiagnosticsParams: %v", err)
	}

	var unmarshaled PublishDiagnosticsParams
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal PublishDiagnosticsParams: %v", err)
	}

	if unmarshaled.URI != "file:///test.crl" {
		t.Errorf("Expected URI 'file:///test.crl', got %s", unmarshaled.URI)
	}

	if len(unmarshaled.Diagnostics) != 1 {
		t.Errorf("Expected 1 diagnostic, got %d", len(unmarshaled.Diagnostics))
	}
}

func TestDocumentSymbolSerialization(t *testing.T) {
	symbols := []DocumentSymbol{
		{
			Name: "Person",
			Kind: SymbolKindClass,
			Range: Range{
				Start: Position{Line: 0, Character: 0},
				End:   Position{Line: 10, Character: 0},
			},
			SelectionRange: Range{
				Start: Position{Line: 0, Character: 5},
				End:   Position{Line: 0, Character: 11},
			},
		},
	}

	data, err := json.Marshal(symbols)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentSymbol array: %v", err)
	}

	var unmarshaled []DocumentSymbol
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal DocumentSymbol array: %v", err)
	}

	if len(unmarshaled) != 1 {
		t.Errorf("Expected 1 symbol, got %d", len(unmarshaled))
	}

	if unmarshaled[0].Name != "Person" {
		t.Errorf("Expected symbol name 'Person', got %s", unmarshaled[0].Name)
	}
}

// Test for optional pointer fields
func TestOptionalFields(t *testing.T) {
	// Test that optional fields can be nil
	params := InitializeParams{
		Capabilities: &ClientCapabilities{},
	}

	if params.ProcessID != nil {
		t.Error("Expected ProcessID to be nil")
	}

	if params.RootURI != nil {
		t.Error("Expected RootURI to be nil")
	}

	if params.ClientInfo != nil {
		t.Error("Expected ClientInfo to be nil")
	}

	// Test JSON marshaling with nil fields
	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal params with nil fields: %v", err)
	}

	var unmarshaled InitializeParams
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal params with nil fields: %v", err)
	}
}
