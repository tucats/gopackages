package rest

import (
	"runtime"

	"github.com/go-resty/resty"
	"github.com/tucats/gopackages/app-cli/ui"
	"github.com/tucats/gopackages/defs"
	"github.com/tucats/gopackages/expressions/data"
	"github.com/tucats/gopackages/expressions/symbols"
)

var allowInsecure = false

func AllowInsecure(b bool) {
	allowInsecure = b
}

func AddAgent(r *resty.Request, agentType string) {
	var version string

	if x, found := symbols.RootSymbolTable.Get(defs.VersionName); found {
		version = data.String(x)
	}

	platform := runtime.Version() + ", " + runtime.GOOS + ", " + runtime.GOARCH
	agent := "Ego " + version + " (" + platform + ") " + agentType

	r.Header.Add("User-Agent", agent)
	ui.Log(ui.RestLogger, "User agent: %s", agent)
}
