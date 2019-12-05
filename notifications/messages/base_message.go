package messages

import "github.com/Jeffail/gabs"

type baseMessage struct {
	Context  interface{}
	DedupKey string
	Payload  *gabs.Container
}

// Init initializes the base properties of a message
func (bm *baseMessage) Init(context interface{}, dedupKey string, payload *gabs.Container) {
	bm.Context = context
	bm.DedupKey = dedupKey
	bm.Payload = payload
}
