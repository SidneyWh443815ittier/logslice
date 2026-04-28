package timerange

import (
	"testing"
)

func TestExtractTimestamp(t *testing.T) {
	cases := []struct {
		line    string
		wantOK  bool
		wantYear int
	}{
		{
			line:     `2024-03-15T10:20:30Z INFO server started`,
			wantOK:   true,
			wantYear: 2024,
		},
		{
			line:     `[2024-03-15 10:20:30.456] ERROR disk full`,
			wantOK:   true,
			wantYear: 2024,
		},
		{
			line:     `192.168.1.1 - - [15/Mar/2024:10:20:30 +0000] "GET / HTTP/1.1" 200`,
			wantOK:   true,
			wantYear: 2024,
		},
		{
			line:    `no timestamp here at all`,
			wantOK:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.line[:min(30, len(tc.line))], func(t *testing.T) {
			got, ok := ExtractTimestamp(tc.line)
			if ok != tc.wantOK {
				t.Fatalf("ok: got %v, want %v", ok, tc.wantOK)
			}
			if ok && got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
		})
	}
}

func TestLineInRange(t *testing.T) {
	r, _ := ParseRange("2024-03-15T09:00:00Z", "2024-03-15T11:00:00Z")

	inside := `2024-03-15T10:00:00Z INFO ping`
	outside := `2024-03-15T12:00:00Z INFO pong`
	noTS := `startup message without timestamp`

	if !LineInRange(inside, r) {
		t.Error("expected inside line to be in range")
	}
	if LineInRange(outside, r) {
		t.Error("expected outside line to be out of range")
	}
	if LineInRange(noTS, r) {
		t.Error("expected no-timestamp line to be excluded when range is bounded")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
