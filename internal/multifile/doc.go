// Package multifile implements chronological merging of multiple log files.
//
// When analysing distributed systems it is common to have one log file per
// service or host. multifile.Merge reads from any number of io.Reader sources
// concurrently and emits MergedLine values in ascending timestamp order so
// that downstream filters and formatters see a single, coherent event stream.
//
// Lines whose timestamps cannot be parsed are sorted by source-arrival order
// and are interleaved conservatively — they are emitted as soon as the heap
// would otherwise advance past them.
//
// Typical usage:
//
//	sources := []multifile.Source{
//		{Name: "api.log",     Reader: f1},
//		{Name: "worker.log", Reader: f2},
//	}
//	out := make(chan multifile.MergedLine, 64)
//	go func() {
//		defer close(out)
//		multifile.Merge(sources, out)
//	}()
//	for ml := range out {
//		fmt.Printf("[%s] %s\n", ml.Source, ml.Line)
//	}
package multifile
