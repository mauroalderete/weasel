package logger

import (
	"fmt"
	"os"
)

type logger struct {
	logfile *os.File
	verbose bool
}

func (l *logger) Write(data []byte) (int, error) {

	var n int
	var err error

	if l.logfile != nil {
		n, err = fmt.Fprintf(l.logfile, "%s", data)
		if err != nil {
			return 0, fmt.Errorf("failed to save in file: %v", err)
		}
		l.logfile.Sync()
	}

	if l.verbose {
		fmt.Printf("%s", data)
	}

	return n, err
}

func New(logfile string, verbose bool) (*logger, error) {

	l := &logger{
		verbose: verbose,
	}

	if len(logfile) != 0 {
		file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed open the logfile '%s' to store the execution log: %v", logfile, err)
		}
		l.logfile = file
	}

	return l, nil
}
