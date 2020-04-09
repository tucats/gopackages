package cli

import (
	"strings"

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
		Value: Options{
			Option{
				LongName:    "key",
				Description: "The key that will be set in the profile. Can be of the form key=value.",
				OptionType:  StringType,
				Required:    true,
			},
			Option{
				LongName:    "value",
				Description: "The value to set for the key. If missing, the key is deleted",
				OptionType:  StringType,
			},
		},
	},
}

// ProfileShow Displays the current contents of the active profile
func ProfileShow(c *Options) error {

	t := tables.New([]string{"Key", "Value"})

	for k, v := range profile.CurrentConfiguration.Items {
		t.AddRowItems(k, v)
	}
	t.SetOrderBy("key")
	t.Underlines(false)
	t.Print(ui.TextTableFormat)

	return nil
}

// ProfileSet uses the first two parameters as a key and value
func ProfileSet(c *Options) error {

	key, _ := GetString(*c, "key")
	value, valueFound := GetString(*c, "value")

	if !valueFound {
		if equals := strings.Index(key, "="); equals >= 0 {
			value = key[equals+1:]
			key = key[:equals]
			valueFound = true
		}
	}

	if valueFound {
		profile.Set(key, value)
		ui.Say("Profile key %s written", key)
	} else {
		profile.Delete(key)
		ui.Say("Profile key %s deleted", key)
	}

	return nil
}
