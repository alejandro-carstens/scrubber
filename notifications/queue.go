package notifications

import (
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications/channels"
	"github.com/alejandro-carstens/scrubber/notifications/messages"

	"github.com/panjf2000/ants/v2"
)

// Queue represents the notifications
// pusher and handler
type Queue struct {
	pool   *ants.PoolWithFunc
	logger *logger.Logger
}

// Init initializes the Queue
func (q *Queue) Init(capacity int, logger *logger.Logger) (*Queue, error) {
	pool, err := ants.NewPoolWithFunc(capacity, func(message interface{}) {
		msg, valid := message.(messages.Sendable)

		if !valid {
			q.logger.Errorf("message not of type messages.Sendable")

			return
		}

		if err := channels.Notify(msg); err != nil {
			q.logger.Errorf(err.Error())
		}
	})

	if err != nil {
		return nil, err
	}

	q.pool = pool
	q.logger = logger

	return q, nil
}

// Push pushes a notification to the queue
func (q *Queue) Push(message messages.Sendable) error {
	return q.pool.Invoke(message)
}

// Release clears the go routine pool
func (q *Queue) Release() {
	q.pool.Release()
}
