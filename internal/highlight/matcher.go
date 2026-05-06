// Package highlight provides utilities for marking matched substrings
// within log lines for terminal or structured output.
package highlight

import (
	"strings"
)

// ANSI escape codes for terminal highlighting.
const (
	AnsiReset  = "\033[0m"
	AnsiBold   = "\033[1m"
	AnsiYellow = "\033[33m"
	AnsiCyan   = "\033[36m"
)

// Match represents a single highlighted region within a line.
type Match struct {
	Start int
	End   int
	Term  string
}

// FindMatches returns all non-overlapping case-insensitive matches of terms
// within line, ordered by their start position.
func FindMatches(line string, terms []string) []Match {
	var matches []Match
	lower := strings.ToLower(line)
	for _, term := range terms {
		if term == "" {
			continue
		}
		lowerTerm := strings.ToLower(term)
		offset := 0
		for {
			idx := strings.Index(lower[offset:], lowerTerm)
			if idx == -1 {
				break
			}
			start := offset + idx
			matches = append(matches, Match{
				Start: start,
				End:   start + len(term),
				Term:  term,
			})
			offset = start + len(term)
		}
	}
	return matches
}

// ApplyANSI wraps matched regions in the line with ANSI bold+yellow escape
// sequences, returning the annotated string. If there are no matches the
// original line is returned unchanged.
func ApplyANSI(line string, matches []Match) string {
	if len(matches) == 0 {
		return line
	}
	var sb strings.Builder
	cursor := 0
	for _, m := range matches {
		if m.Start < cursor {
			// skip overlapping matches
			continue
		}
		sb.WriteString(line[cursor:m.Start])
		sb.WriteString(AnsiBold)
		sb.WriteString(AnsiYellow)
		sb.WriteString(line[m.Start:m.End])
		sb.WriteString(AnsiReset)
		cursor = m.End
	}
	sb.WriteString(line[cursor:])
	return sb.String()
}
