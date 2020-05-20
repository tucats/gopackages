package symbols

import (
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
		switch actual := v.(type) {
		case func([]interface{}) (interface{}, error):
			if !includeBuiltins {
				continue
			}

		case map[string]interface{}:
			skip := false
			for _, k2 := range actual {
				switch k2.(type) {
				case func([]interface{}) (interface{}, error):
					skip = true
					break
				}
			}
			if skip && !includeBuiltins {
				continue
			}
		}
		b.WriteString("   ")
		b.WriteString(k)
		b.WriteString(" = ")
		b.WriteString(util.Format(v))
		b.WriteString("\n")
	}

	if s.Parent != nil {
		sp := s.Parent.Format(includeBuiltins)
		b.WriteString("\n")
		b.WriteString(sp)
	}
	return b.String()
}
