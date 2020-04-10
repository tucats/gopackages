// Package ui contains basic tools for interacting with a user. This includes generating
// informational and debugging messages. It also includes functions for controlling
// whether those messages are displayed or not.
package ui

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// Formatted output types for data more complex than individual messages, such
// as the format for tabular data output. Choices are "text", "json", or default
// which means whatever was set by the command line or profile.
const (

	// DefaultTableFormat means use whatever the default is that may have been set
	// by the global option --output-type, etc.
	DefaultTableFormat = "default"

	// TextTableFormat indicates the output format should be human-readable text
	TextTableFormat = "text"

	// JSONTableFormat indicates the output format should be machine-readable JSON
	JSONTableFormat = "json"
)

// OutputFormat is the default output format if not overridden by a global option
// or explicit call from the user.
var OutputFormat = TextTableFormat

// DebugMode determines if "debug" style messages are output.
var DebugMode = false

// QuietMode determines if optional messaging is performed.
var QuietMode = false

// The sequence number is generated and incremented for each message, in order. The
// associated mutext is used to prevent the sequence from being incremented by a
// separate thread or goroutine.
var sequence = 0
var sequenceMux sync.Mutex

// Debug displays a message if debugging mode is enabled.
func Debug(format string, args ...interface{}) {
	if DebugMode {
		Log("DEBUG", format, args...)
	}
}

// Log displays a message to stdout
func Log(class string, format string, args ...interface{}) {
	pid := os.Getpid()
	s := fmt.Sprintf(format, args...)
	sequenceMux.Lock()
	sequence = sequence + 1
	sequenceString := fmt.Sprintf("%d, %d", pid, sequence)
	fmt.Printf("[%s] %-10s %-7s: %s\n", class, time.Now().Format(time.RFC3339), sequenceString, s)
	sequenceMux.Unlock()

}

// Say displays a message to the user unless we are in "quiet" mode
func Say(format string, args ...interface{}) {
	if !QuietMode {
		s := fmt.Sprintf(format, args...)
		fmt.Println(s)
	}
}
