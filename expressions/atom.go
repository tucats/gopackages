package expressions

import (
	"errors"
	"strconv"
	"strings"
)

func (e *Expression) expressionAtom(symbols map[string]interface{}) (interface{}, error) {

	t := e.Tokens[e.TokenP]

	// Is this a parenthesis expression?
	if t == "(" {
		e.TokenP = e.TokenP + 1
		v, err := e.conditional(symbols)
		if err != nil {
			return nil, err
		}

		if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ")" {
			return nil, errors.New("mismatched parenthesis")
		}
		e.TokenP = e.TokenP + 1
		return v, nil
	}

	// Is this an array constant?
	if t == "[" {
		return e.parseArray(symbols)
	}

	// If the token is a number, convert it
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
		e.TokenP = e.TokenP + 1
		return t[1 : len(t)-1], nil
	}

	if symbol(t) {

		t := strings.ToLower(t)

		// Peek ahead to see if it's the start of a function call...
		if e.TokenP < len(e.Tokens)-1 && e.Tokens[e.TokenP+1] == "(" {
			e.TokenP = e.TokenP + 1
			return e.functionCall(t, symbols)
		}

		// Nope, resolve from the symbol table if possible...
		i, found := symbols[t]
		if found {
			e.TokenP = e.TokenP + 1

			// But before we go, make sure it's not an array reference...
			if e.TokenP < len(e.Tokens)-1 && e.Tokens[e.TokenP] == "[" {
				e.TokenP = e.TokenP + 1
				idx, err := e.conditional(symbols)
				if err != nil {
					return nil, err
				}
				if e.TokenP > len(e.Tokens)-1 || e.Tokens[e.TokenP] != "]" {
					return nil, errors.New("missing ] in array reference")
				}
				switch a := i.(type) {
				case []interface{}:
					i = a[GetInt(idx)-1]
				default:
					return nil, errors.New("invalid array reference")
				}
			}
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

func (e *Expression) parseArray(symbols map[string]interface{}) ([]interface{}, error) {

	args := make([]interface{}, 0)

	var listTerminator = ""
	if e.Tokens[e.TokenP] == "(" {
		listTerminator = ")"
	}
	if e.Tokens[e.TokenP] == "[" {
		listTerminator = "]"
	}
	if listTerminator == "" {
		return args, nil
	}
	e.TokenP = e.TokenP + 1
	for e.Tokens[e.TokenP] != listTerminator {
		v, err := e.conditional(symbols)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
		if e.TokenP >= len(e.Tokens) {
			break
		}
		if e.Tokens[e.TokenP] == listTerminator {
			break
		}
		if e.Tokens[e.TokenP] != "," {
			return nil, errors.New("invalid list")
		}
		e.TokenP = e.TokenP + 1
	}

	e.TokenP = e.TokenP + 1
	return args, nil
}
