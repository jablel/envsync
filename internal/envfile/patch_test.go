package envfile

import (
	"testing"
)

func TestPatch_SetNewKey(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	instructions := []PatchInstruction{
		{Op: PatchSet, Key: "NEW_KEY", Value: "new_value"},
	}
	out, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if !results[0].Applied {
		t.Errorf("expected instruction to be applied")
	}
}

func TestPatch_SetExistingKey(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "old"}}
	instructions := []PatchInstruction{
		{Op: PatchSet, Key: "FOO", Value: "updated"},
	}
	out, _, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "updated" {
		t.Errorf("expected value 'updated', got %q", out[0].Value)
	}
}

func TestPatch_DeleteKey(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}}
	instructions := []PatchInstruction{
		{Op: PatchDelete, Key: "FOO"},
	}
	out, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Key != "BAZ" {
		t.Errorf("expected only BAZ remaining, got %+v", out)
	}
	if !results[0].Applied {
		t.Errorf("expected delete to be applied")
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	instructions := []PatchInstruction{
		{Op: PatchDelete, Key: "MISSING"},
	}
	_, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected delete to not be applied for missing key")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	entries := []Entry{{Key: "OLD_NAME", Value: "value"}}
	instructions := []PatchInstruction{
		{Op: PatchRename, Key: "OLD_NAME", NewKey: "NEW_NAME"},
	}
	out, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Key != "NEW_NAME" {
		t.Errorf("expected key NEW_NAME, got %q", out[0].Key)
	}
	if !results[0].Applied {
		t.Errorf("expected rename to be applied")
	}
}

func TestPatch_RenameConflict(t *testing.T) {
	entries := []Entry{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}}
	instructions := []PatchInstruction{
		{Op: PatchRename, Key: "A", NewKey: "B"},
	}
	_, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Applied {
		t.Errorf("expected rename to fail on conflict")
	}
}

func TestPatch_InvalidOp(t *testing.T) {
	entries := []Entry{{Key: "FOO", Value: "bar"}}
	instructions := []PatchInstruction{
		{Op: "unknown", Key: "FOO"},
	}
	_, _, err := Patch(entries, instructions)
	if err == nil {
		t.Error("expected error for unknown op")
	}
}

func TestPatch_MultipleInstructions(t *testing.T) {
	entries := []Entry{
		{Key: "A", Value: "1"},
		{Key: "B", Value: "2"},
	}
	instructions := []PatchInstruction{
		{Op: PatchSet, Key: "C", Value: "3"},
		{Op: PatchDelete, Key: "A"},
		{Op: PatchRename, Key: "B", NewKey: "B_RENAMED"},
	}
	out, results, err := Patch(entries, instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for _, r := range results {
		if !r.Applied {
			t.Errorf("expected instruction %v to be applied", r.Instruction)
		}
	}
}
