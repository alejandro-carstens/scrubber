package contexts

import (
	"encoding/json"
	"errors"

	"github.com/Jeffail/gabs"
)

const DEFAULT_NUMBER_OF_WORKERS int = 10
const DEFAULT_QUEUE_LENGTH int = 1000

func New(config *gabs.Container) (Contextable, error) {
	action, valid := config.S("action").Data().(string)

	if !valid {
		return nil, errors.New("action not set")
	}

	ctx, err := build(action)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(config.Bytes(), ctx); err != nil {
		return nil, err
	}

	if err := ctx.Config(config); err != nil {
		return nil, err
	}

	if ctx.GetAsync() {
		if ctx.GetNumberOfWorkers() <= 0 {
			ctx.setNumberOfWorkers(DEFAULT_NUMBER_OF_WORKERS)
		}

		if ctx.GetQueueLength() <= 0 {
			ctx.setQueueLength(DEFAULT_QUEUE_LENGTH)
		}
	}

	return ctx, nil
}

func build(action string) (Contextable, error) {
	var ctx Contextable

	switch action {
	case "create_index":
		ctx = new(createIndexContext)
		break
	case "delete_indices":
		ctx = new(DeleteIndicesContext)
		break
	case "create_repository":
		ctx = new(createRepositoryContext)
		break
	case "snapshot":
		ctx = new(snapshotContext)
		break
	case "open_indices":
		ctx = new(openIndicesContext)
		break
	case "close_indices":
		ctx = new(closeIndicesContext)
		break
	case "delete_snapshots":
		ctx = new(deleteSnapshotsContext)
		break
	case "index_settings":
		ctx = new(indexSettingsContext)
		break
	case "alias":
		ctx = new(aliasContext)
		break
	case "restore":
		ctx = new(restoreContext)
		break
	default:
		return nil, errors.New("Invalid action")
	}

	return ctx, nil
}

func isSnapshotAction(action string) bool {
	switch action {
	case "delete_snapshots":
		return true
	case "restore":
		return true
	}

	return false
}
