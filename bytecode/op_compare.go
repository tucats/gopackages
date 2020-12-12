package bytecode

import (
	"reflect"

	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*   C O M P A R E   O P E R A T I O N S   *
*                                         *
\******************************************/

// EqualOpcode implementation
func EqualOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}

	v1, err := c.Pop()
	if err != nil {
		return err
	}
	var r bool

	switch v1.(type) {

	case map[string]interface{}:
		r = reflect.DeepEqual(v1, v2)

	case []interface{}:
		r = reflect.DeepEqual(v1, v2)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) == v2.(int)
		case float64:
			r = v1.(float64) == v2.(float64)
		case string:
			r = v1.(string) == v2.(string)
		case bool:
			r = v1.(bool) == v2.(bool)

		}
	}

	_ = c.Push(r)
	return nil

}

// NotEqualOpcode implementation
func NotEqualOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}

	v1, err := c.Pop()
	if err != nil {
		return err
	}

	var r bool

	switch v1.(type) {

	case map[string]interface{}:
		r = !reflect.DeepEqual(v1, v2)

	case []interface{}:
		r = !reflect.DeepEqual(v1, v2)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) != v2.(int)
		case float64:
			r = v1.(float64) != v2.(float64)
		case string:
			r = v1.(string) != v2.(string)
		case bool:
			r = v1.(bool) != v2.(bool)

		}
	}

	_ = c.Push(r)
	return nil

}

// GreaterThanOpcode implementation
func GreaterThanOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	var r bool

	switch v1.(type) {

	case []interface{}:
		return c.NewError(InvalidTypeError)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) > v2.(int)
		case float64:
			r = v1.(float64) > v2.(float64)
		case string:
			r = v1.(string) > v2.(string)

		default:
			return c.NewError(InvalidTypeError)

		}
	}
	_ = c.Push(r)
	return nil
}

// GreaterThanOrEqualOpcode implementation
func GreaterThanOrEqualOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	var r bool

	switch v1.(type) {

	case []interface{}:
		return c.NewError(InvalidTypeError)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) >= v2.(int)
		case float64:
			r = v1.(float64) >= v2.(float64)
		case string:
			r = v1.(string) >= v2.(string)

		default:
			return c.NewError(InvalidTypeError)

		}
	}
	_ = c.Push(r)
	return nil
}

// LessThanOpcode implementation
func LessThanOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	// Handle nil cases
	if v1 == nil && v2 == nil {
		_ = c.Push(true)
		return nil
	}
	if v1 == nil || v2 == nil {
		_ = c.Push(false)
		return nil
	}

	// Nope, going to have to do type-sensitive compares.
	var r bool

	switch v1.(type) {

	case []interface{}:
		return c.NewError(InvalidTypeError)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) < v2.(int)
		case float64:
			r = v1.(float64) < v2.(float64)
		case string:
			r = v1.(string) < v2.(string)

		default:
			return c.NewError(InvalidTypeError)

		}
	}
	_ = c.Push(r)
	return nil
}

// LessThanOrEqualOpcode implementation
func LessThanOrEqualOpcode(c *Context, i interface{}) error {

	// Terms pushed in reverse order
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	v1, err := c.Pop()
	if err != nil {
		return err
	}

	var r bool

	switch v1.(type) {

	case []interface{}:
		return c.NewError(InvalidTypeError)

	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			r = v1.(int) <= v2.(int)
		case float64:
			r = v1.(float64) <= v2.(float64)
		case string:
			r = v1.(string) <= v2.(string)

		default:
			return c.NewError(InvalidTypeError)

		}
	}
	_ = c.Push(r)
	return nil
}
