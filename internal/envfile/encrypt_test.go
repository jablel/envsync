package envfile

import (
	"strings"
	"testing"
)

func TestEncrypt_BasicValues(t *testing.T) {
	entries := []Entry{
		{Key: "DB_HOST", Value: "localhost"},
		{Key: "DB_PASS", Value: "s3cr3t"},
	}
	res, err := Encrypt(entries, EncryptOptions{Passphrase: "mypassphrase"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Encrypted) != 2 {
		t.Fatalf("expected 2 encrypted, got %d", len(res.Encrypted))
	}
	for _, e := range res.Encrypted {
		if !strings.HasPrefix(e.Value, encPrefix) {
			t.Errorf("key %s value should start with %q, got %q", e.Key, encPrefix, e.Value)
		}
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := Encrypt([]Entry{{Key: "K", Value: "v"}}, EncryptOptions{})
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestEncrypt_SkipsAlreadyEncrypted(t *testing.T) {
	entries := []Entry{
		{Key: "TOKEN", Value: encPrefix + "alreadyencoded"},
	}
	res, err := Encrypt(entries, EncryptOptions{Passphrase: "pass"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Fatalf("expected 1 skipped, got %d", len(res.Skipped))
	}
	if len(res.Encrypted) != 0 {
		t.Fatalf("expected 0 encrypted, got %d", len(res.Encrypted))
	}
}

func TestEncrypt_OnlySensitive(t *testing.T) {
	m := NewMasker(nil)
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_PASSWORD", Value: "secret"},
	}
	res, err := Encrypt(entries, EncryptOptions{
		Passphrase:    "pass",
		OnlySensitive: true,
		Masker:        m,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Encrypted) != 1 {
		t.Fatalf("expected 1 encrypted, got %d", len(res.Encrypted))
	}
	if res.Encrypted[0].Key != "DB_PASSWORD" {
		t.Errorf("expected DB_PASSWORD to be encrypted")
	}
	if len(res.Skipped) != 1 || res.Skipped[0].Key != "APP_NAME" {
		t.Errorf("expected APP_NAME to be skipped")
	}
}

func TestDecrypt_RoundTrip(t *testing.T) {
	orig := []Entry{
		{Key: "SECRET", Value: "topsecret"},
		{Key: "NAME", Value: "alice"},
	}
	pass := "roundtrippass"
	res, err := Encrypt(orig, EncryptOptions{Passphrase: pass})
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	decrypted, err := Decrypt(res.Encrypted, pass)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	for i, e := range decrypted {
		if e.Value != orig[i].Value {
			t.Errorf("key %s: want %q got %q", e.Key, orig[i].Value, e.Value)
		}
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	orig := []Entry{{Key: "K", Value: "val"}}
	res, err := Encrypt(orig, EncryptOptions{Passphrase: "correct"})
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	_, err = Decrypt(res.Encrypted, "wrong")
	if err == nil {
		t.Fatal("expected decryption error with wrong passphrase")
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	_, err := Decrypt([]Entry{{Key: "K", Value: encPrefix + "data"}}, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_LeavesPlaintextUnchanged(t *testing.T) {
	entries := []Entry{{Key: "HOST", Value: "localhost"}}
	out, err := Decrypt(entries, "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out[0].Value != "localhost" {
		t.Errorf("expected unchanged value, got %q", out[0].Value)
	}
}
