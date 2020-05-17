package functions

import (
	"errors"

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

	switch arg := args[0].(type) {

	case map[string]interface{}:
		keys := make([]string, 0)
		for k := range arg {
			keys = append(keys, k)
		}
		return len(keys), nil

	case []interface{}:
		return len(arg), nil
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
