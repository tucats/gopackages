package debugger

import (
	"fmt"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

func Show(s *symbols.SymbolTable, tokens *tokenizer.Tokenizer, line int, text string) error {
	t := tokens.Peek(2)
	var err error
	switch t {

	case "breaks", "breakpoints":
		ShowBreaks()

	case "symbols", "syms":
		fmt.Println(s.Format(false))

	case "line":
		fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)

	default:
		err = fmt.Errorf("unreognized show command: %s", t)
	}
	return err
}
