package envfile

import "testing"

func TestMasker_IsSensitive(t *testing.T) {
	m := NewMasker()

	sensitive := []string{
		"DB_PASSWORD",
		"API_KEY",
		"ACCESS_TOKEN",
		"AWS_SECRET_ACCESS_KEY",
		"GITHUB_AUTH_TOKEN",
	}
	for _, key := range sensitive {
		if !m.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}

	public := []string{
		"APP_NAME",
		"PORT",
		"LOG_LEVEL",
		"DEBUG",
	}
	for _, key := range public {
		if m.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestMasker_MaskValue(t *testing.T) {
	m := NewMasker()

	got := m.MaskValue("DB_PASSWORD", "supersecret")
	if got != maskedValue {
		t.Errorf("expected masked value, got %q", got)
	}

	got = m.MaskValue("PORT", "8080")
	if got != "8080" {
		t.Errorf("expected original value, got %q", got)
	}
}

func TestMasker_MaskEntries(t *testing.T) {
	m := NewMasker()
	entries := []Entry{
		{Key: "APP_NAME", Value: "envsync"},
		{Key: "DB_PASSWORD", Value: "secret123"},
		{Comment: "# a comment"},
		{Key: "API_KEY", Value: "key-abc"},
	}

	masked := m.MaskEntries(entries)

	if masked[0].Value != "envsync" {
		t.Errorf("APP_NAME should not be masked")
	}
	if masked[1].Value != maskedValue {
		t.Errorf("DB_PASSWORD should be masked")
	}
	if masked[2].Value != "" {
		t.Errorf("comment entry value should remain empty")
	}
	if masked[3].Value != maskedValue {
		t.Errorf("API_KEY should be masked")
	}

	// Ensure original entries are not mutated
	if entries[1].Value != "secret123" {
		t.Errorf("original entry should not be mutated")
	}
}

func TestMasker_CustomPatterns(t *testing.T) {
	m := &Masker{Patterns: []string{"INTERNAL"}}
	if !m.IsSensitive("INTERNAL_KEY") {
		t.Error("expected INTERNAL_KEY to be sensitive with custom pattern")
	}
	if m.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom pattern only")
	}
}
