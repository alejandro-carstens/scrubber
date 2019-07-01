package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/responses"
)

// Runnerable represents a contract to be
// implemented by filter runners
type Runnerable interface {
	// RunFilter filters out elements from the actionable list
	RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable)

	// Init initializes the filter runner
	Init(info ...responses.Informable) (Runnerable, error)
}
