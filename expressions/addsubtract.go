package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

func (e *Expression) addSubtract() error {

	err := e.multDivide()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if inList(op, []string{"+", "-", "&"}) {
			e.TokenP = e.TokenP + 1

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
