package cli

import (
	"fmt"

	"github.com/tucats/gopackages/app-cli/tables"
	"github.com/tucats/gopackages/app-cli/ui"
)

// ShowHelp displays help text for the grammar, using a standardized format.
// The help shows subcommands as well as options, including value type cues.
// The output is automatically directed to the stdout console output.
//
// This function uses the tables package to create uniform columns of output.
func ShowHelp(c *Context) {

	if c.Copyright != "" {
		fmt.Printf("%s\n", c.Copyright)
	}

	composedCommand := c.MainProgram + " " + c.Command

	hasSubcommand := false
	hasOptions := false

	for _, option := range c.Grammar {
		if option.OptionType == Subcommand {
			hasSubcommand = true
		} else {
			hasOptions = true
		}
	}
	if hasOptions {
		composedCommand = composedCommand + "[options] "
	}
	if hasSubcommand {
		composedCommand = composedCommand + "[command] "
	}

	g := c.FindGlobal()
	e := g.ExpectedParameterCount
	if g.ParameterDescription > "" {
		composedCommand = composedCommand + " [" + g.ParameterDescription + "]"
	} else if e == 1 {
		composedCommand = composedCommand + " [parameter]"
	} else if e > 1 {
		composedCommand = composedCommand + " [parameters]"
	}

	minimumFirstColumnWidth := len(composedCommand)
	if minimumFirstColumnWidth < 26 {
		minimumFirstColumnWidth = 26
	}
	if c.Parent == nil && c.Version != "" {
		c.Description = c.Description + ", " + c.Version
	}

	fmt.Printf("\nUsage:\n   %-26s   %s\n\n", composedCommand, c.Description)
	headerShown := false

	tc, _ := tables.New([]string{"subcommand", "description"})
	tc.ShowHeadings(false)
	tc.SetIndent(3)
	tc.SetSpacing(3)
	tc.SetMinimumWidth(0, minimumFirstColumnWidth)

	for _, option := range c.Grammar {
		if option.OptionType == Subcommand && !option.Private {
			if !headerShown {
				headerShown = true
				fmt.Printf("Commands:\n")
				tc.AddRow([]string{"help", "Display help text"})
			}
			tc.AddRow([]string{option.LongName, option.Description})
		}
	}
	if headerShown {
		tc.SortRows(0, true)
		tc.Print(ui.TextTableFormat)
		fmt.Printf("\n")
	}

	headerShown = false
	tc, _ = tables.New([]string{"Parameter"})
	tc.ShowHeadings(false)
	tc.SetIndent(3)
	tc.SetMinimumWidth(0, minimumFirstColumnWidth)
	for _, option := range c.Grammar {
		if option.OptionType == ParameterType {
			if !headerShown {
				fmt.Printf("Parameters:\n")
				headerShow = true
				tc.AddRowItems(option.Description)
			}

		}
	}
	if headerShown {
		tc.Print("text")
	}

	to, _ := tables.New([]string{"option", "description"})
	to.ShowHeadings(false)
	to.SetIndent(3)
	to.SetSpacing(3)
	to.SetMinimumWidth(0, minimumFirstColumnWidth)

	for _, option := range c.Grammar {
		if option.Private {
			continue
		}
		if option.OptionType != Subcommand {

			name := ""
			if option.LongName > "" {
				name = "--" + option.LongName
			}
			if option.ShortName > "" {
				if name > "" {
					name = name + ", "
				}
				name = name + "-" + option.ShortName
			}
			switch option.OptionType {
			case IntType:
				name = name + " <integer>"

			case StringType:
				name = name + " <string>"

			case BooleanValueType:
				name = name + " <boolean>"

			case StringListType:
				name = name + " <list>"
			}

			fullDescription := option.Description
			if option.EnvironmentVariable != "" {
				fullDescription = fullDescription + " [" + option.EnvironmentVariable + "]"
			}
			to.AddRow([]string{name, fullDescription})
		}
	}

	fmt.Printf("Options:\n")
	to.AddRow([]string{"--help, -h", "Show this help text"})
	to.SortRows(0, true)
	to.Print(ui.TextTableFormat)
}
