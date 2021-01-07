package debugger

import (
	"fmt"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

func Show(s *symbols.SymbolTable, tokens *tokenizer.Tokenizer, line int, tx *tokenizer.Tokenizer) error {
	t := tokens.Peek(2)
	var err error
	switch t {

	case "breaks", "breakpoints":
		ShowBreaks()

	case "symbols", "syms":
		fmt.Println(s.Format(false))

	case "line":
		text := tx.GetLine(line)
		fmt.Printf("%s:\n\t%5d, %s\n", stepTo, line, text)

	case "scope":
		syms := s
		fmt.Printf("Symbol table scope:\n")
		for syms.Parent != nil {
			fmt.Printf("\t%s\n", syms.Name)
			syms = syms.Parent
		}

	case "source":
		for i, t := range tx.Source {
			fmt.Printf("%-5d %s\n", i+1, t)
		}

	default:
		err = fmt.Errorf("unreognized show command: %s", t)
	}
	return err
}
