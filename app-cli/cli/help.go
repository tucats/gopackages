package cli

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/cli/tables"
	"github.com/tucats/gopackages/cli/ui"
)

// Copyright contains the copyright string (if any) used in help output
var Copyright string

// SetCopyright sets the copyright string. If not set, then no copyright
// message is displayed as part of help.
func SetCopyright(s string) {
	Copyright = s
}

// ShowHelp displays help text for the grammar, using a standardized format.
// The help shows subcommands as well as options, including value type cues.
// The output is automatically directed to the stdout console output.
//
// This function uses the tables package to create uniform columns of output.
func ShowHelp(c *Context) {

	if Copyright != "" {
		name := c.MainProgram
		if MainProgramDescription > "" {
			name = name + " - " + MainProgramDescription
		}
		fmt.Printf("%s - %s\n\n", name, Copyright)
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
	if minimumFirstColumnWidth < 20 {
		minimumFirstColumnWidth = 20
	}

	fmt.Printf("Usage:\n   %-20s   %s\n\n", composedCommand, c.Description)
	headerShown := false

	t := tables.New([]string{"subcommand", "description"})
	t.SuppressHeadings(true)
	t.SetIndent(3)
	t.SetSpacing(3)
	t.SetMinimumWidth(0, minimumFirstColumnWidth)

	for _, option := range c.Grammar {
		if option.OptionType == Subcommand && !option.Private {
			if !headerShown {
				headerShown = true
				fmt.Printf("Commands:\n")
				t.AddRow([]string{"help", "Display help text"})
			}
			t.AddRow([]string{option.LongName, option.Description})
		}
	}
	if headerShown {
		t.SortRows(0, true)
		t.Print(ui.TextTableFormat)
		fmt.Printf("\n")
	}

	t = tables.New([]string{"option", "description"})
	t.SuppressHeadings(true)
	t.SetIndent(3)
	t.SetSpacing(3)
	t.SetMinimumWidth(0, minimumFirstColumnWidth)

	fmt.Printf("Options:\n")

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

			t.AddRow([]string{name, option.Description})
		}
	}

	t.AddRow([]string{"--help, -h", "Show this help text"})
	t.SortRows(0, true)
	t.Print(ui.TextTableFormat)

	os.Exit(0)
}
