package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	input := `spell test_match(value):
    match value:
        case 1:
            return "one"`
	
	fmt.Printf("Testing match formatting:\n%s\n\n", input)
	
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}
	
	formatter := analyzer.NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)
	
	if err != nil {
		fmt.Printf("❌ Formatting failed: %v\n", err)
		return
	}
	
	if len(edits) == 0 {
		fmt.Printf("✅ No formatting changes needed\n")
		fmt.Printf("Original: %s\n", input)
		return
	}
	
	result := edits[0].NewText
	fmt.Printf("✅ Formatted successfully:\n%s\n", result)
}