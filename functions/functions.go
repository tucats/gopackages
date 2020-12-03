package functions

import (
	"reflect"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
)

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
	Name string
	Pkg  string
	Min  int
	Max  int
	F    interface{}
}

// Any is a constant that defines that a function can have as many arguments
// as desired.
const Any = 999999

// FunctionDictionary is the dictionary of functions
var FunctionDictionary = map[string]FunctionDefinition{
	"array":             FunctionDefinition{Min: 1, Max: 2, F: Array},
	"bool":              FunctionDefinition{Min: 1, Max: 1, F: Bool},
	"error":             FunctionDefinition{Min: 1, Max: 1, F: Signal},
	"float":             FunctionDefinition{Min: 1, Max: 1, F: Float},
	"index":             FunctionDefinition{Min: 2, Max: 2, F: Index},
	"int":               FunctionDefinition{Min: 1, Max: 1, F: Int},
	"len":               FunctionDefinition{Min: 1, Max: 1, F: Length},
	"max":               FunctionDefinition{Min: 1, Max: Any, F: Max},
	"members":           FunctionDefinition{Min: 1, Max: 1, F: Members},
	"min":               FunctionDefinition{Min: 1, Max: Any, F: Min},
	"new":               FunctionDefinition{Min: 1, Max: 1, F: New},
	"sort":              FunctionDefinition{Min: 1, Max: Any, F: Sort},
	"string":            FunctionDefinition{Min: 1, Max: 1, F: String},
	"sum":               FunctionDefinition{Min: 1, Max: Any, F: Sum},
	"type":              FunctionDefinition{Min: 1, Max: 1, F: Type},
	"cipher.decrypt":    FunctionDefinition{Min: 2, Max: 2, F: Decrypt},
	"cipher.encrypt":    FunctionDefinition{Min: 2, Max: 2, F: Encrypt},
	"cipher.hash":       FunctionDefinition{Min: 1, Max: 1, F: Hash},
	"io.close":          FunctionDefinition{Min: 1, Max: 1, F: Close},
	"io.delete":         FunctionDefinition{Min: 1, Max: 1, F: DeleteFile},
	"io.expand":         FunctionDefinition{Min: 1, Max: 2, F: Expand},
	"io.open":           FunctionDefinition{Min: 1, Max: 2, F: Open},
	"io.readdir":        FunctionDefinition{Min: 1, Max: 1, F: ReadDir},
	"io.readfile":       FunctionDefinition{Min: 1, Max: 1, F: ReadFile},
	"io.readstring":     FunctionDefinition{Min: 1, Max: 1, F: ReadString},
	"io.split":          FunctionDefinition{Min: 1, Max: 1, F: Split},
	"io.writefile":      FunctionDefinition{Min: 2, Max: 2, F: WriteFile},
	"io.writestring":    FunctionDefinition{Min: 1, Max: 2, F: WriteString},
	"json.decode":       FunctionDefinition{Min: 1, Max: 1, F: Decode},
	"json.encode":       FunctionDefinition{Min: 1, Max: Any, F: Encode},
	"json.format":       FunctionDefinition{Min: 1, Max: Any, F: EncodeFormatted},
	"math.abs":          FunctionDefinition{Min: 1, Max: 1, F: Abs},
	"math.log":          FunctionDefinition{Min: 1, Max: 1, F: Log},
	"math.sqrt":         FunctionDefinition{Min: 1, Max: 1, F: Sqrt},
	"profile.delete":    FunctionDefinition{Min: 1, Max: 1, F: ProfileDelete},
	"profile.get":       FunctionDefinition{Min: 1, Max: 1, F: ProfileGet},
	"profile.keys":      FunctionDefinition{Min: 0, Max: 0, F: ProfileKeys},
	"profile.set":       FunctionDefinition{Min: 1, Max: 2, F: ProfileSet},
	"strings.chars":     FunctionDefinition{Min: 1, Max: 1, F: Chars},
	"strings.format":    FunctionDefinition{Min: 0, Max: Any, F: Format},
	"strings.index":     FunctionDefinition{Min: 2, Max: 2, F: Index},
	"strings.ints":      FunctionDefinition{Min: 1, Max: 1, F: Ints},
	"strings.left":      FunctionDefinition{Min: 2, Max: 2, F: Left},
	"strings.lower":     FunctionDefinition{Min: 1, Max: 1, F: Lower},
	"strings.right":     FunctionDefinition{Min: 2, Max: 2, F: Right},
	"strings.string":    FunctionDefinition{Min: 1, Max: Any, F: ToString},
	"strings.substring": FunctionDefinition{Min: 3, Max: 3, F: Substring},
	"strings.template":  FunctionDefinition{Min: 1, Max: 2, F: Template},
	"strings.tokenize":  FunctionDefinition{Min: 1, Max: 1, F: Tokenize},
	"strings.upper":     FunctionDefinition{Min: 1, Max: 1, F: Upper},
	"time.add":          FunctionDefinition{Min: 2, Max: 2, F: TimeAdd},
	"time.now":          FunctionDefinition{Min: 0, Max: 0, F: TimeNow},
	"time.sleep":        FunctionDefinition{Min: 1, Max: 1, F: Sleep},
	"time.subtract":     FunctionDefinition{Min: 2, Max: 2, F: TimeSub},
	"util.coerce":       FunctionDefinition{Min: 2, Max: 2, F: Coerce},
	"util.exit":         FunctionDefinition{Min: 0, Max: 1, F: Exit},
	"util.getenv":       FunctionDefinition{Min: 1, Max: 1, F: GetEnv},
	"util.normalize":    FunctionDefinition{Min: 2, Max: 2, F: Normalize},
	"util.symbols":      FunctionDefinition{Min: 0, Max: 1, F: FormatSymbols},
	"util.uuid":         FunctionDefinition{Min: 0, Max: 0, F: UUID},
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *symbols.SymbolTable) {

	ui.Debug("+++ Adding in builtin functions to symbol table %s", symbols.Name)
	for n, d := range FunctionDictionary {

		if dot := strings.Index(n, "."); dot >= 0 {
			d.Pkg = n[:dot]
			n = n[dot+1:]
		}

		if d.Pkg == "" {
			symbols.SetAlways(n, d.F)
		} else {
			// Does package already exist? IF not, make it. The package
			// is just a struct containing where each member is a function
			// definition.
			p, found := symbols.Get(d.Pkg)
			if !found {
				p = map[string]interface{}{}
				p.(map[string]interface{})["readonly"] = true
				ui.Debug("    AddBuiltins creating new package %s", d.Pkg)
			}
			p.(map[string]interface{})[n] = d.F
			ui.Debug("    adding builtin %s to %s", n, d.Pkg)
			symbols.SetAlways(d.Pkg, p)
		}
	}
}

// FindFunction returns the function definition associated with the
// provided function pointer, if one is found.
func FindFunction(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) *FunctionDefinition {

	sf1 := reflect.ValueOf(f)

	for _, d := range FunctionDictionary {
		sf2 := reflect.ValueOf(d.F)
		if sf1.Pointer() == sf2.Pointer() {
			return &d
		}
	}
	return nil
}

// FindName returns the name of a function from the dictionary if one is found
func FindName(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) string {

	sf1 := reflect.ValueOf(f)

	for name, d := range FunctionDictionary {
		sf2 := reflect.ValueOf(d.F)
		if sf1.Pointer() == sf2.Pointer() {
			return name
		}
	}

	return ""
}
