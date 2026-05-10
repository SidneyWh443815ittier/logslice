package aggregate

import (
	"strings"
	"testing"
)

func TestNew_EmptyCounts(t *testing.T) {
	c := New("level")
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 total, got %d", got)
	}
	if got := len(c.Entries()); got != 0 {
		t.Fatalf("expected 0 entries, got %d", got)
	}
}

func TestAdd_CountsField(t *testing.T) {
	c := New("level")
	c.Add(`ts=2024-01-01 level=info msg="started"`)
	c.Add(`ts=2024-01-01 level=info msg="running"`)
	c.Add(`ts=2024-01-01 level=error msg="oops"`)

	if got := c.Total(); got != 3 {
		t.Fatalf("expected total 3, got %d", got)
	}
	entries := c.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// highest count first
	if entries[0].Key != "info" || entries[0].Count != 2 {
		t.Errorf("unexpected top entry: %+v", entries[0])
	}
	if entries[1].Key != "error" || entries[1].Count != 1 {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestAdd_MissingField_EmptyKey(t *testing.T) {
	c := New("level")
	c.Add(`ts=2024-01-01 msg="no level here"`)
	entries := c.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Key != "" {
		t.Errorf("expected empty key, got %q", entries[0].Key)
	}
}

func TestAdd_QuotedValue(t *testing.T) {
	c := New("service")
	c.Add(`level=info service="my app" msg=ok`)
	entries := c.Entries()
	if len(entries) != 1 || entries[0].Key != "my app" {
		t.Errorf("expected key 'my app', got %+v", entries)
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	c := New("level")
	c.Add(`level=info msg=a`)
	c.Add(`level=error msg=b`)
	c.Reset()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestEntries_SortedByCountThenKey(t *testing.T) {
	c := New("level")
	for i := 0; i < 3; i++ {
		c.Add(`level=warn`)
	}
	for i := 0; i < 3; i++ {
		c.Add(`level=info`)
	}
	c.Add(`level=error`)
	entries := c.Entries()
	// warn and info both have count 3; alphabetically info < warn
	if entries[0].Key != "info" {
		t.Errorf("expected info first (alpha tie-break), got %q", entries[0].Key)
	}
	if entries[1].Key != "warn" {
		t.Errorf("expected warn second, got %q", entries[1].Key)
	}
	if entries[2].Key != "error" {
		t.Errorf("expected error last, got %q", entries[2].Key)
	}
}

func TestFormatTable_ContainsHeaders(t *testing.T) {
	c := New("level")
	c.Add(`level=info`)
	out := FormatTable(c.Entries())
	if !strings.Contains(out, "KEY") || !strings.Contains(out, "COUNT") {
		t.Errorf("table missing headers: %s", out)
	}
	if !strings.Contains(out, "info") {
		t.Errorf("table missing entry 'info': %s", out)
	}
}
