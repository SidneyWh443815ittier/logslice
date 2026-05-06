package index

import (
	"bufio"
	"io"
	"time"

	"github.com/logslice/logslice/internal/timerange"
)

// Entry represents an indexed log line with its byte offset and timestamp.
type Entry struct {
	Offset    int64
	LineNum   int
	Timestamp time.Time
	HasTime   bool
}

// Index holds all indexed entries for a log file.
type Index struct {
	Entries []Entry
}

// Build scans a reader and builds an in-memory index of line offsets and timestamps.
func Build(r io.ReadSeeker) (*Index, error) {
	if _, err := r.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	idx := &Index{}
	scanner := bufio.NewScanner(r)
	var offset int64
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		entry := Entry{
			Offset:  offset,
			LineNum: lineNum,
		}

		if ts, ok := timerange.ExtractTimestamp(line); ok {
			entry.Timestamp = ts
			entry.HasTime = true
		}

		idx.Entries = append(idx.Entries, entry)
		offset += int64(len(line)) + 1 // +1 for newline
	}

	return idx, scanner.Err()
}

// FilterByRange returns entries whose timestamps fall within the given range.
func (idx *Index) FilterByRange(r timerange.Range) []Entry {
	var result []Entry
	for _, e := range idx.Entries {
		if !e.HasTime || timerange.LineInRange(e.Timestamp, r) {
			result = append(result, e)
		}
	}
	return result
}
