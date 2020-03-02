package options

import "github.com/Jeffail/gabs"

type Optionable interface {
	FillFromContainer(container *gabs.Container) error

	Validate() error

	ValidateNotifiableOptions() error

	IsNotifiable() bool

	GetDisableAction() bool

	GetContainer() *gabs.Container

	Exists(value string) bool

	IsSnapshot() bool

	Get(value string) interface{}

	String(value string) string
}
