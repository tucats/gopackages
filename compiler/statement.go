package compiler

import (
	"errors"
	"fmt"

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

	if c.t.IsNext("print") {
		return c.Print()
	}

	if c.t.IsNext("call") {
		return c.Call()
	}

	if c.t.IsNext("return") {
		return c.Return()
	}

	if c.t.IsNext("function") {
		return c.Function()
	}

	return fmt.Errorf("unrecognized statement: %s", c.t.Peek(1))
}
