// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine.
package app

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/profile"
)

// Run sets up required data structures and executes the parse.
// When completed, the command line functionality will have been
// run. It is up to the caller (typically the main() fucntion)
// to handle any post-processing cleanup, etc.
//
// * The grammar is cli.Options array of cli.Option structures.
//   Each element describes a parsable token in the command line grammar.
//   This grammar is extended to include the automatic built-in
//   commands for profile management, etc.
//
// * The appName is the name of the CLI application.
//   This is used in --help output and in determining the name of
//   the configuration file used.
//
// * The appDescription is used as the text description of the application.
//   This is displayed in --help output
//
func Run(grammar []cli.Option, appName string, appDescription string) error {

	// Prepend the default supplied options
	grammar = append([]cli.Option{
		cli.Option{
			LongName:    "profile",
			Aliases:     []string{"prof"},
			OptionType:  cli.Subcommand,
			Description: "Manage the default profile",
			Value:       profile.Grammar,
		},
		cli.Option{
			ShortName:   "p",
			LongName:    "use-profile",
			Description: "Name of profile to use",
			OptionType:  cli.StringType,
			Action:      UseProfileAction,
		},
		cli.Option{
			ShortName:   "d",
			LongName:    "debug",
			Description: "Are we debugging?",
			OptionType:  cli.BooleanType,
			Action:      DebugAction,
		},
		cli.Option{
			LongName:    "output-format",
			Description: "Specify text or json output format",
			OptionType:  cli.StringType,
			Action:      OutputFormatAction,
		},
		cli.Option{
			ShortName:   "q",
			LongName:    "quiet",
			Description: "If specified, suppress extra messaging",
			OptionType:  cli.BooleanType,
			Action:      QuietAction,
		}}, grammar...)

	// Load the active profile, if any
	profile.Load(appName, "default")

	// Parse the grammar and call the actions (essentially, execute
	// the function of the CLI)
	context := cli.Context{Grammar: grammar}
	err := context.Parse(appDescription)

	// If no errors, then write out an updated profile as needed.
	if err == nil {
		err = profile.Save()
	}

	// If something went wrong, report it to the user and force an exit
	// status of 1. @TOMCOLE later this should be extended to allow an error
	// code to carry along the desired exit code to support multiple types
	// of errors.
	if err != nil {
		fmt.Printf("Error. %s\n", err.Error())
		os.Exit(1)
	}
	return err
}

// SetCopyright sets the copyright string used in the help output.
func SetCopyright(copyright string) {
	cli.SetCopyright("(c) 2020 Tom Cole, fernwood.org")
}
