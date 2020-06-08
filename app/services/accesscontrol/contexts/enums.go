package contexts

import "scrubber/app/repositories"

var availableActions []string = []string{
	repositories.ACCESS_CONTROL_ALIAS_ACTION,
	repositories.ACCESS_CONTROL_CLOSE_INDICES_ACTION,
	repositories.ACCESS_CONTROL_CREATE_INDEX_ACTION,
	repositories.ACCESS_CONTROL_CREATE_REPOSITORY_ACTION,
	repositories.ACCESS_CONTROL_DELETE_INDICES_ACTION,
	repositories.ACCESS_CONTROL_DELETE_REPOSITORY_ACTION,
	repositories.ACCESS_CONTROL_DELETE_SNAPSHOTS_ACTION,
	repositories.ACCESS_CONTROL_DUMP_ACTION,
	repositories.ACCESS_CONTROL_IMPORT_DUMP_ACTION,
	repositories.ACCESS_CONTROL_IMDEX_SETTINGS_ACTION,
	repositories.ACCESS_CONTROL_MUTATE_ACTION,
	repositories.ACCESS_CONTROL_OPEN_INDICES_ACTION,
	repositories.ACCESS_CONTROL_RESTORE_ACTION,
	repositories.ACCESS_CONTROL_ROLLOVER_ACTION,
	repositories.ACCESS_CONTROL_SNAPSHOT_ACTION,
	repositories.ACCESS_CONTROL_WATCH_ACTION,
	repositories.ACCESS_CONTROL_ALL_ACTIONS,
}

var availableScopes []string = []string{
	repositories.ACCESS_CONTROL_READ_SCOPE,
	repositories.ACCESS_CONTROL_WRITE_SCOPE,
	repositories.ACCESS_CONTROL_NO_ACCESS_SCOPE,
}

var availableAllActionScopes []string = []string{
	repositories.ACCESS_CONTROL_WRITE_SCOPE,
	repositories.ACCESS_CONTROL_NO_ACCESS_SCOPE,
}
