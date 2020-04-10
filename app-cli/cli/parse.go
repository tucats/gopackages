package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/cli/ui"
)

// MainProgramDescription is the description of the main program
var MainProgramDescription string

// Parameters stores those command line items not parsed as part of the
// grammar.
var Parameters []string

// expectedParameters describes how many parameters are expected, which
// comes from any parameter value processed in an subcommand. A negative
// number means any number up to the absolute value.
var expectedParameters = 0

// parameterDescription is a string that can be specified in a subcommand
// option to describe parameters. For example, the set-output builtin
// command uses a description of "type"
var parameterDescription = ""

// Action is the function to invoke on the last subcommand found in parsing.
var Action func(c *Context) error

// Globals is a copy of the outermost grammar definition, and is used to find
// global values later.
var Globals *[]Option

// Parse accepts a grammar definition and parses the current argument
// list against that grammar. Unrecognized options or subcommands, as
// well as invalid values are reported as an error. If there is an
// action routine associated with an option or a subcommand, that
// action is executed.
func (c *Context) Parse(description string) error {

	args := os.Args
	c.MainProgram = filepath.Base(args[0])
	c.Description = ""
	MainProgramDescription = description
	c.Command = ""
	Action = nil

	Globals = &c.Grammar

	// If there are no arguments other than the main program name, dump out the help by default.
	if len(args) == 1 {
		ShowHelp(c)
	}

	// Start parsing using the top-level grammar.
	return c.ParseGrammar(args[1:])
}

// ParseGrammar accepts an argument list and a grammar definition, and parses
func (c *Context) ParseGrammar(args []string) error {

	lastArg := len(args)
	var err error
	parametersOnly := false
	helpVerb := true

	for currentArg := 0; currentArg < lastArg; currentArg++ {

		option := args[currentArg]
		ui.Debug("Processing token: %s", option)

		var location *Option
		var name string
		var value string
		isShort := false

		// Are we now only eating parameter values?
		if parametersOnly {
			Parameters = append(Parameters, option)
			count := len(Parameters)
			ui.Debug(fmt.Sprintf("added parameter %d", count))
			continue
		}

		// Handle the special cases automatically.
		if (helpVerb && option == "help") || option == "-h" || option == "--help" {
			ShowHelp(c)
		}
		if option == "--" {
			parametersOnly = true
			helpVerb = false
			continue
		}

		// If it's a long-name option, search for it.
		if len(option) > 2 && option[:2] == "--" {
			name = option[2:]
		} else if len(option) > 1 && option[:1] == "-" {
			name = option[1:]
			isShort = true
		}

		value = ""
		hasValue := false
		if equals := strings.Index(name, "="); equals >= 0 {
			value = name[equals+1:]
			name = name[:equals]
			hasValue = true
		}

		location = nil
		if name > "" {
			for n, entry := range c.Grammar {

				if (isShort && entry.ShortName == name) || (!isShort && entry.LongName == name) {
					location = &(c.Grammar[n])
					break
				}
			}
		}

		// If it was an option (short or long) and not found, this is an error.
		if name != "" && location == nil {
			return errors.New("Unknown command line option: " + option)
		}
		// It could be a parameter, or a subcommand.
		if location == nil {

			// Is it a subcommand?
			for _, entry := range c.Grammar {

				// Is it one of the aliases permitted?
				isAlias := false
				for _, n := range entry.Aliases {
					if option == n {
						isAlias = true
						break
					}
				}
				if (isAlias || entry.LongName == option) && entry.OptionType == Subcommand {

					// We're doing a subcommand! Create a new context that defines the
					// next level down. It should include the current context information,
					// and an updated grammar tree, command text, and description adapted
					// for this subcommand.
					subContext := *c
					if entry.Value != nil {
						subContext.Grammar = entry.Value.([]Option)
					} else {
						subContext.Grammar = []Option{}
					}
					subContext.Command = c.Command + entry.LongName + " "
					subContext.Description = entry.Description

					entry.Found = true
					expectedParameters = entry.Parameters
					parameterDescription = entry.ParameterDescription

					if entry.Action != nil {
						Action = entry.Action
						ui.Debug("Adding action routine")
					}
					ui.Debug("Transferring control to subgrammar for %s", entry.LongName)
					return subContext.ParseGrammar(args[currentArg+1:])
				}
			}

			// Not a subcommand, just save it as an unclaimed parameter
			Parameters = append(Parameters, option)
			count := len(Parameters)
			ui.Debug(fmt.Sprintf("Unclaimed token added parameter %d", count))

		} else {
			ui.Debug("processing option")

			location.Found = true

			// If it's not a boolean type, see it already has a value from the = construct.
			// If not, claim the next argument as the value.
			if location.OptionType != BooleanType {

				if !hasValue {
					currentArg = currentArg + 1
					if currentArg >= lastArg {
						return errors.New("Error, missing option value for " + name)
					}
					value = args[currentArg]
					hasValue = true
				}
			}

			// Validate the argument type and store the appropriately typed value.
			switch location.OptionType {
			case BooleanType:
				location.Value = true

			case BooleanValueType:
				valid := false
				for _, x := range []string{"1", "true", "t", "yes", "y"} {
					if strings.ToLower(value) == x {
						location.Value = true
						valid = true
						break
					}
				}

				if !valid {
					for _, x := range []string{"0", "false", "f", "no", "n"} {
						if strings.ToLower(value) == x {
							location.Value = false
							valid = true
							break
						}
					}
				}

				if !valid {
					return errors.New("option --" + location.LongName + ": invalid boolean value \"" + value + "\"")
				}

			case StringType:
				location.Value = value

			case StringListType:
				// The value is a comma-separated list of items. Parse each individual item
				// and make a string array of the list as the value.

				list := strings.Split(value, ",")
				for n := 0; n < len(list); n++ {
					list[n] = strings.TrimSpace(list[n])
				}
				location.Value = list
			case IntType:
				i, err := strconv.Atoi(value)
				if err != nil {
					return errors.New("option --" + location.LongName + ": invalid integer")
				}
				location.Value = i
			}

			// After parsing the option value, if there is an action routine, call it
			if location.Action != nil {
				err = location.Action(c)
				if err != nil {
					break
				}
			}
		}
	}

	// Whew! Everything parsed and in it's place. Before we wind up, let's verify that
	// all required options were in fact found.

	for _, entry := range c.Grammar {

		if entry.Required && !entry.Found {
			err = errors.New("Required option " + entry.LongName + " not found")
			break
		}
	}

	// If the parse went okay, let's check to make sure we don't have dangling
	// parameters, and then call the action if there is one.

	if err == nil {

		if expectedParameters == 0 && len(Parameters) > 0 {
			return errors.New("unexpected parameters on command line")
		}
		if expectedParameters < 0 {
			if len(Parameters) > -expectedParameters {
				return errors.New("too many parameters on command line")
			}
		} else {
			if len(Parameters) != expectedParameters {
				return errors.New("incorrect number of parameters on command line")
			}
		}
		// Did we ever find an action routine? If so, let's run it. Otherwise,
		// there wasn't enough command to determine what to do, so show the help.
		if Action != nil {
			err = Action(c)
		} else {
			ShowHelp(c)
		}
	}
	return err
}
