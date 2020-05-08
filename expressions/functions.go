package expressions

import (
	"errors"
	"strings"
)

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
	min int
	max int
	f   interface{}
}

// FunctionDictionary is the dictionary of functions
var FunctionDictionary = map[string]FunctionDefinition{
	"int":       FunctionDefinition{min: 1, max: 1, f: functionInt},
	"bool":      FunctionDefinition{min: 1, max: 1, f: functionBool},
	"float":     FunctionDefinition{min: 1, max: 1, f: functionFloat},
	"string":    FunctionDefinition{min: 1, max: 1, f: functionString},
	"len":       FunctionDefinition{min: 1, max: 1, f: functionLen},
	"left":      FunctionDefinition{min: 2, max: 2, f: functionLeft},
	"right":     FunctionDefinition{min: 2, max: 2, f: functionRight},
	"substring": FunctionDefinition{min: 3, max: 3, f: functionSubstring},
	"index":     FunctionDefinition{min: 2, max: 2, f: functionIndex},
	"upper":     FunctionDefinition{min: 1, max: 1, f: functionUpper},
	"lower":     FunctionDefinition{min: 1, max: 1, f: functionLower},
	"min":       FunctionDefinition{min: 1, max: 99999, f: functionMin},
	"max":       FunctionDefinition{min: 1, max: 99999, f: functionMax},
	"sum":       FunctionDefinition{min: 1, max: 99999, f: functionSum},
}

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
	fd, found := FunctionDictionary[fname]
	if found && ((len(args) < fd.min) || (len(args) > fd.max)) {
		return nil, errors.New("incorred number of arguments for " + fname + "()")
	}

	// Call the function
	return f.(func([]interface{}) (interface{}, error))(args)
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols map[string]interface{}) {

	for n, d := range FunctionDictionary {
		symbols[n+"()"] = d.f
	}
}

//functionInt implements the int() function
func functionInt(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], 1)
	if v == nil {
		return nil, errors.New("invalid value to coerce to integer type")
	}
	return v.(int), nil
}

//functionFloat implements the float() function
func functionFloat(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], 1.0)
	if v == nil {
		return nil, errors.New("invalid value to coerce to float type")
	}
	return v.(float64), nil
}

//functionString implements the string() function
func functionString(args []interface{}) (interface{}, error) {

	return GetString(args[0]), nil
}

//functionBool implements the bool() function
func functionBool(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], true)
	if v == nil {
		return nil, errors.New("invalid value to coerce to bool type")
	}
	return v.(bool), nil
}

//functionLeft implements the left() function
func functionLeft(args []interface{}) (interface{}, error) {

	v := GetString(args[0])
	p := GetInt(args[1])

	if p <= 0 {
		return "", nil
	}
	if p >= len(v) {
		return v, nil
	}
	return v[:p], nil
}

//functionRight implements the right() function
func functionRight(args []interface{}) (interface{}, error) {
	v := GetString(args[0])
	p := GetInt(args[1])

	if p <= 0 {
		return "", nil
	}
	if p >= len(v) {
		return v, nil
	}
	return v[len(v)-p:], nil
}

// functionIndex implements the index() function
func functionIndex(args []interface{}) (interface{}, error) {

	v := GetString(args[0])
	p := GetString(args[1])

	return strings.Index(v, p) + 1, nil
}

// functionSubstring implements the substring() function
func functionSubstring(args []interface{}) (interface{}, error) {

	v := GetString(args[0])
	p1 := GetInt(args[1])
	p2 := GetInt(args[2])

	if p1 < 1 {
		p1 = 1
	}
	if p2 == 0 {
		return "", nil
	}
	if p2+p1 > len(v) {
		p2 = len(v) - p2
	}

	s := v[p1-1 : p1+p2-1]
	return s, nil
}

// functionLen implements the len() function
func functionLen(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], "")
	return len(v.(string)), nil

}

// functionLower implements the lower() function
func functionLower(args []interface{}) (interface{}, error) {
	return strings.ToLower(GetString(args[0])), nil
}

// functionUpper implements the upper() function
func functionUpper(args []interface{}) (interface{}, error) {
	return strings.ToUpper(GetString(args[0])), nil
}

// functionMin implements the min() function
func functionMin(args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		return args[0], nil
	}

	r := args[0]

	for _, v := range args[1:] {
		v = Coerce(v, r)
		switch r.(type) {
		case int:
			if v.(int) < r.(int) {
				r = v
			}

		case float64:
			if v.(float64) < r.(float64) {
				r = v
			}

		case string:
			if v.(string) < r.(string) {
				r = v
			}

		case bool:
			if v.(bool) == false {
				r = v
			}
		}
	}
	return r, nil
}

// functionMax implements the max() function
func functionMax(args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		return args[0], nil
	}

	r := args[0]

	for _, v := range args[1:] {
		v = Coerce(v, r)
		switch r.(type) {
		case int:
			if v.(int) > r.(int) {
				r = v
			}

		case float64:
			if v.(float64) > r.(float64) {
				r = v
			}

		case string:
			if v.(string) > r.(string) {
				r = v
			}

		case bool:
			if v.(bool) == true {
				r = v
			}
		}
	}
	return r, nil
}

// functionSum implements the sum() function
func functionSum(args []interface{}) (interface{}, error) {

	base := args[0]
	for _, addend := range args[1:] {
		addend = expressions.Coerce(addend, base)
		switch addend.(type) {
		case int:
			base = base.(int) + addend.(int)
		case float64:
			base = base.(float64) + addend.(float64)
		case string:
			base = base.(string) + addend.(string)

		case bool:
			base = base.(bool) || addend.(bool)
		}
	}
	return base, nil
}
