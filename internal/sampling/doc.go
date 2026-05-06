// Package sampling implements log-line sampling strategies for logslice.
//
// When processing very large log files it is often useful to inspect only a
// representative subset of lines rather than every entry.  Two strategies are
// supported:
//
//   - StrategyNth  – deterministic; keeps every N-th line (e.g. Rate=10 keeps
//     lines 1, 11, 21, …).  Useful when you need reproducible results.
//
//   - StrategyRandom – probabilistic; each line is independently retained with
//     probability 1/Rate.  Useful when you want an unbiased statistical sample.
//
// StrategyNone (the default) disables sampling so all lines pass through,
// preserving backward-compatible behaviour.
//
// Usage:
//
//	s := sampling.New(sampling.Config{
//		Strategy: sampling.StrategyNth,
//		Rate:     100,
//	})
//	for scanner.Scan() {
//		if s.Keep() {
//			process(scanner.Text())
//		}
//	}
package sampling
