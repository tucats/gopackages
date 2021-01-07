package symbols

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tucats/gopackages/util"
)

// Format formats a symbol table into a string for printing/display
func (s *SymbolTable) Format(includeBuiltins bool) string {

	var b strings.Builder

	b.WriteString("Symbol table")
	if s.Name != "" {
		b.WriteString(" \"")
		b.WriteString(s.Name)
		b.WriteString("\"")
	}
	b.WriteString(":\n")

	// Iterate over the members to get a list of the keys

	keys := make([]string, 0)
	for k := range s.Symbols {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Now iterate over the keys in sorted order
	for _, k := range keys {
		v := s.Symbols[k]
		skip := false
		typeString := "package"

		switch actual := v.(type) {
		case func(*SymbolTable, []interface{}) (interface{}, error):
			if !includeBuiltins {
				continue
			}

		case map[string]interface{}:
			for kk, k2 := range actual {
				fmt.Printf("DEBUG: k = %s, kk = \"%s\"; k2 = %v\n", k, kk, k2)
				if kk == "__type" {
					typeString, _ = k2.(string)
				}
				if _, ok := k2.(func(*SymbolTable, []interface{}) (interface{}, error)); ok {
					skip = true
				}
			}
			if skip && !includeBuiltins {
				continue
			}
		}
		b.WriteString("   ")
		b.WriteString(k)
		b.WriteString(" = ")
		if skip {
			b.WriteString("(")
			b.WriteString(typeString)
			b.WriteString(") ")
		}

		// Any variable named _password or _token has it's value obscured
		if k == "_password" || k == "_token" {
			b.WriteString("\"******\"")
		} else {
			b.WriteString(util.Format(v))
		}
		b.WriteString("\n")
	}

	if s.Parent != nil {
		sp := s.Parent.Format(includeBuiltins)
		b.WriteString("\n")
		b.WriteString(sp)
	}
	return b.String()
}
