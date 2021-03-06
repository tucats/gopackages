package functions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/datatypes"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// ReadFile reads a file contents into a string value
func ReadFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	name := util.GetString(args[0])
	if name == "." {
		return ui.Prompt(""), nil
	}

	content, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}

	// Convert []byte to string
	return string(content), nil
}

// Split splits a string into lines
func Split(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	src := util.GetString(args[0])

	// Are we seeing Windows-style line endings? If so, use that as
	// the split boundary.
	if strings.Index(src, "\r\n") > 0 {
		return strings.Split(src, "\r\n"), nil

	}

	// Otherwise, simple split by new-line works fine.
	v := strings.Split(src, "\n")

	// We must recopy this into an array of interfaces to adopt Ego typelessness.
	r := make([]interface{}, 0)
	for _, n := range v {
		r = append(r, n)
	}
	return r, nil
}

// Tokenize splits a string into tokens
func Tokenize(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	src := util.GetString(args[0])
	t := tokenizer.New(src)

	// We must recopy this into an array of interfaces to adopt Ego typelessness.
	r := make([]interface{}, 0)
	for _, n := range t.Tokens {
		r = append(r, n)
	}
	return r, nil
}

// WriteFile writes a string to a file
func WriteFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	fname := util.GetString(args[0])
	text := util.GetString(args[1])

	err := ioutil.WriteFile(fname, []byte(text), 0777)
	return len(text), err
}

// DeleteFile delete a file
func DeleteFile(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	fname := util.GetString(args[0])
	err := os.Remove(fname)
	return err == nil, err
}

// Expand expands a list of file or path names into a list of files.
func Expand(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	path := util.GetString(args[0])
	ext := ""
	if len(args) > 1 {
		ext = util.GetString(args[1])
	}
	list, err := ExpandPath(path, ext)

	// Rewrap as an interface array
	result := []interface{}{}
	for _, item := range list {
		result = append(result, item)
	}
	return result, err
}

// ExpandPath is used to expand a path into a list of fie names
func ExpandPath(path, ext string) ([]string, error) {

	names := []string{}

	// Can we read this as a directory?
	fi, err := ioutil.ReadDir(path)
	if err != nil {
		fn := path
		_, err := ioutil.ReadFile(fn)
		if err != nil {
			fn = path + ext
			_, err = ioutil.ReadFile(fn)
		}
		if err != nil {
			return names, err
		}
		// If we have a default suffix, make sure the pattern matches
		if ext != "" && !strings.HasSuffix(fn, ext) {
			return names, nil
		}
		names = append(names, fn)
		return names, nil
	}

	// Read as a directory
	for _, f := range fi {
		fn := filepath.Join(path, f.Name())
		list, err := ExpandPath(fn, ext)
		if err != nil {
			return names, err
		}
		names = append(names, list...)
	}
	return names, nil
}

// ReadDir implmeents the io.readdir() function
func ReadDir(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	path := util.GetString(args[0])
	result := []interface{}{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return result, err
	}

	for _, file := range files {
		entry := map[string]interface{}{}
		entry["name"] = file.Name()
		entry["directory"] = file.IsDir()
		entry["mode"] = file.Mode().String()
		entry["size"] = int(file.Size())
		entry["modified"] = file.ModTime().String()
		result = append(result, entry)
	}
	return result, nil
}

// This is the generic close() which can be used to close a channel, and maybe
// later other items as well.
func CloseAny(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	switch arg := args[0].(type) {

	case *datatypes.Channel:
		return arg.Close(), nil

	default:
		return nil, NewError("close", InvalidTypeError)
	}
}
