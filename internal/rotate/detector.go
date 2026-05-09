// Package rotate detects log file rotation events such as truncation or
// inode replacement, allowing consumers to reopen files as needed.
package rotate

import (
	"errors"
	"os"
	"time"
)

// Event describes the kind of rotation that was detected.
type Event int

const (
	// EventNone indicates no rotation has occurred.
	EventNone Event = iota
	// EventTruncated indicates the file was truncated (size decreased).
	EventTruncated
	// EventReplaced indicates the file's inode changed (renamed/replaced).
	EventReplaced
)

//go:generate stringer -type=Event

// Detector watches a single file path for rotation events.
type Detector struct {
	path     string
	pollInterval time.Duration
	lastSize int64
	lastIno  uint64
}

// New creates a Detector for the given path, sampling at pollInterval.
// A zero or negative pollInterval defaults to one second.
func New(path string, pollInterval time.Duration) (*Detector, error) {
	if pollInterval <= 0 {
		pollInterval = time.Second
	}
	d := &Detector{path: path, pollInterval: pollInterval}
	if err := d.snapshot(); err != nil {
		return nil, err
	}
	return d, nil
}

// Check returns the rotation Event observed since the last call to Check
// (or since construction). EventNone is returned when the file is unchanged.
func (d *Detector) Check() (Event, error) {
	info, err := os.Stat(d.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return EventReplaced, nil
		}
		return EventNone, err
	}
	ino := inode(info)
	size := info.Size()

	var ev Event
	switch {
	case ino != d.lastIno:
		ev = EventReplaced
	case size < d.lastSize:
		ev = EventTruncated
	default:
		ev = EventNone
	}

	d.lastIno = ino
	d.lastSize = size
	return ev, nil
}

// PollInterval returns the configured polling interval.
func (d *Detector) PollInterval() time.Duration { return d.pollInterval }

func (d *Detector) snapshot() error {
	info, err := os.Stat(d.path)
	if err != nil {
		return err
	}
	d.lastIno = inode(info)
	d.lastSize = info.Size()
	return nil
}
