// Package main contains the main program, which is used as a test driver to validate features
// as they are added to the app-cli package. This is not intended to be run as a useful CLI
// program, as it has uninteresting and limited features.
package main

import (
	"github.com/tucats/gopackages/app-cli/app"
	"github.com/tucats/gopackages/cli-driver/commands"
)

func main() {

	app.SetCopyright("(c) 2020 Tom Cole. All rights reserved.")
	app.SetVersion([]int{1, 1, 35})
	app.Run(commands.Grammar, "cli-driver", "test driver for CLI development")

}
