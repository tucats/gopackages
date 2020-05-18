package symbols

import (
	"sort"
	"strings"

	"github.com/tucats/gopackages/util"
)

// Format formats a symbol table into a string for printing/display
func (s *SymbolTable) Format() string {

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

		b.WriteString("   ")
		b.WriteString(k)
		b.WriteString(" = ")
		b.WriteString(util.Format(v))
		b.WriteString("\n")
	}

	if s.Parent != nil {
		sp := s.Parent.Format()
		b.WriteString("\n")
		b.WriteString(sp)
	}
	return b.String()
}
