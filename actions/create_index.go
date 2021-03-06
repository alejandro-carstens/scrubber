package actions

import (
	"errors"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type createIndex struct {
	action
	options *options.CreateIndexOptions
}

// ApplyOptions implementation of the Actionable interface
func (ci *createIndex) ApplyOptions() Actionable {
	ci.options = ci.context.Options().(*options.CreateIndexOptions)

	ci.indexer.SetOptions(&golastic.IndexOptions{Timeout: ci.options.TimeoutInSeconds()})

	return ci
}

// Perform implementation of the actionable interface
func (ci *createIndex) Perform() Actionable {
	exists, err := ci.indexer.Exists(ci.options.Name)

	if err != nil {
		ci.errorContainer.push(ci.options.Name, ci.name, err)

		return ci
	}

	if exists {
		ci.errorContainer.push(ci.options.Name, ci.name, errors.New("Index already exists"))

		return ci
	}

	schema, err := mapToString(ci.options.ExtraSettings)

	if err != nil {
		ci.errorContainer.push(ci.options.Name, ci.name, err)

		return ci
	}

	if err := ci.indexer.CreateIndex(ci.options.Name, schema); err != nil {
		ci.errorContainer.push(ci.options.Name, ci.name, err)
	}

	if len(ci.errorContainer.list()) > 0 && ci.retryCount < ci.context.GetRetryCount() {
		ci.retryCount = ci.retryCount + 1
		
		ci.Perform()
	}

	ci.notifiableList = append(ci.notifiableList, ci.options.Name)

	return ci
}

// ApplyFilters implementation of the actionable interface
func (ci *createIndex) ApplyFilters() error {
	return nil
}
