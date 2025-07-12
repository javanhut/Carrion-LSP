package protocol

// LSP Protocol types specific to Carrion Language Server

// Initialize request structures
type InitializeParams struct {
	ProcessID             *int                    `json:"processId"`
	ClientInfo            *ClientInfo             `json:"clientInfo,omitempty"`
	Locale                *string                 `json:"locale,omitempty"`
	RootPath              *string                 `json:"rootPath,omitempty"`
	RootURI               *string                 `json:"rootUri"`
	Capabilities          *ClientCapabilities     `json:"capabilities"`
	InitializationOptions interface{}             `json:"initializationOptions,omitempty"`
	WorkspaceFolders      []WorkspaceFolder       `json:"workspaceFolders,omitempty"`
}

type ClientInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Window       *WindowClientCapabilities       `json:"window,omitempty"`
	General      *GeneralClientCapabilities      `json:"general,omitempty"`
}

type WorkspaceClientCapabilities struct {
	ApplyEdit              *bool                              `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditClientCapabilities   `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DidChangeConfigurationCapabilities `json:"didChangeConfiguration,omitempty"`
	DidChangeWatchedFiles  *DidChangeWatchedFilesCapabilities `json:"didChangeWatchedFiles,omitempty"`
	Symbol                 *WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`
	ExecuteCommand         *ExecuteCommandClientCapabilities  `json:"executeCommand,omitempty"`
	Configuration          *bool                              `json:"configuration,omitempty"`
	WorkspaceFolders       *bool                              `json:"workspaceFolders,omitempty"`
}

type TextDocumentClientCapabilities struct {
	Synchronization    *TextDocumentSyncClientCapabilities    `json:"synchronization,omitempty"`
	Completion         *CompletionClientCapabilities          `json:"completion,omitempty"`
	Hover              *HoverClientCapabilities               `json:"hover,omitempty"`
	SignatureHelp      *SignatureHelpClientCapabilities       `json:"signatureHelp,omitempty"`
	Declaration        *DeclarationClientCapabilities         `json:"declaration,omitempty"`
	Definition         *DefinitionClientCapabilities          `json:"definition,omitempty"`
	TypeDefinition     *TypeDefinitionClientCapabilities      `json:"typeDefinition,omitempty"`
	Implementation     *ImplementationClientCapabilities      `json:"implementation,omitempty"`
	References         *ReferenceClientCapabilities           `json:"references,omitempty"`
	DocumentHighlight  *DocumentHighlightClientCapabilities   `json:"documentHighlight,omitempty"`
	DocumentSymbol     *DocumentSymbolClientCapabilities      `json:"documentSymbol,omitempty"`
	CodeAction         *CodeActionClientCapabilities          `json:"codeAction,omitempty"`
	CodeLens           *CodeLensClientCapabilities            `json:"codeLens,omitempty"`
	DocumentLink       *DocumentLinkClientCapabilities        `json:"documentLink,omitempty"`
	ColorProvider      *DocumentColorClientCapabilities       `json:"colorProvider,omitempty"`
	Formatting         *DocumentFormattingClientCapabilities  `json:"formatting,omitempty"`
	RangeFormatting    *DocumentRangeFormattingClientCapabilities `json:"rangeFormatting,omitempty"`
	OnTypeFormatting   *DocumentOnTypeFormattingClientCapabilities `json:"onTypeFormatting,omitempty"`
	Rename             *RenameClientCapabilities              `json:"rename,omitempty"`
	PublishDiagnostics *PublishDiagnosticsClientCapabilities  `json:"publishDiagnostics,omitempty"`
	FoldingRange       *FoldingRangeClientCapabilities        `json:"foldingRange,omitempty"`
	SelectionRange     *SelectionRangeClientCapabilities      `json:"selectionRange,omitempty"`
	SemanticTokens     *SemanticTokensClientCapabilities      `json:"semanticTokens,omitempty"`
}

type WindowClientCapabilities struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
	ShowMessage      *ShowMessageRequestClientCapabilities `json:"showMessage,omitempty"`
	ShowDocument     *ShowDocumentClientCapabilities       `json:"showDocument,omitempty"`
}

type GeneralClientCapabilities struct {
	RegularExpressions *RegularExpressionsClientCapabilities `json:"regularExpressions,omitempty"`
	Markdown           *MarkdownClientCapabilities           `json:"markdown,omitempty"`
}

// Initialize response structures
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string  `json:"name"`
	Version *string `json:"version,omitempty"`
}

type ServerCapabilities struct {
	TextDocumentSync           interface{}                     `json:"textDocumentSync,omitempty"`
	CompletionProvider         *CompletionOptions              `json:"completionProvider,omitempty"`
	HoverProvider              interface{}                     `json:"hoverProvider,omitempty"`
	SignatureHelpProvider      *SignatureHelpOptions           `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider        interface{}                     `json:"declarationProvider,omitempty"`
	DefinitionProvider         interface{}                     `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider     interface{}                     `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider     interface{}                     `json:"implementationProvider,omitempty"`
	ReferencesProvider         interface{}                     `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider  interface{}                     `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider     interface{}                     `json:"documentSymbolProvider,omitempty"`
	CodeActionProvider         interface{}                     `json:"codeActionProvider,omitempty"`
	CodeLensProvider           *CodeLensOptions               `json:"codeLensProvider,omitempty"`
	DocumentLinkProvider       *DocumentLinkOptions           `json:"documentLinkProvider,omitempty"`
	ColorProvider              interface{}                     `json:"colorProvider,omitempty"`
	DocumentFormattingProvider interface{}                     `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider interface{}               `json:"documentRangeFormattingProvider,omitempty"`
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider             interface{}                     `json:"renameProvider,omitempty"`
	FoldingRangeProvider       interface{}                     `json:"foldingRangeProvider,omitempty"`
	ExecuteCommandProvider     *ExecuteCommandOptions         `json:"executeCommandProvider,omitempty"`
	SelectionRangeProvider     interface{}                     `json:"selectionRangeProvider,omitempty"`
	LinkedEditingRangeProvider interface{}                     `json:"linkedEditingRangeProvider,omitempty"`
	CallHierarchyProvider      interface{}                     `json:"callHierarchyProvider,omitempty"`
	SemanticTokensProvider     *SemanticTokensOptions         `json:"semanticTokensProvider,omitempty"`
	MonikerProvider            interface{}                     `json:"monikerProvider,omitempty"`
	WorkspaceSymbolProvider    interface{}                     `json:"workspaceSymbolProvider,omitempty"`
	Workspace                  *WorkspaceServerCapabilities   `json:"workspace,omitempty"`
	Experimental               interface{}                     `json:"experimental,omitempty"`
}

// Text Document Sync
type TextDocumentSyncOptions struct {
	OpenClose bool                        `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind        `json:"change,omitempty"`
	WillSave  bool                        `json:"willSave,omitempty"`
	WillSaveWaitUntil bool                `json:"willSaveWaitUntil,omitempty"`
	Save      *SaveOptions                `json:"save,omitempty"`
}

type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone        TextDocumentSyncKind = 0
	TextDocumentSyncKindFull        TextDocumentSyncKind = 1
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

// Completion
type CompletionOptions struct {
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	ResolveProvider     bool     `json:"resolveProvider,omitempty"`
	CompletionItem      *CompletionOptionsCompletionItem `json:"completionItem,omitempty"`
}

type CompletionOptionsCompletionItem struct {
	LabelDetailsSupport bool `json:"labelDetailsSupport,omitempty"`
}

// Semantic Tokens
type SemanticTokensOptions struct {
	Legend SemanticTokensLegend `json:"legend"`
	Range  interface{}          `json:"range,omitempty"`
	Full   interface{}          `json:"full,omitempty"`
}

type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

// Position and Range
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// Text Document notifications
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         *string                `json:"text,omitempty"`
}

type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// Diagnostics
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
	Range              Range                    `json:"range"`
	Severity           DiagnosticSeverity       `json:"severity,omitempty"`
	Code               interface{}              `json:"code,omitempty"`
	CodeDescription    *CodeDescription         `json:"codeDescription,omitempty"`
	Source             string                   `json:"source,omitempty"`
	Message            string                   `json:"message"`
	Tags               []DiagnosticTag          `json:"tags,omitempty"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	Data               interface{}              `json:"data,omitempty"`
}

type DiagnosticSeverity int

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

type DiagnosticTag int

type CodeDescription struct {
	Href string `json:"href"`
}

type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// Completion
type CompletionParams struct {
	TextDocumentPositionParams
	Context *CompletionContext `json:"context,omitempty"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter *string               `json:"triggerCharacter,omitempty"`
}

type CompletionTriggerKind int

const (
	CompletionTriggerKindInvoked                CompletionTriggerKind = 1
	CompletionTriggerKindTriggerCharacter       CompletionTriggerKind = 2
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
	Label               string                 `json:"label"`
	LabelDetails        *CompletionItemLabelDetails `json:"labelDetails,omitempty"`
	Kind                CompletionItemKind     `json:"kind,omitempty"`
	Tags                []CompletionItemTag    `json:"tags,omitempty"`
	Detail              string                 `json:"detail,omitempty"`
	Documentation       interface{}            `json:"documentation,omitempty"`
	Deprecated          bool                   `json:"deprecated,omitempty"`
	Preselect           bool                   `json:"preselect,omitempty"`
	SortText            string                 `json:"sortText,omitempty"`
	FilterText          string                 `json:"filterText,omitempty"`
	InsertText          string                 `json:"insertText,omitempty"`
	InsertTextFormat    InsertTextFormat       `json:"insertTextFormat,omitempty"`
	InsertTextMode      InsertTextMode         `json:"insertTextMode,omitempty"`
	TextEdit            *TextEdit              `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit             `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string               `json:"commitCharacters,omitempty"`
	Command             *Command               `json:"command,omitempty"`
	Data                interface{}            `json:"data,omitempty"`
}

type CompletionItemLabelDetails struct {
	Detail      string `json:"detail,omitempty"`
	Description string `json:"description,omitempty"`
}

type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

type CompletionItemTag int

type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet   InsertTextFormat = 2
)

type InsertTextMode int

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments,omitempty"`
}

// Hover
type HoverParams struct {
	TextDocumentPositionParams
}

type Hover struct {
	Contents interface{} `json:"contents"`
	Range    *Range      `json:"range,omitempty"`
}

// Definition
type DefinitionParams struct {
	TextDocumentPositionParams
}

// References
type ReferenceParams struct {
	TextDocumentPositionParams
	Context ReferenceContext `json:"context"`
}

type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

// Document Symbol
type DocumentSymbolParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type DocumentSymbol struct {
	Name           string             `json:"name"`
	Detail         string             `json:"detail,omitempty"`
	Kind           SymbolKind         `json:"kind"`
	Tags           []SymbolTag        `json:"tags,omitempty"`
	Deprecated     bool               `json:"deprecated,omitempty"`
	Range          Range              `json:"range"`
	SelectionRange Range              `json:"selectionRange"`
	Children       []DocumentSymbol   `json:"children,omitempty"`
}

type SymbolKind int

const (
	SymbolKindFile          SymbolKind = 1
	SymbolKindModule        SymbolKind = 2
	SymbolKindNamespace     SymbolKind = 3
	SymbolKindPackage       SymbolKind = 4
	SymbolKindClass         SymbolKind = 5
	SymbolKindMethod        SymbolKind = 6
	SymbolKindProperty      SymbolKind = 7
	SymbolKindField         SymbolKind = 8
	SymbolKindConstructor   SymbolKind = 9
	SymbolKindEnum          SymbolKind = 10
	SymbolKindInterface     SymbolKind = 11
	SymbolKindFunction      SymbolKind = 12
	SymbolKindVariable      SymbolKind = 13
	SymbolKindConstant      SymbolKind = 14
	SymbolKindString        SymbolKind = 15
	SymbolKindNumber        SymbolKind = 16
	SymbolKindBoolean       SymbolKind = 17
	SymbolKindArray         SymbolKind = 18
	SymbolKindObject        SymbolKind = 19
	SymbolKindKey           SymbolKind = 20
	SymbolKindNull          SymbolKind = 21
	SymbolKindEnumMember    SymbolKind = 22
	SymbolKindStruct        SymbolKind = 23
	SymbolKindEvent         SymbolKind = 24
	SymbolKindOperator      SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

type SymbolTag int

// Semantic Tokens
type SemanticTokensParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

type SemanticTokens struct {
	ResultID string `json:"resultId,omitempty"`
	Data     []int  `json:"data"`
}

// Formatting
type DocumentFormattingParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

type FormattingOptions struct {
	TabSize                int                    `json:"tabSize"`
	InsertSpaces           bool                   `json:"insertSpaces"`
	TrimTrailingWhitespace bool                   `json:"trimTrailingWhitespace,omitempty"`
	InsertFinalNewline     bool                   `json:"insertFinalNewline,omitempty"`
	TrimFinalNewlines      bool                   `json:"trimFinalNewlines,omitempty"`
	AdditionalProperties   map[string]interface{} `json:",inline"`
}

// Placeholder types for unimplemented capabilities
type WorkspaceEditClientCapabilities struct{}
type DidChangeConfigurationCapabilities struct{}
type DidChangeWatchedFilesCapabilities struct{}
type WorkspaceSymbolClientCapabilities struct{}
type ExecuteCommandClientCapabilities struct{}
type TextDocumentSyncClientCapabilities struct{}
type CompletionClientCapabilities struct{}
type HoverClientCapabilities struct{}
type SignatureHelpClientCapabilities struct{}
type DeclarationClientCapabilities struct{}
type DefinitionClientCapabilities struct{}
type TypeDefinitionClientCapabilities struct{}
type ImplementationClientCapabilities struct{}
type ReferenceClientCapabilities struct{}
type DocumentHighlightClientCapabilities struct{}
type DocumentSymbolClientCapabilities struct{}
type CodeActionClientCapabilities struct{}
type CodeLensClientCapabilities struct{}
type DocumentLinkClientCapabilities struct{}
type DocumentColorClientCapabilities struct{}
type DocumentFormattingClientCapabilities struct{}
type DocumentRangeFormattingClientCapabilities struct{}
type DocumentOnTypeFormattingClientCapabilities struct{}
type RenameClientCapabilities struct{}
type PublishDiagnosticsClientCapabilities struct{}
type FoldingRangeClientCapabilities struct{}
type SelectionRangeClientCapabilities struct{}
type SemanticTokensClientCapabilities struct{}
type ShowMessageRequestClientCapabilities struct{}
type ShowDocumentClientCapabilities struct{}
type RegularExpressionsClientCapabilities struct{}
type MarkdownClientCapabilities struct{}
type SignatureHelpOptions struct{}
type CodeLensOptions struct{}
type DocumentLinkOptions struct{}
type DocumentOnTypeFormattingOptions struct{}
type ExecuteCommandOptions struct{}
type WorkspaceServerCapabilities struct{}