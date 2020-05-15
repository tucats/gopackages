package bytecode

import (
	"fmt"
	"strconv"
	"strings"
)

// Error contains an error generated from the execution context
type Error struct {
	text   string
	module string
	line   int
	token  string
}

// NewError generates a new error
func (c *Context) NewError(msg string) *Error {

	return &Error{
		text:   msg,
		module: c.Name,
		line:   c.line,
		token:  "",
	}
}

// NewStringError generates a new error that includes the
// current token as part of the error information.
func (c *Context) NewStringError(msg string, parm string) *Error {

	return &Error{
		text:   msg,
		module: c.Name,
		line:   c.line,
		token:  parm,
	}
}

// NewIntError generates a new error that includes the
// current token as part of the error information.
func (c *Context) NewIntError(msg string, parm int) *Error {

	return &Error{
		text:   msg,
		module: c.Name,
		line:   c.line,
		token:  strconv.Itoa(parm),
	}
}

// Error produces an error string from this object.
func (e Error) Error() string {

	var b strings.Builder

	b.WriteString("execution error, ")

	if len(e.module) > 0 {
		b.WriteString("in ")
		b.WriteString(e.module)
		b.WriteString(", ")
	}
	if e.line > 0 {
		b.WriteString(fmt.Sprintf("at line %d, ", e.line))
	}
	b.WriteString(e.text)
	if len(e.token) > 0 {
		b.WriteString(": ")
		b.WriteString(e.token)
	}
	return b.String()
}
