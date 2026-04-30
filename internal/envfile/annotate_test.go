package envfile

import (
	"testing"
)

func TestAnnotate_AppliesComment(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "API_KEY", Value: "secret"},
	}
	annotations := []Annotation{
		{Key: "DB_HOST", Comment: "database hostname"},
	}
	out, results, err := Annotate(entries, annotations, AnnotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Comment != "database hostname" {
		t.Errorf("expected comment 'database hostname', got %q", out[0].Comment)
	}
	if len(results) != 1 || !results[0].Applied {
		t.Errorf("expected one applied result")
	}
}

func TestAnnotate_AppliesTags(t *testing.T) {
	entries := []Entry{
		{Key: "API_KEY", Value: "abc123"},
	}
	annotations := []Annotation{
		{Key: "API_KEY", Tags: []string{"secret", "rotation"}},
	}
	out, _, err := Annotate(entries, annotations, AnnotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Comment != "tags:secret,rotation" {
		t.Errorf("unexpected comment: %q", out[0].Comment)
	}
}

func TestAnnotate_SkipsExistingWithoutOverwrite(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost", Comment: "existing comment"},
	}
	annotations := []Annotation{
		{Key: "DB_HOST", Comment: "new comment"},
	}
	out, results, err := Annotate(entries, annotations, AnnotateOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Comment != "existing comment" {
		t.Errorf("expected original comment preserved, got %q", out[0].Comment)
	}
	if len(results) != 1 || !results[0].Skipped {
		t.Errorf("expected skipped result")
	}
}

func TestAnnotate_OverwriteExisting(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost", Comment: "old"},
	}
	annotations := []Annotation{
		{Key: "DB_HOST", Comment: "new comment"},
	}
	out, results, err := Annotate(entries, annotations, AnnotateOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Comment != "new comment" {
		t.Errorf("expected 'new comment', got %q", out[0].Comment)
	}
	if len(results) != 1 || !results[0].Applied {
		t.Errorf("expected applied result")
	}
}

func TestAnnotate_EmptyKeyError(t *testing.T) {
	entries := []Entry{{Key: "X", Value: "1"}}
	_, _, err := Annotate(entries, []Annotation{{Key: "", Comment: "bad"}}, AnnotateOptions{})
	if err == nil {
		t.Error("expected error for empty annotation key")
	}
}

func TestAnnotate_CombinesCommentAndTags(t *testing.T) {
	entries := []Entry{{Key: "TOKEN", Value: "xyz"}}
	annotations := []Annotation{
		{Key: "TOKEN", Comment: "auth token", Tags: []string{"sensitive"}},
	}
	out, _, err := Annotate(entries, annotations, AnnotateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "auth token | tags:sensitive"
	if out[0].Comment != expected {
		t.Errorf("expected %q, got %q", expected, out[0].Comment)
	}
}

func TestAnnotatedKeys_ReturnsCommented(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1", Comment: "has comment"},
		{Key: "B", Value: "2"},
		{Key: "C", Value: "3", Comment: "also has comment"},
	}
	keys := AnnotatedKeys(entries)
	if len(keys) != 2 {
		t.Errorf("expected 2 annotated keys, got %d", len(keys))
	}
}
