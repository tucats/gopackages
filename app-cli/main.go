// Package main contains the main program, which is used as a test driver to validate features
// as they are added to the app-cli package. This is not intended to be run as a standalone
// program, as it has uninterestly-limited features.
package main

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/commands"
	"github.com/tucats/gopackages/cli/profile"
)

func main() {

	profile.Load("default")

	cli.SetCopyright("(c) 2020 Tom Cole, fernwood.org")
	status := cli.Parse(commands.Grammar, "test driver for CLI package")

	if status != nil {
		fmt.Printf("Error, %s\n", status.Error())
		os.Exit(1)
	}

	profile.Save()

}
