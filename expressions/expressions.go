// Package expressions is a simple expression evaluator. It supports
// a rudementary symbol table with scoping, and knows about four data
// types (string, integer, double, and boolean). It does type casting as
// need automatically.
//
// The general pattern of use is:
//
//    e := expressions.New("expression string")
//    v, err := expressions.eval(symbolTableMap)
//    i := GetInt(v)
//    f := GetFloag(v)
//    s := GetString(v)
//    b := GetBool(v)
//
package expressions

import "github.com/tucats/gopackages/bytecode"

// Expression is the type for an instance of the expresssion evaluator.
type Expression struct {
	Source   string
	Tokens   []string
	TokenPos []int
	TokenP   int
	b        *bytecode.ByteCode
	err      error
}

// New creates a new Expression object. The expression to evaluate is
// provided.
func New(expr string) *Expression {

	// Create a new bytecode object, and then use it to create a new
	// expression object.
	b := bytecode.New(expr)
	e := NewWithBytecode(b)

	// tokenize
	e.Parse(expr)

	// compile
	e.err = e.conditional()

	return e

}

// NewWithBytecode allocates an expression object and
// attaches the provided bytecode structure.
func NewWithBytecode(b *bytecode.ByteCode) *Expression {
	var e = Expression{}
	var ep = &e
	ep.b = b
	return ep

}

// Error returns the last error seen on the expression object.
func (e *Expression) Error() error {
	return e.err
}

// Disasm calls the bytecode disassembler.
func (e *Expression) Disasm() {
	e.b.Disasm()
}
