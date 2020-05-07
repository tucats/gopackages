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
	symbols["int()"] = functionInt
	symbols["float()"] = functionFloat
	symbols["string()"] = functionString
	symbols["bool()"] = functionBool
}

// functionLen implements the len() function
func functionLen(args []interface{}) (interface{}, error) {

	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to len() function")
	}

	v := Coerce(args[0], "")
	return len(v.(string)), nil

}

//functionInt implements the int() function
func functionInt(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to int() function")
	}

	v := Coerce(args[0], 1)
	if v == nil {
		return nil, errors.New("invalid value to coerce to integer type")
	}
	return v.(int), nil
}

//functionFloat implements the float() function
func functionFloat(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to float() function")
	}

	v := Coerce(args[0], 1.0)
	if v == nil {
		return nil, errors.New("invalid value to coerce to float type")
	}
	return v.(float64), nil
}

//functionString implements the string() function
func functionString(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to int() function")
	}

	v := Coerce(args[0], "")
	return v.(string), nil
}

//functionBool implements the bool() function
func functionBool(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of arguments to bool() function")
	}

	v := Coerce(args[0], true)
	if v == nil {
		return nil, errors.New("invalid value to coerce to bool type")
	}
	return v.(bool), nil
}
