package filter

import (
	"fmt"
	"strings"
)

// Op represents a comparison operator.
type Op int

const (
	OpEquals Op = iota
	OpContains
	OpPrefix
)

// Condition represents a single key-value filter condition.
type Condition struct {
	Field string
	Value string
	Op    Op
}

// Query holds a set of conditions that must all match (AND semantics).
type Query struct {
	Conditions []Condition
}

// ParseQuery parses a slice of filter expressions of the form:
//   field=value     (exact match)
//   field~value     (contains)
//   field^value     (prefix match)
func ParseQuery(exprs []string) (*Query, error) {
	q := &Query{}
	for _, expr := range exprs {
		var c Condition
		switch {
		case strings.Contains(expr, "~"):
			parts := strings.SplitN(expr, "~", 2)
			c = Condition{Field: parts[0], Value: parts[1], Op: OpContains}
		case strings.Contains(expr, "^"):
			parts := strings.SplitN(expr, "^", 2)
			c = Condition{Field: parts[0], Value: parts[1], Op: OpPrefix}
		case strings.Contains(expr, "="):
			parts := strings.SplitN(expr, "=", 2)
			c = Condition{Field: parts[0], Value: parts[1], Op: OpEquals}
		default:
			return nil, fmt.Errorf("invalid filter expression: %q", expr)
		}
		if c.Field == "" {
			return nil, fmt.Errorf("empty field name in expression: %q", expr)
		}
		q.Conditions = append(q.Conditions, c)
	}
	return q, nil
}

// Matches reports whether the given key-value fields satisfy all conditions.
func (q *Query) Matches(fields map[string]string) bool {
	for _, c := range q.Conditions {
		v, ok := fields[c.Field]
		if !ok {
			return false
		}
		switch c.Op {
		case OpEquals:
			if v != c.Value {
				return false
			}
		case OpContains:
			if !strings.Contains(v, c.Value) {
				return false
			}
		case OpPrefix:
			if !strings.HasPrefix(v, c.Value) {
				return false
			}
		}
	}
	return true
}
