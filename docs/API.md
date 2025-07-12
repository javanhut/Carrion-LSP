# Carrion-LSP API Reference

This document provides detailed information about the Carrion Language Server Protocol implementation, including LSP methods, custom extensions, and internal APIs.

## Table of Contents

- [Standard LSP Methods](#standard-lsp-methods)
- [Carrion-Specific Extensions](#carrion-specific-extensions)
- [Internal APIs](#internal-apis)
- [Data Structures](#data-structures)
- [Error Codes](#error-codes)

## Standard LSP Methods

### Lifecycle Methods

#### `initialize`
Initializes the language server with client capabilities.

**Request:**
```json
{
  "method": "initialize",
  "params": {
    "processId": 12345,
    "rootUri": "file:///path/to/workspace",
    "capabilities": {
      "textDocument": {
        "completion": { "dynamicRegistration": true },
        "hover": { "dynamicRegistration": true },
        "definition": { "dynamicRegistration": true }
      }
    },
    "initializationOptions": {
      "dynamicLoading": true,
      "packageDiscovery": true,
      "debug": false
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "capabilities": {
      "textDocumentSync": 1,
      "completionProvider": {
        "triggerCharacters": [".", "("],
        "resolveProvider": false
      },
      "hoverProvider": true,
      "definitionProvider": true,
      "documentSymbolProvider": true,
      "referencesProvider": false,
      "documentFormattingProvider": true,
      "semanticTokensProvider": {
        "legend": {
          "tokenTypes": ["keyword", "string", "number", "operator", "variable"],
          "tokenModifiers": []
        }
      }
    },
    "serverInfo": {
      "name": "carrion-lsp",
      "version": "2.0.0"
    }
  }
}
```

#### `initialized`
Sent after the initialize response.

**Notification:**
```json
{
  "method": "initialized",
  "params": {}
}
```

#### `shutdown`
Requests server shutdown.

**Request:**
```json
{
  "method": "shutdown",
  "params": null
}
```

**Response:**
```json
{
  "result": null
}
```

#### `exit`
Exits the server process.

**Notification:**
```json
{
  "method": "exit",
  "params": null
}
```

### Text Document Synchronization

#### `textDocument/didOpen`
Notifies the server that a document was opened.

**Notification:**
```json
{
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl",
      "languageId": "carrion",
      "version": 1,
      "text": "spell greet(name: string):\n    return \"Hello, \" + name"
    }
  }
}
```

#### `textDocument/didChange`
Notifies the server of document changes.

**Notification:**
```json
{
  "method": "textDocument/didChange",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl",
      "version": 2
    },
    "contentChanges": [
      {
        "range": {
          "start": { "line": 1, "character": 0 },
          "end": { "line": 1, "character": 12 }
        },
        "text": "    print(\"Hello, \" + name)"
      }
    ]
  }
}
```

#### `textDocument/didClose`
Notifies the server that a document was closed.

**Notification:**
```json
{
  "method": "textDocument/didClose",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    }
  }
}
```

### Language Features

#### `textDocument/completion`
Provides auto-completion suggestions.

**Request:**
```json
{
  "method": "textDocument/completion",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    },
    "position": {
      "line": 2,
      "character": 15
    },
    "context": {
      "triggerKind": 2,
      "triggerCharacter": "."
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "isIncomplete": false,
    "items": [
      {
        "label": "greet",
        "kind": 2,
        "detail": "spell greet() -> string",
        "documentation": {
          "kind": "markdown",
          "value": "Greets the person with a message"
        },
        "insertText": "greet(${1})",
        "insertTextFormat": 2,
        "sortText": "0001",
        "filterText": "greet"
      },
      {
        "label": "length",
        "kind": 2,
        "detail": "spell length() -> int",
        "documentation": "Get string length",
        "insertText": "length()",
        "insertTextFormat": 1
      }
    ]
  }
}
```

#### `textDocument/hover`
Provides hover information for symbols.

**Request:**
```json
{
  "method": "textDocument/hover",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    },
    "position": {
      "line": 1,
      "character": 10
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "contents": {
      "kind": "markdown",
      "value": "**greet**: Spell\n\n```carrion\nspell greet(name: string) -> string\n```\n\nGreets the person with a personalized message."
    },
    "range": {
      "start": { "line": 1, "character": 5 },
      "end": { "line": 1, "character": 10 }
    }
  }
}
```

#### `textDocument/definition`
Navigates to symbol definitions.

**Request:**
```json
{
  "method": "textDocument/definition",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    },
    "position": {
      "line": 5,
      "character": 8
    }
  }
}
```

**Response:**
```json
{
  "result": [
    {
      "uri": "file:///path/to/file.crl",
      "range": {
        "start": { "line": 0, "character": 0 },
        "end": { "line": 0, "character": 25 }
      }
    }
  ]
}
```

#### `textDocument/references`
Finds all references to a symbol.

**Request:**
```json
{
  "method": "textDocument/references",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    },
    "position": {
      "line": 1,
      "character": 10
    },
    "context": {
      "includeDeclaration": true
    }
  }
}
```

**Response:**
```json
{
  "result": [
    {
      "uri": "file:///path/to/file.crl",
      "range": {
        "start": { "line": 0, "character": 6 },
        "end": { "line": 0, "character": 11 }
      }
    },
    {
      "uri": "file:///path/to/file.crl",
      "range": {
        "start": { "line": 5, "character": 8 },
        "end": { "line": 5, "character": 13 }
      }
    }
  ]
}
```

#### `textDocument/documentSymbol`
Returns document outline structure.

**Request:**
```json
{
  "method": "textDocument/documentSymbol",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    }
  }
}
```

**Response:**
```json
{
  "result": [
    {
      "name": "Person",
      "kind": 5,
      "range": {
        "start": { "line": 0, "character": 0 },
        "end": { "line": 10, "character": 0 }
      },
      "selectionRange": {
        "start": { "line": 0, "character": 5 },
        "end": { "line": 0, "character": 11 }
      },
      "children": [
        {
          "name": "greet",
          "kind": 6,
          "range": {
            "start": { "line": 3, "character": 4 },
            "end": { "line": 5, "character": 0 }
          },
          "selectionRange": {
            "start": { "line": 3, "character": 10 },
            "end": { "line": 3, "character": 15 }
          }
        }
      ]
    }
  ]
}
```

#### `textDocument/formatting`
Formats the entire document.

**Request:**
```json
{
  "method": "textDocument/formatting",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    },
    "options": {
      "tabSize": 4,
      "insertSpaces": true,
      "trimTrailingWhitespace": true,
      "insertFinalNewline": true
    }
  }
}
```

**Response:**
```json
{
  "result": [
    {
      "range": {
        "start": { "line": 0, "character": 0 },
        "end": { "line": 5, "character": 0 }
      },
      "newText": "grim Person:\n    init(name: string):\n        self.name = name\n\n    spell greet():\n        return \"Hello, \" + self.name\n"
    }
  ]
}
```

#### `textDocument/semanticTokens/full`
Provides semantic highlighting information.

**Request:**
```json
{
  "method": "textDocument/semanticTokens/full",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.crl"
    }
  }
}
```

**Response:**
```json
{
  "result": {
    "data": [0, 0, 4, 0, 0, 0, 5, 6, 4, 0, 1, 4, 4, 0, 0]
  }
}
```

## Carrion-Specific Extensions

### `carrion/refreshDynamicData`
Manually refreshes dynamically loaded language features.

**Request:**
```json
{
  "method": "carrion/refreshDynamicData",
  "params": {}
}
```

**Response:**
```json
{
  "result": {
    "builtinsLoaded": 75,
    "grimoiresLoaded": 22,
    "refreshTime": "2024-01-15T10:30:00Z"
  }
}
```

### `carrion/loadPackage`
Loads a specific bifrost package.

**Request:**
```json
{
  "method": "carrion/loadPackage",
  "params": {
    "packageName": "json-utils",
    "version": "1.2.0"
  }
}
```

**Response:**
```json
{
  "result": {
    "success": true,
    "packagePath": "/home/user/.carrion/packages/json-utils/1.2.0",
    "symbolsAdded": 15
  }
}
```

### `carrion/getAvailablePackages`
Gets list of available bifrost packages.

**Request:**
```json
{
  "method": "carrion/getAvailablePackages",
  "params": {}
}
```

**Response:**
```json
{
  "result": {
    "packages": [
      {
        "name": "json-utils",
        "version": "1.2.0",
        "path": "/home/user/.carrion/packages/json-utils",
        "description": "JSON parsing and manipulation utilities"
      },
      {
        "name": "http-client",
        "version": "2.1.0",
        "path": "/home/user/.carrion/packages/http-client",
        "description": "HTTP client for web requests"
      }
    ]
  }
}
```

### `carrion/getBuiltinInfo`
Gets information about built-in functions and grimoires.

**Request:**
```json
{
  "method": "carrion/getBuiltinInfo",
  "params": {
    "category": "grimoires"
  }
}
```

**Response:**
```json
{
  "result": {
    "grimoires": [
      {
        "name": "String",
        "description": "String manipulation grimoire",
        "methods": [
          {
            "name": "length",
            "signature": "length() -> int",
            "description": "Get string length"
          },
          {
            "name": "lower",
            "signature": "lower() -> string",
            "description": "Convert to lowercase"
          }
        ]
      }
    ]
  }
}
```

### `carrion/analyzeWorkspace`
Analyzes the entire workspace for symbols.

**Request:**
```json
{
  "method": "carrion/analyzeWorkspace",
  "params": {
    "includeTests": true,
    "includePackages": false
  }
}
```

**Response:**
```json
{
  "result": {
    "filesAnalyzed": 25,
    "grimoiresFound": 8,
    "spellsFound": 45,
    "variablesFound": 120,
    "analysisTime": 150
  }
}
```

## Internal APIs

### Analyzer Interface

```go
type Analyzer interface {
    // Document management
    UpdateDocument(uri, content string, program *ast.Program) *Document
    RemoveDocument(uri string)
    GetDocument(uri string) *Document
    
    // Language features
    GetCompletions(uri string, position protocol.Position) []protocol.CompletionItem
    GetHover(uri string, position protocol.Position) *protocol.Hover
    GetDefinition(uri string, position protocol.Position) []protocol.Location
    GetReferences(uri string, position protocol.Position, includeDecl bool) []protocol.Location
    GetDocumentSymbols(uri string) []protocol.DocumentSymbol
    FormatDocument(uri string, options protocol.FormattingOptions) []protocol.TextEdit
    
    // Dynamic features
    RefreshDynamicData()
    LoadBifrostPackage(packagePath string) error
    GetAvailablePackages() map[string]string
}
```

### Dynamic Loader Interface

```go
type DynamicLoader interface {
    // Feature loading
    GetBuiltins() map[string]*BuiltinInfo
    GetGrimoires() map[string]*GrimoireInfo
    RefreshDynamicData()
    
    // Package management
    LoadBifrostPackage(packagePath string) error
}
```

### Bifrost Integration Interface

```go
type BifrostIntegration interface {
    // Package discovery
    DiscoverAvailablePackages() map[string]string
    LoadPackage(packageName string) error
    LoadPackageFromImport(importPath string) error
    
    // Completion support
    GetPackageCompletions() []string
    AutoLoadImports(doc *Document) error
}
```

## Data Structures

### Symbol Information

```go
type SymbolTable struct {
    Grimoires map[string]*GrimoireSymbol
    Spells    map[string]*SpellSymbol
    Variables map[string]*VariableSymbol
    Imports   map[string]*ImportSymbol
}

type GrimoireSymbol struct {
    Name        string
    Range       protocol.Range
    InitSpell   *SpellSymbol
    Spells      map[string]*SpellSymbol
    IsArcane    bool
    Inherits    string
    DocString   string
}

type SpellSymbol struct {
    Name        string
    Range       protocol.Range
    Parameters  []Parameter
    ReturnType  string
    IsInit      bool
    IsStatic    bool
    IsPrivate   bool
    IsProtected bool
    DocString   string
    Grimoire    string
}

type VariableSymbol struct {
    Name     string
    Range    protocol.Range
    Type     string
    Value    string
    IsGlobal bool
}
```

### Built-in Information

```go
type BuiltinInfo struct {
    Name        string
    Type        string
    Description string
    Parameters  []Parameter
    ReturnType  string
}

type GrimoireInfo struct {
    Name        string
    Description string
    Spells      map[string]*BuiltinInfo
    IsStatic    bool
}

type Parameter struct {
    Name         string
    TypeHint     string
    DefaultValue string
    Range        protocol.Range
}
```

### Document Structure

```go
type Document struct {
    URI     string
    Content string
    AST     *ast.Program
    Tokens  []token.Token
    Symbols *SymbolTable
    Version int
}
```

## Error Codes

### LSP Standard Errors

| Code | Name | Description |
|------|------|-------------|
| -32700 | ParseError | Invalid JSON |
| -32600 | InvalidRequest | Invalid request object |
| -32601 | MethodNotFound | Method not found |
| -32602 | InvalidParams | Invalid method parameters |
| -32603 | InternalError | Internal JSON-RPC error |

### Carrion-Specific Errors

| Code | Name | Description |
|------|------|-------------|
| -33001 | CarrionParseError | Carrion syntax error |
| -33002 | PackageNotFound | Bifrost package not found |
| -33003 | RuntimeError | Carrion runtime error |
| -33004 | DynamicLoadError | Dynamic loading failed |
| -33005 | TypeInferenceError | Type inference failed |

### Error Response Format

```json
{
  "error": {
    "code": -33001,
    "message": "Carrion syntax error",
    "data": {
      "line": 5,
      "character": 10,
      "details": "Expected ':' after spell definition",
      "suggestions": ["Add ':' at end of line"]
    }
  }
}
```

## Usage Examples

### Basic Completion Request

```bash
# Using curl to test completion
curl -X POST http://localhost:9999 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "textDocument/completion",
    "params": {
      "textDocument": {"uri": "file:///test.crl"},
      "position": {"line": 2, "character": 15}
    }
  }'
```

### Dynamic Package Loading

```bash
# Load a package dynamically
curl -X POST http://localhost:9999 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "carrion/loadPackage",
    "params": {
      "packageName": "json-utils"
    }
  }'
```

### Getting Available Packages

```bash
# Get list of available packages
curl -X POST http://localhost:9999 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "carrion/getAvailablePackages",
    "params": {}
  }'
```

## Performance Considerations

### Completion Performance

- **Response Time**: < 50ms for most completion requests
- **Memory Usage**: Symbol tables cached for active documents
- **Concurrent Requests**: Handled with read-write locks

### Dynamic Loading Performance

- **Initial Load**: < 500ms for full runtime discovery
- **Refresh Time**: < 100ms for incremental updates
- **Package Loading**: < 200ms for average-sized packages

### Optimization Tips

1. **Enable Caching**: Use `cacheSymbols: true` in configuration
2. **Limit Scope**: Exclude large directories with `excludePaths`
3. **Reduce Completions**: Set `maxItems` to reasonable number
4. **Disable Features**: Turn off `packageDiscovery` for faster startup

---

This API reference provides comprehensive information for integrating with and extending the Carrion-LSP server. For additional examples and tutorials, see the [documentation](../docs/) directory.