package debugger

import (
	"fmt"
	"strings"

	"github.com/chzyer/readline"
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
	stepTo = "Stepped to"
)

type breakPoint struct {
	kind  breakPointType
	line  int
	hit   bool
	value interface{}
}

var singleStep bool = true
var reader *readline.Instance = nil

// This is called on AtLine to offer the chance for the debugger to take control.
func Debugger(s *symbols.SymbolTable, line int, text string) error {
	var err error
	prompt := false
	// Are we in single-step mode?
	if singleStep {
		fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)
		prompt = true
	}

	for prompt {
		if reader == nil {
			reader, err = readline.New("debug> ")
			if err != nil {
				return err
			}
		}
		var tokens *tokenizer.Tokenizer

		for {
			cmd, err := reader.Readline()
			if err != nil {
				return err
			}
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

			case "print":
				err = compiler.Run(s, tokens)

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
