# logslice

A fast log filtering utility that supports time-range extraction and structured field queries on large log files.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

---

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" app.log

# Query by structured field
logslice --field level=error app.log

# Combine time range and field query
logslice --from "2024-01-15T08:00:00" --to "2024-01-15T09:00:00" --field service=api app.log

# Read from stdin
cat app.log | logslice --field level=warn
```

### Flags

| Flag | Description |
|------|-------------|
| `--from` | Start of time range (RFC3339 format) |
| `--to` | End of time range (RFC3339 format) |
| `--field` | Structured field filter in `key=value` format |
| `--output` | Output file path (defaults to stdout) |

---

## Features

- Efficient line-by-line streaming for large files
- Supports JSON and logfmt structured log formats
- Binary search on time-sorted log files for fast range extraction
- Composable filters via multiple `--field` flags

---

## License

MIT © 2024 [yourusername](https://github.com/yourusername)