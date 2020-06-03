package bytecode

import "github.com/tucats/gopackages/util"

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
