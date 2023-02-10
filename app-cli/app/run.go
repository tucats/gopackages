package app

import (
	"os"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/settings"
)

// Run sets up required data structures and parses the command line. It then
// automatically calls any action routines specfied in the grammar, which do
// the work of the command.
func runFromContext(context *cli.Context) error {

	// Add the user-provided grammar.
	applicationGrammar = append(applicationGrammar, context.Grammar...)

	// Load the active profile, if any from the profile for this application.
	_ = settings.Load(context.AppName, "default")

	context.Grammar = applicationGrammar

	// If we are to dump the grammar (a diagnostic function) do that,
	// then just pack it in and go home.
	if os.Getenv("APP_DUMP_GRAMMAR") != "" {
		cli.DumpGrammar(context)
		os.Exit(0)
	}

	// Parse the grammar and call the actions (essentially, execute
	// the function of the CLI). If it goes poorly, error out.
	if err := context.Parse(); err != nil {
		return err
	} else {
		// If no errors, then write out an updated profile as needed.
		err = settings.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
