// Package context provides utilities for bounding log processing
// to a surrounding window of lines around each match.
package context

// Limiter captures N lines before and after each matched line,
// emitting a contiguous, de-duplicated window to the caller.
type Limiter struct {
	before int
	after  int
	buf    []string // circular buffer for pre-match lines
	head   int
	count  int
	pending int     // lines remaining to emit after a match
}

// New creates a Limiter that retains `before` lines of look-back
// and emits `after` lines following each match.
func New(before, after int) *Limiter {
	if before < 0 {
		before = 0
	}
	if after < 0 {
		after = 0
	}
	cap := before
	if cap == 0 {
		cap = 1 // avoid zero-length slice
	}
	return &Limiter{
		before:  before,
		after:   after,
		buf:     make([]string, cap),
		pending: 0,
	}
}

// Feed records a line and whether it matched the current filter.
// It returns the lines that should be emitted as a result of this call
// (may be empty, the line itself, or the line plus buffered context).
func (l *Limiter) Feed(line string, matched bool) []string {
	if matched {
		var out []string
		// Flush the pre-match circular buffer in insertion order.
		if l.before > 0 && l.count > 0 {
			n := l.count
			if n > l.before {
				n = l.before
			}
			start := (l.head - n + len(l.buf)) % len(l.buf)
			for i := 0; i < n; i++ {
				out = append(out, l.buf[(start+i)%len(l.buf)])
			}
		}
		out = append(out, line)
		l.pending = l.after
		l.resetBuffer()
		return out
	}

	if l.pending > 0 {
		l.pending--
		return []string{line}
	}

	// Buffer for potential future match context.
	if l.before > 0 {
		l.buf[l.head] = line
		l.head = (l.head + 1) % len(l.buf)
		if l.count < l.before {
			l.count++
		}
	}
	return nil
}

// Reset clears all internal state, allowing the Limiter to be reused.
func (l *Limiter) Reset() {
	l.resetBuffer()
	l.pending = 0
}

func (l *Limiter) resetBuffer() {
	l.head = 0
	l.count = 0
}
