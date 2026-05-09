package rotate_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/rotate"
)

func writeTmp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-rotate-*.log")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestNew_MissingFile(t *testing.T) {
	_, err := rotate.New("/nonexistent/path/logslice.log", time.Second)
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestNew_DefaultPollInterval(t *testing.T) {
	path := writeTmp(t, "hello\n")
	d, err := rotate.New(path, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.PollInterval() != time.Second {
		t.Errorf("expected 1s default, got %v", d.PollInterval())
	}
}

func TestCheck_NoChange(t *testing.T) {
	path := writeTmp(t, "line1\n")
	d, err := rotate.New(path, time.Second)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	ev, err := d.Check()
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if ev != rotate.EventNone {
		t.Errorf("expected EventNone, got %v", ev)
	}
}

func TestCheck_Truncated(t *testing.T) {
	path := writeTmp(t, "some longer content here\n")
	d, err := rotate.New(path, time.Second)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	// Truncate the file.
	if err := os.WriteFile(path, []byte("x\n"), 0o644); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	ev, err := d.Check()
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if ev != rotate.EventTruncated {
		t.Errorf("expected EventTruncated, got %v", ev)
	}
}

func TestCheck_Replaced(t *testing.T) {
	path := writeTmp(t, "original\n")
	d, err := rotate.New(path, time.Second)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	// Remove and recreate to change inode.
	if err := os.Remove(path); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if err := os.WriteFile(path, []byte("new content\n"), 0o644); err != nil {
		t.Fatalf("recreate: %v", err)
	}
	ev, err := d.Check()
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if ev != rotate.EventReplaced {
		t.Errorf("expected EventReplaced, got %v", ev)
	}
}

func TestCheck_DeletedFile(t *testing.T) {
	path := writeTmp(t, "data\n")
	d, err := rotate.New(path, time.Second)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	os.Remove(path)
	ev, err := d.Check()
	if err != nil {
		t.Fatalf("unexpected error on deleted file: %v", err)
	}
	if ev != rotate.EventReplaced {
		t.Errorf("expected EventReplaced for deleted file, got %v", ev)
	}
}
