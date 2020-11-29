package expressions

import (
	bc "github.com/tucats/gopackages/bytecode"
)

func (e *Expression) functionCall() error {

	// Note, caller already consumed the opening paren
	argc := 0

	for e.t.Peek(1) != ")" {
		err := e.conditional()
		if err != nil {
			return err
		}
		argc = argc + 1
		if e.t.AtEnd() {
			break
		}
		if e.t.Peek(1) == ")" {
			break
		}
		if e.t.Peek(1) != "," {
			return e.NewError(InvalidListError)
		}
		e.t.Advance(1)
	}

	// Ensure trailing parenthesis
	if e.t.AtEnd() || e.t.Peek(1) != ")" {
		return e.NewError(MissingParenthesisError)
	}
	e.t.Advance(1)

	// Call the function
	e.b.Emit(bc.Call, argc)
	return nil
}

// Function compiles a function call. The value of the
// function has been pushed to the top of the stack.
func (e *Expression) Function() error {

	// Get the atom
	err := e.reference()
	if err != nil {
		return err
	}
	// Peek ahead to see if it's the start of a function call...
	if e.t.IsNext("(") {
		return e.functionCall()
	}
	return nil
}
