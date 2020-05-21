package compiler

import (
	"io/ioutil"

	"github.com/tucats/gopackages/tokenizer"
)

// Package compiles a package statement
func (c *Compiler) Package() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewTokenError("invalid package name")
	}

	if c.PackageName != "" {
		return c.NewError("package already defined")
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

	name := c.t.Next()
	if len(name) > 2 && name[:1] == "\"" {
		name = name[1 : len(name)-1]
	}
	if c.loops != nil {
		return c.NewError("cannot import inside a loop")
	}

	// Save some state
	savedPackageName := c.PackageName
	savedTokenizer := c.t
	savedBlockDepth := c.blockDepth

	// Read the imported object as a file path

	content, err := ioutil.ReadFile(name)
	if err != nil {
		content, err = ioutil.ReadFile(name + ".solve")
		if err != nil {
			return c.NewTokenError("unable to read file")
		}
	}
	// Convert []byte to string
	text := string(content)

	c.t = tokenizer.New(text)
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

	return nil
}
