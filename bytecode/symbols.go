package bytecode

// SymbolTable contains an abstract symbol table
type SymbolTable struct {
	Name    string
	Parent  *SymbolTable
	Symbols map[string]interface{}
}

// RootSymbolTable is the parent of all other tables.
var RootSymbolTable = SymbolTable{
	Name:   "Root Symbol Table",
	Parent: nil,
	Symbols: map[string]interface{}{
		"author":    "Tom Cole",
		"copyright": "(c) Copyright 2020",
	},
}

// NewSymbolTable generates a new symbol table
func NewSymbolTable(name string) *SymbolTable {

	symbols := SymbolTable{
		Name:    name,
		Parent:  &RootSymbolTable,
		Symbols: map[string]interface{}{},
	}
	return &symbols
}

// NewChildSymbolTable generates a new symbol table with an assigned
// parent table.
func NewChildSymbolTable(name string, parent *SymbolTable) *SymbolTable {

	symbols := SymbolTable{
		Name:    name,
		Parent:  parent,
		Symbols: map[string]interface{}{},
	}
	return &symbols
}

// Get retrieves a symbol from the current table or any parent
// table that exists
func (s *SymbolTable) Get(name string) (interface{}, bool) {

	v, f := s.Symbols[name]
	if !f && s.Parent != nil {
		return s.Parent.Get(name)
	}
	return v, f
}

// Set stores a symbol value in the local table. No value in
// any parent table is affected.
func (s *SymbolTable) Set(name string, v interface{}) {
	s.Symbols[name] = v
}