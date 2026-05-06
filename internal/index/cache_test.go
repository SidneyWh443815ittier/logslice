package index

import (
	"strings"
	"testing"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()
	idx := &Index{Entries: []Entry{{LineNum: 1}}}
	c.Set("key1", idx)

	got, ok := c.Get("key1")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(got.Entries))
	}
}

func TestCache_Miss(t *testing.T) {
	c := NewCache()
	_, ok := c.Get("missing")
	if ok {
		t.Error("expected cache miss")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := NewCache()
	c.Set("key1", &Index{})
	c.Invalidate("key1")
	_, ok := c.Get("key1")
	if ok {
		t.Error("expected key to be invalidated")
	}
}

func TestCache_GetOrBuild_Builds(t *testing.T) {
	c := NewCache()
	r := strings.NewReader("2024-01-15T10:00:00Z level=info msg=hello\n")

	idx, err := c.GetOrBuild("file1", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(idx.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(idx.Entries))
	}
}

func TestCache_GetOrBuild_UsesCached(t *testing.T) {
	c := NewCache()
	pre := &Index{Entries: []Entry{{LineNum: 99}}}
	c.Set("file2", pre)

	r := strings.NewReader("this should not be read\n")
	idx, err := c.GetOrBuild("file2", r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if idx.Entries[0].LineNum != 99 {
		t.Errorf("expected cached entry, got LineNum=%d", idx.Entries[0].LineNum)
	}
}
