package expressions

import (
	"errors"
	"fmt"
	"strconv"
)

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) Eval(symbols map[string]interface{}) (interface{}, error) {

	e.TokenP = 0

	v1, err := e.eval1(symbols)
	if err != nil {
		return nil, err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if op == "+" || op == "-" {
			e.TokenP = e.TokenP + 1

			v2, err := e.eval1(symbols)
			if err != nil {
				return nil, err
			}

			v1, v2 = Normalize(v1, v2)
			switch op {

			case "+":
				switch v1.(type) {
				case int:
					v1 = v1.(int) + v2.(int)
				case string:
					v1 = v1.(string) + v2.(string)
				case float64:
					v1 = v1.(float64) + v2.(float64)
				case bool:
					v1 = v1.(bool) && v2.(bool)
				}

			case "-":
				switch v1.(type) {
				case int:
					v1 = v1.(int) - v2.(int)
				case float64:
					v1 = v1.(float64) - v2.(float64)
				default:
					return nil, errors.New("invlid type for '-' operator")
				}

			}

		} else {
			parsing = false
		}
	}
	return v1, nil
}

// Normalize accepts two different values and promotes them to
// the most compatable format
func Normalize(v1 interface{}, v2 interface{}) (interface{}, interface{}) {

	// Same type? we're done here

	switch v1.(type) {

	case string:
		switch v2.(type) {
		case string:
			return v1, v2
		case int:
			return v1, strconv.Itoa(v2.(int))
		case float64:
			return v1, fmt.Sprintf("%v", v2.(float64))
		case bool:
			if v2.(bool) {
				return v1, "true"
			}
			return v1, "false"
		}

	case float64:
		switch v2.(type) {
		case string:
			return fmt.Sprintf("%v", v1.(float64)), v2
		case int:
			return v1, float64(v2.(int))
		case float64:
			return v1, v2
		case bool:
			if v2.(bool) {
				return v1, 1.0
			}
			return v1, 0.0
		}

	case int:
		switch v2.(type) {
		case string:
			return strconv.Itoa(v1.(int)), v2
		case int:
			return v1, v2
		case float64:
			return float64(v1.(int)), v2
		case bool:
			if v2.(bool) {
				return v1, 1
			}
			return v1, 0
		}
	}
	return v1, v2
}

func (e *Expression) eval1(symbols map[string]interface{}) (interface{}, error) {

	// If the token is a number, convert it

	t := e.Tokens[e.TokenP]
	if i, err := strconv.Atoi(t); err == nil {
		e.TokenP = e.TokenP + 1
		return i, nil
	}

	if i, err := strconv.ParseFloat(t, 64); err == nil {
		e.TokenP = e.TokenP + 1
		return i, nil
	}

	if i, err := strconv.ParseBool(t); err == nil {
		e.TokenP = e.TokenP + 1
		return i, nil
	}

	if i, err := strconv.Atoi(t); err == nil {
		e.TokenP = e.TokenP + 1
		return i, nil
	}

	runeValue := t[0:1]
	if runeValue == "\"" {
		return t[1 : len(t)-1], nil
	}

	if symbol(runeValue) {
		i, found := symbols[t]
		if found {
			return i, nil
		}
		return nil, errors.New("symbol not found: " + t)
	}

	return t, nil
}

func symbol(s string) bool {

	for n, c := range s {
		if isLetter(c) {
			continue
		}

		if isDigit(c) && n > 0 {
			continue
		}
		return false
	}
	return true
}

func isLetter(c rune) bool {
	for _, d := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if c == d {
			return true
		}
	}
	return false
}

func isDigit(c rune) bool {

	for _, d := range "0123456789" {
		if c == d {
			return true
		}
	}
	return false
}
