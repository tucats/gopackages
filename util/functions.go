package util

import (
	"errors"
	"reflect"
	"strings"
)

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
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
	"left":      FunctionDefinition{Min: 2, Max: 2, F: FunctionLeft},
	"right":     FunctionDefinition{Min: 2, Max: 2, F: FunctionRight},
	"substring": FunctionDefinition{Min: 3, Max: 3, F: FunctionSubstring},
	"index":     FunctionDefinition{Min: 2, Max: 2, F: FunctionIndex},
	"upper":     FunctionDefinition{Min: 1, Max: 1, F: FunctionUpper},
	"lower":     FunctionDefinition{Min: 1, Max: 1, F: FunctionLower},
	"min":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMin},
	"max":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionMax},
	"sum":       FunctionDefinition{Min: 1, Max: 99999, F: FunctionSum},
}

// FunctionInt implements the int() function
func FunctionInt(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], 1)
	if v == nil {
		return nil, errors.New("invalid value to coerce to integer type")
	}
	return v.(int), nil
}

// FunctionFloat implements the float() function
func FunctionFloat(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], 1.0)
	if v == nil {
		return nil, errors.New("invalid value to coerce to float type")
	}
	return v.(float64), nil
}

// FunctionString implements the string() function
func FunctionString(args []interface{}) (interface{}, error) {

	return GetString(args[0]), nil
}

// FunctionBool implements the bool() function
func FunctionBool(args []interface{}) (interface{}, error) {

	v := Coerce(args[0], true)
	if v == nil {
		return nil, errors.New("invalid value to coerce to bool type")
	}
	return v.(bool), nil
}

// FunctionLeft implements the left() function
func FunctionLeft(args []interface{}) (interface{}, error) {

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

// FunctionRight implements the right() function
func FunctionRight(args []interface{}) (interface{}, error) {
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

// FunctionIndex implements the index() function
func FunctionIndex(args []interface{}) (interface{}, error) {

	switch arg := args[0].(type) {

	case []interface{}:
		for n, v := range arg {
			if reflect.DeepEqual(v, args[1]) {
				return n + 1, nil
			}
		}
		return 0, nil

	default:
		v := GetString(args[0])
		p := GetString(args[1])

		return strings.Index(v, p) + 1, nil
	}
}

// FunctionSubstring implements the substring() function
func FunctionSubstring(args []interface{}) (interface{}, error) {

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

// FunctionLen implements the len() function
func FunctionLen(args []interface{}) (interface{}, error) {

	switch arg := args[0].(type) {

	case map[string]interface{}:
		keys := make([]string, 0)
		for k := range arg {
			keys = append(keys, k)
		}
		return len(keys), nil

	case []interface{}:
		return len(arg), nil
	default:
		v := Coerce(args[0], "")
		return len(v.(string)), nil
	}
}

// FunctionLower implements the lower() function
func FunctionLower(args []interface{}) (interface{}, error) {
	return strings.ToLower(GetString(args[0])), nil
}

// FunctionUpper implements the upper() function
func FunctionUpper(args []interface{}) (interface{}, error) {
	return strings.ToUpper(GetString(args[0])), nil
}

// FunctionMin implements the min() function
func FunctionMin(args []interface{}) (interface{}, error) {

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

// FunctionMax implements the max() function
func FunctionMax(args []interface{}) (interface{}, error) {

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

// FunctionSum implements the sum() function
func FunctionSum(args []interface{}) (interface{}, error) {

	base := args[0]
	for _, addend := range args[1:] {
		addend = Coerce(addend, base)
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
