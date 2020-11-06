package compiler

import (
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/expressions"
	"github.com/tucats/gopackages/tokenizer"
)

// Directive processes a compiler directive. These become symbols generated
// at compile time that are copied to the compiler's symbol table for processing
// elsewhere.
func (c *Compiler) Directive() error {

	var err error

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewStringError("invalid directive name", name)
	}

	switch name {

	case "template":
		return c.Template()

	default:
		// Assume it is to be stored in the global directives structure

		value, err := expressions.NewWithTokenizer(c.t).Eval(c.s)
		if err == nil {

			v, f := c.s.Get(DirectiveStructureName)
			if !f {
				v = map[string]interface{}{}
			}
			m := v.(map[string]interface{})
			m[name] = value
			c.s.SetAlways(DirectiveStructureName, m)
		}
	}
	return err
}

// Template implements the template compiler directive
func (c *Compiler) Template() error {

	// Get the template name
	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewStringError("invalid directive name", name)
	}

	// Get the template string definition
	bc, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Emit2(bytecode.Template, name)
	c.b.Emit2(bytecode.SymbolCreate, name)
	c.b.Emit2(bytecode.Store, name)

	return nil
}
