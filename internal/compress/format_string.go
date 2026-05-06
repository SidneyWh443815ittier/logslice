package compress

import "fmt"

// String returns a human-readable name for the Format.
func (f Format) String() string {
	switch f {
	case FormatPlain:
		return "plain"
	case FormatGzip:
		return "gzip"
	default:
		return fmt.Sprintf("Format(%d)", int(f))
	}
}

// Extension returns the canonical file extension associated with the format,
// including the leading dot. FormatPlain returns an empty string.
func (f Format) Extension() string {
	switch f {
	case FormatGzip:
		return ".gz"
	default:
		return ""
	}
}
