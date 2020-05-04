package commands

import (
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/profile"
	"github.com/tucats/gopackages/app-cli/tables"
	"github.com/tucats/gopackages/app-cli/ui"
)

// ListGrammar is the grammar definition for the list command. It
// defines each of the command line options, the option type and
// value type if appropriate. There are no actions defined in this
// grammar, as the action was defined in the parent grammer for the
// subcommand itself in the parent grammar.
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

// ListAction is the command handler to list objects. It uses option
// values parsed by the app-cli framework to determine how to display
// a short table of values. The action creates the table and then
// uses command line options to set the characteristics of the table
// output in the Table object.
func ListAction(c *cli.Context) error {

	ui.Debug("In the LIST action")

	t, _ := tables.New([]string{"Name", "Age"})

	// Add the rows to the table representing the information to be printed out
	t.AddRow([]string{"Tom", "60"})
	t.AddRow([]string{"Mary", "59"})
	t.AddRow([]string{"Sarah", "25"})
	t.AddRow([]string{"Chelsea", "27"})
	t.AddRowItems("Anna", 25)
	t.AddCSVRow("Claire,17")

	// Add formatting and other control settings the user might specify.
	_ = t.SetAlignment(1, tables.AlignmentRight)

	t.ShowHeadings(!c.GetBool("no-headings"))
	t.ShowRowNumbers(c.GetBool("row-numbers"))

	if name, present := c.GetString("order-by"); present {
		if err := t.SetOrderBy(name); err != nil {
			return err
		}
	}

	if startingRow, present := c.GetInteger("start"); present {
		if err := t.SetStartingRow(startingRow); err != nil {
			return err
		}
	}

	if limit, present := c.GetInteger("limit"); present {
		t.RowLimit(limit)
	}

	// Print the table in the user-requested format.
	return t.Print(profile.Get("output-format"))

}
