package compiler

import (
	"strings"

	"github.com/tucats/gopackages/app-cli/settings"
	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/expressions/builtins"
	"github.com/tucats/gopackages/expressions/bytecode"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
	"github.com/tucats/gopackages/expressions/tokenizer"
)

// flagSet contains flags that generally identify the state of
// the compiler at any given moment. For example, when parsing
// something like a switch conditional value, the value cannot
// be a struct initializer, though that is allowed elsewhere.
type flagSet struct {
	normalizedIdentifiers bool
	strictTypes           bool
}

// Compiler is a structure defining what we know about the compilation.
type Compiler struct {
	activePackageName string
	b                 *bytecode.ByteCode
	t                 *tokenizer.Tokenizer
	s                 *symbols.SymbolTable
	flags             flagSet // Use to hold parser state flags
}

// New creates a new compiler instance.
func New(name string) *Compiler {
	typeChecking := settings.GetBool(defs.StaticTypesSetting)
	if v, ok := symbols.RootSymbolTable.Get(defs.TypeCheckingVariable); ok {
		typeChecking = (data.Int(v) == defs.StrictTypeEnforcement)
	}

	cInstance := Compiler{
		b: bytecode.New(name),
		t: nil,
		s: symbols.NewRootSymbolTable(name),
		flags: flagSet{
			normalizedIdentifiers: false,
			strictTypes:           typeChecking,
		},
	}

	return &cInstance
}

// NormalizedIdentifiers returns true if this instance of the compiler is folding
// all identifiers to a common (lower) case.
func (c *Compiler) NormalizedIdentifiers() bool {
	return c.flags.normalizedIdentifiers
}

// SetNormalizedIdentifiers sets the flag indicating if this compiler instance is
// folding all identifiers to a common case. This function supports attribute
// chaining for a compiler instance.
func (c *Compiler) SetNormalizedIdentifiers(flag bool) *Compiler {
	c.flags.normalizedIdentifiers = flag

	return c
}

// Set the given symbol table as the default symbol table for
// compilation. This mostly affects how builtins are processed.
// This function supports attribute chaining for a compiler instance.
func (c *Compiler) WithSymbols(s *symbols.SymbolTable) *Compiler {
	c.s = s

	return c
}

// WithTokens supplies the token stream to a compiler. This function supports
// attribute chaining for a compiler instance.
func (c *Compiler) WithTokens(t *tokenizer.Tokenizer) *Compiler {
	c.t = t

	return c
}

// WithNormalization sets the normalization flag. This function supports
// attribute chaining for a compiler instance.
func (c *Compiler) WithNormalization(f bool) *Compiler {
	c.flags.normalizedIdentifiers = f

	return c
}

// AddStandard adds the package-independent standard functions (like len() or make()) to the
// given symbol table.
func (c *Compiler) AddStandard(s *symbols.SymbolTable) bool {
	added := false

	if s == nil {
		return false
	}

	ui.Log(ui.CompilerLogger, "Adding standard functions to %s (%v)", s.Name, s.ID())

	for name, f := range builtins.FunctionDictionary {
		if dot := strings.Index(name, "."); dot < 0 {
			_ = s.SetConstant(name, f.F)
		}
	}

	return added
}
