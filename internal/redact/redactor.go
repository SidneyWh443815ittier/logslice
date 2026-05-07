// Package redact provides utilities for masking sensitive fields in log lines
// before output, supporting regex-based and field-key-based redaction rules.
package redact

import (
	"regexp"
	"strings"
)

// Rule describes a single redaction rule.
type Rule struct {
	// Pattern is the compiled regex to match sensitive data.
	Pattern *regexp.Regexp
	// Replacement is the string substituted for matched content.
	Replacement string
}

// Redactor applies a set of Rules to log lines.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor from the provided rules.
func New(rules []Rule) *Redactor {
	return &Redactor{rules: rules}
}

// NewFromPatterns compiles a Redactor from raw regex pattern strings.
// Each pattern is paired with the default mask "[REDACTED]".
func NewFromPatterns(patterns []string) (*Redactor, error) {
	rules := make([]Rule, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{Pattern: re, Replacement: "[REDACTED]"})
	}
	return New(rules), nil
}

// Apply returns a copy of line with all rule patterns replaced by their
// respective replacement strings.
func (r *Redactor) Apply(line string) string {
	if len(r.rules) == 0 {
		return line
	}
	result := line
	for _, rule := range r.rules {
		result = rule.Pattern.ReplaceAllString(result, rule.Replacement)
	}
	return result
}

// ApplyAll applies the redactor to every line in the slice and returns a new
// slice with redacted content.
func (r *Redactor) ApplyAll(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = r.Apply(l)
	}
	return out
}

// HasRules reports whether any redaction rules are configured.
func (r *Redactor) HasRules() bool {
	return len(r.rules) > 0
}

// defaultFieldPatterns maps common sensitive field names to capture-group
// patterns that redact their values in key=value or "key":"value" forms.
var defaultFieldPatterns = map[string]string{
	"password": `(?i)(password["']?\s*[:=]\s*)[^\s,}&"']+`,
	"token":    `(?i)(token["']?\s*[:=]\s*)[^\s,}&"']+`,
	"secret":   `(?i)(secret["']?\s*[:=]\s*)[^\s,}&"']+`,
	"apikey":   `(?i)(api_?key["']?\s*[:=]\s*)[^\s,}&"']+`,
}

// NewDefaultRedactor returns a Redactor pre-loaded with rules for common
// sensitive fields (password, token, secret, api_key).
func NewDefaultRedactor() (*Redactor, error) {
	rules := make([]Rule, 0, len(defaultFieldPatterns))
	for _, pattern := range defaultFieldPatterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		rules = append(rules, Rule{
			Pattern:     re,
			Replacement: "${1}[REDACTED]",
		})
	}
	_ = strings.ToLower // imported for future use
	return New(rules), nil
}
