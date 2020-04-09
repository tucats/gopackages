package commands

import (
	"app-cli/cli"
	"app-cli/tables"
	"app-cli/ui"
)

// ListGrammar is the grammar definition for the list command
var ListGrammar cli.Options = []cli.Option{
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

// ListActions is the command handler to list objects.
func ListActions(c *cli.Options) error {

	t := tables.New([]string{"Name", "Age"})
	_ = t.SetAlignment(1, tables.AlignmentRight)

	// If an order-by is given, tell the table to sort the data
	if name, present := cli.GetString(*c, "order-by"); present {
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
	t.SuppressHeadings(cli.GetBool(*c, "no-headings"))
	t.RowNumbers(cli.GetBool(*c, "row-numbers"))
	if limit, present := cli.GetInteger(*c, "limit"); present {
		t.RowLimit(limit)
	}
	if startingRow, present := cli.GetInteger(*c, "start"); present {
		if err := t.SetStartingRow(startingRow); err != nil {
			return err
		}
	}

	// Print the table in the user-requested format.
	return t.Print(ui.OutputFormat)

}
