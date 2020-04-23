// Package main contains the main program, which is used as a test driver to validate features
// as they are added to the app-cli package. This is not intended to be run as a useful CLI
// program, as it has uninteresting and limited features.
package main

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/app-cli/app"
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/cli-driver/commands"
)

func main() {

	app.SetCopyright("(c) 2020 Tom Cole. All rights reserved.")
	app.SetVersion([]int{1, 1, 1})
	err := app.Run(commands.Grammar, os.Args, "cli-driver: test driver for CLI development")

	// If something went wrong, report it to the user and force an exit
	// status of 1. @TOMCOLE later this should be extended to allow an error
	// code to carry along the desired exit code to support multiple types
	// of errors.
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		if e2, ok := err.(cli.ExitError); ok {
			os.Exit(e2.ExitStatus)
		}
		os.Exit(1)
	}
}
