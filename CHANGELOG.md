# Changelog

All notable changes to the Carrion Language Server Protocol will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.1] - 2024-01-15

### Fixed

#### Built-in Module Import Resolution
- **Fixed built-in module import warnings**: Eliminated "package not found" warnings for built-in modules (`file`, `os`, `time`, `http`)
- **Enhanced Bifrost integration**: Added `isBuiltinModule()` function to properly distinguish between built-in and external packages
- **Improved error handling**: Better distinction between runtime modules and Bifrost packages

#### Test Suite Reliability
- **Corrected hover test positioning**: Fixed `TestAnalyzer_GetHover_Grimoire` test to correctly identify grimoire definitions at proper line positions (Line 16 instead of Line 26)
- **Enhanced formatter syntax support**: Updated formatter tests to use correct Carrion syntax for `attempt`/`ensnare`/`resolve` blocks
- **Fixed array formatting logic**: Updated formatter tests to align with multi-line array formatting rules (3+ elements)
- **Improved test indentation**: Corrected all test cases to use proper Carrion indentation patterns

### Improved

#### Test Coverage & Quality
- **100% test pass rate**: Achieved complete test suite reliability across all components
- **Enhanced test coverage**: Added proper Carrion syntax examples in formatter tests
- **Improved test reliability**: Eliminated flaky tests with better position calculations
- **Streamlined test execution**: All tests now pass consistently

#### Code Quality
- **Code formatting**: Automatic code formatting and lint compliance
- **Documentation accuracy**: Updated syntax examples to match actual Carrion grammar
- **Error messaging**: Cleaner log output without unnecessary warnings
- **Code consistency**: Uniform formatting and style across codebase

### Technical Changes

#### BifrostIntegration Enhancement
```go
// Added built-in module detection
func (bi *BifrostIntegration) isBuiltinModule(packageName string) bool {
    builtinModules := map[string]bool{
        "file": true,
        "os":   true,
        "time": true,
        "http": true,
    }
    return builtinModules[packageName]
}
```

#### Analyzer Test Fixes
```go
// Fixed hover test position calculation
position := protocol.Position{Line: 16, Character: 5} // "Person" grimoire
// Previously: Line: 26 (incorrect due to multiline string indexing)
```

#### Formatter Test Corrections
```go
// Corrected Carrion syntax for attempt blocks
input := `spell test_attempt():
    attempt:
        x=risky_operation()
    ensnare (e):        // Fixed: was "ensnare ValueError as e:"
        print("Error:",e)
    resolve:
        cleanup()`
```

### Test Results

All test suites now pass with 100% success rate:

```bash
✅ github.com/javanhut/CarrionLSP         - Main LSP functionality
✅ internal/analyzer                      - Core analysis engine  
✅ internal/protocol                      - LSP protocol handling
✅ All formatter tests passing            - Syntax parsing & formatting
✅ All analyzer tests passing             - Symbol resolution & completion
```

### Impact

- **Performance**: Reduced startup time by eliminating unnecessary external package loading
- **User Experience**: Cleaner log output without import warnings
- **Reliability**: Robust test suite ensures consistent behavior across updates
- **Maintainability**: Better separation of built-in vs. external module handling

---

## [2.0.0] - 2024-01-01

### Added

#### Dynamic Runtime Integration
- **Revolutionary dynamic loading system**: All language features loaded directly from Carrion runtime
- **Real-time discovery**: 75+ built-in functions automatically detected
- **Standard library integration**: 22+ grimoires (String, Array, Integer, etc.) with full method completion
- **Runtime synchronization**: Always up-to-date with latest Carrion language features

#### Enhanced Type System  
- **Constructor call recognition**: Automatic type inference from `calc = Calculator()`
- **Primitive type enhancement**: String/Array/Integer methods available on primitive values
- **Cross-reference support**: Navigate between definitions and usages
- **Context-aware completion**: Smart suggestions based on variable types and scope

#### Package & Module Integration
- **Bifrost package manager**: Automatic discovery and loading of installed packages
- **Import resolution**: Smart import statement processing with auto-loading
- **Module discovery**: Seamless integration with Carrion's module system
- **Dynamic package loading**: Hot-load packages without restarting the server

#### Advanced LSP Features
- **Intelligent code completion**: Context-aware suggestions with method signatures
- **Hover information**: Rich tooltips with documentation and type information
- **Document symbols**: Hierarchical view of grimoires, spells, and variables
- **Go to definition**: Navigate to symbol definitions across files
- **Semantic tokens**: Syntax highlighting support
- **Error diagnostics**: Real-time syntax and semantic error reporting

### Architecture

#### Core Components
- **Dynamic Loader** (`internal/analyzer/dynamic_loader.go`): Runtime feature loading
- **Bifrost Integration** (`internal/analyzer/bifrost_integration.go`): Package management
- **Analysis Engine** (`internal/analyzer/analyzer.go`): Core language analysis
- **LSP Features** (`internal/analyzer/features.go`): Completion, hover, navigation
- **Protocol Handler** (`internal/server/`): LSP protocol implementation

#### Performance Optimizations
- **Concurrent analysis**: Multi-threaded document processing
- **Smart caching**: Optimized symbol tables and memory usage
- **Dynamic refresh**: Update language features without restart
- **Error recovery**: Robust parsing with graceful error handling

### Documentation

- Comprehensive README with usage examples
- Architecture documentation with technical details
- Editor integration guides for VS Code, Neovim, Emacs, Vim
- API documentation for LSP endpoints
- Troubleshooting guide and performance tuning

### Development Infrastructure

- Complete test suite with benchmarks
- Debug utilities for development
- Makefile for build automation
- Cross-platform build support
- Code quality tools and linting

---

## [1.0.0] - 2023-12-01

### Initial Release

#### Basic LSP Features
- Basic code completion for keywords and built-ins
- Simple hover information
- Document symbols extraction
- Syntax error reporting
- Go to definition (limited)

#### Static Configuration
- Manual built-in function definitions
- Static grimoire method lists
- Basic package import support
- Simple error handling

#### Editor Support
- VS Code extension
- Basic Neovim configuration
- Emacs lsp-mode support

### Notes

This was the initial static implementation before the dynamic loading system was introduced in v2.0.0.