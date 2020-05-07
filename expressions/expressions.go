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

// Expression is the type for an instance of the expresssion evaluator.
type Expression struct {
	Source   string
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
	ep.Parse()
	return ep
}

// GetInt takes a generic interface and returns the integer value, using
// type coercion if needed.
func GetInt(v interface{}) int {
	return Coerce(v, 1).(int)
}

// GetBool takes a generic interface and returns the boolean value, using
// type coercion if needed.
func GetBool(v interface{}) bool {
	return Coerce(v, true).(bool)
}

// GetString takes a generic interface and returns the string value, using
// type coercion if needed.
func GetString(v interface{}) string {
	return Coerce(v, "").(string)
}

// GetFloat takes a generic interface and returns the float64 value, using
// type coercion if needed.
func GetFloat(v interface{}) float64 {
	return Coerce(v, float64(0)).(float64)
}
