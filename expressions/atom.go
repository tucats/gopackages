package expressions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	bc "github.com/tucats/gopackages/bytecode"
)

func (e *Expression) expressionAtom() error {

	t := e.Tokens[e.TokenP]

	// Is this a parenthesis expression?
	if t == "(" {
		e.TokenP = e.TokenP + 1
		err := e.conditional()
		if err != nil {
			return err
		}

		if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ")" {
			return errors.New("mismatched parenthesis")
		}
		e.TokenP = e.TokenP + 1
		return nil
	}

	// Is this an array constant?
	if t == "[" {
		return e.parseArray()
	}

	// Is it a struct constant?
	if t == "{" {
		return e.parseStruct()
	}
	// If the token is a number, convert it
	if i, err := strconv.Atoi(t); err == nil {
		e.TokenP = e.TokenP + 1
		e.b.Emit(bc.Push, i)
		return nil
	}

	if i, err := strconv.ParseFloat(t, 64); err == nil {
		e.TokenP = e.TokenP + 1
		e.b.Emit(bc.Push, i)
		return nil
	}

	if i, err := strconv.ParseBool(t); err == nil {
		e.TokenP = e.TokenP + 1
		e.b.Emit(bc.Push, i)
		return nil
	}

	runeValue := t[0:1]
	if runeValue == "\"" {
		e.TokenP = e.TokenP + 1
		e.b.Emit(bc.Push, t[1:len(t)-1])
		return nil
	}

	if symbol(t) {

		t := strings.ToLower(t)

		// Peek ahead to see if it's the start of a function call...
		if e.TokenP < len(e.Tokens)-1 && e.Tokens[e.TokenP+1] == "(" {
			e.TokenP = e.TokenP + 1
			return e.functionCall(t)
		}

		// Nope, probably name from the symbol table
		e.TokenP = e.TokenP + 1
		e.b.Emit(bc.Load, t)

		// But before we go, make sure it's not an array reference...
		if e.TokenP < len(e.Tokens)-1 && e.Tokens[e.TokenP] == "[" {
			e.TokenP = e.TokenP + 1
			err := e.conditional()
			if err != nil {
				return err
			}
			if e.TokenP > len(e.Tokens)-1 || e.Tokens[e.TokenP] != "]" {
				return errors.New("missing ] in array reference")
			}
			e.b.Emit(bc.Index, 0)
		}

		return nil

	}

	e.b.Emit(bc.Push, t)
	return nil
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
	for _, d := range "_abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" {
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

func (e *Expression) parseArray() error {

	var listTerminator = ""
	if e.Tokens[e.TokenP] == "(" {
		listTerminator = ")"
	}
	if e.Tokens[e.TokenP] == "[" {
		listTerminator = "]"
	}
	if listTerminator == "" {
		return nil
	}
	e.TokenP = e.TokenP + 1
	count := 0

	for e.Tokens[e.TokenP] != listTerminator {
		err := e.conditional()
		if err != nil {
			return err
		}
		count = count + 1
		if e.TokenP >= len(e.Tokens) {
			break
		}
		if e.Tokens[e.TokenP] == listTerminator {
			break
		}
		if e.Tokens[e.TokenP] != "," {
			return errors.New("invalid list")
		}
		e.TokenP = e.TokenP + 1
	}

	e.b.Emit(bc.Array, count)

	e.TokenP = e.TokenP + 1
	return nil
}

func (e *Expression) parseStruct() error {

	var listTerminator = "}"

	e.TokenP = e.TokenP + 1
	count := 0

	for e.Tokens[e.TokenP] != listTerminator {

		// First element: name

		name := e.Tokens[e.TokenP]
		if !symbol(name) {
			return fmt.Errorf("invalid member name: %v", name)
		}

		// Second element: colon
		e.TokenP = e.TokenP + 1
		if e.Tokens[e.TokenP] != ":" {
			return errors.New("missing colon")
		}

		// Third element: value, which is emitted.
		e.TokenP = e.TokenP + 1
		err := e.conditional()
		if err != nil {
			return err
		}
		// Now write the name as a string.
		e.b.Emit(bc.Push, name)

		count = count + 1
		if e.TokenP >= len(e.Tokens) {
			break
		}
		if e.Tokens[e.TokenP] == listTerminator {
			break
		}
		if e.Tokens[e.TokenP] != "," {
			return errors.New("invalid list")
		}
		e.TokenP = e.TokenP + 1
	}

	e.b.Emit(bc.Struct, count)

	e.TokenP = e.TokenP + 1
	return nil
}
