package bytecode

import "errors"

// SymbolTable contains an abstract symbol table
type SymbolTable map[string]interface{}

// Context holds the runtime information about an instance of bytecode being
// executed.
type Context struct {
	bc      *ByteCode
	pc      int
	stack   []interface{}
	sp      int
	running bool
	symbols SymbolTable
}

// NewSymbolTable creates a new symbol table object
func NewSymbolTable() SymbolTable {
	return map[string]interface{}{}
}

// NewContext generates a new context. It must be passed a symbol table and a bytecode
// array. A context holds the runtime state of a given execution unit (program counter,
// runtime stack, symbol table) and is used to actually run bytecode. The bytecode
// can continue to be modified after it is associated with a context.
// @TOMCOLE Is this a good idea? Should a context take a snapshot of the bytecode at
// the time so it is immutable?
func NewContext(s SymbolTable, b *ByteCode) *Context {
	ctx := Context{
		bc:      b,
		pc:      0,
		stack:   make([]interface{}, InitialStackSize),
		sp:      0,
		running: false,
		symbols: s,
	}
	ctxp := &ctx
	ctxp.SetByteCode(b)

	return ctxp
}

// SetByteCode attaches a new bytecode object to the current run context.
func (c *Context) SetByteCode(b *ByteCode) {
	c.bc = b
}

// Get retrieves a symbol value from the symbol table
func (c *Context) Get(name string) (interface{}, bool) {

	v, found := c.symbols[name]
	return v, found
}

// Set sets a symbol value in the symbol table
func (c *Context) Set(name string, value interface{}) {
	c.symbols[name] = value
}

// Pop removes the top-most item from the stack
func (c *Context) Pop() (interface{}, error) {
	if c.sp <= 0 || len(c.stack) < c.sp {
		return nil, errors.New("stack underflow")
	}

	c.sp = c.sp - 1
	v := c.stack[c.sp]
	return v, nil
}

// Push puts a new items on the stack
func (c *Context) Push(v interface{}) error {

	if c.sp >= len(c.stack) {
		c.stack = append(c.stack, make([]interface{}, GrowStackBy)...)
	}
	c.stack[c.sp] = v
	c.sp = c.sp + 1
	return nil
}
