package compiler

import (
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

const (
	indexLoopType = 1
	rangeLoopType = 2
)

// Loop is a structure that defines a loop type.
type Loop struct {
	Parent *Loop
	Type   int
	// Fixup locations for break or continue statements in a
	// loop. These are the addresses that must be fixed up with
	// the target address.
	breaks    []int
	continues []int
}

// FunctionDictionary is a list of functions and the bytecode or native function pointer
type FunctionDictionary map[string]interface{}

// PackageDictionary is a list of packages each with a FunctionDictionary
type PackageDictionary map[string]FunctionDictionary

// Compiler is a structure defining what we know about the
// compilation
type Compiler struct {
	PackageName    string
	b              *bytecode.ByteCode
	t              *tokenizer.Tokenizer
	s              *symbols.SymbolTable
	loops          *Loop
	coerce         *bytecode.ByteCode
	constants      []string
	packages       PackageDictionary
	blockDepth     int
	statementCount int
}

// New creates a new compiler instance
func New() *Compiler {
	cInstance := Compiler{
		b:         nil,
		t:         nil,
		s:         &symbols.SymbolTable{Name: "compile-unit"},
		constants: make([]string, 0),
		packages:  PackageDictionary{},
	}
	c := &cInstance
	return c
}

// CompileString turns a string into a compilation unit. This is a helper function
// around the Compile() operation that removes the need for the caller
// to provide a tokenizer.
func (c *Compiler) CompileString(source string) (*bytecode.ByteCode, error) {
	t := tokenizer.New(source)
	return c.Compile(t)
}

// Compile starts a compilation unit, and returns a bytecode
// of the compiled material.
func (c *Compiler) Compile(t *tokenizer.Tokenizer) (*bytecode.ByteCode, error) {

	c.b = bytecode.New("")
	c.t = t

	c.t.Reset()

	for !c.t.AtEnd() {
		err := c.Statement()
		if err != nil {
			return nil, err
		}
	}

	// Merge in any package definitions
	c.AddPackageToSymbols(c.b.Symbols)

	// Also merge in any other symbols created for this function
	c.b.Symbols.Merge(c.Symbols())

	return c.b, nil
}

// AddBuiltins adds the builtins for the named package (or prebuilt builtins if the package name
// is empty)
func (c *Compiler) AddBuiltins(pkgname string) bool {

	added := false
	for name, f := range functions.FunctionDictionary {

		if dot := strings.Index(name, "."); dot >= 0 {
			f.Pkg = name[:dot]
			name = name[dot+1:]
			ui.Debug("=== n=%s, p=%s\n", name, f.Pkg)
		}

		if f.Pkg == pkgname {
			c.AddPackageFunction(pkgname, name, f.F)
			added = true
		}
	}
	return added
}

// AddPackageFunction adds a new package function to the compiler's package dictionary. If the
// package name does not yet exist, it is created. The function name and interface are then used
// to add an entry for that package.
func (c *Compiler) AddPackageFunction(pkgname string, name string, function interface{}) error {

	fd, found := c.packages[pkgname]
	if !found {
		fd = FunctionDictionary{}
	}

	if _, found := fd[name]; found {
		return c.NewError("Duplicate function definition")
	}
	fd[name] = function
	c.packages[pkgname] = fd

	return nil
}

// AddPackageToSymbols adds all the defined packages for this compilation to the named symbol table.
func (c *Compiler) AddPackageToSymbols(s *symbols.SymbolTable) {

	for pkgname, dict := range c.packages {

		m := map[string]interface{}{}
		for k, v := range dict {

			// If the package name is empty, we add the individual items
			if pkgname == "" {
				s.SetAlways(k, v)
			} else {
				// Otherwise, copy the entire map
				m[k] = v
			}
		}
		if pkgname != "" {
			s.SetAlways(pkgname, m)
		}
	}
}

// StatementEnd returns true when the next token is
// the end-of-statement boundary
func (c *Compiler) StatementEnd() bool {
	next := c.t.Peek(1)

	if next == tokenizer.EndOfTokens {
		return true
	}

	return (next == ";") || (next == "}")
}

// Symbols returns the symbol table map from compilation
func (c *Compiler) Symbols() *symbols.SymbolTable {
	return c.s
}
