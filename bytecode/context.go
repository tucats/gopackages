package bytecode

import (
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
	sym "github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// Context holds the runtime information about an instance of bytecode being
// executed.
type Context struct {
	Name      string
	bc        *ByteCode
	pc        int
	stack     []interface{}
	sp        int
	running   bool
	line      int
	symbols   *sym.SymbolTable
	Tracing   bool
	tokenizer *tokenizer.Tokenizer
	try       []int
	output    *strings.Builder
	this      interface{}
	result    interface{}
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
		this:    "",
		try:     make([]int, 0),
	}
	ctxp := &ctx
	ctxp.SetByteCode(b)

	// If we weren't given a table, create an empty temp table.
	if s == nil {
		s = symbols.NewSymbolTable("")
	}

	// Append any bytecode symbols into the symbol table.
	if b.Symbols != nil {
		for k, v := range b.Symbols.Symbols {
			_ = s.SetAlways(k, v)
		}
	}

	return ctxp
}

// EnableConsoleOutput tells the context to begin capturing all output normally generated
// from Print and Newline into a buffer instead of going to stdout
func (c *Context) EnableConsoleOutput(flag bool) {

	ui.Debug(ui.AppLogger, ">>> Console output set to %v", flag)
	if !flag {
		var b strings.Builder
		c.output = &b
	} else {
		c.output = nil
	}
}

// GetOutput retrieves the output buffer
func (c *Context) GetOutput() string {
	if c.output != nil {
		return c.output.String()
	}
	return ""
}

// SetTokenizer sets a tokenizer in the current context for tracing and debugging.
func (c *Context) SetTokenizer(t *tokenizer.Tokenizer) {
	c.tokenizer = t
}

// GetTokenizer gets the tokenizer in the current context for tracing and debugging.
func (c *Context) GetTokenizer() *tokenizer.Tokenizer {
	return c.tokenizer
}

// AppendSymbols appends a symbol table to the current
// context. This is used to add in compiler maps, for
// example.
func (c *Context) AppendSymbols(s symbols.SymbolTable) {
	for k, v := range s.Symbols {
		_ = c.symbols.SetAlways(k, v)
	}
}

// SetByteCode attaches a new bytecode object to the current run context.
func (c *Context) SetByteCode(b *ByteCode) {
	c.bc = b
}

// SetConstant is a helper function to define a constant value
func (c *Context) SetConstant(name string, v interface{}) error {
	return c.symbols.SetConstant(name, v)
}

// IsConstant is a helper function to define a constant value
func (c *Context) IsConstant(name string) bool {
	return c.symbols.IsConstant(name)
}

// Get is a helper function that retrieves a symbol value from the associated
// symbol table
func (c *Context) Get(name string) (interface{}, bool) {

	v, found := c.symbols.Get(name)
	return v, found
}

// Set is a helper function that sets a symbol value in the associated
// symbol table
func (c *Context) Set(name string, value interface{}) error {
	return c.symbols.Set(name, value)
}

// SetAlways is a helper function that sets a symbol value in the associated
// symbol table
func (c *Context) SetAlways(name string, value interface{}) error {
	return c.symbols.SetAlways(name, value)
}

// Delete deletes a symbol from the current context
func (c *Context) Delete(name string) error {
	return c.symbols.Delete(name)
}

// Create creates a symbol
func (c *Context) Create(name string) error {
	return c.symbols.Create(name)
}

// Pop removes the top-most item from the stack
func (c *Context) Pop() (interface{}, error) {
	if c.sp <= 0 || len(c.stack) < c.sp {
		return nil, c.NewError(StackUnderflowError)
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

// FormatStack formats the stack for tracing output
func FormatStack(s []interface{}, newlines bool) string {

	if len(s) == 0 {
		return "<empty>"
	}
	var b strings.Builder
	if newlines {
		b.WriteString("\n")
	}

	for n := len(s) - 1; n >= 0; n = n - 1 {

		if n < len(s)-1 {
			b.WriteString(", ")
			if newlines {
				b.WriteString("\n")
			}
		}

		b.WriteString(util.Format(s[n]))
		if !newlines && b.Len() > 50 {
			return b.String()[:50] + "..."
		}
	}
	return b.String()
}

// GetConfig retrieves a runtime configuration item from the
// __config data structure. If not found, it also queries
// the persistence layer.
func (c *Context) GetConfig(name string) interface{} {
	var i interface{}

	if config, ok := c.Get("_config"); ok {
		if cfgMap, ok := config.(map[string]interface{}); ok {
			if cfgValue, ok := cfgMap[name]; ok {
				i = cfgValue
			}
		}
	}
	return i
}
