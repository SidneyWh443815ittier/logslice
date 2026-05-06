// Package multifile provides utilities for merging and processing
// multiple log files in chronological order.
package multifile

import (
	"bufio"
	"container/heap"
	"io"

	"github.com/user/logslice/internal/timerange"
)

// entry holds a single log line along with its source file and parsed timestamp.
type entry struct {
	line      string
	lineNum   int
	source    string
	timestamp int64 // Unix nano; -1 if unparseable
	readerIdx int
}

// entryHeap is a min-heap of entries ordered by timestamp.
type entryHeap []*entry

func (h entryHeap) Len() int { return len(h) }
func (h entryHeap) Less(i, j int) bool {
	if h[i].timestamp != h[j].timestamp {
		return h[i].timestamp < h[j].timestamp
	}
	// Stable: preserve source order on tie.
	return h[i].readerIdx < h[j].readerIdx
}
func (h entryHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *entryHeap) Push(x any)   { *h = append(*h, x.(*entry)) }
func (h *entryHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// Source describes a single log source to be merged.
type Source struct {
	Name   string
	Reader io.Reader
}

// MergedLine is emitted by Merge for each log line in chronological order.
type MergedLine struct {
	Line    string
	Source  string
	LineNum int
}

// Merge reads from multiple Sources and emits lines in timestamp order.
// Lines whose timestamps cannot be parsed are emitted in source-arrival order
// relative to other unparseable lines from the same source.
func Merge(sources []Source, out chan<- MergedLine) error {
	scanners := make([]*bufio.Scanner, len(sources))
	lineNums := make([]int, len(sources))
	for i, s := range sources {
		scanners[i] = bufio.NewScanner(s.Reader)
	}

	h := &entryHeap{}
	heap.Init(h)

	// Seed the heap with the first line from each source.
	for i, sc := range scanners {
		if sc.Scan() {
			lineNums[i]++
			e := makeEntry(sc.Text(), lineNums[i], sources[i].Name, i)
			heap.Push(h, e)
		}
		if err := sc.Err(); err != nil {
			return err
		}
	}

	for h.Len() > 0 {
		min := heap.Pop(h).(*entry)
		out <- MergedLine{Line: min.line, Source: min.source, LineNum: min.lineNum}

		sc := scanners[min.readerIdx]
		if sc.Scan() {
			lineNums[min.readerIdx]++
			e := makeEntry(sc.Text(), lineNums[min.readerIdx], sources[min.readerIdx].Name, min.readerIdx)
			heap.Push(h, e)
		}
		if err := sc.Err(); err != nil {
			return err
		}
	}
	return nil
}

func makeEntry(line string, lineNum int, source string, idx int) *entry {
	ts, err := timerange.ExtractTimestamp(line)
	var nano int64 = -1
	if err == nil {
		nano = ts.UnixNano()
	}
	return &entry{line: line, lineNum: lineNum, source: source, timestamp: nano, readerIdx: idx}
}
