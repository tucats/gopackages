package debugger

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/compiler"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

type breakPointType int

const (
	BreakDisabled breakPointType = 0
	BreakAlways   breakPointType = iota
	BreakValue
)

type breakPoint struct {
	kind   breakPointType
	module string
	line   int
	hit    int
	expr   *bytecode.ByteCode
	text   string
}

var breakPoints = []breakPoint{}

func Break(c *bytecode.Context, t *tokenizer.Tokenizer) error {
	var err error
	var line int
	t.Advance(1)

	for t.Peek(1) != tokenizer.EndOfTokens {
		switch t.Next() {
		case "when":
			text := t.GetTokens(2, len(t.Tokens), true)
			ec := compiler.New().WithTokens(tokenizer.New(text))
			bc, err := ec.Expression()
			if err == nil {
				err = breakWhen(bc, text)
				if err != nil {
					return err
				}
			}
			t.Advance(999)

		case "at":
			name := t.Next()
			if t.Peek(1) == ":" {
				t.Advance(1)
			} else {
				name = c.GetModuleName()
				t.Advance(-1)
			}
			line, err = strconv.Atoi(t.Next())
			if err == nil {
				err = breakAtLine(name, line)
			}
		default:
			err = errors.New(InvalidBreakClauseError)
		}

		if err != nil {
			break
		}
	}
	return err
}

func breakAtLine(module string, line int) error {
	b := breakPoint{
		module: module,
		line:   line,
		hit:    0,
		kind:   BreakAlways,
	}
	breakPoints = append(breakPoints, b)
	fmt.Printf("Added break %s\n", FormatBreakpoint(b))
	return nil
}

func breakWhen(expression *bytecode.ByteCode, text string) error {
	b := breakPoint{
		module: "expression",
		hit:    0,
		kind:   BreakValue,
		expr:   expression,
		text:   text,
	}
	breakPoints = append(breakPoints, b)
	fmt.Printf("Added break %s\n", FormatBreakpoint(b))
	return nil
}

func ShowBreaks() {
	if len(breakPoints) == 0 {
		fmt.Printf("No breakpoints defined\n")
	} else {
		for _, b := range breakPoints {
			fmt.Printf("break %s\n", FormatBreakpoint(b))
		}
	}

}

func FormatBreakpoint(b breakPoint) string {
	switch b.kind {
	case BreakAlways:
		return fmt.Sprintf("at %s:%d", b.module, b.line)
	case BreakValue:
		return fmt.Sprintf("when %s", b.text)
	default:
		return fmt.Sprintf("(undefined) %v", b)
	}
}

// Using the current execution state, determine if a breakpoint has
// been encountered.
func EvaluateBreakpoint(c *bytecode.Context) bool {
	s := c.GetSymbols()
	msg := ""
	prompt := false
	for _, b := range breakPoints {
		switch b.kind {
		case BreakValue:

			// fmt.Printf("DEBUG: break eval sym table\n%s\n", s.Format(true))
			// If we already hit this, don't do it again on each statement. Pass.
			if b.hit > 0 {
				break
			}
			ctx := bytecode.NewContext(s, b.expr)
			ctx.SetDebug(false)
			err := ctx.Run()
			if err != nil {
				if err.Error() == StepOver.Error() {
					err = nil
					ctx.StepOver(true)
				}
				if err.Error() == SignalDebugger.Error() {
					err = nil
				}
			}
			//fmt.Printf("Break expression status = %v\n", err)
			if err == nil {
				if v, err := ctx.Pop(); err == nil {
					//fmt.Printf("Break expression result = %v\n", v)
					prompt = util.GetBool(v)
					if prompt {
						b.hit++
					} else {
						b.hit = 0
					}
				}
			}
			msg = "Break when " + b.text

		case BreakAlways:
			line := c.GetLine()
			module := c.GetModuleName()
			// fmt.Printf("Evaluating %s:%d = %s\n", module, line, text)
			if module == b.module && line == b.line {
				prompt = true
				text := c.GetTokenizer().GetLine(line)
				msg = fmt.Sprintf("%s:\n\t%5d, %s", breakAt, line, text)
				b.hit++
				break
			}
		}
	}
	if prompt {
		fmt.Printf("%s\n", msg)
	}
	return prompt
}
