package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp env: %v", err)
	}
	return path
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnv(t, "APP_NAME=envsync\nPORT=8080\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["APP_NAME"] != "envsync" {
		t.Errorf("expected APP_NAME=envsync, got %q", m["APP_NAME"])
	}
	if m["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", m["PORT"])
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnv(t, `DB_URL="postgres://localhost/dev"
SECRET='mysecret'
`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := env.ToMap()
	if m["DB_URL"] != "postgres://localhost/dev" {
		t.Errorf("unexpected DB_URL: %q", m["DB_URL"])
	}
	if m["SECRET"] != "mysecret" {
		t.Errorf("unexpected SECRET: %q", m["SECRET"])
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnv(t, "# comment\n\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 2 {
		t.Errorf("expected 2 entries (comment + kv), got %d", len(env.Entries))
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnv(t, "BADLINE\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestStripQuotes(t *testing.T) {
	cases := []struct{ in, want string }{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`plain`, "plain"},
		{`"`, `"`},
	}
	for _, c := range cases {
		got := stripQuotes(c.in)
		if got != c.want {
			t.Errorf("stripQuotes(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
