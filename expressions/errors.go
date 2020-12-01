package expressions

import (
	"fmt"
	"strings"

	"github.com/tucats/gopackages/util"
)

// Error message strings
const (
	BlockQuoteError         = "invalid block quote terminator"
	GeneralExpressionError  = "general expression error"
	InvalidListError        = "invalid list"
	InvalidRangeError       = "invalid array range"
	InvalidSymbolError      = "invalid symbol name"
	MismatchedQuoteError    = "mismatched quote error"
	MissingBracketError     = "missing or invalid '[]'"
	MissingColonError       = "missing ':'"
	MissingParenthesisError = "missing parenthesis"
	MissingTermError        = "missing term"
	UnexpectedTokenError    = "unexpected token"
)

// Error contains an error generated from the compiler
type Error struct {
	code   int
	text   string
	line   int
	column int
	token  string
}

// NewError creates an error object
func (e *Expression) NewError(msg string, args ...interface{}) *Error {
	token := ""
	if len(args) > 0 {
		token = util.GetString(args[0])
	}

	p := e.t.TokenP
	if p < 0 {
		p = 0
	}
	if p >= len(e.t.Tokens) {
		p = len(e.t.Tokens) - 1
	}
	return &Error{
		text:   msg,
		line:   e.t.Line[p],
		column: e.t.Pos[p],
		token:  token,
	}
}

// Error produces an error string from this object.
func (e *Error) Error() string {

	var b strings.Builder

	b.WriteString("compile error ")
	if e.line > 0 {
		b.WriteString(fmt.Sprintf(util.LineColumnFormat, e.line, e.column))
	}
	b.WriteString(", ")
	b.WriteString(e.text)
	if len(e.token) > 0 {
		b.WriteString(": ")
		b.WriteString(e.token)
	}
	return b.String()
}
