package functions

import (
	"errors"
	"fmt"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// FunctionInt implements the int() function
func FunctionInt(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], 1)
	if v == nil {
		return nil, errors.New("invalid value to coerce to integer type")
	}
	return v.(int), nil
}

// FunctionFloat implements the float() function
func FunctionFloat(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], 1.0)
	if v == nil {
		return nil, errors.New("invalid value to coerce to float type")
	}
	return v.(float64), nil
}

// FunctionString implements the string() function
func FunctionString(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return util.GetString(args[0]), nil
}

// FunctionBool implements the bool() function
func FunctionBool(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	v := util.Coerce(args[0], true)
	if v == nil {
		return nil, errors.New("invalid value to coerce to bool type")
	}
	return v.(bool), nil
}

// FunctionNew implements the new() function
func FunctionNew(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	r := DeepCopy(args[0])

	// IF there was a type in the source, make the clone point back to it

	switch v := r.(type) {

	case nil:
		return nil, errors.New("cannot use nil as arg to new()")

	case symbols.SymbolTable:
		return nil, errors.New("cannot new() a symbol table")

	case func(*symbols.SymbolTable, []interface{}) (interface{}, error):
		return nil, errors.New("cannot new() a native function")

	case int:
	case string:
	case float64:
	case []interface{}:

	case map[string]interface{}:
		if _, found := v["__parent"]; found {
			r.(map[string]interface{})["__parent"] = args[0]
		}

	default:
		return nil, errors.New("unsupported new() type " + fmt.Sprintf("%#v", v))
	}

	return r, nil
}

// DeepCopy makes a deep copy of a Solve data type
func DeepCopy(source interface{}) interface{} {

	switch v := source.(type) {

	case int:
		return v
	case string:
		return v
	case float64:
		return v
	case bool:
		return v

	case []interface{}:

		r := make([]interface{}, 0)
		for _, d := range v {
			r = append(r, DeepCopy(d))
		}
		return r

	case map[string]interface{}:
		r := map[string]interface{}{}
		for k, d := range v {
			r[k] = DeepCopy(d)
		}
		return r

	default:
		ui.Debug("DeepCopy of uncopyable type: %#v\n", v)
		return v
	}
}
