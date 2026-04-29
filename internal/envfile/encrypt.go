package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// EncryptOptions controls encryption behaviour.
type EncryptOptions struct {
	// Passphrase is used to derive the AES-256 key.
	Passphrase string
	// OnlySensitive limits encryption to sensitive keys.
	OnlySensitive bool
	// Masker is used when OnlySensitive is true.
	Masker *Masker
}

// EncryptResult holds the outcome of an Encrypt call.
type EncryptResult struct {
	Encrypted []Entry
	Skipped   []Entry
}

const encPrefix = "enc:"

// Encrypt encrypts entry values using AES-256-GCM derived from the passphrase.
func Encrypt(entries []Entry, opts EncryptOptions) (EncryptResult, error) {
	if opts.Passphrase == "" {
		return EncryptResult{}, errors.New("encrypt: passphrase must not be empty")
	}
	key := deriveKey(opts.Passphrase)
	var result EncryptResult
	for _, e := range entries {
		if opts.OnlySensitive && opts.Masker != nil && !opts.Masker.IsSensitive(e.Key) {
			result.Skipped = append(result.Skipped, e)
			continue
		}
		if strings.HasPrefix(e.Value, encPrefix) {
			// already encrypted
			result.Skipped = append(result.Skipped, e)
			continue
		}
		ciphertext, err := aesgcmEncrypt(key, []byte(e.Value))
		if err != nil {
			return EncryptResult{}, fmt.Errorf("encrypt: key %s: %w", e.Key, err)
		}
		e.Value = encPrefix + base64.StdEncoding.EncodeToString(ciphertext)
		result.Encrypted = append(result.Encrypted, e)
	}
	return result, nil
}

// Decrypt decrypts entry values that were previously encrypted with Encrypt.
func Decrypt(entries []Entry, passphrase string) ([]Entry, error) {
	if passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}
	key := deriveKey(passphrase)
	out := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if !strings.HasPrefix(e.Value, encPrefix) {
			out = append(out, e)
			continue
		}
		encoded := strings.TrimPrefix(e.Value, encPrefix)
		ciphertext, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("decrypt: key %s: base64 decode: %w", e.Key, err)
		}
		plaintext, err := aesgcmDecrypt(key, ciphertext)
		if err != nil {
			return nil, fmt.Errorf("decrypt: key %s: %w", e.Key, err)
		}
		e.Value = string(plaintext)
		out = append(out, e)
	}
	return out, nil
}

func deriveKey(passphrase string) []byte {
	h := sha256.Sum256([]byte(passphrase))
	return h[:]
}

func aesgcmEncrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func aesgcmDecrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return nil, errors.New("ciphertext too short")
	}
	return gcm.Open(nil, data[:ns], data[ns:], nil)
}
