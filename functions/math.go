package functions

import (
	"math"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// FunctionMin implements the min() function
func FunctionMin(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		return args[0], nil
	}

	r := args[0]

	for _, v := range args[1:] {
		v = util.Coerce(v, r)
		if v == nil {
			return nil, NewError("min", InvalidTypeError)
		}
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
		default:
			return nil, NewError("min", InvalidTypeError)

		}
	}
	return r, nil
}

// FunctionMax implements the max() function
func FunctionMax(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		return args[0], nil
	}

	r := args[0]

	for _, v := range args[1:] {
		v = util.Coerce(v, r)
		if v == nil {
			return nil, NewError("max", InvalidTypeError)
		}
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

		default:
			return nil, NewError("max", InvalidTypeError)
		}
	}
	return r, nil
}

// FunctionSum implements the sum() function
func FunctionSum(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	base := args[0]
	for _, addend := range args[1:] {
		addend = util.Coerce(addend, base)
		if addend == nil {
			return nil, NewError("sum", InvalidTypeError)
		}
		switch addend.(type) {
		case int:
			base = base.(int) + addend.(int)
		case float64:
			base = base.(float64) + addend.(float64)
		case string:
			base = base.(string) + addend.(string)

		case bool:
			base = base.(bool) || addend.(bool)
		default:
			return nil, NewError("sum", InvalidTypeError)

		}
	}
	return base, nil
}

// FunctionSqrt implements the sqrt() function
func FunctionSqrt(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	f := util.GetFloat(args[0])
	return math.Sqrt(f), nil
}

// FunctionAbs implements the abs() function
func FunctionAbs(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	f := util.GetFloat(args[0])
	return math.Abs(f), nil
}

// FunctionLog is the log() function
func FunctionLog(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return math.Log(util.GetFloat(args[0])), nil
}
