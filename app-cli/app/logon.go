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

const (
	// LogonEndpoint is the endpoint for the logon service
	LogonEndpoint = "/services/logon"

	// LogonServerSetting is the name of the profile item that
	// describes the URL of the logon server (less the endpoint)
	LogonServerSetting = "ego.logon.server"

	// LogonTokenSetting is th ename of the profile item that
	// contains the logon token recieved from a succesful logon
	LogonTokenSetting = "ego.logon.token"
)

// LogonGrammar describes the login subcommand
var LogonGrammar = []cli.Option{
	{
		LongName:            "username",
		ShortName:           "u",
		OptionType:          cli.StringType,
		Description:         "Username for login",
		EnvironmentVariable: "CLI_USERNAME",
	},
	{
		LongName:            "password",
		ShortName:           "p",
		OptionType:          cli.StringType,
		Description:         "Password for login",
		EnvironmentVariable: "CLI_PASSWORD",
	},
	{
		LongName:            "logon-server",
		ShortName:           "l",
		OptionType:          cli.StringType,
		Description:         "URL of logon server",
		EnvironmentVariable: "CLI_LOGON_SERVER",
	},
}

// Logon handles the logon subcommand
func Logon(c *cli.Context) error {

	// Do we know where the logon server is? Start with the default from
	// the profile, but if it was explicitly set on the command line, use
	// the command line item and update the saved profile setting.
	url := persistence.Get(LogonServerSetting)
	if c.WasFound("logon-server") {
		url, _ = c.GetString("logon-server")
		persistence.Set(LogonServerSetting, url)
	}
	if url == "" {
		return errors.New("no --logon-server specified")
	}

	// Get the username. If not supplied by the user, prompt until provided.
	user, _ := c.GetString("username")
	for user == "" {
		user = ui.Prompt("Username: ")
	}

	// Get the password. If not supplied by the user, prompt until provided.
	pass, _ := c.GetString("password")
	for pass == "" {
		pass = ui.PromptPassword("Password: ")
	}

	// Turn logon server address and endpoint into full URL
	url = strings.TrimSuffix(url, "/") + LogonEndpoint

	// Call the endpoint
	r, err := resty.New().SetDisableWarn(true).SetBasicAuth(user, pass).NewRequest().Get(url)

	// If the call was successful and the server responded with Success, remove any trailing
	// newline from the result body and store the string as the new token value.
	if err == nil && r.StatusCode() == 200 {
		token := strings.TrimSuffix(string(r.Body()), "\n")
		persistence.Set(LogonTokenSetting, token)
		err = persistence.Save()
		if err == nil {
			ui.Say("Successfully logged in as \"%s\"", user)
		}
		return err
	}

	// If there was an  HTTP error condition, let's report it now.
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
