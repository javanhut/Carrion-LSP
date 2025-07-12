//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	fmt.Println("Debug: General completion")

	a := analyzer.New()

	// Test general completion
	code := "spell test():\n    result = p"
	a.UpdateDocument("test.crl", code, nil)

	// Check what built-ins are loaded
	fmt.Printf("Built-ins loaded in analyzer: %d\n", len(a.GetBuiltins()))
	for name := range a.GetBuiltins() {
		fmt.Printf("  - %s\n", name)
		break // Just show first one as example
	}

	// Check what grimoires are loaded
	fmt.Printf("Grimoires loaded in analyzer: %d\n", len(a.GetGrimoires()))
	for name := range a.GetGrimoires() {
		fmt.Printf("  - %s\n", name)
		break // Just show first one as example
	}

	completions := a.GetCompletions("test.crl", protocol.Position{Line: 1, Character: 14})

	fmt.Printf("General completions found: %d\n", len(completions))
	for i, completion := range completions {
		if i < 10 { // Show first 10
			fmt.Printf("  - %s (kind=%d): %s\n", completion.Label, completion.Kind, completion.Detail)
		}
	}
	if len(completions) > 10 {
		fmt.Printf("  ... and %d more\n", len(completions)-10)
	}
}
