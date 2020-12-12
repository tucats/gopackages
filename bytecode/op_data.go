package bytecode

import (
	"errors"

	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*           D A T A  T Y P E S            *
*         A N D   S T O R A G E           *
*                                         *
\******************************************/

// MakeArrayOpcode implementation
func MakeArrayOpcode(c *Context, i interface{}) error {

	parms := util.GetInt(i)

	if parms == 2 {
		initialValue, err := c.Pop()
		if err != nil {
			return err
		}
		sv, err := c.Pop()
		if err != nil {
			return err
		}
		size := util.GetInt(sv)
		if size < 0 {
			size = 0
		}
		array := make([]interface{}, size)
		for n := 0; n < size; n++ {
			array[n] = initialValue
		}
		_ = c.Push(array)
		return nil
	}

	// No initializer, so get the size and make it
	// a non-negative integer
	sv, err := c.Pop()
	if err != nil {
		return err
	}

	size := util.GetInt(sv)
	if size < 0 {
		size = 0
	}
	array := make([]interface{}, size)
	_ = c.Push(array)

	return nil
}

// ArrayOpcode implementation
func ArrayOpcode(c *Context, i interface{}) error {

	count := util.GetInt(i)
	array := make([]interface{}, count)

	for n := 0; n < count; n++ {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		array[(count-n)-1] = v
	}

	_ = c.Push(array)
	return nil
}

// StructOpcode implementation. The operand is a count
// of elements on the stack. These are pulled off in pairs,
// where the first value is the name of the struct field and
// the second value is the value of the struct field.
func StructOpcode(c *Context, i interface{}) error {

	count := util.GetInt(i)

	m := map[string]interface{}{}

	for n := 0; n < count; n++ {
		name, err := c.Pop()
		if err != nil {
			return err
		}
		value, err := c.Pop()
		if err != nil {
			return err
		}
		m[util.GetString(name)] = value
	}

	_ = c.Push(m)
	return nil
}

// CoerceOpcode implementation
func CoerceOpcode(c *Context, i interface{}) error {
	t := util.GetInt(i)
	v, err := c.Pop()
	if err != nil {
		return err
	}
	switch t {
	case IntType:
		v = util.GetInt(v)
	case FloatType:
		v = util.GetFloat(v)
	case StringType:
		v = util.GetString(v)
	case BoolType:
		v = util.GetBool(v)
	case ArrayType:
		// If it's  not already an array, wrap it in one.
		if _, ok := v.([]interface{}); !ok {
			v = []interface{}{v}
		}
	case StructType:
		// If it's not a struct, we can't do anything so fail
		if _, ok := v.(map[string]interface{}); !ok {
			return c.NewError(InvalidTypeError)
		}
	case UndefinedType:
		// No work at all to do here.

	default:
		return c.NewError(InvalidTypeError)
	}

	_ = c.Push(v)
	return nil
}

/******************************************\
*                                         *
*         D A T A   A C C E S S           *
*                                         *
\******************************************/

// StoreOpcode implementation
func StoreOpcode(c *Context, i interface{}) error {

	v, err := c.Pop()
	if err != nil {
		return err
	}

	// Get the name. If it is the reserved name "_" it means
	// to just discard the value.
	varname := util.GetString(i)
	if varname == "_" {
		return nil
	}
	err = c.Set(varname, v)
	if err != nil {
		return c.NewError(err.Error())
	}

	// Is this a readonly variable that is a structure? If so, mark it
	// with the embedded readonly flag.

	if len(varname) > 1 && varname[0:1] == "_" {
		switch a := v.(type) {
		case map[string]interface{}:
			a["__readonly"] = true
		}
	}
	return err
}

// LoadOpcode implementation
func LoadOpcode(c *Context, i interface{}) error {

	name := util.GetString(i)
	if len(name) == 0 {
		return c.NewError(InvalidIdentifierError, name)
	}
	v, found := c.Get(util.GetString(i))
	if !found {
		return c.NewError(UnknownIdentifierError, name)
	}

	_ = c.Push(v)
	return nil
}

// MemberOpcode implementation. This pops two values from
// the stack (the first must be a string and the second a
// map) and indexes into the map to get the matching value
// and puts back on the stack.
func MemberOpcode(c *Context, i interface{}) error {

	var name string
	if i != nil {
		name = util.GetString(i)
	} else {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		name = util.GetString(v)
	}

	m, err := c.Pop()
	if err != nil {
		return err
	}

	// The only the type that is supported is a map
	var v interface{}
	found := false

	mv, ok := m.(map[string]interface{})
	if ok {
		v, found = mv[name]
		if !found {
			return c.NewError(UnknownMemberError, name)
		}
		c.this = m // Remember where we loaded this from
	} else {
		return c.NewError(InvalidTypeError)
	}
	_ = c.Push(v)
	return nil
}

// ClassMemberOpcode implementation. This pops two values from
// the stack (the first must be a string and the second a
// map) and indexes into the map to get the matching value
// and puts back on the stack.
//
// If the member does not exist, but there is a __parent
// member in the structure, we also search the __parent field
// for the value. This supports calling packages based on
// a given object value.
func ClassMemberOpcode(c *Context, i interface{}) error {

	var name string
	if i != nil {
		name = util.GetString(i)
	} else {
		v, err := c.Pop()
		if err != nil {
			return err
		}
		name = util.GetString(v)
	}

	m, err := c.Pop()
	if err != nil {
		return err
	}

	// The only the type that is supported is a map
	switch mv := m.(type) {
	case map[string]interface{}:

		if _, found := mv["__parent"]; !found {
			return c.NewError(NotATypeError)
		}
		v, found := mv[name]
		if !found {

			v, found := searchParents(mv, name)
			if found {
				return c.Push(v)
			}
			return c.NewError(UnknownMemberError, name)
		}
		_ = c.Push(v)

	default:
		return c.NewError(InvalidTypeError)
	}
	return nil
}

func searchParents(mv map[string]interface{}, name string) (interface{}, bool) {

	// Is there a parent we should check?
	if t, found := mv["__parent"]; found {
		switch tv := t.(type) {
		case map[string]interface{}:
			v, found := tv[name]
			if !found {
				return searchParents(tv, name)
			}
			return v, true

		case string:
			return nil, false

		default:
			return nil, false
		}
	}
	return nil, false
}

// LoadIndexOpcode implementation
func LoadIndexOpcode(c *Context, i interface{}) error {

	index, err := c.Pop()
	if err != nil {
		return err
	}

	array, err := c.Pop()
	if err != nil {
		return err
	}

	switch a := array.(type) {

	// Index into map is just member access
	case map[string]interface{}:
		subscript := util.GetString(index)
		v, f := a[subscript]
		if !f {
			return c.NewError(UnknownMemberError, subscript)
		}
		_ = c.Push(v)
		c.this = a

	case []interface{}:
		subscript := util.GetInt(index)
		if subscript < 0 || subscript >= len(a) {
			return c.NewError(InvalidArrayIndexError, subscript)
		}
		v := a[subscript]
		_ = c.Push(v)

	default:
		return c.NewError(InvalidTypeError)
	}

	return nil
}

// LoadSliceOpcode implementation
func LoadSliceOpcode(c *Context, i interface{}) error {

	index2, err := c.Pop()
	if err != nil {
		return err
	}

	index1, err := c.Pop()
	if err != nil {
		return err
	}

	array, err := c.Pop()
	if err != nil {
		return err
	}

	switch a := array.(type) {

	// Array of objects means we retrieve a slice.
	case []interface{}:
		subscript1 := util.GetInt(index1)
		if subscript1 < 0 || subscript1 >= len(a) {
			return c.NewError(InvalidSliceIndexError, subscript1)
		}
		subscript2 := util.GetInt(index2)
		if subscript2 < subscript1 || subscript2 >= len(a) {
			return c.NewError(InvalidSliceIndexError, subscript2)
		}
		v := a[subscript1 : subscript2+1]
		_ = c.Push(v)

	default:
		return c.NewError(InvalidTypeError)
	}

	return nil
}

// StoreIndexOpcode implementation
func StoreIndexOpcode(c *Context, i interface{}) error {

	index, err := c.Pop()
	if err != nil {
		return err
	}

	array, err := c.Pop()
	if err != nil {
		return err
	}

	v, err := c.Pop()
	if err != nil {
		return err
	}

	switch a := array.(type) {

	// Index into map is just member access. Make sure it's not
	// a read-only member or a function pointer...
	case map[string]interface{}:
		subscript := util.GetString(index)

		// Does this member have a flag marking it as readonly?
		old, found := a["__readonly"]
		if found {
			if util.GetBool(old) {
				return c.NewError(ReadOnlyError)
			}
		}
		// Does this item already exist and is readonly?
		old, found = a[subscript]
		if found {
			if subscript[0:1] == "_" {
				return c.NewError(ReadOnlyError)
			}

			// Check to be sure this isn't a restricted (function code) type

			switch old.(type) {

			case func(*symbols.SymbolTable, []interface{}) (interface{}, error):
				return errors.New(ReadOnlyError)

			}
		}

		a[subscript] = v
		// If we got a true argument, push the result back on the stack also. This
		// is needed to create TYPE definitions.
		if util.GetBool(i) {
			_ = c.Push(a)
		}

	// Index into array is integer index
	case []interface{}:
		subscript := util.GetInt(index)
		if subscript < 0 || subscript >= len(a) {
			return c.NewError(InvalidArrayIndexError, subscript)
		}
		a[subscript] = v
		_ = c.Push(a)

	default:
		return c.NewError(InvalidTypeError)
	}

	return nil
}

// ThisOpcode implements the This opcode
func ThisOpcode(c *Context, i interface{}) error {
	c.this = util.GetString(i)
	v, err := c.Pop()
	if err != nil {
		return err
	}
	if this, ok := c.this.(string); ok {
		return c.SetAlways(this, v)
	}
	return c.NewError(InvalidThisError)
}
