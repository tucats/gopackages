package cli

import (
	"errors"

	"github.com/tucats/gopackages/cli/profile"
	"github.com/tucats/gopackages/cli/tables"
	"github.com/tucats/gopackages/cli/ui"
)

// ProfileGrammar contains profile subcommands
var ProfileGrammar = Options{
	Option{
		LongName:    "show",
		Description: "Show the current profile",
		Action:      ProfileShow,
		OptionType:  Subcommand,
	},
	Option{
		LongName:    "set",
		Description: "Set a profile value",
		Action:      ProfileSet,
		OptionType:  Subcommand,
	},
}

// ProfileShow Displays the current contents of the active profile
func ProfileShow(c *Options) error {

	t := tables.New([]string{"Key", "Value"})

	for k, v := range profile.CurrentConfiguration.Items {
		t.AddRowItems(k, v)
	}
	t.SetOrderBy("key")
	t.Print(ui.DefaultTableFormat)

	return nil
}

// ProfileSet uses the first two parameters as a key and value
func ProfileSet(c *Options) error {

	if len(Parameters) != 2 {
		return errors.New("Missing key and value")
	}

	profile.Set(Parameters[0], Parameters[1])
	return nil
}
