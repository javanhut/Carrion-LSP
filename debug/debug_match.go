package main

import (
	"fmt"

	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

func main() {
	input := `spell test_match(value):
    match value:
        case 1:
            return "one"`
	
	fmt.Printf("Testing match parsing:\n%s\n\n", input)
	
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	
	fmt.Printf("Parse errors: %d\n", len(p.Errors()))
	for _, err := range p.Errors() {
		fmt.Printf("  - %s\n", err)
	}
	
	if len(p.Errors()) == 0 {
		fmt.Printf("Parsed successfully!\n")
		fmt.Printf("AST: %s\n", program.String())
	}
}