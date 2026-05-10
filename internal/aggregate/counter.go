// Package aggregate provides log line aggregation and counting utilities.
// It supports grouping matched lines by a field value and counting occurrences.
package aggregate

import (
	"fmt"
	"sort"
	"strings"
)

// Entry holds the count for a single group key.
type Entry struct {
	Key   string
	Count int
}

// Counter accumulates counts keyed by a field value extracted from log lines.
type Counter struct {
	field  string
	counts map[string]int
}

// New creates a Counter that groups lines by the given field name.
// Field values are extracted using simple key=value or key="value" scanning.
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add records one occurrence of the value found for the configured field in line.
// If the field is absent the line is counted under the empty-string key.
func (c *Counter) Add(line string) {
	key := extractField(line, c.field)
	c.counts[key]++
}

// Entries returns all accumulated entries sorted by count descending,
// with ties broken alphabetically by key.
func (c *Counter) Entries() []Entry {
	entries := make([]Entry, 0, len(c.counts))
	for k, v := range c.counts {
		entries = append(entries, Entry{Key: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Reset clears all accumulated counts.
func (c *Counter) Reset() {
	c.counts = make(map[string]int)
}

// Total returns the sum of all recorded occurrences.
func (c *Counter) Total() int {
	n := 0
	for _, v := range c.counts {
		n += v
	}
	return n
}

// extractField scans line for field=value or field="value" and returns the value.
func extractField(line, field string) string {
	prefix := field + "="
	idx := strings.Index(line, prefix)
	if idx == -1 {
		return ""
	}
	rest := line[idx+len(prefix):]
	if len(rest) == 0 {
		return ""
	}
	if rest[0] == '"' {
		end := strings.Index(rest[1:], "\"")
		if end == -1 {
			return rest[1:]
		}
		return rest[1 : end+1]
	}
	end := strings.IndexAny(rest, " \t\n\r")
	if end == -1 {
		return rest
	}
	return rest[:end]
}

// FormatTable renders entries as a simple two-column text table.
func FormatTable(entries []Entry) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-40s %s\n", "KEY", "COUNT"))
	sb.WriteString(strings.Repeat("-", 48) + "\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%-40s %d\n", e.Key, e.Count))
	}
	return sb.String()
}
