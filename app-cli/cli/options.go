package cli

const (

	// StringType accepts a string (in quotes if it contains spaces or punctuation)
	StringType = 1

	// IntType accepts a signed integer value
	IntType = 2

	// BooleanType is true if present, or false if not present
	BooleanType = 3

	// BooleanValueType is representation of the boolean value (true/false)
	BooleanValueType = 4

	// Subcommand specifies that the LongName is a command name, and parsing continues with the SubGrammar
	Subcommand = 6

	// StringListType is a string value or a list of string values, separated by commas and enclosed in quotes
	StringListType = 7
)

// Option defines the structure of each option that can be parsed.
type Option struct {
	ShortName            string
	LongName             string
	Aliases              []string
	Description          string
	OptionType           int
	Parameters           int
	ParameterDescription string
	Found                bool
	Required             bool
	Private              bool
	SubGrammar           []Option
	Value                interface{}
	Action               func(c *Context) error
}

// Context is a simple array of Option types, and is used to express
// a grammar.
type Context struct {
	AppName                string
	MainProgram            string
	Description            string
	Command                string
	Grammar                []Option
	Parent                 *Context
	Parameters             []string
	ParameterCount         int
	ExpectedParameterCount int
	ParameterDescription   string
}

// FindGlobal locates the top-most context structure in the chain
// of nested contexts.
func (c *Context) FindGlobal() *Context {
	if c.Parent != nil {
		return c.Parent.FindGlobal()
	}
	return c
}
