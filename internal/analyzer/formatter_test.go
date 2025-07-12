package analyzer

import (
	"strings"
	"testing"

	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func TestCarrionFormatter_New(t *testing.T) {
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	
	if formatter == nil {
		t.Fatal("Expected formatter to be created")
	}
	
	if formatter.tabSize != 4 {
		t.Errorf("Expected tab size 4, got %d", formatter.tabSize)
	}
	
	if !formatter.insertSpaces {
		t.Error("Expected insert spaces to be true")
	}
}

func TestCarrionFormatter_FormatSimpleFunction(t *testing.T) {
	input := `spell greet(name):
    return "Hello, " + name`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	if len(edits) != 1 {
		t.Fatalf("Expected 1 edit, got %d", len(edits))
	}
	
	result := edits[0].NewText
	
	// Check basic formatting
	if !strings.Contains(result, "spell greet(name):") {
		t.Error("Expected function declaration")
	}
	
	if !strings.Contains(result, `"Hello, " + name`) {
		t.Error("Expected spacing around operators")
	}
}

func TestCarrionFormatter_FormatFunctionWithDocstring(t *testing.T) {
	input := `spell greet(name:string):
    """Greet someone by name"""
    return "Hello, " + name`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	if len(edits) != 1 {
		t.Fatalf("Expected 1 edit, got %d", len(edits))
	}
	
	result := edits[0].NewText
	// Check that docstring and proper spacing are present
	if !strings.Contains(result, `"""Greet someone by name"""`) {
		t.Error("Expected docstring to be preserved")
	}
	
	if !strings.Contains(result, "name: string") {
		t.Error("Expected parameter type spacing")
	}
}

func TestCarrionFormatter_FormatGrimoire(t *testing.T) {
	input := `grim Person:
    init(name:string):
        self.name=name
    spell greet():
        return "Hello, I'm "+self.name`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check basic formatting
	if !strings.Contains(result, "grim Person:") {
		t.Error("Expected grimoire declaration")
	}
	
	if !strings.Contains(result, "init(") {
		t.Error("Expected init method")
	}
	
	// Check spacing around operators
	if !strings.Contains(result, "self.name = name") {
		t.Error("Expected spacing around assignment operator")
	}
}

func TestCarrionFormatter_FormatIfStatement(t *testing.T) {
	input := `if x>5:
print("big")
otherwise x==3:
print("three")
else:
print("small")`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check operator spacing
	if !strings.Contains(result, "x > 5") {
		t.Error("Expected spacing around comparison operator")
	}
	
	if !strings.Contains(result, "x == 3") {
		t.Error("Expected spacing around equality operator")
	}
	
	// Check indentation
	if !strings.Contains(result, "    print(") {
		t.Error("Expected print statements to be indented")
	}
}

func TestCarrionFormatter_FormatForLoop(t *testing.T) {
	input := `for i in range(5):
print(i)
if i>2:
break`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check nested indentation
	if !strings.Contains(result, "    print(i)") {
		t.Error("Expected first level indentation")
	}
	
	if !strings.Contains(result, "        break") {
		t.Error("Expected second level indentation")
	}
}

func TestCarrionFormatter_FormatAttemptBlock(t *testing.T) {
	input := `attempt:
x=risky_operation()
ensnare ValueError as e:
print("Error:",e)
resolve:
cleanup()`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check that all blocks are properly indented
	if !strings.Contains(result, "    x = risky_operation()") {
		t.Error("Expected attempt block content to be indented")
	}
	
	if !strings.Contains(result, "    cleanup()") {
		t.Error("Expected resolve block content to be indented")
	}
}

func TestCarrionFormatter_FormatArrayAndHash(t *testing.T) {
	input := `spell test():
    arr=[1,2,3,4,5]
    hash={"key":"value","other":42}`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check basic formatting
	if !strings.Contains(result, "spell test():") {
		t.Error("Expected function declaration")
	}
	
	// Check spacing in collections
	if !strings.Contains(result, "[1, 2, 3, 4, 5]") {
		t.Error("Expected spacing in array literal")
	}
	
	if !strings.Contains(result, `{"key": "value", "other": 42}`) {
		t.Error("Expected spacing in hash literal")
	}
}

func TestCarrionFormatter_FormatCallExpression(t *testing.T) {
	input := `spell test():
    result=function_name(arg1,arg2,arg3)
    obj.method(x,y)`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check basic formatting
	if !strings.Contains(result, "spell test():") {
		t.Error("Expected function declaration")
	}
	
	// Check spacing in function calls
	if !strings.Contains(result, "(arg1, arg2, arg3)") {
		t.Error("Expected spacing in function call arguments")
	}
	
	if !strings.Contains(result, "obj.method(x, y)") {
		t.Error("Expected spacing in method call arguments")
	}
}

func TestCarrionFormatter_FormatMainBlock(t *testing.T) {
	input := `main:
print("Starting program")
result=calculate(5,10)
print("Result:",result)`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check main block indentation
	if !strings.Contains(result, "    print(") {
		t.Error("Expected main block content to be indented")
	}
	
	if !strings.Contains(result, "result = calculate(5, 10)") {
		t.Error("Expected proper spacing in assignment and function call")
	}
}

func TestCarrionFormatter_FormatWithTabs(t *testing.T) {
	input := `spell test():
    return 42`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: false, // Use tabs
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check that tabs are used instead of spaces
	if !strings.Contains(result, "\treturn 42") {
		t.Error("Expected tab indentation when InsertSpaces is false")
	}
}

func TestCarrionFormatter_FormatComplexProgram(t *testing.T) {
	input := `import "os" as OS
grim Calculator:
"""A simple calculator grimoire"""
init():
self.history=[]
spell add(a,b):
result=a+b
self.history.append({"op":"add","args":[a,b],"result":result})
return result
spell get_history():
return self.history
main:
calc=Calculator()
result=calc.add(5,3)
print("5 + 3 =",result)
history=calc.get_history()
print("History:",history)`

	options := protocol.FormattingOptions{
		TabSize:                4,
		InsertSpaces:           true,
		TrimTrailingWhitespace: true,
		InsertFinalNewline:     true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		t.Fatalf("Formatting failed: %v", err)
	}
	
	result := edits[0].NewText
	
	// Check various formatting aspects
	tests := []string{
		`import "os" as OS`,                    // Import formatting
		`grim Calculator:`,                     // Grimoire declaration
		`"""A simple calculator grimoire"""`,  // Docstring
		`    init():`,                          // Method indentation
		`        self.history = []`,            // Assignment spacing
		`    spell add(a, b):`,                 // Method signature
		`        result = a + b`,               // Operator spacing
		`        return result`,                // Return statement
		`main:`,                               // Main block
		`    calc = Calculator()`,             // Constructor call
		`    result = calc.add(5, 3)`,         // Method call
	}
	
	for _, expected := range tests {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected to find: %s\nIn result:\n%s", expected, result)
		}
	}
	
	// Check that result ends with newline
	if !strings.HasSuffix(result, "\n") {
		t.Error("Expected result to end with newline")
	}
}

func TestCarrionFormatter_ParseError(t *testing.T) {
	input := `spell broken syntax here...`
	
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err == nil {
		t.Error("Expected formatting to fail with parse error")
	}
	
	if edits != nil {
		t.Error("Expected no edits when parsing fails")
	}
}

func TestAnalyzer_FormatDocument_Integration(t *testing.T) {
	analyzer := New()
	
	input := `spell greet(name):
    return "Hello, "+name`
	
	// Update document first
	analyzer.UpdateDocument("test.crl", input, nil)
	
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	edits := analyzer.FormatDocument("test.crl", options)
	
	if len(edits) == 0 {
		t.Error("Expected formatting edits to be returned")
	}
	
	if len(edits) > 0 {
		result := edits[0].NewText
		
		// Check basic formatting
		if !strings.Contains(result, "spell greet(name):") {
			t.Error("Expected function declaration")
		}
		
		if !strings.Contains(result, `"Hello, " + name`) {
			t.Error("Expected spacing around operators")
		}
	}
}

func TestAnalyzer_FormatDocument_NonExistentFile(t *testing.T) {
	analyzer := New()
	
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	edits := analyzer.FormatDocument("nonexistent.crl", options)
	
	if edits != nil {
		t.Error("Expected nil edits for non-existent document")
	}
}

// Benchmark formatting performance
func BenchmarkCarrionFormatter_FormatDocument(b *testing.B) {
	input := `
import "os" as OS

grim Calculator:
    """A comprehensive calculator grimoire"""
    
    init():
        self.history = []
        self.precision = 2
    
    spell add(a, b):
        """Add two numbers"""
        result = a + b
        self.history.append({"op": "add", "args": [a, b], "result": result})
        return result
    
    spell multiply(a, b):
        """Multiply two numbers"""
        result = a * b
        self.history.append({"op": "multiply", "args": [a, b], "result": result})
        return result
    
    spell get_history():
        """Get calculation history"""
        return self.history

main:
    calc = Calculator()
    
    for i in range(10):
        result = calc.add(i, i * 2)
        print(f"Iteration {i}: result = {result}")
    
    history = calc.get_history()
    print(f"Total calculations: {len(history)}")
`

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := NewCarrionFormatter(options)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := formatter.FormatDocument(input)
		if err != nil {
			b.Fatalf("Formatting failed: %v", err)
		}
	}
}