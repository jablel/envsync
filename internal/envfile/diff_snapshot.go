package envfile

// DiffSnapshot compares a live set of entries against a previously saved
// snapshot and returns a DiffResult describing what changed.
func DiffSnapshot(current []Entry, snap Snapshot) DiffResult {
	snapshotEntries := snap.Entries
	return Diff(snapshotEntries, current)
}

// DiffSnapshotByPath loads a snapshot from disk and diffs it against the
// provided entries. Returns an error if the snapshot cannot be loaded.
func DiffSnapshotByPath(current []Entry, path string) (DiffResult, error) {
	snap, err := LoadSnapshot(path)
	if err != nil {
		return DiffResult{}, err
	}
	return DiffSnapshot(current, snap), nil
}

// ApplySnapshotDiff applies all changes described by a DiffResult to a base
// set of entries, returning the updated slice. Added and modified entries from
// the diff are merged in; removed entries are dropped.
func ApplySnapshotDiff(base []Entry, d DiffResult) []Entry {
	baseMap := make(map[string]Entry, len(base))
	for _, e := range base {
		baseMap[e.Key] = e
	}

	for _, e := range d.Added {
		baseMap[e.Key] = e
	}
	for _, e := range d.Modified {
		baseMap[e.Key] = e.New
	}
	for _, e := range d.Removed {
		delete(baseMap, e.Key)
	}

	result := make([]Entry, 0, len(baseMap))
	for _, e := range baseMap {
		result = append(result, e)
	}
	sortEntries(result)
	return result
}
