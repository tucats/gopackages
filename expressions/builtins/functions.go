package builtins

import (
	"fmt"
	"math"
	"reflect"
	"sort"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/errors"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
	"github.com/tucats/gopackages/i18n"
)

// FunctionDefinition is an element in the function dictionary. This
// defines each function that is implemented as native Go code (a
// "builtin" function).
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string
	// Pkg is the package that contains the function, if it is
	// a builtin package member.
	Pkg string

	// Min is the minimum number of arguments the function can accept.
	Min int

	// Max is the maximum number of arguments the function can accept.
	Max int

	// ErrReturn is true if the function returns a tuple containing the
	// function result and an error return.
	ErrReturn bool

	// FullScope indicates if this function is allowed to access the
	// entire scope tree of the running program.
	FullScope bool

	// F is the address of the function implementation
	F interface{}

	// V is a value constant associated with this name.
	V interface{}

	// D is a function declaration object that details the
	// parameter and return types.
	D *data.Declaration
}

// Any is a constant that defines that a function can have as many arguments
// as desired.
const Any = math.MaxInt32

// FunctionDictionary is the dictionary of functions. As functions are determined
// to allow the return of both a value and an error as multi-part results, add the
// ErrReturn:true flag to each function definition.
var FunctionDictionary = map[string]FunctionDefinition{
	"$new": {Min: 1, Max: 1, F: New},
	"index": {Min: 2, Max: 2, F: Index, D: &data.Declaration{
		Name: "index",
		Parameters: []data.Parameter{
			{
				Name: "item",
				Type: data.InterfaceType,
			},
			{
				Name: "index",
				Type: data.InterfaceType,
			},
		},
		Returns: []*data.Type{data.IntType},
	}},
	"len": {Min: 1, Max: 1, F: Length, D: &data.Declaration{
		Name: "len",
		Parameters: []data.Parameter{
			{
				Name: "item",
				Type: data.InterfaceType,
			},
		},
		Returns: []*data.Type{data.IntType},
	}},
	"make": {Min: 2, Max: 2, F: Make, D: &data.Declaration{
		Name: "make",
		Parameters: []data.Parameter{
			{
				Name: "t",
				Type: data.TypeType,
			},
			{
				Name: "count",
				Type: data.IntType,
			},
		},
		Returns: []*data.Type{data.IntType},
	}},
	"sizeof": {Min: 1, Max: 1, F: SizeOf, D: &data.Declaration{
		Name: "sizeof",
		Parameters: []data.Parameter{
			{
				Name: "item",
				Type: data.InterfaceType,
			},
		},
		Returns: []*data.Type{data.IntType},
	}},
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbolTable *symbols.SymbolTable) {
	ui.Log(ui.CompilerLogger, "+++ Adding in builtin functions to symbol table %s", symbolTable.Name)

	functionNames := make([]string, 0)
	for k := range FunctionDictionary {
		functionNames = append(functionNames, k)
	}

	sort.Strings(functionNames)

	for _, n := range functionNames {
		d := FunctionDictionary[n]

		if d.D != nil {
			data.RegisterDeclaration(d.D)
		}

		if dot := strings.Index(n, "."); dot >= 0 {
			d.Pkg = n[:dot]
			n = n[dot+1:]
		}

		_ = symbolTable.SetWithAttributes(n, d.F, symbols.SymbolAttribute{Readonly: true})
	}
}

// FindFunction returns the function definition associated with the
// provided function pointer, if one is found.
func FindFunction(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) *FunctionDefinition {
	sf1 := reflect.ValueOf(f)

	for _, d := range FunctionDictionary {
		if d.F != nil { // Only function entry points have an F value
			sf2 := reflect.ValueOf(d.F)
			if sf1.Pointer() == sf2.Pointer() {
				return &d
			}
		}
	}

	return nil
}

// FindName returns the name of a function from the dictionary if one is found.
func FindName(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) string {
	sf1 := reflect.ValueOf(f)

	for name, d := range FunctionDictionary {
		if d.F != nil {
			sf2 := reflect.ValueOf(d.F)
			if sf1.Pointer() == sf2.Pointer() {
				return name
			}
		}
	}

	return ""
}

func CallBuiltin(s *symbols.SymbolTable, name string, args ...interface{}) (interface{}, error) {
	var fdef = FunctionDefinition{}

	found := false

	for fn, d := range FunctionDictionary {
		if fn == name {
			fdef = d
			found = true
		}
	}

	if !found {
		return nil, errors.ErrInvalidFunctionName.Context(name)
	}

	if len(args) < fdef.Min || len(args) > fdef.Max {
		return nil, errors.ErrPanic.Context(i18n.E("arg.count"))
	}

	fn, ok := fdef.F.(func(*symbols.SymbolTable, []interface{}) (interface{}, error))
	if !ok {
		return nil, errors.ErrPanic.Context(fmt.Errorf(i18n.E("function.pointer",
			map[string]interface{}{"ptr": fdef.F})))
	}

	return fn(s, args)
}

func AddFunction(s *symbols.SymbolTable, fd FunctionDefinition) error {
	// Make sure not a collision
	if _, ok := FunctionDictionary[fd.Name]; ok {
		return errors.ErrFunctionAlreadyExists
	}

	FunctionDictionary[fd.Name] = fd

	return nil
}

func stubFunction(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	return nil, errors.ErrInvalidFunctionName
}

// extensions retrieves the boolean indicating if extensions are supported. This can
// be used to do runtime checks for etended featues of builtins.
func extensions() bool {
	f := false
	if v, ok := symbols.RootSymbolTable.Get(defs.ExtensionsVariable); ok {
		f = data.Bool(v)
	}

	return f
}
