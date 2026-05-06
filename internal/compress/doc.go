// Package compress provides transparent decompression support for logslice.
//
// It detects the compression format of a log file based on its file extension
// and returns an io.Reader that decompresses on the fly, allowing the rest of
// the pipeline to treat compressed and plain log files uniformly.
//
// Supported formats:
//
//   - Plain text (.log, .txt, and any unrecognised extension)
//   - Gzip       (.gz, .gzip)
//
// Usage:
//
//	r, closer, err := compress.OpenReader("app.log.gz")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer closer.Close()
//	scanner := bufio.NewScanner(r)
package compress
