package runners

import (
	"scrubber/actions/criterias"
	"scrubber/actions/responses"
)

type Runnerable interface {
	RunFilter(channel chan *FilterResponse, criteria criterias.Criteriable)

	Init(info ...responses.Informable) (Runnerable, error)
}
