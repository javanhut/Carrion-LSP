package analyzer

import (
	"fmt"
	"strings"

	"github.com/javanhut/CarrionLSP/internal/protocol"
	"github.com/javanhut/TheCarrionLanguage/src/token"
)

// GetCompletions provides auto-completion suggestions for Carrion code
func (a *Analyzer) GetCompletions(uri string, position protocol.Position) []protocol.CompletionItem {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil {
		return nil
	}

	var completions []protocol.CompletionItem

	// Get text at cursor position to determine context
	line := position.Line
	character := position.Character

	lines := strings.Split(doc.Content, "\n")
	if line >= len(lines) {
		return completions
	}

	currentLine := lines[line]
	prefix := ""
	if character <= len(currentLine) {
		prefix = currentLine[:character]
	}

	// Determine completion context
	if strings.HasSuffix(prefix, ".") {
		// Method/property completion
		return a.getMethodCompletions(doc, prefix)
	} else if strings.HasSuffix(prefix, "(") {
		// Function parameter completion
		return a.getParameterCompletions(doc, prefix)
	} else {
		// General completion
		return a.getGeneralCompletions(doc, prefix)
	}
}

func (a *Analyzer) getMethodCompletions(doc *Document, prefix string) []protocol.CompletionItem {
	var completions []protocol.CompletionItem

	// Extract object before the dot
	parts := strings.Split(strings.TrimSuffix(prefix, "."), " ")
	if len(parts) == 0 {
		return completions
	}

	objectName := parts[len(parts)-1]

	// Check if it's a known built-in grimoire (like File, OS, Time)
	if grimoire, exists := a.carriongGrimoires[objectName]; exists {
		for spellName, spell := range grimoire.Spells {
			completions = append(completions, protocol.CompletionItem{
				Label:            spellName,
				Kind:             protocol.CompletionItemKindMethod,
				Detail:           fmt.Sprintf("spell %s(%s) -> %s", spell.Name, a.formatParameters(spell.Parameters), spell.ReturnType),
				Documentation:    spell.Description,
				InsertText:       fmt.Sprintf("%s(${1})", spell.Name),
				InsertTextFormat: protocol.InsertTextFormatSnippet,
			})
		}
	}

	// Check if it's an instance variable with a known type
	if variable, exists := doc.Symbols.Variables[objectName]; exists {
		// Check user-defined grimoires first
		if grimoire, exists := doc.Symbols.Grimoires[variable.Type]; exists {
			for spellName, spell := range grimoire.Spells {
				completions = append(completions, protocol.CompletionItem{
					Label:            spellName,
					Kind:             protocol.CompletionItemKindMethod,
					Detail:           fmt.Sprintf("spell %s(%s) -> %s", spell.Name, a.formatSpellParameters(spell.Parameters), spell.ReturnType),
					Documentation:    spell.DocString,
					InsertText:       fmt.Sprintf("%s(${1})", spell.Name),
					InsertTextFormat: protocol.InsertTextFormatSnippet,
				})
			}
		}

		// Check built-in grimoires for the variable's type
		if grimoire, exists := a.carriongGrimoires[variable.Type]; exists {
			for spellName, spell := range grimoire.Spells {
				completions = append(completions, protocol.CompletionItem{
					Label:            spellName,
					Kind:             protocol.CompletionItemKindMethod,
					Detail:           fmt.Sprintf("spell %s(%s) -> %s", spell.Name, a.formatParameters(spell.Parameters), spell.ReturnType),
					Documentation:    spell.Description,
					InsertText:       fmt.Sprintf("%s(${1})", spell.Name),
					InsertTextFormat: protocol.InsertTextFormatSnippet,
				})
			}
		}

		// Handle primitive types with their respective grimoires
		primitiveToGrimoire := map[string]string{
			"string": "String",
			"int":    "Integer",
			"float":  "Float",
			"bool":   "Boolean",
			"array":  "Array",
		}

		if grimoireName, exists := primitiveToGrimoire[variable.Type]; exists {
			if grimoire, exists := a.carriongGrimoires[grimoireName]; exists {
				for spellName, spell := range grimoire.Spells {
					completions = append(completions, protocol.CompletionItem{
						Label:            spellName,
						Kind:             protocol.CompletionItemKindMethod,
						Detail:           fmt.Sprintf("spell %s(%s) -> %s", spell.Name, a.formatParameters(spell.Parameters), spell.ReturnType),
						Documentation:    spell.Description,
						InsertText:       fmt.Sprintf("%s(${1})", spell.Name),
						InsertTextFormat: protocol.InsertTextFormatSnippet,
					})
				}
			}
		}
	}

	return completions
}

func (a *Analyzer) getParameterCompletions(doc *Document, prefix string) []protocol.CompletionItem {
	// TODO: Implement parameter hint completion
	return nil
}

func (a *Analyzer) getGeneralCompletions(doc *Document, prefix string) []protocol.CompletionItem {
	var completions []protocol.CompletionItem

	// Extract the last token from the prefix for matching
	matchToken := a.extractLastToken(prefix)

	// Carrion keywords
	keywords := []string{
		"spell", "grim", "init", "self", "if", "otherwise", "else", "for", "in", "while",
		"return", "attempt", "ensnare", "resolve", "raise", "import", "as", "match", "case",
		"stop", "skip", "ignore", "True", "False", "None", "and", "or", "not", "main",
		"global", "autoclose", "arcane", "arcanespell", "super", "check",
	}

	for _, keyword := range keywords {
		if strings.HasPrefix(keyword, strings.ToLower(matchToken)) {
			kind := protocol.CompletionItemKindKeyword
			insertText := keyword

			// Add snippets for structural keywords
			switch keyword {
			case "spell":
				insertText = "spell ${1:name}(${2:params}):\n\t${3:body}"
			case "grim":
				insertText = "grim ${1:ClassName}:\n\tinit(${2:params}):\n\t\t${3:body}"
			case "if":
				insertText = "if ${1:condition}:\n\t${2:body}"
			case "for":
				insertText = "for ${1:var} in ${2:iterable}:\n\t${3:body}"
			case "while":
				insertText = "while ${1:condition}:\n\t${2:body}"
			case "attempt":
				insertText = "attempt:\n\t${1:try_body}\nensnare:\n\t${2:except_body}"
			case "autoclose":
				insertText = "autoclose ${1:resource} as ${2:var}:\n\t${3:body}"
			}

			if strings.Contains(insertText, "${") {
				completions = append(completions, protocol.CompletionItem{
					Label:            keyword,
					Kind:             kind,
					InsertText:       insertText,
					InsertTextFormat: protocol.InsertTextFormatSnippet,
				})
			} else {
				completions = append(completions, protocol.CompletionItem{
					Label:      keyword,
					Kind:       kind,
					InsertText: insertText,
				})
			}
		}
	}

	// Built-in functions
	for name, builtin := range a.builtins {
		if strings.HasPrefix(name, matchToken) {
			completions = append(completions, protocol.CompletionItem{
				Label:            name,
				Kind:             protocol.CompletionItemKindFunction,
				Detail:           fmt.Sprintf("%s(%s) -> %s", builtin.Name, a.formatParameters(builtin.Parameters), builtin.ReturnType),
				Documentation:    builtin.Description,
				InsertText:       fmt.Sprintf("%s(${1})", name),
				InsertTextFormat: protocol.InsertTextFormatSnippet,
			})
		}
	}

	// Built-in grimoires
	for name, grimoire := range a.carriongGrimoires {
		if strings.HasPrefix(name, matchToken) {
			completions = append(completions, protocol.CompletionItem{
				Label:         name,
				Kind:          protocol.CompletionItemKindClass,
				Detail:        fmt.Sprintf("grim %s", name),
				Documentation: grimoire.Description,
			})
		}
	}

	// Document symbols
	if doc.Symbols != nil {
		// Grimoires
		for name, grimoire := range doc.Symbols.Grimoires {
			if strings.HasPrefix(name, matchToken) {
				completions = append(completions, protocol.CompletionItem{
					Label:         name,
					Kind:          protocol.CompletionItemKindClass,
					Detail:        fmt.Sprintf("grim %s", name),
					Documentation: grimoire.DocString,
				})
			}
		}

		// Spells
		for name, spell := range doc.Symbols.Spells {
			if strings.HasPrefix(name, matchToken) {
				completions = append(completions, protocol.CompletionItem{
					Label:            name,
					Kind:             protocol.CompletionItemKindFunction,
					Detail:           fmt.Sprintf("spell %s(%s) -> %s", spell.Name, a.formatSpellParameters(spell.Parameters), spell.ReturnType),
					Documentation:    spell.DocString,
					InsertText:       fmt.Sprintf("%s(${1})", name),
					InsertTextFormat: protocol.InsertTextFormatSnippet,
				})
			}
		}

		// Variables
		for name, variable := range doc.Symbols.Variables {
			if strings.HasPrefix(name, matchToken) {
				completions = append(completions, protocol.CompletionItem{
					Label:  name,
					Kind:   protocol.CompletionItemKindVariable,
					Detail: fmt.Sprintf("%s: %s", name, variable.Type),
				})
			}
		}
	}

	return completions
}

// GetHover provides hover information for symbols
func (a *Analyzer) GetHover(uri string, position protocol.Position) *protocol.Hover {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil {
		return nil
	}

	// Find word at position
	word := a.getWordAtPosition(doc.Content, position)
	if word == "" {
		return nil
	}

	// Check built-ins
	if builtin, exists := a.builtins[word]; exists {
		return &protocol.Hover{
			Contents: fmt.Sprintf("**%s**: %s\n\n```carrion\n%s(%s) -> %s\n```\n\n%s",
				builtin.Name, builtin.Type, builtin.Name, a.formatParameters(builtin.Parameters), builtin.ReturnType, builtin.Description),
		}
	}

	// Check grimoires
	if grimoire, exists := a.carriongGrimoires[word]; exists {
		return &protocol.Hover{
			Contents: fmt.Sprintf("**%s**: Grimoire\n\n%s", grimoire.Name, grimoire.Description),
		}
	}

	// Check document symbols
	if doc.Symbols != nil {
		if grimoire, exists := doc.Symbols.Grimoires[word]; exists {
			content := fmt.Sprintf("**%s**: Grimoire", grimoire.Name)
			if grimoire.DocString != "" {
				content += "\n\n" + grimoire.DocString
			}
			if grimoire.Inherits != "" {
				content += fmt.Sprintf("\n\nInherits from: %s", grimoire.Inherits)
			}
			return &protocol.Hover{Contents: content}
		}

		if spell, exists := doc.Symbols.Spells[word]; exists {
			content := fmt.Sprintf("**%s**: Spell\n\n```carrion\nspell %s(%s) -> %s\n```",
				spell.Name, spell.Name, a.formatSpellParameters(spell.Parameters), spell.ReturnType)
			if spell.DocString != "" {
				content += "\n\n" + spell.DocString
			}
			return &protocol.Hover{Contents: content}
		}

		if variable, exists := doc.Symbols.Variables[word]; exists {
			return &protocol.Hover{
				Contents: fmt.Sprintf("**%s**: Variable\n\nType: %s", variable.Name, variable.Type),
			}
		}
	}

	return nil
}

// GetDefinition finds symbol definitions
func (a *Analyzer) GetDefinition(uri string, position protocol.Position) []protocol.Location {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil {
		return nil
	}

	word := a.getWordAtPosition(doc.Content, position)
	if word == "" {
		return nil
	}

	var locations []protocol.Location

	// Check document symbols
	if doc.Symbols != nil {
		if grimoire, exists := doc.Symbols.Grimoires[word]; exists {
			locations = append(locations, protocol.Location{
				URI:   uri,
				Range: grimoire.Range,
			})
		}

		if spell, exists := doc.Symbols.Spells[word]; exists {
			locations = append(locations, protocol.Location{
				URI:   uri,
				Range: spell.Range,
			})
		}

		if variable, exists := doc.Symbols.Variables[word]; exists {
			locations = append(locations, protocol.Location{
				URI:   uri,
				Range: variable.Range,
			})
		}
	}

	return locations
}

// GetReferences finds all references to a symbol
func (a *Analyzer) GetReferences(uri string, position protocol.Position, includeDeclaration bool) []protocol.Location {
	// TODO: Implement reference finding
	return nil
}

// GetDocumentSymbols returns document outline
func (a *Analyzer) GetDocumentSymbols(uri string) []protocol.DocumentSymbol {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil || doc.Symbols == nil {
		return nil
	}

	var symbols []protocol.DocumentSymbol

	// Add grimoires
	for _, grimoire := range doc.Symbols.Grimoires {
		symbol := protocol.DocumentSymbol{
			Name:           grimoire.Name,
			Kind:           protocol.SymbolKindClass,
			Range:          grimoire.Range,
			SelectionRange: grimoire.Range,
		}

		// Add spells as children
		for _, spell := range grimoire.Spells {
			symbol.Children = append(symbol.Children, protocol.DocumentSymbol{
				Name:           spell.Name,
				Kind:           protocol.SymbolKindMethod,
				Range:          spell.Range,
				SelectionRange: spell.Range,
			})
		}

		symbols = append(symbols, symbol)
	}

	// Add standalone spells
	for _, spell := range doc.Symbols.Spells {
		if spell.Grimoire == "" {
			symbols = append(symbols, protocol.DocumentSymbol{
				Name:           spell.Name,
				Kind:           protocol.SymbolKindFunction,
				Range:          spell.Range,
				SelectionRange: spell.Range,
			})
		}
	}

	return symbols
}

// GetSemanticTokens provides semantic token information
func (a *Analyzer) GetSemanticTokens(uri string) *protocol.SemanticTokens {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil {
		return nil
	}

	var data []int

	for _, tok := range doc.Tokens {
		tokenType := a.mapTokenToSemanticType(tok.Type)
		if tokenType >= 0 {
			// LSP semantic tokens format: [deltaLine, deltaStart, length, tokenType, tokenModifiers]
			data = append(data, tok.Line, tok.Column, len(tok.Literal), tokenType, 0)
		}
	}

	return &protocol.SemanticTokens{
		Data: data,
	}
}

// FormatDocument provides document formatting
func (a *Analyzer) FormatDocument(uri string, options protocol.FormattingOptions) []protocol.TextEdit {
	a.mu.RLock()
	defer a.mu.RUnlock()

	doc := a.documents[uri]
	if doc == nil {
		return nil
	}

	// Create formatter with options
	formatter := NewCarrionFormatter(options)

	// Format the document
	edits, err := formatter.FormatDocument(doc.Content)
	if err != nil {
		// If formatting fails, return no edits
		return nil
	}

	return edits
}

// Helper functions

func (a *Analyzer) getWordAtPosition(content string, position protocol.Position) string {
	lines := strings.Split(content, "\n")
	if position.Line >= len(lines) {
		return ""
	}

	line := lines[position.Line]
	if position.Character >= len(line) {
		return ""
	}

	// Find word boundaries
	start := position.Character
	end := position.Character

	// Go backwards to find start
	for start > 0 && (isAlphaNumeric(rune(line[start-1])) || line[start-1] == '_') {
		start--
	}

	// Go forwards to find end
	for end < len(line) && (isAlphaNumeric(rune(line[end])) || line[end] == '_') {
		end++
	}

	if start >= end {
		return ""
	}

	return line[start:end]
}

func (a *Analyzer) formatParameters(params []Parameter) string {
	var parts []string
	for _, param := range params {
		part := param.Name
		if param.TypeHint != "" {
			part += ": " + param.TypeHint
		}
		if param.DefaultValue != "" {
			part += " = " + param.DefaultValue
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, ", ")
}

func (a *Analyzer) formatSpellParameters(params []Parameter) string {
	return a.formatParameters(params)
}

func (a *Analyzer) mapTokenToSemanticType(tokenType token.TokenType) int {
	switch tokenType {
	case token.SPELL, token.GRIMOIRE, token.IF, token.ELSE, token.FOR, token.WHILE, token.RETURN, token.IMPORT, token.AS:
		return 0 // keyword
	case token.STRING, token.DOCSTRING:
		return 1 // string
	case token.INT, token.FLOAT:
		return 2 // number
	case token.ASSIGN, token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.EQ, token.NOT_EQ:
		return 3 // operator
	case token.IDENT:
		return 4 // variable
	default:
		return -1 // skip
	}
}

func isAlphaNumeric(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')
}

// extractLastToken extracts the last identifier token from a text string
func (a *Analyzer) extractLastToken(text string) string {
	// Find the last identifier by working backwards from the end
	end := len(text)

	// Skip trailing whitespace
	for end > 0 && (text[end-1] == ' ' || text[end-1] == '\t') {
		end--
	}

	if end == 0 {
		return ""
	}

	// Find the start of the identifier
	start := end
	for start > 0 && (isAlphaNumeric(rune(text[start-1])) || text[start-1] == '_') {
		start--
	}

	return text[start:end]
}
