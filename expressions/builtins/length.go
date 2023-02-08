package builtins

import (
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// Length implements the len() function.
func Length(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if args[0] == nil {
		return 0, nil
	}

	switch arg := args[0].(type) {
	case error:
		return len(arg.Error()), nil

	case nil:
		return 0, nil

	case string:
		return len(arg), nil

	default:
		// Extensions have to be enabled and we must not be in strict
		// type checking mode to return length of the stringified argument.
		if v, found := s.Get(defs.ExtensionsVariable); found {
			if data.Bool(v) {
				if v, found := s.Get(defs.TypeCheckingVariable); found {
					if data.Int(v) > 0 {
						return len(data.String(arg)), nil
					}
				}
			}
		}

		// Otherwise, invalid type.
		return 0, errors.ErrArgumentType.In("len").Context(data.TypeOf(args[0]))
	}
}

// SizeOf returns the size in bytes of an arbibrary object.
func SizeOf(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	size := data.SizeOf(args[0])

	return size, nil
}
