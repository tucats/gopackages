package bytecode

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/util"
)

// StopOpcode bytecode implementation
func StopOpcode(b *ByteCode, i *I) error {
	b.running = false
	return nil
}

// ArrayOpcode implementation
func ArrayOpcode(b *ByteCode, i *I) error {

	count := util.GetInt(i.Operand)
	array := make([]interface{}, count)

	for n := 0; n < count; n++ {
		v, err := b.Pop()
		if err != nil {
			return err
		}
		array[(count-n)-1] = v
	}

	b.Push(array)
	return nil
}

// StoreOpcode implementation
func StoreOpcode(b *ByteCode, i *I) error {

	v, err := b.Pop()
	if err != nil {
		return err
	}

	b.Set(util.GetString(i.Operand), v)
	return nil
}

// LoadOpcode implementation
func LoadOpcode(b *ByteCode, i *I) error {

	name := util.GetString(i.Operand)
	if len(name) == 0 {
		return fmt.Errorf("invalid symbol name: %v", name)
	}
	v := b.Get(util.GetString(i.Operand))
	if v == nil {
		return fmt.Errorf("unknown symbol: %v", name)
	}

	b.Push(v)
	return nil
}

// CallOpcode bytecode implementation.
func CallOpcode(b *ByteCode, i *I) error {

	var fname string
	var err error
	var v interface{}

	// Argument count is in operand
	argc := i.Operand.(int)

	// Function name is last item on stack
	v, err = b.Pop()
	if err != nil {
		return err
	}
	fname = util.GetString(v)

	// Arguments are in reverse order on stack.
	args := make([]interface{}, argc)
	for n := 0; n < argc; n = n + 1 {
		v, err = b.Pop()
		if err != nil {
			return err
		}
		args[(argc-n)-1] = v
	}

	// Is it in the dictionary?
	fn, found := util.FunctionDictionary[fname]
	if found {
		if argc > fn.Max || argc < fn.Min {
			return errors.New("incorrect number of function arguments")
		}

		f := fn.F
		v, err = f.(func([]interface{}) (interface{}, error))(args)
	} else {

		// How about as a user-defined function? These are in the symbol
		// table with "()" as the suffix.
		f, found := b.symbols[fname+"()"]
		if !found {
			return fmt.Errorf("undefined function: %v", fname)
		}
		v, err = f.(func([]interface{}) (interface{}, error))(args)
	}

	if err != nil {
		return err
	}
	b.Push(v)
	return nil
}

// PushOpcode bytecode implementation
func PushOpcode(b *ByteCode, i *I) error {
	return b.Push(i.Operand)
}

// AddOpcode bytecode implementation
func AddOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
			return b.Push(newArray)

		// Everything else is a simple append.
		default:
			newArray := append(vx, v2)
			return b.Push(newArray)
		}
		// All other types are scalar math
	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			return b.Push(v1.(int) + v2.(int))
		case float64:
			return b.Push(v1.(float64) + v2.(float64))
		case string:
			return b.Push(v1.(string) + v2.(string))
		case bool:
			return b.Push(v1.(bool) && v2.(bool))
		default:
			return errors.New("unsupported datatype")
		}
	}
}

// AndOpcode bytecode implementation
func AndOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v1, err := b.Pop()
	if err != nil {
		return err
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}

	return b.Push(util.GetBool(v1) && util.GetBool(v2))

}

// OrOpcode bytecode implementation
func OrOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v1, err := b.Pop()
	if err != nil {
		return err
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}

	return b.Push(util.GetBool(v1) || util.GetBool(v2))

}

// SubOpcode bytecode implementation
func SubOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
		return b.Push(newArray)

	// Everything else is a scalar subtraction
	default:
		v1, v2 = util.Normalize(v1, v2)
		switch v1.(type) {
		case int:
			return b.Push(v1.(int) - v2.(int))
		case float64:
			return b.Push(v1.(float64) - v2.(float64))
		case string:
			s := strings.ReplaceAll(v1.(string), v2.(string), "")
			return b.Push(s)
		default:
			return errors.New("unsupported datatype")
		}
	}
}

// MulOpcode bytecode implementation
func MulOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
	if err != nil {
		return err
	}

	v1, v2 = util.Normalize(v1, v2)
	switch v1.(type) {
	case int:
		return b.Push(v1.(int) * v2.(int))
	case float64:
		return b.Push(v1.(float64) * v2.(float64))
	case bool:
		return b.Push(v1.(bool) || v2.(bool))
	default:
		return errors.New("unsupported datatype")
	}
}

// DivOpcode bytecode implementation
func DivOpcode(b *ByteCode, i *I) error {

	if b.sp < 1 {
		return errors.New("stack underflow")
	}
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
	if err != nil {
		return err
	}

	v1, v2 = util.Normalize(v1, v2)
	switch v1.(type) {
	case int:
		if v2.(int) == 0 {
			return errors.New("divide by zero")
		}
		return b.Push(v1.(int) / v2.(int))
	case float64:
		if v2.(float64) == 0 {
			return errors.New("divide by zero")
		}
		return b.Push(v1.(float64) / v2.(float64))
	default:
		return errors.New("unsupported datatype")
	}
}

// BranchFalseOpcode bytecode implementation
func BranchFalseOpcode(b *ByteCode, i *I) error {

	// Get test value
	v, err := b.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i.Operand)
	if address < 0 || address > b.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	if !util.GetBool(v) {
		b.pc = address
	}
	return nil
}

// BranchOpcode bytecode implementation
func BranchOpcode(b *ByteCode, i *I) error {

	// Get destination
	address := util.GetInt(i.Operand)
	if address < 0 || address > b.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	b.pc = address
	return nil
}

// BranchTrueOpcode bytecode implementation
func BranchTrueOpcode(b *ByteCode, i *I) error {

	// Get test value
	v, err := b.Pop()
	if err != nil {
		return err
	}

	// Get destination
	address := util.GetInt(i.Operand)
	if address < 0 || address > b.emitPos {
		return errors.New("invalid destination address: " + strconv.Itoa(address))
	}

	if util.GetBool(v) {
		b.pc = address
	}
	return nil
}

// EqualOpcode implementation
func EqualOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}

	v1, err := b.Pop()
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

	b.Push(r)
	return nil

}

// NotEqualOpcode implementation
func NotEqualOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}

	v1, err := b.Pop()
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

	b.Push(r)
	return nil

}

// GreaterThanOpcode implementation
func GreaterThanOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
	b.Push(r)
	return nil
}

// GreaterThanOrEqualOpcode implementation
func GreaterThanOrEqualOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
	b.Push(r)
	return nil
}

// LessThanOpcode implementation
func LessThanOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
	b.Push(r)
	return nil
}

// LessThanOrEqualOpcode implementation
func LessThanOrEqualOpcode(b *ByteCode, i *I) error {

	// Terms pushed in reverse order
	v2, err := b.Pop()
	if err != nil {
		return err
	}
	v1, err := b.Pop()
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
	b.Push(r)
	return nil
}

// IndexOpcode implementation
func IndexOpcode(b *ByteCode, i *I) error {

	index, err := b.Pop()
	if err != nil {
		return err
	}

	array, err := b.Pop()
	if err != nil {
		return err
	}

	subscript := util.GetInt(index)
	switch a := array.(type) {
	case []interface{}:
		if subscript < 1 || subscript > len(a) {
			return fmt.Errorf("invalid array index: %v", subscript)
		}
		v := a[subscript-1]
		b.Push(v)

	default:
		return fmt.Errorf("invalid type for index operation")
	}

	return nil
}

// NegateOpcode implementation
func NegateOpcode(b *ByteCode, i *I) error {

	v, err := b.Pop()
	if err != nil {
		return err
	}

	switch value := v.(type) {
	case bool:
		b.Push(!value)

	case int:
		b.Push(-value)
	case float64:
		b.Push(0.0 - value)

	case string:
		return errors.New("invalid data type for negation")
	}
	return nil
}
