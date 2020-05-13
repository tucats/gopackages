package expressions

import (
	"errors"
	"fmt"

	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/util"
)

func (e *Expression) functionCall(fname string) error {

	// Note, caller already consumed the opening paren
	argc := 0

	for e.t.Peek() != ")" {
		err := e.conditional()
		if err != nil {
			return err
		}
		argc = argc + 1
		if e.t.AtEnd() {
			break
		}
		if e.t.Peek() == ")" {
			break
		}
		if e.t.Peek() != "," {
			return errors.New("invalid argument list")
		}
		e.t.Advance(1)
	}

	// Ensure trailing parenthesis
	if e.t.AtEnd() || e.t.Peek() != ")" {
		return errors.New("mismatched parenthesis in argument list")
	}
	e.t.Advance(1)

	// Quick sanity check on argument count for builtin functions
	fd, found := util.FunctionDictionary[fname]
	if found && ((argc < fd.Min) || (argc > fd.Max)) {
		return fmt.Errorf("incorred number of arguments for %s()", fname)
	}

	// Call the function
	e.b.Emit(bc.Push, fname)
	e.b.Emit(bc.Call, argc)
	return nil
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *bc.SymbolTable) {

	for n, d := range util.FunctionDictionary {
		symbols.Set(n+"()", d.F)
	}
}
