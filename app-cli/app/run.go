package app

import (
	"os"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/app-cli/profile"
	"github.com/tucats/gopackages/app-cli/ui"
)

// Run sets up required data structures and executes the parse.
// When completed, the command line functionality will have been
// run. It is up to the caller (typically the main() function)
// to handle any post-processing cleanup, etc.
func runFromContext(context *cli.Context) error {

	// Create a new grammar which prepends the default supplied options
	// to the caller's grammar definition.
	grammar := []cli.Option{
		{
			LongName:    "profile",
			Aliases:     []string{"prof"},
			OptionType:  cli.Subcommand,
			Description: "Manage the default profile",
			Value:       profile.Grammar,
		},
		{
			LongName:    "logon",
			Aliases:     []string{"login"},
			OptionType:  cli.Subcommand,
			Description: "Log on to a remote server",
			Action:      Logon,
			Value:       LogonGrammar,
		},
		{
			ShortName:           "p",
			LongName:            "profile",
			Description:         "Name of profile to use",
			OptionType:          cli.StringType,
			Action:              UseProfileAction,
			EnvironmentVariable: "CLI_PROFILE",
		},
		{
			ShortName:   "d",
			LongName:    "debug",
			Description: "Debug loggers to enable",
			OptionType:  cli.StringListType,
			Action:      DebugAction,
		},
		{
			LongName:            "output-format",
			Description:         "Specify text or json output format",
			OptionType:          cli.StringType,
			Action:              OutputFormatAction,
			EnvironmentVariable: "CLI_OUTPUT_FORMAT",
		},
		{
			ShortName:   "v",
			LongName:    "version",
			Description: "Show version number of command line tool",
			OptionType:  cli.BooleanType,
			Action:      ShowVersionAction,
		},
		{
			ShortName:           "q",
			LongName:            "quiet",
			Description:         "If specified, suppress extra messaging",
			OptionType:          cli.BooleanType,
			Action:              QuietAction,
			EnvironmentVariable: "CLI_QUIET",
		},
	}

	// Add the user-provided grammar
	grammar = append(grammar, context.Grammar...)

	// Load the active profile, if any from the profile for this application.
	_ = persistence.Load(context.AppName, "default")

	// If the CLI_DEBUG environment variable is set, then turn on
	// debugging now, so messages will come out before that particular
	// option is processed.
	ui.SetLogger(ui.DebugLogger, os.Getenv("CLI_DEBUG") != "")

	// Parse the grammar and call the actions (essentially, execute
	// the function of the CLI)
	context.Grammar = grammar
	err := context.Parse()

	// If no errors, then write out an updated profile as needed.
	if err == nil {
		err = persistence.Save()
	}

	return err
}
