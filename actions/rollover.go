package actions

import (
	"errors"
	"fmt"
	"time"

	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/options"
)

type rollover struct {
	action
	options *options.RolloverOptions
}

// ApplyOptions implementation of the Actionable interface
func (r *rollover) ApplyOptions() Actionable {
	r.options = r.context.Options().(*options.RolloverOptions)

	r.indexer.SetOptions(&golastic.IndexOptions{Timeout: r.options.TimeoutInSeconds()})

	return r
}

// Perform implementation of the Actionable interface
func (r *rollover) Perform() Actionable {
	if err := r.verifyRollableIndex(); err != nil {
		r.errorReportMap.push(r.name, r.options.Name, err)

		return r
	}

	newIndex := r.options.NewIndex

	if len(newIndex) == 0 {
		newIndex = r.generateNewIndexName()
	}

	response, err := r.indexer.Rollover(
		r.options.Name,
		newIndex,
		r.options.MaxAge,
		r.options.MaxSize,
		int64(r.options.MaxDocs),
		r.options.ExtraSettings,
	)

	if err != nil {
		r.errorReportMap.push(r.name, r.options.Name, err)

		return r
	}

	if rolledOver, _ := response.S("rolled_over").Data().(bool); !rolledOver {
		r.errorReportMap.push(r.name, r.options.Name, errors.New("rollover failed: "+response.String()))

		return r
	}

	return r
}

// ApplyFilters implementation of the Actionable interface
func (r *rollover) ApplyFilters() error {
	return nil
}

func (r *rollover) verifyRollableIndex() error {
	response, err := r.indexer.AliasesCat()

	if err != nil {
		return err
	}

	aliases, err := response.Children()

	if err != nil {
		return err
	}

	indices := []string{}

	for _, alias := range aliases {
		aliasName := alias.S("alias").Data().(string)
		indexName := alias.S("index").Data().(string)

		if aliasName == r.options.Name {
			indices = append(indices, indexName)
		}

		if len(indices) > 1 {
			return errors.New("there can only be one index associated to the alias of the index to be rolled over")
		}
	}

	if len(indices) == 0 {
		return errors.New("no results found for alias: " + r.options.Name)
	}

	if len(r.options.NewIndex) > 0 {
		return nil
	}

	last2chars := indices[0][len(indices[0])-2:]

	if string(last2chars[0]) == "-" && isDigit(string(last2chars[1])) {
		return nil
	}

	if isDigit(string(last2chars[0])) && isDigit(string(last2chars[1])) {
		return nil
	}

	return errors.New("one of the last 2 characters of the index is not a digit or a hyphen followed by a digit")
}

func (r *rollover) generateNewIndexName() string {
	return fmt.Sprintf("github.com/alejandro-carstens/scrubber-rollover-%v", time.Now().Format(time.RFC3339))
}
