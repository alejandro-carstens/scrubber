package logging

type loggable interface {
	Close() error

	Noticef(format string, v ...interface{})

	Warnf(format string, v ...interface{})

	Errorf(format string, v ...interface{})

	Fatalf(format string, v ...interface{})

	Debugf(format string, v ...interface{})

	Tracef(format string, v ...interface{})
}
