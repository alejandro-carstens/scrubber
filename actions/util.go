package actions

import (
	"encoding/json"
	"errors"
	"scrubber/actions/contexts"
	"scrubber/logging"

	"github.com/Jeffail/gabs"
)

func Create(context contexts.Contextable, logger *logging.SrvLogger) (Actionable, error) {
	action, err := build(context.Action())

	if err != nil {
		return nil, err
	}

	if err := action.Init(context, logger); err != nil {
		return nil, err
	}

	if err := action.ApplyOptions().ApplyFilters(); err != nil {
		return nil, err
	}

	return action, nil
}

func newReporter(logger *logging.SrvLogger) *reporter {
	return &reporter{
		logger: logger,
	}
}

func newErrorReport(action string, name string, err error) *errorReport {
	errorReport := new(errorReport)

	errorReport.errs = append([]error{}, err)
	errorReport.name = name
	errorReport.action = action

	if isSnapshotAction(action) {
		errorReport.actionType = "snapshot"
	} else {
		errorReport.actionType = "index"
	}

	return errorReport
}

func newErrorReportMap() *errorReportMap {
	errorReportMap := new(errorReportMap)

	errorReportMap.reports = map[string]*errorReport{}

	return errorReportMap
}

func build(name string) (Actionable, error) {
	var action Actionable

	switch name {
	case "create_index":
		action = new(createIndex)
		break
	case "delete_indices":
		action = new(deleteIndices)
		break
	case "snapshot":
		action = new(snapshot)
		break
	case "create_repository":
		action = new(createRepository)
		break
	case "open_indices":
		action = new(openIndices)
		break
	case "close_indices":
		action = new(closeIndices)
		break
	case "delete_snapshots":
		action = new(deleteSnapshots)
		break
	case "index_settings":
		action = new(indexSettings)
		break
	case "alias":
		action = new(alias)
		break
	case "restore":
		action = new(restore)
		break
	default:
		return nil, errors.New("Invalid action type")
	}

	return action, nil
}

func containerToMap(container *gabs.Container) (map[string]interface{}, error) {
	res := map[string]interface{}{}

	if err := json.Unmarshal([]byte(container.String()), &res); err != nil {
		return nil, err
	}

	return res, nil
}

func mapToString(val map[string]interface{}) (string, error) {
	b, err := json.Marshal(val)

	if err != nil {
		return "", err
	}

	container, err := gabs.ParseJSON(b)

	if err != nil {
		return "", err
	}

	return container.String(), nil
}

func isSnapshotAction(action string) bool {
	switch action {
	case "restore":
		return true
	case "delete_snapshots":
		return true
	}

	return false
}
