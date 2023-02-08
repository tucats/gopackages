package bytecode

import (
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
)

// loadByteCode instruction processor.
func loadByteCode(c *Context, i interface{}) error {
	name := data.String(i)
	if len(name) == 0 {
		return c.error(errors.ErrInvalidIdentifier).Context(name)
	}

	v, found := c.get(name)
	if !found {
		return c.error(errors.ErrUnknownIdentifier).Context(name)
	}

	return c.push(data.UnwrapConstant(v))
}

// explodeByteCode implements Explode. This accepts a struct on the top of
// the stack, and creates local variables for each of the members of the
// struct by their name.
func explodeByteCode(c *Context, i interface{}) error {
	var err error

	var v interface{}

	v, err = c.Pop()
	if err != nil {
		return err
	}

	if isStackMarker(v) {
		return c.error(errors.ErrFunctionReturnedVoid)
	}

	err = c.error(errors.ErrInvalidType).Context(data.TypeOf(v).String())

	return err
}
