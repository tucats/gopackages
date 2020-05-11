package expressions

import (
	"errors"

	"github.com/tucats/gopackages/util"
)

// conditional handles parsing the ?: trinary operator. The first term is
// converted to a boolean value, and if true the second term is returned, else
// the third term. All terms must be present.
func (e *Expression) conditional(symbols map[string]interface{}) (interface{}, error) {

	// Parse the conditional
	v, err := e.relations(symbols)
	if err != nil {
		return nil, err
	}

	// If this is not a conditional, we're done.

	if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != "?" {
		return v, nil
	}

	// Parse both parts of the alternate values
	e.TokenP = e.TokenP + 1
	v1, err := e.relations(symbols)
	if e.TokenP >= len(e.Tokens) || e.Tokens[e.TokenP] != ":" {
		return nil, errors.New("missing colon in conditional")
	}
	e.TokenP = e.TokenP + 1
	v2, err := e.relations(symbols)

	if util.GetBool(v) {
		return v1, nil
	}
	return v2, nil

}
