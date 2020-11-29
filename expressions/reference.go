package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

// reference parses a structure or array reference
func (e *Expression) reference() error {

	// Parse the function call or exprssion atom
	err := e.expressionAtom()
	if err != nil {
		return err
	}

	// is there a trailing structure or array reference?
	for !e.t.AtEnd() {

		op := e.t.Peek(1)
		switch op {

		// Struct reference
		case "->":
			e.t.Advance(1)
			name := e.t.Next()
			e.b.Emit(bc.Dup)
			e.b.Emit(bc.Push, name)
			e.b.Emit(bc.ClassMember)

		// Map member reference
		case ".":
			e.t.Advance(1)
			name := e.t.Next()
			e.b.Emit(bc.Push, name)
			e.b.Emit(bc.Member)

		// Array index reference
		case "[":
			e.t.Advance(1)
			err := e.conditional()
			if err != nil {
				return err
			}

			// is it a slice instead of an index?
			if e.t.IsNext(":") {
				err := e.conditional()
				if err != nil {
					return err
				}
				e.b.Emit(bc.LoadSlice)
				if e.t.Next() != "]" {
					return e.NewError(MissingBracketError)
				}
			} else {
				// Nope, singular index
				if e.t.Next() != "]" {
					return e.NewError(MissingBracketError)
				}
				e.b.Emit(bc.LoadIndex)
			}

		// Nothing else, term is complete
		default:
			return nil
		}
	}
	return nil
}
