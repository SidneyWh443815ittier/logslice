package aggregate

import (
	"errors"
	"fmt"
)

// Strategy controls how aggregation results are presented.
type Strategy int

const (
	// StrategyCount groups lines by field value and reports counts.
	StrategyCount Strategy = iota
	// StrategyTop limits output to the N most frequent entries.
	StrategyTop
)

// Config holds user-supplied options for the aggregate module.
type Config struct {
	// Field is the log field to aggregate on (required).
	Field string
	// Strategy selects the presentation mode.
	Strategy Strategy
	// TopN is the maximum number of entries to return when Strategy == StrategyTop.
	// Values <= 0 are treated as "all".
	TopN int
}

// Validate returns an error if the configuration is invalid.
func (c Config) Validate() error {
	if c.Field == "" {
		return errors.New("aggregate: field must not be empty")
	}
	if c.Strategy == StrategyTop && c.TopN < 0 {
		return fmt.Errorf("aggregate: top_n must be >= 0, got %d", c.TopN)
	}
	return nil
}

// Build constructs a Counter from the configuration.
// Validate must pass before calling Build.
func (c Config) Build() (*Counter, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return New(c.Field), nil
}

// Apply filters entries according to the strategy.
func (c Config) Apply(entries []Entry) []Entry {
	if c.Strategy == StrategyTop && c.TopN > 0 && len(entries) > c.TopN {
		return entries[:c.TopN]
	}
	return entries
}
