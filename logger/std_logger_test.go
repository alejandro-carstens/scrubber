package logger

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestStdLoggerNotice(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, false, false, false, false)
		logger.Noticef("foo")
	}, "[INF] foo\n")
}

func TestStdLoggerNoticeWithColor(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, false, false, true, false)
		logger.Noticef("foo")
	}, "[\x1b[32mINF\x1b[0m] foo\n")
}

func TestStdLoggerDebug(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, true, false, false, false)
		logger.Debugf("foo %s", "bar")
	}, "[DBG] foo bar\n")
}

func TestStdLoggerDebugWithOutDebug(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, false, false, false, false)
		logger.Debugf("foo")
	}, "")
}

func TestStdLoggerTrace(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, false, true, false, false)
		logger.Tracef("foo")
	}, "[TRC] foo\n")
}

func TestStdLoggerTraceWithOutDebug(t *testing.T) {
	expectOutput(t, func() {
		logger := newStdLogger(false, false, false, false, false)
		logger.Tracef("foo")
	}, "")
}

func expectOutput(t *testing.T, f func(), expected string) {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	outC := make(chan string)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Stderr.Close()
	os.Stderr = old

	if out := <-outC; out != expected {
		t.Fatalf("Expected '%s', received '%s'\n", expected, out)
	}
}
