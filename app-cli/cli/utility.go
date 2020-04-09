package cli

import (
	"errors"
	"strings"

	"github.com/tucats/gopackages/cli/profile"
	"github.com/tucats/gopackages/cli/ui"
)

// ValidKeyword does a case-insensitive compare of a string containing
// a keyword against a list of possible stirng values.
func ValidKeyword(test string, valid []string) bool {

	for _, v := range valid {
		if strings.ToLower(test) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

// FindKeyword does a case-insensitive compare of a string containing
// a keyword against a list of possible string values. If the keyword
// is found, it's position in the list is returned. If it was not found,
// the value returned is -1
func FindKeyword(test string, valid []string) int {

	for n, v := range valid {
		if strings.ToLower(test) == strings.ToLower(v) {
			return n
		}
	}
	return -1
}

// SetDebugMessaging is an action routine to set the global debug status if specified
func SetDebugMessaging(c *Options) error {
	ui.DebugMode = GetBool(*c, "debug")
	return nil
}

// SetOutputFormat sets the default output format to use.
func SetOutputFormat(c *Options) error {

	if formatString, present := GetString(*c, "output-format"); present {
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

// SetQuietMode is an action routine to set the global debug status if specified
func SetQuietMode(c *Options) error {
	ui.QuietMode = GetBool(*c, "quiet")
	return nil
}

// SetDefaultProfile is the action routine when --profile is specified as a global
// option. It's string value is used as the name of the active profile.
func SetDefaultProfile(c *Options) error {
	name, _ := GetString(*c, "profile")
	profile.UseProfile(name)
	return nil
}
