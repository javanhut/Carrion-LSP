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
	fmt.Println("Debug: Completion context detection")

	a := analyzer.New()

	// Test the exact string completion that was working
	code := `spell test():
    message = "hello"
    message.`

	a.UpdateDocument("test.crl", code, nil)

	// Test completion at the exact position
	position := protocol.Position{Line: 2, Character: 12}

	// Debug the prefix extraction
	lines := strings.Split(code, "\n")
	if position.Line < len(lines) {
		currentLine := lines[position.Line]
		prefix := ""
		if position.Character <= len(currentLine) {
			prefix = currentLine[:position.Character]
		}

		fmt.Printf("Current line: '%s'\n", currentLine)
		fmt.Printf("Prefix: '%s'\n", prefix)
		fmt.Printf("Has suffix '.': %v\n", strings.HasSuffix(prefix, "."))
	}

	completions := a.GetCompletions("test.crl", position)

	fmt.Printf("Found %d completions:\n", len(completions))
	for _, completion := range completions {
		fmt.Printf("  - %s (kind=%d): %s\n", completion.Label, completion.Kind, completion.Detail)
	}

	// Also test general completion (no dot)
	position2 := protocol.Position{Line: 1, Character: 13}
	completions2 := a.GetCompletions("test.crl", position2)

	fmt.Printf("\nGeneral completions (no dot) found %d:\n", len(completions2))
	for i, completion := range completions2 {
		if i < 5 { // Show first 5
			fmt.Printf("  - %s (kind=%d)\n", completion.Label, completion.Kind)
		}
	}
	if len(completions2) > 5 {
		fmt.Printf("  ... and %d more\n", len(completions2)-5)
	}
}
