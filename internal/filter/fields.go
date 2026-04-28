package filter

import (
	"strings"
)

// ParseFields extracts structured key=value pairs from a log line.
// It supports space-separated tokens of the form key=value.
// Quoted values (key="some value") are handled as well.
func ParseFields(line string) map[string]string {
	fields := make(map[string]string)
	tokens := tokenize(line)
	for _, tok := range tokens {
		idx := strings.IndexByte(tok, '=')
		if idx <= 0 || idx == len(tok)-1 {
			continue
		}
		key := tok[:idx]
		val := tok[idx+1:]
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		fields[key] = val
	}
	return fields
}

// tokenize splits a line into tokens respecting double-quoted strings.
func tokenize(line string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false

	for i := 0; i < len(line); i++ {
		ch := line[i]
		switch {
		case ch == '"':
			inQuote = !inQuote
			cur.WriteByte(ch)
		case ch == ' ' && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(ch)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}
