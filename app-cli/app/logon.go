package app

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty"
	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/app-cli/persistence"
	"github.com/tucats/gopackages/app-cli/ui"
)

const LogonEndpoint = "/services/logon"
const LogonServerSetting = "logon-server"
const LogonTokenSetting = "logon-token"

// LogonGrammar describes the login subcommand
var LogonGrammar = []cli.Option{
	{
		LongName:            "username",
		ShortName:           "u",
		OptionType:          cli.StringType,
		Description:         "Username for login",
		EnvironmentVariable: "EGO_USERNAME",
	},
	{
		LongName:            "password",
		ShortName:           "p",
		OptionType:          cli.StringType,
		Description:         "Password for login",
		EnvironmentVariable: "EGO_PASSWORD",
	},
	{
		LongName:            "logon-server",
		ShortName:           "l",
		OptionType:          cli.StringType,
		Description:         "URL of logon server",
		EnvironmentVariable: "EGO_LOGON_SERVER",
	},
}

// Logon handles the logon subcommand
func Logon(c *cli.Context) error {

	// Do we know where the logon server is?
	url := persistence.Get(LogonServerSetting)
	if c.WasFound("logon-server") {
		url, _ = c.GetString("logon-server")
		persistence.Set(LogonServerSetting, url)
	}
	if url == "" {
		return errors.New("no --logon-server specified")
	}

	user, _ := c.GetString("username")
	pass, _ := c.GetString("password")

	for user == "" {
		user = ui.Prompt("Username: ")
	}
	for pass == "" {
		pass = ui.PromptPassword("Password: ")
	}

	// Turn logon server into full URL
	url = strings.TrimSuffix(url, "/") + LogonEndpoint

	// Call the endpoint
	r, err := resty.New().SetDisableWarn(true).SetBasicAuth(user, pass).NewRequest().Get(url)
	if err == nil && r.StatusCode() == 200 {
		token := strings.TrimSuffix(string(r.Body()), "\n")
		persistence.Set(LogonTokenSetting, token)
		err = persistence.Save()
		if err == nil {
			ui.Say("Successfully logged in as \"%s\"", user)
		}
		return err
	}

	// IF there was an error condition, let's report it now.
	if err == nil {
		switch r.StatusCode() {
		case 401:
			err = errors.New("no credentials provided")
		case 403:
			err = errors.New("invalid credentials")
		case 404:
			err = errors.New("logon endpoint not found")
		default:
			err = fmt.Errorf("HTTP %d", r.StatusCode())
		}
	}

	return err
}
