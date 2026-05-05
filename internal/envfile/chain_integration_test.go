package envfile_test

import (
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestChain_ParseNormalizeFilterRoundTrip(t *testing.T) {
	raw := `
debug_mode = true
  API_KEY=secret123  
DB_HOST=localhost
PREFIX_ONLY=1
`
	entries, err := envfile.Parse(strings.NewReader(raw))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	steps := []envfile.ChainStep{
		{
			Name: "normalize",
			Fn: func(e []envfile.Entry) ([]envfile.Entry, error) {
				return envfile.Normalize(e, envfile.DefaultNormalizeOptions())
			},
		},
		{
			Name: "filter-prefix",
			Fn: func(e []envfile.Entry) ([]envfile.Entry, error) {
				return envfile.Filter(e, envfile.FilterOptions{Prefixes: []string{"DB_", "API_"}})
			},
		},
	}

	r, err := envfile.Chain(entries, steps, envfile.ChainOptions{StopOnError: true})
	if err != nil {
		t.Fatalf("chain: %v", err)
	}
	if r.HasErrors() {
		t.Error("expected no step errors")
	}

	keys := make(map[string]bool)
	for _, e := range r.Entries {
		keys[e.Key] = true
	}
	if !keys["API_KEY"] {
		t.Error("expected API_KEY in result")
	}
	if !keys["DB_HOST"] {
		t.Error("expected DB_HOST in result")
	}
	if keys["DEBUG_MODE"] {
		t.Error("DEBUG_MODE should have been filtered out")
	}

	summary := envfile.ChainSummary(r)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
	formatted := envfile.FormatChainSteps(r)
	if formatted == "" {
		t.Error("expected non-empty step format")
	}
}
