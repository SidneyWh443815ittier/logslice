package redact

import (
	"regexp"
	"strings"
	"testing"
)

func TestNew_NoRules_PassThrough(t *testing.T) {
	r := New(nil)
	line := "user=alice password=secret123"
	if got := r.Apply(line); got != line {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestApply_SingleRule(t *testing.T) {
	re := regexp.MustCompile(`\d{4}-\d{4}-\d{4}-\d{4}`)
	r := New([]Rule{{Pattern: re, Replacement: "[CARD]"}})
	line := "card=1234-5678-9012-3456 processed"
	got := r.Apply(line)
	if strings.Contains(got, "1234-5678-9012-3456") {
		t.Errorf("card number not redacted: %q", got)
	}
	if !strings.Contains(got, "[CARD]") {
		t.Errorf("expected [CARD] placeholder, got %q", got)
	}
}

func TestApply_MultipleRules(t *testing.T) {
	r, err := NewFromPatterns([]string{
		`\d{3}-\d{2}-\d{4}`,  // SSN
		`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`, // email
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "ssn=123-45-6789 email=user@example.com"
	got := r.Apply(line)
	if strings.Contains(got, "123-45-6789") {
		t.Errorf("SSN not redacted: %q", got)
	}
	if strings.Contains(got, "user@example.com") {
		t.Errorf("email not redacted: %q", got)
	}
}

func TestNewFromPatterns_InvalidRegex(t *testing.T) {
	_, err := NewFromPatterns([]string{`[invalid`})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestApplyAll(t *testing.T) {
	re := regexp.MustCompile(`secret`)
	r := New([]Rule{{Pattern: re, Replacement: "***"}})
	lines := []string{"ok line", "has secret here", "another secret"}
	got := r.ApplyAll(lines)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
	if strings.Contains(got[1], "secret") || strings.Contains(got[2], "secret") {
		t.Errorf("secrets not redacted: %v", got)
	}
	if got[0] != "ok line" {
		t.Errorf("unaffected line changed: %q", got[0])
	}
}

func TestHasRules(t *testing.T) {
	if New(nil).HasRules() {
		t.Error("expected HasRules=false for empty redactor")
	}
	re := regexp.MustCompile(`x`)
	if !New([]Rule{{Pattern: re, Replacement: ""}}).HasRules() {
		t.Error("expected HasRules=true")
	}
}

func TestNewDefaultRedactor_RedactsPassword(t *testing.T) {
	r, err := NewDefaultRedactor()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := `{"password":"hunter2","user":"alice"}`
	got := r.Apply(line)
	if strings.Contains(got, "hunter2") {
		t.Errorf("password not redacted: %q", got)
	}
}
