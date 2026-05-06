// Package tail provides live tailing of log files for logslice.
//
// A Watcher polls a file at a configurable interval and emits new lines
// written after the watcher was started. It is intended to be used alongside
// the existing reader pipeline for real-time log monitoring use cases.
//
// Basic usage:
//
//	w := tail.NewWatcher("/var/log/app.log", 100*time.Millisecond)
//	if err := w.Start(); err != nil {
//		log.Fatal(err)
//	}
//	defer w.Stop()
//
//	for line := range w.Lines() {
//		fmt.Println(line)
//	}
package tail
