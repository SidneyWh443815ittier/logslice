// Package truncate provides line truncation utilities for limiting output
// line length when processing large log files with very long lines.
package truncate

import "unicode/utf8"

// Mode controls how truncation is applied.
type Mode int

const (
	// ModeNone disables truncation.
	ModeNone Mode = iota
	// ModeBytes truncates at a byte boundary.
	ModeBytes
	// ModeRunes truncates at a rune boundary, preserving valid UTF-8.
	ModeRunes
)

const defaultSuffix = "..."

// Truncator applies length-based truncation to log lines.
type Truncator struct {
	mode   Mode
	limit  int
	suffix string
}

// New creates a Truncator with the given mode and character/byte limit.
// A limit of 0 disables truncation regardless of mode.
func New(mode Mode, limit int) *Truncator {
	return &Truncator{
		mode:   mode,
		limit:  limit,
		suffix: defaultSuffix,
	}
}

// WithSuffix returns a copy of the Truncator that appends the given suffix
// to truncated lines instead of the default "...".
func (t *Truncator) WithSuffix(suffix string) *Truncator {
	return &Truncator{mode: t.mode, limit: t.limit, suffix: suffix}
}

// Apply truncates line according to the configured mode and limit.
// If the line does not exceed the limit, it is returned unchanged.
func (t *Truncator) Apply(line string) string {
	if t.mode == ModeNone || t.limit <= 0 {
		return line
	}
	switch t.mode {
	case ModeBytes:
		return t.truncateBytes(line)
	case ModeRunes:
		return t.truncateRunes(line)
	default:
		return line
	}
}

func (t *Truncator) truncateBytes(line string) string {
	if len(line) <= t.limit {
		return line
	}
	cutoff := t.limit - len(t.suffix)
	if cutoff < 0 {
		cutoff = 0
	}
	// Walk back to a valid UTF-8 boundary.
	for cutoff > 0 && !utf8.RuneStart(line[cutoff]) {
		cutoff--
	}
	return line[:cutoff] + t.suffix
}

func (t *Truncator) truncateRunes(line string) string {
	count := utf8.RuneCountInString(line)
	if count <= t.limit {
		return line
	}
	suffixRunes := utf8.RuneCountInString(t.suffix)
	target := t.limit - suffixRunes
	if target < 0 {
		target = 0
	}
	i, n := 0, 0
	for n < target {
		_, size := utf8.DecodeRuneInString(line[i:])
		i += size
		n++
	}
	return line[:i] + t.suffix
}
