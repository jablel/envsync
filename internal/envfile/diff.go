package envfile

// DiffStatus represents the type of difference between two env files.
type DiffStatus string

const (
	StatusAdded    DiffStatus = "added"
	StatusRemoved  DiffStatus = "removed"
	StatusModified DiffStatus = "modified"
	StatusUnchanged DiffStatus = "unchanged"
)

// DiffEntry represents a single key difference between two env maps.
type DiffEntry struct {
	Key      string
	Status   DiffStatus
	OldValue string
	NewValue string
}

// Diff compares two env maps (base vs target) and returns a slice of DiffEntry.
// Keys present only in base are "removed", only in target are "added",
// present in both but with different values are "modified", otherwise "unchanged".
func Diff(base, target map[string]string) []DiffEntry {
	seen := make(map[string]bool)
	var entries []DiffEntry

	for k, baseVal := range base {
		seen[k] = true
		if targetVal, ok := target[k]; !ok {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusRemoved,
				OldValue: baseVal,
				NewValue: "",
			})
		} else if baseVal != targetVal {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusModified,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		} else {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusUnchanged,
				OldValue: baseVal,
				NewValue: targetVal,
			})
		}
	}

	for k, targetVal := range target {
		if !seen[k] {
			entries = append(entries, DiffEntry{
				Key:      k,
				Status:   StatusAdded,
				OldValue: "",
				NewValue: targetVal,
			})
		}
	}

	return entries
}

// HasChanges returns true if any entry in the diff is not unchanged.
func HasChanges(entries []DiffEntry) bool {
	for _, e := range entries {
		if e.Status != StatusUnchanged {
			return true
		}
	}
	return false
}
