package compiler

import (
	"fmt"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
)

// Function compiles a function definition
func (c *Compiler) Function() error {

	parameters := []string{}

	fname := c.t.Next()
	if !expressions.Symbol(fname) {
		return fmt.Errorf("invalid function name: %s", fname)
	}

	// Process parameter names
	if c.t.IsNext("(") {
		for !c.t.IsNext(")") {
			if c.t.AtEnd() {
				break
			}
			name := c.t.Peek(1)
			if expressions.Symbol(name) {
				parameters = append(parameters, name)
			} else {
				return fmt.Errorf("invalid parameter: %s", name)
			}
			if c.t.IsNext(",") {
				// No action
			}
		}
	}

	b := bytecode.New(fname)

	// Generate the parameter assignments
	for n, name := range parameters {
		b.Emit(bytecode.Push, "_args")
		b.Emit(bytecode.Push, n)
		b.Emit0(bytecode.Index)
		b.Emit(bytecode.Store, name)
	}

	// Now compile a statement or block into the function body.
	cInstance := Compiler{b: b, t: c.t}
	cx := &cInstance

	err := cx.Block()
	if err != nil {
		return err
	}

	// Store the compiled code is the compiler's symbol table
	c.s.Set(fname+"()", b)
	return nil
}
