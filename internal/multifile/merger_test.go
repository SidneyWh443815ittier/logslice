package multifile

import (
	"strings"
	"testing"
)

const (
	log1 = `2024-01-15T10:00:01Z level=info msg="service started"
2024-01-15T10:00:05Z level=info msg="request received"
2024-01-15T10:00:09Z level=warn msg="slow query"`

	log2 = `2024-01-15T10:00:02Z level=info msg="worker ready"
2024-01-15T10:00:06Z level=error msg="connection refused"
2024-01-15T10:00:10Z level=info msg="retry succeeded"`
)

func collectMerge(t *testing.T, sources []Source) []MergedLine {
	t.Helper()
	out := make(chan MergedLine, 64)
	var mergeErr error
	go func() {
		defer close(out)
		mergeErr = Merge(sources, out)
	}()
	var lines []MergedLine
	for ml := range out {
		lines = append(lines, ml)
	}
	if mergeErr != nil {
		t.Fatalf("Merge returned error: %v", mergeErr)
	}
	return lines
}

func TestMerge_ChronologicalOrder(t *testing.T) {
	sources := []Source{
		{Name: "api", Reader: strings.NewReader(log1)},
		{Name: "worker", Reader: strings.NewReader(log2)},
	}
	lines := collectMerge(t, sources)

	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d", len(lines))
	}

	want := []string{
		"service started",
		"worker ready",
		"request received",
		"connection refused",
		"slow query",
		"retry succeeded",
	}
	for i, ml := range lines {
		if !strings.Contains(ml.Line, want[i]) {
			t.Errorf("line %d: expected to contain %q, got %q", i, want[i], ml.Line)
		}
	}
}

func TestMerge_SourceAttribution(t *testing.T) {
	sources := []Source{
		{Name: "api", Reader: strings.NewReader(log1)},
		{Name: "worker", Reader: strings.NewReader(log2)},
	}
	lines := collectMerge(t, sources)

	if lines[0].Source != "api" {
		t.Errorf("first line source: want api, got %s", lines[0].Source)
	}
	if lines[1].Source != "worker" {
		t.Errorf("second line source: want worker, got %s", lines[1].Source)
	}
}

func TestMerge_SingleSource(t *testing.T) {
	sources := []Source{
		{Name: "only", Reader: strings.NewReader(log1)},
	}
	lines := collectMerge(t, sources)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
}

func TestMerge_EmptySources(t *testing.T) {
	lines := collectMerge(t, []Source{})
	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestMerge_UnparseableTimestamps(t *testing.T) {
	raw := "no timestamp here\nstill no timestamp"
	sources := []Source{
		{Name: "raw", Reader: strings.NewReader(raw)},
	}
	lines := collectMerge(t, sources)
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}
