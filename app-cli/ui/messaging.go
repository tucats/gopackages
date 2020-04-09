package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (

	// DefaultTableFormat means use whatever the default is that may have been set
	// by the global option --output-type, etc.
	DefaultTableFormat = 0

	// TextTableFormat indicates the output format should be human-readable text
	TextTableFormat = 1

	// JSONTableFormat indicates the output format should be machine-readable JSON
	JSONTableFormat = 2
)

// OutputFormat is the default output format if not overridden by a global option
// or explicit call from the user.
var OutputFormat int = TextTableFormat

// DebugMode determines if "debug" style messages are output.
var DebugMode = false

// QuietMode determines if optional messaging is performed.
var QuietMode = false

var sequence = 0
var sequenceMux sync.Mutex

// Debug displays a message if debugging mode is enabled.
func Debug(format string, args ...interface{}) {
	if DebugMode {
		pid := os.Getpid()

		s := fmt.Sprintf(format, args...)
		sequenceMux.Lock()
		sequence = sequence + 1
		sequenceString := fmt.Sprintf("%d, %d", pid, sequence)
		fmt.Printf("[%s] %-10s DEBUG: %s\n", time.Now().Format(time.RFC3339), sequenceString, s)
		sequenceMux.Unlock()
	}
}

// Say displays a message to the user unless we are in "quiet" mode
func Say(format string, args ...interface{}) {
	if !QuietMode {
		s := fmt.Sprintf(format, args...)
		fmt.Println(s)
	}
}
