package dedupe

import (
	"fmt"
	"testing"
)

func TestStrategyNone_AllowsAll(t *testing.T) {
	d := New(StrategyNone, 0)
	for i := 0; i < 5; i++ {
		if d.IsDuplicate("same line") {
			t.Fatalf("StrategyNone should never report a duplicate (iteration %d)", i)
		}
	}
}

func TestStrategyExact_DetectsDuplicate(t *testing.T) {
	d := New(StrategyExact, 0)
	if d.IsDuplicate("hello") {
		t.Fatal("first occurrence should not be a duplicate")
	}
	if !d.IsDuplicate("hello") {
		t.Fatal("second occurrence should be a duplicate")
	}
}

func TestStrategyExact_DistinctLines(t *testing.T) {
	d := New(StrategyExact, 0)
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		if d.IsDuplicate(l) {
			t.Fatalf("distinct line %q should not be duplicate", l)
		}
	}
}

func TestStrategyExact_Reset(t *testing.T) {
	d := New(StrategyExact, 0)
	d.IsDuplicate("line")
	d.Reset()
	if d.IsDuplicate("line") {
		t.Fatal("after Reset, line should not be considered duplicate")
	}
}

func TestStrategyWindow_WithinWindow(t *testing.T) {
	d := New(StrategyWindow, 4)
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	if !d.IsDuplicate("a") {
		t.Fatal("'a' is within the window and should be detected as duplicate")
	}
}

func TestStrategyWindow_OutsideWindow(t *testing.T) {
	const windowSize = 3
	d := New(StrategyWindow, windowSize)
	d.IsDuplicate("target")
	// Push 'target' out of the window.
	for i := 0; i < windowSize; i++ {
		d.IsDuplicate(fmt.Sprintf("filler-%d", i))
	}
	if d.IsDuplicate("target") {
		t.Fatal("'target' should have been evicted from the window")
	}
}

func TestStrategyWindow_ZeroSize_AllowsAll(t *testing.T) {
	d := New(StrategyWindow, 0)
	for i := 0; i < 3; i++ {
		if d.IsDuplicate("repeat") {
			t.Fatal("zero-size window should never report a duplicate")
		}
	}
}

func TestStrategyWindow_Reset(t *testing.T) {
	d := New(StrategyWindow, 4)
	d.IsDuplicate("x")
	d.Reset()
	if d.IsDuplicate("x") {
		t.Fatal("after Reset, 'x' should not be in the window")
	}
}
