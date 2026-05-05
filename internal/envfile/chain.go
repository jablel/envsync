package envfile

// ChainOptions controls how a pipeline of operations is applied.
type ChainOptions struct {
	// StopOnError aborts the chain on the first error.
	StopOnError bool
}

// ChainStep is a single named operation in a chain.
type ChainStep struct {
	Name string
	Fn   func([]Entry) ([]Entry, error)
}

// ChainResult holds the outcome of executing a chain.
type ChainResult struct {
	Entries []Entry
	Steps   []ChainStepResult
}

// ChainStepResult records what happened during a single step.
type ChainStepResult struct {
	Name    string
	Applied bool
	Err     error
}

// HasErrors returns true if any step recorded an error.
func (r ChainResult) HasErrors() bool {
	for _, s := range r.Steps {
		if s.Err != nil {
			return true
		}
	}
	return false
}

// Chain executes a sequence of transformation steps against entries.
// If StopOnError is true the pipeline halts on the first error and returns
// the entries as they were before that step.
func Chain(entries []Entry, steps []ChainStep, opts ChainOptions) (ChainResult, error) {
	result := ChainResult{
		Entries: copyEntries(entries),
		Steps:   make([]ChainStepResult, 0, len(steps)),
	}

	for _, step := range steps {
		out, err := step.Fn(result.Entries)
		sr := ChainStepResult{Name: step.Name, Applied: err == nil, Err: err}
		result.Steps = append(result.Steps, sr)
		if err != nil {
			if opts.StopOnError {
				return result, err
			}
			continue
		}
		result.Entries = out
	}

	return result, nil
}

func copyEntries(src []Entry) []Entry {
	dst := make([]Entry, len(src))
	copy(dst, src)
	return dst
}
