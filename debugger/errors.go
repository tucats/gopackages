package debugger

import (
	"errors"
)

const (
	InvalidBreakClauseError = "invalid break clause"
)

var SignalDebugger = errors.New("signal")

func InvokeDebugger(e error) bool {
	return e != nil && e.Error() == SignalDebugger.Error()
}
