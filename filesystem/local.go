package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Local represents the configuration
// for the local filesystem
type Local struct {
	Path string
}

// Validate implementation of
// the Configurable interface
func (l *Local) Validate() error {
	if len(l.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	return nil
}

// FillFromEnvs implementation of
// the Configurable interface
func (l *Local) FillFromEnvs() Configurable {
	return l
}

// Name implementation of the
// Configurable interface
func (l *Local) Name() string {
	return "local"
}

type local struct {
	path       string
	streamFile *os.File
}

// Init implementation of the Configurable interface
func (l *local) Init(configuration Configurable) (Storeable, error) {
	config := configuration.(*Local)

	l.path = config.Path

	return l, nil
}

// Put implementation of the Storeable interface
func (l *local) Put(name string, reader io.Reader) error {
	path := filepath.Join(l.path, filepath.FromSlash(name))

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, reader)

	return err
}

// OpenStream implementation of the Storeable interface
func (l *local) OpenStream(name string) error {
	path := filepath.Join(l.path, filepath.FromSlash(name))

	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	l.streamFile = file

	return nil
}

// Stream implementation of the Storeable interface
func (l *local) Stream(content chan string) error {
	if l.streamFile == nil {
		return errors.New("please open a stream file")
	}

	defer l.streamFile.Close()

	for line := range content {
		n, err := l.streamFile.Write([]byte(line))

		if err != nil {
			return err
		}

		if n != len(line) {
			return errors.New("failed to write data")
		}
	}

	return nil
}

// Remove implementation of the Storeable interface
func (l *local) Remove(name string) error {
	return os.Remove(filepath.Join(l.path, filepath.FromSlash(name)))
}
