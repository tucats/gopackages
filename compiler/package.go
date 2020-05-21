package compiler

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/tucats/gopackages/tokenizer"
)

// Package compiles a package statement
func (c *Compiler) Package() error {

	if c.statementCount > 1 {
		return c.NewError("package statement must be first")
	}
	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewTokenError("invalid package name")
	}

	if (c.PackageName != "") && (c.PackageName != name) {
		return c.NewError("cannot redefine package name")
	}
	c.PackageName = name

	// Create a named struct that can be initialized with the symbol names
	c.s.Set(name, map[string]interface{}{})

	return nil
}

// Import handles the import statement
func (c *Compiler) Import() error {

	if c.blockDepth > 0 {
		return c.NewError("cannot import inside a block")
	}
	if c.loops != nil {
		return c.NewError("cannot import inside a loop")
	}

	fileName := c.t.Next()
	if len(fileName) > 2 && fileName[:1] == "\"" {
		fileName = fileName[1 : len(fileName)-1]
	}
	if c.loops != nil {
		return c.NewError("cannot import inside a loop")
	}

	// Get the package name from the given string. If this is
	// a file system name, remove the extension if present.
	packageName := filepath.Base(fileName)
	if filepath.Ext(packageName) != "" {
		packageName = packageName[:len(filepath.Ext(packageName))]
	}

	// Save some state
	savedPackageName := c.PackageName
	savedTokenizer := c.t
	savedBlockDepth := c.blockDepth
	savedStatementCount := c.statementCount

	// Read the imported object as a file path
	text, err := c.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Set up the new compiler settings
	c.statementCount = 0
	c.t = tokenizer.New(text)
	c.PackageName = packageName

	for !c.t.AtEnd() {
		err := c.Statement()
		if err != nil {
			return err
		}
	}

	// Reset the token stream we were working on
	c.t = savedTokenizer
	c.PackageName = savedPackageName
	c.blockDepth = savedBlockDepth
	c.statementCount = savedStatementCount
	return nil
}

// ReadFile reads the text from a file into a string
func (c *Compiler) ReadFile(name string) (string, error) {

	content, err := ioutil.ReadFile(name)
	if err != nil {
		content, err = ioutil.ReadFile(name + ".solve")
		if err != nil {
			r := os.Getenv("SOLVE_PATH")
			fn := filepath.Join(r, "lib", name+".solve")
			content, err = ioutil.ReadFile(fn)
			if err != nil {
				c.t.Advance(-1)
				return "", c.NewStringError("unable to read import file", err.Error())
			}
		}
	}

	// Convert []byte to string
	return string(content), nil
}
