package cli

import (
	"fmt"
	"os"

	"github.com/tucats/gopackages/app-cli/ui"
	"gitlab.com/tucats/gopackages/cli/tables"
)

// Copyright contains the copyright string (if any) used in help output
var Copyright string

// SetCopyright sets the copyright string. If not set, then no copyright
// message is displayed as part of help.
func SetCopyright(s string) {
	Copyright = s
}

// ShowHelp displays help text for the grammar
func ShowHelp(grammar Options) {

	if Copyright != "" {
		name := MainProgram
		if MainProgramDescription > "" {
			name = name + " - " + MainProgramDescription
		}
		fmt.Printf("%s - %s\n\n", name, Copyright)
	}

	composedCommand := MainProgram + " " + CommandRoot

	hasSubcommand := false
	hasOptions := false

	for _, option := range grammar {
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

	minimumFirstColumnWidth := len(composedCommand)
	if minimumFirstColumnWidth < 20 {
		minimumFirstColumnWidth = 20
	}

	fmt.Printf("Usage:\n   %-20s   %s\n\n", composedCommand, CurrentVerbDescription)
	headerShown := false

	t := tables.New([]string{"subcommand", "description"})
	t.SuppressHeadings(true)
	t.SetIndent(3)
	t.SetSpacing(3)
	t.SetMinimumWidth(0, minimumFirstColumnWidth)

	for _, option := range grammar {
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

	for n, option := range grammar {
		if n == 0 {
			fmt.Printf("Options:\n")
		}
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
