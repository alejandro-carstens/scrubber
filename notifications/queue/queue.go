package queue

import (
	"scrubber/logger"
	"scrubber/notifications/channels"
	"scrubber/notifications/messages"

	"github.com/panjf2000/ants"
)

// Enqueuer represents the notifications
// pusher and handler
type Enqueuer struct {
	pool   *ants.PoolWithFunc
	logger *logger.Logger
}

// Init initializes the Enqueuer
func (e *Enqueuer) Init(capacity int, logger *logger.Logger) (*Enqueuer, error) {
	pool, err := ants.NewPoolWithFunc(capacity, func(message interface{}) {
		msg, valid := message.(messages.Sendable)

		if !valid {
			e.logger.Errorf("message not of type messages.Sendable")

			return
		}

		if err := channels.Notify(msg); err != nil {
			e.logger.Errorf(err.Error())
		}
	})

	if err != nil {
		return nil, err
	}

	e.pool = pool
	e.logger = logger

	return e, nil
}

// Push pushes a notification to the queue
func (e *Enqueuer) Push(message messages.Sendable) error {
	return e.pool.Invoke(message)
}

// Release clears the go routine pool
func (e *Enqueuer) Release() {
	e.pool.Release()
}
