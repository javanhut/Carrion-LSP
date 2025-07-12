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
	fmt.Printf("Input: %s\n", input)

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
	// Test different ensnare syntaxes
	testSyntax("Ensnare with colon", `ensnare Error:
    print("failed")`)

	testSyntax("Ensnare without as", `ensnare(Error):
    print("failed")`)

	testSyntax("Ensnare with as", `ensnare(Error) as e:
    print("failed")`)

	// Test complete attempt block
	testSyntax("Complete attempt", `attempt:
    dangerous()
ensnare(Error) as e:
    print("failed")`)
}
