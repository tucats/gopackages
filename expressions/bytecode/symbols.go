package bytecode

import (
	"strconv"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

/******************************************\
*                                         *
*   S Y M B O L S   A N D  T A B L E S    *
*                                         *
\******************************************/

// pushScopeByteCode instruction processor.
func pushScopeByteCode(c *Context, i interface{}) error {
	oldName := c.symbols.Name

	c.mux.Lock()
	defer c.mux.Unlock()

	c.blockDepth++
	c.symbols = symbols.NewChildSymbolTable("block "+strconv.Itoa(c.blockDepth), c.symbols).Shared(false)

	ui.Log(ui.SymbolLogger, "(%d) push symbol table \"%s\" <= \"%s\"",
		c.threadID, c.symbols.Name, oldName)

	return nil
}

// symbolCreateByteCode instruction processor.
func createAndStoreByteCode(c *Context, i interface{}) error {
	var value interface{}

	var err error

	var name string

	if operands, ok := i.([]interface{}); ok && len(operands) == 2 {
		name = data.String(operands[0])
		value = operands[1]
	} else {
		name = data.String(i)
		value, err = c.Pop()
		if err != nil {
			return err
		}
	}

	if c.isConstant(name) {
		return c.error(errors.ErrReadOnly)
	}

	err = c.create(name)
	if err != nil {
		return c.error(err)
	}

	// If the name starts with "_" it is implicitly a readonly
	// variable.  In this case, make a copy of the value to
	// be stored, and mark it as a readonly value if it is
	// a complex type. Then, store the copy as a constant with
	// the given name.
	if len(name) > 1 && name[0:1] == defs.ReadonlyVariablePrefix {
		constantValue := data.DeepCopy(value)

		err = c.setConstant(name, constantValue)
	} else {
		err = c.set(name, value)
	}

	return err
}

// symbolCreateByteCode instruction processor.
func symbolCreateByteCode(c *Context, i interface{}) error {
	n := data.String(i)
	if c.isConstant(n) {
		return c.error(errors.ErrReadOnly)
	}

	err := c.create(n)
	if err != nil {
		err = c.error(err)
	}

	return err
}

// symbolCreateIfByteCode instruction processor.
func symbolCreateIfByteCode(c *Context, i interface{}) error {
	n := data.String(i)
	if c.isConstant(n) {
		return c.error(errors.ErrReadOnly)
	}

	sp := c.symbols
	if _, found := sp.GetLocal(n); found {
		return nil
	}

	err := c.symbols.Create(n)
	if err != nil {
		err = c.error(err)
	}

	return err
}

// symbolDeleteByteCode instruction processor.
func symbolDeleteByteCode(c *Context, i interface{}) error {
	n := data.String(i)

	err := c.delete(n)
	if err != nil {
		return c.error(err)
	}

	return nil
}

// constantByteCode instruction processor.
func constantByteCode(c *Context, i interface{}) error {
	v, err := c.Pop()
	if err != nil {
		return err
	}

	if isStackMarker(v) {
		return c.error(errors.ErrFunctionReturnedVoid)
	}

	varname := data.String(i)

	err = c.setConstant(varname, v)
	if err != nil {
		return c.error(err)
	}

	return err
}
