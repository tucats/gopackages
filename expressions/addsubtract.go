package expressions

import (
	"errors"
	"reflect"
)

func (e *Expression) addSubtract(symbols map[string]interface{}) (interface{}, error) {

	v1, err := e.multDivide(symbols)
	if err != nil {
		return nil, err
	}

	var parsing = true
	for parsing {
		if e.TokenP >= len(e.Tokens) {
			break
		}
		op := e.Tokens[e.TokenP]
		if inList(op, []string{"+", "-", "&"}) {
			e.TokenP = e.TokenP + 1

			v2, err := e.multDivide(symbols)
			if err != nil {
				return nil, err
			}

			// Let's handle special case of an array, which just
			// appends the item to the array.
			switch a := v1.(type) {
			case []interface{}:

				switch op {
				case "+":

					switch element := v2.(type) {

					// If second item also an array, append elements
					case []interface{}:
						v1 = append(a, element...)

					// Else append the opaque object.
					default:
						v1 = append(a, v2)
					}

				case "-":
					vNew := make([]interface{}, 0)
					for _, t := range a {
						if !reflect.DeepEqual(t, v2) {
							vNew = append(vNew, t)
						}
					}
					v1 = vNew
				default:
					return nil, errors.New("Unsupported operation on array")
				}
				return v1, nil
			}

			// Otherwise, normalize the two items and go...
			v1, v2 = Normalize(v1, v2)
			switch op {

			case "+":
				switch v1.(type) {
				case int:
					v1 = v1.(int) + v2.(int)
				case string:
					v1 = v1.(string) + v2.(string)
				case float64:
					v1 = v1.(float64) + v2.(float64)
				case bool:
					v1 = v1.(bool) && v2.(bool)
				}

			case "-":
				switch v1.(type) {
				case int:
					v1 = v1.(int) - v2.(int)
				case float64:
					v1 = v1.(float64) - v2.(float64)
				default:
					return nil, errors.New("invlid type for '-' operator")
				}

			case "&":
				v1 = Coerce(v1, true)
				v2 = Coerce(v2, true)
				if v1 == nil || v2 == nil {
					return nil, errors.New("invalid value for coercion to bool")
				}
				v1 = v1.(bool) && v2.(bool)
			}

		} else {
			parsing = false
		}
	}
	return v1, nil
}
