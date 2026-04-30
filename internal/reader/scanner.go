package reader

import (
	"bufio"
	"io"
	"os"
)

// Scanner wraps a buffered line-by-line reader for log files.
type Scanner struct {
	scanner *bufio.Scanner
	line    string
	lineNum int
}

// NewFileScanner opens a file and returns a Scanner for it.
func NewFileScanner(path string) (*Scanner, io.Closer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return NewScanner(f), f, nil
}

// NewScanner creates a Scanner from any io.Reader.
func NewScanner(r io.Reader) *Scanner {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &Scanner{scanner: s}
}

// Scan advances to the next line. Returns false when done or on error.
func (s *Scanner) Scan() bool {
	if s.scanner.Scan() {
		s.line = s.scanner.Text()
		s.lineNum++
		return true
	}
	return false
}

// Line returns the current line text.
func (s *Scanner) Line() string {
	return s.line
}

// LineNumber returns the 1-based current line number.
func (s *Scanner) LineNumber() int {
	return s.lineNum
}

// Err returns any scanner error (excluding io.EOF).
func (s *Scanner) Err() error {
	return s.scanner.Err()
}
