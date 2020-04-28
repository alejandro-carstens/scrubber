package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/alejandro-carstens/golastic"
	"github.com/alejandro-carstens/scrubber/actions/contexts"
	"github.com/alejandro-carstens/scrubber/actions/options"
	"github.com/alejandro-carstens/scrubber/filesystem"
	"github.com/alejandro-carstens/scrubber/logger"
	"github.com/alejandro-carstens/scrubber/notifications"
)

type timer struct {
	done  bool
	timer *time.Timer
}

func (t *timer) start(seconds int64) *timer {
	t.timer = time.NewTimer(time.Duration(seconds) * time.Second)

	return t
}

func (t *timer) expired() bool {
	if t.done {
		return t.done
	}

	select {
	case <-t.timer.C:
		t.done = true

		return t.done
	default:
		return false
	}
}

// SNAPSHOT_ACTION_TYPE (self explanatory)
const SNAPSHOT_ACTION_TYPE string = "snapshot"

// INDEX_ACTION_TYPE (self explanatory)
const INDEX_ACTION_TYPE string = "index"

// SECONDS_IN_A_DAY (self explanatory)
const SECONDS_IN_A_DAY int64 = 86400

// SECONDS_IN_A_MONTH (self explanatory)
const SECONDS_IN_A_MONTH int64 = 2628000

// SECONDS_IN_A_YEAR (self explanatory)
const SECONDS_IN_A_YEAR int64 = 31540000

var availableNumericTypes []string = []string{
	"long",
	"integer",
	"short",
	"byte",
	"double",
	"float",
	"half_float",
	"scaled_float",
}

// Create builds an Actionable action
func Create(
	context contexts.Contextable,
	logger *logger.Logger,
	connection *golastic.Connection,
	queue *notifications.Queue,
	ctx context.Context,
) (Actionable, error) {
	action, err := build(context.Action())

	if err != nil {
		return nil, err
	}

	if err := action.Init(context, logger, connection, queue, ctx); err != nil {
		return nil, err
	}

	if err := action.ApplyOptions().ApplyFilters(); err != nil {
		return nil, err
	}

	return action, nil
}

func newReporter(logger *logger.Logger) *reporter {
	r := new(reporter)
	r.logger = logger

	return r
}

func newErrorReport(action string, name string, err error) *errorReport {
	errorReport := new(errorReport)

	errorReport.errors = append([]error{}, err)
	errorReport.name = name
	errorReport.action = action

	if isSnapshotAction(action) {
		errorReport.actionType = SNAPSHOT_ACTION_TYPE
	} else {
		errorReport.actionType = INDEX_ACTION_TYPE
	}

	return errorReport
}

func newErrorContainer() *errorContainer {
	errorContainer := new(errorContainer)

	errorContainer.reports = map[string]*errorReport{}

	return errorContainer
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
	case "mutate":
		action = new(mutate)
		break
	case "dump":
		action = new(dump)
		break
	case "import_dump":
		action = new(importDump)
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

func buildQuery(builder *golastic.Builder, queryCriteria []*options.QueryCriteria) {
	for _, criteria := range queryCriteria {
		switch criteria.Clause {
		case "where":
			builder.Where(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "where_nested":
			builder.WhereNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "where_in":
			builder.WhereIn(criteria.Key, criteria.Values)
			break
		case "where_in_nested":
			builder.WhereInNested(criteria.Key, criteria.Values)
			break
		case "where_not_in":
			builder.WhereNotIn(criteria.Key, criteria.Values)
			break
		case "where_not_in_nested":
			builder.WhereNotInNested(criteria.Key, criteria.Values)
			break
		case "filter":
			builder.Filter(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "filter_nested":
			builder.FilterNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "filter_in":
			builder.FilterIn(criteria.Key, criteria.Values)
			break
		case "filter_in_nested":
			builder.FilterInNested(criteria.Key, criteria.Values)
			break
		case "match":
			builder.Match(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_nested":
			builder.MatchNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_in":
			builder.MatchIn(criteria.Key, criteria.Values)
			break
		case "match_in_nested":
			builder.MatchInNested(criteria.Key, criteria.Values)
			break
		case "match_not_in":
			builder.MatchNotIn(criteria.Key, criteria.Values)
			break
		case "match_not_in_nested":
			builder.MatchNotInNested(criteria.Key, criteria.Values)
			break
		case "match_phrase":
			builder.MatchPhrase(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_phrase_nested":
			builder.MatchPhraseNested(criteria.Key, criteria.Operator, criteria.Value)
			break
		case "match_phrase_in":
			builder.MatchPhraseIn(criteria.Key, criteria.Values)
			break
		case "match_phrase_in_nested":
			builder.MatchPhraseInNested(criteria.Key, criteria.Values)
			break
		case "match_phrase_not_in":
			builder.MatchPhraseNotIn(criteria.Key, criteria.Values)
			break
		case "match_phrase_not_in_nested":
			builder.MatchPhraseNotInNested(criteria.Key, criteria.Values)
			break
		case "limit":
			builder.Limit(criteria.Limit)
			break
		case "order_by":
			builder.OrderBy(criteria.Key, criteria.Order)
			break
		case "order_by_nested":
			builder.OrderByNested(criteria.Key, criteria.Order)
			break
		}
	}
}

func fileToJSON(fs filesystem.Storeable, path string) (*gabs.Container, error) {
	r, err := fs.Open(path)

	if err != nil {
		return nil, err
	}

	if _, isFile := r.(*os.File); isFile {
		defer r.(*os.File).Close()
	}

	b := bytes.NewBuffer(nil)

	if _, err := io.Copy(b, r); err != nil {
		return nil, err
	}

	return gabs.ParseJSON(b.Bytes())
}

func extractSource(s string) (map[string]interface{}, error) {
	data := map[string]interface{}{}

	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil, err
	}

	id, valid := data["_id"].(string)

	if !valid {
		return nil, errors.New("could not extract id from document.")
	}

	source, valid := data["_source"].(map[string]interface{})

	if !valid {
		return nil, errors.New("could not retrive source from document")
	}

	if _, valid := source["id"].(string); !valid {
		source["id"] = id
	}

	return source, nil
}

func addToMap(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	for key, value := range m2 {
		if value == nil {
			m1[key] = map[string]interface{}{}

			continue
		}

		m1[key] = value
	}

	return m1
}

func removeFromMap(m map[string]interface{}, keys ...string) map[string]interface{} {
	for _, key := range keys {
		delete(m, key)
	}

	return m
}
