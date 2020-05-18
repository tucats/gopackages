package bytecode

import (
	"github.com/tucats/gopackages/symbols"
	sym "github.com/tucats/gopackages/symbols"
)

// Context holds the runtime information about an instance of bytecode being
// executed.
type Context struct {
	Name    string
	bc      *ByteCode
	pc      int
	stack   []interface{}
	sp      int
	running bool
	line    int
	symbols *sym.SymbolTable
	Tracing bool
}

// NewContext generates a new context. It must be passed a symbol table and a bytecode
// array. A context holds the runtime state of a given execution unit (program counter,
// runtime stack, symbol table) and is used to actually run bytecode. The bytecode
// can continue to be modified after it is associated with a context.
// @TOMCOLE Is this a good idea? Should a context take a snapshot of the bytecode at
// the time so it is immutable?
func NewContext(s *symbols.SymbolTable, b *ByteCode) *Context {
	ctx := Context{
		Name:    b.Name,
		bc:      b,
		pc:      0,
		stack:   make([]interface{}, InitialStackSize),
		sp:      0,
		running: false,
		line:    0,
		symbols: s,
		Tracing: false,
	}
	ctxp := &ctx
	ctxp.SetByteCode(b)

	// Append the bytecode symbols into the symbol table.
	for k, v := range b.Symbols.Symbols {
		s.Set(k, v)
	}
	return ctxp
}

// AppendSymbols appends a symbol table to the current
// context. This is used to add in compiler maps, for
// example.
func (c *Context) AppendSymbols(s symbols.SymbolTable) {
	for k, v := range s.Symbols {
		c.symbols.Set(k, v)
	}
}

// SetByteCode attaches a new bytecode object to the current run context.
func (c *Context) SetByteCode(b *ByteCode) {
	c.bc = b
}

// Get is a helper function that retrieves a symbol value from the associated
// symbol table
func (c *Context) Get(name string) (interface{}, bool) {

	v, found := c.symbols.Get(name)
	return v, found
}

// Set is a helper function that sets a symbol value in the associated
// symbol table
func (c *Context) Set(name string, value interface{}) {
	c.symbols.Set(name, value)
}

// Pop removes the top-most item from the stack
func (c *Context) Pop() (interface{}, error) {
	if c.sp <= 0 || len(c.stack) < c.sp {
		return nil, c.NewError("stack underflow")
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
