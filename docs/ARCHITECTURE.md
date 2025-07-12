# Carrion-LSP Architecture

This document provides a comprehensive overview of the Carrion Language Server Protocol architecture, design decisions, and implementation details.

## Table of Contents

- [Overview](#overview)
- [System Architecture](#system-architecture)
- [Core Components](#core-components)
- [Dynamic Loading System](#dynamic-loading-system)
- [Data Flow](#data-flow)
- [Threading Model](#threading-model)
- [Memory Management](#memory-management)
- [Performance Considerations](#performance-considerations)
- [Extension Points](#extension-points)

## Overview

Carrion-LSP is a Language Server Protocol implementation that provides intelligent language support for the Carrion programming language. The architecture is designed around the principle of **dynamic runtime integration**, ensuring that the LSP server always reflects the current state of the Carrion language runtime.

### Key Design Principles

1. **Dynamic Over Static**: Load language features from runtime rather than static definitions
2. **Performance First**: Optimize for low-latency responses and efficient memory usage
3. **Extensibility**: Provide clear extension points for new features
4. **Reliability**: Graceful error handling and recovery mechanisms
5. **Simplicity**: Clear separation of concerns and minimal dependencies

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Editor / IDE                                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   VS Code   │  │   Neovim    │  │    Emacs    │  │     Vim     │  ...   │
│  └─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────┬───────────────────────────────────────────────────┘
                          │ LSP Protocol (JSON-RPC)
                          │ Transport: stdio/TCP
                          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Carrion-LSP Server                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│  LSP Server Layer                                                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Protocol Handler│  │ Message Router  │  │ Response Manager│             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
├─────────────────────────────────────────────────────────────────────────────┤
│  Analysis Engine                                                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Document Manager│  │ Symbol Tables   │  │ Type Inference  │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Completion      │  │ Hover Provider  │  │ Navigation      │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
├─────────────────────────────────────────────────────────────────────────────┤
│  Dynamic Loading System                                                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Runtime Loader  │  │ Bifrost Bridge  │  │ Package Manager │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
├─────────────────────────────────────────────────────────────────────────────┤
│  Foundation Layer                                                           │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ AST Integration │  │ Error Handling  │  │ Logging System  │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
└─────────────────────────┬───────────────────────────────────────────────────┘
                          │ Direct Runtime Integration
                          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                      TheCarrionLanguage Runtime                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Lexer & Parser  │  │ Evaluator       │  │ Object System   │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐             │
│  │ Built-ins       │  │ Munin Stdlib    │  │ Module System   │             │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘             │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Core Components

### 1. LSP Server Layer

#### Protocol Handler (`internal/server/`)
- **Responsibility**: LSP protocol compliance and message handling
- **Key Files**:
  - `server.go`: Main server implementation
  - `handler.go`: Request/response handling
  - `transport.go`: Communication transport (stdio/TCP)

```go
type Server struct {
    analyzer    *analyzer.Analyzer
    transport   Transport
    capabilities ServerCapabilities
    mu          sync.RWMutex
}

func (s *Server) HandleRequest(ctx context.Context, req *protocol.Request) *protocol.Response {
    switch req.Method {
    case "textDocument/completion":
        return s.handleCompletion(ctx, req)
    case "textDocument/hover":
        return s.handleHover(ctx, req)
    // ... more handlers
    }
}
```

#### Message Router
- Routes incoming LSP requests to appropriate handlers
- Manages request/response correlation
- Handles notifications vs. requests

#### Response Manager
- Formats responses according to LSP specification
- Handles error responses and edge cases
- Manages response caching where appropriate

### 2. Analysis Engine

#### Document Manager (`internal/analyzer/analyzer.go`)
- **Responsibility**: Document lifecycle and symbol table management
- **Key Features**:
  - Document parsing using TheCarrionLanguage parser
  - Symbol table construction and maintenance
  - Multi-document analysis coordination

```go
type Analyzer struct {
    mu                sync.RWMutex
    documents         map[string]*Document
    builtins          map[string]*BuiltinInfo
    carriongGrimoires map[string]*GrimoireInfo
    dynamicLoader     *DynamicLoader
    bifrostIntegration *BifrostIntegration
}

func (a *Analyzer) UpdateDocument(uri, content string, program *ast.Program) *Document {
    // Parse with TheCarrionLanguage parser
    if program == nil {
        l := lexer.New(content)
        p := parser.New(l)
        program = p.ParseProgram()
    }
    
    // Build symbol table
    symbols := a.buildSymbolTable(program)
    
    // Create document
    doc := &Document{
        URI:     uri,
        Content: content,
        AST:     program,
        Symbols: symbols,
        Version: a.getNextVersion(uri),
    }
    
    // Auto-load imports
    a.bifrostIntegration.AutoLoadImports(doc)
    
    return doc
}
```

#### Symbol Tables
- **Grimoires**: Class-like structures with methods (spells)
- **Spells**: Functions/methods with parameter information  
- **Variables**: Variable declarations with type inference
- **Imports**: Module imports with aliases

```go
type SymbolTable struct {
    Grimoires map[string]*GrimoireSymbol
    Spells    map[string]*SpellSymbol
    Variables map[string]*VariableSymbol
    Imports   map[string]*ImportSymbol
}
```

#### Type Inference (`internal/analyzer/analyzer.go`)
- **Context-Aware**: Uses symbol table for accurate inference
- **Constructor Recognition**: Detects grimoire instantiation
- **Primitive Mapping**: Maps primitives to grimoire types

```go
func (a *Analyzer) inferTypeWithContext(expr ast.Expression, symbols *SymbolTable) string {
    switch node := expr.(type) {
    case *ast.CallExpression:
        // Constructor call detection
        if ident, ok := node.Function.(*ast.Identifier); ok {
            if _, exists := symbols.Grimoires[ident.Value]; exists {
                return ident.Value // Grimoire type
            }
        }
    case *ast.StringLiteral:
        return "string" // Maps to String grimoire
    case *ast.IntegerLiteral:
        return "int"    // Maps to Integer grimoire
    }
    return "unknown"
}
```

### 3. Language Features

#### Completion Provider (`internal/analyzer/features.go`)
- **Context Detection**: Determines completion type (method vs. general)
- **Dynamic Suggestions**: Uses runtime-loaded features
- **Smart Filtering**: Prefix-based matching with token extraction

```go
func (a *Analyzer) GetCompletions(uri string, position protocol.Position) []protocol.CompletionItem {
    doc := a.documents[uri]
    lines := strings.Split(doc.Content, "\n")
    currentLine := lines[position.Line]
    prefix := currentLine[:position.Character]
    
    if strings.HasSuffix(prefix, ".") {
        return a.getMethodCompletions(doc, prefix)
    } else {
        return a.getGeneralCompletions(doc, prefix)
    }
}
```

#### Hover Provider
- **Symbol Information**: Rich tooltips with signatures
- **Documentation**: Extracted from docstrings and runtime
- **Type Information**: Shows variable types and function signatures

#### Navigation Provider
- **Go-to-Definition**: AST-based symbol location
- **Document Symbols**: Hierarchical structure extraction
- **Reference Finding**: Cross-file symbol usage (planned)

### 4. Dynamic Loading System

#### Runtime Loader (`internal/analyzer/dynamic_loader.go`)
- **Built-in Discovery**: Loads functions from Carrion evaluator
- **Grimoire Detection**: Discovers standard library grimoires
- **Module Integration**: Integrates File, OS, Time, HTTP modules

```go
type DynamicLoader struct {
    env       *object.Environment
    builtins  map[string]*BuiltinInfo
    grimoires map[string]*GrimoireInfo
}

func (dl *DynamicLoader) loadBuiltins() {
    // Get built-ins from Carrion runtime
    runtimeBuiltins := evaluator.GetBuiltins()
    
    for name := range runtimeBuiltins {
        dl.builtins[name] = &BuiltinInfo{
            Name:        name,
            Type:        "function",
            Description: dl.getBuiltinDescription(name),
            Parameters:  dl.extractBuiltinParameters(name),
            ReturnType:  dl.inferBuiltinReturnType(name),
        }
    }
}
```

#### Bifrost Integration (`internal/analyzer/bifrost_integration.go`)
- **Package Discovery**: Scans standard package locations
- **Import Resolution**: Handles import statement processing
- **Auto-Loading**: Automatically loads required packages

```go
func (bi *BifrostIntegration) LoadPackageFromImport(importPath string) error {
    if strings.HasPrefix(importPath, "./") {
        return bi.loadRelativePackage(importPath)
    }
    
    packageName := strings.Split(importPath, "/")[0]
    return bi.LoadPackage(packageName)
}
```

## Dynamic Loading System

The dynamic loading system is the core innovation of Carrion-LSP, providing real-time synchronization with the Carrion runtime.

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Dynamic Loader                           │
├─────────────────────────────────────────────────────────────┤
│  Runtime Discovery                                         │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ Built-in Loader │  │ Grimoire Loader │                  │
│  └─────────────────┘  └─────────────────┘                  │
│                                                             │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ Module Scanner  │  │ Stdlib Loader   │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│  Package Integration                                        │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ Bifrost Bridge  │  │ Import Resolver │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────┬───────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│              Carrion Runtime Environment                   │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ Evaluator       │  │ Object System   │                  │
│  └─────────────────┘  └─────────────────┘                  │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ Built-ins Map   │  │ Munin Stdlib    │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### Loading Process

1. **Environment Creation**: Create Carrion environment
2. **Stdlib Loading**: Load munin standard library
3. **Built-in Extraction**: Extract functions from evaluator
4. **Grimoire Discovery**: Scan environment for grimoires
5. **Module Integration**: Load File, OS, Time, HTTP modules
6. **Cache Construction**: Build optimized lookup tables

### Benefits

- **Always Current**: Reflects latest runtime features
- **Zero Maintenance**: No manual updates required
- **Complete Coverage**: Discovers all available features
- **Type Accuracy**: Uses actual runtime type information

## Data Flow

### Completion Request Flow

```
Editor Request
      │
      ▼
┌─────────────────┐
│ LSP Server      │ 
│ - Validate      │
│ - Parse position│
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Analyzer        │
│ - Get document  │
│ - Extract prefix│
│ - Determine type│
└─────────────────┘
      │
      ▼
┌─────────────────┐     ┌─────────────────┐
│ Method Completion│ OR  │General Completion│
│ - Find variable │     │ - Built-ins     │
│ - Get type      │     │ - Keywords      │
│ - Find methods  │     │ - Grimoires     │
└─────────────────┘     └─────────────────┘
      │                         │
      └──────────┬──────────────┘
                 ▼
┌─────────────────┐
│ Dynamic Loader  │
│ - Runtime data  │
│ - Cached info   │
│ - Type mapping  │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Response Format │
│ - LSP protocol  │
│ - Sort & filter │
│ - Add metadata  │
└─────────────────┘
      │
      ▼
Editor Response
```

### Document Update Flow

```
Document Change
      │
      ▼
┌─────────────────┐
│ Document Manager│
│ - Parse content │
│ - Update version│
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ AST Generation  │
│ - Carrion lexer │
│ - Carrion parser│
│ - Error handling│
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Symbol Analysis │
│ - Grimoires     │
│ - Spells        │
│ - Variables     │
│ - Imports       │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Type Inference  │
│ - Context-aware │
│ - Constructor   │
│ - Primitive map │
└─────────────────┘
      │
      ▼
┌─────────────────┐
│ Bifrost Auto-   │
│ Load Imports    │
│ - Scan imports  │
│ - Load packages │
│ - Update symbols│
└─────────────────┘
      │
      ▼
Symbol Table Updated
```

## Threading Model

### Concurrency Strategy

Carrion-LSP uses a **reader-writer lock pattern** for safe concurrent access:

```go
type Analyzer struct {
    mu sync.RWMutex  // Protects documents and symbol tables
    // ... other fields
}

// Read operations (completions, hover, etc.)
func (a *Analyzer) GetCompletions(uri string, pos protocol.Position) []protocol.CompletionItem {
    a.mu.RLock()         // Multiple readers allowed
    defer a.mu.RUnlock()
    // ... safe read operations
}

// Write operations (document updates)
func (a *Analyzer) UpdateDocument(uri, content string, program *ast.Program) *Document {
    a.mu.Lock()          // Exclusive write access
    defer a.mu.Unlock()
    // ... document modification
}
```

### Thread Safety Guarantees

1. **Document Access**: All document operations are thread-safe
2. **Symbol Tables**: Protected by read-write locks
3. **Dynamic Loading**: Atomic updates to prevent inconsistency
4. **Cache Management**: Lock-free where possible, synchronized where needed

### Performance Implications

- **Read Scalability**: Multiple completion requests can be processed concurrently
- **Write Serialization**: Document updates are serialized but fast
- **Lock Granularity**: Fine-grained locking minimizes contention

## Memory Management

### Memory Usage Patterns

```
┌─────────────────────────────────────────────────────────────┐
│                    Memory Allocation                       │
├─────────────────────────────────────────────────────────────┤
│  Document Storage (40%)                                     │
│  ├── AST Trees                                             │
│  ├── Token Arrays                                          │
│  └── Content Strings                                       │
├─────────────────────────────────────────────────────────────┤
│  Symbol Tables (30%)                                       │
│  ├── Grimoire Definitions                                  │
│  ├── Spell Signatures                                      │
│  └── Variable Information                                  │
├─────────────────────────────────────────────────────────────┤
│  Dynamic Cache (20%)                                       │
│  ├── Built-in Functions                                    │
│  ├── Runtime Grimoires                                     │
│  └── Package Information                                   │
├─────────────────────────────────────────────────────────────┤
│  Runtime Objects (10%)                                     │
│  ├── Carrion Environment                                   │
│  ├── Evaluator State                                       │
│  └── Temporary Objects                                     │
└─────────────────────────────────────────────────────────────┘
```

### Optimization Strategies

1. **Document Caching**: Keep frequently accessed documents in memory
2. **Symbol Sharing**: Share common symbol definitions across documents
3. **Lazy Loading**: Load packages only when needed
4. **Memory Pooling**: Reuse objects for temporary operations
5. **Garbage Collection**: Periodic cleanup of unused documents

### Memory Limits

- **Document Limit**: Configurable maximum number of cached documents
- **Symbol Limit**: Automatic cleanup of unused symbols
- **Package Limit**: LRU eviction for package cache

## Performance Considerations

### Latency Targets

| Operation | Target | Typical |
|-----------|--------|---------|
| Completion Request | < 50ms | 20-30ms |
| Hover Request | < 20ms | 5-10ms |
| Document Update | < 100ms | 30-50ms |
| Dynamic Refresh | < 200ms | 80-120ms |

### Optimization Techniques

#### 1. Prefix Caching
```go
type CompletionCache struct {
    prefixCache map[string][]protocol.CompletionItem
    maxEntries  int
    lastClean   time.Time
}
```

#### 2. Symbol Indexing
```go
type SymbolIndex struct {
    byName     map[string][]*Symbol
    byType     map[string][]*Symbol
    byLocation map[string]*Symbol
}
```

#### 3. Incremental Parsing
- Parse only changed regions when possible
- Reuse AST nodes from previous versions
- Smart invalidation of dependent symbols

### Performance Monitoring

- **Request Timing**: Track response times for all operations
- **Memory Usage**: Monitor heap size and garbage collection
- **Cache Hit Rates**: Track effectiveness of caching strategies
- **Dynamic Load Times**: Monitor runtime discovery performance

## Extension Points

### Adding New Language Features

#### 1. LSP Method Handler
```go
func (s *Server) registerCustomHandlers() {
    s.handlers["textDocument/customFeature"] = s.handleCustomFeature
}

func (s *Server) handleCustomFeature(ctx context.Context, req *protocol.Request) *protocol.Response {
    // Implementation
}
```

#### 2. Analyzer Extension
```go
func (a *Analyzer) GetCustomInfo(uri string, position protocol.Position) *CustomInfo {
    // Access documents and symbols
    // Provide custom analysis
}
```

#### 3. Dynamic Loader Extension
```go
func (dl *DynamicLoader) loadCustomFeatures() {
    // Discover custom runtime features
    // Add to built-ins or grimoires
}
```

### Editor Integration

#### 1. Protocol Extensions
- Add custom LSP methods for editor-specific features
- Implement editor-specific optimizations
- Handle editor-specific configuration

#### 2. Transport Customization
```go
type CustomTransport struct {
    // Custom transport implementation
}

func (ct *CustomTransport) Read() (*protocol.Message, error) {
    // Custom message reading
}
```

### Package System Integration

#### 1. Custom Package Loaders
```go
type CustomPackageLoader interface {
    LoadPackage(name string) (*PackageInfo, error)
    DiscoverPackages() []string
}
```

#### 2. Import Resolvers
```go
type CustomImportResolver interface {
    ResolveImport(importPath string) (string, error)
    GetPackagePath(packageName string) string
}
```

---

This architecture enables Carrion-LSP to provide intelligent, real-time language support while maintaining high performance and extensibility. The dynamic loading system ensures that the LSP server always reflects the current state of the Carrion language, making it a reliable and future-proof tool for Carrion developers.