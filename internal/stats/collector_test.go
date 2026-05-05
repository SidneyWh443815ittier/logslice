package stats

import (
	"strings"
	"testing"
	"time"
)

func TestNew_InitializesStartTime(t *testing.T) {
	before := time.Now()
	c := New()
	after := time.Now()

	if c.StartTime.Before(before) || c.StartTime.After(after) {
		t.Errorf("StartTime %v not in expected range [%v, %v]", c.StartTime, before, after)
	}
}

func TestRecordLine_Matched(t *testing.T) {
	c := New()
	c.RecordLine(100, true)
	c.RecordLine(200, true)

	if c.LinesRead != 2 {
		t.Errorf("expected LinesRead=2, got %d", c.LinesRead)
	}
	if c.LinesMatched != 2 {
		t.Errorf("expected LinesMatched=2, got %d", c.LinesMatched)
	}
	if c.LinesSkipped != 0 {
		t.Errorf("expected LinesSkipped=0, got %d", c.LinesSkipped)
	}
	if c.BytesRead != 300 {
		t.Errorf("expected BytesRead=300, got %d", c.BytesRead)
	}
}

func TestRecordLine_Skipped(t *testing.T) {
	c := New()
	c.RecordLine(50, false)

	if c.LinesSkipped != 1 {
		t.Errorf("expected LinesSkipped=1, got %d", c.LinesSkipped)
	}
	if c.LinesMatched != 0 {
		t.Errorf("expected LinesMatched=0, got %d", c.LinesMatched)
	}
}

func TestRecordParseError(t *testing.T) {
	c := New()
	c.RecordParseError()
	c.RecordParseError()

	if c.ParseErrors != 2 {
		t.Errorf("expected ParseErrors=2, got %d", c.ParseErrors)
	}
}

func TestSummary_ContainsKeyFields(t *testing.T) {
	c := New()
	c.RecordLine(1024, true)
	c.RecordLine(512, false)
	c.RecordParseError()

	// small sleep so elapsed is non-zero
	time.Sleep(2 * time.Millisecond)

	s := c.Summary()
	for _, want := range []string{"2", "1 matched", "1 skipped", "1 parse errors"} {
		if !strings.Contains(s, want) {
			t.Errorf("Summary() missing %q in: %s", want, s)
		}
	}
}

func TestElapsed_PositiveDuration(t *testing.T) {
	c := New()
	time.Sleep(1 * time.Millisecond)
	if c.Elapsed() <= 0 {
		t.Error("expected positive elapsed duration")
	}
}
