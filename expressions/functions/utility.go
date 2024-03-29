package functions

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/symbols"
	"github.com/tucats/gopackages/util"
)

// Sleep implements util.sleep()
func Sleep(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	duration, err := time.ParseDuration(util.GetString(args[0]))
	if err == nil {
		time.Sleep(duration)
	}
	return true, err
}

// ProfileGet implements the profile.get() function
func ProfileGet(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	key := util.GetString(args[0])
	return persistence.Get(key), nil

}

// ProfileSet implements the profile.set() function
func ProfileSet(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	key := util.GetString(args[0])

	// If the value is an empty string, delete the key else
	// store the value for the key.
	value := util.GetString(args[1])
	if value == "" {
		persistence.Delete(key)
	} else {
		persistence.Set(key, value)
	}

	return nil, persistence.Save()
}

// ProfileDelete implements the profile.delete() function
func ProfileDelete(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	key := util.GetString(args[0])
	persistence.Delete(key)
	return nil, nil
}

// ProfileKeys implements the profile.keys() function
func ProfileKeys(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	keys := persistence.Keys()
	result := make([]interface{}, len(keys))
	for i, key := range keys {
		result[i] = key
	}
	return result, nil
}

// UUID implements the uuid() function
func UUID(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	u := uuid.New()
	return u.String(), nil
}

// Length implements the len() function
func Length(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if args[0] == nil {
		return 0, nil
	}

	switch arg := args[0].(type) {

	case error:
		return len(arg.Error()), nil

	case map[string]interface{}:
		keys := make([]string, 0)
		for k := range arg {
			if !strings.HasPrefix(k, "__") {
				keys = append(keys, k)
			}
		}
		return len(keys), nil

	case []interface{}:
		return len(arg), nil

	case nil:
		return 0, nil

	default:
		v := util.Coerce(args[0], "")
		if v == nil {
			return 0, nil
		}
		return len(v.(string)), nil
	}
}

// Array implements the array() function, which creates
// an empty array of the given size. IF there are two parameters,
// the first must be an existing array which is resized to match
// the new array
func Array(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	var array []interface{}
	count := 0

	if len(args) == 2 {
		switch v := args[0].(type) {
		case []interface{}:
			count = util.GetInt(args[1])
			if count < len(v) {
				array = v[:count]
			} else if count == len(v) {
				array = v
			} else {
				array = append(v, make([]interface{}, count-len(v))...)
			}
		default:
			return nil, errors.ErrInvalidType
		}
	} else {
		count = util.GetInt(args[0])
		array = make([]interface{}, count)
	}
	return array, nil

}

// GetEnv implements the util.getenv() function which reads
// an environment variable from the os.
func GetEnv(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return os.Getenv(util.GetString(args[0])), nil
}

// GetMode implements the util.Mode() function which reports the runtime mode
func GetMode(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	m, ok := symbols.Get("__exec_mode")
	if !ok {
		m = "run"
	}
	return m, nil
}

// Members gets an array of the names of the fields in a structure
func Members(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	switch v := args[0].(type) {
	case map[string]interface{}:

		keys := make([]string, 0)
		for k := range v {
			if !strings.HasPrefix(k, "__") {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		a := make([]interface{}, len(keys))
		for n, k := range keys {
			a[n] = k
		}
		return a, nil

	default:
		return nil, errors.ErrInvalidType
	}
}

// Sort implements the sort() function.
func Sort(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	// Make a master array of the values presented
	var array []interface{}
	for _, a := range args {
		switch v := a.(type) {
		case []interface{}:
			array = append(array, v...)
		default:
			array = append(array, v)
		}
	}

	if len(array) == 0 {
		return array, nil
	}

	v1 := array[0]
	switch v1.(type) {

	case int:
		intArray := make([]int, 0)
		for _, i := range array {
			intArray = append(intArray, util.GetInt(i))
		}
		sort.Ints(intArray)
		resultArray := make([]interface{}, len(array))
		for n, i := range intArray {
			resultArray[n] = i
		}
		return resultArray, nil

	case float64:
		floatArray := make([]float64, 0)
		for _, i := range array {
			floatArray = append(floatArray, util.GetFloat(i))
		}
		sort.Float64s(floatArray)
		resultArray := make([]interface{}, len(array))
		for n, i := range floatArray {
			resultArray[n] = i
		}
		return resultArray, nil

	case string:
		stringArray := make([]string, 0)
		for _, i := range array {
			stringArray = append(stringArray, util.GetString(i))
		}
		sort.Strings(stringArray)
		resultArray := make([]interface{}, len(array))
		for n, i := range stringArray {
			resultArray[n] = i
		}
		return resultArray, nil

	default:
		return nil, errors.ErrInvalidType
	}
}

// Exit implements the util.exit() function
func Exit(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	// If no arguments, just do a simple exit
	if len(args) == 0 {
		os.Exit(0)
	}

	switch v := args[0].(type) {

	case int:
		os.Exit(v)

	case string:
		return nil, errors.NewMessage(v)

	default:
		return nil, errors.ErrInvalidType
	}

	return nil, nil
}

// FormatSymbols implements the util.symbols() function
func FormatSymbols(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return syms.Format(false), nil
}

// Type implements the type() function
func Type(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	switch v := args[0].(type) {
	case nil:
		return "nil", nil
	case error:
		return "error", nil
	case int:
		return "int", nil
	case float64, float32:
		return "float", nil
	case string:
		return "string", nil
	case bool:
		return "bool", nil
	case []interface{}:
		return "array", nil

	default:
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Func {
			return "builtin", nil
		}
		if vv.Kind() == reflect.Ptr {
			ts := vv.String()
			if ts == "<*bytecode.ByteCode Value>" {
				return "func", nil
			}
			return fmt.Sprintf("ptr %s", ts), nil
		}
		return "unknown", nil
	}
}

// Append implements the builtin append() function, which concatenates all the items
// together as an array. The first argument is flattened into the result, and then each
// additional argument is added to the array as-is.
func Append(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	result := []interface{}{}
	for i, j := range args {
		if array, ok := j.([]interface{}); ok && i == 0 {
			result = append(result, array...)
		} else {
			result = append(result, j)
		}
	}
	return result, nil
}

// Delete can be used three ways. To delete a member from a structure, to delete
// an element from an array by index number, or to delete a symbol entirely. The
// first form requires a string name, the second form requires an integer index,
// and the third form does not have a second parameter.
func Delete(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if _, ok := args[0].(string); ok && len(args) != 1 {
		return nil, errors.ErrArgumentCount
	} else {
		if len(args) != 2 {
			return nil, errors.ErrArgumentCount
		}
	}
	switch v := args[0].(type) {

	case string:
		return nil, s.Delete(v, true)

	case map[string]interface{}:
		key := util.GetString(args[1])
		delete(v, key)
		return v, nil

	case []interface{}:
		i := util.GetInt(args[1])
		if i < 0 || i >= len(v) {
			return nil, errors.ErrArrayIndex
		}
		r := append(v[:i], v[i+1:]...)
		return r, nil

	default:
		return nil, errors.ErrArgumentCount
	}
}

// GetArgs implements util.Args() which fetches command-line arguments from
// the Ego command invocation, if any.
func GetArgs(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	r, found := s.Get("__cli_args")
	if !found {
		r = []interface{}{}
	}
	return r, nil
}

// Make implements the make() function. The first argument must be a model of the
// array type (using the Go native version), and the second argument is the size.
func Make(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	kind := args[0]
	size := util.GetInt(args[1])
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
	case map[string]interface{}:
		for i := range array {
			array[i] = map[string]interface{}{}
		}

	default:
		fmt.Printf("DEBUG: v = %#v\n", kind)
	}
	return array, nil

}
