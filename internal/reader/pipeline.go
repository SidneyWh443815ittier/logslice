package reader

import (
	"io"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/timerange"
)

// Options holds filtering parameters for a scan pipeline.
type Options struct {
	Range   *timerange.Range
	Queries []*filter.Query
}

// Result represents a matched log line.
type Result struct {
	LineNumber int
	Text       string
}

// Run scans the reader r, applying time range and field query filters,
// and sends matching lines to the returned channel. The channel is
// closed when scanning is complete or the reader is exhausted.
func Run(r io.Reader, opts Options) (<-chan Result, <-chan error) {
	results := make(chan Result, 64)
	errs := make(chan error, 1)

	go func() {
		defer close(results)
		defer close(errs)

		s := NewScanner(r)
		for s.Scan() {
			line := s.Line()

			if opts.Range != nil && !timerange.LineInRange(line, opts.Range) {
				continue
			}

			matched := true
			for _, q := range opts.Queries {
				if !q.Matches(line) {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}

			results <- Result{LineNumber: s.LineNumber(), Text: line}
		}

		if err := s.Err(); err != nil {
			errs <- err
		}
	}()

	return results, errs
}
