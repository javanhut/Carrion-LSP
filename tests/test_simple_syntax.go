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
	// Test minimal match
	testSyntax("Simple Match", `match x:
    case 1:
        print("one")`)

	// Test minimal attempt
	testSyntax("Simple Attempt", `attempt:
    dangerous()
ensnare Error:
    print("failed")`)

	// Test what we know works
	testSyntax("Known Working If", `if x > 5:
    print("big")`)
}