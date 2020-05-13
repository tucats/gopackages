package util

import (
	"fmt"
	"sort"
	"strings"
)

// Format converts the given object into a string representation.
// In particular, this varies from a simple "%v" format in Go because
// it puts commas in the array list output to match the syntax of an
// array constant and puts quotes around string values.
func Format(arg interface{}) string {

	switch v := arg.(type) {

	case map[string]interface{}:
		var b strings.Builder

		keys := make([]string, 0)
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		b.WriteRune('{')
		for n, k := range keys {
			i := v[k]
			if n > 0 {
				b.WriteString(", ")
			}
			b.WriteString(k)
			b.WriteRune(':')
			b.WriteString(Format(i))
		}
		b.WriteRune('}')
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

	default:
		return fmt.Sprintf("%v", v)
	}
}
