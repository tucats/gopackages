package util

import (
	"fmt"
	"sort"
	"strings"
)

// FormatUnquoted formats a value but does not
// put quotes on strings.
func FormatUnquoted(arg interface{}) string {
	switch v := arg.(type) {
	case string:
		return v
	default:
		return Format(v)
	}
}

// Format converts the given object into a string representation.
// In particular, this varies from a simple "%v" format in Go because
// it puts commas in the array list output to match the syntax of an
// array constant and puts quotes around string values.
func Format(arg interface{}) string {

	switch v := arg.(type) {

	case int:
		return fmt.Sprintf("%d", v)

	case bool:
		if v {
			return "true"
		}
		return "false"

	case float64:
		return fmt.Sprintf("%v", v)

	case map[string]interface{}:
		var b strings.Builder

		keys := make([]string, 0)
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteString("{")
		for n, k := range keys {
			i := v[k]
			if n > 0 {
				b.WriteString(",")
			}
			b.WriteRune(' ')
			b.WriteString(k)
			b.WriteString(": ")
			b.WriteString(Format(i))
		}
		b.WriteString(" }")
		return b.String()

	case []interface{}:

		var b strings.Builder
		b.WriteRune('[')

		for n, i := range v {
			if n > 0 {
				b.WriteString(", ")
			}
			b.WriteString(Format(i))
		}
		b.WriteRune(']')
		return b.String()

	case string:
		return "\"" + v + "\""

	case func([]interface{}) (interface{}, error):
		return fmt.Sprintf("<builtin %v>", arg)

	default:
		return "<bytecode>"
	}
}
