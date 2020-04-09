package main

import (
	"app-cli/cli"
	"app-cli/commands"
	"fmt"
	"os"
)

func main() {

	cli.SetCopyright("(c) 2020 Tom Cole, fernwood.org")
	status := cli.Parse(commands.Grammar, "test driver for CLI package")

	if status != nil {
		fmt.Printf("Error, %s\n", status.Error())
		os.Exit(1)
	}

}
