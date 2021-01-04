package tokenizer

import "github.com/tucats/gopackages/util"

//ReservedWords is the list of reserved words in the _Ego_ language
var ReservedWords []string = []string{
	"break",
	"const",
	"defer",
	"else",
	"for",
	"func",
	"if",
	"import",
	"nil",
	"package",
	"return",
	"int",
	"float",
	"string",
	"bool",
	"struct",
}

var ExtendedReservedWords = []string{
	"array",
	"call",
	"catch",
	"print",
	"try",
}

// IsReserved indicates if a name is a reserved word.
func IsReserved(name string, includeExtensions bool) bool {
	r := util.InList(name, ReservedWords...)
	if includeExtensions {
		r = r || util.InList(name, ExtendedReservedWords...)
	}
	return r
}
