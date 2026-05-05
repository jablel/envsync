package envfile

import (
	"errors"
	"testing"
)

func makeEntries(kvs ...string) []Entry {
	var out []Entry
	for i := 0; i+1 < len(kvs); i += 2 {
		out = append(out, Entry{Key: kvs[i], Value: kvs[i+1]})
	}
	return out
}

func TestChain_AllStepsApplied(t *testing.T) {
	entries := makeEntries("A", "1", "B", "2")
	steps := []ChainStep{
		{Name: "uppercase", Fn: func(e []Entry) ([]Entry, error) { return UppercaseKeys(e) }},
		{Name: "trim", Fn: func(e []Entry) ([]Entry, error) { return TrimValues(e) }},
	}
	r, err := Chain(entries, steps, ChainOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Steps) != 2 {
		t.Fatalf("expected 2 step results, got %d", len(r.Steps))
	}
	for _, s := range r.Steps {
		if !s.Applied {
			t.Errorf("step %q should be applied", s.Name)
		}
	}
}

func TestChain_StopOnError(t *testing.T) {
	entries := makeEntries("A", "1")
	boom := errors.New("boom")
	steps := []ChainStep{
		{Name: "fail", Fn: func(e []Entry) ([]Entry, error) { return nil, boom }},
		{Name: "should-skip", Fn: func(e []Entry) ([]Entry, error) { return e, nil }},
	}
	r, err := Chain(entries, steps, ChainOptions{StopOnError: true})
	if !errors.Is(err, boom) {
		t.Fatalf("expected boom error, got %v", err)
	}
	if len(r.Steps) != 1 {
		t.Errorf("expected 1 step recorded, got %d", len(r.Steps))
	}
}

func TestChain_ContinueOnError(t *testing.T) {
	entries := makeEntries("A", "1")
	steps := []ChainStep{
		{Name: "fail", Fn: func(e []Entry) ([]Entry, error) { return nil, errors.New("oops") }},
		{Name: "ok", Fn: func(e []Entry) ([]Entry, error) { return e, nil }},
	}
	r, err := Chain(entries, steps, ChainOptions{StopOnError: false})
	if err != nil {
		t.Fatalf("unexpected hard error: %v", err)
	}
	if !r.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
	if len(r.Steps) != 2 {
		t.Errorf("expected 2 steps, got %d", len(r.Steps))
	}
}

func TestChain_EmptySteps(t *testing.T) {
	entries := makeEntries("X", "y")
	r, err := Chain(entries, nil, ChainOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(r.Entries))
	}
}

func TestChainSummary_NoErrors(t *testing.T) {
	r := ChainResult{
		Entries: makeEntries("A", "1"),
		Steps: []ChainStepResult{
			{Name: "s1", Applied: true},
			{Name: "s2", Applied: true},
		},
	}
	s := ChainSummary(r)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
