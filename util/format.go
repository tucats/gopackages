package util

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
)

// LineColumnFormat describes the format string for the portion
// of formatted messages that include a line and column designation
const LineColumnFormat = "at %d:%d"

// LineFormat describes the format string for a message that contains
// just a line number.
const LineFormat = "at %d"

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

	if arg == nil {
		return "<nil>"
	}

	switch v := arg.(type) {

	case error:
		return fmt.Sprintf("%v", v)

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

		// Make a list of the keys, ignoring hidden members whose name
		// starts with "__"
		keys := make([]string, 0)
		for k := range v {
			if len(k) < 2 || k[0:2] != "__" {
				keys = append(keys, k)
			}
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

	default:
		vv := reflect.ValueOf(v)
		if vv.Kind() == reflect.Func {
			return "builtin"
		}
		if vv.Kind() == reflect.Ptr {
			ts := vv.String()
			if ts == "<*bytecode.ByteCode Value>" {
				return "func"
			}
			return fmt.Sprintf("ptr %s", ts)
		}
		if ui.DebugMode {
			return fmt.Sprintf("kind %v <%#v>", vv.Kind(), v)
		}
		return fmt.Sprintf("kind %v <%v>", vv.Kind(), v)
	}
}
