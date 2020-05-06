// Package expressions is a simple expression evaluator. It supports
// a rudementary symbol table with scoping, and knows about four data
// types (string, integer, double, and boolean). It does type casting as
// need automatically.
package expressions

// ValueType is the type of an expression result.
type ValueType int

const (
	// StringType represents a string value type
	StringType = 1

	// IntegerType represents an int value type
	IntegerType = 2

	// DoubleType represents a float64 value type
	DoubleType = 3

	// BoolType represents a bool type
	BoolType = 4
)

// Expression is the type for an instance of the expresssion evaluator.
type Expression struct {
	Source   string
	Type     ValueType
	Value    interface{}
	Tokens   []string
	TokenPos []int
	TokenP   int
}

// New creates a new Expression object. The expression to evaluate is
// provided.
func New(expr string) *Expression {

	var e = Expression{
		Source: expr,
	}
	var ep = &e
	return ep
}
