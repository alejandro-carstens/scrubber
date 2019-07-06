package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/infos"
)

// Runnerable represents a contract to be
// implemented by filter runners
type Runnerable interface {
	// Init initializes the filter runner
	Init(info ...infos.Informable) (Runnerable, error)

	// RunFilter filters out elements from the actionable list
	RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable)
}
