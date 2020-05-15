package expressions

import (
	"strconv"
	"strings"

	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

func (e *Expression) expressionAtom() error {

	t := e.t.Peek(1)

	// Is this a parenthesis expression?
	if t == "(" {
		e.t.Advance(1)
		err := e.conditional()
		if err != nil {
			return err
		}

		if e.t.Next() != ")" {
			return e.NewError("mismatched parenthesis")
		}
		return nil
	}

	// Is this an array constant?
	if t == "[" {
		return e.parseArray()
	}

	// Is it a map constant?
	if t == "{" {
		return e.parseStruct()
	}
	// If the token is a number, convert it
	if i, err := strconv.Atoi(t); err == nil {
		e.t.Advance(1)
		e.b.Emit2(bc.Push, i)
		return nil
	}

	if i, err := strconv.ParseFloat(t, 64); err == nil {
		e.t.Advance(1)
		e.b.Emit2(bc.Push, i)
		return nil
	}

	if i, err := strconv.ParseBool(t); err == nil {
		e.t.Advance(1)
		e.b.Emit2(bc.Push, i)
		return nil
	}

	runeValue := t[0:1]
	if runeValue == "\"" {
		e.t.Advance(1)
		e.b.Emit2(bc.Push, t[1:len(t)-1])
		return nil
	}

	if tokenizer.IsSymbol(t) {

		e.t.Advance(1)
		t := strings.ToLower(t)

		// Peek ahead to see if it's the start of a function call...
		if e.t.IsNext("(") {
			return e.functionCall(t)
		}

		// Nope, probably name from the symbol table
		e.b.Emit2(bc.Load, t)

		return nil

	}

	e.b.Emit2(bc.Push, t)
	return nil
}

func (e *Expression) parseArray() error {

	var listTerminator = ""
	if e.t.Peek(1) == "(" {
		listTerminator = ")"
	}
	if e.t.Peek(1) == "[" {
		listTerminator = "]"
	}
	if listTerminator == "" {
		return nil
	}
	e.t.Advance(1)
	count := 0

	for e.t.Peek(1) != listTerminator {
		err := e.conditional()
		if err != nil {
			return err
		}
		count = count + 1
		if e.t.AtEnd() {
			break
		}
		if e.t.Peek(1) == listTerminator {
			break
		}
		if e.t.Peek(1) != "," {
			return e.NewError("invalid list")
		}
		e.t.Advance(1)
	}

	e.b.Emit2(bc.Array, count)

	e.t.Advance(1)
	return nil
}

func (e *Expression) parseStruct() error {

	var listTerminator = "}"

	e.t.Advance(1)
	count := 0

	for e.t.Peek(1) != listTerminator {

		// First element: name

		name := e.t.Next()
		if !tokenizer.IsSymbol(name) {
			return e.NewTokenError("invalid member name")
		}

		// Second element: colon
		if e.t.Next() != ":" {
			return e.NewError("missing colon")
		}

		// Third element: value, which is emitted.
		err := e.conditional()
		if err != nil {
			return err
		}
		// Now write the name as a string.
		e.b.Emit2(bc.Push, name)

		count = count + 1
		if e.t.AtEnd() {
			break
		}
		if e.t.Peek(1) == listTerminator {
			break
		}
		if e.t.Peek(1) != "," {
			return e.NewError("invalid list")
		}
		e.t.Advance(1)
	}

	e.b.Emit2(bc.Struct, count)
	e.t.Advance(1)
	return nil
}
