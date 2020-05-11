package expressions

import (
	"errors"

	"github.com/tucats/gopackages/util"
)

func (e *Expression) functionCall(fname string, symbols map[string]interface{}) (interface{}, error) {

	// validate this is a function

	if e.Tokens[e.TokenP] != "(" {
		return nil, errors.New("invalid function call format")
	}
	f, found := symbols[fname+"()"]
	if !found {
		return nil, errors.New("function not found: " + fname)
	}

	// Parse the argument list, if any
	var args []interface{}

	e.TokenP = e.TokenP + 1

	for e.Tokens[e.TokenP] != ")" {
		v, err := e.conditional(symbols)
		if err != nil {
			return v, err
		}
		args = append(args, v)
		if e.TokenP >= len(e.Tokens) {
			break
		}
		if e.Tokens[e.TokenP] == ")" {
			break
		}
		if e.Tokens[e.TokenP] != "," {
			return nil, errors.New("invalid argument list")
		}
		e.TokenP = e.TokenP + 1
	}

	// Ensure trailing parenthesis
	if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ")" {
		return nil, errors.New("mismatched parenthesis in argument list")
	}
	e.TokenP = e.TokenP + 1

	// Quick sanity check on argument count for builtin functions
	fd, found := util.FunctionDictionary[fname]
	if found && ((len(args) < fd.Min) || (len(args) > fd.Max)) {
		return nil, errors.New("incorred number of arguments for " + fname + "()")
	}

	// Call the function
	return f.(func([]interface{}) (interface{}, error))(args)
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols map[string]interface{}) {

	for n, d := range util.FunctionDictionary {
		symbols[n+"()"] = d.F
	}
}
