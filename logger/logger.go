package logger

import "sync"

type logfn func(logger loggable, format string, v ...interface{})

type opts struct {
	filename string
	time     bool
	debug    bool
	trace    bool
	pid      bool
}

type Logger struct {
	sync.RWMutex
	opts   *opts
	logger loggable
}

func (sl *Logger) Noticef(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Noticef(format, v...)
	}, format, v...)
}

func (sl *Logger) Errorf(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Errorf(format, v...)
	}, format, v...)
}

func (sl *Logger) Warnf(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Warnf(format, v...)
	}, format, v...)
}

func (sl *Logger) Fatalf(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Fatalf(format, v...)
	}, format, v...)
}

func (sl *Logger) Debugf(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Debugf(format, v...)
	}, format, v...)
}

func (sl *Logger) Tracef(format string, v ...interface{}) {
	sl.executeLogCall(func(logger loggable, format string, v ...interface{}) {
		logger.Tracef(format, v...)
	}, format, v...)
}

func (sl *Logger) executeLogCall(fn logfn, format string, args ...interface{}) {
	sl.RLock()

	defer sl.RUnlock()

	if &sl.logger == nil {
		return
	}

	fn(sl.logger, format, args...)
}

func (sl *Logger) ReOpenLogFile() {
	sl.RLock()

	logger := sl.logger

	sl.RUnlock()

	if logger == nil {
		sl.Noticef("File log re-open ignored, no logger")

		return
	}

	if sl.opts.filename == "" {
		sl.Noticef("File log re-open ignored, not a file logger")

		return
	}

	sl.logger = newFileLogger(sl.opts.filename, sl.opts.time, sl.opts.debug, sl.opts.trace, true)
	sl.Noticef("File log re-opened")
}
