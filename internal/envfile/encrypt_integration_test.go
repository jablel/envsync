package envfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envsync/internal/envfile"
)

func TestEncrypt_ParseEncryptDecryptRoundTrip(t *testing.T) {
	content := `DB_HOST=localhost
DB_PASSWORD=supersecret
APP_ENV=production
`
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	entries, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	pass := "integration-test-pass"
	res, err := envfile.Encrypt(entries, envfile.EncryptOptions{Passphrase: pass})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	if len(res.Encrypted) != 3 {
		t.Fatalf("expected 3 encrypted entries, got %d", len(res.Encrypted))
	}
	for _, e := range res.Encrypted {
		if !strings.HasPrefix(e.Value, "enc:") {
			t.Errorf("entry %s not encrypted", e.Key)
		}
	}

	decrypted, err := envfile.Decrypt(res.Encrypted, pass)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	origMap := make(map[string]string)
	for _, e := range entries {
		origMap[e.Key] = e.Value
	}
	for _, e := range decrypted {
		if origMap[e.Key] != e.Value {
			t.Errorf("key %s: want %q got %q", e.Key, origMap[e.Key], e.Value)
		}
	}
}

func TestEncryptSummary_Output(t *testing.T) {
	entries := []envfile.Entry{
		{Key: "SECRET_KEY", Value: "abc"},
		{Key: "PLAIN", Value: "xyz"},
	}
	res, err := envfile.Encrypt(entries, envfile.EncryptOptions{Passphrase: "pass"})
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	summary := envfile.EncryptSummary(res)
	if !strings.Contains(summary, "Encrypted: 2") {
		t.Errorf("summary missing encrypted count: %q", summary)
	}
	formatted := envfile.FormatEncryptedEntries(res.Encrypted)
	for _, e := range res.Encrypted {
		if !strings.Contains(formatted, e.Key) {
			t.Errorf("formatted output missing key %s", e.Key)
		}
	}
	if !strings.Contains(formatted, "<encrypted>") {
		t.Errorf("formatted output should mask encrypted values")
	}
}
