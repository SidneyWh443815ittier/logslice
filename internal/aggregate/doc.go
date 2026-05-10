// Package aggregate provides field-based aggregation and counting for log lines.
//
// It is designed to answer questions such as "how many lines appeared per log
// level?" or "which service emitted the most errors?" without requiring the
// caller to parse structured formats — values are extracted directly from raw
// log text using lightweight key=value scanning compatible with both logfmt and
// common structured log styles.
//
// # Usage
//
//	c := aggregate.New("level")
//	for _, line := range lines {
//		c.Add(line)
//	}
//	for _, e := range c.Entries() {
//		fmt.Printf("%s\t%d\n", e.Key, e.Count)
//	}
//
// Entries are returned sorted by count descending so the most frequent values
// appear first. Ties are broken alphabetically.
package aggregate
