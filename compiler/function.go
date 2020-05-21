package compiler

import (
	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

// Function compiles a function definition
func (c *Compiler) Function() error {

	parameters := []string{}

	fname := c.t.Next()
	if !tokenizer.IsSymbol(fname) {
		return c.NewTokenError("invalid function name")
	}

	// Process parameter names
	if c.t.IsNext("(") {
		for !c.t.IsNext(")") {
			if c.t.AtEnd() {
				break
			}
			name := c.t.Next()
			if tokenizer.IsSymbol(name) {
				parameters = append(parameters, name)
			} else {
				return c.NewTokenError("invalid parameter")
			}
			if c.t.IsNext(",") {
				// No action
			}
		}
	}

	b := bytecode.New(fname)

	// Generate the parameter assignments. These are extracted
	// from the automatic array named _args which is generated
	// as part of the function call during bytecode exectuion.
	// Note that the array is 1-based.
	for n, name := range parameters {
		b.Emit2(bytecode.Load, "_args")
		b.Emit2(bytecode.Push, n+1)
		b.Emit1(bytecode.LoadIndex)
		b.Emit2(bytecode.SymbolCreate, name)
		b.Emit2(bytecode.Store, name)
	}

	// Now compile a statement or block into the function body.
	cInstance := Compiler{b: b, t: c.t, s: c.s}
	cx := &cInstance

	err := cx.Statement()
	if err != nil {
		return err
	}

	// Store anchor to the function, either in the current
	// table or package.

	if c.PackageName == "" {
		c.s.SetAlways(fname, b)
	} else {
		v, found := c.s.Get(fname)
		if !found {
			v := map[string]interface{}{}
			v[fname] = b
			v["__readonly"] = true
			c.s.SetAlways(c.PackageName, v)
		} else {
			switch a := v.(type) {
			case map[string]interface{}:
				a[fname] = b
				c.s.SetAlways(c.PackageName, a)
			default:
				return c.NewError("invalid package object")
			}
		}
	}

	if ui.DebugMode {
		b.Disasm()
	}
	return nil
}
