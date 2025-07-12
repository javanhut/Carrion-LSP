//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	fmt.Println("Debug: String method completion with dynamic loading")

	// Create analyzer
	a := analyzer.New()

	// Simple string test
	code := `spell test():
    message = "hello"
    message.`

	// Update document
	doc := a.UpdateDocument("test.crl", code, nil)

	fmt.Printf("Variables found:\n")
	for name, variable := range doc.Symbols.Variables {
		fmt.Printf("  - %s: %s\n", name, variable.Type)
	}

	// Test completion
	position := protocol.Position{Line: 2, Character: 12}
	completions := a.GetCompletions("test.crl", position)

	fmt.Printf("\nCompletions found: %d\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s (kind=%d): %s\n", completion.Label, completion.Kind, completion.Detail)
	}
}
