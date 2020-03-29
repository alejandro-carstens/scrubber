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

// Validate implementation of the Configurable interface
func (l *Local) Validate() error {
	if len(l.Path) == 0 {
		return errors.New("path cannot be empty")
	}

	return nil
}

// Name implementation of the Configurable interface
func (l *Local) Name() string {
	return "local"
}

type local struct {
	path          string
	streamFile    *os.File
	streamChannel chan string
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
	l.streamChannel = make(chan string)

	return nil
}

// Stream implementation of the Storeable interface
func (l *local) Stream() error {
	if l.streamFile == nil || l.streamChannel == nil {
		return errors.New("please open a stream")
	}

	for line := range l.streamChannel {
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

// Channel implementation of the Storeable interface
func (l *local) Channel(data string) {
	l.streamChannel <- data
}

// Remove implementation of the Storeable interface
func (l *local) Remove(name string) error {
	return os.Remove(filepath.Join(l.path, filepath.FromSlash(name)))
}

// Close implementation of the Storeable interface
func (l *local) Close() error {
	if l.streamFile == nil || l.streamChannel == nil {
		return errors.New("please open a stream")
	}

	close(l.streamChannel)

	return l.streamFile.Close()
}
