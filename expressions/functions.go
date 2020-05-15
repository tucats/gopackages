package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/functions"
)

func (e *Expression) functionCall(fname string) error {

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

	// Quick sanity check on argument count for builtin functions
	fd, found := functions.FunctionDictionary[fname]
	if found && ((argc < fd.Min) || (argc > fd.Max)) {
		return e.NewStringError("incorrect number of arguments for function", fname)
	}

	// Call the function
	e.b.Emit2(bc.Push, fname)
	e.b.Emit2(bc.Call, argc)
	return nil
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *bc.SymbolTable) {

	for n, d := range functions.FunctionDictionary {
		symbols.Set(n+"()", d.F)
	}
}
