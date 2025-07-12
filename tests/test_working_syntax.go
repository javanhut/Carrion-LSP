//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

func testSyntax(name, input string) {
	fmt.Printf("\n=== Testing %s ===\n", name)

	l := lexer.New(input)
	p := parser.New(l)
	_ = p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("❌ Parsing errors:\n")
		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Printf("✅ Parsed successfully\n")
	}
}

func main() {
	// Test simple attempt without ensnare
	testSyntax("Simple attempt", `attempt:
    dangerous()`)

	// Test attempt with resolve only
	testSyntax("Attempt with resolve", `attempt:
    dangerous()
resolve:
    cleanup()`)

	// Test just ensnare by itself (to see what parser expects)
	testSyntax("Just ensnare", `ensnare:
    print("caught")`)

	// Test attempt with simple ensnare (no parentheses)
	testSyntax("Attempt with simple ensnare", `attempt:
    dangerous()
ensnare:
    print("caught")`)

	// Test what the parser might actually expect based on the logic
	testSyntax("Attempt with ensnare alias", `attempt:
    dangerous()
ensnare(e):
    print("caught")`)
}
