package functions

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
)

// FunctionDefinition is an element in the function dictionary
type FunctionDefinition struct {
	Name string
	Pkg  string
	Min  int
	Max  int
	F    interface{}
}

// Any is a constant that defines that a function can have as many arguments
// as desired.
const Any = 999999

// FunctionDictionary is the dictionary of functions
var FunctionDictionary = map[string]FunctionDefinition{
	"array":                {Min: 1, Max: 2, F: Array},
	"bool":                 {Min: 1, Max: 1, F: Bool},
	"error":                {Min: 1, Max: 1, F: Signal},
	"float":                {Min: 1, Max: 1, F: Float},
	"index":                {Min: 2, Max: 2, F: Index},
	"int":                  {Min: 1, Max: 1, F: Int},
	"len":                  {Min: 1, Max: 1, F: Length},
	"max":                  {Min: 1, Max: Any, F: Max},
	"members":              {Min: 1, Max: 1, F: Members},
	"min":                  {Min: 1, Max: Any, F: Min},
	"new":                  {Min: 1, Max: 1, F: New},
	"sort":                 {Min: 1, Max: Any, F: Sort},
	"string":               {Min: 1, Max: 1, F: String},
	"sum":                  {Min: 1, Max: Any, F: Sum},
	"type":                 {Min: 1, Max: 1, F: Type},
	"cipher.Create":        {Min: 1, Max: 2, F: CreateToken},
	"cipher.Decrypt":       {Min: 2, Max: 2, F: Decrypt},
	"cipher.Encrypt":       {Min: 2, Max: 2, F: Encrypt},
	"cipher.Hash":          {Min: 1, Max: 1, F: Hash},
	"cipher.Token":         {Min: 1, Max: 2, F: Extract},
	"cipher.Validate":      {Min: 1, Max: 2, F: Validate},
	"fmt.Print":            {Min: 1, Max: Any, F: Print},
	"fmt.Printf":           {Min: 1, Max: Any, F: Printf},
	"fmt.Println":          {Min: 0, Max: Any, F: Println},
	"fmt.Sprintf":          {Min: 1, Max: Any, F: Sprintf},
	"io.Close":             {Min: 1, Max: 1, F: Close},
	"io.Delete":            {Min: 1, Max: 1, F: DeleteFile},
	"io.Expand":            {Min: 1, Max: 2, F: Expand},
	"io.Open":              {Min: 1, Max: 2, F: Open},
	"io.ReadDir":           {Min: 1, Max: 1, F: ReadDir},
	"io.ReadFile":          {Min: 1, Max: 1, F: ReadFile},
	"io.ReadString":        {Min: 1, Max: 1, F: ReadString},
	"io.Split":             {Min: 1, Max: 1, F: Split},
	"io.WriteFile":         {Min: 2, Max: 2, F: WriteFile},
	"io.WriteString":       {Min: 1, Max: 2, F: WriteString},
	"json.UnMarshal":       {Min: 1, Max: 1, F: Decode},
	"json.Marshal":         {Min: 1, Max: Any, F: Encode},
	"json.MarshalIndented": {Min: 1, Max: Any, F: EncodeFormatted},
	"math.Abs":             {Min: 1, Max: 1, F: Abs},
	"math.Log":             {Min: 1, Max: 1, F: Log},
	"math.Sqrt":            {Min: 1, Max: 1, F: Sqrt},
	"profile.Delete":       {Min: 1, Max: 1, F: ProfileDelete},
	"profile.Get":          {Min: 1, Max: 1, F: ProfileGet},
	"profile.Keys":         {Min: 0, Max: 0, F: ProfileKeys},
	"profile.Set":          {Min: 1, Max: 2, F: ProfileSet},
	"strings.Chars":        {Min: 1, Max: 1, F: Chars},
	"strings.Format":       {Min: 0, Max: Any, F: Format},
	"strings.Index":        {Min: 2, Max: 2, F: Index},
	"strings.Ints":         {Min: 1, Max: 1, F: Ints},
	"strings.Left":         {Min: 2, Max: 2, F: Left},
	"strings.ToLower":      {Min: 1, Max: 1, F: Lower},
	"strings.Right":        {Min: 2, Max: 2, F: Right},
	"strings.String":       {Min: 1, Max: Any, F: ToString},
	"strings.Substring":    {Min: 3, Max: 3, F: Substring},
	"strings.Template":     {Min: 1, Max: 2, F: Template},
	"strings.Tokenize":     {Min: 1, Max: 1, F: Tokenize},
	"strings.ToUpper":      {Min: 1, Max: 1, F: Upper},
	"time.Add":             {Min: 2, Max: 2, F: TimeAdd},
	"time.Now":             {Min: 0, Max: 0, F: TimeNow},
	"time.Sleep":           {Min: 1, Max: 1, F: Sleep},
	"time.Subtract":        {Min: 2, Max: 2, F: TimeSub},
	"util.Coerce":          {Min: 2, Max: 2, F: Coerce},
	"util.Exit":            {Min: 0, Max: 1, F: Exit},
	"util.Getenv":          {Min: 1, Max: 1, F: GetEnv},
	"util.Normalize":       {Min: 2, Max: 2, F: Normalize},
	"util.Symbols":         {Min: 0, Max: 1, F: FormatSymbols},
	"util.UUID":            {Min: 0, Max: 0, F: UUID},
}

// AddBuiltins adds or overrides the default function library in the symbol map.
// Function names are distinct in the map because they always have the "()"
// suffix for the key.
func AddBuiltins(symbols *symbols.SymbolTable) {

	ui.Debug(ui.CompilerLogger, "+++ Adding in builtin functions to symbol table %s", symbols.Name)
	for n, d := range FunctionDictionary {

		if dot := strings.Index(n, "."); dot >= 0 {
			d.Pkg = n[:dot]
			n = n[dot+1:]
		}

		if d.Pkg == "" {
			_ = symbols.SetAlways(n, d.F)
		} else {
			// Does package already exist? IF not, make it. The package
			// is just a struct containing where each member is a function
			// definition.
			p, found := symbols.Get(d.Pkg)
			if !found {
				p = map[string]interface{}{}
				p.(map[string]interface{})["readonly"] = true
				ui.Debug(ui.CompilerLogger, "    AddBuiltins creating new package %s", d.Pkg)
			}
			p.(map[string]interface{})[n] = d.F
			ui.Debug(ui.CompilerLogger, "    adding builtin %s to %s", n, d.Pkg)
			_ = symbols.SetAlways(d.Pkg, p)
		}
	}
}

// FindFunction returns the function definition associated with the
// provided function pointer, if one is found.
func FindFunction(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) *FunctionDefinition {

	sf1 := reflect.ValueOf(f)

	for _, d := range FunctionDictionary {
		sf2 := reflect.ValueOf(d.F)
		if sf1.Pointer() == sf2.Pointer() {
			return &d
		}
	}
	return nil
}

// FindName returns the name of a function from the dictionary if one is found
func FindName(f func(*symbols.SymbolTable, []interface{}) (interface{}, error)) string {

	sf1 := reflect.ValueOf(f)

	for name, d := range FunctionDictionary {
		sf2 := reflect.ValueOf(d.F)
		if sf1.Pointer() == sf2.Pointer() {
			return name
		}
	}

	return ""
}

func CallBuiltin(s *symbols.SymbolTable, name string, args ...interface{}) (interface{}, error) {

	// Search the dictionary for a name match
	var fdef = FunctionDefinition{}
	found := false
	for fn, d := range FunctionDictionary {
		if fn == name {
			fdef = d
			found = true
		}
	}
	if !found {
		return nil, errors.New("no such function: " + name)
	}

	if len(args) < fdef.Min || len(args) > fdef.Max {
		return nil, errors.New("incorrect number of arguments")
	}

	fn, ok := fdef.F.(func(*symbols.SymbolTable, []interface{}) (interface{}, error))
	if !ok {
		return nil, fmt.Errorf("unable to convert %#v to function pointer", fdef.F)
	}
	return fn(s, args)
}
