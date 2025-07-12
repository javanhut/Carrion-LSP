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
	fmt.Println("Debug: Prefix extraction")

	a := analyzer.New()

	// Test different code scenarios
	testCases := []struct {
		name string
		code string
		line int
		char int
	}{
		{
			name: "After equals with p",
			code: "spell test():\n    result = p",
			line: 1,
			char: 14,
		},
		{
			name: "After equals no letter",
			code: "spell test():\n    result = ",
			line: 1,
			char: 13,
		},
		{
			name: "Beginning of line",
			code: "spell test():\n    p",
			line: 1,
			char: 5,
		},
	}

	for _, testCase := range testCases {
		fmt.Printf("\n=== %s ===\n", testCase.name)

		a.UpdateDocument("test.crl", testCase.code, nil)

		// Extract prefix manually like the function does
		lines := strings.Split(testCase.code, "\n")
		currentLine := lines[testCase.line]
		prefix := ""
		if testCase.char <= len(currentLine) {
			prefix = currentLine[:testCase.char]
		}

		fmt.Printf("Code: %q\n", testCase.code)
		fmt.Printf("Line %d: %q\n", testCase.line, currentLine)
		fmt.Printf("Prefix (char %d): %q\n", testCase.char, prefix)
		fmt.Printf("Has suffix '.': %v\n", strings.HasSuffix(prefix, "."))
		fmt.Printf("Has suffix '(': %v\n", strings.HasSuffix(prefix, "("))

		completions := a.GetCompletions("test.crl", protocol.Position{Line: testCase.line, Character: testCase.char})

		fmt.Printf("Completions found: %d\n", len(completions))

		// Check what built-ins should match this prefix
		builtins := a.GetBuiltins()
		lastToken := extractLastToken(prefix)
		fmt.Printf("Last token from prefix: %q\n", lastToken)

		matchingBuiltins := 0
		for name := range builtins {
			if strings.HasPrefix(name, lastToken) {
				matchingBuiltins++
			}
		}
		fmt.Printf("Built-ins that should match: %d\n", matchingBuiltins)
	}
}

func extractLastToken(text string) string {
	// Split by whitespace and operators to get the last token
	tokens := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '=' || r == '(' || r == ')' || r == ','
	})
	if len(tokens) > 0 {
		return tokens[len(tokens)-1]
	}
	return text
}
