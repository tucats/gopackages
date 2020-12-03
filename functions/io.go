package functions

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/tucats/gopackages/app-cli/ui"
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
		return "", err
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
	return err == nil, err
}

// Open opens a file
func Open(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

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

// Close closes a file
func Close(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, NewError("close", InvalidFileIdentifierError, args[0])
	}

	file := handle.(map[string]interface{})
	f := file["id"].(*os.File)
	err := f.Close()

	s.DeleteAlways(id)

	return err == nil, err
}

// ReadString closes a file
func ReadString(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, NewError("read", InvalidFileIdentifierError, args[0])
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

// WriteString closes a file
func WriteString(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {

	id := util.GetString(args[0])
	handle, found := s.Get(id)

	if !found {
		return false, NewError("write", InvalidFileIdentifierError, args[0])
	}

	file := handle.(map[string]interface{})
	f := file["id"].(*os.File)

	l, err := f.WriteString(util.GetString(args[1]) + "\n")
	return l, err

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
		// ui.Debug("+++ scan file      \"%s\"", fn)

		names = append(names, fn)
		return names, nil
	}
	// ui.Debug("+++ scan directory \"%s\"", path)

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
