// Package dedupe provides line-level deduplication for log streams.
// It supports exact-match and windowed deduplication strategies.
package dedupe

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

// Strategy controls how duplicate detection is performed.
type Strategy int

const (
	// StrategyNone disables deduplication; all lines pass through.
	StrategyNone Strategy = iota
	// StrategyExact suppresses any line seen before in the current session.
	StrategyExact
	// StrategyWindow suppresses a line only if it appeared within the last N lines.
	StrategyWindow
)

// Deduplicator filters repeated log lines according to a chosen strategy.
type Deduplicator struct {
	strategy Strategy
	windowSize int

	mu     sync.Mutex
	seen   map[string]struct{} // used by StrategyExact
	window []string            // ring buffer of hashes for StrategyWindow
	pos    int
}

// New creates a Deduplicator. For StrategyWindow, windowSize sets how many
// recent line hashes are retained; it is ignored for other strategies.
func New(strategy Strategy, windowSize int) *Deduplicator {
	d := &Deduplicator{
		strategy:   strategy,
		windowSize: windowSize,
	}
	if strategy == StrategyExact {
		d.seen = make(map[string]struct{})
	}
	if strategy == StrategyWindow && windowSize > 0 {
		d.window = make([]string, windowSize)
	}
	return d
}

// IsDuplicate reports whether line has been seen according to the active
// strategy. It updates internal state as a side-effect.
func (d *Deduplicator) IsDuplicate(line string) bool {
	if d.strategy == StrategyNone {
		return false
	}
	h := hash(line)
	d.mu.Lock()
	defer d.mu.Unlock()

	switch d.strategy {
	case StrategyExact:
		if _, ok := d.seen[h]; ok {
			return true
		}
		d.seen[h] = struct{}{}
		return false

	case StrategyWindow:
		if d.windowSize <= 0 {
			return false
		}
		for _, v := range d.window {
			if v == h {
				return true
			}
		}
		d.window[d.pos%d.windowSize] = h
		d.pos++
		return false
	}
	return false
}

// Reset clears all retained state, allowing previously seen lines to pass again.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
	if d.windowSize > 0 {
		d.window = make([]string, d.windowSize)
	}
	d.pos = 0
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
