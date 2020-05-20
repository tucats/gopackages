package symbols

import (
	"errors"
)

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
		"_author":    "Tom Cole",
		"_copyright": "(c) Copyright 2020",
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

// SetAlways stores a symbol value in the local table. No value in
// any parent table is affected. This can be used for functions and
// readonly values.
func (s *SymbolTable) SetAlways(name string, v interface{}) error {
	if s.Symbols == nil {
		s.Symbols = map[string]interface{}{}
	}

	s.Symbols[name] = v
	return nil
}

// Set stores a symbol value in the table where it was found.
func (s *SymbolTable) Set(name string, v interface{}) error {
	if s.Symbols == nil {
		s.Symbols = map[string]interface{}{}
	}

	old, found := s.Symbols[name]

	if found {
		if name[0:1] == "_" {
			return errors.New("readonly symbol")
		}

		// Check to be sure this isn't a restricted (function code) type

		switch old.(type) {

		case func([]interface{}) (interface{}, error):
			return errors.New("readonly builtin symbol")

		}
	} else {

		// If there are no more tables, we have an error.
		if s.Parent == nil {
			return errors.New("unknown symbol")
		}
		// Otherwise, ask the parent to try to set the value.
		return s.Parent.Set(name, v)
	}

	s.Symbols[name] = v
	return nil
}

// Delete removes a symbol from the table. Search from the local symbol
// up the parent tree until you find the symbol to delete.
func (s *SymbolTable) Delete(name string) error {

	if len(name) == 0 {
		return errors.New("invalid symbol")
	}
	if name[:1] == "_" {
		return errors.New("readonly symbol")
	}
	if s.Symbols == nil {
		return errors.New("SymbolDelete of " + name + " when there are no symbols")
	}

	_, f := s.Symbols[name]
	if !f {
		if s.Parent == nil {
			return errors.New("symbol " + name + " not found")
		}
		return s.Parent.Delete(name)
	}
	delete(s.Symbols, name)
	return nil
}

// Create creates a symbol name in the table
func (s *SymbolTable) Create(name string) error {

	if len(name) == 0 {
		return errors.New("invalid symbol")
	}

	_, found := s.Symbols[name]
	if found {
		return errors.New("symbol already exists")
	}
	s.Symbols[name] = nil
	return nil
}
