package functions

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
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

// FunctionWriteFile writes a string to a file
func FunctionWriteFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	fname := util.GetString(args[0])
	text := util.GetString(args[1])

	err := ioutil.WriteFile(fname, []byte(text), 0777)
	return err == nil, err
}

// FunctionOpen opens a file
func FunctionOpen(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	fname := util.GetString(args[0])
	outputFile := false
	if len(args) > 1 {
		mode := strings.ToLower(util.GetString(args[1]))
		if mode == "true" || mode == "create" || mode == "output" {
			outputFile = true
		}
	}

	var f *os.File
	var err error

	if outputFile {
		f, err = os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Open(fname)
	}
	if err != nil {
		return nil, err
	}

	id := "__file-" + uuid.New().String()

	file := map[string]interface{}{}
	file["id"] = f
	s.SetAlways(id, file)
	return id, nil
}

// FunctionClose closes a file
func FunctionClose(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, errors.New("invalid file identifier")
	}

	file := handle.(map[string]interface{})
	f := file["id"].(*os.File)
	err := f.Close()

	s.DeleteAlways(id)

	return err == nil, err
}

// FunctionReadString closes a file
func FunctionReadString(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, errors.New("invalid file identifier")
	}

	file := handle.(map[string]interface{})
	f := file["id"].(*os.File)

	var scanner *bufio.Scanner

	scanX, found := file["scanner"]
	if !found {
		scanner = bufio.NewScanner(f)
		file["scanner"] = scanner
		s.Set(id, file)
	} else {
		scanner = scanX.(*bufio.Scanner)
	}
	scanner.Scan()
	return scanner.Text(), nil
}

// FunctionWriteString closes a file
func FunctionWriteString(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, errors.New("invalid file identifier")
	}

	file := handle.(map[string]interface{})
	f := file["id"].(*os.File)

	l, err := f.WriteString(util.GetString(args[1]) + "\n")
	return l, err

}

// FunctionDeleteFile delete a file
func FunctionDeleteFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	fname := util.GetString(args[0])

	err := os.Remove(fname)
	return err == nil, err
}
