package bytecode

import (
	"fmt"
	"text/template"

	"github.com/tucats/gopackages/app-cli/ui"
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
		s := util.FormatUnquoted(v)
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

// SayOpcode implementation. This can be used in place
// of NewLine to end buffered output, but the output is
// only displayed if we are not in --quiet mode.
func SayOpcode(c *Context, i interface{}) error {
	ui.Say("%s\n", c.output.String())
	c.output = nil
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

/******************************************\
*                                         *
*           T E M P L A T E S             *
*                                         *
\******************************************/

// TemplateOpcode compiles a template string from the
// stack and stores it in the template manager for the
// context.
func TemplateOpcode(c *Context, i interface{}) error {

	name := util.GetString(i)
	t, err := c.Pop()
	if err == nil {
		t, err = template.New(name).Parse(util.GetString(t))
		if err == nil {
			err = c.Push(t)
		}
	}
	return err
}
