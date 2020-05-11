package expressions

import (
	"errors"
	"fmt"

	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/util"
)

func (e *Expression) functionCall(fname string) error {

	// validate this is a function
	if e.Tokens[e.TokenP] != "(" {
		return errors.New("invalid function call format")
	}

	e.TokenP = e.TokenP + 1
	argc := 0

	for e.Tokens[e.TokenP] != ")" {
		err := e.conditional()
		if err != nil {
			return err
		}
		argc = argc + 1
		if e.TokenP >= len(e.Tokens) {
			break
		}
		if e.Tokens[e.TokenP] == ")" {
			break
		}
		if e.Tokens[e.TokenP] != "," {
			return errors.New("invalid argument list")
		}
		e.TokenP = e.TokenP + 1
	}

	// Ensure trailing parenthesis
	if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ")" {
		return errors.New("mismatched parenthesis in argument list")
	}
	e.TokenP = e.TokenP + 1

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
func AddBuiltins(symbols map[string]interface{}) {

	for n, d := range util.FunctionDictionary {
		symbols[n+"()"] = d.F
	}
}
