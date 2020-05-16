package compiler

import (
	"errors"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// Statement parses a single statement
func (c *Compiler) Statement() error {

	// We just eat statement separators
	if c.t.IsNext(";") {
		return nil
	}

	if c.t.IsNext(tokenizer.EndOfTokens) {
		return nil
	}

	// Statement block
	if c.t.IsNext("{") {
		return c.Block()
	}

	if c.t.IsNext("function") {
		return c.Function()
	}

	// It's a single statement, so let's drop down
	// a linenumber marker
	c.b.Emit2(bytecode.AtLine, c.t.Line[c.t.TokenP])

	// Crude assignment statement test
	if c.IsLValue() {

		lv, err := c.LValue()
		if err != nil {
			return err
		}
		if !c.t.IsNext(":=") {
			return errors.New("expected := not found")
		}

		bc, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		c.b.Append(bc)
		c.b.Append(lv)
		return nil
	}

	if c.t.IsNext("if") {
		return c.If()
	}

	if c.t.IsNext("for") {
		return c.For()
	}

	if c.t.IsNext("print") {
		return c.Print()
	}

	if c.t.IsNext("call") {
		return c.Call()
	}

	if c.t.IsNext("return") {
		return c.Return()
	}

	if c.t.IsNext("array") {
		return c.Array()
	}

	c.t.Next()
	return c.NewTokenError("unrecognized or unexpected token")
}
