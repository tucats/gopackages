package bytecode

import (
	"encoding/json"
	"fmt"

	"github.com/tucats/gopackages/util"
)

type StackMarker struct {
	Desc string
}

func NewStackMarker(label string, count int) StackMarker {
	return StackMarker{
		Desc: fmt.Sprintf("%s %d items", label, count),
	}
}

/******************************************\
*                                         *
*    S T A C K   M A N A G E M E N T      *
*                                         *
\******************************************/

// DropToMarkerOpcode discards items on the stack until it
// finds a marker value, at which point it stops. This is
// used to discard unused return values on the stack. IF there
// is no marker, this drains the stack.
func DropToMarkerOpcode(c *Context, i interface{}) error {
	found := false
	for !found {
		v, err := c.Pop()
		if err != nil {
			break
		}
		_, found = v.(StackMarker)
	}
	return nil
}

// StackCheckOpcode has an integer argument, and verifies
// that there are this many items on the stack, which is
// used to verify that multiple return-values on the stack
// are present.
func StackCheckOpcode(c *Context, i interface{}) error {
	count := util.GetInt(i)
	if c.sp <= count {
		return c.NewError(IncorrectReturnValueCount)
	}

	// The marker is an instance of a StackMarker object.
	v := c.stack[c.sp-(count+1)]
	if _, ok := v.(StackMarker); ok {
		return nil
	}
	return c.NewError(IncorrectReturnValueCount)
}

// PushOpcode bytecode implementation
func PushOpcode(c *Context, i interface{}) error {
	return c.Push(i)
}

// DropOpcode implementation
func DropOpcode(c *Context, i interface{}) error {

	count := 1
	if i != nil {
		count = util.GetInt(i)
	}
	for n := 0; n < count; n = n + 1 {
		_, err := c.Pop()
		if err != nil {
			return nil
		}
	}
	return nil
}

// DupOpcode implementation
func DupOpcode(c *Context, i interface{}) error {
	v, err := c.Pop()
	if err != nil {
		return err
	}
	_ = c.Push(v)
	_ = c.Push(v)
	return nil
}

// SwapOpcode implementation
func SwapOpcode(c *Context, i interface{}) error {
	v1, err := c.Pop()
	if err != nil {
		return err
	}
	v2, err := c.Pop()
	if err != nil {
		return err
	}
	_ = c.Push(v1)
	_ = c.Push(v2)
	return nil
}

// CopyOpcode implementation
func CopyOpcode(c *Context, i interface{}) error {
	v, err := c.Pop()
	if err != nil {
		return err
	}
	_ = c.Push(v)

	// Use JSON as a reflection-based cloner
	var v2 interface{}
	byt, _ := json.Marshal(v)
	err = json.Unmarshal(byt, &v2)

	_ = c.Push(2)
	return err
}
