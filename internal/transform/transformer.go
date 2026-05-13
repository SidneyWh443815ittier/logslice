// Package transform provides line transformation utilities for logslice,
// allowing field extraction, renaming, and value substitution on log lines
// before they are passed to output formatters.
package transform

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Op represents a single transformation operation.
type Op struct {
	Kind  string // "rename", "drop", "set"
	Field string
	Value string // used by "rename" (new name) and "set" (literal value)
}

// Transformer applies a sequence of Ops to log lines.
type Transformer struct {
	ops []Op
}

// New creates a Transformer from the provided operations.
func New(ops []Op) *Transformer {
	return &Transformer{ops: ops}
}

// Apply transforms a single log line. If the line is valid JSON the ops are
// applied field-by-field and the result is re-serialised. Plain-text lines
// are returned unchanged unless a "set" op appends a key=value suffix.
func (t *Transformer) Apply(line string) (string, error) {
	if len(t.ops) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err == nil {
		for _, op := range t.ops {
			switch op.Kind {
			case "rename":
				if v, ok := obj[op.Field]; ok {
					obj[op.Value] = v
					delete(obj, op.Field)
				}
			case "drop":
				delete(obj, op.Field)
			case "set":
				obj[op.Field] = op.Value
			}
		}
		b, err := json.Marshal(obj)
		if err != nil {
			return line, fmt.Errorf("transform: marshal: %w", err)
		}
		return string(b), nil
	}

	// Plain text — only "set" ops are applicable (appended as key=value).
	var parts []string
	for _, op := range t.ops {
		if op.Kind == "set" {
			parts = append(parts, op.Field+"="+op.Value)
		}
	}
	if len(parts) > 0 {
		return line + " " + strings.Join(parts, " "), nil
	}
	return line, nil
}
