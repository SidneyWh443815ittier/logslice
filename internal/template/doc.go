// Package template provides Go-template-based rendering for log lines.
//
// It supports two modes:
//
//   - Plain mode: the template receives a single map with the key "Line"
//     containing the raw log line string.
//
//   - JSON mode: the log line is first decoded as a JSON object and the
//     resulting map[string]any is passed directly to the template, giving
//     access to individual structured fields. If JSON decoding fails the
//     package falls back to plain mode, populating "Line" with the raw
//     string and "_parseError" with the decode error message.
//
// Example usage:
//
//	f, err := template.New(`{{ index . "level" | upper }}: {{ index . "msg" }}`, true)
//	if err != nil { ... }
//	out, err := f.Format(line)
//
// Templates are compiled once via New and are safe for concurrent use
// across multiple goroutines.
package template
