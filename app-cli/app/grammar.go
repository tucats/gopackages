package app

import (
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/config"
	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
)

// Create a new applicationGrammar which prepends the default supplied options
// to the caller's applicationGrammar definition.
var applicationGrammar = []cli.Option{
	{
		LongName:            "insecure",
		ShortName:           "k",
		OptionType:          cli.BooleanType,
		Description:         "insecure",
		Action:              InsecureAction,
		EnvironmentVariable: "APP_INSECURE_CLIENT",
	},
	{
		LongName:    "config",
		Aliases:     []string{"configuration", "profile", "prof"},
		OptionType:  cli.Subcommand,
		Description: "app.config",
		Value:       config.Grammar,
	},
	{
		LongName:    "logon",
		Aliases:     []string{"login"},
		OptionType:  cli.Subcommand,
		Description: "app.logon",
		Action:      Logon,
		Value:       LogonGrammar,
	},
	{
		ShortName:           "p",
		LongName:            "profile",
		Description:         "global.profile",
		OptionType:          cli.StringType,
		Action:              UseProfileAction,
		EnvironmentVariable: "APP_PROFILE",
	},
	{
		LongName:            "log",
		ShortName:           "l",
		Description:         "global.log",
		OptionType:          cli.StringListType,
		Action:              LogAction,
		EnvironmentVariable: defs.DefaultLogging,
	},
	{
		LongName:            "log-file",
		Description:         "global.log.file",
		OptionType:          cli.StringType,
		Action:              LogFileAction,
		EnvironmentVariable: defs.DefaultLogFileName,
	},
	{
		LongName:            "format",
		ShortName:           "f",
		Description:         "global.format",
		OptionType:          cli.KeywordType,
		Keywords:            []string{ui.JSONFormat, ui.JSONIndentedFormat, ui.TextFormat},
		Action:              OutputFormatAction,
		EnvironmentVariable: "APP_OUTPUT_FORMAT",
	},
	{
		ShortName:   "v",
		LongName:    "version",
		Description: "global.version",
		OptionType:  cli.BooleanType,
		Action:      ShowVersionAction,
	},
	{
		ShortName:           "q",
		LongName:            "quiet",
		Description:         "global.quiet",
		OptionType:          cli.BooleanType,
		Action:              QuietAction,
		EnvironmentVariable: "APP_QUIET",
	},
	{
		LongName:    "version",
		Description: "global.version",
		OptionType:  cli.Subcommand,
		Action:      VersionAction,
	},
}

// MakePrivate allows an application to disable a grammar item
// in the built-in application commands like "config" or "logon",
// for those applications that do not want ot support these
// options. For example, an application that does not make use
// of a remote server does not need the "logon" builtin command.
//
// The function returns true if the option was found and disabled,
// and returns false if the option name was not found in the built-in
// grammar.
func MakePrivate(name string) bool {
	for n, option := range applicationGrammar {
		if option.LongName == name {
			applicationGrammar[n].Private = true

			return true
		}
	}

	return false
}
