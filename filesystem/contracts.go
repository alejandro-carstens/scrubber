package filesystem

import (
	"io"
	"os"
)

// Storeable represents the contract to be implemented by
// different filesystems in order to store, retrieve,
// stream and delete files.
type Storeable interface {
	// Init initializes the store with the proper configs
	Init(config Configurable) (Storeable, error)

	// Put stores a file with the specified name
	Put(name string, reader io.Reader) error

	// Stream streams data into a file
	Stream() error

	// OpenStream initializes the stream
	OpenStream(name string) error

	// Channel returns the stream channel
	Channel(data string)

	// Close releases the streaming resources
	Close() error

	// List lists all the file/directory names in a directory
	List(name string) ([]string, error)

	// Open retrieves a file
	Open(name string) (*os.File, error)
}

// Configurable represents the contract to be implemented
// in order to comply with setting the configurations for
// the different notification channels
type Configurable interface {
	// Validate validates the configuration for a given channel
	Validate() error

	// Name returns the configuration name
	Name() string
}
