package profile

import (
	"errors"
	"strings"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/tables"
	"github.com/tucats/gopackages/cli/ui"
)

// Grammar describes profile subcommands
var Grammar = []cli.Option{
	cli.Option{
		LongName:    "list",
		Description: "List all profiles",
		Action:      ListAction,
		OptionType:  cli.Subcommand,
	},
	cli.Option{
		LongName:    "show",
		Description: "Show the current profile",
		Action:      ShowAction,
		OptionType:  cli.Subcommand,
	},
	cli.Option{
		LongName:             "set-output",
		OptionType:           cli.Subcommand,
		Description:          "Set the default output type (text or json)",
		ParameterDescription: "type",
		Action:               SetOutputAction,
		Parameters:           1,
	},
	cli.Option{
		LongName:             "set-description",
		OptionType:           cli.Subcommand,
		Description:          "Set the profile description",
		ParameterDescription: "text",
		Parameters:           1,
		Action:               SetDescriptionAction,
	},
	cli.Option{
		LongName:             "delete",
		OptionType:           cli.Subcommand,
		Description:          "Delete a key from the profile",
		Action:               DeleteAction,
		Parameters:           1,
		ParameterDescription: "key",
	},
	cli.Option{
		LongName:    "set",
		Description: "Set a profile value",
		Action:      SetAction,
		OptionType:  cli.Subcommand,
		Value: []cli.Option{
			cli.Option{
				LongName:    "key",
				Description: "The key that will be set in the profile. Can be of the form key=value.",
				OptionType:  cli.StringType,
				Required:    true,
			},
			cli.Option{
				LongName:    "value",
				Description: "The value to set for the key. If missing, the key is deleted",
				OptionType:  cli.StringType,
			},
		},
	},
}

// ShowAction Displays the current contents of the active profile
func ShowAction(c *cli.Context) error {

	t := tables.New([]string{"Key", "Value"})

	for k, v := range CurrentConfiguration.Items {
		t.AddRowItems(k, v)
	}
	t.SetOrderBy("key")
	t.Underlines(false)
	t.Print(ui.TextTableFormat)

	return nil
}

// ListAction Displays the current contents of the active profile
func ListAction(c *cli.Context) error {

	t := tables.New([]string{"Name", "Description"})

	for k, v := range Configurations {
		t.AddRowItems(k, v.Description)
	}
	t.SetOrderBy("name")
	t.Underlines(false)
	t.Print(ui.TextTableFormat)

	return nil
}

// SetOutputAction is the action handler for the set-output subcommand.
func SetOutputAction(c *cli.Context) error {

	if c.GetParameterCount() == 1 {
		outputType := c.GetParameter(0)
		if outputType == "text" || outputType == "json" {
			Set("output-format", outputType)
			return nil
		}
		return errors.New("Invalid output type: " + outputType)
	}
	return errors.New("Missing output type")
}

// SetAction uses the first two parameters as a key and value
func SetAction(c *cli.Context) error {

	// Generic --key and --value specification
	key, _ := c.GetString("key")
	value, valueFound := c.GetString("value")

	if !valueFound {
		if equals := strings.Index(key, "="); equals >= 0 {
			value = key[equals+1:]
			key = key[:equals]
			valueFound = true
		}
	}

	if valueFound {
		Set(key, value)
		ui.Say("Profile key %s written", key)
	} else {
		Delete(key)
		ui.Say("Profile key %s deleted", key)
	}

	return nil
}

// DeleteAction deletes a named key value
func DeleteAction(c *cli.Context) error {

	key := c.GetParameter(0)
	Delete(key)
	ui.Say("Profile key %s deleted", key)

	return nil
}

// SetDescriptionAction sets the profile description string
func SetDescriptionAction(c *cli.Context) error {

	config := Configurations[ProfileName]
	config.Description = c.GetParameter(0)
	Configurations[ProfileName] = config
	profileDirty = true

	return nil
}
