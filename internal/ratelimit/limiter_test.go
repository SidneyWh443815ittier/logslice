package ratelimit

import (
	"testing"
	"time"
)

func TestStrategyNone_AllowsAll(t *testing.T) {
	l := New(StrategyNone, 0)
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatalf("StrategyNone blocked line %d", i)
		}
	}
}

func TestStrategyLinesPerSecond_AllowsUpToRate(t *testing.T) {
	l := New(StrategyLinesPerSecond, 5)
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed, got %d", allowed)
	}
}

func TestStrategyLinesPerSecond_BlocksOverRate(t *testing.T) {
	l := New(StrategyLinesPerSecond, 3)
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected 4th call to be blocked")
	}
}

func TestReset_ClearsCount(t *testing.T) {
	l := New(StrategyLinesPerSecond, 2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("expected block before reset")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("expected allow after reset")
	}
}

func TestWindowRollover_AllowsAfterSecond(t *testing.T) {
	l := New(StrategyLinesPerSecond, 2)
	l.Allow()
	l.Allow()
	// Force the window to appear expired.
	l.window = time.Now().Add(-2 * time.Second)
	if !l.Allow() {
		t.Fatal("expected allow after window rollover")
	}
}

func TestStats_ReturnsCountAndRate(t *testing.T) {
	l := New(StrategyLinesPerSecond, 10)
	l.Allow()
	l.Allow()
	allowed, rate := l.Stats()
	if allowed != 2 {
		t.Fatalf("expected allowed=2, got %d", allowed)
	}
	if rate != 10 {
		t.Fatalf("expected rate=10, got %d", rate)
	}
}

func TestNew_ZeroRate_PassThrough(t *testing.T) {
	l := New(StrategyLinesPerSecond, 0)
	for i := 0; i < 50; i++ {
		if !l.Allow() {
			t.Fatal("zero rate should behave as no-op")
		}
	}
}
