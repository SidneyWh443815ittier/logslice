// Package ratelimit provides output rate limiting for log streaming,
// allowing callers to cap the number of lines emitted per second.
package ratelimit

import (
	"time"
)

// Strategy controls how rate limiting is applied.
type Strategy int

const (
	// StrategyNone disables rate limiting; all lines pass through.
	StrategyNone Strategy = iota
	// StrategyLinesPerSecond limits output to a fixed number of lines per second.
	StrategyLinesPerSecond
)

// Limiter gates line emission according to a configured strategy.
type Limiter struct {
	strategy Strategy
	rate     int
	ticker   *time.Ticker
	count    int
	window   time.Time
}

// New constructs a Limiter. If rate <= 0 or strategy is StrategyNone,
// the limiter becomes a no-op pass-through.
func New(strategy Strategy, rate int) *Limiter {
	if strategy == StrategyNone || rate <= 0 {
		return &Limiter{strategy: StrategyNone}
	}
	return &Limiter{
		strategy: strategy,
		rate:     rate,
		window:   time.Now(),
	}
}

// Allow reports whether the next line should be emitted.
// It is safe to call Allow from a single goroutine.
func (l *Limiter) Allow() bool {
	if l.strategy == StrategyNone {
		return true
	}
	now := time.Now()
	if now.Sub(l.window) >= time.Second {
		l.window = now
		l.count = 0
	}
	if l.count < l.rate {
		l.count++
		return true
	}
	return false
}

// Reset clears internal counters, starting a fresh window immediately.
func (l *Limiter) Reset() {
	l.count = 0
	l.window = time.Now()
}

// Stats returns the number of lines allowed in the current window.
func (l *Limiter) Stats() (allowed int, rate int) {
	return l.count, l.rate
}
