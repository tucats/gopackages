package bytecode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/util"
)

// StopOpcode bytecode implementation
func StopOpcode(c *Context, i interface{}) error {
	c.running = false
	return nil
}

// PrintOpcode implementation. If the operand
// is given, it represents the number of items
// to remove from the stack.
func PrintOpcode(c *Context, i interface{}) error {

	count := 1
	if i != nil {
		count = util.GetInt(i)
	}

	for n := 0; n < count; n = n + 1 {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		fmt.Printf("%s", util.FormatUnquoted(v))
	}
	return nil
}

// NewlineOpcode implementation.
func NewlineOpcode(c *Context, i interface{}) error {
	fmt.Printf("\n")
	return nil
}

// ArrayOpcode implementation
func ArrayOpcode(c *Context, i interface{}) error {

	count := util.GetInt(i)
	array := make([]interface{}, count)

	for n := 0; n < count; n++ {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		array[(count-n)-1] = v
	}

	c.Push(array)
	return nil
}

// StructOpcode implementation. The operand is a count
// of elements on the stack. These are pulled off in pairs,
// where the first value is the name of the struct field and
// the second value is the value of the struct field.
func StructOpcode(c *Context, i interface{}) error {

	count := util.GetInt(i)

	m := map[string]interface{}{}

	for n := 0; n < count; n++ {
		name, err := c.Pop()
		if err != nil {
			return err
		}
		value, err := c.Pop()
		if err != nil {
			return err
		}
		m[util.GetString(name)] = value
	}

	c.Push(m)
	return nil
}

// MemberOpcode implementation. This pops two values from
// the stack (the first must be a string and the second a
// map) and indexes into the map to get the matching value
// and puts back on the stack.
func MemberOpcode(c *Context, i interface{}) error {

	var name string
	if i != nil {
		name = util.GetString(i)
	} else {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		name = util.GetString(v)
	}

	m, err := c.Pop()
	if err != nil {
		return err
	}

	// The only the type that is supported is a map
	switch mv := m.(type) {
	case map[string]interface{}:
		v, found := mv[name]
		if !found {
			return fmt.Errorf("no such member: %s", name)
		}
		c.Push(v)

	default:
		return errors.New("not a map")
	}
	return nil
}

// StoreOpcode implementation
func StoreOpcode(c *Context, i interface{}) error {

	v, err := c.Pop()
	if err != nil {
		return err
	}

	c.Set(util.GetString(i), v)
	return nil
}

// LoadOpcode implementation
func LoadOpcode(c *Context, i interface{}) error {

	name := util.GetString(i)
	if len(name) == 0 {
		return fmt.Errorf("invalid symbol name: %v", name)
	}
	v, found := c.Get(util.GetString(i))
	if !found {
		return fmt.Errorf("unknown symbol: %v", name)
	}

	c.Push(v)
	return nil
}

// CallOpcode bytecode implementation.
func CallOpcode(c *Context, i interface{}) error {

	var fname string
	var err error
	var v interface{}

	// Argument count is in operand
	argc := i.(int)

	// Function name is last item on stack
	v, err = c.Pop()
	if err != nil {
		return err
	}
	fname = util.GetString(v)

	// Arguments are in reverse order on stack.
	args := make([]interface{}, argc)
	for n := 0; n < argc; n = n + 1 {
		v, err = c.Pop()
		if err != nil {
			return err
		}
		args[(argc-n)-1] = v
	}

	// Is it in the dictionary?
	fn, found := functions.FunctionDictionary[fname]
	if found {
		if argc > fn.Max || argc < fn.Min {
			return errors.New("incorrect number of function arguments")
		}

		f := fn.F
		v, err = f.(func([]interface{}) (interface{}, error))(args)
	} else {

		// How about as a user-defined function? These are in the symbol
		// table with "()" as the suffix.
		f, found := c.symbols.Get(fname + "()")
		if !found {
			return fmt.Errorf("undefined function: %v", fname)
		}

		// Depends on the type here as to what we call...

		switch af := f.(type) {
		case *ByteCode:

			// Make a new symbol table for the fucntion to run with,
			// and a new execution context. Store the argument list in
			// the child table.
			sf := NewChildSymbolTable("Function "+fname, c.symbols)
			cx := NewContext(sf, af)
			sf.Set("_args", args)

			// Run the function. If it doesn't get an error, then
			// extract the stop stack item as the result
			err = cx.Run()
			if err == nil {
				v, err = cx.Pop()
			}

		default:
			v, err = f.(func([]interface{}) (interface{}, error))(args)
		}
	}

	if err != nil {
		return err
	}
	c.Push(v)
	return nil
}

// PushOpcode bytecode implementation
func PushOpcode(c *Context, i interface{}) error {
	return c.Push(i)
}

// DropOpcode implementation
func DropOpcode(c *Context, i interface{}) error {
	_, err := c.Pop()
	return err
}

// AddOpcode bytecode implementation
func AddOpcode(c *Context, i interface{}) error {

	if c.sp < 1 {
		return errors.New("stack underflow")
	}
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
			return errors.New("unsupported datatype")
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
			return errors.New("unsupported datatype")
		}
	}
}

// AndOpcode bytecode implementation
func AndOpcode(c *Context, i interface{}) error {

	if c.sp < 1 {
		return errors.New("stack underflow")
	}
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

	if c.sp < 1 {
		return errors.New("stack underflow")
	}
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

	if c.sp < 1 {
		return errors.New("stack underflow")
	}
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
			return errors.New("unsupported datatype")
		}
	}
}

// MulOpcode bytecode implementation
func MulOpcode(c *Context, i interface{}) error {

	if c.sp < 1 {
		return errors.New("stack underflow")
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
		return c.Push(v1.(int) * v2.(int))
	case float64:
		return c.Push(v1.(float64) * v2.(float64))
	case bool:
		return c.Push(v1.(bool) || v2.(bool))
	default:
		return errors.New("unsupported datatype")
	}
}

// DivOpcode bytecode implementation
func DivOpcode(c *Context, i interface{}) error {

	if c.sp < 1 {
		return errors.New("stack underflow")
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
			return errors.New("divide by zero")
		}
		return c.Push(v1.(int) / v2.(int))
	case float64:
		if v2.(float64) == 0 {
			return errors.New("divide by zero")
		}
		return c.Push(v1.(float64) / v2.(float64))
	default:
		return errors.New("unsupported datatype")
	}
}

// BranchFalseOpcode bytecode implementation
func BranchFalseOpcode(c *Context, i interface{}) error {

	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	if !util.GetBool(v) {
		c.pc = address
	}
	return nil
}

// BranchOpcode bytecode implementation
func BranchOpcode(c *Context, i interface{}) error {

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	c.pc = address
	return nil
}

// BranchTrueOpcode bytecode implementation
func BranchTrueOpcode(c *Context, i interface{}) error {

	// Get test value
	v, err := c.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i)
	if address < 0 || address > c.bc.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	if util.GetBool(v) {
		c.pc = address
	}
	return nil
}

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

	c.Push(r)
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

	c.Push(r)
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
		return errors.New("unsupported array operation")

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
			return errors.New("unsupported type for operation")

		}
	}
	c.Push(r)
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
		return errors.New("unsupported array operation")

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
			return errors.New("unsupported type for operation")

		}
	}
	c.Push(r)
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

	var r bool

	switch v1.(type) {

	case []interface{}:
		return errors.New("unsupported array operation")

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
			return errors.New("unsupported type for operation")

		}
	}
	c.Push(r)
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
		return errors.New("unsupported array operation")

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
			return errors.New("unsupported type for operation")

		}
	}
	c.Push(r)
	return nil
}

// LoadIndexOpcode implementation
func LoadIndexOpcode(c *Context, i interface{}) error {

	index, err := c.Pop()
	if err != nil {
		return err
	}

	array, err := c.Pop()
	if err != nil {
		return err
	}

	switch a := array.(type) {

	// Index into map is just member access
	case map[string]interface{}:
		subscript := util.GetString(index)
		v, f := a[subscript]
		if !f {
			return fmt.Errorf("member not found: %s", subscript)
		}
		c.Push(v)

	// Index into array is integer index (1-based)
	case []interface{}:
		subscript := util.GetInt(index)
		if subscript < 1 || subscript > len(a) {
			return fmt.Errorf("invalid array index: %v", subscript)
		}
		v := a[subscript-1]
		c.Push(v)

	default:
		return fmt.Errorf("invalid type for index operation")
	}

	return nil
}

// StoreIndexOpcode implementation
func StoreIndexOpcode(c *Context, i interface{}) error {

	index, err := c.Pop()
	if err != nil {
		return err
	}

	array, err := c.Pop()
	if err != nil {
		return err
	}

	v, err := c.Pop()
	if err != nil {
		return err
	}

	switch a := array.(type) {

	// Index into map is just member access
	case map[string]interface{}:
		subscript := util.GetString(index)
		a[subscript] = v
		c.Push(a)

	// Index into array is integer index (1-based)
	case []interface{}:
		subscript := util.GetInt(index)
		if subscript < 1 || subscript > len(a) {
			return fmt.Errorf("invalid array index: %v", subscript)
		}
		a[subscript-1] = v
		c.Push(a)

	default:
		return fmt.Errorf("invalid type for index operation")
	}

	return nil
}

// NegateOpcode implementation
func NegateOpcode(c *Context, i interface{}) error {

	v, err := c.Pop()
	if err != nil {
		return err
	}

	switch value := v.(type) {
	case bool:
		c.Push(!value)

	case int:
		c.Push(-value)
	case float64:
		c.Push(0.0 - value)

	case string:
		return errors.New("invalid data type for negation")
	}
	return nil
}
