package bytecode

import (
	"fmt"

	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*           B A S I C   I / O             *
*                                         *
\******************************************/

// PrintOpcode implementation. If the operand
// is given, it represents the number of items
// to remove from the stack.
func PrintOpcode(c *Context, i interface{}) error {

	count := 1
	if i != nil {
		count = util.GetInt(i)
	}

	for n := 0; n < count; n = n + 1 {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		s := fmt.Sprintf("%s", util.FormatUnquoted(v))
		if c.output == nil {
			fmt.Printf("%s", s)
		} else {
			c.output.WriteString(s)
		}
	}

	// If we are instruction tracing, print out a newline anyway so the trace
	// display isn't made illegible.
	if c.output == nil && c.Tracing {
		fmt.Println()
	}

	return nil
}

// NewlineOpcode implementation.
func NewlineOpcode(c *Context, i interface{}) error {

	if c.output == nil {
		fmt.Printf("\n")
	} else {
		c.output.WriteString("\n")
	}
	return nil
}
