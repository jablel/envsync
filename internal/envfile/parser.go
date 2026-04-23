package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from an .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
}

// Parse reads and parses the .env file at the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	env := &EnvFile{Path: path}
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			env.Entries = append(env.Entries, Entry{
				Comment: trimmed,
				Line:    lineNum,
			})
			continue
		}

		key, value, err := parseLine(trimmed)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		env.Entries = append(env.Entries, Entry{
			Key:   key,
			Value: value,
			Line:  lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return env, nil
}

// ToMap converts the env file entries into a key-value map.
func (e *EnvFile) ToMap() map[string]string {
	m := make(map[string]string, len(e.Entries))
	for _, entry := range e.Entries {
		if entry.Key != "" {
			m[entry.Key] = entry.Value
		}
	}
	return m
}

func parseLine(line string) (string, string, error) {
	idx := strings.IndexByte(line, '=')
	if idx < 0 {
		return "", "", fmt.Errorf("invalid line %q: missing '='" , line)
	}
	key := strings.TrimSpace(line[:idx])
	value := strings.TrimSpace(line[idx+1:])
	value = stripQuotes(value)
	return key, value, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
