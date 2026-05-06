package index

import (
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/timerange"
)

const sampleLog = `2024-01-15T10:00:00Z level=info msg="starting service"
2024-01-15T10:01:00Z level=debug msg="connected to db"
2024-01-15T10:02:00Z level=error msg="request failed"
no-timestamp-line
2024-01-15T10:03:00Z level=info msg="recovered"
`

func TestBuild_EntryCount(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(idx.Entries))
	}
}

func TestBuild_TimestampParsed(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !idx.Entries[0].HasTime {
		t.Error("expected first entry to have a timestamp")
	}
	if idx.Entries[3].HasTime {
		t.Error("expected no-timestamp line to have HasTime=false")
	}
}

func TestBuild_LineNumbers(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, e := range idx.Entries {
		if e.LineNum != i+1 {
			t.Errorf("entry %d: expected LineNum %d, got %d", i, i+1, e.LineNum)
		}
	}
}

func TestFilterByRange(t *testing.T) {
	r := strings.NewReader(sampleLog)
	idx, err := Build(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	start := time.Date(2024, 1, 15, 10, 1, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 10, 2, 0, 0, time.UTC)
	rng := timerange.Range{Start: start, End: end}

	results := idx.FilterByRange(rng)
	// entries with timestamps in [10:01, 10:02] plus the no-timestamp line
	if len(results) == 0 {
		t.Error("expected at least one result")
	}
}
