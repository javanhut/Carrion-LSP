//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	// Read the unformatted input
	input, err := ioutil.ReadFile("test_format_input.crl")
	if err != nil {
		log.Fatal(err)
	}

	// Read the expected output
	expected, err := ioutil.ReadFile("correct_format.crl")
	if err != nil {
		log.Fatal(err)
	}

	// Create formatter
	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}

	formatter := analyzer.NewCarrionFormatter(options)

	// Format the input
	edits, err := formatter.FormatDocument(string(input))
	if err != nil {
		log.Fatalf("Formatting failed: %v", err)
	}

	if len(edits) == 0 {
		fmt.Println("No formatting changes needed")
		return
	}

	result := edits[0].NewText

	fmt.Println("=== INPUT ===")
	fmt.Println(string(input))
	fmt.Println("\n=== FORMATTER OUTPUT ===")
	fmt.Println(result)
	fmt.Println("\n=== EXPECTED ===")
	fmt.Println(string(expected))

	fmt.Println("\n=== COMPARISON ===")
	if result == string(expected) {
		fmt.Println("✅ PERFECT MATCH! Formatter produces exactly the expected output.")
	} else {
		fmt.Println("❌ MISMATCH - Formatter output differs from expected.")
		fmt.Printf("Expected length: %d, Got length: %d\n", len(expected), len(result))

		// Show ending characters in detail
		fmt.Printf("Expected ends with: %q\n", string(expected)[len(expected)-10:])
		fmt.Printf("Result ends with: %q\n", result[len(result)-10:])
	}
}
