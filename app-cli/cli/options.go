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
	SubGrammar           Options
	Value                interface{}
	Action               func(grammar *Options) error
}

// Options is a simple array of Option types, and is used to express
// a grammar.
type Options []Option
