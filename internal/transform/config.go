package transform

import (
	"fmt"
	"strings"
)

// Config holds the user-facing configuration for the transform module.
type Config struct {
	// Ops is the ordered list of raw operation strings.
	// Supported syntax:
	//   rename:<old>:<new>
	//   drop:<field>
	//   set:<field>:<value>
	Ops []string `json:"ops" yaml:"ops"`
}

// IsEnabled reports whether any operations are configured.
func (c Config) IsEnabled() bool {
	return len(c.Ops) > 0
}

// Build parses the raw Op strings and returns a ready-to-use Transformer.
// An error is returned if any op string is malformed.
func (c Config) Build() (*Transformer, error) {
	ops := make([]Op, 0, len(c.Ops))
	for _, raw := range c.Ops {
		op, err := parseOp(raw)
		if err != nil {
			return nil, err
		}
		ops = append(ops, op)
	}
	return New(ops), nil
}

// parseOp converts a single raw string into an Op.
func parseOp(raw string) (Op, error) {
	parts := strings.SplitN(raw, ":", 3)
	if len(parts) < 2 {
		return Op{}, fmt.Errorf("transform: invalid op %q: expected kind:field[:value]", raw)
	}
	kind := strings.ToLower(parts[0])
	switch kind {
	case "drop":
		return Op{Kind: kind, Field: parts[1]}, nil
	case "rename", "set":
		if len(parts) < 3 {
			return Op{}, fmt.Errorf("transform: op %q requires a value", raw)
		}
		return Op{Kind: kind, Field: parts[1], Value: parts[2]}, nil
	default:
		return Op{}, fmt.Errorf("transform: unknown op kind %q", kind)
	}
}
