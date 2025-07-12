# Contributing to Carrion-LSP

Thank you for your interest in contributing to the Carrion Language Server Protocol! This document provides guidelines and information for contributors.

## How to Contribute

### Reporting Issues

1. **Search existing issues** to avoid duplicates
2. **Use issue templates** when available
3. **Provide clear reproduction steps**
4. **Include system information** (OS, Go version, Carrion version)
5. **Attach logs** when possible

#### Bug Report Template

```markdown
**Bug Description**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Create file with content '...'
2. Position cursor at '...'
3. Request completion
4. See error

**Expected Behavior**
What you expected to happen.

**Environment**
- OS: [e.g. Ubuntu 22.04]
- Go Version: [e.g. 1.21.0]
- Carrion Version: [e.g. 1.0.0]
- LSP Version: [e.g. 2.0.0]

**Logs**
```
Paste relevant log output here
```
```

### Suggesting Features

1. **Check the roadmap** to see if it's already planned
2. **Describe the use case** clearly
3. **Explain the expected behavior**
4. **Consider implementation complexity**

## Development Setup

### Prerequisites

- **Go 1.21+**: Required for building the LSP server
- **TheCarrionLanguage**: Required for dynamic runtime integration
- **Git**: For version control
- **Make**: For build automation
- **golangci-lint**: For code linting (optional but recommended)

### Initial Setup

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/yourusername/CarrionLSP.git
cd CarrionLSP

# 3. Add upstream remote
git remote add upstream https://github.com/javanhut/CarrionLSP.git

# 4. Install dependencies
go mod download

# 5. Verify setup
make test
```

### Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/your-feature-name

# 2. Make your changes
# Edit files...

# 3. Run tests and checks
make check          # Runs fmt, vet, lint, and test
make test-coverage  # Run tests with coverage

# 4. Commit your changes
git add .
git commit -m "feat: add your feature description"

# 5. Push to your fork
git push origin feature/your-feature-name

# 6. Create pull request on GitHub
```

## Code Guidelines

### Code Style

- **Follow Go conventions**: Use `gofmt`, `go vet`, and `golangci-lint`
- **Use meaningful names**: Variables, functions, and types should be self-documenting
- **Write comments**: Especially for complex logic and public APIs
- **Keep functions small**: Aim for single responsibility
- **Handle errors**: Always check and handle errors appropriately

#### Example Good Code

```go
// analyzeVariable extracts variable information from assignment statements
// and adds it to the symbol table with inferred type information.
func (a *Analyzer) analyzeVariable(node *ast.AssignStatement, symbols *SymbolTable) error {
    ident, ok := node.Name.(*ast.Identifier)
    if !ok {
        return fmt.Errorf("expected identifier, got %T", node.Name)
    }

    variable := &VariableSymbol{
        Name:  ident.Value,
        Range: a.astNodeToRange(node),
    }

    // Infer type from assignment value
    if node.Value != nil {
        variable.Type = a.inferTypeWithContext(node.Value, symbols)
    }

    symbols.Variables[variable.Name] = variable
    return nil
}
```

### Project Structure

```
CarrionLSP/
├── main.go                # LSP server entry point
├── internal/              # Private application code
│   ├── analyzer/          # Core analysis engine
│   ├── protocol/          # LSP protocol implementation
│   ├── server/            # LSP server implementation
│   └── formatter/         # Code formatting
├── tests/                 # Test files and test data
├── debug/                 # Debug utilities and examples
├── build/                 # Build artifacts
├── docs/                  # Documentation
│   ├── API.md            # API reference
│   └── ARCHITECTURE.md   # Architecture documentation
└── editors/               # Editor integrations (future)
```

### Testing Guidelines

#### Unit Tests

- **Test all public functions**
- **Use table-driven tests** for multiple scenarios
- **Mock external dependencies**
- **Test error conditions**

```go
func TestAnalyzeVariable(t *testing.T) {
    tests := []struct {
        name     string
        code     string
        expected string
        wantErr  bool
    }{
        {
            name:     "string literal",
            code:     `message = "hello"`,
            expected: "string",
            wantErr:  false,
        },
        {
            name:     "integer literal",
            code:     `count = 42`,
            expected: "int",
            wantErr:  false,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Integration Tests

- **Test LSP protocol compliance**
- **Test editor integration scenarios**
- **Use real Carrion code examples**

```go
func TestCompletionIntegration(t *testing.T) {
    server := setupTestServer(t)
    
    // Send textDocument/didOpen
    server.SendNotification("textDocument/didOpen", DidOpenParams{
        TextDocument: TextDocumentItem{
            URI:  "file:///test.crl",
            Text: `spell test(): message = "hello"; message.`,
        },
    })
    
    // Request completion
    response := server.SendRequest("textDocument/completion", CompletionParams{
        TextDocument: TextDocumentIdentifier{URI: "file:///test.crl"},
        Position:     Position{Line: 0, Character: 40},
    })
    
    // Verify response
    assert.Contains(t, response.Items, "lower")
    assert.Contains(t, response.Items, "upper")
}
```

### Documentation

- **Update README.md** for user-facing changes
- **Add godoc comments** for all public APIs
- **Include code examples** in documentation
- **Update CHANGELOG.md** for notable changes

#### Godoc Example

```go
// DynamicLoader provides dynamic loading of Carrion runtime components.
// It connects to the Carrion evaluator to discover built-in functions,
// standard library grimoires, and module definitions at runtime.
//
// Example usage:
//
//     loader := NewDynamicLoader()
//     builtins := loader.GetBuiltins()
//     for name, info := range builtins {
//         fmt.Printf("Function: %s -> %s\n", name, info.ReturnType)
//     }
type DynamicLoader struct {
    env       *object.Environment
    builtins  map[string]*BuiltinInfo
    grimoires map[string]*GrimoireInfo
}
```

## 🏗 Architecture Overview

Understanding the LSP architecture will help you contribute effectively:

### Core Components

1. **LSP Server** (`internal/server/`)
   - Handles LSP protocol messages
   - Manages client communication
   - Dispatches requests to analyzer

2. **Analyzer** (`internal/analyzer/`)
   - Parses Carrion code using TheCarrionLanguage
   - Builds and maintains symbol tables
   - Provides language features (completion, hover, etc.)

3. **Dynamic Loader** (`internal/analyzer/dynamic_loader.go`)
   - Loads features from Carrion runtime
   - Discovers built-ins and grimoires
   - Provides real-time updates

4. **Bifrost Integration** (`internal/analyzer/bifrost_integration.go`)
   - Integrates with package manager
   - Handles import resolution
   - Manages package dependencies

### Data Flow

```
Editor Request → LSP Server → Analyzer → Dynamic Loader → Carrion Runtime
                                     ↓
Editor Response ← LSP Server ← Completion/Hover/etc. ← Symbol Tables
```

## 🎯 Contribution Areas

### High Priority

1. **Performance Improvements**
   - Optimize symbol table operations
   - Implement incremental parsing
   - Add more efficient caching

2. **LSP Feature Completion**
   - Implement find references
   - Add rename symbol support
   - Enhance go-to-definition

3. **Error Handling**
   - Improve error messages
   - Add recovery from parse errors
   - Better diagnostics reporting

### Medium Priority

1. **Editor Integrations**
   - VS Code extension improvements
   - Vim/Neovim plugin enhancements
   - Emacs package updates

2. **Testing & Quality**
   - Increase test coverage
   - Add performance benchmarks
   - Improve integration tests

3. **Documentation**
   - API documentation
   - Tutorial content
   - Video guides

### Nice to Have

1. **Advanced Features**
   - Code actions and quick fixes
   - Refactoring support
   - Call hierarchy

2. **Developer Experience**
   - Debug mode improvements
   - Better logging
   - Configuration UI

## Testing Your Changes

### Pre-submission Checklist

- [ ] **Code compiles** without warnings
- [ ] **All tests pass** (`make test`)
- [ ] **Linting passes** (`make lint`)
- [ ] **Code is formatted** (`make fmt`)
- [ ] **Documentation updated** if needed
- [ ] **Changelog entry added** for notable changes

### Manual Testing

1. **Build and run** the LSP server
   ```bash
   make build
   ./build/carrion-lsp --stdio --debug --log=/tmp/test.log
   ```

2. **Use debug utilities** for specific testing
   ```bash
   # Test completion context detection
   go run debug/debug_completion_context.go
   
   # Test dynamic loading features
   go run debug/debug_dynamic_string.go
   
   # Test method completion
   go run debug/debug_method_completion.go
   ```

3. **Test with editor** integration
   - Open a `.crl` file
   - Test completions, hover, navigation
   - Check error reporting

4. **Test edge cases**
   - Empty files
   - Syntax errors
   - Large files
   - Package imports

### Performance Testing

```bash
# Run benchmarks
make bench

# Use debug utilities for testing specific scenarios
go run debug/debug_completion_context.go
go run debug/debug_dynamic_string.go

# Test with large files
echo 'spell test(): message = "hello"; message.' > large-test.crl
# Test LSP with large-test.crl
```

## Pull Request Process

### Before Submitting

1. **Rebase on latest main**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Squash commits** if needed
   ```bash
   git rebase -i HEAD~n  # n = number of commits to squash
   ```

3. **Write good commit messages**
   ```
   feat: add dynamic loading for bifrost packages
   
   - Implement package discovery in standard locations
   - Add auto-loading of import statements
   - Support for relative and absolute package paths
   
   Fixes #123
   ```

### PR Template

```markdown
## Description
Brief description of the changes.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review of code performed
- [ ] Code is commented, particularly in hard-to-understand areas
- [ ] Corresponding changes to documentation made
- [ ] Changes generate no new warnings
- [ ] Tests added that prove fix is effective or feature works
- [ ] New and existing tests pass locally
```

### Review Process

1. **Automated checks** must pass (CI/CD)
2. **Code review** by maintainers
3. **Address feedback** promptly
4. **Final approval** before merge

## Recognition

Contributors will be recognized in:

- **README.md** acknowledgments section
- **CHANGELOG.md** for significant contributions
- **GitHub releases** for version contributions
- **Project documentation** for major features

## Questions?

- **GitHub Discussions**: Ask questions and discuss ideas
- **GitHub Issues**: Report bugs and request features
- **Email**: Contact maintainers directly for sensitive issues

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- **Be respectful** and inclusive
- **Accept constructive criticism** gracefully
- **Focus on what's best** for the community
- **Help others** when possible

Thank you for contributing to Carrion-LSP!