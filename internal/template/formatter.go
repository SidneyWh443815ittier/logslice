// Package template provides go-template based log line formatting.
package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	gotemplate "text/template"
)

// Formatter applies a Go template to each log line, optionally parsing
// the line as JSON to expose structured fields inside the template.
type Formatter struct {
	tmpl    *gotemplate.Template
	jsonMode bool
}

// New compiles the given template string and returns a Formatter.
// If jsonMode is true, each line is parsed as JSON before rendering;
// the template receives a map[string]any. In plain mode the template
// receives a single string value bound to .Line.
func New(tmplStr string, jsonMode bool) (*Formatter, error) {
	if strings.TrimSpace(tmplStr) == "" {
		return nil, fmt.Errorf("template: template string must not be empty")
	}
	t, err := gotemplate.New("logslice").Option("missingkey=zero").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template: parse error: %w", err)
	}
	return &Formatter{tmpl: t, jsonMode: jsonMode}, nil
}

// Format renders the template against the given log line.
// It returns the rendered string and any execution error.
func (f *Formatter) Format(line string) (string, error) {
	var data any
	if f.jsonMode {
		m := make(map[string]any)
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			// Fall back to wrapping the raw line so the template still runs.
			m["Line"] = line
			m["_parseError"] = err.Error()
		}
		data = m
	} else {
		data = map[string]any{"Line": line}
	}

	var buf bytes.Buffer
	if err := f.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template: execute error: %w", err)
	}
	return buf.String(), nil
}
