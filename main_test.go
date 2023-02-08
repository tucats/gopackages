package main

import (
	"fmt"
	"testing"

	"github.com/tucats/gopackages/app-cli/app"
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/tables"
)

var grammar = []cli.Option{
	{
		LongName:   "test",
		OptionType: cli.StringType,
		Action:     setTest,
	},
}

func setTest(c *cli.Context) error {
	fmt.Println("--test activated")

	return nil
}

func defaultAction(c *cli.Context) error {
	fmt.Println("In default action")
	if c.WasFound("test") {
		fmt.Println(c.String("test"))
	}

	f, _ := tables.New([]string{"Name", "Age"})

	_ = f.AddRowItems("Tom", 63)
	_ = f.AddRowItems("Mary", 58)
	_ = f.AddRowItems("David", 29)
	_ = f.AddRowItems("Bonnie", 73)

	f.SetWhere("age < 60")

	return f.Print("text")
}

func TestMain(t *testing.T) {
	app := app.New("test driver").SetDefaultAction(defaultAction)

	args := []string{"driver", "--test"}
	err := app.Run(grammar, args[1:])
	if err != nil {
		fmt.Println(err)
	}
}
