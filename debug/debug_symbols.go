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
	doc := analyzer.UpdateDocument("test.crl", code, nil)

	fmt.Printf("Document parsed. Symbols found:\n")
	if doc.Symbols != nil {
		fmt.Printf("Variables: %d\n", len(doc.Symbols.Variables))
		for name, variable := range doc.Symbols.Variables {
			fmt.Printf("  - %s: %s\n", name, variable.Type)
		}
		
		fmt.Printf("Spells: %d\n", len(doc.Symbols.Spells))
		for name := range doc.Symbols.Spells {
			fmt.Printf("  - %s\n", name)
		}
		
		fmt.Printf("Grimoires: %d\n", len(doc.Symbols.Grimoires))
		for name := range doc.Symbols.Grimoires {
			fmt.Printf("  - %s\n", name)
		}
	} else {
		fmt.Printf("No symbols found!\n")
	}

	// Test completion right after the dot
	position := protocol.Position{
		Line:      2, // message.
		Character: 12, // Right after the dot
	}

	fmt.Printf("\nTesting completion at position %d:%d\n", position.Line, position.Character)
	fmt.Printf("Code:\n%s\n", code)

	completions := analyzer.GetCompletions("test.crl", position)

	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s (%s): %s\n", completion.Label, completion.Kind, completion.Detail)
	}
}