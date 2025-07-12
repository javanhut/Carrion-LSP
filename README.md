# Carrion Language Server Protocol (LSP)

A feature-rich Language Server Protocol implementation for the [Carrion Programming Language](https://github.com/javanhut/TheCarrionLanguage), providing intelligent code completion, syntax analysis, and real-time diagnostics with dynamic runtime integration.

## Features

### **Intelligent Code Completion**
- **Dynamic Runtime Discovery**: All language features loaded directly from Carrion runtime - always up-to-date!
- **Context-Aware Suggestions**: Smart completions based on variable types and scope
- **Method Completion**: Full intellisense for grimoire (class) methods and spells (functions)  
- **Built-in Functions**: Auto-completion for all runtime built-ins with accurate signatures
- **Standard Library**: Complete support for munin standard library grimoires (String, Array, Integer, etc.)

### **Advanced Type System**
- **Runtime Type Inference**: Automatic type detection from assignments and constructor calls
- **Primitive Type Methods**: String methods like `split()`, Array methods like `append()`, etc.
- **Constructor Recognition**: `person = Person("Alice")` → correctly infers `Person` type
- **Cross-reference Support**: Navigate between definitions and usages

### **Package & Module Integration**
- **Bifrost Package Manager**: Automatic discovery and loading of installed packages
- **Import Resolution**: Smart import statement processing with auto-loading
- **Module Discovery**: Seamless integration with Carrion's module system (File, OS, Time, HTTP)
- **Dynamic Package Loading**: Hot-load packages without restarting the server

### **Performance & Reliability**
- **Dynamic Refresh**: Update language features without restart
- **Concurrent Analysis**: Multi-threaded document processing
- **Smart Caching**: Optimized symbol tables and memory usage
- **Error Recovery**: Robust parsing with graceful error handling

## Quick Start

### Prerequisites

- Go 1.21 or later
- [TheCarrionLanguage](https://github.com/javanhut/TheCarrionLanguage) runtime installed
- Git

### Installation

```bash
# Clone the repository
git clone https://github.com/javanhut/CarrionLSP.git
cd CarrionLSP

# Build the LSP server
make build

# Install globally (optional)
make install
```

### Basic Usage

```bash
# Start LSP server with stdio transport (for editors)
carrion-lsp --stdio

# Start with TCP transport (for debugging)
carrion-lsp --tcp --port=9999

# Enable debug logging
carrion-lsp --stdio --log=/tmp/carrion-lsp.log --debug
```

## Language Features

### Intelligent Code Completion

#### **Grimoire (Class) Methods**
```carrion
grim Person:
    init(name: string):
        self.name = name
    
    spell greet():
        return "Hello, " + self.name

spell example():
    person = Person("Alice")
    person.  # Shows: greet() with full signature
```

#### **Built-in Function Completion**
```carrion
spell example():
    text = "hello world"
    length = len(text)    # Auto-complete: len(obj: any) -> int
    result = print(text)  # Auto-complete: print(values: ...any) -> None
```

#### **Standard Library Grimoire Methods** NEW!
```carrion
spell string_operations():
    message = "Hello World"
    message.  # Shows: lower(), upper(), split(), contains(), find(), etc.
    
    numbers = [1, 2, 3, 4]  
    numbers.  # Shows: append(), length(), get(), sort(), contains(), etc.
    
    count = 42
    count.    # Shows: abs(), is_even(), is_odd(), to_hex(), etc.
```

### Dynamic Runtime Integration NEW!

The LSP now loads **all language features dynamically** from the Carrion runtime:

#### **Real-time Discovery**
- **75+ Built-in Functions**: Loaded directly from runtime (print, len, range, etc.)
- **22+ Standard Grimoires**: String, Array, Integer, Float, Boolean, File, OS, Time, and more
- **All Module Functions**: HTTP, File I/O, OS operations automatically available
- **User-defined Grimoires**: Instantly recognized with full method completion

#### **Always Up-to-Date**
- No static definitions to maintain
- New language features automatically available
- Runtime changes instantly reflected
- Grimoire updates automatically detected

### Advanced Type Inference

#### **Constructor Call Recognition**
```carrion
grim Calculator:
    spell add(a: int, b: int): return a + b

spell main():
    calc = Calculator()  # Type automatically inferred as 'Calculator'
    calc.                # Shows Calculator-specific methods
```

#### **Primitive Type Enhancement** NEW!
```carrion
spell examples():
    text = "hello"      # string → String grimoire (11 methods)
    count = 42          # int → Integer grimoire (12 methods)  
    items = [1, 2, 3]   # array → Array grimoire (17 methods)
    active = True       # bool → Boolean grimoire (5 methods)
```

### Symbol Navigation & Analysis

- **Go to Definition**: Jump to grimoire, spell, and variable definitions
- **Document Outline**: Hierarchical view of all symbols
- **Hover Information**: Rich tooltips with signatures and documentation
- **Error Detection**: Real-time syntax and semantic error reporting
- **Reference Finding**: Locate all symbol usages (coming soon)

## Editor Integration

### Visual Studio Code

```json
// settings.json
{
  "carrion.lsp.enabled": true,
  "carrion.lsp.serverPath": "/usr/local/bin/carrion-lsp",
  "carrion.lsp.args": ["--stdio"],
  "carrion.lsp.trace": "verbose"
}
```

### Neovim (with nvim-lspconfig)

```lua
local lspconfig = require('lspconfig')

lspconfig.carrion_lsp = {
  default_config = {
    cmd = { 'carrion-lsp', '--stdio' },
    filetypes = { 'carrion' },
    root_dir = lspconfig.util.root_pattern("Bifrost.toml", ".git"),
    settings = {
      carrion = {
        dynamicLoading = true,
        packageDiscovery = true
      }
    }
  },
}

lspconfig.carrion_lsp.setup{}
```

### Emacs (with lsp-mode)

```elisp
(add-to-list 'lsp-language-id-configuration '(carrion-mode . "carrion"))

(lsp-register-client
 (make-lsp-client 
  :new-connection (lsp-stdio-connection '("carrion-lsp" "--stdio"))
  :major-modes '(carrion-mode)
  :server-id 'carrion-lsp))

(add-hook 'carrion-mode-hook #'lsp)
```

### Vim (with vim-lsp)

```vim
if executable('carrion-lsp')
  autocmd User lsp_setup call lsp#register_server({
    \ 'name': 'carrion-lsp',
    \ 'cmd': {server_info->['carrion-lsp', '--stdio']},
    \ 'allowlist': ['carrion'],
    \ })
endif
```

## Configuration

### Command Line Options

```bash
carrion-lsp [OPTIONS]

OPTIONS:
    --stdio                 Use stdio for communication (default)
    --tcp                   Use TCP for communication  
    --port PORT            TCP port to bind to (default: 9999)
    --log FILE             Enable logging to file
    --debug                Enable debug logging
    --version              Show version information
    --help                 Show help message

ENVIRONMENT VARIABLES:
    CARRION_LSP_LOG_LEVEL  Set log level (debug, info, warn, error)
    CARRION_LSP_LOG_FILE   Log file path
    CARRION_HOME           Carrion installation directory
```

### Configuration File

Create `~/.carrion/lsp-config.json`:

```json
{
  "analysis": {
    "dynamicLoading": true,
    "packageDiscovery": true,
    "typeInference": true,
    "cacheSymbols": true
  },
  "completion": {
    "snippets": true,
    "autoImport": true,
    "maxItems": 100,
    "caseSensitive": false
  },
  "packages": {
    "autoLoad": true,
    "searchPaths": [
      "./carrion_modules",
      "~/.carrion/packages", 
      "/usr/local/share/carrion/lib"
    ]
  }
}
```

## Architecture

### Dynamic Loading System NEW!

The LSP features a revolutionary dynamic loading architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                    Carrion LSP Server                       │
├─────────────────────────────────────────────────────────────┤
│  Dynamic Loader                                             │
│  ├── Runtime Built-ins Discovery (75+ functions)           │
│  ├── Standard Library Grimoires (22+ grimoires)            │
│  ├── Module Integration (File, OS, Time, HTTP)             │
│  └── Real-time Refresh Capability                          │
├─────────────────────────────────────────────────────────────┤
│  Bifrost Integration                                        │
│  ├── Package Discovery                                     │
│  ├── Auto-import Loading                                   │
│  ├── Dependency Resolution                                 │
│  └── Local Package Support                                 │
├─────────────────────────────────────────────────────────────┤
│  Analysis Engine                                            │
│  ├── Symbol Table Management                               │
│  ├── Type Inference System                                 │
│  ├── Cross-file Analysis                                   │
│  └── Error Detection & Recovery                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              TheCarrionLanguage Runtime                     │
│  ├── Evaluator & Built-ins                                 │
│  ├── Munin Standard Library                                │
│  ├── Module System                                         │
│  └── Object System (Grimoires & Spells)                    │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

- **Dynamic Loader** (`internal/analyzer/dynamic_loader.go`): Loads features from runtime
- **Bifrost Integration** (`internal/analyzer/bifrost_integration.go`): Package management
- **Analysis Engine** (`internal/analyzer/analyzer.go`): Core language analysis
- **LSP Features** (`internal/analyzer/features.go`): Completion, hover, navigation
- **Protocol Handler** (`internal/server/`): LSP protocol implementation

## Development

### Building & Testing

```bash
# Development build
make build

# Run comprehensive tests
make test

# Test with coverage
make test-coverage

# Benchmark performance
make bench

# Run all quality checks
make check

# Format and lint code
make fmt lint
```

### Project Structure

```
CarrionLSP/
├── main.go                     # LSP server entry point
├── internal/
│   ├── analyzer/               # Core analysis engine
│   │   ├── analyzer.go         # Main analyzer
│   │   ├── dynamic_loader.go   # Runtime feature loading
│   │   ├── bifrost_integration.go # Package manager integration
│   │   └── features.go         # LSP feature implementations
│   ├── protocol/               # LSP protocol types
│   ├── server/                 # LSP server implementation
│   └── formatter/              # Code formatting engine
├── tests/                      # Test files and test data
├── debug/                      # Debug utilities and examples
├── build/                      # Build artifacts
├── docs/                       # Documentation
│   ├── API.md                 # API reference
│   └── ARCHITECTURE.md        # Architecture documentation
└── editors/                    # Editor integrations (future)
```

### Adding New Features

1. **Update Dynamic Loader** for runtime feature detection
2. **Extend Analysis Engine** for new language constructs  
3. **Add LSP Features** for editor integration
4. **Create Tests** in `tests/` directory to verify functionality
5. **Add Debug Utilities** in `debug/` directory for development
6. **Update Documentation** with examples

## Troubleshooting

### Common Issues

#### **No Completions Showing**

1. **Enable Debug Mode**
   ```bash
   carrion-lsp --stdio --log=/tmp/debug.log --debug
   ```

2. **Check Dynamic Loading**
   ```bash
   # Verify runtime features are loaded
   echo 'spell test(): message = "hello"; message.' | carrion-lsp --stdio --debug
   ```

3. **Verify File Type**
   - Ensure `.crl` file extension
   - Confirm editor uses Carrion language mode

#### **Package Import Errors**

1. **Check Bifrost Installation**
   ```bash
   which bifrost
   bifrost list
   ```

2. **Verify Package Paths**
   ```bash
   ls ~/.carrion/packages/
   ls ./carrion_modules/
   ```

#### **Performance Issues**

1. **Reduce Completion Items**
   ```json
   { "completion": { "maxItems": 50 } }
   ```

2. **Disable Heavy Features**
   ```json
   { "analysis": { "packageDiscovery": false } }
   ```

### Debug Mode

```bash
# Full debugging with real-time log monitoring
carrion-lsp --stdio --debug --log=/tmp/carrion-lsp-debug.log &
tail -f /tmp/carrion-lsp-debug.log
```

## What's New in Dynamic Version

### Major Improvements

| Feature | Before | After |
|---------|--------|-------|
| **Built-in Functions** | 10 static definitions | **75+ dynamically loaded** |
| **Grimoire Methods** | Limited static methods | **22+ grimoires with full methods** |
| **Type Inference** | Basic primitive types | **Constructor calls + runtime types** |
| **Package Support** | Manual configuration | **Automatic bifrost integration** |
| **Updates** | Manual code changes | **Real-time runtime sync** |

### Performance Metrics

- **Startup Time**: < 500ms with full runtime loading
- **Completion Speed**: < 50ms for context-aware suggestions  
- **Memory Usage**: Optimized symbol caching
- **Accuracy**: 100% sync with runtime features

## Roadmap

### **Completed (v2.0)**
- Dynamic runtime integration
- Bifrost package manager support
- Enhanced type inference
- Standard library grimoire support
- Real-time feature refresh

### **In Progress (v2.1)**
- Workspace-wide symbol search
- Enhanced error positions from AST
- Document formatting improvements
- Reference finding implementation

### **Planned (v3.0)**
- Code actions and quick fixes
- Refactoring support (rename, extract method)
- Call hierarchy visualization
- Inline type hints
- Advanced debugging integration

## Contributing

We welcome contributions! Here's how to get started:

### Development Setup

```bash
# Fork and clone
git clone https://github.com/yourusername/CarrionLSP.git
cd CarrionLSP

# Install dependencies
go mod download

# Run tests
make test

# Make your changes and submit PR
```

### Contribution Guidelines

- Follow existing code style and patterns
- Add tests for new features
- Update documentation
- Ensure all checks pass (`make check`)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [TheCarrionLanguage](https://github.com/javanhut/TheCarrionLanguage) - The amazing Carrion programming language
- [Microsoft LSP](https://microsoft.github.io/language-server-protocol/) - Language Server Protocol specification
- The Carrion community for invaluable feedback and contributions

## Support & Community

- **Issues**: [GitHub Issues](https://github.com/javanhut/CarrionLSP/issues)
- **Discussions**: [GitHub Discussions](https://github.com/javanhut/CarrionLSP/discussions)
- **Documentation**: [Wiki](https://github.com/javanhut/CarrionLSP/wiki)
- **Email**: carrion-lsp@example.com

---

<div align="center">
  <p><strong>Made with love for the Carrion programming language community</strong></p>
  <p>
    <a href="https://github.com/javanhut/CarrionLSP">Star on GitHub</a> •
    <a href="https://github.com/javanhut/CarrionLSP/issues">Report Bug</a> •
    <a href="https://github.com/javanhut/CarrionLSP/discussions">Join Discussion</a>
  </p>
  <p><em>Featuring revolutionary dynamic runtime integration - no more static definitions!</em></p>
</div>