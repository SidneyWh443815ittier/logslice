package ratelimit

import "fmt"

// Config holds user-facing rate limit configuration.
type Config struct {
	// Enabled controls whether rate limiting is active.
	Enabled bool
	// LinesPerSecond is the maximum number of lines emitted per second.
	// Values <= 0 are treated as unlimited when Enabled is true.
	LinesPerSecond int
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if c.Enabled && c.LinesPerSecond < 0 {
		return fmt.Errorf("ratelimit: LinesPerSecond must be >= 0, got %d", c.LinesPerSecond)
	}
	return nil
}

// Build constructs a Limiter from the configuration.
func (c Config) Build() (*Limiter, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if !c.Enabled || c.LinesPerSecond <= 0 {
		return New(StrategyNone, 0), nil
	}
	return New(StrategyLinesPerSecond, c.LinesPerSecond), nil
}
