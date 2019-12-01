package messages

import "github.com/Jeffail/gabs"

// Sendable represents the contract of a message
// that can be sent over a Notifiable channel
type Sendable interface {
	// Format sets the format of a Sendable message
	Format() error

	// Type returns the message type
	Type() string

	// Init initializes the base properties of a message
	Init(context interface{}, dedupKey string, payload *gabs.Container)
}
