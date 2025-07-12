//go:build ignore
// +build ignore

package main

import (
	"fmt"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func testKeyword(name, input string) {
	fmt.Printf("\n=== Testing %s ===\n", name)

	options := protocol.FormattingOptions{
		TabSize:      4,
		InsertSpaces: true,
	}

	formatter := analyzer.NewCarrionFormatter(options)
	edits, err := formatter.FormatDocument(input)

	if err != nil {
		fmt.Printf("❌ FAILED: %v\n", err)
		return
	}

	if len(edits) == 0 {
		fmt.Printf("✅ PARSED: No formatting changes needed\n")
		fmt.Printf("Input: %s\n", input)
		return
	}

	result := edits[0].NewText
	fmt.Printf("✅ FORMATTED:\n%s\n", result)
}

func main() {
	// Test spell functions
	testKeyword("Spell Functions", `spell test_function(param1, param2):
    result = param1 + param2
    return result`)

	// Test grimoires
	testKeyword("Grimoires", `grim TestClass:
    init(name):
        self.name = name
    spell get_name():
        return self.name`)

	// Test if/otherwise/else
	testKeyword("If/Otherwise/Else", `spell test_conditions(x):
    if x > 10:
        print("big")
    otherwise x == 5:
        print("medium")
    else:
        print("small")`)

	// Test match/case
	testKeyword("Match/Case", `spell test_match(value):
    match value:
        case 1:
            return "one"
        case 2:
            return "two"`)

	// Test autoclose
	testKeyword("Autoclose", `spell test_autoclose():
    autoclose open("file.txt") as f:
        content = f.read()
        return content`)

	// Test for loop
	testKeyword("For Loop", `spell test_for():
    for i in range(10):
        print(i)
        if i > 5:
            skip
        print("processed:", i)`)

	// Test while loop
	testKeyword("While Loop", `spell test_while():
    count = 0
    while count < 10:
        print(count)
        count = count + 1
        if count == 5:
            stop`)

	// Test attempt/ensnare/resolve
	testKeyword("Attempt/Ensnare/Resolve", `spell test_errors():
    attempt:
        risky_operation()
    ensnare:
        print("Error caught")
    resolve:
        cleanup()`)

	// Test main block
	testKeyword("Main Block", `main:
    print("Starting program")
    result = test_function(1, 2)
    print("Result:", result)`)

	// Test nested constructs
	testKeyword("Complex Nesting", `grim ComplexExample:
    spell complex_method():
        for i in range(5):
            if i % 2 == 0:
                attempt:
                    process(i)
                ensnare:
                    print("Error processing:", i)
                    skip
            otherwise i == 3:
                stop
        return "done"`)
}
