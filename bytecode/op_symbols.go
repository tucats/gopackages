package bytecode

import (
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*   S Y M B O L S   A N D  T A B L E S    *
*                                         *
\******************************************/

// PushScopeOpcode implementation
func PushScopeOpcode(c *Context, i interface{}) error {

	s := symbols.NewChildSymbolTable("statement block", c.symbols)
	c.symbols = s
	return nil
}

// PopScopeOpcode implementation
func PopScopeOpcode(c *Context, i interface{}) error {
	c.symbols = c.symbols.Parent
	return nil
}

// SymbolCreateOpcode implementation
func SymbolCreateOpcode(c *Context, i interface{}) error {

	n := util.GetString(i)
	if c.IsConstant(n) {
		return c.NewError("attmpt to write to constant")
	}
	err := c.Create(n)
	if err != nil {
		err = c.NewError(err.Error())
	}
	return err
}

// SymbolDeleteOpcode implementation
func SymbolDeleteOpcode(c *Context, i interface{}) error {

	n := util.GetString(i)
	err := c.Delete(n)
	if err != nil {
		err = c.NewError(err.Error())
	}
	return err
}

// ConstantOpcode implementation
func ConstantOpcode(c *Context, i interface{}) error {

	v, err := c.Pop()
	if err != nil {
		return err
	}

	varname := util.GetString(i)
	err = c.SetConstant(varname, v)
	if err != nil {
		return c.NewError(err.Error())
	}

	return err
}
