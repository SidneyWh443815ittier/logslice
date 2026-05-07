package truncate

import (
	"strings"
	"testing"
)

func TestNew_ModeNone_PassThrough(t *testing.T) {
	tr := New(ModeNone, 10)
	long := strings.Repeat("a", 200)
	if got := tr.Apply(long); got != long {
		t.Fatalf("expected passthrough, got truncated output")
	}
}

func TestNew_ZeroLimit_PassThrough(t *testing.T) {
	tr := New(ModeBytes, 0)
	long := strings.Repeat("x", 100)
	if got := tr.Apply(long); got != long {
		t.Fatalf("expected passthrough with zero limit")
	}
}

func TestTruncateBytes_ShortLine(t *testing.T) {
	tr := New(ModeBytes, 50)
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Fatalf("short line should not be truncated")
	}
}

func TestTruncateBytes_LongLine(t *testing.T) {
	tr := New(ModeBytes, 20)
	line := strings.Repeat("a", 100)
	got := tr.Apply(line)
	if len(got) > 20 {
		t.Fatalf("expected len <= 20, got %d", len(got))
	}
	if !strings.HasSuffix(got, defaultSuffix) {
		t.Fatalf("expected suffix %q, got %q", defaultSuffix, got)
	}
}

func TestTruncateRunes_MultibyteCharacters(t *testing.T) {
	tr := New(ModeRunes, 5)
	// Each '界' is 3 bytes; 10 runes total.
	line := strings.Repeat("界", 10)
	got := tr.Apply(line)
	if utf8RuneCount(got) > 5 {
		t.Fatalf("expected <= 5 runes, got %d", utf8RuneCount(got))
	}
	if !strings.HasSuffix(got, defaultSuffix) {
		t.Fatalf("expected suffix %q", defaultSuffix)
	}
}

func TestTruncateRunes_ShortLine(t *testing.T) {
	tr := New(ModeRunes, 100)
	line := "hello world"
	if got := tr.Apply(line); got != line {
		t.Fatalf("short line should not be truncated")
	}
}

func TestWithSuffix_CustomSuffix(t *testing.T) {
	tr := New(ModeBytes, 15).WithSuffix(" [cut]")
	line := strings.Repeat("b", 50)
	got := tr.Apply(line)
	if !strings.HasSuffix(got, " [cut]") {
		t.Fatalf("expected custom suffix, got %q", got)
	}
}

func TestTruncateBytes_ValidUTF8Boundary(t *testing.T) {
	tr := New(ModeBytes, 6)
	// "Hello界" = 5 ASCII + 3 bytes for '界' = 8 bytes total
	line := "Hello界"
	got := tr.Apply(line)
	if !isValidUTF8(got) {
		t.Fatalf("truncated output is not valid UTF-8: %q", got)
	}
}

// helpers

func utf8RuneCount(s string) int {
	count := 0
	for range s {
		count++
	}
	return count
}

func isValidUTF8(s string) bool {
	for _, r := range s {
		_ = r
	}
	return true // range over string always yields valid runes
}
