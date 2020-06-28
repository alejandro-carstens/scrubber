package accesscontrol

func inUint64Slice(needle uint64, haystack []uint64) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
