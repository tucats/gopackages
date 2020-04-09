// Package main contains the main program, which is used as a test driver to validate features
// as they are added to the app-cli package. This is not intended to be run as a standalone
// program, as it has uninterestly-limited features.
package main

import (
	"github.com/tucats/gopackages/cli/app"
	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/commands"
)

func main() {

	cli.SetCopyright("(c) 2020 Tom Cole, fernwood.org")
	app.Run(commands.Grammar, "app-cli", "test driver for CLI development")

}
