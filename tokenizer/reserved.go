package tokenizer

import "github.com/tucats/gopackages/util"

//ReservedWords is the list of reserved words in the _Ego_ language
var ReservedWords []string = []string{
	"array",
	"break",
	"call",
	"catch",
	"const",
	"defer",
	"else",
	"for",
	"func",
	"if",
	"import",
	"nil",
	"package",
	"print",
	"return",
	"try",
	"int",
	"float",
	"string",
	"bool",
	"struct",
}

// IsReserved indicates if a name is a reserved word.
func IsReserved(name string) bool {
	return util.InList(name, ReservedWords...)

}
