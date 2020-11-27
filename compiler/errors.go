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

// Compiler errors. Currently these are the string values. They will eventually
// be converted to identifiers for localized assets.
const (
	FunctionAlreadyExistsError     = "function already defined"
	InvalidConstantError           = "invalid constant expression"
	InvalidDirectiveError          = "invalid directive name"
	InvalidFunctionArgument        = "invalid function argument"
	InvalidFunctionCall            = "invalid function call"
	InvalidFunctionName            = "invalid function name"
	InvalidImportError             = "import not permitted inside a block or loop"
	InvalidLoopControlError        = "loop control statement outside of for-loop"
	InvalidLoopIndexError          = "invalid loop index variable"
	InvalidReturnValueError        = "invalid return value for void function"
	InvalidSymbolError             = "invalid symbol name"
	InvalidTypeNameError           = "invalid type name"
	MissingAssignmentError         = "missing '=' or ':='"
	MissingBracketError            = "missing array bracket"
	MissingBlockError              = "missing '{'"
	MissingCaseError               = "missing 'case'"
	MissingCatchError              = "missing 'catch' clause"
	MissingColonError              = "missing ':'"
	MissingEndOfBlockError         = "missing '}'"
	MissingEqualError              = "missing '='"
	MissingForLoopInitializerError = "missing for-loop initializer"
	MissingFunctionTypeError       = "missing function return type"
	MissingLoopAssignmentError     = "missing ':='"
	MissingSemicolonError          = "missing ';'"
	PackageRedefinitionError       = "cannot redefine existing package"
	TestingAssertError             = "testing @assert failure"
	UnexpectedTokenError           = "unexpected token"
)

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
