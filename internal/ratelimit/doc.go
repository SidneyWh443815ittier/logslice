// Package ratelimit implements output rate limiting for logslice.
//
// When tailing or streaming large log files, it can be useful to cap the
// number of lines emitted per second to avoid overwhelming downstream
// consumers or terminals.
//
// # Usage
//
//	cfg := ratelimit.Config{Enabled: true, LinesPerSecond: 200}
//	limiter, err := cfg.Build()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, line := range lines {
//	    if limiter.Allow() {
//	        fmt.Println(line)
//	    }
//	}
//
// # Strategies
//
//   - StrategyNone: no limiting; all lines pass through (default).
//   - StrategyLinesPerSecond: allows at most N lines per one-second window.
package ratelimit
