package messages

// Sendable represents the contract of a message
// that can be sent over a Notifiable channel
type Sendable interface {
	// Format sets the format of a Sendable message
	Format() error

	// Type returns the message type
	Type() string
}
