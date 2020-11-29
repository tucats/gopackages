package compiler

import (
	"github.com/tucats/gopackages/expressions"
)

// Assignment compiles an assignment statement.
func (c *Compiler) Assignment() error {

	storeLValue, err := c.LValue()
	if err != nil {
		return err
	}
	if !c.t.AnyNext([]string{":=", "="}) {
		return c.NewError(MissingAssignmentError)
	}

	expressionCode, err := expressions.Compile(c.t)
	if err != nil {
		return err
	}
	c.b.Append(expressionCode)
	c.b.Append(storeLValue)
	return nil

}
