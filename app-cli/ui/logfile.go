package ui

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tucats/gopackages/errors"
)

var logFile *os.File
var baseLogFileName string
var currentLogFileName string

// LogRetainCount is the number of roll-over log versions to keep in the
// logging directory.
var LogRetainCount = -1

func OpenLogFile(userLogFileName string, withTimeStamp bool) error {
	if LogRetainCount < 1 {
		LogRetainCount = 3
	}

	err := openLogFile(userLogFileName, withTimeStamp)
	if err != nil {
		return errors.NewError(err)
	}

	if withTimeStamp {
		PurgeLogs()

		go rollOverTask()
	}

	return nil
}

// Return the path of the current log file being written to.
func CurrentLogFile() string {
	if logFile == nil {
		return ""
	}

	return currentLogFileName
}

// Internal routine that actually opens a log file.
func openLogFile(path string, withTimeStamp bool) error {
	var err error

	_ = SaveLastLog()

	var fileName string

	if withTimeStamp {
		fileName = timeStampLogFileName(path)
	} else {
		fileName, err = filepath.Abs(path)
		if err != nil {
			return errors.NewError(err)
		}
	}

	logFile, err = os.Create(fileName)
	if err != nil {
		logFile = nil

		return errors.NewError(err)
	}

	baseLogFileName, _ = filepath.Abs(path)
	currentLogFileName, _ = filepath.Abs(fileName)

	WriteLog(InfoLogger, "New log file opened: %s", currentLogFileName)

	return nil
}

// Schedule roll-over operations for the log. We calculate when the next start-of-date + 24 hours
// is, and sleep until then. We then roll over the log file and sleep again.
func rollOverTask() {
	for {
		year, month, day := time.Now().Date()
		beginningOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
		wakeTime := beginningOfDay.Add(24*time.Hour + time.Second)
		sleepUntil := time.Until(wakeTime)
		WriteLog(InfoLogger, "Log rollover scheduled for %s", wakeTime.String())
		time.Sleep(sleepUntil)
		RollOverLog()
	}
}

// Roll over the open log. Close the current log, and rename it to include a timestamp of when
// it was created. Then create a new log file.
func RollOverLog() {
	if err1 := SaveLastLog(); err1 != nil {
		WriteLog(InternalLogger, "ERROR: RollOverLog() unable to roll over log file; %v", err1)

		return
	}

	err := openLogFile(baseLogFileName, true)
	if err != nil {
		WriteLog(InternalLogger, "ERROR: RollOverLog() unable to open new log file; %v", err)

		return
	}

	PurgeLogs()
}

func timeStampLogFileName(path string) string {
	logStarted := time.Now()
	dateStamp := logStarted.Format("_2006-01-02-150405")
	newName, _ := filepath.Abs(strings.TrimSuffix(path, ".log") + dateStamp + ".log")

	return newName
}

// Save the current (last) log file to the archive name with the timestamp of when the log
// was initialized.
func SaveLastLog() error {
	if logFile != nil {
		WriteLog(InfoLogger, "Log file being rolled over")

		sequenceMux.Lock()
		defer sequenceMux.Unlock()
		logFile.Close()

		logFile = nil
	}

	return nil
}

func PurgeLogs() int {
	count := 0
	keep := LogRetainCount
	searchPath := path.Dir(CurrentLogFile())
	names := []string{}

	Log(ServerLogger, "Purging all but %d logs from %s", keep, searchPath)

	files, err := ioutil.ReadDir(searchPath)
	if err != nil {
		Log(ServerLogger, "Error making list of log files, %s", err.Error())

		return count
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "ego-server_") && !file.IsDir() {
			names = append(names, file.Name())
		}
	}

	if len(names) <= 1 {
		return 0
	}

	sort.Strings(names)

	for n := 0; n < len(names)-keep; n++ {
		name := names[n]
		fileName := path.Join(searchPath, name)

		err := os.Remove(fileName)
		if err != nil {
			Log(ServerLogger, "Error purging log file, %v", err)
		} else {
			Log(ServerLogger, "Purged log file %s", fileName)
			count++
		}
	}

	return count
}
