package expressions

import (
	"fmt"
	"strings"
)

// Error contains an error generated from the compiler
type Error struct {
	text   string
	line   int
	column int
	token  string
}

// NewError generates a new error
func (e *Expression) NewError(msg string) *Error {
	return e.NewStringError(msg, "")
}

// NewStringError generates a new error with a string parameter
func (e *Expression) NewStringError(msg string, parm string) *Error {

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
		token:  parm,
	}
}

// NewTokenError generates a new error that includes the
// current token as part of the error information.
func (e *Expression) NewTokenError(msg string) *Error {

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
		token:  e.t.Tokens[p],
	}
}

// Error produces an error string from this object.
func (e *Error) Error() string {

	var b strings.Builder

	b.WriteString("compile error, ")
	if e.line > 0 {
		b.WriteString(fmt.Sprintf("at line %d, column %d, ", e.line, e.column))
	}
	b.WriteString(e.text)
	if len(e.token) > 0 {
		b.WriteString(": ")
		b.WriteString(e.token)
	}
	return b.String()
}
