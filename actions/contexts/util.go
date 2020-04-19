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
		ctx = new(CreateIndexContext)
		break
	case "delete_indices":
		ctx = new(DeleteIndicesContext)
		break
	case "create_repository":
		ctx = new(CreateRepositoryContext)
		break
	case "snapshot":
		ctx = new(SnapshotContext)
		break
	case "rollover":
		ctx = new(RolloverContext)
	case "open_indices":
		ctx = new(OpenIndicesContext)
		break
	case "close_indices":
		ctx = new(CloseIndicesContext)
		break
	case "delete_snapshots":
		ctx = new(DeleteSnapshotsContext)
		break
	case "index_settings":
		ctx = new(IndexSettingsContext)
		break
	case "alias":
		ctx = new(AliasContext)
		break
	case "restore":
		ctx = new(RestoreContext)
		break
	case "list_indices":
		ctx = new(ListIndicesContext)
		break
	case "list_snapshots":
		ctx = new(ListSnapshotsContext)
		break
	case "delete_repositories":
		ctx = new(DeleteRepositoriesContext)
		break
	case "watch":
		ctx = new(WatchContext)
		break
	case "mutate":
		ctx = new(MutateContext)
		break
	case "dump":
		ctx = new(DumpContext)
		break
	case "import_dump":
		ctx = new(ImportDumpContext)
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
	case "list_snapshots":
		return true
	}

	return false
}
