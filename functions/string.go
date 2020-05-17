package functions

import (
	"errors"
	"reflect"
	"strings"

	"github.com/tucats/gopackages/util"
)

// FunctionLower implements the lower() function
func FunctionLower(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of function arguments")
	}

	return strings.ToLower(util.GetString(args[0])), nil
}

// FunctionUpper implements the upper() function
func FunctionUpper(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of function arguments")
	}
	return strings.ToUpper(util.GetString(args[0])), nil
}

// FunctionLeft implements the left() function
func FunctionLeft(args []interface{}) (interface{}, error) {

	if len(args) != 2 {
		return nil, errors.New("incorrect number of function arguments")
	}

	v := util.GetString(args[0])
	p := util.GetInt(args[1])

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
	if len(args) != 2 {
		return nil, errors.New("incorrect number of function arguments")
	}

	v := util.GetString(args[0])
	p := util.GetInt(args[1])

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
	if len(args) != 2 {
		return nil, errors.New("incorrect number of function arguments")
	}

	switch arg := args[0].(type) {

	case []interface{}:
		for n, v := range arg {
			if reflect.DeepEqual(v, args[1]) {
				return n + 1, nil
			}
		}
		return 0, nil

	case map[string]interface{}:
		key := util.GetString(args[1])
		_, found := arg[key]
		return found, nil

	default:
		v := util.GetString(args[0])
		p := util.GetString(args[1])

		return strings.Index(v, p) + 1, nil
	}
}

// FunctionSubstring implements the substring() function
func FunctionSubstring(args []interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, errors.New("incorrect number of function arguments")
	}

	v := util.GetString(args[0])
	p1 := util.GetInt(args[1])
	p2 := util.GetInt(args[2])

	if p1 < 1 {
		p1 = 1
	}
	if p2 == 0 {
		return "", nil
	}
	if p2+p1 > len(v) {
		p2 = len(v) - p1 + 1
	}

	s := v[p1-1 : p1+p2-1]
	return s, nil
}
