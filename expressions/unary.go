package expressions

import "errors"

func (e *Expression) unary(symbols map[string]interface{}) (interface{}, error) {

	// Check for unary negation or not before passing into top-level diadic operators.

	for e.TokenP < len(e.Tokens) {

		t := e.Tokens[e.TokenP]
		switch t {
		case "-":
			e.TokenP = e.TokenP + 1
			v, err := e.expressionAtom(symbols)
			if err != nil {
				return nil, err
			}
			switch value := v.(type) {
			case bool:
				return !value, nil

			case int:
				return -value, nil
			case float64:
				return 0.0 - value, nil

			case string:
				return nil, errors.New("invalid data type for negation")
			}

		case "!":
			e.TokenP = e.TokenP + 1
			v, err := e.expressionAtom(symbols)
			if err != nil {
				return nil, err
			}

			return !(Coerce(v, true).(bool)), nil

		default:
			return e.expressionAtom(symbols)

		}
	}
	return e.expressionAtom(symbols)
}
