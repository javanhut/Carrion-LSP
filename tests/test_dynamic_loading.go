//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	fmt.Println("Testing dynamic loading system...")

	// Create analyzer with dynamic loading
	a := analyzer.New()

	// Test code with string variable
	code := `spell test():
    message = "hello world"
    message.`

	// Update the document
	doc := a.UpdateDocument("test.crl", code, nil)

	fmt.Printf("Document parsed with dynamic loading.\n")
	if doc.Symbols != nil {
		fmt.Printf("Variables found: %d\n", len(doc.Symbols.Variables))
		for name, variable := range doc.Symbols.Variables {
			fmt.Printf("  - %s: %s\n", name, variable.Type)
		}
	}

	// Test completion
	position := protocol.Position{Line: 2, Character: 12}
	completions := a.GetCompletions("test.crl", position)

	fmt.Printf("\nDynamic method completions for string variable:\n")
	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s: %s\n", completion.Label, completion.Detail)
	}

	// Test refresh
	fmt.Printf("\nTesting dynamic refresh...\n")
	a.RefreshDynamicData()

	// Test again after refresh
	completions = a.GetCompletions("test.crl", position)
	fmt.Printf("After refresh - Found %d completions\n", len(completions))

	fmt.Printf("\nDynamic loading test completed!\n")
}
