package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// Constant compiles a constant block
func (c *Compiler) Constant() error {

	terminator := ""

	if c.t.IsNext("(") {
		terminator = ")"
	}

	for terminator == "" || !c.t.IsNext(terminator) {
		name := c.t.Next()
		if !tokenizer.IsSymbol(name) {
			return c.NewTokenError(InvalidSymbolError)
		}

		if !c.t.IsNext("=") {
			return c.NewTokenError(MissingEqualError)
		}
		vx, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}

		// Search to make sure it doesn't contain a load statement that isn't for another
		// constant

		for _, i := range vx.Opcodes() {
			if i.Opcode == bytecode.Load && !tokenizer.InList(util.GetString(i.Operand), c.constants) {
				return c.NewError(InvalidConstantError)
			}
		}
		c.constants = append(c.constants, name)

		c.b.Append(vx)
		c.b.Emit2(bytecode.Constant, name)

		if terminator == "" {
			break
		}

	}
	return nil
}
