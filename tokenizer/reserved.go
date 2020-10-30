package tokenizer

//ReservedWords is the list of reserved words in the _Ego_ language
var ReservedWords []string = []string{
	"array",
	"break",
	"call",
	"catch",
	"const",
	"else",
	"for",
	"function",
	"if",
	"import",
	"package",
	"print",
	"return",
	"try",
}

// IsReserved indicates if a name is a reserved word.
func IsReserved(name string) bool {
	return InList(name, ReservedWords)

}
