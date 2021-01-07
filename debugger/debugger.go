package debugger

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tucats/ego/io"
	"github.com/tucats/gopackages/compiler"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

const (
	stepTo  = "Stepped to"
	breakAt = "Break at"
)

var singleStep bool = true

// This is called on AtLine to offer the chance for the debugger to take control.
func Debugger(s *symbols.SymbolTable, module string, line int, text string) error {
	var err error
	prompt := false
	// Are we in single-step mode?
	if singleStep {
		fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)
		prompt = true
	} else {
		prompt = EvaluateBreakpoint(s, module, line, text)
	}

	for prompt {
		var tokens *tokenizer.Tokenizer

		for {
			// cmd := "break when a == 55"
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
				singleStep = false
				prompt = false

			case "step":
				singleStep = true
				prompt = false

			case "show":
				err = Show(s, tokens, line, text)

			case "set":
				err = runAfterFirstToken(s, tokens)

			case "call":
				err = runAfterFirstToken(s, tokens)

			case "print":
				text := "fmt.Println(" + strings.Replace(tokens.GetSource(), "print", "", 1) + ")"
				t2 := tokenizer.New(text)
				err = compiler.Run("debugger", s, t2)

			case "break":
				err = Break(tokens)

			case "exit":
				return errors.New("stop")

			default:
				fmt.Printf("Unrecognized command: %s\n", t)
			}
			if err != nil {
				if err.Error() != "stop" {
					fmt.Printf("Debugger error, %v\n", err)
				}
				err = nil
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
