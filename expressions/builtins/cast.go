package builtins

import (
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// Compiler-generate casting; generally always array types. This is used to
// convert numeric arrays to a different kind of array, to convert a string
// to an array of integer (rune) values, etc.  It is called from within
// the Call bytecode when the function is really a type.
func Cast(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	// Target t is the last parameter
	t := data.TypeOf(args[len(args)-1])
	source := args[len(args)-1]

	if t.IsString() {

		// If the source is a []byte type, we can just fetch the bytes and do a direct convesion.
		// If the source is a []int type, we can convert each integer to a rune and add it to a
		// string builder. Otherwise, just format it as a string value.
		return data.FormatUnquoted(source), nil

	}

	switch source.(type) {

	case string:
		return data.Coerce(source, data.InstanceOfType(t)), nil

	default:
		v := data.Coerce(source, data.InstanceOfType(t))
		if v != nil {
			return v, nil
		}

		return nil, errors.ErrInvalidType.Context(data.TypeOf(source).String())
	}
}
