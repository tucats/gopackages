package cli

import (
	"app-cli/ui"
	"errors"
	"strings"
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
	}
	return nil
}

// SetQuietMode is an action routine to set the global debug status if specified
func SetQuietMode(c *Options) error {
	ui.QuietMode = GetBool(*c, "quiet")
	return nil
}
