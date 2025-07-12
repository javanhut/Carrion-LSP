package analyzer

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/javanhut/CarrionLSP/internal/protocol"
	"github.com/javanhut/TheCarrionLanguage/src/ast"
	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
	"github.com/javanhut/TheCarrionLanguage/src/token"
)

type Analyzer struct {
	mu                 sync.RWMutex
	documents          map[string]*Document
	builtins           map[string]*BuiltinInfo
	carriongGrimoires  map[string]*GrimoireInfo
	dynamicLoader      *DynamicLoader
	bifrostIntegration *BifrostIntegration
}

type Document struct {
	URI     string
	Content string
	AST     *ast.Program
	Tokens  []token.Token
	Symbols *SymbolTable
	Version int
}

type SymbolTable struct {
	Grimoires map[string]*GrimoireSymbol // Classes in Carrion
	Spells    map[string]*SpellSymbol    // Functions/Methods in Carrion
	Variables map[string]*VariableSymbol
	Imports   map[string]*ImportSymbol
}

type GrimoireSymbol struct {
	Name      string
	Range     protocol.Range
	InitSpell *SpellSymbol
	Spells    map[string]*SpellSymbol
	IsArcane  bool // Static class
	Inherits  string
	DocString string
}

type SpellSymbol struct {
	Name        string
	Range       protocol.Range
	Parameters  []Parameter
	ReturnType  string
	IsInit      bool
	IsStatic    bool
	IsPrivate   bool
	IsProtected bool
	DocString   string
	Grimoire    string // Parent class name
}

type Parameter struct {
	Name         string
	TypeHint     string
	DefaultValue string
	Range        protocol.Range
}

type VariableSymbol struct {
	Name     string
	Range    protocol.Range
	Type     string
	Value    string
	IsGlobal bool
}

type ImportSymbol struct {
	Name      string
	Range     protocol.Range
	Path      string
	Alias     string
	ClassName string
}

type BuiltinInfo struct {
	Name        string
	Type        string
	Description string
	Parameters  []Parameter
	ReturnType  string
}

type GrimoireInfo struct {
	Name        string
	Description string
	Spells      map[string]*BuiltinInfo
	IsStatic    bool
}

type Workspace struct {
	RootPath  string
	Documents map[string]*Document
}

func New() *Analyzer {
	loader := NewDynamicLoader()
	analyzer := &Analyzer{
		documents:         make(map[string]*Document),
		builtins:          loader.GetBuiltins(),
		carriongGrimoires: loader.GetGrimoires(),
		dynamicLoader:     loader,
	}
	analyzer.bifrostIntegration = NewBifrostIntegration(analyzer)
	return analyzer
}

func NewWorkspace(rootPath string) *Workspace {
	return &Workspace{
		RootPath:  rootPath,
		Documents: make(map[string]*Document),
	}
}

// Dynamic loading and analysis using TheCarrionLanguage parser
func (a *Analyzer) UpdateDocument(uri, content string, program *ast.Program) *Document {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Use provided AST program or parse if nil
	if program == nil {
		l := lexer.New(content)
		p := parser.New(l)
		program = p.ParseProgram()
	}

	// Tokenize for semantic analysis
	tokens := a.tokenizeDocument(content)

	// Build symbol table
	symbols := a.buildSymbolTable(program)

	doc := &Document{
		URI:     uri,
		Content: content,
		AST:     program,
		Tokens:  tokens,
		Symbols: symbols,
		Version: a.getNextVersion(uri),
	}

	a.documents[uri] = doc

	// Auto-load imports from bifrost packages
	if a.bifrostIntegration != nil {
		a.bifrostIntegration.AutoLoadImports(doc)
	}

	return doc
}

func (a *Analyzer) RemoveDocument(uri string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.documents, uri)
}

func (a *Analyzer) GetDocument(uri string) *Document {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.documents[uri]
}

func (a *Analyzer) tokenizeDocument(content string) []token.Token {
	l := lexer.New(content)
	var tokens []token.Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}

	return tokens
}

func (a *Analyzer) buildSymbolTable(program *ast.Program) *SymbolTable {
	symbols := &SymbolTable{
		Grimoires: make(map[string]*GrimoireSymbol),
		Spells:    make(map[string]*SpellSymbol),
		Variables: make(map[string]*VariableSymbol),
		Imports:   make(map[string]*ImportSymbol),
	}

	for _, stmt := range program.Statements {
		switch node := stmt.(type) {
		case *ast.GrimoireDefinition:
			a.analyzeGrimoire(node, symbols)
		case *ast.FunctionDefinition:
			a.analyzeSpell(node, symbols)
		case *ast.AssignStatement:
			a.analyzeVariable(node, symbols)
		case *ast.ImportStatement:
			a.analyzeImport(node, symbols)
		}
	}

	return symbols
}

func (a *Analyzer) analyzeGrimoire(node *ast.GrimoireDefinition, symbols *SymbolTable) {
	grimoire := &GrimoireSymbol{
		Name:     node.Name.Value,
		Range:    a.astNodeToRange(node),
		Spells:   make(map[string]*SpellSymbol),
		IsArcane: false, // TODO: Check for arcane keyword
	}

	if node.Inherits != nil {
		grimoire.Inherits = node.Inherits.Value
	}

	if node.DocString != nil {
		grimoire.DocString = node.DocString.Value
	}

	// Analyze init method
	if node.InitMethod != nil {
		initSpell := &SpellSymbol{
			Name:       "init",
			Range:      a.astNodeToRange(node.InitMethod),
			Parameters: a.extractParameters(node.InitMethod.Parameters),
			IsInit:     true,
			Grimoire:   grimoire.Name,
		}
		if node.InitMethod.DocString != nil {
			initSpell.DocString = node.InitMethod.DocString.Value
		}
		// Analyze init method body
		if node.InitMethod.Body != nil {
			a.analyzeBlockStatement(node.InitMethod.Body, symbols)
		}
		grimoire.InitSpell = initSpell
		symbols.Spells["init"] = initSpell
	}

	// Analyze regular methods
	for _, method := range node.Methods {
		spell := &SpellSymbol{
			Name:       method.Name.Value,
			Range:      a.astNodeToRange(method),
			Parameters: a.extractParameters(method.Parameters),
			Grimoire:   grimoire.Name,
		}
		if method.DocString != nil {
			spell.DocString = method.DocString.Value
		}
		// Analyze method body
		if method.Body != nil {
			a.analyzeBlockStatement(method.Body, symbols)
		}
		grimoire.Spells[spell.Name] = spell
		symbols.Spells[spell.Name] = spell
	}

	symbols.Grimoires[grimoire.Name] = grimoire
}

func (a *Analyzer) analyzeSpell(node *ast.FunctionDefinition, symbols *SymbolTable) {
	spell := &SpellSymbol{
		Name:       node.Name.Value,
		Range:      a.astNodeToRange(node),
		Parameters: a.extractParameters(node.Parameters),
	}

	if node.DocString != nil {
		spell.DocString = node.DocString.Value
	}

	symbols.Spells[spell.Name] = spell

	// Also analyze the function body for local variables
	if node.Body != nil {
		a.analyzeBlockStatement(node.Body, symbols)
	}
}

func (a *Analyzer) analyzeVariable(node *ast.AssignStatement, symbols *SymbolTable) {
	// Extract variable name from assignment
	if ident, ok := node.Name.(*ast.Identifier); ok {
		variable := &VariableSymbol{
			Name:  ident.Value,
			Range: a.astNodeToRange(node),
		}

		// Try to infer type from value
		if node.Value != nil {
			variable.Type = a.inferTypeWithContext(node.Value, symbols)
		}

		symbols.Variables[variable.Name] = variable
	}
}

func (a *Analyzer) analyzeImport(node *ast.ImportStatement, symbols *SymbolTable) {
	importSym := &ImportSymbol{
		Range: a.astNodeToRange(node),
	}

	if node.FilePath != nil {
		importSym.Path = node.FilePath.Value
		// Use filename as default name
		importSym.Name = filepath.Base(strings.TrimSuffix(node.FilePath.Value, ".crl"))
	}

	if node.ClassName != nil {
		importSym.ClassName = node.ClassName.Value
		importSym.Name = node.ClassName.Value
	}

	if node.Alias != nil {
		importSym.Alias = node.Alias.Value
		importSym.Name = node.Alias.Value
	}

	symbols.Imports[importSym.Name] = importSym
}

// analyzeBlockStatement recursively analyzes statements within a block
func (a *Analyzer) analyzeBlockStatement(block *ast.BlockStatement, symbols *SymbolTable) {
	for _, stmt := range block.Statements {
		switch node := stmt.(type) {
		case *ast.AssignStatement:
			a.analyzeVariable(node, symbols)
		case *ast.BlockStatement:
			a.analyzeBlockStatement(node, symbols)
		case *ast.IfStatement:
			if node.Consequence != nil {
				a.analyzeBlockStatement(node.Consequence, symbols)
			}
			if node.Alternative != nil {
				a.analyzeBlockStatement(node.Alternative, symbols)
			}
		case *ast.ForStatement:
			if node.Body != nil {
				a.analyzeBlockStatement(node.Body, symbols)
			}
		case *ast.WhileStatement:
			if node.Body != nil {
				a.analyzeBlockStatement(node.Body, symbols)
			}
			// Add more statement types as needed
		}
	}
}

func (a *Analyzer) extractParameters(params []ast.Expression) []Parameter {
	var parameters []Parameter

	for _, param := range params {
		switch p := param.(type) {
		case *ast.Identifier:
			parameters = append(parameters, Parameter{
				Name:  p.Value,
				Range: a.astNodeToRange(p),
			})
		case *ast.Parameter:
			parameter := Parameter{
				Name:  p.Name.Value,
				Range: a.astNodeToRange(p),
			}

			if p.TypeHint != nil {
				parameter.TypeHint = p.TypeHint.String()
			}

			if p.DefaultValue != nil {
				parameter.DefaultValue = p.DefaultValue.String()
			}

			parameters = append(parameters, parameter)
		}
	}

	return parameters
}

func (a *Analyzer) inferType(expr ast.Expression) string {
	switch node := expr.(type) {
	case *ast.IntegerLiteral:
		return "int"
	case *ast.FloatLiteral:
		return "float"
	case *ast.StringLiteral:
		return "string"
	case *ast.Boolean:
		return "bool"
	case *ast.ArrayLiteral:
		return "array"
	case *ast.HashLiteral:
		return "hash"
	case *ast.TupleLiteral:
		return "tuple"
	case *ast.NoneLiteral:
		return "None"
	case *ast.Identifier:
		// Return the identifier name for lookups
		return node.Value
	default:
		return "unknown"
	}
}

// inferTypeWithContext performs type inference with access to symbol table context
func (a *Analyzer) inferTypeWithContext(expr ast.Expression, symbols *SymbolTable) string {
	switch node := expr.(type) {
	case *ast.IntegerLiteral:
		return "int"
	case *ast.FloatLiteral:
		return "float"
	case *ast.StringLiteral:
		return "string"
	case *ast.Boolean:
		return "bool"
	case *ast.ArrayLiteral:
		return "array"
	case *ast.HashLiteral:
		return "hash"
	case *ast.TupleLiteral:
		return "tuple"
	case *ast.NoneLiteral:
		return "None"
	case *ast.CallExpression:
		// Check if it's a constructor call (grimoire instantiation)
		if ident, ok := node.Function.(*ast.Identifier); ok {
			// Check if the function name matches a known grimoire
			if _, exists := symbols.Grimoires[ident.Value]; exists {
				return ident.Value // Return the grimoire name as the type
			}
			// Check built-in grimoires
			if _, exists := a.carriongGrimoires[ident.Value]; exists {
				return ident.Value
			}
		}
		return "unknown"
	case *ast.Identifier:
		// Check if this identifier refers to a variable with known type
		if variable, exists := symbols.Variables[node.Value]; exists {
			return variable.Type
		}
		return node.Value
	default:
		return "unknown"
	}
}

func (a *Analyzer) astNodeToRange(node ast.Node) protocol.Range {
	// For now, return a placeholder range
	// TODO: Extract actual position information from AST nodes
	return protocol.Range{
		Start: protocol.Position{Line: 0, Character: 0},
		End:   protocol.Position{Line: 0, Character: 0},
	}
}

func (a *Analyzer) getNextVersion(uri string) int {
	if doc, exists := a.documents[uri]; exists {
		return doc.Version + 1
	}
	return 1
}

// RefreshDynamicData reloads built-ins and grimoires from the Carrion runtime
func (a *Analyzer) RefreshDynamicData() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.dynamicLoader.RefreshDynamicData()
	a.builtins = a.dynamicLoader.GetBuiltins()
	a.carriongGrimoires = a.dynamicLoader.GetGrimoires()
}

// LoadBifrostPackage attempts to load a bifrost package
func (a *Analyzer) LoadBifrostPackage(packagePath string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	err := a.dynamicLoader.LoadBifrostPackage(packagePath)
	if err != nil {
		return err
	}

	// Update our local caches
	a.builtins = a.dynamicLoader.GetBuiltins()
	a.carriongGrimoires = a.dynamicLoader.GetGrimoires()

	return nil
}

// GetAvailablePackages returns a list of available bifrost packages
func (a *Analyzer) GetAvailablePackages() map[string]string {
	if a.bifrostIntegration != nil {
		return a.bifrostIntegration.DiscoverAvailablePackages()
	}
	return make(map[string]string)
}

// LoadPackage manually loads a bifrost package
func (a *Analyzer) LoadPackage(packageName string) error {
	if a.bifrostIntegration != nil {
		return a.bifrostIntegration.LoadPackage(packageName)
	}
	return fmt.Errorf("bifrost integration not available")
}

// GetBuiltins returns the currently loaded built-in functions
func (a *Analyzer) GetBuiltins() map[string]*BuiltinInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.builtins
}

// GetGrimoires returns the currently loaded grimoires
func (a *Analyzer) GetGrimoires() map[string]*GrimoireInfo {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.carriongGrimoires
}

// Legacy static initialization - now replaced by dynamic loading
// These functions are kept for reference but should not be used
func initializeBuiltins_DEPRECATED() map[string]*BuiltinInfo {
	builtins := make(map[string]*BuiltinInfo)

	// Core built-in functions
	builtins["print"] = &BuiltinInfo{
		Name:        "print",
		Type:        "function",
		Description: "Print values to output",
		Parameters: []Parameter{
			{Name: "values", TypeHint: "...any"},
		},
		ReturnType: "None",
	}

	builtins["input"] = &BuiltinInfo{
		Name:        "input",
		Type:        "function",
		Description: "Read user input with optional prompt",
		Parameters: []Parameter{
			{Name: "prompt", TypeHint: "string", DefaultValue: "\"\""},
		},
		ReturnType: "string",
	}

	builtins["len"] = &BuiltinInfo{
		Name:        "len",
		Type:        "function",
		Description: "Get length of strings, arrays, or hashes",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "int",
	}

	builtins["type"] = &BuiltinInfo{
		Name:        "type",
		Type:        "function",
		Description: "Get the type of an object",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "string",
	}

	builtins["range"] = &BuiltinInfo{
		Name:        "range",
		Type:        "function",
		Description: "Generate a sequence of numbers",
		Parameters: []Parameter{
			{Name: "start", TypeHint: "int"},
			{Name: "stop", TypeHint: "int", DefaultValue: "None"},
			{Name: "step", TypeHint: "int", DefaultValue: "1"},
		},
		ReturnType: "array",
	}

	builtins["int"] = &BuiltinInfo{
		Name:        "int",
		Type:        "function",
		Description: "Convert to integer",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "any"},
		},
		ReturnType: "int",
	}

	builtins["float"] = &BuiltinInfo{
		Name:        "float",
		Type:        "function",
		Description: "Convert to float",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "any"},
		},
		ReturnType: "float",
	}

	builtins["str"] = &BuiltinInfo{
		Name:        "str",
		Type:        "function",
		Description: "Convert to string",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "any"},
		},
		ReturnType: "string",
	}

	builtins["bool"] = &BuiltinInfo{
		Name:        "bool",
		Type:        "function",
		Description: "Convert to boolean",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "any"},
		},
		ReturnType: "bool",
	}

	builtins["list"] = &BuiltinInfo{
		Name:        "list",
		Type:        "function",
		Description: "Convert to array",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "any"},
		},
		ReturnType: "array",
	}

	builtins["open"] = &BuiltinInfo{
		Name:        "open",
		Type:        "function",
		Description: "Open a file and return File grimoire instance",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
			{Name: "mode", TypeHint: "string", DefaultValue: "\"r\""},
		},
		ReturnType: "File",
	}

	// Add more core built-ins
	builtins["max"] = &BuiltinInfo{
		Name:        "max",
		Type:        "function",
		Description: "Find maximum value",
		Parameters: []Parameter{
			{Name: "values", TypeHint: "...any"},
		},
		ReturnType: "any",
	}

	builtins["abs"] = &BuiltinInfo{
		Name:        "abs",
		Type:        "function",
		Description: "Get absolute value",
		Parameters: []Parameter{
			{Name: "value", TypeHint: "number"},
		},
		ReturnType: "number",
	}

	builtins["enumerate"] = &BuiltinInfo{
		Name:        "enumerate",
		Type:        "function",
		Description: "Enumerate arrays with indices",
		Parameters: []Parameter{
			{Name: "array", TypeHint: "array"},
		},
		ReturnType: "array",
	}

	builtins["pairs"] = &BuiltinInfo{
		Name:        "pairs",
		Type:        "function",
		Description: "Extract key-value pairs from hashes",
		Parameters: []Parameter{
			{Name: "hash", TypeHint: "hash"},
		},
		ReturnType: "array",
	}

	builtins["len"] = &BuiltinInfo{
		Name:        "len",
		Type:        "function",
		Description: "Get length of collection",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "int",
	}

	builtins["str"] = &BuiltinInfo{
		Name:        "str",
		Type:        "function",
		Description: "Convert value to string",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "string",
	}

	builtins["int"] = &BuiltinInfo{
		Name:        "int",
		Type:        "function",
		Description: "Convert value to integer",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "int",
	}

	builtins["float"] = &BuiltinInfo{
		Name:        "float",
		Type:        "function",
		Description: "Convert value to float",
		Parameters: []Parameter{
			{Name: "obj", TypeHint: "any"},
		},
		ReturnType: "float",
	}

	builtins["input"] = &BuiltinInfo{
		Name:        "input",
		Type:        "function",
		Description: "Read input from user",
		Parameters: []Parameter{
			{Name: "prompt", TypeHint: "string", DefaultValue: "\"\""},
		},
		ReturnType: "string",
	}

	builtins["range"] = &BuiltinInfo{
		Name:        "range",
		Type:        "function",
		Description: "Generate range of numbers",
		Parameters: []Parameter{
			{Name: "start", TypeHint: "int"},
			{Name: "stop", TypeHint: "int", DefaultValue: "None"},
			{Name: "step", TypeHint: "int", DefaultValue: "1"},
		},
		ReturnType: "array",
	}

	builtins["pairs"] = &BuiltinInfo{
		Name:        "pairs",
		Type:        "function",
		Description: "Get key-value pairs from hash",
		Parameters: []Parameter{
			{Name: "hash", TypeHint: "hash"},
			{Name: "mode", TypeHint: "string", DefaultValue: "\"both\""},
		},
		ReturnType: "array",
	}

	return builtins
}

// Legacy static initialization - now replaced by dynamic loading
func initializeCarrionGrimoires_DEPRECATED() map[string]*GrimoireInfo {
	grimoires := make(map[string]*GrimoireInfo)

	// String grimoire
	stringSpells := make(map[string]*BuiltinInfo)
	stringSpells["length"] = &BuiltinInfo{Name: "length", Type: "method", Description: "Get string length", ReturnType: "int"}
	stringSpells["lower"] = &BuiltinInfo{Name: "lower", Type: "method", Description: "Convert to lowercase", ReturnType: "string"}
	stringSpells["upper"] = &BuiltinInfo{Name: "upper", Type: "method", Description: "Convert to uppercase", ReturnType: "string"}
	stringSpells["reverse"] = &BuiltinInfo{Name: "reverse", Type: "method", Description: "Reverse string", ReturnType: "string"}
	stringSpells["find"] = &BuiltinInfo{Name: "find", Type: "method", Description: "Find substring index", Parameters: []Parameter{{Name: "substring", TypeHint: "string"}}, ReturnType: "int"}
	stringSpells["contains"] = &BuiltinInfo{Name: "contains", Type: "method", Description: "Check if contains substring", Parameters: []Parameter{{Name: "substring", TypeHint: "string"}}, ReturnType: "bool"}
	stringSpells["char_at"] = &BuiltinInfo{Name: "char_at", Type: "method", Description: "Get character at index", Parameters: []Parameter{{Name: "index", TypeHint: "int"}}, ReturnType: "string"}
	stringSpells["split"] = &BuiltinInfo{Name: "split", Type: "method", Description: "Split string by separator", Parameters: []Parameter{{Name: "separator", TypeHint: "string"}}, ReturnType: "array"}
	stringSpells["join"] = &BuiltinInfo{Name: "join", Type: "method", Description: "Join array of strings", Parameters: []Parameter{{Name: "string_list", TypeHint: "array"}}, ReturnType: "string"}
	stringSpells["strip"] = &BuiltinInfo{Name: "strip", Type: "method", Description: "Remove characters from ends", Parameters: []Parameter{{Name: "characters", TypeHint: "string", DefaultValue: "\" \""}}, ReturnType: "string"}

	grimoires["String"] = &GrimoireInfo{
		Name:        "String",
		Description: "String manipulation grimoire",
		Spells:      stringSpells,
		IsStatic:    false,
	}

	// Array grimoire
	arraySpells := make(map[string]*BuiltinInfo)
	arraySpells["length"] = &BuiltinInfo{Name: "length", Type: "method", Description: "Get array length", ReturnType: "int"}
	arraySpells["append"] = &BuiltinInfo{Name: "append", Type: "method", Description: "Add element to end", Parameters: []Parameter{{Name: "value", TypeHint: "any"}}, ReturnType: "None"}
	arraySpells["get"] = &BuiltinInfo{Name: "get", Type: "method", Description: "Get element at index", Parameters: []Parameter{{Name: "index", TypeHint: "int"}}, ReturnType: "any"}
	arraySpells["set"] = &BuiltinInfo{Name: "set", Type: "method", Description: "Set element at index", Parameters: []Parameter{{Name: "index", TypeHint: "int"}, {Name: "value", TypeHint: "any"}}, ReturnType: "None"}
	arraySpells["is_empty"] = &BuiltinInfo{Name: "is_empty", Type: "method", Description: "Check if empty", ReturnType: "bool"}
	arraySpells["contains"] = &BuiltinInfo{Name: "contains", Type: "method", Description: "Check if contains value", Parameters: []Parameter{{Name: "value", TypeHint: "any"}}, ReturnType: "bool"}
	arraySpells["index_of"] = &BuiltinInfo{Name: "index_of", Type: "method", Description: "Find index of first occurrence", Parameters: []Parameter{{Name: "value", TypeHint: "any"}}, ReturnType: "int"}
	arraySpells["remove"] = &BuiltinInfo{Name: "remove", Type: "method", Description: "Remove first occurrence", Parameters: []Parameter{{Name: "value", TypeHint: "any"}}, ReturnType: "bool"}
	arraySpells["clear"] = &BuiltinInfo{Name: "clear", Type: "method", Description: "Remove all elements", ReturnType: "None"}
	arraySpells["first"] = &BuiltinInfo{Name: "first", Type: "method", Description: "Get first element", ReturnType: "any"}
	arraySpells["last"] = &BuiltinInfo{Name: "last", Type: "method", Description: "Get last element", ReturnType: "any"}
	arraySpells["slice"] = &BuiltinInfo{Name: "slice", Type: "method", Description: "Extract subarray", Parameters: []Parameter{{Name: "start", TypeHint: "int"}, {Name: "end", TypeHint: "int"}}, ReturnType: "array"}
	arraySpells["reverse"] = &BuiltinInfo{Name: "reverse", Type: "method", Description: "Create reversed copy", ReturnType: "array"}
	arraySpells["sort"] = &BuiltinInfo{Name: "sort", Type: "method", Description: "Create sorted copy", ReturnType: "array"}

	grimoires["Array"] = &GrimoireInfo{
		Name:        "Array",
		Description: "Array manipulation grimoire",
		Spells:      arraySpells,
		IsStatic:    false,
	}

	// Integer grimoire
	integerSpells := make(map[string]*BuiltinInfo)
	integerSpells["to_bin"] = &BuiltinInfo{Name: "to_bin", Type: "method", Description: "Convert to binary string", ReturnType: "string"}
	integerSpells["to_oct"] = &BuiltinInfo{Name: "to_oct", Type: "method", Description: "Convert to octal string", ReturnType: "string"}
	integerSpells["to_hex"] = &BuiltinInfo{Name: "to_hex", Type: "method", Description: "Convert to hexadecimal string", ReturnType: "string"}
	integerSpells["abs"] = &BuiltinInfo{Name: "abs", Type: "method", Description: "Absolute value", ReturnType: "int"}
	integerSpells["pow"] = &BuiltinInfo{Name: "pow", Type: "method", Description: "Power operation", Parameters: []Parameter{{Name: "exponent", TypeHint: "int"}}, ReturnType: "int"}
	integerSpells["is_even"] = &BuiltinInfo{Name: "is_even", Type: "method", Description: "Check if even", ReturnType: "bool"}
	integerSpells["is_odd"] = &BuiltinInfo{Name: "is_odd", Type: "method", Description: "Check if odd", ReturnType: "bool"}
	integerSpells["is_prime"] = &BuiltinInfo{Name: "is_prime", Type: "method", Description: "Check if prime number", ReturnType: "bool"}

	grimoires["Integer"] = &GrimoireInfo{
		Name:        "Integer",
		Description: "Integer operations grimoire",
		Spells:      integerSpells,
		IsStatic:    false,
	}

	// Float grimoire
	floatSpells := make(map[string]*BuiltinInfo)
	floatSpells["round"] = &BuiltinInfo{Name: "round", Type: "method", Description: "Round to decimal places", Parameters: []Parameter{{Name: "decimals", TypeHint: "int"}}, ReturnType: "float"}
	floatSpells["floor"] = &BuiltinInfo{Name: "floor", Type: "method", Description: "Floor operation", ReturnType: "int"}
	floatSpells["ceil"] = &BuiltinInfo{Name: "ceil", Type: "method", Description: "Ceiling operation", ReturnType: "int"}
	floatSpells["abs"] = &BuiltinInfo{Name: "abs", Type: "method", Description: "Absolute value", ReturnType: "float"}
	floatSpells["sqrt"] = &BuiltinInfo{Name: "sqrt", Type: "method", Description: "Square root", ReturnType: "float"}
	floatSpells["pow"] = &BuiltinInfo{Name: "pow", Type: "method", Description: "Power operation", Parameters: []Parameter{{Name: "exponent", TypeHint: "float"}}, ReturnType: "float"}
	floatSpells["sin"] = &BuiltinInfo{Name: "sin", Type: "method", Description: "Sine (Taylor series)", ReturnType: "float"}
	floatSpells["cos"] = &BuiltinInfo{Name: "cos", Type: "method", Description: "Cosine (Taylor series)", ReturnType: "float"}
	floatSpells["is_integer"] = &BuiltinInfo{Name: "is_integer", Type: "method", Description: "Check if whole number", ReturnType: "bool"}

	grimoires["Float"] = &GrimoireInfo{
		Name:        "Float",
		Description: "Float operations grimoire",
		Spells:      floatSpells,
		IsStatic:    false,
	}

	// Boolean grimoire
	booleanSpells := make(map[string]*BuiltinInfo)
	booleanSpells["to_int"] = &BuiltinInfo{Name: "to_int", Type: "method", Description: "Convert to integer (1/0)", ReturnType: "int"}
	booleanSpells["negate"] = &BuiltinInfo{Name: "negate", Type: "method", Description: "Logical NOT", ReturnType: "bool"}
	booleanSpells["and_with"] = &BuiltinInfo{Name: "and_with", Type: "method", Description: "Logical AND", Parameters: []Parameter{{Name: "other", TypeHint: "bool"}}, ReturnType: "bool"}
	booleanSpells["or_with"] = &BuiltinInfo{Name: "or_with", Type: "method", Description: "Logical OR", Parameters: []Parameter{{Name: "other", TypeHint: "bool"}}, ReturnType: "bool"}
	booleanSpells["xor_with"] = &BuiltinInfo{Name: "xor_with", Type: "method", Description: "Logical XOR", Parameters: []Parameter{{Name: "other", TypeHint: "bool"}}, ReturnType: "bool"}

	grimoires["Boolean"] = &GrimoireInfo{
		Name:        "Boolean",
		Description: "Boolean operations grimoire",
		Spells:      booleanSpells,
		IsStatic:    false,
	}

	// File grimoire (static methods)
	fileSpells := make(map[string]*BuiltinInfo)
	fileSpells["read"] = &BuiltinInfo{
		Name:        "read",
		Type:        "method",
		Description: "Read file content",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
		},
		ReturnType: "string",
	}
	fileSpells["write"] = &BuiltinInfo{
		Name:        "write",
		Type:        "method",
		Description: "Write content to file",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
			{Name: "content", TypeHint: "string"},
		},
		ReturnType: "None",
	}
	fileSpells["append"] = &BuiltinInfo{
		Name:        "append",
		Type:        "method",
		Description: "Append content to file",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
			{Name: "content", TypeHint: "string"},
		},
		ReturnType: "None",
	}
	fileSpells["exists"] = &BuiltinInfo{
		Name:        "exists",
		Type:        "method",
		Description: "Check if file exists",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
		},
		ReturnType: "bool",
	}
	fileSpells["open"] = &BuiltinInfo{
		Name:        "open",
		Type:        "method",
		Description: "Open file for reading/writing",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
			{Name: "mode", TypeHint: "string", DefaultValue: "\"r\""},
		},
		ReturnType: "FileHandle",
	}

	grimoires["File"] = &GrimoireInfo{
		Name:        "File",
		Description: "File operations grimoire",
		Spells:      fileSpells,
		IsStatic:    true,
	}

	// OS grimoire (static methods)
	osSpells := make(map[string]*BuiltinInfo)
	osSpells["cwd"] = &BuiltinInfo{
		Name:        "cwd",
		Type:        "method",
		Description: "Get current working directory",
		Parameters:  []Parameter{},
		ReturnType:  "string",
	}
	osSpells["listdir"] = &BuiltinInfo{
		Name:        "listdir",
		Type:        "method",
		Description: "List directory contents",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
		},
		ReturnType: "array",
	}
	osSpells["mkdir"] = &BuiltinInfo{
		Name:        "mkdir",
		Type:        "method",
		Description: "Create directory",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
		},
		ReturnType: "None",
	}
	osSpells["remove"] = &BuiltinInfo{
		Name:        "remove",
		Type:        "method",
		Description: "Remove file or directory",
		Parameters: []Parameter{
			{Name: "path", TypeHint: "string"},
		},
		ReturnType: "None",
	}
	osSpells["getenv"] = &BuiltinInfo{
		Name:        "getenv",
		Type:        "method",
		Description: "Get environment variable",
		Parameters: []Parameter{
			{Name: "name", TypeHint: "string"},
		},
		ReturnType: "string",
	}
	osSpells["setenv"] = &BuiltinInfo{
		Name:        "setenv",
		Type:        "method",
		Description: "Set environment variable",
		Parameters: []Parameter{
			{Name: "name", TypeHint: "string"},
			{Name: "value", TypeHint: "string"},
		},
		ReturnType: "None",
	}
	osSpells["run"] = &BuiltinInfo{
		Name:        "run",
		Type:        "method",
		Description: "Run system command",
		Parameters: []Parameter{
			{Name: "command", TypeHint: "string"},
			{Name: "args", TypeHint: "array", DefaultValue: "[]"},
			{Name: "capture", TypeHint: "bool", DefaultValue: "False"},
		},
		ReturnType: "string",
	}

	grimoires["OS"] = &GrimoireInfo{
		Name:        "OS",
		Description: "Operating system operations grimoire",
		Spells:      osSpells,
		IsStatic:    true,
	}

	// Time grimoire
	timeSpells := make(map[string]*BuiltinInfo)
	timeSpells["now"] = &BuiltinInfo{Name: "now", Type: "method", Description: "Current Unix timestamp", ReturnType: "int"}
	timeSpells["now_nano"] = &BuiltinInfo{Name: "now_nano", Type: "method", Description: "Current nanosecond timestamp", ReturnType: "int"}
	timeSpells["sleep"] = &BuiltinInfo{Name: "sleep", Type: "method", Description: "Sleep for duration", Parameters: []Parameter{{Name: "seconds", TypeHint: "float"}}, ReturnType: "None"}
	timeSpells["format"] = &BuiltinInfo{Name: "format", Type: "method", Description: "Format timestamp to string", Parameters: []Parameter{{Name: "timestamp", TypeHint: "int"}, {Name: "format", TypeHint: "string"}}, ReturnType: "string"}
	timeSpells["parse"] = &BuiltinInfo{Name: "parse", Type: "method", Description: "Parse time string to timestamp", Parameters: []Parameter{{Name: "format", TypeHint: "string"}, {Name: "time_str", TypeHint: "string"}}, ReturnType: "int"}
	timeSpells["date"] = &BuiltinInfo{Name: "date", Type: "method", Description: "Get date components [year, month, day]", Parameters: []Parameter{{Name: "timestamp", TypeHint: "int"}}, ReturnType: "array"}
	timeSpells["add_duration"] = &BuiltinInfo{Name: "add_duration", Type: "method", Description: "Add duration to timestamp", Parameters: []Parameter{{Name: "timestamp", TypeHint: "int"}, {Name: "seconds", TypeHint: "int"}}, ReturnType: "int"}
	timeSpells["diff"] = &BuiltinInfo{Name: "diff", Type: "method", Description: "Calculate time difference", Parameters: []Parameter{{Name: "timestamp1", TypeHint: "int"}, {Name: "timestamp2", TypeHint: "int"}}, ReturnType: "int"}

	grimoires["Time"] = &GrimoireInfo{
		Name:        "Time",
		Description: "Time operations grimoire",
		Spells:      timeSpells,
		IsStatic:    true,
	}

	return grimoires
}
