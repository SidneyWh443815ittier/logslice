package highlight

import (
	"strings"
	"testing"
)

func TestFindMatches_SingleTerm(t *testing.T) {
	matches := FindMatches("error: disk full", []string{"error"})
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0].Start != 0 || matches[0].End != 5 {
		t.Errorf("unexpected match bounds: %+v", matches[0])
	}
}

func TestFindMatches_CaseInsensitive(t *testing.T) {
	matches := FindMatches("ERROR: disk full", []string{"error"})
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
}

func TestFindMatches_MultipleTerms(t *testing.T) {
	matches := FindMatches("warn: low memory, error: oom", []string{"warn", "error"})
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestFindMatches_NoMatch(t *testing.T) {
	matches := FindMatches("all good here", []string{"error"})
	if len(matches) != 0 {
		t.Errorf("expected no matches, got %d", len(matches))
	}
}

func TestFindMatches_EmptyTerm(t *testing.T) {
	matches := FindMatches("some line", []string{""})
	if len(matches) != 0 {
		t.Errorf("empty term should produce no matches")
	}
}

func TestFindMatches_RepeatedTerm(t *testing.T) {
	matches := FindMatches("foo foo foo", []string{"foo"})
	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
}

func TestApplyANSI_WrapsMatch(t *testing.T) {
	line := "error: disk full"
	matches := FindMatches(line, []string{"error"})
	result := ApplyANSI(line, matches)
	if !strings.Contains(result, AnsiBold) {
		t.Error("expected ANSI bold code in output")
	}
	if !strings.Contains(result, AnsiReset) {
		t.Error("expected ANSI reset code in output")
	}
	if !strings.Contains(result, "error") {
		t.Error("original term should still appear in output")
	}
}

func TestApplyANSI_NoMatches_ReturnsOriginal(t *testing.T) {
	line := "all good here"
	result := ApplyANSI(line, nil)
	if result != line {
		t.Errorf("expected original line, got %q", result)
	}
}

func TestApplyANSI_PreservesNonMatchedText(t *testing.T) {
	line := "prefix error suffix"
	matches := FindMatches(line, []string{"error"})
	result := ApplyANSI(line, matches)
	if !strings.Contains(result, "prefix") || !strings.Contains(result, "suffix") {
		t.Error("non-matched text should be preserved")
	}
}
