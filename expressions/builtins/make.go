package builtins

import (
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

// Make implements the make() function. The first argument must be a model of the
// array type (using the Go native version), and the second argument is the size.
func Make(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	kind := args[0]
	size := data.Int(args[1])

	// if it's an Ego type, get the model for the type.
	if v, ok := kind.(*data.Type); ok {
		kind = data.InstanceOfType(v)
	}

	array := make([]interface{}, size)

	if v, ok := kind.([]interface{}); ok {
		if len(v) > 0 {
			kind = v[0]
		}
	}

	// If the model is a type we know about, let's go ahead and populate the array
	// with specific values.
	switch kind.(type) {
	case []int, int:
		for i := range array {
			array[i] = 0
		}

	case []bool, bool:
		for i := range array {
			array[i] = false
		}

	case []string, string:
		for i := range array {
			array[i] = ""
		}

	case []float64, float64:
		for i := range array {
			array[i] = 0.0
		}
	}

	return array, nil
}
