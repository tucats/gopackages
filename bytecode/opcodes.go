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

// StoreOpcode implementation
func StoreOpcode(b *ByteCode, i *I) error {

	v, err := b.Pop()
	if err != nil {
		return err
	}

	b.Set(util.GetString(i.operand), v)
	return nil
}

// LoadOpcode implementation
func LoadOpcode(b *ByteCode, i *I) error {

	name := util.GetString(i.operand)
	if len(name) == 0 {
		return fmt.Errorf("invalid symbol name: %v", name)
	}
	v := b.Get(util.GetString(i.operand))
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

	// Argument count is in operand
	argc := i.operand.(int)

	// Function name is last item on stack
	v, err := b.Pop()
	if err != nil {
		return err
	}
	fname = util.GetString(v)

	// Arguments are in reverse order on stack.
	args := make([]interface{}, argc)
	for n := 0; n < argc; n = n + 1 {
		v, err := b.Pop()
		if err != nil {
			return err
		}
		args[(argc-n)-1] = v
	}

	fn, found := util.FunctionDictionary[fname]
	if !found {
		return errors.New("undefined function: " + fname)
	}
	if argc > fn.Max || argc < fn.Min {
		return errors.New("incorrect number of function arguments")
	}

	f := fn.F
	v, err = f.(func([]interface{}) (interface{}, error))(args)

	if err != nil {
		return err
	}
	b.Push(v)
	return nil
}

// PushOpcode bytecode implementation
func PushOpcode(b *ByteCode, i *I) error {
	return b.Push(i.operand)
}

// AddOpcode bytecode implementation
func AddOpcode(b *ByteCode, i *I) error {

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

// SubOpcode bytecode implementation
func SubOpcode(b *ByteCode, i *I) error {

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

// MulOpcode bytecode implementation
func MulOpcode(b *ByteCode, i *I) error {

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
	v1, err := b.Pop()
	if err != nil {
		return err
	}
	v2, err := b.Pop()
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
	address := util.GetInt(i.operand)
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
	address := util.GetInt(i.operand)
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
	address := util.GetInt(i.operand)
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
