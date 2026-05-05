package reader

import (
	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/stats"
	"github.com/user/logslice/internal/timerange"
)

// RunOptions configures a pipeline execution.
type RunOptions struct {
	Scanner   *Scanner
	Range     *timerange.Range
	Query     *filter.Query
	Formatter *output.Formatter
	Stats     *stats.Collector
}

// Run reads lines from the scanner, applies time-range and field filters,
// and writes matching lines via the formatter. It returns the number of
// matched lines and any write error encountered.
func Run(opts RunOptions) (int64, error) {
	var matched int64

	for opts.Scanner.Scan() {
		line := opts.Scanner.Text()
		lineNum := opts.Scanner.LineNumber()

		// Time-range filter
		if opts.Range != nil && !timerange.LineInRange(line, opts.Range) {
			if opts.Stats != nil {
				opts.Stats.RecordLine(len(line), false)
			}
			continue
		}

		// Field query filter
		if opts.Query != nil && !opts.Query.Matches(line) {
			if opts.Stats != nil {
				opts.Stats.RecordLine(len(line), false)
			}
			continue
		}

		if opts.Stats != nil {
			opts.Stats.RecordLine(len(line), true)
		}

		if err := opts.Formatter.Write(line, lineNum); err != nil {
			return matched, err
		}
		matched++
	}

	return matched, opts.Scanner.Err()
}
