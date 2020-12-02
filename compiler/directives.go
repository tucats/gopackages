package compiler

import (
	"fmt"

	"github.com/tucats/gopackages/bytecode"
	"github.com/tucats/gopackages/functions"
	"github.com/tucats/gopackages/symbols"
	"github.com/tucats/gopackages/tokenizer"
	"github.com/tucats/gopackages/util"
)

// Directive processes a compiler directive. These become symbols generated
// at compile time that are copied to the compiler's symbol table for processing
// elsewhere.
func (c *Compiler) Directive() error {
	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewError(InvalidDirectiveError, name)
	}

	switch name {
	case "error":
		return c.Error()
	case "template":
		return c.Template()
	case "pass":
		return c.TestPass()
	case "assert":
		return c.Assert()
	case "fail":
		return c.Fail()
	case "test":
		return c.Test()
	default:
		return c.NewError(InvalidDirectiveError, name)
	}
}

// Template implements the template compiler directive
func (c *Compiler) Template() error {

	// Get the template name
	name := c.t.Next()
	if !tokenizer.IsSymbol(name) {
		return c.NewError(InvalidSymbolError, name)
	}

	// Get the template string definition
	bc, err := c.Expression()
	if err != nil {
		return err
	}
	c.b.Append(bc)
	c.b.Emit(bytecode.Template, name)
	c.b.Emit(bytecode.SymbolCreate, name)
	c.b.Emit(bytecode.Store, name)

	return nil
}

// Test compiles the @test directive
func (c *Compiler) Test() error {

	s := c.t.Next()
	if s[:1] == "\"" {
		s = s[1 : len(s)-1]
	}

	test := map[string]interface{}{}
	//test["__readonly"] = true
	test["assert"] = TestAssert
	test["fail"] = TestFail
	test["isType"] = TestIsType
	test["description"] = s

	c.s.SetAlways("testing", test)

	// Generate code to update the description (this is required for the
	// cases of the ego test command running multiple tests as a single
	// stream)
	c.b.Emit(bytecode.Push, s)
	c.b.Emit(bytecode.Load, "testing")
	c.b.Emit(bytecode.Push, "description")
	c.b.Emit(bytecode.StoreIndex)

	// Generate code to report that the test is starting.
	c.b.Emit(bytecode.Push, "TEST: ")
	c.b.Emit(bytecode.Print)
	c.b.Emit(bytecode.Load, "testing")
	c.b.Emit(bytecode.Push, "description")
	c.b.Emit(bytecode.Member)
	c.b.Emit(bytecode.Print)
	c.b.Emit(bytecode.Newline)

	return nil
}

// TestAssert implements the testing.assert() function
func TestAssert(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) < 1 || len(args) > 2 {
		return nil, functions.NewError("assert", functions.ArgumentCountError)
	}

	// Figure out the test name. If not found, use "test"
	name := "test"
	if m, ok := s.Get("testing"); ok {
		if structMap, ok := m.(map[string]interface{}); ok {
			if nameString, ok := structMap["description"]; ok {
				name = util.GetString(nameString)
			}
		}
	}

	b := util.GetBool(args[0])
	if !b {
		msg := TestingAssertError
		if len(args) > 1 {
			msg = util.GetString(args[1])
		}
		return nil, fmt.Errorf("%s in %s", msg, name)
	}
	return true, nil
}

// TestIsType implements the testing.assert() function
func TestIsType(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, functions.NewError("istype", functions.ArgumentCountError)
	}

	// Figure out the test name. If not found, use "test"
	name := "test"
	if m, ok := s.Get("testing"); ok {
		if structMap, ok := m.(map[string]interface{}); ok {
			if nameString, ok := structMap["name"]; ok {
				name = util.GetString(nameString)
			}
		}
	}

	// Use the Type() function to get a string representation of the type
	got, _ := functions.FunctionType(s, args[0:1])
	expected := util.GetString(args[1])
	b := (expected == got)
	if !b {
		msg := fmt.Sprintf("testing.isType(\"%s\" != \"%s\") failure", got, expected)
		if len(args) > 2 {
			msg = util.GetString(args[2])
		}
		return nil, fmt.Errorf("%s in %s", msg, name)
	}
	return true, nil
}

// TestFail implements the testing.fail() function which generates a fatal
// error.
func TestFail(s *symbols.SymbolTable, args []interface{}) (interface{}, error) {
	msg := "testing.fail()"
	if len(args) == 1 {
		msg = util.GetString(args[0])
	}

	// Figure out the test name. If not found, use "test"
	name := "test"
	if m, ok := s.Get("testing"); ok {
		fmt.Printf("DEBUG: found testing package\n")
		if structMap, ok := m.(map[string]interface{}); ok {
			fmt.Printf("DEBUG: found map\n")
			if nameString, ok := structMap["description"]; ok {
				fmt.Printf("DEBUG: found name member\n")
				name = util.GetString(nameString)
			}
		}
	}

	return nil, fmt.Errorf("%s in %s", msg, name)
}

// Assert implements the @assert directive
func (c *Compiler) Assert() error {

	c.b.Emit(bytecode.Load, "testing")
	c.b.Emit(bytecode.Push, "assert")
	c.b.Emit(bytecode.Member)

	argCount := 1
	code, err := c.Expression()
	if err != nil {
		return err
	}
	c.b.Append(code)

	next := c.t.Peek(1)
	if next != "@" && next != ";" && next != tokenizer.EndOfTokens {
		code, err := c.Expression()
		if err != nil {
			return err
		}
		c.b.Append(code)
		argCount = 2
	}

	c.b.Emit(bytecode.Call, argCount)

	return nil
}

// Fail implements the @fail directive
func (c *Compiler) Fail() error {
	next := c.t.Peek(1)
	if next != "@" && next != ";" && next != tokenizer.EndOfTokens {
		code, err := c.Expression()
		if err != nil {
			return err
		}
		c.b.Append(code)
	} else {
		c.b.Emit(bytecode.Push, "@fail error signal")
	}
	c.b.Emit(bytecode.Panic, true)
	return nil
}

// TestPass implements the @pass directive
func (c *Compiler) TestPass() error {

	c.b.Emit(bytecode.Push, "PASS: ")
	c.b.Emit(bytecode.Print)

	c.b.Emit(bytecode.Load, "testing")
	c.b.Emit(bytecode.Push, "description")
	c.b.Emit(bytecode.Member)
	c.b.Emit(bytecode.Print)
	c.b.Emit(bytecode.Newline)
	return nil
}

// Error implements the @error directive
func (c *Compiler) Error() error {
	c.b.Emit(bytecode.AtLine, c.t.Line[c.t.TokenP])
	code, err := c.Expression()
	if err == nil {
		c.b.Append(code)
		c.b.Emit(bytecode.Panic, false) // Does not cause fatal error
	}
	return err
}
