package ratelimit

import "testing"

func TestConfig_Validate_Valid(t *testing.T) {
	c := Config{Enabled: true, LinesPerSecond: 100}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfig_Validate_NegativeRate(t *testing.T) {
	c := Config{Enabled: true, LinesPerSecond: -1}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestConfig_Validate_DisabledNegativeRate(t *testing.T) {
	// Disabled configs skip rate validation.
	c := Config{Enabled: false, LinesPerSecond: -5}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error for disabled config: %v", err)
	}
}

func TestConfig_Build_Disabled(t *testing.T) {
	c := Config{Enabled: false}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.strategy != StrategyNone {
		t.Fatal("expected StrategyNone for disabled config")
	}
}

func TestConfig_Build_Enabled(t *testing.T) {
	c := Config{Enabled: true, LinesPerSecond: 50}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.strategy != StrategyLinesPerSecond {
		t.Fatal("expected StrategyLinesPerSecond")
	}
	if l.rate != 50 {
		t.Fatalf("expected rate=50, got %d", l.rate)
	}
}

func TestConfig_Build_ZeroRate_PassThrough(t *testing.T) {
	c := Config{Enabled: true, LinesPerSecond: 0}
	l, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.strategy != StrategyNone {
		t.Fatalf("expected StrategyNone for zero rate, got %v", l.strategy)
	}
}

func TestConfig_Build_InvalidReturnsError(t *testing.T) {
	c := Config{Enabled: true, LinesPerSecond: -10}
	_, err := c.Build()
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
