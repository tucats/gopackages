package expressions

import (
	"strconv"
	"strings"

	"github.com/tucats/gopackages/bytecode"
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

	if t == "true" || t == "false" {
		e.t.Advance(1)
		e.b.Emit2(bc.Push, (t == "true"))
		return nil
	}

	runeValue := t[0:1]
	if runeValue == "\"" {
		e.t.Advance(1)
		s, err := strconv.Unquote(t)
		e.b.Emit2(bc.Push, s)
		return err
	}

	if tokenizer.IsSymbol(t) {

		e.t.Advance(1)
		t := strings.ToLower(t)

		// Nope, probably name from the symbol table
		e.b.Emit2(bc.Load, t)

		return nil

	}

	return e.NewTokenError("unrecognized or unexpected token")
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
	t1 := 1
	var err error

	// Let's experimenally see if this is a range constant expression. This can be
	// of the form [start:end] which creates an array of integers between the start
	// and end values (inclusive). It can also be of the form [:end] which assumes
	// a start number of 1.

	if e.t.Peek(1) == ":" {
		err = nil
		e.t.Advance(-1)
	} else {
		t1, err = strconv.Atoi(e.t.Peek(1))
	}
	if err == nil {
		if e.t.Peek(2) == ":" {
			t2, err := strconv.Atoi(e.t.Peek(3))
			if err == nil {
				e.t.Advance(3)
				count := t2 - t1 + 1

				if count < 0 {
					count = (-count) + 2

					for n := t1; n >= t2; n = n - 1 {
						e.b.Emit2(bytecode.Push, n)
					}

				} else {
					for n := t1; n <= t2; n = n + 1 {
						e.b.Emit2(bytecode.Push, n)
					}
				}
				e.b.Emit2(bytecode.Array, count)
				if !e.t.IsNext("]") {
					return e.NewError("invalid array range constant")
				}
				return nil
			}
		}
	}
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
	var err error

	e.t.Advance(1)
	count := 0

	for e.t.Peek(1) != listTerminator {

		// First element: name

		name := e.t.Next()

		if len(name) > 2 && name[0:1] == "\"" {
			name, err = strconv.Unquote(name)
			if err != nil {
				return err
			}
		} else {
			if !tokenizer.IsSymbol(name) {
				return e.NewTokenError("invalid member name")
			}
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
	return err
}
