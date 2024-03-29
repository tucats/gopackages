package bytecode

import (
	"github.com/tucats/gopackages/errors"
)

// error is a helper function used to generate a new error based
// on the runtime context. The current module name and line number
// from the context are stored in the new error object, along with
// the message and context.
func (c *Context) error(err error, context ...interface{}) *errors.Error {
	if err == nil {
		return nil
	}

	var r *errors.Error

	if e, ok := err.(*errors.Error); ok {
		if !e.HasIn() {
			e = e.In(c.name)
		}

		if !e.HasAt() {
			e = e.At(c.GetLine(), 0)
		}

		r = e
	} else {
		r = errors.NewError(err).In(c.name).At(c.GetLine(), 0)
	}

	if len(context) > 0 {
		r = r.Context(context[0])
	}

	return r
}
