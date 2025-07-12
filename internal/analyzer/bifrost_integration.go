package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/javanhut/TheCarrionLanguage/src/evaluator"
	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

// BifrostIntegration provides integration with the Bifrost package manager
type BifrostIntegration struct {
	analyzer      *Analyzer
	packagePaths  []string
	loadedPackages map[string]bool
}

// NewBifrostIntegration creates a new Bifrost integration
func NewBifrostIntegration(analyzer *Analyzer) *BifrostIntegration {
	return &BifrostIntegration{
		analyzer:       analyzer,
		packagePaths:   getCarrionPackagePaths(),
		loadedPackages: make(map[string]bool),
	}
}

// getCarrionPackagePaths returns the standard Carrion package search paths
func getCarrionPackagePaths() []string {
	var paths []string
	
	// Current directory's carrion_modules
	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, "carrion_modules"))
	}
	
	// User packages directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".carrion", "packages"))
	}
	
	// Global packages directory
	paths = append(paths, "/usr/local/share/carrion/lib")
	
	return paths
}

// DiscoverAvailablePackages scans for available packages in the search paths
func (bi *BifrostIntegration) DiscoverAvailablePackages() map[string]string {
	packages := make(map[string]string)
	
	for _, searchPath := range bi.packagePaths {
		if entries, err := os.ReadDir(searchPath); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					packageName := entry.Name()
					packagePath := filepath.Join(searchPath, packageName)
					packages[packageName] = packagePath
				}
			}
		}
	}
	
	return packages
}

// LoadPackage loads a specific package by name
func (bi *BifrostIntegration) LoadPackage(packageName string) error {
	// Skip if already loaded
	if bi.loadedPackages[packageName] {
		return nil
	}
	
	// Find the package
	packagePath := bi.findPackage(packageName)
	if packagePath == "" {
		return fmt.Errorf("package not found: %s", packageName)
	}
	
	// Load the package's main file
	mainFile := filepath.Join(packagePath, "src", "main.crl")
	if _, err := os.Stat(mainFile); os.IsNotExist(err) {
		// Try alternative locations
		alternatives := []string{
			filepath.Join(packagePath, "main.crl"),
			filepath.Join(packagePath, packageName+".crl"),
		}
		
		found := false
		for _, alt := range alternatives {
			if _, err := os.Stat(alt); err == nil {
				mainFile = alt
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("no main file found for package: %s", packageName)
		}
	}
	
	// Read and parse the file
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("failed to read package file: %v", err)
	}
	
	// Parse the package
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	
	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in package %s: %v", packageName, p.Errors())
	}
	
	// Evaluate in the dynamic loader's environment
	evaluator.Eval(program, bi.analyzer.dynamicLoader.env, nil)
	
	// Refresh the analyzer's data
	bi.analyzer.RefreshDynamicData()
	
	// Mark as loaded
	bi.loadedPackages[packageName] = true
	
	return nil
}

// findPackage searches for a package in the search paths
func (bi *BifrostIntegration) findPackage(packageName string) string {
	for _, searchPath := range bi.packagePaths {
		packagePath := filepath.Join(searchPath, packageName)
		if info, err := os.Stat(packagePath); err == nil && info.IsDir() {
			return packagePath
		}
	}
	return ""
}

// LoadPackageFromImport loads a package based on an import statement
func (bi *BifrostIntegration) LoadPackageFromImport(importPath string) error {
	// Parse the import path
	// Examples: "json-utils", "http-client/request", "./local-module"
	
	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		// Relative import - handle differently
		return bi.loadRelativePackage(importPath)
	}
	
	// Extract package name from path
	parts := strings.Split(importPath, "/")
	packageName := parts[0]
	
	return bi.LoadPackage(packageName)
}

// loadRelativePackage loads a package from a relative path
func (bi *BifrostIntegration) loadRelativePackage(relativePath string) error {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	
	// Resolve the relative path
	absolutePath := filepath.Join(cwd, relativePath)
	
	// Check if it's a .crl file
	if strings.HasSuffix(absolutePath, ".crl") {
		return bi.loadCarrionFile(absolutePath)
	}
	
	// Check if it's a directory with a main file
	if info, err := os.Stat(absolutePath); err == nil && info.IsDir() {
		mainFile := filepath.Join(absolutePath, "main.crl")
		if _, err := os.Stat(mainFile); err == nil {
			return bi.loadCarrionFile(mainFile)
		}
	}
	
	return fmt.Errorf("relative package not found: %s", relativePath)
}

// loadCarrionFile loads and evaluates a single .crl file
func (bi *BifrostIntegration) loadCarrionFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}
	
	l := lexer.New(string(content))
	p := parser.New(l)
	program := p.ParseProgram()
	
	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors in file %s: %v", filePath, p.Errors())
	}
	
	// Evaluate in the dynamic loader's environment
	evaluator.Eval(program, bi.analyzer.dynamicLoader.env, nil)
	
	// Refresh the analyzer's data
	bi.analyzer.RefreshDynamicData()
	
	return nil
}

// GetPackageCompletions returns completion suggestions for package names
func (bi *BifrostIntegration) GetPackageCompletions() []string {
	packages := bi.DiscoverAvailablePackages()
	var completions []string
	
	for packageName := range packages {
		completions = append(completions, packageName)
	}
	
	return completions
}

// AutoLoadImports scans a document for import statements and loads the packages
func (bi *BifrostIntegration) AutoLoadImports(doc *Document) error {
	if doc.Symbols == nil {
		return nil
	}
	
	// Load packages from import statements
	for _, importSym := range doc.Symbols.Imports {
		if importSym.Path != "" {
			err := bi.LoadPackageFromImport(importSym.Path)
			if err != nil {
				// Log the error but don't fail completely
				fmt.Printf("Warning: Failed to load import %s: %v\n", importSym.Path, err)
			}
		}
	}
	
	return nil
}