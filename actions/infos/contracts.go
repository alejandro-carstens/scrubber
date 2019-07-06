package infos

import "github.com/Jeffail/gabs"

// Informable is a contract implemented by structs that hold
// important information about an index or snapshot
type Informable interface {
	// Marshal takes in a raw *gabs.Container and parses it to
	// an Informable implementation
	Marshal(container *gabs.Container) (Informable, error)

	// IsSnapshotInfo signals if the Informable instance
	// corresponds to a snpashot or an index
	IsSnapshotInfo() bool

	// Name returns the name of the index or snaphsot
	// depending on the Informable implementation
	Name() string

	// Returns the creation_date as a string of a given
	// Informablt implementation
	CreationDate() string
}
