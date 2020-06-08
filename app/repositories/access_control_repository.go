package repositories

import (
	"scrubber/app/models"

	"github.com/jinzhu/gorm"
)

const ACCESS_CONTROL_READ_SCOPE string = "read"
const ACCESS_CONTROL_WRITE_SCOPE string = "write"
const ACCESS_CONTROL_NO_ACCESS_SCOPE string = "no_access"

const ACCESS_CONTROL_ALIAS_ACTION string = "alias"
const ACCESS_CONTROL_CLOSE_INDICES_ACTION string = "close_indices"
const ACCESS_CONTROL_CREATE_INDEX_ACTION string = "create_index"
const ACCESS_CONTROL_CREATE_REPOSITORY_ACTION string = "create_repository"
const ACCESS_CONTROL_DELETE_INDICES_ACTION string = "delete_indices"
const ACCESS_CONTROL_DELETE_REPOSITORY_ACTION string = "delete_repository"
const ACCESS_CONTROL_DELETE_SNAPSHOTS_ACTION string = "delete_snapshots"
const ACCESS_CONTROL_DUMP_ACTION string = "dump"
const ACCESS_CONTROL_IMPORT_DUMP_ACTION string = "import_dump"
const ACCESS_CONTROL_IMDEX_SETTINGS_ACTION string = "index_settings"
const ACCESS_CONTROL_MUTATE_ACTION string = "mutate"
const ACCESS_CONTROL_OPEN_INDICES_ACTION string = "open_indices"
const ACCESS_CONTROL_RESTORE_ACTION string = "restore"
const ACCESS_CONTROL_ROLLOVER_ACTION string = "rollover"
const ACCESS_CONTROL_SNAPSHOT_ACTION string = "snapshot"
const ACCESS_CONTROL_WATCH_ACTION string = "watch"
const ACCESS_CONTROL_ALL_ACTIONS string = "all"

func NewAcessControlRepository() *AccessControlRepository {
	return repo(&models.AccessControl{}, nil).(*AccessControlRepository)
}

type AccessControlRepository struct {
	repository
}

func (acr *AccessControlRepository) FromTx(tx *gorm.DB) *AccessControlRepository {
	return repo(&models.AccessControl{}, tx).(*AccessControlRepository)
}
