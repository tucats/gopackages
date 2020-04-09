package app

import (
	"errors"
	"strings"

	"github.com/tucats/gopackages/cli/cli"
	"github.com/tucats/gopackages/cli/profile"
	"github.com/tucats/gopackages/cli/ui"
)

// SetOutputFormat sets the default output format to use.
func SetOutputFormat(c *cli.Options) error {

	if formatString, present := cli.GetString(*c, "output-format"); present {
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

// SetDebugMessaging is an action routine to set the global debug status if specified
func SetDebugMessaging(c *cli.Options) error {
	ui.DebugMode = cli.GetBool(*c, "debug")
	return nil
}

// SetQuietMode is an action routine to set the global debug status if specified
func SetQuietMode(c *cli.Options) error {
	ui.QuietMode = cli.GetBool(*c, "quiet")
	return nil
}

// SetDefaultProfile is the action routine when --profile is specified as a global
// option. It's string value is used as the name of the active profile.
func SetDefaultProfile(c *cli.Options) error {
	name, _ := cli.GetString(*c, "use-profile")
	ui.Debug("Using profile %s", name)
	profile.UseProfile(name)
	return nil
}
