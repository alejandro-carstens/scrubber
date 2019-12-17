package actions

import (
	"errors"
	"scrubber/actions/options"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
)

type alias struct {
	filterAction
	options *options.AliasOptions
}

// ApplyOptions implementation of the Actionable interface
func (a *alias) ApplyOptions() Actionable {
	a.options = a.context.Options().(*options.AliasOptions)
	a.indexer.SetOptions(&golastic.IndexOptions{Timeout: a.options.TimeoutInSeconds()})

	return a
}

// Perform implementation of the Actionable interface
func (a *alias) Perform() Actionable {
	a.exec(func(index string) error {
		if a.options.Type == "add" {
			return a.add(index)
		}

		return a.remove(index)
	})

	return a
}

func (a *alias) add(index string) error {
	if a.options.ExtraSettings == nil {
		return a.checkResponse(a.indexer.AddAlias(index, a.options.Name))
	}

	aliasAddAction := a.indexer.AliasAddAction(a.options.Name).Index(index)

	if len(a.options.ExtraSettings.Routing) > 0 {
		aliasAddAction.Routing(a.options.ExtraSettings.Routing)
	}

	if len(a.options.ExtraSettings.SearchRouting) > 0 {
		aliasAddAction.SearchRouting(strings.Split(a.options.ExtraSettings.SearchRouting, ",")...)
	}

	if a.options.ExtraSettings.Filter != nil {
		filter, err := mapToString(a.options.ExtraSettings.Filter)

		if err != nil {
			return err
		}

		aliasAddAction.Filter(a.connection.Builder(index).RawQuery(filter))
	}

	return a.checkResponse(a.indexer.AddAliasByAction(aliasAddAction))
}

func (a *alias) remove(index string) error {
	return a.checkResponse(a.indexer.RemoveIndexFromAlias(index, a.options.Name))
}

func (a *alias) checkResponse(response *gabs.Container, err error) error {
	if err != nil {
		return err
	}

	if acknowledged, _ := response.S("acknowledged").Data().(bool); !acknowledged {
		return errors.New("add alias action was not acknowledged")
	}

	return nil
}
