# Overview

This project contains a set of packages inended to support a command-line tool written in Go.
They allow for the definition of the command line grammar (including type checking on option
values, missing arguments, etc) and a defined action routine called when a subcommand is
processed successfully.

# Introduction
A simple command line tool defines a grammar for the commands and subcommands, and their 
options. It then calls the app package Run() method which handles parsing and execution
control from then on.

The command line tool developer:

* Defines the grammar structure
* Provides action routines for each subcommand that results in execution of the command function.
* Action routines can use support packages for these additional common functions:
  * Query and set configuration values
  * Handle basic messaging to the console
  * Enable debug logging as needed
  * Generate consistently formatted tabular output
  * Generate JSON tabular output

# Grammar Definition

A grammar definition is just an array of the following structure:

    type Option struct {
        ShortName            string
        LongName             string
        Aliases              []string
        Description          string
        OptionType           int
        Parameters           int
        ParameterDescription string
        Required             bool
        Private              bool
        SubGrammar           Options
        Value                interface{}
        Action               func(grammar *Options) error
    }

### ShortName
The `ShortName` field describes the short option value. For example, you might have an option
named "--debug", but you can also express it as "-d" as the short form of the name. This field
contains the (typically single character) short name, without the leading dash.

If the `ShortName` is not specified or is an empty string, then there is no short option name.

### LongName
The `LongName` field describes the fully-spelled-out name. For example, in the case of "--debug" 
and "-d" as described above in the `ShortName` field, the long name is "debug". This is the 
name displayed as the primary option name in the help, and is always the name used by your
CLI to query the existence or state of the option.


### Aliases
The `Aliases` array is an optional list of alternative spellings of the option. For example,
the built-in `profile` command that manages persistent defaults can be abbreviated `prof` on
the command line; this is implemented by specifying an alias. There can be multiple aliases
for a command. The alias values are never displayed in the help output.


### Description
The `Description` field contains a text sentence describing the option. It does not need to
include the option name or type information, but may contain allowed keywords.


### OptionType
The `OptionType` field describes the type of each grammar item. This can indicate if a value
is expected and what the type of that value should be. It can also indicate a grammar item
is a command with a grammar specific to that command.

The allowed values for `OptionType` are as follows:

| Type             | Description                                          | 
| ---------------- | ---------------------------------------------------- | 
| BooleanType      | The option has no value, and is true if the option is specified, and false if it is not specified.      |
| BooleanValueType | The option value must be an expression of a boolean value: "0", "1", "true", "false", "yes", "no", etc. |
| StringType       | The option value is any string. If the string contains blanks or punctuation marks, it should be in double quotes |
| StringListType   | The option value is a list of strings, separated by commas. |
| IntegerType      | The option must be a valid signed integer value |
| Subcommand       | The value is a keyword that branches to a subsequent grammar definition to continue parsing. |

### Parameters
The `Parameters` field is only used on an item of type `Subcommand`. It tells the grammar processor how many parameters the
subcommand expects to find after parsing the command line. If zero, then no parameters are allowed, which is the default.
A positive number means there must be exactly that number of parameters. A negative number means a variable number of parameters
are allowed, with a maximum being equal to the absolute value of the field. That is, a value of -5 means up to five parameters
are permitted.

### ParameterDescription
The `ParameterDescription` field contains a text string that can optionally assist the formatted help output. If not
specified, the `Parameters` field is used to determine if the output should indicate parameter or perameters are permitted.
But if the parameter is a specific value, like a keyword or file, then this field can create a text desription for the
keyword. For example, the built-in profile command can set the default table output format, using values of "text" or "json".
Because the parameter is very specific in type, the `ParameterDescription` is "type" indicating an output format type.

### Required

### Private

### SubGrammar

### Action

### Value
