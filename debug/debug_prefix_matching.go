package main

import (
	"fmt"
	"strings"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	fmt.Println("Debug: Prefix matching for completions")
	
	a := analyzer.New()

	// Check what built-ins start with common prefixes
	builtins := a.GetBuiltins()
	fmt.Printf("Built-ins starting with 'p': ")
	pCount := 0
	for name := range builtins {
		if strings.HasPrefix(name, "p") {
			fmt.Printf("%s ", name)
			pCount++
		}
	}
	fmt.Printf("(%d total)\n", pCount)
	
	fmt.Printf("Built-ins starting with 'l': ")
	lCount := 0
	for name := range builtins {
		if strings.HasPrefix(name, "l") {
			fmt.Printf("%s ", name)
			lCount++
		}
	}
	fmt.Printf("(%d total)\n", lCount)

	// Test with empty prefix
	code := "spell test():\n    result = "
	a.UpdateDocument("test.crl", code, nil)
	completions := a.GetCompletions("test.crl", protocol.Position{Line: 1, Character: 13})

	fmt.Printf("\nEmpty prefix completions: %d\n", len(completions))
	
	var builtinCompletions []string
	var grimoireCompletions []string
	
	for _, completion := range completions {
		if completion.Kind == protocol.CompletionItemKindFunction {
			builtinCompletions = append(builtinCompletions, completion.Label)
		} else if completion.Kind == protocol.CompletionItemKindClass {
			grimoireCompletions = append(grimoireCompletions, completion.Label)
		}
	}
	
	fmt.Printf("Built-in functions: %d - %v\n", len(builtinCompletions), builtinCompletions[:min(5, len(builtinCompletions))])
	fmt.Printf("Grimoire classes: %d - %v\n", len(grimoireCompletions), grimoireCompletions[:min(5, len(grimoireCompletions))])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}