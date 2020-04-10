package app

import (
	"errors"
	"strings"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/profile"
	"github.com/tucats/gopackages/cli/ui"
)

// OutputFormatAction sets the default output format to use.
func OutputFormatAction(c *cli.Context) error {

	if formatString, present := c.FindGlobal().GetString("output-format"); present {
		switch strings.ToLower(formatString) {
		case "text":
			ui.OutputFormat = ui.TextTableFormat

		case "json":
			ui.OutputFormat = ui.JSONTableFormat

		default:
			return errors.New("Invalid output format specified: " + formatString)
		}
		profile.SetDefault("output-format", strings.ToLower(formatString))
	}
	return nil
}

// DebugAction is an action routine to set the global debug status if specified
func DebugAction(c *cli.Context) error {
	ui.DebugMode = c.FindGlobal().GetBool("debug")
	return nil
}

// QuietAction is an action routine to set the global debug status if specified
func QuietAction(c *cli.Context) error {
	ui.QuietMode = c.FindGlobal().GetBool("quiet")
	return nil
}

// UseProfileAction is the action routine when --profile is specified as a global
// option. It's string value is used as the name of the active profile.
func UseProfileAction(c *cli.Context) error {
	name, _ := c.GetString("use-profile")
	ui.Debug("Using profile %s", name)
	profile.UseProfile(name)
	return nil
}
