package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

func (e *Expression) relations() error {

	err := e.addSubtract()
	if err != nil {
		return err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if op == "=" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">=" {
			e.TokenP = e.TokenP + 1

			err := e.addSubtract()
			if err != nil {
				return err
			}

			switch op {

			case "=":
				e.b.Emit(bc.Equal, nil)

			case "!=":
				e.b.Emit(bc.NotEqual, nil)

			case "<":
				e.b.Emit(bc.LessThan, nil)

			case "<=":
				e.b.Emit(bc.LessThanOrEqual, nil)

			case ">":
				e.b.Emit(bc.GreaterThan, nil)

			case ">=":
				e.b.Emit(bc.GreaterThanOrEqual, nil)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
