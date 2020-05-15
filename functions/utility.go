package functions

import (
	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/util"
)

// FunctionProfile implements the profile() function
func FunctionProfile(args []interface{}) (interface{}, error) {

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
	u := uuid.New()
	return u.String(), nil
}

// FunctionLen implements the len() function
func FunctionLen(args []interface{}) (interface{}, error) {

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
