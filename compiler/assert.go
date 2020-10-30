package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Assert processes an assert statment. This is an expression that
// must be resolvable to true else a fatal error is generated
func (c *Compiler) Assert() error {

	c.b.Emit1(bytecode.PushScope)
	c.blockDepth = c.blockDepth + 1

	exprStart := c.t.Mark()

	expressionCode, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(expressionCode)
	fixup := c.b.Mark()
	c.b.Emit2(bytecode.BranchTrue, 0)

	if c.StatementEnd() {
		exprEnd := c.t.Mark()
		msg := c.t.GetTokens(exprStart, exprEnd, false)
		c.b.Emit2(bytecode.Push, "assertion failed: "+msg)
		c.b.Emit2(bytecode.Panic, true)
	} else {
		c.t.IsNext(",")
		expressionCode, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		c.b.Append(expressionCode)
		c.b.Emit2(bytecode.Panic, true)
	}
	c.b.SetAddress(fixup, c.b.Mark())
	c.b.Emit1(bytecode.PopScope)
	c.blockDepth = c.blockDepth - 1
	return nil
}
