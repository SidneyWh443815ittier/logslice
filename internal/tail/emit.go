package tail

import "bytes"

// emitLines splits raw bytes on newlines and sends each complete line to ch.
func emitLines(data []byte, ch chan<- string) {
	for len(data) > 0 {
		idx := bytes.IndexByte(data, '\n')
		if idx == -1 {
			// No newline yet — partial line, emit as-is.
			ch <- string(data)
			return
		}
		line := data[:idx]
		// Strip trailing carriage return for Windows-style line endings.
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1]
		}
		if len(line) > 0 {
			ch <- string(line)
		}
		data = data[idx+1:]
	}
}
