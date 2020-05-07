package expressions

import "errors"

// Eval evaluates the parsed expression. This can be called multiple times
// with the same scanned string, but with different symbols.
func (e *Expression) multDivide(symbols map[string]interface{}) (interface{}, error) {

	v1, err := e.expressionAtom(symbols)
	if err != nil {
		return nil, err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if inList(op, []string{"*", "/", "|"}) {
			e.TokenP = e.TokenP + 1

			v2, err := e.expressionAtom(symbols)
			if err != nil {
				return nil, err
			}

			v1, v2 = Normalize(v1, v2)
			switch op {

			case "*":
				switch v1.(type) {
				case int:
					v1 = v1.(int) * v2.(int)
				case float64:
					v1 = v1.(float64) * v2.(float64)
				case bool:
					v1 = v1.(bool) || v2.(bool)
				default:
					return nil, errors.New("Invalid operand types for *")
				}

			case "/":
				switch v1.(type) {
				case int:
					v1 = v1.(int) / v2.(int)
				case float64:
					v1 = v1.(float64) / v2.(float64)
				default:
					return nil, errors.New("invalid type for '/' operator")
				}

			case "|":
				v1 = Coerce(v1, true)
				v1 = Coerce(v2, true)
				v1 = v1.(bool) || v2.(bool)
			}

		} else {
			parsing = false
		}
	}
	return v1, nil
}
