package functions

import (
	"strings"

	"github.com/tucats/gopackages/util"
)

// Runtime error messages
const (
	ArgumentCountError            = "incorrect function argument count"
	ArgumentTypeError             = "incorrect function argument type"
	DivisionByZeroError           = "division by zero"
	InvalidArgCheckError          = "invalid ArgCheck array"
	InvalidArrayIndexError        = "invalid array index"
	InvalidBytecodeAddress        = "invalid bytecode address"
	InvalidFileIdentifierError    = "invalid file identifier"
	InvalidFunctionCallError      = "invalid function call"
	InvalidIdentifierError        = "invalid identifier"
	InvalidNewValueError          = "invalid argument to new()"
	InvalidSliceIndexError        = "invalid slice index"
	InvalidTemplateNameError      = "invalid template reference"
	InvalidThisError              = "invalid _this_ identifier"
	InvalidTypeError              = "invalid or unsupported data type for this operation"
	InvalidValueError             = "invalid value for this operation"
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
	token  string
}

// NewError generates a new error
func NewError(fn, msg string, args ...interface{}) *Error {

	token := ""
	if len(args) > 0 {
		token = util.GetString(args[0])
	}
	return &Error{
		text:   msg,
		token:  token,
		module: fn,
	}
}

// Error produces an error string from this object.
func (e Error) Error() string {

	var b strings.Builder

	b.WriteString("function error, in ")
	b.WriteString(e.module)
	b.WriteString("(), ")
	b.WriteString(e.text)
	if len(e.token) > 0 {
		b.WriteString(": ")
		b.WriteString(e.token)
	}
	return b.String()
}
