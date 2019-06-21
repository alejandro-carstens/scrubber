package responses

import "github.com/Jeffail/gabs"

type Informable interface {
	Marshal(container *gabs.Container) (Informable, error)

	IsSnapshotInfo() bool

	Name() string

	CreationDate() string
}
