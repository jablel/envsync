package envfile

import (
	"strings"
	"testing"
)

func tagEntries() []Entry {
	return []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "secret", Comment: "# existing comment"},
		{Key: "PORT", Value: "8080"},
	}
}

func TestTag_AppliesTagToAllEntries(t *testing.T) {
	entries := tagEntries()
	out, results, err := Tag(entries, TagOptions{Tags: []string{"env:prod"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, e := range out {
		if !strings.Contains(e.Comment, "# @env:prod") {
			t.Errorf("key %s missing tag, comment: %q", e.Key, e.Comment)
		}
	}
}

func TestTag_AppliesTagToSpecificKeys(t *testing.T) {
	out, _, err := Tag(tagEntries(), TagOptions{
		Keys: []string{"DB_HOST"},
		Tags: []string{"team:infra"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out[0].Comment, "# @team:infra") {
		t.Errorf("DB_HOST should be tagged")
	}
	if strings.Contains(out[1].Comment, "# @team:infra") {
		t.Errorf("API_KEY should not be tagged")
	}
}

func TestTag_SkipsExistingTagWithoutOverwrite(t *testing.T) {
	entries := tagEntries()
	entries[0].Comment = "# @env:prod"
	_, results, err := Tag(entries, TagOptions{
		Keys: []string{"DB_HOST"},
		Tags: []string{"env:prod"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 || len(results[0].Skipped) == 0 {
		t.Errorf("expected tag to be skipped")
	}
}

func TestTag_OverwriteReplacesExistingTag(t *testing.T) {
	entries := tagEntries()
	entries[0].Comment = "# @env:staging"
	out, results, err := Tag(entries, TagOptions{
		Keys:      []string{"DB_HOST"},
		Tags:      []string{"env:prod"},
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out[0].Comment, "env:staging") {
		t.Errorf("old tag should be replaced")
	}
	if !strings.Contains(out[0].Comment, "env:prod") {
		t.Errorf("new tag should be present")
	}
	if len(results) == 0 || len(results[0].Tagged) == 0 {
		t.Errorf("expected tagged result")
	}
}

func TestTag_EmptyTagsReturnsError(t *testing.T) {
	_, _, err := Tag(tagEntries(), TagOptions{})
	if err == nil {
		t.Error("expected error for empty tags")
	}
}

func TestTaggedKeys_ReturnsMatchingKeys(t *testing.T) {
	entries := tagEntries()
	entries[0].Comment = "# @env:prod"
	entries[2].Comment = "# @env:prod"
	keys := TaggedKeys(entries, "env:prod")
	if len(keys) != 2 {
		t.Errorf("expected 2 tagged keys, got %d", len(keys))
	}
}
