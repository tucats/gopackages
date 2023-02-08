package compiler

import (
	"sort"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/settings"
	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/builtins"
	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/data"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

// requiredPackages is the list of packages that are always imported, regardless
// of user import statements or auto-import profile settings.
var requiredPackages []string = []string{
	"os",
	"profile",
}

// loop is a structure that defines a loop type.
type loop struct {
	parent   *loop
	loopType int
	// Fixup locations for break or continue statements in a
	// loop. These are the addresses that must be fixed up with
	// a target address pointing to exit point or start of the loop.
	breaks    []int
	continues []int
}

// flagSet contains flags that generally identify the state of
// the compiler at any given moment. For example, when parsing
// something like a switch conditional value, the value cannot
// be a struct initializer, though that is allowed elsewhere.
type flagSet struct {
	disallowStructInits   bool
	extensionsEnabled     bool
	normalizedIdentifiers bool
	strictTypes           bool
	testMode              bool
	mainSeen              bool
}

// Compiler is a structure defining what we know about the compilation.
type Compiler struct {
	activePackageName string
	sourceFile        string
	id                string
	b                 *bytecode.ByteCode
	t                 *tokenizer.Tokenizer
	s                 *symbols.SymbolTable
	rootTable         *symbols.SymbolTable
	loops             *loop
	coercions         []*bytecode.ByteCode
	constants         []string
	deferQueue        []int
	packages          map[string]*data.Package
	packageMutex      sync.Mutex
	types             map[string]*data.Type
	functionDepth     int
	blockDepth        int
	statementCount    int
	flags             flagSet // Use to hold parser state flags
	exitEnabled       bool    // Only true in interactive mode
}

// New creates a new compiler instance.
func New(name string) *Compiler {
	extensions := settings.GetBool(defs.ExtensionsEnabledSetting)
	if v, ok := symbols.RootSymbolTable.Get(defs.ExtensionsEnabledSetting); ok {
		extensions = data.Bool(v)
	}

	typeChecking := settings.GetBool(defs.StaticTypesSetting)
	if v, ok := symbols.RootSymbolTable.Get(defs.TypeCheckingVariable); ok {
		typeChecking = (data.Int(v) == defs.StrictTypeEnforcement)
	}

	cInstance := Compiler{
		b:            bytecode.New(name),
		t:            nil,
		s:            symbols.NewRootSymbolTable(name),
		id:           uuid.NewString(),
		constants:    make([]string, 0),
		deferQueue:   make([]int, 0),
		types:        map[string]*data.Type{},
		packageMutex: sync.Mutex{},
		packages:     map[string]*data.Package{},
		flags: flagSet{
			normalizedIdentifiers: false,
			extensionsEnabled:     extensions,
			strictTypes:           typeChecking,
		},
		rootTable: &symbols.RootSymbolTable,
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

// Override the default root symbol table for this compilation. This determines
// where package names are stored/found, for example. This is overridden by the
// web service handlers as they have per-call instances of root. This function
// supports attribute chaining for a compiler instance.
func (c *Compiler) SetRoot(s *symbols.SymbolTable) *Compiler {
	c.rootTable = s
	c.s.SetParent(s)

	return c
}

// If set to true, the compiler allows the "exit" statement. This function supports
// attribute chaining for a compiler instance.
func (c *Compiler) ExitEnabled(b bool) *Compiler {
	c.exitEnabled = b

	return c
}

// TesetMode returns whether the compiler is being used under control
// of the Ego "test" command, which has slightly different rules for
// block constructs.
func (c *Compiler) TestMode() bool {
	return c.flags.testMode
}

// MainSeen indicates if a "package main" has been seen in this
// compilation.
func (c *Compiler) MainSeen() bool {
	return c.flags.mainSeen
}

// SetTestMode is used to set the test mode indicator for the compiler.
// This is set to true only when running in Ego "test" mode. This
// function supports attribute chaining for a compiler instance.
func (c *Compiler) SetTestMode(b bool) *Compiler {
	c.flags.testMode = b

	return c
}

// Set the given symbol table as the default symbol table for
// compilation. This mostly affects how builtins are processed.
// This function supports attribute chaining for a compiler instance.
func (c *Compiler) WithSymbols(s *symbols.SymbolTable) *Compiler {
	c.s = s

	return c
}

// If set to true, the compiler allows the PRINT, TRY/CATCH, etc. statements.
// This function supports attribute chaining for a compiler instance.
func (c *Compiler) ExtensionsEnabled(b bool) *Compiler {
	c.flags.extensionsEnabled = b

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

// AddBuiltins adds the builtins for the named package (or prebuilt builtins if the package name
// is empty).
func (c *Compiler) AddBuiltins(pkgname string) bool {
	added := false

	pkg, _ := bytecode.GetPackage(pkgname)
	symV, _ := pkg.Get(data.SymbolsMDKey)
	syms := symV.(*symbols.SymbolTable)

	ui.Log(ui.CompilerLogger, "### Adding builtin packages to %s package", pkgname)

	functionNames := make([]string, 0)
	for k := range builtins.FunctionDictionary {
		functionNames = append(functionNames, k)
	}

	sort.Strings(functionNames)

	for _, name := range functionNames {
		f := builtins.FunctionDictionary[name]

		if dot := strings.Index(name, "."); dot >= 0 {
			f.Pkg = name[:dot]
			f.Name = name[dot+1:]
			name = f.Name
		} else {
			f.Name = name
		}

		if f.Pkg == pkgname {
			if ui.IsActive(ui.CompilerLogger) {
				debugName := name
				if f.Pkg != "" {
					debugName = f.Pkg + "." + name
				}

				ui.Log(ui.CompilerLogger, "... processing builtin %s", debugName)
			}

			added = true

			if pkgname == "" && c.s != nil {
				syms.SetAlways(name, f.F)
				pkg.Set(name, f.F)
			} else {
				if f.F != nil {
					syms.SetAlways(name, f.F)
					pkg.Set(name, f.F)
				} else {
					syms.SetAlways(name, f.V)
					pkg.Set(name, f.V)
				}
			}
		}
	}

	return added
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

// Get retrieves a compile-time symbol value.
func (c *Compiler) Get(name string) (interface{}, bool) {
	return c.s.Get(name)
}

// normalize performs case-normalization based on the current
// compiler settings.
func (c *Compiler) normalize(name string) string {
	if c.flags.normalizedIdentifiers {
		return strings.ToLower(name)
	}

	return name
}

// normalizeToken performs case-normalization based on the current
// compiler settings for an identifier token.
func (c *Compiler) normalizeToken(t tokenizer.Token) tokenizer.Token {
	if t.IsIdentifier() && c.flags.normalizedIdentifiers {
		return tokenizer.NewIdentifierToken(strings.ToLower(t.Spelling()))
	}

	return t
}

// SetInteractive indicates if the compilation is happening in interactive
// (i.e. REPL) mode. This function supports attribute chaining for a compiler
// instance.
func (c *Compiler) SetInteractive(b bool) *Compiler {
	if b {
		c.functionDepth++
	}

	return c
}

var packageMerge sync.Mutex

// AddPackageToSymbols adds all the defined packages for this compilation
// to the given symbol table. This function supports attribute chaining
// for a compiler instance.
func (c *Compiler) AddPackageToSymbols(s *symbols.SymbolTable) *Compiler {
	ui.Log(ui.CompilerLogger, "Adding compiler packages to %s(%v)", s.Name, s.ID())
	packageMerge.Lock()
	defer packageMerge.Unlock()

	for packageName, packageDictionary := range c.packages {
		// Skip over any metadata
		if strings.HasPrefix(packageName, data.MetadataPrefix) {
			continue
		}

		m := data.NewPackage(packageName)

		keys := packageDictionary.Keys()
		if len(keys) == 0 {
			continue
		}

		for _, k := range keys {
			v, _ := packageDictionary.Get(k)
			// Do we already have a package of this name defined?
			_, found := s.Get(k)
			if found {
				ui.Log(ui.CompilerLogger, "Duplicate package %s already in table", k)
			}

			// If the package name is empty, we add the individual items
			if packageName == "" {
				_ = s.SetConstant(k, v)
			} else {
				// Otherwise, copy the entire map
				m.Set(k, v)
			}
		}
		// Make sure the package is marked as readonly so the user can't modify
		// any function definitions, etc. that are built in.
		m.Set(data.TypeMDKey, data.PackageType(packageName))
		m.Set(data.ReadonlyMDKey, true)

		if packageName != "" {
			s.SetAlways(packageName, m)
		}
	}

	return c
}

// isStatementEnd returns true when the next token is
// the end-of-statement boundary.
func (c *Compiler) isStatementEnd() bool {
	next := c.t.Peek(1)

	return tokenizer.InList(next, tokenizer.EndOfTokens, tokenizer.SemicolonToken, tokenizer.BlockEndToken)
}

// Symbols returns the symbol table map from compilation.
func (c *Compiler) Symbols() *symbols.SymbolTable {
	return c.s
}

// Clone makes a new copy of the current compiler. The withLock flag
// indicates if the clone should respect symbol table locking. This
// function supports attribute chaining for a compiler instance.
func (c *Compiler) Clone(withLock bool) *Compiler {
	cx := Compiler{
		activePackageName: c.activePackageName,
		sourceFile:        c.sourceFile,
		b:                 c.b,
		t:                 c.t,
		s:                 c.s.Clone(withLock),
		rootTable:         c.s.Clone(withLock),
		coercions:         c.coercions,
		constants:         c.constants,
		packageMutex:      sync.Mutex{},
		deferQueue:        []int{},
		flags: flagSet{
			normalizedIdentifiers: c.flags.normalizedIdentifiers,
			extensionsEnabled:     c.flags.extensionsEnabled,
		},
		exitEnabled: c.exitEnabled,
	}

	packages := map[string]*data.Package{}

	c.packageMutex.Lock()
	defer c.packageMutex.Unlock()

	for n, m := range c.packages {
		packageDef := data.NewPackage(n)

		keys := m.Keys()
		for _, k := range keys {
			v, _ := m.Get(k)
			packageDef.Set(k, v)
		}

		packages[n] = packageDef
	}

	// Put the newly created data in the copy of the compiler, with
	// it's own mutex
	cx.packageMutex = sync.Mutex{}
	cx.packages = packages

	return &cx
}
