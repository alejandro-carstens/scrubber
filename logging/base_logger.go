package logging

import (
	"fmt"
	"log"
)

type baseLogger struct {
	debug      bool
	trace      bool
	infoLabel  string
	warnLabel  string
	errorLabel string
	fatalLabel string
	debugLabel string
	traceLabel string
	logger     *log.Logger
}

func (bl *baseLogger) Noticef(format string, v ...interface{}) {
	bl.logger.Printf(bl.infoLabel+format, v...)
}

func (bl *baseLogger) Warnf(format string, v ...interface{}) {
	bl.logger.Printf(bl.warnLabel+format, v...)
}

func (bl *baseLogger) Errorf(format string, v ...interface{}) {
	bl.logger.Printf(bl.errorLabel+format, v...)
}

func (bl *baseLogger) Fatalf(format string, v ...interface{}) {
	bl.logger.Fatalf(bl.fatalLabel+format, v...)
}

func (bl *baseLogger) Debugf(format string, v ...interface{}) {
	if bl.debug {
		bl.logger.Printf(bl.debugLabel+format, v...)
	}
}

func (bl *baseLogger) Tracef(format string, v ...interface{}) {
	if bl.trace {
		bl.logger.Printf(bl.traceLabel+format, v...)
	}
}

func (bl *baseLogger) setDebug(debug bool) {
	bl.debug = debug
}

func (bl *baseLogger) setTrace(trace bool) {
	bl.trace = trace
}

func (bl *baseLogger) setLogger(logger *log.Logger) {
	bl.logger = logger
}

func (bl *baseLogger) setPlainLabelFormats() {
	bl.infoLabel = "[INF] "
	bl.debugLabel = "[DBG] "
	bl.warnLabel = "[WRN] "
	bl.errorLabel = "[ERR] "
	bl.fatalLabel = "[FTL] "
	bl.traceLabel = "[TRC] "
}

func (bl *baseLogger) setColoredLabelFormats() {
	colorFormat := "[\x1b[%sm%s\x1b[0m] "

	bl.infoLabel = fmt.Sprintf(colorFormat, "32", "INF")
	bl.debugLabel = fmt.Sprintf(colorFormat, "36", "DBG")
	bl.warnLabel = fmt.Sprintf(colorFormat, "0;93", "WRN")
	bl.errorLabel = fmt.Sprintf(colorFormat, "31", "ERR")
	bl.fatalLabel = fmt.Sprintf(colorFormat, "31", "FTL")
	bl.traceLabel = fmt.Sprintf(colorFormat, "33", "TRC")
}
