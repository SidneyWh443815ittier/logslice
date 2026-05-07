package context

import (
	"testing"
)

func TestNew_DefaultsNegativeToZero(t *testing.T) {
	l := New(-1, -3)
	if l.before != 0 || l.after != 0 {
		t.Fatalf("expected before=0 after=0, got before=%d after=%d", l.before, l.after)
	}
}

func TestFeed_NoContext_MatchOnly(t *testing.T) {
	l := New(0, 0)
	got := l.Feed("match", true)
	if len(got) != 1 || got[0] != "match" {
		t.Fatalf("expected [match], got %v", got)
	}
	got = l.Feed("skip", false)
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestFeed_AfterContext(t *testing.T) {
	l := New(0, 2)
	l.Feed("a", false)
	out := l.Feed("match", true)
	if len(out) != 1 || out[0] != "match" {
		t.Fatalf("unexpected match output: %v", out)
	}
	out = l.Feed("after1", false)
	if len(out) != 1 || out[0] != "after1" {
		t.Fatalf("expected after1, got %v", out)
	}
	out = l.Feed("after2", false)
	if len(out) != 1 || out[0] != "after2" {
		t.Fatalf("expected after2, got %v", out)
	}
	out = l.Feed("beyond", false)
	if len(out) != 0 {
		t.Fatalf("expected empty after window, got %v", out)
	}
}

func TestFeed_BeforeContext(t *testing.T) {
	l := New(2, 0)
	l.Feed("b1", false)
	l.Feed("b2", false)
	l.Feed("b3", false) // oldest, should be evicted
	out := l.Feed("match", true)
	if len(out) != 3 {
		t.Fatalf("expected 3 lines (2 before + match), got %d: %v", len(out), out)
	}
	if out[0] != "b2" || out[1] != "b3" || out[2] != "match" {
		t.Fatalf("unexpected order: %v", out)
	}
}

func TestFeed_BeforeAndAfterContext(t *testing.T) {
	l := New(1, 1)
	l.Feed("pre", false)
	out := l.Feed("match", true)
	if len(out) != 2 || out[0] != "pre" || out[1] != "match" {
		t.Fatalf("unexpected output: %v", out)
	}
	out = l.Feed("post", false)
	if len(out) != 1 || out[0] != "post" {
		t.Fatalf("expected post, got %v", out)
	}
}

func TestReset_ClearsState(t *testing.T) {
	l := New(2, 2)
	l.Feed("x", false)
	l.Feed("match", true)
	l.Reset()
	// After reset the pre-buffer should be empty.
	out := l.Feed("match2", true)
	if len(out) != 1 || out[0] != "match2" {
		t.Fatalf("expected only match2 after reset, got %v", out)
	}
}
