package builtins

import (
	"reflect"
	"strings"

	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// Index implements the index() function.
func Index(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if !extensions() {
		return nil, errors.ErrExtension.Context("index")
	}

	switch arg := args[0].(type) {
	case []interface{}:
		for n, v := range arg {
			if reflect.DeepEqual(v, args[1]) {
				return n, nil
			}
		}

		return -1, nil

	default:
		v := data.String(args[0])
		p := data.String(args[1])

		return strings.Index(v, p) + 1, nil
	}
}
