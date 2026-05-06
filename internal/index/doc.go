// Package index provides an in-memory indexing layer for logslice.
//
// It scans log files and records the byte offset, line number, and parsed
// timestamp of each line. This allows downstream components to seek directly
// to relevant portions of large files rather than scanning from the beginning.
//
// Usage:
//
//	f, _ := os.Open("app.log")
//	idx, _ := index.Build(f)
//
//	// Filter entries by a time range before reading lines.
//	rng, _ := timerange.ParseRange("2024-01-15T10:00:00Z", "2024-01-15T11:00:00Z")
//	entries := idx.FilterByRange(rng)
//
// For repeated access to the same file, use Cache to avoid redundant scans:
//
//	cache := index.NewCache()
//	idx, _ := cache.GetOrBuild(filePath, f)
package index
