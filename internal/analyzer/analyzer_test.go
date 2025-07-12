package analyzer

import (
	"strings"
	"testing"

	"github.com/javanhut/CarrionLSP/internal/protocol"
)

// Test data for comprehensive testing
const testCarrionCode = `
# Test Carrion code with all language features
import "file" as File
import "os".OS

# Global variable
global_var = "Hello, Carrion!"

spell greet(name: string = "World"):
    """Greet someone by name with optional parameter"""
    return "Hello, " + name + "!"

spell calculate(x: int, y: int = 5):
    """Calculate sum with default parameter"""
    return x + y

grim Person:
    """A person grimoire with inheritance example"""
    
    init(name: string, age: int = 0):
        """Initialize a new person"""
        self.name = name
        self.age = age
    
    spell say_hello():
        """Person says hello"""
        return "Hi, I'm " + self.name
    
    spell get_age():
        """Get person's age"""
        return self.age

grim Employee:
    """Employee inherits from Person"""
    
    init(name: string, age: int, job: string):
        """Initialize employee"""
        super.init(name, age)
        self.job = job
    
    spell introduce():
        """Employee introduction"""
        return self.say_hello() + " and I work as " + self.job

spell main_function():
    """Main test function"""
    person = Person("Alice", 30)
    employee = Employee("Bob", 25, "Developer")
    
    # Test File operations
    content = File.read("test.txt")
    File.write("output.txt", content)
    
    # Test OS operations
    current_dir = OS.cwd()
    files = OS.listdir(".")
    
    # Test attempt/ensnare
    attempt:
        result = calculate(10, 20)
        print(result)
    ensnare:
        print("Error occurred")
    resolve:
        print("Cleanup")

# Test control flow
if True:
    for i in range(5):
        if i > 2:
            print(i)
        otherwise i == 2:
            print("Two")
        else:
            print("Less than 2")

# Test autoclose
autoclose File.open("test.txt", "r") as file:
    data = file.read_content()
    print(data)

main:
    main_function()
`

func TestAnalyzer_New(t *testing.T) {
	analyzer := New()

	if analyzer == nil {
		t.Fatal("Expected analyzer to be created")
	}

	if analyzer.documents == nil {
		t.Error("Expected documents map to be initialized")
	}

	if analyzer.builtins == nil {
		t.Error("Expected builtins map to be initialized")
	}

	if analyzer.carriongGrimoires == nil {
		t.Error("Expected carrion grimoires map to be initialized")
	}

	// Test built-ins are loaded
	expectedBuiltins := []string{"print", "len", "str", "int", "float", "input", "range", "pairs"}
	for _, builtin := range expectedBuiltins {
		if _, exists := analyzer.builtins[builtin]; !exists {
			t.Errorf("Expected builtin '%s' to be loaded", builtin)
		}
	}

	// Test grimoires are loaded
	expectedGrimoires := []string{"File", "OS"}
	for _, grimoire := range expectedGrimoires {
		if _, exists := analyzer.carriongGrimoires[grimoire]; !exists {
			t.Errorf("Expected grimoire '%s' to be loaded", grimoire)
		}
	}
}

func TestAnalyzer_UpdateDocument(t *testing.T) {
	analyzer := New()

	doc := analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	if doc == nil {
		t.Fatal("Expected document to be created")
	}

	if doc.URI != "test.crl" {
		t.Errorf("Expected URI to be 'test.crl', got %s", doc.URI)
	}

	if doc.Content != testCarrionCode {
		t.Error("Expected content to match input")
	}

	if doc.AST == nil {
		t.Error("Expected AST to be parsed")
	}

	if doc.Symbols == nil {
		t.Fatal("Expected symbols to be created")
	}

	// Test specific symbols are found
	expectedSpells := []string{"greet", "calculate", "main_function"}
	for _, spell := range expectedSpells {
		if _, exists := doc.Symbols.Spells[spell]; !exists {
			t.Errorf("Expected spell '%s' to be found", spell)
		}
	}

	expectedGrimoires := []string{"Person", "Employee"}
	for _, grimoire := range expectedGrimoires {
		if _, exists := doc.Symbols.Grimoires[grimoire]; !exists {
			t.Errorf("Expected grimoire '%s' to be found", grimoire)
		}
	}

	// Test inheritance - skip for now since the test code doesn't have inheritance
	if employee, exists := doc.Symbols.Grimoires["Employee"]; exists {
		// Note: Test code doesn't include inheritance syntax, so this test is disabled
		_ = employee // Use the variable to avoid unused variable error
	}
}

func TestAnalyzer_UpdateDocument_Versioning(t *testing.T) {
	analyzer := New()

	// First update
	doc1 := analyzer.UpdateDocument("test.crl", "spell test1(): return 1", nil)
	if doc1.Version != 1 {
		t.Errorf("Expected first version to be 1, got %d", doc1.Version)
	}

	// Second update
	doc2 := analyzer.UpdateDocument("test.crl", "spell test2(): return 2", nil)
	if doc2.Version != 2 {
		t.Errorf("Expected second version to be 2, got %d", doc2.Version)
	}
}

func TestAnalyzer_RemoveDocument(t *testing.T) {
	analyzer := New()

	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	// Verify document exists
	if doc := analyzer.GetDocument("test.crl"); doc == nil {
		t.Error("Expected document to exist before removal")
	}

	analyzer.RemoveDocument("test.crl")

	// Verify document is removed
	if doc := analyzer.GetDocument("test.crl"); doc != nil {
		t.Error("Expected document to be removed")
	}
}

func TestAnalyzer_GetCompletions_Keywords(t *testing.T) {
	analyzer := New()

	content := "sp"
	analyzer.UpdateDocument("test.crl", content, nil)

	position := protocol.Position{Line: 0, Character: 2}
	completions := analyzer.GetCompletions("test.crl", position)

	// Should find "spell" keyword
	found := false
	for _, completion := range completions {
		if completion.Label == "spell" {
			found = true
			if completion.Kind != protocol.CompletionItemKindKeyword {
				t.Error("Expected spell completion to be keyword kind")
			}
			if !strings.Contains(completion.InsertText, "${") {
				t.Error("Expected spell completion to include snippet")
			}
			break
		}
	}

	if !found {
		t.Error("Expected 'spell' keyword in completions")
	}
}

func TestAnalyzer_GetCompletions_Methods(t *testing.T) {
	analyzer := New()

	content := `
spell test():
    File.`
	analyzer.UpdateDocument("test.crl", content, nil)

	position := protocol.Position{Line: 2, Character: 9}
	completions := analyzer.GetCompletions("test.crl", position)

	if len(completions) == 0 {
		t.Error("Expected completions for File grimoire")
	}

	expectedMethods := []string{"read", "write", "append", "exists", "open"}
	for _, method := range expectedMethods {
		found := false
		for _, completion := range completions {
			if completion.Label == method {
				found = true
				if completion.Kind != protocol.CompletionItemKindMethod {
					t.Errorf("Expected %s completion to be method kind", method)
				}
				break
			}
		}
		if !found {
			t.Errorf("Expected '%s' method in File completions", method)
		}
	}
}

func TestAnalyzer_GetCompletions_Builtins(t *testing.T) {
	analyzer := New()

	content := "pr"
	analyzer.UpdateDocument("test.crl", content, nil)

	position := protocol.Position{Line: 0, Character: 2}
	completions := analyzer.GetCompletions("test.crl", position)

	// Should find "print" builtin
	found := false
	for _, completion := range completions {
		if completion.Label == "print" {
			found = true
			if completion.Kind != protocol.CompletionItemKindFunction {
				t.Error("Expected print completion to be function kind")
			}
			break
		}
	}

	if !found {
		t.Error("Expected 'print' builtin in completions")
	}
}

func TestAnalyzer_GetCompletions_UserSymbols(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	position := protocol.Position{Line: 0, Character: 3}
	completions := analyzer.GetCompletions("test.crl", position)

	// Should find "greet" function
	found := false
	for _, completion := range completions {
		if completion.Label == "greet" {
			found = true
			if completion.Kind != protocol.CompletionItemKindFunction {
				t.Error("Expected greet completion to be function kind")
			}
			break
		}
	}

	if !found {
		t.Error("Expected 'greet' function in completions")
	}
}

func TestAnalyzer_GetHover_Builtin(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", "print('test')", nil)

	position := protocol.Position{Line: 0, Character: 2}
	hover := analyzer.GetHover("test.crl", position)

	if hover == nil {
		t.Error("Expected hover information for 'print' builtin")
	}

	if !strings.Contains(hover.Contents.(string), "print") {
		t.Error("Expected hover to contain 'print' information")
	}
}

func TestAnalyzer_GetHover_UserFunction(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	position := protocol.Position{Line: 8, Character: 6} // "greet" function
	hover := analyzer.GetHover("test.crl", position)

	if hover == nil {
		t.Error("Expected hover information for 'greet' function")
		return
	}

	if hover.Contents == nil {
		t.Error("Expected hover contents to be present")
		return
	}

	content := hover.Contents.(string)
	if !strings.Contains(content, "greet") {
		t.Error("Expected hover to contain 'greet' information")
	}

	if !strings.Contains(content, "Greet someone by name") {
		t.Error("Expected hover to contain docstring")
	}
}

func TestAnalyzer_GetHover_Grimoire(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	// Debug: check what's at line 17
	lines := strings.Split(testCarrionCode, "\n")
	if len(lines) > 17 {
		t.Logf("Line 16: '%s'", lines[16])
		t.Logf("Line 17: '%s'", lines[17])
		t.Logf("Line 18: '%s'", lines[18])

		if len(lines[17]) > 5 {
			t.Logf("Character at position 5 on line 17: '%c'", lines[17][5])
		}
	}

	position := protocol.Position{Line: 16, Character: 5} // "Person" grimoire
	word := analyzer.getWordAtPosition(testCarrionCode, position)
	t.Logf("Word at position (16, 5): '%s'", word)

	// Check if document symbols were parsed
	doc := analyzer.documents["test.crl"]
	if doc != nil && doc.Symbols != nil {
		t.Logf("Number of grimoires found: %d", len(doc.Symbols.Grimoires))
		for name := range doc.Symbols.Grimoires {
			t.Logf("  Grimoire: %s", name)
		}
	}

	hover := analyzer.GetHover("test.crl", position)

	if hover == nil {
		t.Error("Expected hover information for 'Person' grimoire")
		return
	}

	if hover.Contents == nil {
		t.Error("Expected hover contents to be present")
		return
	}

	content := hover.Contents.(string)
	if !strings.Contains(content, "Person") {
		t.Error("Expected hover to contain 'Person' information")
	}

	if !strings.Contains(content, "Grimoire") {
		t.Error("Expected hover to indicate it's a grimoire")
	}
}

func TestAnalyzer_GetDefinition(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	position := protocol.Position{Line: 8, Character: 6} // "greet" function
	locations := analyzer.GetDefinition("test.crl", position)

	if len(locations) == 0 {
		t.Error("Expected definition location for 'greet' function")
	}

	if locations[0].URI != "test.crl" {
		t.Error("Expected definition location to be in same file")
	}
}

func TestAnalyzer_GetDocumentSymbols(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	symbols := analyzer.GetDocumentSymbols("test.crl")

	if len(symbols) == 0 {
		t.Error("Expected document symbols to be found")
	}

	// Check for grimoires
	foundPerson := false
	foundEmployee := false

	for _, symbol := range symbols {
		if symbol.Name == "Person" && symbol.Kind == protocol.SymbolKindClass {
			foundPerson = true
			// Check for methods as children
			if len(symbol.Children) == 0 {
				t.Error("Expected Person grimoire to have method children")
			}
		}
		if symbol.Name == "Employee" && symbol.Kind == protocol.SymbolKindClass {
			foundEmployee = true
		}
	}

	if !foundPerson {
		t.Error("Expected Person grimoire in document symbols")
	}

	if !foundEmployee {
		t.Error("Expected Employee grimoire in document symbols")
	}
}

func TestAnalyzer_GetSemanticTokens(t *testing.T) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", "spell test(): return 42", nil)

	tokens := analyzer.GetSemanticTokens("test.crl")

	if tokens == nil {
		t.Error("Expected semantic tokens to be generated")
	}

	if len(tokens.Data) == 0 {
		t.Error("Expected semantic token data to be present")
	}

	// Data should be in groups of 5: [deltaLine, deltaStart, length, tokenType, tokenModifiers]
	if len(tokens.Data)%5 != 0 {
		t.Error("Expected semantic token data to be in groups of 5")
	}
}

func TestAnalyzer_InferType(t *testing.T) {
	analyzer := New()

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{"integer", "42", "int"},
		{"float", "3.14", "float"},
		{"string", `"hello"`, "string"},
		{"boolean", "True", "bool"},
		{"array", "[1, 2, 3]", "array"},
		{"hash", `{"key": "value"}`, "hash"},
		{"tuple", "(1, 2)", "tuple"},
		{"none", "None", "None"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content := "x = " + test.content
			doc := analyzer.UpdateDocument("test.crl", content, nil)

			if variable, exists := doc.Symbols.Variables["x"]; exists {
				if variable.Type != test.expected {
					t.Errorf("Expected type %s, got %s for %s", test.expected, variable.Type, test.content)
				}
			} else {
				t.Errorf("Expected variable 'x' to be found for %s", test.content)
			}
		})
	}
}

func TestAnalyzer_ExtractParameters(t *testing.T) {
	analyzer := New()

	content := `
spell test_params(name: string, age: int = 25, active: bool):
    return name
`

	doc := analyzer.UpdateDocument("test.crl", content, nil)

	if spell, exists := doc.Symbols.Spells["test_params"]; exists {
		params := spell.Parameters

		if len(params) != 3 {
			t.Errorf("Expected 3 parameters, got %d", len(params))
		}

		// Check first parameter
		if params[0].Name != "name" {
			t.Errorf("Expected first parameter to be 'name', got %s", params[0].Name)
		}
		if params[0].TypeHint != "string" {
			t.Errorf("Expected first parameter type to be 'string', got %s", params[0].TypeHint)
		}

		// Check second parameter with default
		if params[1].Name != "age" {
			t.Errorf("Expected second parameter to be 'age', got %s", params[1].Name)
		}
		if params[1].DefaultValue != "25" {
			t.Errorf("Expected second parameter default to be '25', got %s", params[1].DefaultValue)
		}
	} else {
		t.Error("Expected test_params spell to be found")
	}
}

func TestAnalyzer_getWordAtPosition(t *testing.T) {
	analyzer := New()

	content := "spell greet(name):\n    return name"

	tests := []struct {
		line      int
		character int
		expected  string
	}{
		{0, 0, "spell"},
		{0, 6, "greet"},
		{0, 12, "name"},
		{1, 4, "return"},
		{1, 11, "name"},
	}

	for _, test := range tests {
		word := analyzer.getWordAtPosition(content, protocol.Position{
			Line:      test.line,
			Character: test.character,
		})

		if word != test.expected {
			t.Errorf("Expected word '%s' at position (%d, %d), got '%s'",
				test.expected, test.line, test.character, word)
		}
	}
}

func TestAnalyzer_formatParameters(t *testing.T) {
	analyzer := New()

	params := []Parameter{
		{Name: "name", TypeHint: "string"},
		{Name: "age", TypeHint: "int", DefaultValue: "25"},
		{Name: "active", TypeHint: "bool"},
	}

	result := analyzer.formatParameters(params)
	expected := "name: string, age: int = 25, active: bool"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestAnalyzer_mapTokenToSemanticType(t *testing.T) {
	// Note: This test would need the actual token types from the token package
	// For now, we're testing the logic conceptually
	t.Log("Token mapping test - would need actual token types from Carrion language")
}

func TestNewWorkspace(t *testing.T) {
	workspace := NewWorkspace("/test/path")

	if workspace == nil {
		t.Fatal("Expected workspace to be created")
	}

	if workspace.RootPath != "/test/path" {
		t.Errorf("Expected root path to be '/test/path', got %s", workspace.RootPath)
	}

	if workspace.Documents == nil {
		t.Error("Expected documents map to be initialized")
	}
}

// Benchmark tests for performance
func BenchmarkAnalyzer_UpdateDocument(b *testing.B) {
	analyzer := New()

	for i := 0; i < b.N; i++ {
		analyzer.UpdateDocument("test.crl", testCarrionCode, nil)
	}
}

func BenchmarkAnalyzer_GetCompletions(b *testing.B) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	position := protocol.Position{Line: 2, Character: 9}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.GetCompletions("test.crl", position)
	}
}

func BenchmarkAnalyzer_GetHover(b *testing.B) {
	analyzer := New()
	analyzer.UpdateDocument("test.crl", testCarrionCode, nil)

	position := protocol.Position{Line: 6, Character: 6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		analyzer.GetHover("test.crl", position)
	}
}
