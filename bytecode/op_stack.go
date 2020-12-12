package bytecode

import (
	"encoding/json"

	"github.com/tucats/gopackages/util"
)

/******************************************\
*                                         *
*    S T A C K   M A N A G E M E N T      *
*                                         *
\******************************************/

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
