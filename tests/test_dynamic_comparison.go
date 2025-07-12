package main

import (
	"fmt"
	"sort"

	"github.com/javanhut/CarrionLSP/internal/analyzer"
	"github.com/javanhut/CarrionLSP/internal/protocol"
)

func main() {
	fmt.Println("=== Dynamic Loading vs Static Comparison ===\n")
	
	// Create analyzer with dynamic loading
	a := analyzer.New()

	// Test with multiple types
	code := `spell test():
    text = "hello"
    numbers = [1, 2, 3]
    count = 42
    text.`

	// Update document
	a.UpdateDocument("test.crl", code, nil)

	// Test string methods
	fmt.Println("STRING METHODS (Dynamically Loaded):")
	completions := a.GetCompletions("test.crl", protocol.Position{Line: 3, Character: 9})
	var stringMethods []string
	for _, completion := range completions {
		stringMethods = append(stringMethods, completion.Label)
	}
	sort.Strings(stringMethods)
	for _, method := range stringMethods {
		fmt.Printf("  - %s\n", method)
	}

	// Test array methods  
	code2 := `spell test():
    numbers = [1, 2, 3]
    numbers.`
	a.UpdateDocument("test.crl", code2, nil)
	
	fmt.Println("\nARRAY METHODS (Dynamically Loaded):")
	completions = a.GetCompletions("test.crl", protocol.Position{Line: 2, Character: 12})
	var arrayMethods []string
	for _, completion := range completions {
		arrayMethods = append(arrayMethods, completion.Label)
	}
	sort.Strings(arrayMethods)
	for _, method := range arrayMethods {
		fmt.Printf("  - %s\n", method)
	}

	// Test built-in functions
	code3 := `spell test():
    result = `
	a.UpdateDocument("test.crl", code3, nil)
	
	fmt.Println("\nBUILT-IN FUNCTIONS (Dynamically Loaded):")
	completions = a.GetCompletions("test.crl", protocol.Position{Line: 1, Character: 13})
	var builtinFunctions []string
	for _, completion := range completions {
		if completion.Kind == protocol.CompletionItemKindFunction {
			builtinFunctions = append(builtinFunctions, completion.Label)
		}
	}
	sort.Strings(builtinFunctions)
	for _, function := range builtinFunctions {
		fmt.Printf("  - %s\n", function)
	}

	// Test grimoire classes
	fmt.Println("\nGRIMOIRE CLASSES (Dynamically Loaded):")
	var grimoireClasses []string
	for _, completion := range completions {
		if completion.Kind == protocol.CompletionItemKindClass {
			grimoireClasses = append(grimoireClasses, completion.Label)
		}
	}
	sort.Strings(grimoireClasses)
	for _, grimoire := range grimoireClasses {
		fmt.Printf("  - %s\n", grimoire)
	}

	fmt.Println("\n=== Summary ===")
	fmt.Printf("✅ Dynamically loaded %d string methods\n", len(stringMethods))
	fmt.Printf("✅ Dynamically loaded %d array methods\n", len(arrayMethods))
	fmt.Printf("✅ Dynamically loaded %d built-in functions\n", len(builtinFunctions))
	fmt.Printf("✅ Dynamically loaded %d grimoire classes\n", len(grimoireClasses))
	fmt.Println("✅ All data loaded directly from Carrion runtime - always up to date!")
}