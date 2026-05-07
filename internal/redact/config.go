package redact

import (
	"fmt"
	"regexp"
)

// Config holds the user-supplied redaction configuration, typically parsed
// from CLI flags or a config file.
type Config struct {
	// Patterns is a list of raw regular expression strings.
	Patterns []string `json:"patterns" yaml:"patterns"`
	// Fields enables the built-in sensitive-field rules (password, token, etc.).
	Fields bool `json:"fields" yaml:"fields"`
	// Mask overrides the default "[REDACTED]" replacement string.
	Mask string `json:"mask" yaml:"mask"`
}

// DefaultMask is used when Config.Mask is empty.
const DefaultMask = "[REDACTED]"

// Build constructs a Redactor from the Config, returning an error if any
// pattern fails to compile.
func (c Config) Build() (*Redactor, error) {
	mask := c.Mask
	if mask == "" {
		mask = DefaultMask
	}

	var rules []Rule

	if c.Fields {
		for name, pattern := range defaultFieldPatterns {
			re, err := regexp.Compile(pattern)
			if err != nil {
				return nil, fmt.Errorf("redact: built-in field %q: %w", name, err)
			}
			rules = append(rules, Rule{Pattern: re, Replacement: "${1}" + mask})
		}
	}

	for _, p := range c.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("redact: pattern %q: %w", p, err)
		}
		rules = append(rules, Rule{Pattern: re, Replacement: mask})
	}

	return New(rules), nil
}

// IsEnabled reports whether any redaction is configured.
func (c Config) IsEnabled() bool {
	return c.Fields || len(c.Patterns) > 0
}
