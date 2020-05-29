package functions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

// FunctionLower implements the lower() function
func FunctionLower(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return strings.ToLower(util.GetString(args[0])), nil
}

// FunctionUpper implements the upper() function
func FunctionUpper(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return strings.ToUpper(util.GetString(args[0])), nil
}

// FunctionLeft implements the left() function
func FunctionLeft(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	v := util.GetString(args[0])
	p := util.GetInt(args[1])

	if p <= 0 {
		return "", nil
	}
	if p >= len(v) {
		return v, nil
	}
	return v[:p], nil
}

// FunctionRight implements the right() function
func FunctionRight(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	v := util.GetString(args[0])
	p := util.GetInt(args[1])

	if p <= 0 {
		return "", nil
	}
	if p >= len(v) {
		return v, nil
	}
	return v[len(v)-p:], nil
}

// FunctionIndex implements the index() function
func FunctionIndex(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	switch arg := args[0].(type) {

	case []interface{}:
		for n, v := range arg {
			if reflect.DeepEqual(v, args[1]) {
				return n + 1, nil
			}
		}
		return 0, nil

	case map[string]interface{}:
		key := util.GetString(args[1])
		_, found := arg[key]
		return found, nil

	default:
		v := util.GetString(args[0])
		p := util.GetString(args[1])

		return strings.Index(v, p) + 1, nil
	}
}

// FunctionSubstring implements the substring() function
func FunctionSubstring(symbols *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	v := util.GetString(args[0])
	p1 := util.GetInt(args[1])
	p2 := util.GetInt(args[2])

	if p1 < 1 {
		p1 = 1
	}
	if p2 == 0 {
		return "", nil
	}
	if p2+p1 > len(v) {
		p2 = len(v) - p1 + 1
	}

	s := v[p1-1 : p1+p2-1]
	return s, nil
}

// FunctionFormat impelments the _strings.format() fucntion
func FunctionFormat(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	if len(args) == 0 {
		return "", nil
	}

	if len(args) == 1 {
		return util.GetString(args[0]), nil
	}

	return fmt.Sprintf(util.GetString(args[0]), args[1:]...), nil
}

// FunctionChars implements the _strings.chars() function. This accepts a string
// value and converts it to an array of characters.
func FunctionChars(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	v := util.GetString(args[0])
	r := make([]interface{}, 0)

	for n := 0; n < len(v); n = n + 1 {
		r = append(r, v[n:n+1])
	}
	return r, nil
}

// FunctionInts implements the _strings.ints() function. This accepts a string
// value and converts it to an array of integer rune values.
func FunctionInts(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	v := util.GetString(args[0])
	r := make([]interface{}, 0)
	i := []rune(v)

	for n := 0; n < len(i); n = n + 1 {
		r = append(r, int(i[n]))
	}
	return r, nil
}

// FunctionToString implements the _strings.string() function, which accepts an array
// of items and converts it to a single long string of each item. Normally , this is
// an array of characters.
func FunctionToString(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	var b strings.Builder

	for _, v := range args {

		switch a := v.(type) {
		case string:
			b.WriteString(a)

		case int:
			b.WriteRune(rune(a))

		case []interface{}:
			for _, c := range a {
				switch k := c.(type) {
				case int:
					b.WriteRune(rune(k))
				case string:
					b.WriteString(util.GetString(c))
				default:
					return nil, errors.New("incorrect argument type")
				}
			}
		default:
			return nil, errors.New("incorrect argument type")
		}
	}
	return b.String(), nil

}
