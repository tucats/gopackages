package debugger

import (
	"errors"
)

var SignalDebugger = errors.New("signal")

func InvokeDebugger(e error) bool {
	return e.Error() == SignalDebugger.Error()
}
