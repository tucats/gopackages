package compiler

import (
	"fmt"
	"strings"
)

// Error contains an error generated from the compiler
type Error struct {
	text   string
	pkg    string
	line   int
	column int
	token  string
}

// NewError generates a new error
func (c *Compiler) NewError(msg string) *Error {

	p := c.t.TokenP
	if p < 0 {
		p = 0
	}
	if p >= len(c.t.Tokens) {
		p = len(c.t.Tokens) - 1
	}
	return &Error{
		text:   msg,
		line:   c.t.Line[p],
		column: c.t.Pos[p],
		token:  "",
		pkg:    c.PackageName,
	}
}

// NewTokenError generates a new error that includes the
// current token as part of the error information.
func (c *Compiler) NewTokenError(msg string) *Error {

	p := c.t.TokenP
	if p < 0 {
		p = 0
	}
	if p >= len(c.t.Tokens) {
		p = len(c.t.Tokens) - 1
	}
	return &Error{
		text:   msg,
		line:   c.t.Line[p],
		column: c.t.Pos[p],
		token:  c.t.Tokens[p],
		pkg:    c.PackageName,
	}
}

// NewStringError generates a new error that includes the
// string argument as part of the error information. The
// token string value is not used.
func (c *Compiler) NewStringError(msg string, s string) *Error {

	p := c.t.TokenP
	if p < 0 {
		p = 0
	}
	if p >= len(c.t.Tokens) {
		p = len(c.t.Tokens) - 1
	}
	return &Error{
		text:   fmt.Sprintf("%s: %v", msg, s),
		line:   c.t.Line[p],
		column: c.t.Pos[p],
		pkg:    c.PackageName,
	}
}

// Error produces an error string from this object.
func (e Error) Error() string {

	var b strings.Builder

	b.WriteString("compile error, ")
	if e.pkg != "" {
		b.WriteString("package ")
		b.WriteString(e.pkg)
		b.WriteString(", ")
	}
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
