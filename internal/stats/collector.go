package stats

import (
	"fmt"
	"time"
)

// Collector tracks processing statistics for a log filtering run.
type Collector struct {
	StartTime     time.Time
	LinesRead     int64
	LinesMatched  int64
	LinesSkipped  int64
	BytesRead     int64
	ParseErrors   int64
}

// New creates a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{
		StartTime: time.Now(),
	}
}

// RecordLine records a scanned line and its byte length.
func (c *Collector) RecordLine(bytes int, matched bool) {
	c.LinesRead++
	c.BytesRead += int64(bytes)
	if matched {
		c.LinesMatched++
	} else {
		c.LinesSkipped++
	}
}

// RecordParseError increments the parse error counter.
func (c *Collector) RecordParseError() {
	c.ParseErrors++
}

// Elapsed returns the duration since the collector was created.
func (c *Collector) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

// Summary returns a human-readable summary string.
func (c *Collector) Summary() string {
	elapsed := c.Elapsed()
	throughput := float64(c.BytesRead) / elapsed.Seconds() / 1024 / 1024
	return fmt.Sprintf(
		"read %d lines (%d matched, %d skipped) | %.2f MB/s | %d parse errors | elapsed %s",
		c.LinesRead, c.LinesMatched, c.LinesSkipped,
		throughput, c.ParseErrors, elapsed.Round(time.Millisecond),
	)
}
