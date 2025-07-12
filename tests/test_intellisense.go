//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	analyzer := analyzer.New()

	// Test Carrion code with various types of objects
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
    
    # Test method calls
    greeting = person.greet()
    info = person.get_info()
    
    return "done"

main:
    test_function()`

	// Update the document
	analyzer.UpdateDocument("test.crl", code, nil)

	// Test different completion scenarios
	testCompletions := []struct {
		name     string
		line     int
		char     int
		expected []string
	}{
		{
			name:     "Person instance methods",
			line:     22, // person.
			char:     11,
			expected: []string{"greet", "get_info"},
		},
		{
			name:     "String methods",
			line:     23, // message.
			char:     12,
			expected: []string{"length", "lower", "upper", "split", "contains"},
		},
		{
			name:     "Array methods",
			line:     24, // numbers.
			char:     12,
			expected: []string{"length", "append", "get", "contains", "sort"},
		},
		{
			name:     "Integer methods",
			line:     25, // count.
			char:     10,
			expected: []string{"abs", "is_even", "is_odd", "to_hex"},
		},
		{
			name:     "Float methods",
			line:     26, // ratio.
			char:     10,
			expected: []string{"round", "floor", "ceil", "sqrt", "abs"},
		},
		{
			name:     "Boolean methods",
			line:     27, // is_active.
			char:     14,
			expected: []string{"negate", "and_with", "or_with", "to_int"},
		},
		{
			name:     "Built-in grimoires",
			line:     0, // Global scope
			char:     0,
			expected: []string{"File", "OS", "Time", "String", "Array"},
		},
	}

	for _, test := range testCompletions {
		fmt.Printf("\n=== Testing %s ===\n", test.name)

		position := protocol.Position{
			Line:      test.line,
			Character: test.char,
		}

		completions := analyzer.GetCompletions("test.crl", position)

		fmt.Printf("Found %d completions:\n", len(completions))

		found := make(map[string]bool)
		for _, completion := range completions {
			found[completion.Label] = true
			fmt.Printf("  - %s (%s): %s\n", completion.Label, completion.Kind, completion.Detail)
		}

		// Check if expected completions are present
		for _, expected := range test.expected {
			if found[expected] {
				fmt.Printf("✅ Found expected completion: %s\n", expected)
			} else {
				fmt.Printf("❌ Missing expected completion: %s\n", expected)
			}
		}
	}

	// Test hover information
	fmt.Printf("\n=== Testing Hover Information ===\n")
	hoverPosition := protocol.Position{Line: 22, Character: 5} // "person" variable
	hover := analyzer.GetHover("test.crl", hoverPosition)

	if hover != nil {
		fmt.Printf("Hover info: %s\n", hover.Contents)
	} else {
		fmt.Printf("No hover information found\n")
	}
}
