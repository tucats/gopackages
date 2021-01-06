package debugger

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tucats/ego/io"
	"github.com/tucats/gopackages/compiler"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

type breakPointType int

const (
	BreakDisabled breakPointType = 0
	BreakAlways   breakPointType = iota
	BreakValue
)

const (
	stepTo  = "Stepped to"
	breakAt = "Break at"
)

type breakPoint struct {
	kind breakPointType
	line int
	hit  int
}

var breakPoints = []breakPoint{}
var singleStep bool = true

// This is called on AtLine to offer the chance for the debugger to take control.
func Debugger(s *symbols.SymbolTable, line int, text string) error {
	var err error
	prompt := false
	// Are we in single-step mode?
	if singleStep {
		fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)
		prompt = true
	} else {
		for _, b := range breakPoints {
			fmt.Printf("Evaluating break %d == %d\n", line, b.line)
			if line == b.line {
				prompt = true
				fmt.Printf("%s:\n\t%5d, %s\n", breakAt, line, text)
				b.hit++
			}
		}
	}

	for prompt {
		var tokens *tokenizer.Tokenizer

		for {
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
			case "go", "continue":
				singleStep = false
				prompt = false

			case "step":
				singleStep = true
				prompt = false

			case "show":
				t := tokens.Peek(2)
				switch t {
				case "symbols", "syms":
					fmt.Println(s.Format(false))
				case "line":
					fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)
				default:
					err = fmt.Errorf("unreognized show command: %s", t)
				}

			case "set":
				err = runAfterFirstToken(s, tokens)

			case "call":
				err = runAfterFirstToken(s, tokens)

			case "print":
				text := "fmt.Println(" + strings.Replace(tokens.GetSource(), "print", "", 1) + ")"
				t2 := tokenizer.New(text)
				err = compiler.Run(s, t2)

			case "break":
				err = Break(tokens)

			case "exit":
				return errors.New("exit from debugger")

			default:
				fmt.Printf("Unrecognized command: %s\n", t)
			}

			if err != nil {
				fmt.Printf("Debugger error, %v\n", err)
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
	return compiler.Run(s, t2)
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

func Break(t *tokenizer.Tokenizer) error {
	var err error
	var line int
	t.Advance(1)

	for t.Peek(1) != tokenizer.EndOfTokens {
		switch t.Next() {
		case "at":
			line, err = strconv.Atoi(t.Next())
			if err == nil {
				err = breakAtLine(line)
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

func breakAtLine(line int) error {
	b := breakPoint{line: line, hit: 0, kind: BreakAlways}
	breakPoints = append(breakPoints, b)
	return nil
}
