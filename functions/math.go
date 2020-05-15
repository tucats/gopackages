package functions

import "github.com/tucats/gopackages/util"

// FunctionMin implements the min() function
func FunctionMin(args []interface{}) (interface{}, error) {

	if len(args) == 1 {
		return args[0], nil
	}

	r := args[0]

	for _, v := range args[1:] {
		v = util.Coerce(v, r)
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
		v = util.Coerce(v, r)
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
		addend = util.Coerce(addend, base)
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
