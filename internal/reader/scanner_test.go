package reader

import (
	"strings"
	"testing"
)

func TestScanner_BasicLines(t *testing.T) {
	input := "line one\nline two\nline three\n"
	s := NewScanner(strings.NewReader(input))

	expected := []string{"line one", "line two", "line three"}
	var got []string
	for s.Scan() {
		got = append(got, s.Line())
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, l := range expected {
		if got[i] != l {
			t.Errorf("line %d: expected %q, got %q", i+1, l, got[i])
		}
	}
}

func TestScanner_LineNumbers(t *testing.T) {
	input := "a\nb\nc"
	s := NewScanner(strings.NewReader(input))
	num := 0
	for s.Scan() {
		num++
		if s.LineNumber() != num {
			t.Errorf("expected line number %d, got %d", num, s.LineNumber())
		}
	}
}

func TestScanner_EmptyInput(t *testing.T) {
	s := NewScanner(strings.NewReader(""))
	if s.Scan() {
		t.Error("expected no lines for empty input")
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScanner_SingleLineNoNewline(t *testing.T) {
	s := NewScanner(strings.NewReader("only line"))
	if !s.Scan() {
		t.Fatal("expected one line")
	}
	if s.Line() != "only line" {
		t.Errorf("unexpected line: %q", s.Line())
	}
	if s.Scan() {
		t.Error("expected no more lines")
	}
}

func TestNewFileScanner_MissingFile(t *testing.T) {
	_, _, err := NewFileScanner("/nonexistent/path/to/file.log")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
