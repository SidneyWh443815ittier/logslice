// Package bookmark provides persistent read-position tracking for log files.
// It allows logslice to resume processing from the last known offset,
// enabling incremental log tailing across restarts.
package bookmark

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
)

// ErrNotFound is returned when no bookmark exists for a given file.
var ErrNotFound = errors.New("bookmark: no entry found")

// Entry records the last processed position in a log file.
type Entry struct {
	FilePath  string    `json:"file_path"`
	Offset    int64     `json:"offset"`
	LineNo    int64     `json:"line_no"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store persists bookmark entries to a JSON file on disk.
type Store struct {
	mu      sync.Mutex
	path    string
	entries map[string]Entry
}

// NewStore opens or creates a bookmark store at the given path.
func NewStore(path string) (*Store, error) {
	s := &Store{
		path:    path,
		entries: make(map[string]Entry),
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return s, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, &s.entries); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns the bookmark entry for the given file path.
func (s *Store) Get(filePath string) (Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[filePath]
	if !ok {
		return Entry{}, ErrNotFound
	}
	return e, nil
}

// Set records or updates the bookmark for the given file path.
func (s *Store) Set(filePath string, offset, lineNo int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[filePath] = Entry{
		FilePath:  filePath,
		Offset:    offset,
		LineNo:    lineNo,
		UpdatedAt: time.Now().UTC(),
	}
	return s.flush()
}

// Delete removes the bookmark for the given file path.
func (s *Store) Delete(filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, filePath)
	return s.flush()
}

// flush writes the current entries to disk. Must be called with mu held.
func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
