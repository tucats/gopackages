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
}

// New creates a new Expression object. The expression to evaluate is
// provided.
func New(expr string) *Expression {
	var e = Expression{
		Source: expr,
	}
	var ep = &e
	ep.Parse()
	return ep
}
