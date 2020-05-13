package expressions

import (
	"errors"

	bc "github.com/tucats/gopackages/bytecode"
)

// conditional handles parsing the ?: trinary operator. The first term is
// converted to a boolean value, and if true the second term is returned, else
// the third term. All terms must be present.
func (e *Expression) conditional() error {

	// Parse the conditional
	err := e.relations()
	if err != nil {
		return err
	}

	// If this is not a conditional, we're done.

	if e.t.AtEnd() || e.t.Peek() != "?" {
		return nil
	}

	m1 := e.b.Mark()
	e.b.Emit(bc.BranchFalse, 0)

	// Parse both parts of the alternate values
	e.t.Advance(1)
	err = e.relations()
	if err != nil {
		return err
	}
	if e.t.AtEnd() || e.t.Peek() != ":" {
		return errors.New("missing colon in conditional")
	}
	m2 := e.b.Mark()
	e.b.Emit(bc.Branch, 0)

	e.b.SetAddressHere(m1)
	e.t.Advance(1)
	err = e.relations()
	if err != nil {
		return err
	}

	// Patch up the forward references.
	e.b.SetAddressHere(m2)

	return nil

}
