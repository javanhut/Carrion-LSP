//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

func testParsing(name, input string) {
	fmt.Printf("\n=== Testing %s ===\n", name)
	fmt.Printf("Input:\n%s\n\n", input)

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Printf("❌ Parsing errors:\n")
		for _, err := range p.Errors() {
			fmt.Printf("  - %s\n", err)
		}
	} else {
		fmt.Printf("✅ Parsed successfully\n")
		fmt.Printf("AST: %s\n", program.String())
	}
}

func main() {
	// Test match/case with proper syntax
	testParsing("Match/Case", `spell test_match(value):
    match value:
        case 1:
            return "one"
        case 2:
            return "two"
        default:
            return "other"`)

	// Test attempt/ensnare/resolve with proper syntax
	testParsing("Attempt/Ensnare/Resolve", `spell test_errors():
    attempt:
        risky_operation()
    ensnare ValueError as e:
        print("Error:", e)
    resolve:
        cleanup()`)
}
