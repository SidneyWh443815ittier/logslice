package tail

import (
	"testing"
)

func TestEmitLines_SplitsOnNewline(t *testing.T) {
	ch := make(chan string, 10)
	emitLines([]byte("alpha\nbeta\ngamma\n"), ch)

	want := []string{"alpha", "beta", "gamma"}
	for i, w := range want {
		got := <-ch
		if got != w {
			t.Errorf("line[%d] = %q, want %q", i, got, w)
		}
	}
}

func TestEmitLines_SkipsEmptyLines(t *testing.T) {
	ch := make(chan string, 10)
	emitLines([]byte("first\n\nsecond\n"), ch)

	if got := <-ch; got != "first" {
		t.Errorf("got %q, want %q", got, "first")
	}
	if got := <-ch; got != "second" {
		t.Errorf("got %q, want %q", got, "second")
	}
	if len(ch) != 0 {
		t.Errorf("expected no remaining lines, got %d", len(ch))
	}
}

func TestEmitLines_WindowsLineEndings(t *testing.T) {
	ch := make(chan string, 10)
	emitLines([]byte("one\r\ntwo\r\n"), ch)

	if got := <-ch; got != "one" {
		t.Errorf("got %q, want %q", got, "one")
	}
	if got := <-ch; got != "two" {
		t.Errorf("got %q, want %q", got, "two")
	}
}

func TestEmitLines_PartialLine(t *testing.T) {
	ch := make(chan string, 10)
	emitLines([]byte("no-newline-here"), ch)

	if got := <-ch; got != "no-newline-here" {
		t.Errorf("got %q, want %q", got, "no-newline-here")
	}
}
