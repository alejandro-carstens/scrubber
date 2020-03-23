package filesystem

import "io"

// Storeable represents the contract to be implemented by
// different filesystems in order to store, retrieve,
// stream and delete files.
type Storeable interface {
	// Init initializes the store with the proper configs
	Init(config Configurable) (Storeable, error)

	// Put stores a file with the specified name
	Put(name string, reader io.Reader) error

	// Remove deletes a file with the specified name
	Remove(name string) error

	// Stream streams data into a file
	Stream(content chan string) error

	// OpenStream initializes the stream
	OpenStream(name string) error
}

// Configurable represents the contract to be implemented
// in order to comply with setting the configurations for
// the different notification channels
type Configurable interface {
	// Validate validates the configuration for a given channel
	Validate() error

	// FillFromEnvs is responsible for setting the configuration
	// for the channel from the respective env variables
	FillFromEnvs() Configurable

	// Name returns the configuration name
	Name() string
}
