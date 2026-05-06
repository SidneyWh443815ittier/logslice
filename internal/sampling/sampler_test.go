package sampling_test

import (
	"testing"

	"logslice/internal/sampling"
)

func TestNew_DefaultRate(t *testing.T) {
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyNone, Rate: 0})
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
	// Rate < 1 should not panic and should keep all lines.
	for i := 0; i < 10; i++ {
		if !s.Keep() {
			t.Error("StrategyNone should always return true")
		}
	}
}

func TestStrategyNone_KeepsAll(t *testing.T) {
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyNone, Rate: 5})
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Errorf("StrategyNone: line %d should be kept", i)
		}
	}
}

func TestStrategyNth_KeepsEveryNth(t *testing.T) {
	rate := 5
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyNth, Rate: rate})
	kept := 0
	total := 100
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	expected := total / rate
	if kept != expected {
		t.Errorf("StrategyNth rate=%d: kept %d, want %d", rate, kept, expected)
	}
}

func TestStrategyNth_Reset(t *testing.T) {
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyNth, Rate: 3})
	// consume some lines
	for i := 0; i < 7; i++ {
		s.Keep()
	}
	s.Reset()
	// after reset, first call should be kept (counter resets to 0, add gives 1)
	if !s.Keep() {
		t.Error("first call after Reset should be kept for Nth strategy")
	}
}

func TestStrategyRandom_ApproximateRate(t *testing.T) {
	rate := 10
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyRandom, Rate: rate})
	kept := 0
	total := 10000
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	// Expect roughly total/rate ± 20%
	expected := total / rate
	margin := expected / 5
	if kept < expected-margin || kept > expected+margin {
		t.Errorf("StrategyRandom rate=%d: kept %d, expected ~%d (±%d)", rate, kept, expected, margin)
	}
}

func TestStrategyRandom_Rate1_KeepsAll(t *testing.T) {
	s := sampling.New(sampling.Config{Strategy: sampling.StrategyRandom, Rate: 1})
	for i := 0; i < 50; i++ {
		if !s.Keep() {
			t.Errorf("rate=1 random: line %d should always be kept", i)
		}
	}
}
