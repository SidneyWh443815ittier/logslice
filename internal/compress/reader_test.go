package compress_test

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"logslice/internal/compress"
)

func TestDetectFormat_Gzip(t *testing.T) {
	for _, name := range []string{"app.log.gz", "app.log.gzip", "APP.LOG.GZ"} {
		if got := compress.DetectFormat(name); got != compress.FormatGzip {
			t.Errorf("DetectFormat(%q) = %v, want FormatGzip", name, got)
		}
	}
}

func TestDetectFormat_Plain(t *testing.T) {
	for _, name := range []string{"app.log", "syslog", "out.txt"} {
		if got := compress.DetectFormat(name); got != compress.FormatPlain {
			t.Errorf("DetectFormat(%q) = %v, want FormatPlain", name, got)
		}
	}
}

func TestOpenReader_Plain(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.log")
	want := "hello world\nsecond line\n"
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatal(err)
	}

	r, closer, err := compress.OpenReader(path)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}
	defer closer.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestOpenReader_Gzip(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.log.gz")
	want := "compressed line one\ncompressed line two\n"

	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gw := gzip.NewWriter(f)
	if _, err := io.Copy(gw, strings.NewReader(want)); err != nil {
		t.Fatal(err)
	}
	gw.Close()
	f.Close()

	r, closer, err := compress.OpenReader(path)
	if err != nil {
		t.Fatalf("OpenReader: %v", err)
	}
	defer closer.Close()

	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestOpenReader_MissingFile(t *testing.T) {
	_, _, err := compress.OpenReader("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
