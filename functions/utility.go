package functions

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/symbols"
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
	return nil, nil
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
			if k != "__readonly" && k != "__type" {
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
			return nil, NewError("array", InvalidTypeError)
		}
	} else {
		count = util.GetInt(args[0])
		array = make([]interface{}, count)
	}
	return array, nil

}

// GetEnv implementes the util.getenv() function which reads
// an environment variable from the os.
func GetEnv(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return os.Getenv(util.GetString(args[0])), nil
}

// Members gets an array of the names of the fields in a structure
func Members(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	switch v := args[0].(type) {
	case map[string]interface{}:

		keys := make([]string, 0)
		for k := range v {
			if k != "__readonly" && k != "__type" {
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
		return nil, NewError("members", InvalidTypeError)
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
		return nil, NewError("sort", InvalidTypeError)
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
		return nil, errors.New(v)

	default:
		return nil, NewError("exit", InvalidTypeError)
	}

	return nil, nil
}

// FormatSymbols implements the util.symbols() function
func FormatSymbols(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return syms.Format(false), nil
}

// Type implements the util.type() function
func Type(syms *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	switch v := args[0].(type) {
	case error:
		return "error", nil
	case int:
		return "int", nil
	case float64:
		return "float", nil
	case string:
		return "string", nil
	case bool:
		return "bool", nil
	case []interface{}:
		return "array", nil
	case map[string]interface{}:
		if sv, ok := v["__type"]; ok {
			return util.GetString(sv), nil
		}
		return "struct", nil

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

// Signal creates an error object based on the
// parameters
func Signal(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return NewError("error", util.GetString(args[0]), args[1:]...), nil
}
