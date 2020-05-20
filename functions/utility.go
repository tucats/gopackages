package functions

import (
	"errors"
	"os"
	"sort"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/util"
)

// FunctionProfile implements the profile() function
func FunctionProfile(args []interface{}) (interface{}, error) {

	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("incorrect number of function arguments")
	}

	key := util.GetString(args[0])

	if len(args) == 1 {
		return persistence.Get(key), nil
	}

	// If the value is an empty string, delete the key else
	// store the value for the key.
	value := util.GetString(args[1])
	if value == "" {
		persistence.Delete(key)
	} else {
		persistence.Set(key, value)
	}
	return true, nil
}

// FunctionUUID implements the uuid() function
func FunctionUUID(args []interface{}) (interface{}, error) {
	if len(args) != 0 {
		return nil, errors.New("incorrect number of function arguments")
	}
	u := uuid.New()
	return u.String(), nil
}

// FunctionLen implements the len() function
func FunctionLen(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of function arguments")
	}

	if args[0] == nil {
		return 0, nil
	}

	switch arg := args[0].(type) {

	case map[string]interface{}:
		keys := make([]string, 0)
		for k := range arg {
			if k != "__readonly" {
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
		return len(v.(string)), nil
	}
}

// FunctionArray implements the array() function, which creates
// an empty array of the given size. IF there are two parameters,
// the first must be an existing array which is resized to match
// the new array
func FunctionArray(args []interface{}) (interface{}, error) {

	if len(args) < 1 || len(args) > 2 {
		return nil, errors.New("incorrect number of function arguments")
	}

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
			return nil, errors.New("first argument must be array")
		}
	} else {
		count = util.GetInt(args[0])
		array = make([]interface{}, count)
	}
	return array, nil

}

// FunctionGetEnv implementes the util.getenv() function which reads
// an environment variable from the os.
func FunctionGetEnv(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of function arguments")
	}

	return os.Getenv(util.GetString(args[0])), nil
}

// FunctionMembers gets an array of the names of the fields in a structure
func FunctionMembers(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, errors.New("incorrect number of function arguments")
	}
	switch v := args[0].(type) {

	case map[string]interface{}:

		keys := make([]string, 0)
		for k := range v {
			if k != "__readonly" {
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
		return nil, errors.New("incorrect data type")
	}
}

// FunctionSort implements the sort() function.
func FunctionSort(args []interface{}) (interface{}, error) {

	switch array := args[0].(type) {

	case []interface{}:
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
			return nil, errors.New("unsupported data type")
		}
	default:
		return nil, errors.New("not an array type")
	}
}

// FunctionExit implements the _util.exit() function
func FunctionExit(args []interface{}) (interface{}, error) {

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
		return nil, errors.New("unsupported exit() type")
	}

	return nil, nil
}
