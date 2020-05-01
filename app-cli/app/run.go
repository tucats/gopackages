// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine.
package app

import (
	"os"
	"strconv"
	"strings"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/profile"
	"github.com/tucats/gopackages/app-cli/ui"
)

// Run sets up required data structures and executes the parse.
// When completed, the command line functionality will have been
// run. It is up to the caller (typically the main() function)
// to handle any post-processing cleanup, etc.
//
// The grammar is cli.Options array of cli.Option structures.
// Each element describes a parsable token in the command line grammar.
// This grammar is extended to include the automatic built-in
// commands for profile management, etc.
//
// The args array is the array of command line arguments to be parsed.
// This could be derived from the os.Args array, or a generated command
// line argument (or one read in from an external source).
//
// The appName is the name of the CLI application, optionally followed
// by a colon character (":") and the text description of the application.
// These are primarily used in presenting --help information.
//
func Run(grammar []cli.Option, args []string, appName string) error {

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
			ShortName:           "p",
			LongName:            "profile",
			Description:         "Name of profile to use",
			OptionType:          cli.StringType,
			Action:              UseProfileAction,
			EnvironmentVariable: "CLI_PROFILE",
		},
		cli.Option{
			ShortName:           "d",
			LongName:            "debug",
			Description:         "Are we debugging?",
			OptionType:          cli.BooleanType,
			Action:              DebugAction,
			EnvironmentVariable: "CLI_DEBUG",
		},
		cli.Option{
			LongName:            "output-format",
			Description:         "Specify text or json output format",
			OptionType:          cli.StringType,
			Action:              OutputFormatAction,
			EnvironmentVariable: "CLI_OUTPUT_FORMAT",
		},
		cli.Option{
			ShortName:   "v",
			LongName:    "version",
			Description: "Show version number of command line tool",
			OptionType:  cli.BooleanType,
			Action:      ShowVersionAction,
		},
		cli.Option{
			ShortName:           "q",
			LongName:            "quiet",
			Description:         "If specified, suppress extra messaging",
			OptionType:          cli.BooleanType,
			Action:              QuietAction,
			EnvironmentVariable: "CLI_QUIET",
		}}, grammar...)

	// Extract the description of the app if it was given

	var appDescription = ""
	if i := strings.Index(appName, ":"); i > 0 {
		appDescription = strings.TrimSpace(appName[i+1:])
		appName = strings.TrimSpace(appName[:i])
	}

	// Load the active profile, if any
	profile.Load(appName, "default")

	// If the CLI_DEBUG environment variable is set, then turn on
	// debugging now, so messages will come out before that particular
	// option is processed.

	if os.Getenv("CLI_DEBUG") != "" {
		ui.DebugMode = true
	}

	// Parse the grammar and call the actions (essentially, execute
	// the function of the CLI)
	context := cli.Context{Description: appDescription, Grammar: grammar, Args: args}
	err := context.Parse(appDescription)

	// If no errors, then write out an updated profile as needed.
	if err == nil {
		err = profile.Save()
	}

	return err
}

// SetCopyright sets the copyright string used in the help output.
func SetCopyright(copyright string) {
	cli.SetCopyright(copyright)
}

// SetVersion sets the version string for the application
func SetVersion(version []int) {
	var v strings.Builder

	v.WriteString("v")
	for i, n := range version {

		if i > 1 && n == 0 {
			break
		}
		if i > 1 {
			v.WriteRune('-')
		} else if i > 0 {
			v.WriteRune('.')
		}
		v.WriteString(strconv.Itoa(n))
	}
	cli.Version = v.String()
}
