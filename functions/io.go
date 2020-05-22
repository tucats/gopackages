package functions

import (
	"io/ioutil"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// FunctionReadFile reads a file contents into a string value
func FunctionReadFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	name := util.GetString(args[0])

	if name == "." {
		return ui.Prompt(""), nil
	}

	content, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}

	// Convert []byte to string
	return string(content), nil
}

// FunctionSplit splits a string into lines
func FunctionSplit(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	src := util.GetString(args[0])

	// Are we seeing Windows-style line endings? If so, use that as
	// the split boundary.
	if strings.Index(src, "\r\n") > 0 {
		return strings.Split(src, "\r\n"), nil

	}

	// Otherwise, simple split by new-line works fine.
	v := strings.Split(src, "\n")

	// We must recopy this into an array of interfaces to adopt Solve typelessness.
	r := make([]interface{}, 0)
	for _, n := range v {
		r = append(r, n)
	}
	return r, nil
}

// FunctionTokenize splits a string into tokens
func FunctionTokenize(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	src := util.GetString(args[0])
	t := tokenizer.New(src)

	// We must recopy this into an array of interfaces to adopt Solve typelessness.
	r := make([]interface{}, 0)
	for _, n := range t.Tokens {
		r = append(r, n)
	}
	return r, nil
}
