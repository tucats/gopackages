package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Print compiles a print statement. The verb is already removed
// from the token stream
func (c *Compiler) Print() error {

	newline := true
	for !c.StatementEnd() {
		if c.t.IsNext(",") {
			return c.NewTokenError(UnexpectedTokenError)
		}
		bc, err := expressions.Compile(c.t)
		if err != nil {
			return err
		}
		newline = true
		c.b.Append(bc)
		c.b.Emit1(bytecode.Print)

		if !c.t.IsNext(",") {
			break
		}
		newline = false
	}
	if newline {
		c.b.Emit1(bytecode.Newline)
	}
	return nil
}
