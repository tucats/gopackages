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
	"int":               FunctionDefinition{Min: 1, Max: 1, F: FunctionInt},
	"bool":              FunctionDefinition{Min: 1, Max: 1, F: FunctionBool},
	"float":             FunctionDefinition{Min: 1, Max: 1, F: FunctionFloat},
	"string":            FunctionDefinition{Min: 1, Max: 1, F: FunctionString},
	"len":               FunctionDefinition{Min: 1, Max: 1, F: FunctionLen},
	"index":             FunctionDefinition{Min: 2, Max: 2, F: FunctionIndex},
	"min":               FunctionDefinition{Min: 1, Max: Any, F: FunctionMin},
	"max":               FunctionDefinition{Min: 1, Max: Any, F: FunctionMax},
	"sum":               FunctionDefinition{Min: 1, Max: Any, F: FunctionSum},
	"array":             FunctionDefinition{Min: 1, Max: 2, F: FunctionArray},
	"new":               FunctionDefinition{Min: 1, Max: 1, F: FunctionNew},
	"math.sqrt":         FunctionDefinition{Min: 1, Max: 1, F: FunctionSqrt},
	"strings.left":      FunctionDefinition{Min: 2, Max: 2, F: FunctionLeft},
	"strings.right":     FunctionDefinition{Min: 2, Max: 2, F: FunctionRight},
	"strings.substring": FunctionDefinition{Min: 3, Max: 3, F: FunctionSubstring},
	"strings.upper":     FunctionDefinition{Min: 1, Max: 1, F: FunctionUpper},
	"strings.lower":     FunctionDefinition{Min: 1, Max: 1, F: FunctionLower},
	"strings.format":    FunctionDefinition{Min: 0, Max: Any, F: FunctionFormat},
	"strings.tokenize":  FunctionDefinition{Min: 1, Max: 1, F: FunctionTokenize},
	"strings.chars":     FunctionDefinition{Min: 1, Max: 1, F: FunctionChars},
	"strings.ints":      FunctionDefinition{Min: 1, Max: 1, F: FunctionInts},
	"strings.string":    FunctionDefinition{Min: 1, Max: Any, F: FunctionToString},
	"strings.index":     FunctionDefinition{Min: 2, Max: 2, F: FunctionIndex},
	"strings.template":  FunctionDefinition{Min: 1, Max: 2, F: FunctionTemplate},
	"util.uuid":         FunctionDefinition{Min: 0, Max: 0, F: FunctionUUID},
	"util.profile":      FunctionDefinition{Min: 1, Max: 2, F: FunctionProfile},
	"util.getenv":       FunctionDefinition{Min: 1, Max: 1, F: FunctionGetEnv},
	"members":           FunctionDefinition{Min: 1, Max: 1, F: FunctionMembers},
	"sort":              FunctionDefinition{Min: 1, Max: Any, F: FunctionSort},
	"util.exit":         FunctionDefinition{Min: 0, Max: 1, F: FunctionExit},
	"util.symbols":      FunctionDefinition{Min: 0, Max: 1, F: FunctionSymbols},
	"type":              FunctionDefinition{Min: 1, Max: 1, F: FunctionType},
	"io.open":           FunctionDefinition{Min: 1, Max: 2, F: FunctionOpen},
	"io.close":          FunctionDefinition{Min: 1, Max: 1, F: FunctionClose},
	"io.readfile":       FunctionDefinition{Min: 1, Max: 1, F: FunctionReadFile},
	"io.readstring":     FunctionDefinition{Min: 1, Max: 1, F: FunctionReadString},
	"io.writestring":    FunctionDefinition{Min: 1, Max: 2, F: FunctionWriteString},
	"io.split":          FunctionDefinition{Min: 1, Max: 1, F: FunctionSplit},
	"io.writefile":      FunctionDefinition{Min: 2, Max: 2, F: FunctionWriteFile},
	"io.delete":         FunctionDefinition{Min: 1, Max: 1, F: FunctionDeleteFile},
	"json.decode":       FunctionDefinition{Min: 1, Max: 1, F: FunctionDecode},
	"json.encode":       FunctionDefinition{Min: 1, Max: Any, F: FunctionEncode},
	"json.format":       FunctionDefinition{Min: 1, Max: Any, F: FunctionEncodeFormatted},
	"cipher.hash":       FunctionDefinition{Min: 1, Max: 1, F: FunctionHash},
	"cipher.encrypt":    FunctionDefinition{Min: 2, Max: 2, F: FunctionEncrypt},
	"cipher.decrypt":    FunctionDefinition{Min: 2, Max: 2, F: FunctionDecrypt},
	"time.now":          FunctionDefinition{Min: 0, Max: 0, F: FunctionTimeNow},
	"time.add":          FunctionDefinition{Min: 2, Max: 2, F: FunctionTimeAdd},
	"time.subtract":     FunctionDefinition{Min: 2, Max: 2, F: FunctionTimeSub},
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
