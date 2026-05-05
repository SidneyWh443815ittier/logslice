package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls how matched log lines are emitted.
type Format int

const (
	FormatRaw  Format = iota // original line, unchanged
	FormatJSON               // wrap line in a JSON envelope
	FormatColor              // highlight matched fields with ANSI codes
)

// Formatter writes log lines to an output destination.
type Formatter struct {
	w      io.Writer
	fmt    Format
	fields []string // fields to highlight in color mode
}

// NewFormatter creates a Formatter that writes to w in the given format.
func NewFormatter(w io.Writer, f Format, fields []string) *Formatter {
	return &Formatter{w: w, fmt: f, fields: fields}
}

// envelope is the JSON wrapper for FormatJSON mode.
type envelope struct {
	Line   uint64 `json:"line"`
	Raw    string `json:"raw"`
}

// Write emits a single log line, applying the configured format.
func (f *Formatter) Write(lineNum uint64, line string) error {
	switch f.fmt {
	case FormatJSON:
		return f.writeJSON(lineNum, line)
	case FormatColor:
		return f.writeColor(line)
	default:
		_, err := fmt.Fprintln(f.w, line)
		return err
	}
}

func (f *Formatter) writeJSON(lineNum uint64, line string) error {
	env := envelope{Line: lineNum, Raw: line}
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(f.w, string(b))
	return err
}

const (
	ansiYellow = "\033[33m"
	ansiReset  = "\033[0m"
)

func (f *Formatter) writeColor(line string) error {
	result := line
	for _, field := range f.fields {
		if field == "" {
			continue
		}
		result = strings.ReplaceAll(result, field, ansiYellow+field+ansiReset)
	}
	_, err := fmt.Fprintln(f.w, result)
	return err
}
