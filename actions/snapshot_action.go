package actions

type snapshotAction struct {
	filterAction
}

func (sa *snapshotAction) checkForSnapshotsInProgress(repository string) (bool, error) {
	response, err := sa.builder.GetSnapshots(repository, "")

	if err != nil {
		return false, err
	}

	snapshots, err := response.S("snapshots").Children()

	if err != nil {
		return false, err
	}

	if len(snapshots) == 0 {
		return false, nil
	}

	for _, snapshot := range snapshots {
		if snapshot.S("state").Data().(string) == IN_PROGRESS_STATUS {
			return true, nil
		}
	}

	return false, nil
}
