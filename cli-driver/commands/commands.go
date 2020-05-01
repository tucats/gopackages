// Package commands contains the grammar definitions for all commands, and can
// optionally contain the implementations of those commands. In this example,
// the actions are stored in separate source files.
package commands

import "github.com/tucats/gopackages/app-cli/cli"

// Grammar is the primary grammar of the command line, which defines all global options
// and any subcommands.
var Grammar = []cli.Option{
	cli.Option{
		LongName:    "list",
		Description: "List contents of a table",
		OptionType:  cli.Subcommand,
		Value:       ListGrammar,
		Action:      ListAction,
	},
	cli.Option{
		LongName:    "weather",
		Description: "Get the current weather",
		OptionType:  cli.Subcommand,
		Value:       WeatherGrammar,
		Action:      WeatherAction,
	},
}
