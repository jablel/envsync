package envfile

import (
	"fmt"
	"strings"
)

// Annotation holds metadata attached to an env entry.
type Annotation struct {
	Key     string
	Comment string
	Tags    []string
}

// AnnotateOptions controls how annotations are applied.
type AnnotateOptions struct {
	// Overwrite replaces existing comments/tags when true.
	Overwrite bool
}

// AnnotateResult describes what happened to each entry.
type AnnotateResult struct {
	Key     string
	Applied bool
	Skipped bool // skipped because entry already had annotation and Overwrite=false
}

// Annotate applies a set of Annotation records to entries.
// Annotations are matched by key. The comment is prepended as an inline
// comment and tags are appended as a comma-separated tag list comment.
func Annotate(entries []Entry, annotations []Annotation, opts AnnotateOptions) ([]Entry, []AnnotateResult, error) {
	if len(entries) == 0 {
		return entries, nil, nil
	}

	annoMap := make(map[string]Annotation, len(annotations))
	for _, a := range annotations {
		if strings.TrimSpace(a.Key) == "" {
			return nil, nil, fmt.Errorf("annotation key must not be empty")
		}
		annoMap[a.Key] = a
	}

	results := make([]AnnotateResult, 0, len(entries))
	out := make([]Entry, len(entries))

	for i, e := range entries {
		anno, ok := annoMap[e.Key]
		if !ok {
			out[i] = e
			continue
		}

		hasExisting := e.Comment != ""
		if hasExisting && !opts.Overwrite {
			out[i] = e
			results = append(results, AnnotateResult{Key: e.Key, Skipped: true})
			continue
		}

		updated := e
		parts := []string{}
		if anno.Comment != "" {
			parts = append(parts, anno.Comment)
		}
		if len(anno.Tags) > 0 {
			parts = append(parts, "tags:"+strings.Join(anno.Tags, ","))
		}
		updated.Comment = strings.Join(parts, " | ")
		out[i] = updated
		results = append(results, AnnotateResult{Key: e.Key, Applied: true})
	}

	return out, results, nil
}

// AnnotatedKeys returns the keys that have a non-empty comment.
func AnnotatedKeys(entries []Entry) []string {
	keys := []string{}
	for _, e := range entries {
		if e.Comment != "" {
			keys = append(keys, e.Key)
		}
	}
	return keys
}
