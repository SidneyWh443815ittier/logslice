// Package sampling provides log line sampling strategies for large files
// where processing every line is unnecessary or too slow.
package sampling

import (
	"math/rand"
	"sync/atomic"
)

// Strategy defines how lines are selected during sampling.
type Strategy int

const (
	// StrategyNone disables sampling (all lines pass through).
	StrategyNone Strategy = iota
	// StrategyRandom retains each line with probability 1/Rate.
	StrategyRandom
	// StrategyNth retains every Nth line deterministically.
	StrategyNth
)

// Config holds sampling configuration.
type Config struct {
	Strategy Strategy
	// Rate is the inverse sampling rate: 1 means keep all, 10 means keep ~1 in 10.
	Rate int
}

// Sampler decides whether each log line should be kept.
type Sampler struct {
	cfg     Config
	counter atomic.Int64
	rng     *rand.Rand
}

// New creates a Sampler from the given Config.
// Rate < 1 is treated as 1 (keep everything).
func New(cfg Config) *Sampler {
	if cfg.Rate < 1 {
		cfg.Rate = 1
	}
	return &Sampler{
		cfg: cfg,
		rng: rand.New(rand.NewSource(42)), //nolint:gosec
	}
}

// Keep reports whether the current line should be retained.
// It is safe to call from a single goroutine; for concurrent use,
// create one Sampler per goroutine.
func (s *Sampler) Keep() bool {
	switch s.cfg.Strategy {
	case StrategyNone:
		return true
	case StrategyNth:
		n := s.counter.Add(1)
		return n%int64(s.cfg.Rate) == 1
	case StrategyRandom:
		return s.rng.Intn(s.cfg.Rate) == 0
	default:
		return true
	}
}

// Reset resets the internal counter (useful between files).
func (s *Sampler) Reset() {
	s.counter.Store(0)
}
