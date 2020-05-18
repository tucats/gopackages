package functions

import (
	"errors"

	"github.com/tucats/gopackages/util"
)

// FunctionInt implements the int() function
func FunctionInt(args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], 1)
	if v == nil {
		return nil, errors.New("invalid value to coerce to integer type")
	}
	return v.(int), nil
}

// FunctionFloat implements the float() function
func FunctionFloat(args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], 1.0)
	if v == nil {
		return nil, errors.New("invalid value to coerce to float type")
	}
	return v.(float64), nil
}

// FunctionString implements the string() function
func FunctionString(args []interface{}) (interface{}, error) {
	return util.GetString(args[0]), nil
}

// FunctionBool implements the bool() function
func FunctionBool(args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], true)
	if v == nil {
		return nil, errors.New("invalid value to coerce to bool type")
	}
	return v.(bool), nil
}
