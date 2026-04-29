package envfile

// StripOptions controls which entries are removed during a strip operation.
type StripOptions struct {
	// RemoveSecrets removes entries whose keys are considered sensitive.
	RemoveSecrets bool
	// RemoveEmpty removes entries with empty values.
	RemoveEmpty bool
	// RemovePrefixes removes entries whose keys start with any of these prefixes.
	RemovePrefixes []string
	// RemoveKeys removes entries with exactly these key names.
	RemoveKeys []string
}

// StripResult holds the outcome of a Strip operation.
type StripResult struct {
	Kept    []Entry
	Removed []Entry
}

// Strip removes entries from the list according to the provided options.
func Strip(entries []Entry, opts StripOptions) StripResult {
	masker := NewMasker()
	keySet := make(map[string]struct{}, len(opts.RemoveKeys))
	for _, k := range opts.RemoveKeys {
		keySet[k] = struct{}{}
	}

	var result StripResult
	for _, e := range entries {
		if shouldStrip(e, opts, masker, keySet) {
			result.Removed = append(result.Removed, e)
		} else {
			result.Kept = append(result.Kept, e)
		}
	}
	return result
}

// StrippedKeys returns just the keys that were removed.
func StrippedKeys(result StripResult) []string {
	keys := make([]string, 0, len(result.Removed))
	for _, e := range result.Removed {
		keys = append(keys, e.Key)
	}
	return keys
}

func shouldStrip(e Entry, opts StripOptions, masker *Masker, keySet map[string]struct{}) bool {
	if _, ok := keySet[e.Key]; ok {
		return true
	}
	if opts.RemoveSecrets && masker.IsSensitive(e.Key) {
		return true
	}
	if opts.RemoveEmpty && e.Value == "" {
		return true
	}
	for _, prefix := range opts.RemovePrefixes {
		if len(e.Key) >= len(prefix) && e.Key[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
