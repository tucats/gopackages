package compiler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
)

// Package compiles a package statement
func (c *Compiler) Package() error {

	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewTokenError("invalid package name")
	}

	name = strings.ToLower(name)

	if (c.PackageName != "") && (c.PackageName != name) {
		return c.NewError("cannot redefine package name")
	}
	c.PackageName = name

	// Create a named struct that can be initialized with the symbol names.
	// This is done by creating a source table and then merging it with the
	// active table.
	tmp := symbols.NewSymbolTable("")
	tmp.SetAlways(name, map[string]interface{}{"__parent": name})
	c.s.Merge(tmp)

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

	isList := false
	if c.t.IsNext("(") {
		isList = true
	}

	parsing := true
	for parsing {

		// Make sure that if this isn't the list format of an import, we only do this once.
		if !isList {
			parsing = false
		}
		fileName := c.t.Next()

		// End of the list? If so, break out
		if isList && fileName == ")" {
			break
		}
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
		packageName = strings.ToLower(packageName)

		// If this is an import of a package already processed, no work to do.
		if _, found := c.s.Get(packageName); found {
			ui.Debug("+++ Previously imported %s, skipping", packageName)
			continue
		}

		// If this is an import of the package we're currently importing, no work to do.
		if packageName == c.PackageName {
			continue
		}

		builtinsAdded := c.AddBuiltins(packageName)
		if builtinsAdded {
			ui.Debug("+++ Adding builtins for package " + packageName)
		} else {
			ui.Debug("+++ No builtins for package " + packageName)
		}

		// Save some state
		savedPackageName := c.PackageName
		savedTokenizer := c.t
		savedBlockDepth := c.blockDepth
		savedStatementCount := c.statementCount

		// Read the imported object as a file path
		text, err := c.ReadFile(fileName)
		if err != nil {

			// If it wasn't found but we did add some builtins, good enough.
			// Skip past the filename that was rejected by c.Readfile()...
			if builtinsAdded {
				c.t.Advance(1)

				if !isList || c.t.IsNext(")") {
					break
				}
				continue
			}

			// Nope, import had no effect.
			return err
		}

		ui.Debug("+++ Adding source for package " + packageName)
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
		if !isList {
			break
		}
		if isList && c.t.Next() == ")" {
			break
		}
	}
	return nil
}

// ReadFile reads the text from a file into a string
func (c *Compiler) ReadFile(name string) (string, error) {

	s, err := c.ReadDirectory(name)
	if err == nil {
		return s, nil
	}
	ui.Debug("+++ Reading package file %s", name)
	// Not a directory, try to read the file
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

// ReadDirectory reads all the files in a directory into a single string.
func (c *Compiler) ReadDirectory(name string) (string, error) {

	var b strings.Builder
	r := os.Getenv("SOLVE_PATH")
	dirname := filepath.Join(r, "lib", name)

	fi, err := ioutil.ReadDir(dirname)
	if err != nil {
		return "", err
	}

	ui.Debug("+++ Reading package directory %s", dirname)
	for _, f := range fi {
		if !f.IsDir() {

			fname := filepath.Join(dirname, f.Name())
			t, err := c.ReadFile(fname)
			if err != nil {
				return "", err
			}
			b.WriteString(t)
			b.WriteString("\n")
		}
	}

	return b.String(), nil
}
