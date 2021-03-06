// Package app provides the top-level framework for CLI execution. This includes
// the Run() method to run the program, plus a number of action routines that can
// be invoked from the grammar or by a user action routine. These support common or
// global actions, like specifying which profile to use.
package app

import (
	"fmt"
	"strings"

	"github.com/tucats/gopackages/app-cli/cli"
	"github.com/tucats/gopackages/symbols"
)

// App is the wrapper type for information needed for a command line application.
// It contains the globals needed for the application as well as the runtime
// context root.
type App struct {
	Name        string
	Description string
	Copyright   string
	Version     string
	Context     *cli.Context
	Action      func(c *cli.Context) error
}

// New creates a new instance of an application context, given the name of the
// application.
func New(appName string) App {
	// Extract the description of the app if it was given
	var appDescription = ""
	if i := strings.Index(appName, ":"); i > 0 {
		appDescription = strings.TrimSpace(appName[i+1:])
		appName = strings.TrimSpace(appName[:i])
	}
	app := App{Name: appName, Description: appDescription}
	return app
}

// SetVersion sets the version number for the application.
func (app *App) SetVersion(major, minor, delta int) {
	app.Version = fmt.Sprintf("%d.%d-%d", major, minor, delta)
	symbols.RootSymbolTable.Symbols["_version"] = app.Version
}

// SetCopyright sets the copyright string (if any) used in the
// help output.
func (app *App) SetCopyright(s string) {
	app.Copyright = s
	symbols.RootSymbolTable.Symbols["_copyright"] = app.Copyright

}

// Parse runs a grammar, and then calls the provided action routine. It is typically
// used in cases where there are no subcommands, and an action should be run after
// parsing options.
func (app *App) Parse(grammar []cli.Option, args []string, action func(c *cli.Context) error) error {
	app.Action = action
	return app.Run(grammar, args)
}

// Run runs a grammar given a set of arguments in the current
// applciation. The grammar must declare action routines for the
// various subcommands, which will be executed by the parser.
func (app *App) Run(grammar []cli.Option, args []string) error {
	app.Context = &cli.Context{
		Description: app.Description,
		Copyright:   app.Copyright,
		Version:     app.Version,
		AppName:     app.Name,
		Grammar:     grammar,
		Args:        args,
		Action:      app.Action,
	}
	return runFromContext(app.Context)
}
