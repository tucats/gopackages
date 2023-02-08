package expressions

import "github.com/tucats/gopackages/expressions/tokenizer"

// Parse parses a text expression
func (e *Expression) Parse(s string) error {

	e.t = tokenizer.New(s, true)

	return nil
}
