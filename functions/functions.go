package functions

import (
	"reflect"

	"github.com/tucats/gopackages/symbols"
)

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
	Pkg string
	Min int
	Max int
	F   interface{}
}

// FunctionDictionary is the dictionary of functions
var FunctionDictionary = map[string]FunctionDefinition{
	"int":       FunctionDefinition{Min: 1, Max: 1, F: FunctionInt},
	"bool":      FunctionDefinition{Min: 1, Max: 1, F: FunctionBool},
	"float":     FunctionDefinition{Min: 1, Max: 1, F: FunctionFloat},
	"string":    FunctionDefinition{Min: 1, Max: 1, F: FunctionString},
	"len":       FunctionDefinition{Min: 1, Max: 1, F: FunctionLen},
	"left":      FunctionDefinition{Min: 2, Max: 2, F: FunctionLeft, Pkg: "_strings"},
	"right":     FunctionDefinition{Min: 2, Max: 2, F: FunctionRight, Pkg: "_strings"},
	"substring": FunctionDefinition{Min: 3, Max: 3, F: FunctionSubstring, Pkg: "_strings"},
	"index":     FunctionDefinition{Min: 2, Max: 2, F: FunctionIndex},
	"upper":     FunctionDefinition{Min: 1, Max: 1, F: FunctionUpper, Pkg: "_strings"},
	"lower":     FunctionDefinition{Min: 1, Max: 1, F: FunctionLower, Pkg: "_strings"},
	"min":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMin},
	"max":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMax},
	"sum":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionSum},
	"uuid":      FunctionDefinition{Min: 0, Max: 0, F: FunctionUUID, Pkg: "_util"},
	"profile":   FunctionDefinition{Min: 1, Max: 2, F: FunctionProfile, Pkg: "_util"},
	"array":     FunctionDefinition{Min: 1, Max: 2, F: FunctionArray},
	"getenv":    FunctionDefinition{Min: 1, Max: 1, F: FunctionGetEnv, Pkg: "_util"},
	"members":   FunctionDefinition{Min: 1, Max: 1, F: FunctionMembers, Pkg: "_util"},
	"sqrt":      FunctionDefinition{Min: 1, Max: 1, F: FunctionSqrt, Pkg: "_math"},
	"sort":      FunctionDefinition{Min: 1, Max: 1, F: FunctionSort, Pkg: "_util"},
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *symbols.SymbolTable) {

	for n, d := range FunctionDictionary {

		if d.Pkg == "" {
			symbols.SetAlways(n, d.F)
		} else {
			// Does package already exist? IF not, make it. The package
			// is just a struct containing where each member is a function
			// definition.
			p, found := symbols.Get(d.Pkg)
			if !found {
				p = map[string]interface{}{}
				p.(map[string]interface{})["__readonly"] = true
			}

			p.(map[string]interface{})[n] = d.F
			symbols.SetAlways(d.Pkg, p)
		}
	}
}

// FindFunction returns the function definition associated with the
// provided function pointer, if one is found.
func FindFunction(f func([]interface{}) (interface{}, error)) *FunctionDefinition {

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
func FindName(f func([]interface{}) (interface{}, error)) string {

	sf1 := reflect.ValueOf(f)

	for name, d := range FunctionDictionary {
		sf2 := reflect.ValueOf(d.F)
		if sf1.Pointer() == sf2.Pointer() {
			return name
		}
	}

	return ""
}
