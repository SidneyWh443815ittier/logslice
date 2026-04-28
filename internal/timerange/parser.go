package timerange

import (
	"fmt"
	"time"
)

// Common log timestamp layouts to attempt parsing
var knownLayouts = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05.000",
	"2006-01-02 15:04:05.000",
	"02/Jan/2006:15:04:05 -0700",
}

// Range holds the start and end of a time window.
type Range struct {
	Start time.Time
	End   time.Time
}

// ParseRange parses two timestamp strings into a Range.
// Either value may be empty to indicate an open-ended range.
func ParseRange(from, to string) (Range, error) {
	var r Range
	var err error

	if from != "" {
		r.Start, err = ParseTimestamp(from)
		if err != nil {
			return r, fmt.Errorf("invalid --from timestamp: %w", err)
		}
	}

	if to != "" {
		r.End, err = ParseTimestamp(to)
		if err != nil {
			return r, fmt.Errorf("invalid --to timestamp: %w", err)
		}
	}

	if !r.Start.IsZero() && !r.End.IsZero() && r.End.Before(r.Start) {
		return r, fmt.Errorf("--to must not be before --from")
	}

	return r, nil
}

// ParseTimestamp tries each known layout until one succeeds.
func ParseTimestamp(s string) (time.Time, error) {
	for _, layout := range knownLayouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognised timestamp format: %q", s)
}

// Contains reports whether t falls within the range.
// An unset Start/End is treated as unbounded.
func (r Range) Contains(t time.Time) bool {
	if !r.Start.IsZero() && t.Before(r.Start) {
		return false
	}
	if !r.End.IsZero() && t.After(r.End) {
		return false
	}
	return true
}
