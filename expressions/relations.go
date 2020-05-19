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
		if e.t.AtEnd() {
			break
		}
		op := e.t.Peek(1)
		if op == "==" || op == "!=" || op == "<" || op == "<=" || op == ">" || op == ">=" {
			e.t.Advance(1)

			err := e.addSubtract()
			if err != nil {
				return err
			}

			switch op {

			case "==":
				e.b.Emit1(bc.Equal)

			case "!=":
				e.b.Emit1(bc.NotEqual)

			case "<":
				e.b.Emit1(bc.LessThan)

			case "<=":
				e.b.Emit1(bc.LessThanOrEqual)

			case ">":
				e.b.Emit1(bc.GreaterThan)

			case ">=":
				e.b.Emit1(bc.GreaterThanOrEqual)

			}

		} else {
			parsing = false
		}
	}
	return nil
}
