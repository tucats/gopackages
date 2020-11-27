package bytecode

import (
	"fmt"
	"strconv"
	"strings"
)

// Runtime error messages
const (
	ArgumentCountError            = "incorrect function argument count"
	ArgumentTypeError             = "incorrect function argument type"
	DivisionByZeroError           = "division by zero"
	InvalidArgCheckError          = "invalid ArgCheck array"
	InvalidArrayIndexError        = "invalid array index"
	InvalidBytecodeAddress        = "invalid bytecode address"
	InvalidFunctionCallError      = "invalid function call"
	InvalidIdentifierError        = "invalid identifier"
	InvalidSliceIndexError        = "invalid slice index"
	InvalidThisError              = "invalid _this_ identifier"
	InvalidTypeError              = "invalid or unsupported data type for this operation"
	NotATypeError                 = "not a type"
	OpcodeAlreadyDefinedError     = "opcode already defined: %d"
	ReadOnlyError                 = "invalid write to read-only item"
	StackUnderflowError           = "stack underflow"
	TryCatchMismatchError         = "try/catch stack error"
	UnimplementedInstructionError = "unimplemented bytecode instruction"
	UnknownIdentifierError        = "unknown identifier"
	UnknownMemberError            = "unknown structure member"
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
