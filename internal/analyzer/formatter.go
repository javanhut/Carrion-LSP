package analyzer

import (
	"fmt"
	"strings"

	"github.com/javanhut/CarrionLSP/internal/protocol"
	"github.com/javanhut/TheCarrionLanguage/src/ast"
	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

// CarrionFormatter handles formatting of Carrion code according to language conventions
type CarrionFormatter struct {
	tabSize      int
	insertSpaces bool
	indentLevel  int
	options      protocol.FormattingOptions
}

// NewCarrionFormatter creates a new formatter with the given options
func NewCarrionFormatter(options protocol.FormattingOptions) *CarrionFormatter {
	return &CarrionFormatter{
		tabSize:      options.TabSize,
		insertSpaces: options.InsertSpaces,
		options:      options,
	}
}

// FormatDocument formats the entire Carrion document
func (f *CarrionFormatter) FormatDocument(content string) ([]protocol.TextEdit, error) {
	// Parse the content
	l := lexer.New(content)
	p := parser.New(l)
	program := p.ParseProgram()

	// Check for parsing errors
	if len(p.Errors()) > 0 {
		return nil, fmt.Errorf("parsing errors: %v", p.Errors())
	}

	// Format the AST
	formatted := f.formatProgram(program)

	// Apply final formatting rules
	formatted = f.applyFinalFormatting(formatted)

	// Create a single text edit that replaces the entire document
	lines := strings.Split(content, "\n")
	endLine := len(lines) - 1
	endChar := 0
	if endLine >= 0 && endLine < len(lines) {
		endChar = len(lines[endLine])
	}

	edit := protocol.TextEdit{
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End:   protocol.Position{Line: endLine, Character: endChar},
		},
		NewText: formatted,
	}

	return []protocol.TextEdit{edit}, nil
}

// formatProgram formats the entire AST program
func (f *CarrionFormatter) formatProgram(program *ast.Program) string {
	var parts []string
	f.indentLevel = 0

	for i, stmt := range program.Statements {
		formatted := f.formatStatement(stmt)
		if formatted != "" {
			parts = append(parts, formatted)
		}

		// Add blank lines between major blocks
		if i < len(program.Statements)-1 {
			switch stmt.(type) {
			case *ast.GrimoireDefinition:
				// Add double blank line before main block
				if _, isNextMain := program.Statements[i+1].(*ast.MainStatement); isNextMain {
					parts = append(parts, "", "")
				}
			case *ast.FunctionDefinition:
				// Add blank line after standalone functions
				parts = append(parts, "")
			}
		}
	}

	return strings.Join(parts, "\n")
}

// formatStatement formats individual statements
func (f *CarrionFormatter) formatStatement(stmt ast.Statement) string {
	if stmt == nil {
		return ""
	}

	switch node := stmt.(type) {
	case *ast.FunctionDefinition:
		return f.formatFunctionDefinition(node)
	case *ast.GrimoireDefinition:
		return f.formatGrimoireDefinition(node)
	case *ast.AssignStatement:
		return f.formatAssignStatement(node)
	case *ast.ExpressionStatement:
		return f.formatExpressionStatement(node)
	case *ast.IfStatement:
		return f.formatIfStatement(node)
	case *ast.ForStatement:
		return f.formatForStatement(node)
	case *ast.WhileStatement:
		return f.formatWhileStatement(node)
	case *ast.ReturnStatement:
		return f.formatReturnStatement(node)
	case *ast.AttemptStatement:
		return f.formatAttemptStatement(node)
	case *ast.WithStatement:
		return f.formatWithStatement(node)
	case *ast.ImportStatement:
		return f.formatImportStatement(node)
	case *ast.MatchStatement:
		return f.formatMatchStatement(node)
	case *ast.MainStatement:
		return f.formatMainStatement(node)
	case *ast.GlobalStatement:
		return f.formatGlobalStatement(node)
	default:
		return f.indent() + stmt.String()
	}
}

// formatFunctionDefinition formats spell definitions
func (f *CarrionFormatter) formatFunctionDefinition(node *ast.FunctionDefinition) string {
	if node == nil {
		return ""
	}

	var parts []string

	// Add docstring if present and it's the first statement
	if node.DocString != nil {
		parts = append(parts, f.indent()+fmt.Sprintf(`"""%s"""`, node.DocString.Value))
	}

	// Format function signature
	var params string
	if node.Parameters != nil {
		params = f.formatParameters(node.Parameters)
	}

	var name string
	if node.Name != nil {
		name = node.Name.Value
	}

	signature := fmt.Sprintf("spell %s(%s):", name, params)
	parts = append(parts, f.indent()+signature)

	// Format body
	f.indentLevel++
	bodyParts := f.formatBlockStatement(node.Body)
	parts = append(parts, bodyParts...)
	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatInitMethod formats init method definitions (without "spell" keyword)
func (f *CarrionFormatter) formatInitMethod(node *ast.FunctionDefinition) string {
	if node == nil {
		return ""
	}

	var parts []string

	// Add docstring if present
	if node.DocString != nil {
		parts = append(parts, f.indent()+fmt.Sprintf(`"""%s"""`, node.DocString.Value))
	}

	// Format init signature (without "spell" keyword)
	var params string
	if node.Parameters != nil {
		params = f.formatParameters(node.Parameters)
	}

	signature := fmt.Sprintf("init(%s):", params)
	parts = append(parts, f.indent()+signature)

	// Format body
	f.indentLevel++
	bodyParts := f.formatBlockStatement(node.Body)
	parts = append(parts, bodyParts...)
	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatGrimoireDefinition formats grimoire (class) definitions
func (f *CarrionFormatter) formatGrimoireDefinition(node *ast.GrimoireDefinition) string {
	var parts []string

	// Add docstring if present
	if node.DocString != nil {
		parts = append(parts, f.indent()+fmt.Sprintf(`"""%s"""`, node.DocString.Value))
	}

	// Format grimoire declaration
	grimoireDecl := fmt.Sprintf("grim %s:", node.Name.Value)
	if node.Inherits != nil {
		grimoireDecl = fmt.Sprintf("grim %s(%s):", node.Name.Value, node.Inherits.Value)
	}
	parts = append(parts, f.indent()+grimoireDecl)

	f.indentLevel++

	// Format init method
	if node.InitMethod != nil {
		parts = append(parts, f.formatInitMethod(node.InitMethod))
	}

	// Format other methods
	for _, method := range node.Methods {
		parts = append(parts, f.formatFunctionDefinition(method))
	}

	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatAssignStatement formats variable assignments
func (f *CarrionFormatter) formatAssignStatement(node *ast.AssignStatement) string {
	if node == nil {
		return ""
	}

	name := f.formatExpression(node.Name)
	value := f.formatExpression(node.Value)
	operator := node.Operator
	if operator == "" {
		operator = "="
	}

	// Add spacing around operators
	return f.indent() + fmt.Sprintf("%s %s %s", name, operator, value)
}

// formatExpressionStatement formats expression statements
func (f *CarrionFormatter) formatExpressionStatement(node *ast.ExpressionStatement) string {
	if node == nil || node.Expression == nil {
		return ""
	}
	return f.indent() + f.formatExpression(node.Expression)
}

// formatIfStatement formats if/otherwise/else statements
func (f *CarrionFormatter) formatIfStatement(node *ast.IfStatement) string {
	var parts []string

	// Main if clause
	condition := f.formatExpression(node.Condition)
	parts = append(parts, f.indent()+fmt.Sprintf("if %s:", condition))

	// If body
	f.indentLevel++
	ifBody := f.formatBlockStatement(node.Consequence)
	parts = append(parts, ifBody...)
	f.indentLevel--

	// Otherwise clauses
	for _, branch := range node.OtherwiseBranches {
		branchCondition := f.formatExpression(branch.Condition)
		parts = append(parts, f.indent()+fmt.Sprintf("otherwise %s:", branchCondition))

		f.indentLevel++
		branchBody := f.formatBlockStatement(branch.Consequence)
		parts = append(parts, branchBody...)
		f.indentLevel--
	}

	// Else clause
	if node.Alternative != nil {
		parts = append(parts, f.indent()+"else:")
		f.indentLevel++
		elseBody := f.formatBlockStatement(node.Alternative)
		parts = append(parts, elseBody...)
		f.indentLevel--
	}

	return strings.Join(parts, "\n")
}

// formatForStatement formats for loops
func (f *CarrionFormatter) formatForStatement(node *ast.ForStatement) string {
	var parts []string

	variable := f.formatExpression(node.Variable)
	iterable := f.formatExpression(node.Iterable)
	parts = append(parts, f.indent()+fmt.Sprintf("for %s in %s:", variable, iterable))

	// For body
	f.indentLevel++
	forBody := f.formatBlockStatement(node.Body)
	parts = append(parts, forBody...)
	f.indentLevel--

	// Else clause for for loop
	if node.Alternative != nil {
		parts = append(parts, f.indent()+"else:")
		f.indentLevel++
		elseBody := f.formatBlockStatement(node.Alternative)
		parts = append(parts, elseBody...)
		f.indentLevel--
	}

	return strings.Join(parts, "\n")
}

// formatWhileStatement formats while loops
func (f *CarrionFormatter) formatWhileStatement(node *ast.WhileStatement) string {
	var parts []string

	condition := f.formatExpression(node.Condition)
	parts = append(parts, f.indent()+fmt.Sprintf("while %s:", condition))

	f.indentLevel++
	body := f.formatBlockStatement(node.Body)
	parts = append(parts, body...)
	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatReturnStatement formats return statements
func (f *CarrionFormatter) formatReturnStatement(node *ast.ReturnStatement) string {
	if node == nil {
		return ""
	}
	if node.ReturnValue == nil {
		return f.indent() + "return"
	}
	value := f.formatExpression(node.ReturnValue)
	return f.indent() + fmt.Sprintf("return %s", value)
}

// formatAttemptStatement formats attempt/ensnare/resolve blocks
func (f *CarrionFormatter) formatAttemptStatement(node *ast.AttemptStatement) string {
	var parts []string

	// Attempt block
	parts = append(parts, f.indent()+"attempt:")
	f.indentLevel++
	attemptBody := f.formatBlockStatement(node.TryBlock)
	parts = append(parts, attemptBody...)
	f.indentLevel--

	// Ensnare clauses
	for _, ensnare := range node.EnsnareClauses {
		ensnareClause := f.indent() + "ensnare"
		if ensnare.Condition != nil {
			condition := f.formatExpression(ensnare.Condition)
			ensnareClause += fmt.Sprintf("(%s)", condition)
		}
		if ensnare.Alias != nil {
			ensnareClause += fmt.Sprintf(" as %s", ensnare.Alias.Value)
		}
		ensnareClause += ":"
		parts = append(parts, ensnareClause)

		f.indentLevel++
		ensnareBody := f.formatBlockStatement(ensnare.Consequence)
		parts = append(parts, ensnareBody...)
		f.indentLevel--
	}

	// Resolve block
	if node.ResolveBlock != nil {
		parts = append(parts, f.indent()+"resolve:")
		f.indentLevel++
		resolveBody := f.formatBlockStatement(node.ResolveBlock)
		parts = append(parts, resolveBody...)
		f.indentLevel--
	}

	return strings.Join(parts, "\n")
}

// formatWithStatement formats autoclose blocks
func (f *CarrionFormatter) formatWithStatement(node *ast.WithStatement) string {
	var parts []string

	expr := f.formatExpression(node.Expression)
	variable := node.Variable.Value
	parts = append(parts, f.indent()+fmt.Sprintf("autoclose %s as %s:", expr, variable))

	f.indentLevel++
	body := f.formatBlockStatement(node.Body)
	parts = append(parts, body...)
	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatImportStatement formats import statements
func (f *CarrionFormatter) formatImportStatement(node *ast.ImportStatement) string {
	importStmt := f.indent() + "import "

	if node.FilePath != nil {
		importStmt += fmt.Sprintf(`"%s"`, node.FilePath.Value)
	}

	if node.ClassName != nil {
		importStmt += fmt.Sprintf(".%s", node.ClassName.Value)
	}

	if node.Alias != nil {
		importStmt += fmt.Sprintf(" as %s", node.Alias.Value)
	}

	return importStmt
}

// formatMatchStatement formats match statements
func (f *CarrionFormatter) formatMatchStatement(node *ast.MatchStatement) string {
	var parts []string

	value := f.formatExpression(node.MatchValue)
	parts = append(parts, f.indent()+fmt.Sprintf("match %s:", value))

	f.indentLevel++
	for _, caseClause := range node.Cases {
		condition := f.formatExpression(caseClause.Condition)
		parts = append(parts, f.indent()+fmt.Sprintf("case %s:", condition))

		f.indentLevel++
		caseBody := f.formatBlockStatement(caseClause.Body)
		parts = append(parts, caseBody...)
		f.indentLevel--
	}

	if node.Default != nil {
		parts = append(parts, f.indent()+"default:")
		f.indentLevel++
		defaultBody := f.formatBlockStatement(node.Default.Body)
		parts = append(parts, defaultBody...)
		f.indentLevel--
	}
	f.indentLevel--

	return strings.Join(parts, "\n")
}

// formatMainStatement formats main blocks
func (f *CarrionFormatter) formatMainStatement(node *ast.MainStatement) string {
	var parts []string
	parts = append(parts, f.indent()+"main:")

	if node.Body != nil {
		f.indentLevel++
		body := f.formatBlockStatement(node.Body)
		parts = append(parts, body...)
		f.indentLevel--
	}

	return strings.Join(parts, "\n")
}

// formatGlobalStatement formats global declarations
func (f *CarrionFormatter) formatGlobalStatement(node *ast.GlobalStatement) string {
	names := make([]string, len(node.Names))
	for i, name := range node.Names {
		names[i] = name.Value
	}
	return f.indent() + fmt.Sprintf("global %s", strings.Join(names, ", "))
}

// formatBlockStatement formats block statements
func (f *CarrionFormatter) formatBlockStatement(block *ast.BlockStatement) []string {
	var parts []string

	if block == nil {
		return parts
	}

	for _, stmt := range block.Statements {
		formatted := f.formatStatement(stmt)
		if formatted != "" {
			parts = append(parts, formatted)
		}
	}

	return parts
}

// formatExpression formats expressions
func (f *CarrionFormatter) formatExpression(expr ast.Expression) string {
	if expr == nil {
		return ""
	}

	switch node := expr.(type) {
	case *ast.Identifier:
		return node.Value
	case *ast.IntegerLiteral:
		return fmt.Sprintf("%d", node.Value)
	case *ast.FloatLiteral:
		return fmt.Sprintf("%g", node.Value)
	case *ast.StringLiteral:
		return fmt.Sprintf(`"%s"`, node.Value)
	case *ast.Boolean:
		if node.Value {
			return "True"
		}
		return "False"
	case *ast.NoneLiteral:
		return "None"
	case *ast.ArrayLiteral:
		return f.formatArrayLiteral(node)
	case *ast.HashLiteral:
		return f.formatHashLiteral(node)
	case *ast.TupleLiteral:
		return f.formatTupleLiteral(node)
	case *ast.CallExpression:
		return f.formatCallExpression(node)
	case *ast.InfixExpression:
		return f.formatInfixExpression(node)
	case *ast.PrefixExpression:
		return f.formatPrefixExpression(node)
	case *ast.IndexExpression:
		return f.formatIndexExpression(node)
	case *ast.DotExpression:
		return f.formatDotExpression(node)
	case *ast.SliceExpression:
		return f.formatSliceExpression(node)
	default:
		return expr.String()
	}
}

// formatParameters formats function parameters
func (f *CarrionFormatter) formatParameters(params []ast.Expression) string {
	var parts []string

	for _, param := range params {
		switch p := param.(type) {
		case *ast.Identifier:
			parts = append(parts, p.Value)
		case *ast.Parameter:
			paramStr := p.Name.Value
			if p.TypeHint != nil {
				paramStr += ": " + f.formatExpression(p.TypeHint)
			}
			if p.DefaultValue != nil {
				paramStr += " = " + f.formatExpression(p.DefaultValue)
			}
			parts = append(parts, paramStr)
		default:
			parts = append(parts, f.formatExpression(param))
		}
	}

	return strings.Join(parts, ", ")
}

// formatArrayLiteral formats array literals
func (f *CarrionFormatter) formatArrayLiteral(node *ast.ArrayLiteral) string {
	if len(node.Elements) == 0 {
		return "[]"
	}

	var elements []string
	for _, elem := range node.Elements {
		elements = append(elements, f.formatExpression(elem))
	}

	if len(elements) <= 3 {
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
	}

	// Multi-line format for longer arrays
	return fmt.Sprintf("[\n%s%s\n%s]",
		f.indentString(f.indentLevel+1),
		strings.Join(elements, ",\n"+f.indentString(f.indentLevel+1)),
		f.indentString(f.indentLevel))
}

// formatHashLiteral formats hash literals
func (f *CarrionFormatter) formatHashLiteral(node *ast.HashLiteral) string {
	if len(node.Pairs) == 0 {
		return "{}"
	}

	var pairs []string
	for key, value := range node.Pairs {
		keyStr := f.formatExpression(key)
		valueStr := f.formatExpression(value)
		pairs = append(pairs, fmt.Sprintf("%s: %s", keyStr, valueStr))
	}

	if len(pairs) <= 2 {
		return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
	}

	// Multi-line format for longer hashes
	return fmt.Sprintf("{\n%s%s\n%s}",
		f.indentString(f.indentLevel+1),
		strings.Join(pairs, ",\n"+f.indentString(f.indentLevel+1)),
		f.indentString(f.indentLevel))
}

// formatTupleLiteral formats tuple literals
func (f *CarrionFormatter) formatTupleLiteral(node *ast.TupleLiteral) string {
	var elements []string
	for _, elem := range node.Elements {
		elements = append(elements, f.formatExpression(elem))
	}
	return fmt.Sprintf("(%s)", strings.Join(elements, ", "))
}

// formatCallExpression formats function calls
func (f *CarrionFormatter) formatCallExpression(node *ast.CallExpression) string {
	function := f.formatExpression(node.Function)
	var args []string
	for _, arg := range node.Arguments {
		args = append(args, f.formatExpression(arg))
	}
	return fmt.Sprintf("%s(%s)", function, strings.Join(args, ", "))
}

// formatInfixExpression formats infix expressions
func (f *CarrionFormatter) formatInfixExpression(node *ast.InfixExpression) string {
	left := f.formatExpression(node.Left)
	right := f.formatExpression(node.Right)
	return fmt.Sprintf("%s %s %s", left, node.Operator, right)
}

// formatPrefixExpression formats prefix expressions
func (f *CarrionFormatter) formatPrefixExpression(node *ast.PrefixExpression) string {
	right := f.formatExpression(node.Right)
	return fmt.Sprintf("%s%s", node.Operator, right)
}

// formatIndexExpression formats index expressions
func (f *CarrionFormatter) formatIndexExpression(node *ast.IndexExpression) string {
	left := f.formatExpression(node.Left)
	index := f.formatExpression(node.Index)
	return fmt.Sprintf("%s[%s]", left, index)
}

// formatDotExpression formats dot expressions
func (f *CarrionFormatter) formatDotExpression(node *ast.DotExpression) string {
	left := f.formatExpression(node.Left)
	right := node.Right.Value
	return fmt.Sprintf("%s.%s", left, right)
}

// formatSliceExpression formats slice expressions
func (f *CarrionFormatter) formatSliceExpression(node *ast.SliceExpression) string {
	left := f.formatExpression(node.Left)
	var start, end string
	if node.Start != nil {
		start = f.formatExpression(node.Start)
	}
	if node.End != nil {
		end = f.formatExpression(node.End)
	}
	return fmt.Sprintf("%s[%s:%s]", left, start, end)
}

// applyFinalFormatting applies final formatting rules
func (f *CarrionFormatter) applyFinalFormatting(content string) string {
	lines := strings.Split(content, "\n")
	var formattedLines []string

	for _, line := range lines {
		// Trim trailing whitespace
		if f.options.TrimTrailingWhitespace {
			line = strings.TrimRight(line, " \t")
		}
		formattedLines = append(formattedLines, line)
	}

	result := strings.Join(formattedLines, "\n")

	// Ensure exactly one newline at the end (matching Carrion style)
	result = strings.TrimRight(result, "\n") + "\n"

	// Apply additional formatting options if specified
	if f.options.InsertFinalNewline && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}

	// Trim final newlines if requested
	if f.options.TrimFinalNewlines {
		result = strings.TrimRight(result, "\n") + "\n"
	}

	return result
}

// indent returns the current indentation string
func (f *CarrionFormatter) indent() string {
	return f.indentString(f.indentLevel)
}

// indentString returns indentation string for the given level
func (f *CarrionFormatter) indentString(level int) string {
	if f.insertSpaces {
		return strings.Repeat(" ", level*f.tabSize)
	}
	return strings.Repeat("\t", level)
}
