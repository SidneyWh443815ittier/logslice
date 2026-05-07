package redact

import (
	"strings"
	"testing"
)

func TestConfig_IsEnabled_False(t *testing.T) {
	c := Config{}
	if c.IsEnabled() {
		t.Error("empty config should not be enabled")
	}
}

func TestConfig_IsEnabled_Patterns(t *testing.T) {
	c := Config{Patterns: []string{`\d+`}}
	if !c.IsEnabled() {
		t.Error("config with patterns should be enabled")
	}
}

func TestConfig_IsEnabled_Fields(t *testing.T) {
	c := Config{Fields: true}
	if !c.IsEnabled() {
		t.Error("config with Fields=true should be enabled")
	}
}

func TestConfig_Build_NoRules(t *testing.T) {
	c := Config{}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.HasRules() {
		t.Error("expected no rules for empty config")
	}
}

func TestConfig_Build_CustomMask(t *testing.T) {
	c := Config{
		Patterns: []string{`\d{4}`},
		Mask:     "<NUM>",
	}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("code=1234")
	if !strings.Contains(got, "<NUM>") {
		t.Errorf("expected custom mask <NUM>, got %q", got)
	}
}

func TestConfig_Build_InvalidPattern(t *testing.T) {
	c := Config{Patterns: []string{`[bad`}}
	_, err := c.Build()
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestConfig_Build_FieldsRedactsToken(t *testing.T) {
	c := Config{Fields: true}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := "token=abc123xyz"
	got := r.Apply(line)
	if strings.Contains(got, "abc123xyz") {
		t.Errorf("token value not redacted: %q", got)
	}
}

func TestConfig_Build_DefaultMask(t *testing.T) {
	c := Config{Patterns: []string{`secret`}}
	r, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Apply("my secret value")
	if !strings.Contains(got, DefaultMask) {
		t.Errorf("expected default mask %q in output %q", DefaultMask, got)
	}
}
