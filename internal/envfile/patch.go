package envfile

import "fmt"

// PatchOp represents a single patch operation type.
type PatchOp string

const (
	PatchSet    PatchOp = "set"
	PatchDelete PatchOp = "delete"
	PatchRename PatchOp = "rename"
)

// PatchInstruction describes one operation to apply to an env file.
type PatchInstruction struct {
	Op      PatchOp
	Key     string
	Value   string // used by PatchSet
	NewKey  string // used by PatchRename
	Comment string
}

// PatchResult holds the outcome of a single patch instruction.
type PatchResult struct {
	Instruction PatchInstruction
	Applied     bool
	Reason      string
}

// Patch applies a slice of PatchInstructions to a set of entries.
// It returns the modified entries and a result report for each instruction.
func Patch(entries []Entry, instructions []PatchInstruction) ([]Entry, []PatchResult, error) {
	results := make([]PatchResult, 0, len(instructions))
	working := make([]Entry, len(entries))
	copy(working, entries)

	for _, inst := range instructions {
		var res PatchResult
		res.Instruction = inst

		switch inst.Op {
		case PatchSet:
			if inst.Key == "" {
				return nil, nil, fmt.Errorf("patch set: key must not be empty")
			}
			working, res.Applied = applySet(working, inst.Key, inst.Value, inst.Comment)
			res.Reason = "key set"

		case PatchDelete:
			if inst.Key == "" {
				return nil, nil, fmt.Errorf("patch delete: key must not be empty")
			}
			working, res.Applied = applyDelete(working, inst.Key)
			if !res.Applied {
				res.Reason = "key not found"
			} else {
				res.Reason = "key deleted"
			}

		case PatchRename:
			if inst.Key == "" || inst.NewKey == "" {
				return nil, nil, fmt.Errorf("patch rename: both key and new_key must be set")
			}
			var err error
			working, err = applyRenameOp(working, inst.Key, inst.NewKey)
			if err != nil {
				res.Reason = err.Error()
			} else {
				res.Applied = true
				res.Reason = "key renamed"
			}

		default:
			return nil, nil, fmt.Errorf("unknown patch op: %q", inst.Op)
		}

		results = append(results, res)
	}

	return working, results, nil
}

func applySet(entries []Entry, key, value, comment string) ([]Entry, bool) {
	for i, e := range entries {
		if e.Key == key {
			entries[i].Value = value
			if comment != "" {
				entries[i].Comment = comment
			}
			return entries, true
		}
	}
	return append(entries, Entry{Key: key, Value: value, Comment: comment}), true
}

func applyDelete(entries []Entry, key string) ([]Entry, bool) {
	for i, e := range entries {
		if e.Key == key {
			return append(entries[:i], entries[i+1:]...), true
		}
	}
	return entries, false
}

func applyRenameOp(entries []Entry, oldKey, newKey string) ([]Entry, error) {
	found := false
	for _, e := range entries {
		if e.Key == newKey {
			return entries, fmt.Errorf("rename conflict: key %q already exists", newKey)
		}
		if e.Key == oldKey {
			found = true
		}
	}
	if !found {
		return entries, fmt.Errorf("key %q not found", oldKey)
	}
	for i, e := range entries {
		if e.Key == oldKey {
			entries[i].Key = newKey
			break
		}
	}
	return entries, nil
}
