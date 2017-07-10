package config

import (
	"fmt"
	"os"
	"time"

	"path/filepath"
)

// Functions

// CreateLog checks for existence of a 'results'
// folder in current directory and creates and
// opens a log file for the current run.
func CreateLog() (*os.File, error) {

	// Retrieve current directory.
	dir, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	// Path to results directory.
	resultsDir := filepath.Join(dir, "results")

	// Name and path to log file for this run.
	logFileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15-04-05"))
	logFilePath := filepath.Join(resultsDir, logFileName)

	// Ensure that a folder 'results' is present.
	if _, err := os.Stat(resultsDir); os.IsNotExist(err) {

		// Create all folders including 'results'.
		err := os.MkdirAll(resultsDir, 0744)
		if err != nil {
			return nil, err
		}
	}

	// Create log file for this run.
	logFile, err := os.OpenFile(logFilePath, (os.O_CREATE | os.O_RDWR), 0755)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
