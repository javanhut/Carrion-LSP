//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"strings"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	analyzer := analyzer.New()

	// Test code with actual positions where we want to test completion
	code := `grim Person:
    init(name: string):
        self.name = name
        self.age = 25
        self.items = []

    spell greet():
        return "Hello, I'm " + self.name

    spell get_info():
        return self.name + " is " + str(self.age)

spell test_function():
    # Create instances and variables
    person = Person("Alice")
    message = "Hello World"
    numbers = [1, 2, 3, 4, 5]
    count = 42
    ratio = 3.14
    is_active = True
    
    # Test method calls - these lines would have dots for testing
    # person.
    # message.
    # numbers.
    # count.
    # ratio.
    # is_active.
    
    return "done"

main:
    test_function()`

	// Update the document
	analyzer.UpdateDocument("test.crl", code, nil)

	// Print line by line with numbers for reference
	fmt.Printf("Code with line numbers:\n")
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		fmt.Printf("Line %2d: %s\n", i, line)
	}

	// Test completion after adding a dot to person variable (line 15)
	modifiedCode := strings.ReplaceAll(code, "person = Person(\"Alice\")", "person = Person(\"Alice\")\n    person.")
	analyzer.UpdateDocument("test.crl", modifiedCode, nil)

	fmt.Printf("\n=== Testing person. completion ===\n")
	completions := analyzer.GetCompletions("test.crl", protocol.Position{Line: 16, Character: 11})
	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s\n", completion.Label)
	}
}
