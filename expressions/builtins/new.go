package builtins

import (
	"reflect"
	"strings"
	"sync"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// New implements the $new() function. This function creates a new
// "zero value" of any given type or object. If an integer type
// number or a string type name is given, the "zero value" for
// that type is returned. For an array, struct, or map, a recursive
// copy is done of the members to a new object which is returned.
func New(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	// Is the type an integer? If so it's a type kind from the native
	// reflection package.
	if typeValue, ok := args[0].(int); ok {
		switch reflect.Kind(typeValue) {
		case reflect.Uint8, reflect.Int8:
			return byte(0), nil

		case reflect.Int32:
			return int32(0), nil

		case reflect.Int, reflect.Int64:
			return 0, nil

		case reflect.String:
			return "", nil

		case reflect.Bool:
			return false, nil

		case reflect.Float32:
			return float32(0), nil

		case reflect.Float64:
			return float64(0), nil

		default:
			return nil, errors.ErrInvalidType.In("new").Context(typeValue)
		}
	}

	// Is it an actual type?
	if typeValue, ok := args[0].(*data.Type); ok {
		return typeValue.InstanceOf(typeValue), nil
	}

	// Is the type a string? If so it's a bult-in scalar type name
	if typeValue, ok := args[0].(string); ok {
		switch strings.ToLower(typeValue) {
		case data.BoolType.Name():
			return false, nil

		case data.ByteType.Name():
			return byte(0), nil

		case data.Int32TypeName:
			return int32(0), nil

		case data.IntTypeName:
			return 0, nil

		case data.Int64TypeName:
			return int64(0), nil

		case data.StringTypeName:
			return "", nil

		case data.Float32TypeName:
			return float32(0), nil

		case data.Float64TypeName:
			return float64(0), nil

		default:
			return nil, errors.ErrInvalidType.In("new").Context(typeValue)
		}
	}

	// If it's a WaitGroup, make a new one. Note, have to use the switch statement
	// form here to prevent Go from complaining that the interface{} is being copied.
	// In reality, we don't care as we don't actually make a copy anyway but instead
	// make a new waitgroup object.
	switch args[0].(type) {
	case sync.WaitGroup:
		return data.InstanceOfType(data.WaitGroupType), nil
	}

	r := DeepCopy(args[0], MaxDeepCopyDepth)

	// If there was a user-defined type in the source, make the clone point back to it
	switch v := r.(type) {
	case nil:
		return nil, errors.ErrInvalidValue.In("new").Context(nil)

	case symbols.SymbolTable:
		return nil, errors.ErrInvalidValue.In("new").Context("symbol table")

	case func(*symbols.SymbolTable, []interface{}) (interface{}, error):
		return v, nil

	// No action for this group
	case byte, int32, int, int64, string, float32, float64:

	default:
		return nil, errors.ErrInvalidType.In("new").Context(v)
	}

	return r, nil
}
