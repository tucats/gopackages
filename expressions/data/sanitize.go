package data

import (
	"strings"
)

// SanitizeName is used to examine a string that is used as a name (a filename,
// a module name, etc.). The function will ensure it has no embedded characters
// that would either reformat a string inappropriately -- such as entries in a
// log -- or allow any kind of unwanted injection.
//
// The function converts all control characters that could affect line ending or
// spacing to a "." character. It also processes other selection punctuation that
// is not allowed in an Ego name.
func SanitizeName(name string) string {
	result := strings.Builder{}
	blackList := []rune{'$', '\\', '/', '.', ';', ':'}

	for _, ch := range name {
		if ch < 26 {
			ch = '.'
		} else {
			for _, badCh := range blackList {
				if ch == badCh {
					ch = '.'

					break
				}
			}
		}

		result.WriteRune(ch)
	}

	return result.String()
}
