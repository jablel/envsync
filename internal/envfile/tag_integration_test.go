package envfile_test

import (
	"strings"
	"testing"

	"github.com/your-org/envsync/internal/envfile"
)

func TestTag_ParseAndTagRoundTrip(t *testing.T) {
	content := `DB_HOST=localhost
API_KEY=supersecret
PORT=8080
`
	tf := writeTempEnv(t, content)

	entries, err := envfile.Parse(tf)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	tagged, results, err := envfile.Tag(entries, envfile.TagOptions{
		Keys: []string{"DB_HOST", "PORT"},
		Tags: []string{"env:prod", "team:infra"},
	})
	if err != nil {
		t.Fatalf("tag error: %v", err)
	}

	sum := envfile.SummarizeTag(results)
	if sum.TotalTagged != 4 { // 2 keys × 2 tags
		t.Errorf("expected 4 tags applied, got %d", sum.TotalTagged)
	}

	// Verify tagged keys are retrievable.
	prodKeys := envfile.TaggedKeys(tagged, "env:prod")
	if len(prodKeys) != 2 {
		t.Errorf("expected 2 prod-tagged keys, got %d", len(prodKeys))
	}

	// API_KEY should be unaffected.
	for _, e := range tagged {
		if e.Key == "API_KEY" && strings.Contains(e.Comment, "env:prod") {
			t.Errorf("API_KEY should not have been tagged")
		}
	}

	// Format and verify output is non-empty.
	out := envfile.FormatTaggedEntries(results)
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
