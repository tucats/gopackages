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

// MainProgram is the name of the main program extracted from arguments.
var MainProgram string

// CurrentVerbDescription is a descriptive string used in help output. It
// is the description of the main program, or the current active verb.
var CurrentVerbDescription string

// MainProgramDescription is the description of the main program
var MainProgramDescription string

// CommandRoot builds out the command components as we parse, for use in
// help output. So the command verb and the chain of subcommands are stored
// here, separated by blanks.
var CommandRoot string

// Parameters stores those command line items not parsed as part of the
// grammar.
var Parameters []string

// Action is the function to invoke on the last subcommand found in parsing.
var Action func(grammar *Options) error

// Globals is a copy of the outermost grammar definition, and is used to find
// global values later.
var Globals *[]Option

// Parse accepts a grammar definition and parses the current argument
// list against that grammar. Unrecognized options or subcommands, as
// well as invalid values are reported as an error. If there is an
// action routine associated with an option or a subcommand, that
// action is executed.
func Parse(grammar []Option, description string) error {

	args := os.Args
	MainProgram = filepath.Base(args[0])
	CurrentVerbDescription = ""
	MainProgramDescription = description
	CommandRoot = ""
	Action = nil

	// Prepend the default supplied options
	grammar = append([]Option{
		Option{
			LongName:    "profile",
			OptionType:  Subcommand,
			Description: "Manage the default profile",
			Value:       ProfileGrammar,
		},
		Option{
			ShortName:   "p",
			LongName:    "use-profile",
			Description: "Name of profile to use",
			OptionType:  StringType,
			Action:      SetDefaultProfile,
		},
		Option{
			ShortName:   "d",
			LongName:    "debug",
			Description: "Are we debugging?",
			OptionType:  BooleanType,
			Action:      SetDebugMessaging,
		},
		Option{
			LongName:    "output-format",
			Description: "Specify text or json output format",
			OptionType:  StringType,
			Action:      SetOutputFormat,
		},
		Option{
			ShortName:   "q",
			LongName:    "quiet",
			Description: "If specified, suppress extra messaging",
			OptionType:  BooleanType,
			Action:      SetQuietMode,
		}}, grammar...)

	Globals = &grammar

	// If there are no arguments other than the main program name, dump out the help by default.
	if len(args) == 1 {
		ShowHelp(grammar)
	}

	// Start parsing using the top-level grammar.
	return ParseGrammar(args[1:], grammar)
}

// ParseGrammar accepts an argument list and a grammar definition, and parses
func ParseGrammar(args []string, grammar Options) error {

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
			ShowHelp(grammar)
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
			for n, entry := range grammar {

				if (isShort && entry.ShortName == name) || (!isShort && entry.LongName == name) {
					location = &grammar[n]
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
			for _, entry := range grammar {
				if entry.LongName == option && entry.OptionType == Subcommand {

					subGrammar := Options{}
					if entry.Value != nil {
						subGrammar = entry.Value.(Options)
					}

					CommandRoot = CommandRoot + entry.LongName + " "
					CurrentVerbDescription = entry.Description
					entry.Found = true
					if entry.Action != nil {
						Action = entry.Action
						ui.Debug("Adding action routine")
					}
					ui.Debug("Transferring control to subgrammar for %s", entry.LongName)
					return ParseGrammar(args[currentArg+1:], subGrammar)
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
				err = location.Action(&grammar)
				if err != nil {
					break
				}
			}
		}
	}

	// Whew! Everything parsed and in it's place. Before we wind up, let's verify that
	// all required options were in fact found.

	for _, entry := range grammar {

		if entry.Required && !entry.Found {
			err = errors.New("Required option " + entry.LongName + " not found")
			break
		}
	}

	// Did we ever find an action routine? If so, let's run it.
	if err == nil && Action != nil {
		err = Action(&grammar)
	} else {
		ShowHelp(grammar)
	}
	return err
}
