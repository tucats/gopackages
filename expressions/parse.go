package expressions

import "github.com/tucats/gopackages/tokenizer"

// Parse parses a text expression
func (e *Expression) Parse(s string) error {

	e.t = tokenizer.New(s)

	return nil
}
