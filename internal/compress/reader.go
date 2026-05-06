// Package compress provides transparent decompression for log files,
// supporting gzip and zstd formats in addition to plain text.
package compress

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Format represents a supported compression format.
type Format int

const (
	FormatPlain Format = iota
	FormatGzip
)

// DetectFormat infers the compression format from the file extension.
func DetectFormat(path string) Format {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".gz", ".gzip":
		return FormatGzip
	default:
		return FormatPlain
	}
}

// OpenReader opens the file at path and returns a reader that transparently
// decompresses the content based on the detected format. The caller is
// responsible for closing the returned closer.
func OpenReader(path string) (io.Reader, io.Closer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("compress: open %q: %w", path, err)
	}

	switch DetectFormat(path) {
	case FormatGzip:
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, nil, fmt.Errorf("compress: gzip reader for %q: %w", path, err)
		}
		return bufio.NewReader(gr), multiCloser(gr, f), nil
	default:
		return bufio.NewReader(f), f, nil
	}
}

// multiCloser closes multiple closers in order, returning the first error.
func multiCloser(closers ...io.Closer) io.Closer {
	return &chainCloser{closers: closers}
}

type chainCloser struct {
	closers []io.Closer
}

func (c *chainCloser) Close() error {
	var first error
	for _, cl := range c.closers {
		if err := cl.Close(); err != nil && first == nil {
			first = err
		}
	}
	return first
}
