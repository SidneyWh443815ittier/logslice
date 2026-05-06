// Package tail provides live log file watching and tailing functionality.
package tail

import (
	"io"
	"os"
	"time"
)

// Watcher monitors a file for new content and emits new lines.
type Watcher struct {
	path     string
	pollInterval time.Duration
	lines    chan string
	errors   chan error
	done     chan struct{}
}

// NewWatcher creates a Watcher for the given file path.
func NewWatcher(path string, pollInterval time.Duration) *Watcher {
	return &Watcher{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan string, 64),
		errors:       make(chan error, 8),
		done:         make(chan struct{}),
	}
}

// Lines returns the channel on which new log lines are emitted.
func (w *Watcher) Lines() <-chan string { return w.lines }

// Errors returns the channel on which watch errors are emitted.
func (w *Watcher) Errors() <-chan error { return w.errors }

// Start begins tailing the file from its current end.
func (w *Watcher) Start() error {
	f, err := os.Open(w.path)
	if err != nil {
		return err
	}
	// Seek to end so we only emit new content.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		f.Close()
		return err
	}
	go w.poll(f)
	return nil
}

// Stop signals the watcher to stop.
func (w *Watcher) Stop() { close(w.done) }

func (w *Watcher) poll(f *os.File) {
	defer f.Close()
	buf := make([]byte, 0, 4096)
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-w.done:
			return
		case <-ticker.C:
			buf = buf[:0]
			tmp := make([]byte, 4096)
			for {
				n, err := f.Read(tmp)
				if n > 0 {
					buf = append(buf, tmp[:n]...)
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					w.errors <- err
					return
				}
			}
			if len(buf) > 0 {
				emitLines(buf, w.lines)
			}
		}
	}
}
