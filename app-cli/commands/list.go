package commands

import (
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/profile"
	"github.com/tucats/gopackages/app-cli/tables"
	"github.com/tucats/gopackages/app-cli/ui"
)

// ListGrammar is the grammar definition for the list command
var ListGrammar = []cli.Option{
	cli.Option{
		LongName:    "no-headings",
		Description: "If specified, do not print headings",
		OptionType:  cli.BooleanType,
	},
	cli.Option{
		LongName:    "row-numbers",
		Description: "If specified, print a column with the row number",
		OptionType:  cli.BooleanType,
	},
	cli.Option{
		LongName:    "start",
		Description: "Specify the row number to start the list",
		OptionType:  cli.IntType,
	},
	cli.Option{
		LongName:    "limit",
		Description: "Specify the number of items to list",
		OptionType:  cli.IntType,
	},
	cli.Option{
		LongName:    "columns",
		ShortName:   "c",
		OptionType:  cli.StringListType,
		Description: "Specify the columns you wish listed",
		Private:     true,
	},
	cli.Option{
		LongName:    "order-by",
		OptionType:  cli.StringType,
		Description: "Specify the column to use to sort the output",
	},
}

// ListAction is the command handler to list objects.
func ListAction(c *cli.Context) error {

	ui.Debug("In the LIST action")

	t := tables.New([]string{"Name", "Age"})
	_ = t.SetAlignment(1, tables.AlignmentRight)

	// If an order-by is given, tell the table to sort the data
	if name, present := c.GetString("order-by"); present {
		if err := t.SetOrderBy(name); err != nil {
			return err
		}
	}

	// Add the rows to the table representing the information to be printed out
	t.AddRow([]string{"Tom", "60"})
	t.AddRow([]string{"Mary", "59"})
	t.AddRow([]string{"Sarah", "25"})
	t.AddRow([]string{"Chelsea", "27"})
	t.AddRowItems("Anna", 25)

	// Add formatting and other control settings the user might specify.
	t.ShowHeadings(!c.GetBool("no-headings"))
	t.ShowRowNumbers(c.GetBool("row-numbers"))
	if limit, present := c.GetInteger("limit"); present {
		t.RowLimit(limit)
	}
	if startingRow, present := c.GetInteger("start"); present {
		if err := t.SetStartingRow(startingRow); err != nil {
			return err
		}
	}

	// Print the table in the user-requested format.
	return t.Print(profile.Get("output-format"))

}
