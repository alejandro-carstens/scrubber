package repositories

const ADMIN_ROLE string = "Admin"

const READ_SCOPE string = "read"
const WRITE_SCOPE string = "write"
const NO_ACCESS_SCOPE string = "no_access"

const ALIAS_ACTION string = "alias"
const CLOSE_INDICES_ACTION string = "close_indices"
const CREATE_INDEX_ACTION string = "create_index"
const CREATE_REPOSITORY_ACTION string = "create_repository"
const DELETE_INDICES_ACTION string = "delete_indices"
const DELETE_REPOSITORY_ACTION string = "delete_repository"
const DELETE_SNAPSHOTS_ACTION string = "delete_snapshots"
const DUMP_ACTION string = "dump"
const IMPORT_DUMP_ACTION string = "import_dump"
const IMDEX_SETTINGS_ACTION string = "index_settings"
const MUTATE_ACTION string = "mutate"
const OPEN_INDICES_ACTION string = "open_indices"
const RESTORE_ACTION string = "restore"
const ROLLOVER_ACTION string = "rollover"
const SNAPSHOT_ACTION string = "snapshot"
const WATCH_ACTION string = "watch"

var AvailableActions []string = []string{
	ALIAS_ACTION,
	CLOSE_INDICES_ACTION,
	CREATE_INDEX_ACTION,
	CREATE_REPOSITORY_ACTION,
	DELETE_INDICES_ACTION,
	DELETE_REPOSITORY_ACTION,
	DELETE_SNAPSHOTS_ACTION,
	DUMP_ACTION,
	IMPORT_DUMP_ACTION,
	IMDEX_SETTINGS_ACTION,
	MUTATE_ACTION,
	OPEN_INDICES_ACTION,
	RESTORE_ACTION,
	ROLLOVER_ACTION,
	SNAPSHOT_ACTION,
	WATCH_ACTION,
}

var AvailableScopes []string = []string{
	READ_SCOPE,
	WRITE_SCOPE,
	NO_ACCESS_SCOPE,
}
