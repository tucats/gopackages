package expressions

import (
	"errors"
	"strconv"
)

func (e *Expression) expressionAtom(symbols map[string]interface{}) (interface{}, error) {

	t := e.Tokens[e.TokenP]

	// Is this a parenthesis expression?
	if t == "(" {
		e.TokenP = e.TokenP + 1
		v, err := e.relations(symbols)
		if err != nil {
			return nil, err
		}

		if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ")" {
			return nil, errors.New("mismatched parenthesis")
		}
		e.TokenP = e.TokenP + 1
		return v, nil
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

	if symbol(runeValue) {

		// Check for special cases of boolean constants, else look up symbol
		switch t {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:

			// @TOMCOLE Will need to peek ahead for () argument list for
			// function calls here.

			if e.TokenP < len(e.Tokens)-1 && e.Tokens[e.TokenP+1] == "(" {
				e.TokenP = e.TokenP + 1
				return e.functionCall(t, symbols)
			}
			i, found := symbols[t]
			if found {
				e.TokenP = e.TokenP + 1
				return i, nil
			}
			return nil, errors.New("symbol not found: " + t)

		}
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
