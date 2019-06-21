package logging

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestFileLogger(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "scrubber")

	if err != nil {
		t.Fatal("Could not create tmp dir")
	}

	defer os.RemoveAll(tmpDir)

	file, err := ioutil.TempFile(tmpDir, "scrubber:log_")

	if err != nil {
		t.Fatalf("Could not create the temp file: %v", err)
	}

	file.Close()

	logger := newFileLogger(file.Name(), false, false, false, false)

	logger.Noticef("foo")

	buf, err := ioutil.ReadFile(file.Name())

	if err != nil {
		t.Fatalf("Could not read logfile: %v", err)
	}

	if len(buf) <= 0 {
		t.Fatal("Expected a non-zero length logfile")
	}

	if string(buf) != "[INF] foo\n" {
		t.Fatalf("Expected '%s', received '%s'\n", "[INFO] foo", string(buf))
	}

	file, err = ioutil.TempFile(tmpDir, "scrubber:log_")

	if err != nil {
		t.Fatalf("Could not create the temp file: %v", err)
	}

	file.Close()

	logger = newFileLogger(file.Name(), true, true, true, true)

	logger.Errorf("foo")

	buf, err = ioutil.ReadFile(file.Name())

	if err != nil {
		t.Fatalf("Could not read logfile: %v", err)
	}

	if len(buf) <= 0 {
		t.Fatal("Expected a non-zero length logfile")
	}

	str := string(buf)
	errMsg := fmt.Sprintf("Expected '%s', received '%s'\n", "[pid] <date> [ERR] foo", str)
	pidEnd := strings.Index(str, " ")
	infoStart := strings.LastIndex(str, "[ERR]")

	if pidEnd == -1 || infoStart == -1 {
		t.Fatalf("%v", errMsg)
	}

	pid := str[0:pidEnd]

	if pid[0] != '[' || pid[len(pid)-1] != ']' {
		t.Fatalf("%v", errMsg)
	}

	if !strings.HasSuffix(str, "[ERR] foo\n") {
		t.Fatalf("%v", errMsg)
	}
}
