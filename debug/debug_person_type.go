package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	analyzer := analyzer.New()

	// Test code with Person grimoire and variable
	code := `grim Person:
    init(name: string):
        self.name = name

    spell greet():
        return "Hello, I'm " + self.name

spell test_function():
    person = Person("Alice")
    person.`

	// Update the document
	doc := analyzer.UpdateDocument("test.crl", code, nil)

	fmt.Printf("Document parsed. Symbols found:\n")
	if doc.Symbols != nil {
		fmt.Printf("Variables: %d\n", len(doc.Symbols.Variables))
		for name, variable := range doc.Symbols.Variables {
			fmt.Printf("  - %s: %s\n", name, variable.Type)
		}
		
		fmt.Printf("Grimoires: %d\n", len(doc.Symbols.Grimoires))
		for name, grimoire := range doc.Symbols.Grimoires {
			fmt.Printf("  - %s (spells: %d)\n", name, len(grimoire.Spells))
			for spellName := range grimoire.Spells {
				fmt.Printf("    * %s\n", spellName)
			}
		}
	}

	// Test completion right after the dot
	position := protocol.Position{
		Line:      9, // person.
		Character: 11, // Right after the dot
	}

	fmt.Printf("\nTesting completion at position %d:%d\n", position.Line, position.Character)
	completions := analyzer.GetCompletions("test.crl", position)

	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s: %s\n", completion.Label, completion.Detail)
	}
}