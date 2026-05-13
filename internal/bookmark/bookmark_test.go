package bookmark_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"logslice/internal/bookmark"
)

func tempStore(t *testing.T) (*bookmark.Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")
	s, err := bookmark.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return s, path
}

func TestNewStore_CreatesEmptyStore(t *testing.T) {
	s, _ := tempStore(t)
	_, err := s.Get("/var/log/app.log")
	if !errors.Is(err, bookmark.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestSet_PersistsEntry(t *testing.T) {
	s, path := tempStore(t)
	if err := s.Set("/var/log/app.log", 1024, 42); err != nil {
		t.Fatalf("Set: %v", err)
	}

	// Reload from disk.
	s2, err := bookmark.NewStore(path)
	if err != nil {
		t.Fatalf("NewStore reload: %v", err)
	}
	e, err := s2.Get("/var/log/app.log")
	if err != nil {
		t.Fatalf("Get after reload: %v", err)
	}
	if e.Offset != 1024 {
		t.Errorf("offset: want 1024, got %d", e.Offset)
	}
	if e.LineNo != 42 {
		t.Errorf("lineNo: want 42, got %d", e.LineNo)
	}
}

func TestSet_OverwritesExisting(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Set("/var/log/app.log", 100, 5)
	_ = s.Set("/var/log/app.log", 200, 10)
	e, err := s.Get("/var/log/app.log")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if e.Offset != 200 {
		t.Errorf("expected offset 200, got %d", e.Offset)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Set("/var/log/app.log", 512, 20)
	if err := s.Delete("/var/log/app.log"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("/var/log/app.log")
	if !errors.Is(err, bookmark.ErrNotFound) {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestNewStore_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bookmarks.json")
	if err := os.WriteFile(path, []byte("not-json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := bookmark.NewStore(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestEntry_UpdatedAt_IsSet(t *testing.T) {
	s, _ := tempStore(t)
	_ = s.Set("/var/log/app.log", 0, 0)
	e, _ := s.Get("/var/log/app.log")
	if e.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}
