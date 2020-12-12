package bytecode

import (
	"math"
	"reflect"
	"strings"

	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*     M A T H   P R I M I T I V E S       *
*                                         *
\******************************************/

// NegateOpcode implementation
func NegateOpcode(c *Context, i interface{}) error {

	v, err := c.Pop()
	if err != nil {
		return err
	}

	switch value := v.(type) {

	case bool:
		_ = c.Push(!value)

	case int:
		_ = c.Push(-value)

	case float64:
		_ = c.Push(0.0 - value)

	case []interface{}:
		// Create an array in inverse order
		r := make([]interface{}, len(value))
		for n, d := range value {
			r[len(value)-n-1] = d
		}
		_ = c.Push(r)

	default:
		return c.NewError(InvalidTypeError)
	}
	return nil
}

// AddOpcode bytecode implementation
func AddOpcode(c *Context, i interface{}) error {

	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	switch vx := v1.(type) {

	// Is it an array we are concatenating to?
	case []interface{}:

		switch vy := v2.(type) {
		// Array requires a deep concatnation
		case []interface{}:
			newArray := append(vx, vy...)
			return c.Push(newArray)

		// Everything else is a simple append.
		default:
			newArray := append(vx, v2)
			return c.Push(newArray)
		}

		// You can add a map to another map
	case map[string]interface{}:

		switch vy := v2.(type) {
		case map[string]interface{}:
			for k, v := range vy {
				vx[k] = v
			}
			return c.Push(vx)

		default:
			return c.NewError(InvalidTypeError)
		}

		// All other types are scalar math
	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			return c.Push(v1.(int) + v2.(int))
		case float64:
			return c.Push(v1.(float64) + v2.(float64))
		case string:
			return c.Push(v1.(string) + v2.(string))
		case bool:
			return c.Push(v1.(bool) && v2.(bool))
		default:
			return c.NewError(InvalidTypeError)
		}
	}
}

// AndOpcode bytecode implementation
func AndOpcode(c *Context, i interface{}) error {

	v1, err := c.Pop()
	if err != nil {
		return err
	}
	v2, err := c.Pop()
	if err != nil {
		return err
	}

	return c.Push(util.GetBool(v1) && util.GetBool(v2))

}

// OrOpcode bytecode implementation
func OrOpcode(c *Context, i interface{}) error {

	v1, err := c.Pop()
	if err != nil {
		return err
	}
	v2, err := c.Pop()
	if err != nil {
		return err
	}

	return c.Push(util.GetBool(v1) || util.GetBool(v2))

}

// SubOpcode bytecode implementation
func SubOpcode(c *Context, i interface{}) error {
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	switch vx := v1.(type) {

	// For an array, make a copy removing the item to be subtracted.
	case []interface{}:
		newArray := make([]interface{}, 0)
		for _, v := range vx {
			if !reflect.DeepEqual(v2, v) {
				newArray = append(newArray, v)
			}
		}
		return c.Push(newArray)

	// Everything else is a scalar subtraction
	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			return c.Push(v1.(int) - v2.(int))
		case float64:
			return c.Push(v1.(float64) - v2.(float64))
		case string:
			s := strings.ReplaceAll(v1.(string), v2.(string), "")
			return c.Push(s)
		default:
			return c.NewError(InvalidTypeError)
		}
	}
}

// MulOpcode bytecode implementation
func MulOpcode(c *Context, i interface{}) error {

	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	v1, v2 = util.Normalize(v1, v2)
	switch v1.(type) {
	case int:
		return c.Push(v1.(int) * v2.(int))
	case float64:
		return c.Push(v1.(float64) * v2.(float64))
	case bool:
		return c.Push(v1.(bool) || v2.(bool))
	default:
		return c.NewError(InvalidTypeError)
	}
}

// ExpOpcode bytecode implementation
func ExpOpcode(c *Context, i interface{}) error {

	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	v1, v2 = util.Normalize(v1, v2)
	switch v1.(type) {
	case int:
		if v2.(int) == 0 {
			return c.Push(0)
		}
		if v2.(int) == 1 {
			return c.Push(v1)
		}
		prod := v1.(int)
		for n := 2; n <= v2.(int); n = n + 1 {
			prod = prod * v1.(int)
		}
		return c.Push(prod)

	case float64:
		return c.Push(math.Pow(v1.(float64), v2.(float64)))
	default:
		return c.NewError(InvalidTypeError)
	}
}

// DivOpcode bytecode implementation
func DivOpcode(c *Context, i interface{}) error {

	if c.sp < 1 {
		return c.NewError(StackUnderflowError)
	}
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	v1, v2 = util.Normalize(v1, v2)
	switch v1.(type) {
	case int:
		if v2.(int) == 0 {
			return c.NewError(DivisionByZeroError)
		}
		return c.Push(v1.(int) / v2.(int))
	case float64:
		if v2.(float64) == 0 {
			return c.NewError(DivisionByZeroError)
		}
		return c.Push(v1.(float64) / v2.(float64))
	default:
		return c.NewError(InvalidTypeError)
	}
}
