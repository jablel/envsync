package envfile

import "fmt"

// SnapshotDiffResult holds the comparison between two snapshots.
type SnapshotDiffResult struct {
	FromLabel string
	ToLabel   string
	Added     []Entry
	Removed   []Entry
	Modified  []SnapshotChange
	Unchanged []Entry
}

// SnapshotChange records a value change between snapshots.
type SnapshotChange struct {
	Key      string
	OldValue string
	NewValue string
}

// DiffSnapshots compares two snapshots and returns what changed.
func DiffSnapshots(from, to Snapshot) SnapshotDiffResult {
	result := SnapshotDiffResult{
		FromLabel: from.Label,
		ToLabel:   to.Label,
	}

	fromMap := from.ToMap()
	toMap := to.ToMap()

	for k, toVal := range toMap {
		if fromVal, exists := fromMap[k]; !exists {
			result.Added = append(result.Added, Entry{Key: k, Value: toVal})
		} else if fromVal != toVal {
			result.Modified = append(result.Modified, SnapshotChange{
				Key: k, OldValue: fromVal, NewValue: toVal,
			})
		} else {
			result.Unchanged = append(result.Unchanged, Entry{Key: k, Value: toVal})
		}
	}

	for k, fromVal := range fromMap {
		if _, exists := toMap[k]; !exists {
			result.Removed = append(result.Removed, Entry{Key: k, Value: fromVal})
		}
	}

	sortEntries(result.Added)
	sortEntries(result.Removed)
	sortEntries(result.Unchanged)
	sortChanges(result.Modified)
	return result
}

// HasSnapshotChanges returns true if the diff contains any additions, removals, or modifications.
func HasSnapshotChanges(d SnapshotDiffResult) bool {
	return len(d.Added)+len(d.Removed)+len(d.Modified) > 0
}

// Summary returns a human-readable summary of the snapshot diff.
func (d SnapshotDiffResult) Summary() string {
	return fmt.Sprintf(
		"snapshot diff [%s -> %s]: +%d added, -%d removed, ~%d modified, %d unchanged",
		d.FromLabel, d.ToLabel,
		len(d.Added), len(d.Removed), len(d.Modified), len(d.Unchanged),
	)
}

func sortChanges(changes []SnapshotChange) {
	for i := 1; i < len(changes); i++ {
		for j := i; j > 0 && changes[j].Key < changes[j-1].Key; j-- {
			changes[j], changes[j-1] = changes[j-1], changes[j]
		}
	}
}
