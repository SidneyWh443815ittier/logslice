package index

import (
	"io"
	"sync"
)

// Cache stores built indexes keyed by a string identifier (e.g. file path + mtime).
type Cache struct {
	mu    sync.RWMutex
	store map[string]*Index
}

// NewCache creates an empty Cache.
func NewCache() *Cache {
	return &Cache{store: make(map[string]*Index)}
}

// Get returns a cached index if present.
func (c *Cache) Get(key string) (*Index, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	idx, ok := c.store[key]
	return idx, ok
}

// Set stores an index under the given key.
func (c *Cache) Set(key string, idx *Index) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = idx
}

// Invalidate removes a key from the cache.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

// GetOrBuild returns a cached index or builds one from the given ReadSeeker.
func (c *Cache) GetOrBuild(key string, r io.ReadSeeker) (*Index, error) {
	if idx, ok := c.Get(key); ok {
		return idx, nil
	}
	idx, err := Build(r)
	if err != nil {
		return nil, err
	}
	c.Set(key, idx)
	return idx, nil
}
