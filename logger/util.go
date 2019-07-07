package logger

import (
	"fmt"
	"log"
	"os"
)

func NewLogger(filename string, time, debug, trace, pid bool) *Logger {
	Logger := &Logger{}

	if len(filename) > 0 {
		Logger.logger = newFileLogger(filename, time, debug, trace, pid)
	} else {
		Logger.logger = newStdLogger(time, debug, trace, true, pid)
	}

	Logger.opts = &opts{
		filename,
		time,
		debug,
		trace,
		pid,
	}

	return Logger
}

func newFileLogger(filename string, time, debug, trace, pid bool) loggable {
	fileflags := os.O_WRONLY | os.O_APPEND | os.O_CREATE

	file, err := os.OpenFile(filename, fileflags, 0660)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	flags := 0

	if time {
		flags = log.LstdFlags | log.Lmicroseconds
	}

	pre := ""

	if pid {
		pre = fmt.Sprintf("[%d] ", os.Getpid())
	}

	fl := new(fileLogger)
	fl.setLogger(log.New(file, pre, flags))
	fl.setDebug(debug)
	fl.setTrace(trace)
	fl.setPlainLabelFormats()

	return fl
}

func newStdLogger(time, debug, trace, colors, pid bool) loggable {
	flags := 0

	if time {
		flags = log.LstdFlags | log.Lmicroseconds
	}

	pre := ""

	if pid {
		pre = fmt.Sprintf("[%d] ", os.Getpid())
	}

	sl := new(stdLogger)
	sl.setLogger(log.New(os.Stderr, pre, flags))
	sl.setDebug(debug)
	sl.setTrace(trace)

	if colors {
		sl.setColoredLabelFormats()
	} else {
		sl.setPlainLabelFormats()
	}

	return sl
}
