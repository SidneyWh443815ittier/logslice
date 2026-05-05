package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatter_Raw(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatRaw, nil)

	if err := f.Write(1, "hello world"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimRight(buf.String(), "\n")
	if got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON, nil)

	if err := f.Write(42, `level=info msg="started"`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var env envelope
	if err := json.Unmarshal([]byte(strings.TrimRight(buf.String(), "\n")), &env); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v", err)
	}
	if env.Line != 42 {
		t.Errorf("expected line 42, got %d", env.Line)
	}
	if !strings.Contains(env.Raw, "started") {
		t.Errorf("raw field missing expected content, got %q", env.Raw)
	}
}

func TestFormatter_Color_HighlightsField(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatColor, []string{"error"})

	if err := f.Write(1, "level=error msg=crash"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, ansiYellow+"error"+ansiReset) {
		t.Errorf("expected ANSI highlight around 'error', got %q", out)
	}
}

func TestFormatter_Color_NoFields(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatColor, nil)

	if err := f.Write(1, "plain line"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := strings.TrimRight(buf.String(), "\n")
	if got != "plain line" {
		t.Errorf("expected unmodified line, got %q", got)
	}
}

func TestFormatter_JSON_LineNumbers(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, FormatJSON, nil)

	for i := uint64(1); i <= 3; i++ {
		buf.Reset()
		_ = f.Write(i, "msg")
		var env envelope
		_ = json.Unmarshal([]byte(strings.TrimRight(buf.String(), "\n")), &env)
		if env.Line != i {
			t.Errorf("line %d: expected %d, got %d", i, i, env.Line)
		}
	}
}
