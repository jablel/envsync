package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot captures the state of an env file at a point in time.
type Snapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Label     string           `json:"label"`
	Entries   []Entry          `json:"entries"`
}

// Entry represents a single key-value pair from an env file.
type Entry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// TakeSnapshot creates a Snapshot from a map of parsed env entries.
func TakeSnapshot(label string, entries map[string]string) Snapshot {
	snap := Snapshot{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Entries:   make([]Entry, 0, len(entries)),
	}
	for k, v := range entries {
		snap.Entries = append(snap.Entries, Entry{Key: k, Value: v})
	}
	sortEntries(snap.Entries)
	return snap
}

// SaveSnapshot writes a Snapshot to a JSON file at the given path.
func SaveSnapshot(path string, snap Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// LoadSnapshot reads a Snapshot from a JSON file.
func LoadSnapshot(path string) (Snapshot, error) {
	var snap Snapshot
	f, err := os.Open(path)
	if err != nil {
		return snap, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return snap, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// ToMap converts a Snapshot's entries back to a key-value map.
func (s Snapshot) ToMap() map[string]string {
	m := make(map[string]string, len(s.Entries))
	for _, e := range s.Entries {
		m[e.Key] = e.Value
	}
	return m
}

func sortEntries(entries []Entry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
