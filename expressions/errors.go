package expressions

import (
	"fmt"
	"strings"
)

// Error message codes
const (
	GeneralExpressionErr = iota
	BlockQuoteErr
)

// Error message strings
const (
	GeneralExpressionError  = "general expression error"
	InvalidListError        = "invalid list"
	InvalidRangeError       = "invalid array range"
	InvalidSymbolError      = "invalid symbol name"
	MismatchedQuoteError    = "mismatched quote error"
	MissingBracketError     = "missing or invalid '[]'"
	MissingColonError       = "missing ':'"
	MissingParenthesisError = "missing parenthesis"
	MissingTermError        = "missing term"
)

// ErrorMessageMap is used to map an error code to a message.
var ErrorMessageMap map[int]string = map[int]string{
	GeneralExpressionErr: "general expression error",
	BlockQuoteErr:        "mismatched quote character",
}

// Error contains an error generated from the compiler
type Error struct {
	code   int
	text   string
	line   int
	column int
	token  string
}

// NewErrorCode creates an error object using a numeric code
// rather than a message text
func (e *Expression) NewErrorCode(code int, parm string) *Error {
	text, found := ErrorMessageMap[code]
	if !found {
		text = ErrorMessageMap[GeneralExpressionErr]
	}
	if len(parm) > 0 {
		text = text + ": "
	}
	err := e.NewStringError(text, parm)
	err.code = code
	return err
}

// NewError generates a new error
func (e *Expression) NewError(msg string) *Error {
	err := e.NewStringError(msg, "")
	err.code = GeneralExpressionErr
	return err
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

// Code returns the numeric code for the error
func (e *Error) Code() int {
	return e.code
}
