// Package commands contains the grammar definitioin for all commands, and can
// also contain the implementations of those commands.
package commands

import "github.com/tucats/gopackages/app-cli/cli"

// Grammar is the primary grammar of the command line, which defines all global options
// and any subcommands.
var Grammar = []cli.Option{
	cli.Option{
		LongName:    "list",
		Description: "Demonstration command to list a table",
		OptionType:  cli.Subcommand,
		Value:       ListGrammar,
		Action:      ListAction,
	},
}
