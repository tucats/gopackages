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
	"package",
	"print",
	"return",
	"try",
}

// IsReserved indicates if a name is a reserved word.
func IsReserved(name string) bool {
	return util.InList(name, ReservedWords...)

}
