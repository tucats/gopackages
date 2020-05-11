package util

import (
	"fmt"
	"strconv"
)

// GetInt takes a generic interface and returns the integer value, using
// type coercion if needed.
func GetInt(v interface{}) int {
	return Coerce(v, 1).(int)
}

// GetBool takes a generic interface and returns the boolean value, using
// type coercion if needed.
func GetBool(v interface{}) bool {
	return Coerce(v, true).(bool)
}

// GetString takes a generic interface and returns the string value, using
// type coercion if needed.
func GetString(v interface{}) string {
	return Coerce(v, "").(string)
}

// GetFloat takes a generic interface and returns the float64 value, using
// type coercion if needed.
func GetFloat(v interface{}) float64 {
	return Coerce(v, float64(0)).(float64)
}

// Coerce returns the value after it has been converted to the type of the
// model value.
func Coerce(v interface{}, model interface{}) interface{} {

	switch model.(type) {

	case int:
		switch value := v.(type) {
		case bool:
			if value {
				return 1
			}
			return 0

		case int:
			return value

		case float64:
			return int(value)

		case string:
			st, err := strconv.Atoi(value)
			if err != nil {
				return nil
			}
			return st
		}

	case float64:
		switch value := v.(type) {
		case bool:
			if value {
				return 1.0
			}
			return 0.0

		case int:
			return float64(value)

		case float64:
			return value

		case string:
			st, _ := strconv.ParseFloat(value, 64)
			return st
		}

	case string:
		switch value := v.(type) {
		case bool:
			if value {
				return "true"
			}
			return "false"

		case int:
			return strconv.Itoa(value)

		case float64:
			return fmt.Sprintf("%v", value)

		case string:
			return value
		}

	case bool:

		switch v.(type) {
		case bool:
			return v

		case int:
			return (v.(int) != 0)

		case float64:
			return v.(float64) != 0.0

		case string:
			switch v.(string) {
			case "true":
				return true
			case "false":
				return false
			default:
				return false
			}
		}
	}

	return nil
}

// Normalize accepts two different values and promotes them to
// the most compatable format
func Normalize(v1 interface{}, v2 interface{}) (interface{}, interface{}) {

	// Same type? we're done here

	switch v1.(type) {

	case string:
		switch v2.(type) {
		case string:
			return v1, v2
		case int:
			return v1, strconv.Itoa(v2.(int))
		case float64:
			return v1, fmt.Sprintf("%v", v2.(float64))
		case bool:
			if v2.(bool) {
				return v1, "true"
			}
			return v1, "false"
		}

	case float64:
		switch v2.(type) {
		case string:
			return fmt.Sprintf("%v", v1.(float64)), v2
		case int:
			return v1, float64(v2.(int))
		case float64:
			return v1, v2
		case bool:
			if v2.(bool) {
				return v1, 1.0
			}
			return v1, 0.0
		}

	case int:
		switch v2.(type) {
		case string:
			return strconv.Itoa(v1.(int)), v2
		case int:
			return v1, v2
		case float64:
			return float64(v1.(int)), v2
		case bool:
			if v2.(bool) {
				return v1, 1
			}
			return v1, 0
		}

	case bool:
		switch v2.(type) {
		case string:
			if v1.(bool) {
				return "true", v2.(string)
			}
			return "false", v2.(string)

		case int:
			if v1.(bool) {
				return 1, v2.(int)
			}
			return 0, v2.(int)

		case float64:
			if v1.(bool) {
				return 1.0, v2.(float64)
			}
			return 0.0, v2.(float64)

		case bool:
			return v1, v2
		}
	}
	return v1, v2
}
