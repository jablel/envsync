package envfile

import (
	"fmt"
	"strings"
)

// PatchSummary returns a human-readable summary of patch results.
func PatchSummary(results []PatchResult) string {
	if len(results) == 0 {
		return "No patch instructions applied."
	}

	var sb strings.Builder
	applied, skipped := 0, 0

	for _, r := range results {
		if r.Applied {
			applied++
		} else {
			skipped++
		}
	}

	fmt.Fprintf(&sb, "Patch summary: %d applied, %d skipped\n", applied, skipped)
	for _, r := range results {
		status := "✓"
		if !r.Applied {
			status = "✗"
		}
		switch r.Instruction.Op {
		case PatchSet:
			fmt.Fprintf(&sb, "  %s set    %s=%s (%s)\n", status, r.Instruction.Key, r.Instruction.Value, r.Reason)
		case PatchDelete:
			fmt.Fprintf(&sb, "  %s delete %s (%s)\n", status, r.Instruction.Key, r.Reason)
		case PatchRename:
			fmt.Fprintf(&sb, "  %s rename %s -> %s (%s)\n", status, r.Instruction.Key, r.Instruction.NewKey, r.Reason)
		}
	}
	return sb.String()
}

// FormatPatchInstructions renders a list of instructions as a readable plan.
func FormatPatchInstructions(instructions []PatchInstruction) string {
	if len(instructions) == 0 {
		return "No instructions."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Patch plan (%d instructions):\n", len(instructions))
	for _, inst := range instructions {
		switch inst.Op {
		case PatchSet:
			fmt.Fprintf(&sb, "  SET    %s = %s\n", inst.Key, inst.Value)
		case PatchDelete:
			fmt.Fprintf(&sb, "  DELETE %s\n", inst.Key)
		case PatchRename:
			fmt.Fprintf(&sb, "  RENAME %s -> %s\n", inst.Key, inst.NewKey)
		}
	}
	return sb.String()
}
