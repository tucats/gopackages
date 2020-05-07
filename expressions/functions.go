package expressions

import "errors"

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
		v, err := e.relations(symbols)
		if err != nil {
			return v, err
		}
		args = append(args, v)
	}

	// Ensure trailing parenthesis
	if e.Tokens[e.TokenP] != ")" {
		return nil, errors.New("mismatched parenthesis in argument list")
	}
	e.TokenP = e.TokenP + 1
	return f.(func([]interface{}) (interface{}, error))(args)
}

// AddBuiltins adds or overrides the default function library to the symbol map
func AddBuiltins(symbols map[string]interface{}) {
	symbols["len()"] = functionLen
}

// functionLen implements the len() function
func functionLen(args []interface{}) (interface{}, error) {

	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to len() function")
	}

	v := Coerce(args[0], "")
	return len(v.(string)), nil

}
