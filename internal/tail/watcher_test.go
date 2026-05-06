package tail_test

import (
	"os"
	"testing"
	"time"

	"logslice/internal/tail"
)

func TestWatcher_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := tail.NewWatcher(f.Name(), 20*time.Millisecond)
	if err := w.Start(); err != nil {
		t.Fatal(err)
	}
	defer w.Stop()

	// Write lines after the watcher has started.
	time.Sleep(30 * time.Millisecond)
	_, _ = f.WriteString("line one\nline two\n")

	collected := []string{}
	timeout := time.After(300 * time.Millisecond)
collect:
	for {
		select {
		case line := <-w.Lines():
			collected = append(collected, line)
			if len(collected) >= 2 {
				break collect
			}
		case err := <-w.Errors():
			t.Fatalf("unexpected error: %v", err)
		case <-timeout:
			break collect
		}
	}

	if len(collected) < 2 {
		t.Fatalf("expected at least 2 lines, got %d: %v", len(collected), collected)
	}
	if collected[0] != "line one" {
		t.Errorf("first line = %q, want %q", collected[0], "line one")
	}
	if collected[1] != "line two" {
		t.Errorf("second line = %q, want %q", collected[1], "line two")
	}
}

func TestWatcher_MissingFile(t *testing.T) {
	w := tail.NewWatcher("/nonexistent/path/log.txt", 10*time.Millisecond)
	if err := w.Start(); err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestWatcher_StopIsIdempotent(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "logslice-tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	w := tail.NewWatcher(f.Name(), 20*time.Millisecond)
	if err := w.Start(); err != nil {
		t.Fatal(err)
	}

	// Calling Stop multiple times should not panic or block.
	w.Stop()
	w.Stop()
}
