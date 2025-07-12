package analyzer

import (
	"fmt"
	"strings"

	"github.com/javanhut/TheCarrionLanguage/src/ast"
	"github.com/javanhut/TheCarrionLanguage/src/evaluator"
	"github.com/javanhut/TheCarrionLanguage/src/lexer"
	"github.com/javanhut/TheCarrionLanguage/src/object"
	"github.com/javanhut/TheCarrionLanguage/src/parser"
)

// DynamicLoader provides dynamic loading of Carrion runtime components
type DynamicLoader struct {
	env       *object.Environment
	builtins  map[string]*BuiltinInfo
	grimoires map[string]*GrimoireInfo
}

// NewDynamicLoader creates a new dynamic loader with a Carrion environment
func NewDynamicLoader() *DynamicLoader {
	env := object.NewEnvironment()

	// Load the munin standard library (includes all grimoires and modules)
	err := evaluator.LoadMuninStdlib(env)
	if err != nil {
		// Fallback to empty environment if loading fails
		fmt.Printf("Warning: Failed to load munin stdlib: %v\n", err)
	}

	loader := &DynamicLoader{
		env:       env,
		builtins:  make(map[string]*BuiltinInfo),
		grimoires: make(map[string]*GrimoireInfo),
	}

	loader.loadBuiltins()
	loader.loadGrimoires()

	return loader
}

// loadBuiltins discovers built-in functions from the Carrion runtime
func (dl *DynamicLoader) loadBuiltins() {
	// Get all built-ins from the evaluator
	runtimeBuiltins := evaluator.GetBuiltins()

	for name := range runtimeBuiltins {
		// Extract function information
		builtinInfo := &BuiltinInfo{
			Name:        name,
			Type:        "function",
			Description: dl.getBuiltinDescription(name),
			Parameters:  dl.extractBuiltinParameters(name),
			ReturnType:  dl.inferBuiltinReturnType(name),
		}

		dl.builtins[name] = builtinInfo
	}

	// Also check environment for additional functions loaded by modules
	dl.loadEnvironmentBuiltins()
}

// loadEnvironmentBuiltins discovers functions loaded into the environment by modules
func (dl *DynamicLoader) loadEnvironmentBuiltins() {
	// Get all symbols from the environment
	envStore := dl.env.GetStore()

	for name, obj := range envStore {
		switch typedObj := obj.(type) {
		case *object.Builtin:
			if _, exists := dl.builtins[name]; !exists {
				dl.builtins[name] = &BuiltinInfo{
					Name:        name,
					Type:        "function",
					Description: dl.getBuiltinDescription(name),
					Parameters:  dl.extractBuiltinParameters(name),
					ReturnType:  dl.inferBuiltinReturnType(name),
				}
			}
		case *object.Function:
			// Handle user-defined functions from modules
			if _, exists := dl.builtins[name]; !exists {
				dl.builtins[name] = &BuiltinInfo{
					Name:        name,
					Type:        "function",
					Description: fmt.Sprintf("Module function: %s", name),
					Parameters:  dl.extractFunctionParameters(typedObj),
					ReturnType:  "unknown",
				}
			}
		}
	}
}

// loadGrimoires discovers grimoire classes from the runtime environment
func (dl *DynamicLoader) loadGrimoires() {
	envStore := dl.env.GetStore()

	for name, obj := range envStore {
		if grimoire, isGrimoire := obj.(*object.Grimoire); isGrimoire {
			grimoireInfo := &GrimoireInfo{
				Name:        name,
				Description: dl.getGrimoireDescription(name),
				Spells:      make(map[string]*BuiltinInfo),
				IsStatic:    grimoire.IsArcane,
			}

			// Extract spells/methods from the grimoire
			for spellName, spellFunc := range grimoire.Methods {
				spell := &BuiltinInfo{
					Name:        spellName,
					Type:        "method",
					Description: dl.getSpellDescription(name, spellName),
					Parameters:  dl.extractFunctionParameters(spellFunc),
					ReturnType:  dl.inferSpellReturnType(name, spellName),
				}
				grimoireInfo.Spells[spellName] = spell
			}

			dl.grimoires[name] = grimoireInfo
		}
	}
}

// getBuiltinDescription provides descriptions for built-in functions
func (dl *DynamicLoader) getBuiltinDescription(name string) string {
	descriptions := map[string]string{
		"print":     "Print values to output",
		"input":     "Read user input with optional prompt",
		"len":       "Get length of strings, arrays, or hashes",
		"type":      "Get the type of an object",
		"range":     "Generate a sequence of numbers",
		"int":       "Convert to integer",
		"float":     "Convert to float",
		"str":       "Convert to string",
		"bool":      "Convert to boolean",
		"list":      "Convert to array",
		"open":      "Open a file and return File grimoire instance",
		"max":       "Find maximum value",
		"abs":       "Get absolute value",
		"enumerate": "Enumerate arrays with indices",
		"pairs":     "Extract key-value pairs from hashes",
		// Time module functions
		"time_now":    "Get current Unix timestamp",
		"time_sleep":  "Sleep for specified duration",
		"time_format": "Format timestamp to string",
		"time_parse":  "Parse time string to timestamp",
		// File module functions
		"file_read":   "Read file content",
		"file_write":  "Write content to file",
		"file_exists": "Check if file exists",
		// OS module functions
		"os_cwd":     "Get current working directory",
		"os_listdir": "List directory contents",
		"os_mkdir":   "Create directory",
		"os_getenv":  "Get environment variable",
		"os_run":     "Run system command",
		// HTTP module functions
		"http_get":    "HTTP GET request",
		"http_post":   "HTTP POST request",
		"http_put":    "HTTP PUT request",
		"http_delete": "HTTP DELETE request",
	}

	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return fmt.Sprintf("Built-in function: %s", name)
}

// getGrimoireDescription provides descriptions for grimoires
func (dl *DynamicLoader) getGrimoireDescription(name string) string {
	descriptions := map[string]string{
		"String":  "String manipulation grimoire",
		"Array":   "Array manipulation grimoire",
		"Integer": "Integer operations grimoire",
		"Float":   "Float operations grimoire",
		"Boolean": "Boolean operations grimoire",
		"File":    "File operations grimoire",
		"OS":      "Operating system operations grimoire",
		"Time":    "Time operations grimoire",
	}

	if desc, exists := descriptions[name]; exists {
		return desc
	}
	return fmt.Sprintf("Grimoire: %s", name)
}

// getSpellDescription provides descriptions for grimoire spells
func (dl *DynamicLoader) getSpellDescription(grimoire, spell string) string {
	// This could be enhanced to read from docstrings or comments
	return fmt.Sprintf("%s method from %s grimoire", spell, grimoire)
}

// extractBuiltinParameters attempts to extract parameter information for built-ins
func (dl *DynamicLoader) extractBuiltinParameters(name string) []Parameter {
	// Static parameter definitions for known built-ins
	// This could be enhanced to use reflection or documentation parsing
	paramMap := map[string][]Parameter{
		"print": {{Name: "values", TypeHint: "...any"}},
		"input": {{Name: "prompt", TypeHint: "string", DefaultValue: "\"\""}},
		"len":   {{Name: "obj", TypeHint: "any"}},
		"type":  {{Name: "obj", TypeHint: "any"}},
		"range": {
			{Name: "start", TypeHint: "int"},
			{Name: "stop", TypeHint: "int", DefaultValue: "None"},
			{Name: "step", TypeHint: "int", DefaultValue: "1"},
		},
		"int":   {{Name: "value", TypeHint: "any"}},
		"float": {{Name: "value", TypeHint: "any"}},
		"str":   {{Name: "value", TypeHint: "any"}},
		"bool":  {{Name: "value", TypeHint: "any"}},
		"open": {
			{Name: "path", TypeHint: "string"},
			{Name: "mode", TypeHint: "string", DefaultValue: "\"r\""},
		},
		"max":       {{Name: "values", TypeHint: "...any"}},
		"abs":       {{Name: "value", TypeHint: "number"}},
		"enumerate": {{Name: "array", TypeHint: "array"}},
		"pairs":     {{Name: "hash", TypeHint: "hash"}},
	}

	if params, exists := paramMap[name]; exists {
		return params
	}
	return []Parameter{} // Empty if unknown
}

// extractFunctionParameters extracts parameters from a Function object
func (dl *DynamicLoader) extractFunctionParameters(fn *object.Function) []Parameter {
	var parameters []Parameter

	if fn.Parameters != nil {
		for _, param := range fn.Parameters {
			// Handle different parameter types
			switch p := param.(type) {
			case *ast.Identifier:
				parameters = append(parameters, Parameter{
					Name:     p.Value,
					TypeHint: "any",
				})
			case *ast.Parameter:
				parameter := Parameter{
					Name:     p.Name.Value,
					TypeHint: "any",
				}
				if p.TypeHint != nil {
					parameter.TypeHint = p.TypeHint.String()
				}
				if p.DefaultValue != nil {
					parameter.DefaultValue = p.DefaultValue.String()
				}
				parameters = append(parameters, parameter)
			default:
				// Fallback for other parameter types
				parameters = append(parameters, Parameter{
					Name:     "param",
					TypeHint: "any",
				})
			}
		}
	}

	return parameters
}

// inferBuiltinReturnType attempts to infer return types for built-ins
func (dl *DynamicLoader) inferBuiltinReturnType(name string) string {
	returnTypes := map[string]string{
		"print":     "None",
		"input":     "string",
		"len":       "int",
		"type":      "string",
		"range":     "array",
		"int":       "int",
		"float":     "float",
		"str":       "string",
		"bool":      "bool",
		"list":      "array",
		"open":      "File",
		"max":       "any",
		"abs":       "number",
		"enumerate": "array",
		"pairs":     "array",
	}

	if returnType, exists := returnTypes[name]; exists {
		return returnType
	}
	return "unknown"
}

// inferSpellReturnType attempts to infer return types for grimoire spells
func (dl *DynamicLoader) inferSpellReturnType(grimoire, spell string) string {
	// Enhanced type inference could be added here
	// For now, use some common patterns
	if strings.Contains(spell, "is_") || strings.Contains(spell, "contains") {
		return "bool"
	}
	if strings.Contains(spell, "length") || strings.Contains(spell, "count") {
		return "int"
	}
	if grimoire == "String" && (spell == "lower" || spell == "upper" || spell == "reverse") {
		return "string"
	}
	return "unknown"
}

// GetBuiltins returns the dynamically loaded built-ins
func (dl *DynamicLoader) GetBuiltins() map[string]*BuiltinInfo {
	return dl.builtins
}

// GetGrimoires returns the dynamically loaded grimoires
func (dl *DynamicLoader) GetGrimoires() map[string]*GrimoireInfo {
	return dl.grimoires
}

// LoadBifrostPackage attempts to load a bifrost package dynamically
func (dl *DynamicLoader) LoadBifrostPackage(packagePath string) error {
	// Read the package's main file
	// This is a placeholder - would need to integrate with bifrost's loading system

	// For now, try to parse and evaluate a .crl file
	code := fmt.Sprintf(`import "%s"`, packagePath)

	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return fmt.Errorf("parse errors: %v", p.Errors())
	}

	// Evaluate in our environment
	evaluator.Eval(program, dl.env, nil)

	// Reload grimoires and builtins to pick up new definitions
	dl.loadBuiltins()
	dl.loadGrimoires()

	return nil
}

// RefreshDynamicData reloads all dynamic data from the runtime
func (dl *DynamicLoader) RefreshDynamicData() {
	dl.builtins = make(map[string]*BuiltinInfo)
	dl.grimoires = make(map[string]*GrimoireInfo)

	dl.loadBuiltins()
	dl.loadGrimoires()
}
