package app

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/profile"
)

// Run sets up required data structures and executes the parse. The caller's
// grammar is extended with pre-defined global verbs and options (such as
// profile or --output-format) common to all CLI applications.
func Run(grammar cli.Options, appName string, appDescription string) error {

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
			Action:      SetDefaultProfile,
		},
		cli.Option{
			ShortName:   "d",
			LongName:    "debug",
			Description: "Are we debugging?",
			OptionType:  cli.BooleanType,
			Action:      SetDebugMessaging,
		},
		cli.Option{
			LongName:    "output-format",
			Description: "Specify text or json output format",
			OptionType:  cli.StringType,
			Action:      SetOutputFormat,
		},
		cli.Option{
			ShortName:   "q",
			LongName:    "quiet",
			Description: "If specified, suppress extra messaging",
			OptionType:  cli.BooleanType,
			Action:      SetQuietMode,
		}}, grammar...)

	// Load the active profile, if any
	profile.Load(appName, "default")

	// Parse the grammar and call the actions (essentially, execute
	// the function of the CLI)
	err := cli.Parse(grammar, appDescription)

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
