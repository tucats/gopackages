package util

import (
	"fmt"
	"reflect"
	"runtime"
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
	case int64:
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

		// IF it's an internal function, show it's name. If it is a standard builtin from the
		// function library, show the short form of the name.
		if vv.Kind() == reflect.Func {
			if ui.IsActive(ui.DebugLogger) {
				name := runtime.FuncForPC(reflect.ValueOf(v).Pointer()).Name()
				name = strings.Replace(name, "github.com/tucats/gopackages/", "", 1)
				name = strings.Replace(name, "github.com/tucats/gopackages/runtime.", "", 1)
				return "builtin " + name
			} else {
				return "builtin"
			}
		}

		// If it's a bytecode.Bytecode pointer, use reflection to get the
		// Name field value and use that with the name. A function literal
		// will have no name.
		if vv.Kind() == reflect.Ptr {
			ts := vv.String()
			if ts == "<*bytecode.ByteCode Value>" {
				e := reflect.ValueOf(v).Elem()
				if ui.IsActive(ui.DebugLogger) {
					name := GetString(e.Field(0).Interface())
					return "func " + name
				} else {
					return "func"
				}
			}
			return fmt.Sprintf("ptr %s", ts)
		}

		if strings.HasPrefix(vv.String(), "<bytecode.StackMarker") {
			e := reflect.ValueOf(v).Field(0)
			name := GetString(e.Interface())
			return fmt.Sprintf("<%s>", name)
		}

		if strings.HasPrefix(vv.String(), "<bytecode.CallFrame") {
			e := reflect.ValueOf(v).Field(0)
			module := GetString(e.Interface())
			e = reflect.ValueOf(v).Field(1)
			line := GetInt(e.Interface())
			return fmt.Sprintf("<frame %s:%d>", module, line)
		}

		if ui.IsActive(ui.DebugLogger) {
			return fmt.Sprintf("kind %v %#v", vv.Kind(), v)
		}
		return fmt.Sprintf("kind %v %v", vv.Kind(), v)
	}
}
