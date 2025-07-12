package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	analyzer := analyzer.New()

	// Simple test code
	code := `spell test():
    message = "hello"
    message.`

	// Update the document
	analyzer.UpdateDocument("test.crl", code, nil)

	// Test completion right after the dot
	position := protocol.Position{
		Line:      2, // message.
		Character: 12, // Right after the dot
	}

	fmt.Printf("Testing completion at position %d:%d\n", position.Line, position.Character)
	fmt.Printf("Code:\n%s\n", code)

	completions := analyzer.GetCompletions("test.crl", position)

	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s (%s): %s\n", completion.Label, completion.Kind, completion.Detail)
	}
}