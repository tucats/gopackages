package bytecode

import (
	"fmt"
	"io/ioutil"
	"text/template"
	"time"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/data"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/tokenizer"
)

/******************************************\
*                                         *
*           B A S I C   I / O             *
*                                         *
\******************************************/

// logByteCode implements the Log directive, which outputs the top stack
// item to the logger named in the operand. The operand can either by a logger
// by name or by class id.
func logByteCode(c *Context, i interface{}) error {
	var class int

	if id, ok := i.(int); ok {
		class = id
	} else {
		class = ui.LoggerByName(data.String(i))
	}

	if class <= ui.NoSuchLogger {
		return c.error(errors.ErrInvalidLoggerName).Context(i)
	}

	msg, err := c.Pop()
	if err != nil {
		return err
	}

	ui.Log(class, "%v", msg)

	return nil
}

// sayByteCode instruction processor. If the operand is true, output the string as-is,
// else output it adding a trailing newline. The Say opcode  can be used in place
// of NewLine to end buffered output, but the output is only displayed if we are
// not in --quiet mode.
//
// This is used by the code generated from @test and @pass, for example, to allow
// test logging to be quiet if necessary.
func sayByteCode(c *Context, i interface{}) error {
	msg := ""
	if c.output != nil {
		msg = c.output.String()
		c.output = nil
	}

	fmt := "%s\n"
	if data.Bool(i) && len(msg) > 0 {
		fmt = "%s"
	}

	ui.Say(fmt, msg)

	return nil
}

// newlineByteCode instruction processor generates a newline character to stdout.
func newlineByteCode(c *Context, i interface{}) error {
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

// templateByteCode compiles a template string from the stack and stores it in
// the template manager for the execution context.
func templateByteCode(c *Context, i interface{}) error {
	name := data.String(i)

	t, err := c.Pop()
	if err == nil {
		if isStackMarker(t) {
			return c.error(errors.ErrFunctionReturnedVoid)
		}

		t, e2 := template.New(name).Parse(data.String(t))
		if e2 == nil {
			err = c.push(t)
		} else {
			err = c.error(e2)
		}
	}

	return err
}

/******************************************\
*                                         *
*             U T I L I T Y               *
*                                         *
\******************************************/

// fromFileByteCode loads the context tokenizer with the
// source from a file if it does not already exist and
// we are in debug mode.
func fromFileByteCode(c *Context, i interface{}) error {
	if !c.debugging {
		return nil
	}

	if b, err := ioutil.ReadFile(data.String(i)); err == nil {
		c.tokenizer = tokenizer.New(string(b), false)

		return nil
	} else {
		return errors.NewError(err)
	}
}

func timerByteCode(c *Context, i interface{}) error {
	mode := data.Int(i)
	switch mode {
	case 0:
		t := time.Now()
		c.timerStack = append(c.timerStack, t)

		return nil

	case 1:
		timerStack := len(c.timerStack)
		if timerStack == 0 {
			return c.error(errors.ErrInvalidTimer)
		}

		t := c.timerStack[timerStack-1]
		c.timerStack = c.timerStack[:timerStack-1]
		now := time.Now()
		elapsed := now.Sub(t)
		ms := elapsed.Milliseconds()
		unit := "s"

		// If the unit scale is too large or too small, then
		// adjust it down to millisends or up to minutes.
		if ms == 0 {
			ms = elapsed.Microseconds()
			unit = "ms"
		} else if ms > 60000 {
			ms = ms / 1000
			unit = "m"
		}

		msText := fmt.Sprintf("%4.3f%s", float64(ms)/1000.0, unit)

		return c.push(msText)

	default:
		return c.error(errors.ErrInvalidTimer).Context(mode)
	}
}
