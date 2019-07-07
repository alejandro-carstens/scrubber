package logger

import "os"

type fileLogger struct {
	baseLogger
	logFile *os.File
}

func (fl *fileLogger) Close() error {
	if file := fl.logFile; file != nil {
		fl.logFile = nil

		return file.Close()
	}

	return nil
}
