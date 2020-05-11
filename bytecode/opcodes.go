package bytecode

import (
	"errors"
	"strings"

	"github.com/tucats/gopackages/util"
)

// StopOpcode bytecode implementation
func StopOpcode(b *ByteCode, i *I) error {
	b.running = false
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
