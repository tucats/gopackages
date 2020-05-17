package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/functions"
)

func (e *Expression) functionCall() error {

	// Note, caller already consumed the opening paren
	argc := 0

	for e.t.Peek(1) != ")" {
		err := e.conditional()
		if err != nil {
			return err
		}
		argc = argc + 1
		if e.t.AtEnd() {
			break
		}
		if e.t.Peek(1) == ")" {
			break
		}
		if e.t.Peek(1) != "," {
			return e.NewError("invalid argument list")
		}
		e.t.Advance(1)
	}

	// Ensure trailing parenthesis
	if e.t.AtEnd() || e.t.Peek(1) != ")" {
		return e.NewError("mismatched parenthesis in argument list")
	}
	e.t.Advance(1)

	// Call the function
	e.b.Emit2(bc.Call, argc)
	return nil
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *bc.SymbolTable) {

	for n, d := range functions.FunctionDictionary {

		if d.Pkg == "" {
			symbols.Set(n, d.F)
		} else {
			// Does package already exist? IF not, make it. The package
			// is just a struct containing where each member is a function
			// definition.
			p, found := symbols.Get(d.Pkg)
			if !found {
				p = map[string]interface{}{}
			}

			p.(map[string]interface{})[n] = d.F
			symbols.Set(d.Pkg, p)
		}
	}
}

// Function compiles a function call. The value of the
// function has been pushed to the top of the stack.
func (e *Expression) Function() error {

	// Get the atom
	err := e.reference()
	if err != nil {
		return err
	}
	// Peek ahead to see if it's the start of a function call...
	if e.t.IsNext("(") {
		return e.functionCall()
	}
	return nil
}
