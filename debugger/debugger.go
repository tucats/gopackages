package debugger

import (
	"fmt"
	"strings"

	"github.com/tucats/ego/io"
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/compiler"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

const (
	stepTo  = "Stepped to"
	breakAt = "Break at"
)

// Run a context but allow the debugger to take control as
// needed.
func Run(c *bytecode.Context) error {
	return RunFrom(c, 0)
}

func RunFrom(c *bytecode.Context, pc int) error {
	var err error
	c.SetPC(pc)

	for err == nil {
		err = c.Resume()
		if err != nil && err.Error() == SignalDebugger.Error() {
			err = Debugger(c)
		}
		if err != nil && err.Error() == Stop.Error() {
			return nil
		}
	}
	return err
}

// This is called on AtLine to offer the chance for the debugger to take control.
func Debugger(c *bytecode.Context) error {
	var err error

	line := c.GetLine()
	text := c.GetTokenizer().GetLine(line)
	s := c.GetSymbols()

	prompt := false
	// Are we in single-step mode?
	if c.SingleStep() {
		fmt.Printf("%s:\n  %s %3d, %s\n", stepTo, c.GetModuleName(), line, text)
		prompt = true
	} else {
		prompt = EvaluateBreakpoint(c)
	}

	for prompt {
		var tokens *tokenizer.Tokenizer

		for {
			// cmd := ""
			cmd := getLine()
			if len(strings.TrimSpace(cmd)) == 0 {
				cmd = "step"
			}
			tokens = tokenizer.New(cmd)
			if !tokens.AtEnd() {
				break
			}
		}

		// We have a command now in the tokens buffer.
		if err == nil {
			t := tokens.Peek(1)
			switch t {
			case "help":
				_ = Help()

			case "go", "continue":
				c.SetSingleStep(false)
				prompt = false

			case "step":
				c.SetSingleStep(true)
				c.SetStepOver(false)
				prompt = false
				if tokens.Peek(2) == "over" {
					c.SetStepOver(true)
				} else {
					if tokens.Peek(2) == "into" {
						c.SetStepOver(false)
					} else {
						if tokens.Peek(2) != tokenizer.EndOfTokens {
							err = fmt.Errorf("unrecognized step type: %s", tokens.Peek(2))
							c.SetSingleStep(false)
							prompt = true
						}
					}
				}

			case "show":
				err = Show(s, tokens, line, c)

			case "set":
				err = runAfterFirstToken(s, tokens)

			case "call":
				err = runAfterFirstToken(s, tokens)

			case "print":
				text := "fmt.Println(" + strings.Replace(tokens.GetSource(), "print", "", 1) + ")"
				t2 := tokenizer.New(text)
				err = compiler.Run("debugger", s, t2)
				if err != nil && err.Error() == Stop.Error() {
					err = nil
				}
			case "break":
				err = Break(c, tokens)

			case "exit":
				return Stop

			default:
				err = fmt.Errorf("unrecognized command: %s", t)
			}
			if err != nil && err.Error() != Stop.Error() && err.Error() != StepOver.Error() {
				fmt.Printf("Debugger error, %v\n", err)
				err = nil
			}
			if err != nil && err.Error() == Stop.Error() {
				err = nil
				prompt = false
			}
		}
	}
	return err
}

func runAfterFirstToken(s *symbols.SymbolTable, t *tokenizer.Tokenizer) error {
	verb := t.GetTokens(0, 1, false)
	text := strings.TrimPrefix(strings.TrimSpace(t.GetSource()), verb)
	t2 := tokenizer.New(text)
	return compiler.Run("debugger", s, t2)
}

// getLine reads a line of text from the console, and requires that it contain matching
// tick-quotes and braces.
func getLine() string {

	text := io.ReadConsoleText("debug> ")
	if len(strings.TrimSpace(text)) == 0 {
		return ""
	}

	t := tokenizer.New(text)
	for {
		braceCount := 0
		parenCount := 0
		bracketCount := 0
		openTick := false
		lastToken := t.Tokens[len(t.Tokens)-1]
		if lastToken[0:1] == "`" && lastToken[len(lastToken)-1:] != "`" {
			openTick = true
		}
		if !openTick {
			for _, v := range t.Tokens {
				switch v {
				case "[":
					bracketCount++
				case "]":
					bracketCount--
				case "(":
					parenCount++
				case ")":
					parenCount--
				case "{":
					braceCount++
				case "}":
					braceCount--
				}
			}
		}
		if braceCount > 0 || parenCount > 0 || bracketCount > 0 || openTick {
			text = text + io.ReadConsoleText(".....> ")
			t = tokenizer.New(text)
			continue
		} else {
			break
		}
	}
	return text
}
