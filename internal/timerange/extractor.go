package timerange

import (
	"regexp"
	"strings"
	"time"
)

// timestampRe matches ISO-8601 / common log date patterns embedded in a line.
var timestampRe = regexp.MustCompile(
	`(\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:[.,]\d+)?(?:Z|[+-]\d{2}:?\d{2})?|` +
		`\d{2}/[A-Za-z]{3}/\d{4}:\d{2}:\d{2}:\d{2} [+-]\d{4})`,
)

// ExtractTimestamp scans a log line and returns the first parseable timestamp
// it finds, along with a boolean indicating success.
func ExtractTimestamp(line string) (time.Time, bool) {
	matches := timestampRe.FindAllString(line, -1)
	for _, m := range matches {
		// Normalise separator between date and time
		normalised := strings.Replace(m, " ", "T", 1)
		t, err := ParseTimestamp(normalised)
		if err == nil {
			return t, true
		}
		// Try original form too
		t, err = ParseTimestamp(m)
		if err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// LineInRange returns true if the line contains a timestamp that falls within r.
// Lines with no parseable timestamp are excluded when the range is bounded.
func LineInRange(line string, r Range) bool {
	t, ok := ExtractTimestamp(line)
	if !ok {
		// If the range is completely open, pass through lines without timestamps.
		return r.Start.IsZero() && r.End.IsZero()
	}
	return r.Contains(t)
}
