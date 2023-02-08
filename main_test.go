package main

import (
	"fmt"
	"testing"

	"github.com/tucats/gopackages/app-cli/app"
	"github.com/tucats/gopackages/app-cli/cli"
)

var grammar = []cli.Option{
	{
		LongName:   "test",
		OptionType: cli.BooleanType,
		Action:     setTest,
	},
}

func setTest(c *cli.Context) error {
	fmt.Println("--test activated")

	return nil
}

func TestMain(t *testing.T) {
	app := app.New("test driver")

	args := []string{"driver", "--test"}
	err := app.Run(grammar, args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
