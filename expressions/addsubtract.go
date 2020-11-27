package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/tokenizer"
)

func (e *Expression) addSubtract() error {

	err := e.multDivide()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if e.t.AtEnd() {
			break
		}
		op := e.t.Peek(1)
		if tokenizer.InList(op, []string{"+", "-", "&"}) {
			e.t.Advance(1)

			if e.t.IsNext(tokenizer.EndOfTokens) {
				return e.NewError(MissingTermError)
			}

			err := e.multDivide()
			if err != nil {
				return err
			}

			switch op {

			case "+":
				e.b.Emit1(bc.Add)

			case "-":
				e.b.Emit1(bc.Sub)

			case "&":
				e.b.Emit1(bc.And)
			}

		} else {
			parsing = false
		}
	}
	return nil
}
