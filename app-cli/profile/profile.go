package profile

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/app-cli/tables"
	"github.com/tucats/gopackages/app-cli/ui"
)

// Grammar describes profile subcommands
var Grammar = []cli.Option{
	{
		LongName:    "list",
		Description: "List all profiles",
		Action:      ListAction,
		OptionType:  cli.Subcommand,
	},
	{
		LongName:    "show",
		Description: "Show the current profile",
		Action:      ShowAction,
		OptionType:  cli.Subcommand,
	},
	{
		LongName:             "set-output",
		OptionType:           cli.Subcommand,
		Description:          "Set the default output type (text or json)",
		ParameterDescription: "type",
		Action:               SetOutputAction,
		ParametersExpected:   1,
	},
	{
		LongName:             "set-description",
		OptionType:           cli.Subcommand,
		Description:          "Set the profile description",
		ParameterDescription: "text",
		ParametersExpected:   1,
		Action:               SetDescriptionAction,
	},
	{
		LongName:             "delete",
		OptionType:           cli.Subcommand,
		Description:          "Delete a key from the profile",
		Action:               DeleteAction,
		ParametersExpected:   1,
		ParameterDescription: "key",
	},
	{
		LongName:    "set",
		Description: "Set a profile value",
		Action:      SetAction,
		OptionType:  cli.Subcommand,
		Value: []cli.Option{
			{
				LongName:    "key",
				Description: "The key that will be set in the profile. Can be of the form key=value.",
				OptionType:  cli.StringType,
				Required:    true,
			},
			{
				LongName:    "value",
				Description: "The value to set for the key. If missing, the key is deleted",
				OptionType:  cli.StringType,
			},
		},
	},
}

// ShowAction Displays the current contents of the active profile
func ShowAction(c *cli.Context) error {

	t, _ := tables.New([]string{"Key", "Value"})

	for k, v := range persistence.CurrentConfiguration.Items {
		if len(fmt.Sprintf("%v", v)) > 60 {
			v = fmt.Sprintf("%v", v)[:60] + "..."
		}
		_ = t.AddRowItems(k, v)
	}
	_ = t.SetOrderBy("key")
	t.ShowUnderlines(false)
	t.Print(ui.TextTableFormat)

	return nil
}

// ListAction Displays the current contents of the active profile
func ListAction(c *cli.Context) error {

	t, _ := tables.New([]string{"Name", "Description"})

	for k, v := range persistence.Configurations {
		_ = t.AddRowItems(k, v.Description)
	}
	_ = t.SetOrderBy("name")
	t.ShowUnderlines(false)
	t.Print(ui.TextTableFormat)

	return nil
}

// SetOutputAction is the action handler for the set-output subcommand.
func SetOutputAction(c *cli.Context) error {

	if c.GetParameterCount() == 1 {
		outputType := c.GetParameter(0)
		if outputType == "text" || outputType == "json" {
			persistence.Set("output-format", outputType)
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
		persistence.Set(key, value)
		ui.Say("Profile key %s written", key)
	} else {
		persistence.Delete(key)
		ui.Say("Profile key %s deleted", key)
	}

	return nil
}

// DeleteAction deletes a named key value
func DeleteAction(c *cli.Context) error {

	key := c.GetParameter(0)
	persistence.Delete(key)
	ui.Say("Profile key %s deleted", key)

	return nil
}

// SetDescriptionAction sets the profile description string
func SetDescriptionAction(c *cli.Context) error {

	config := persistence.Configurations[persistence.ProfileName]
	config.Description = c.GetParameter(0)
	persistence.Configurations[persistence.ProfileName] = config
	persistence.ProfileDirty = true

	return nil
}
