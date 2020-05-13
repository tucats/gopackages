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
		op := e.t.Peek()
		if tokenizer.InList(op, []string{"+", "-", "&"}) {
			e.t.Advance(1)

			err := e.multDivide()
			if err != nil {
				return err
			}

			switch op {

			case "+":
				e.b.Emit(bc.Add, nil)

			case "-":
				e.b.Emit(bc.Sub, nil)

			case "&":
				e.b.Emit(bc.And, nil)
			}

		} else {
			parsing = false
		}
	}
	return nil
}
