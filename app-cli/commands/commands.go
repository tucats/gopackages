package commands

import "github.com/tucats/gopackages/cli/cli"

// Grammar is the primary grammar of the command line, which defines all global options
// and any subcommands.
var Grammar cli.Options = cli.Options{
	cli.Option{
		LongName:    "list",
		Description: "List stuff",
		OptionType:  cli.Subcommand,
		Value:       ListGrammar,
		Action:      ListActions,
	},
}
