package actions

import (
	"encoding/json"
	"errors"
	"scrubber/actions/contexts"
	"scrubber/logger"
	"strconv"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
)

const SNAPSHOT_ACTION_TYPE string = "snapshot"
const INDEX_ACTION_TYPE string = "index"
const SECONDS_IN_A_DAY int64 = 86400
const SECONDS_IN_A_MONTH int64 = 2628000
const SECONDS_IN_A_YEAR int64 = 31540000

func Create(context contexts.Contextable, logger *logger.Logger, connection *golastic.Connection) (Actionable, error) {
	action, err := build(context.Action())

	if err != nil {
		return nil, err
	}

	if err := action.Init(context, logger, connection); err != nil {
		return nil, err
	}

	if err := action.ApplyOptions().ApplyFilters(); err != nil {
		return nil, err
	}

	return action, nil
}

func newReporter(logger *logger.Logger) *reporter {
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
		errorReport.actionType = SNAPSHOT_ACTION_TYPE
	} else {
		errorReport.actionType = INDEX_ACTION_TYPE
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
	case "rollover":
		action = new(rollover)
		break
	case "list_indices":
		action = new(listIndices)
		break
	case "list_snapshots":
		action = new(listSnapshots)
		break
	case "delete_repositories":
		action = new(deleteRepositories)
		break
	case "watch":
		action = new(watch)
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
	case "list_snapshots":
		return true
	}

	return false
}

func isDigit(digit string) bool {
	_, err := strconv.Atoi(digit)

	return err == nil
}

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}

func intervalToSeconds(interval int64, unit string) int64 {
	switch unit {
	case "minutes":
		return interval * 60
	case "hours":
		return interval * 3600
	case "days":
		return interval * SECONDS_IN_A_DAY
	case "months":
		return interval * SECONDS_IN_A_MONTH
	case "years":
		return interval * SECONDS_IN_A_YEAR
	}

	return interval
}
