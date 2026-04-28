package timerange

import (
	"testing"
	"time"
)

func TestParseTimestamp_KnownFormats(t *testing.T) {
	cases := []struct {
		input string
		wantYear int
	}{
		{"2024-03-15T10:20:30Z", 2024},
		{"2024-03-15 10:20:30", 2024},
		{"2024-03-15T10:20:30.123", 2024},
		{"15/Mar/2024:10:20:30 +0000", 2024},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseTimestamp(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.Year() != tc.wantYear {
				t.Errorf("year: got %d, want %d", got.Year(), tc.wantYear)
			}
		})
	}
}

func TestParseTimestamp_Unknown(t *testing.T) {
	_, err := ParseTimestamp("not-a-date")
	if err == nil {
		t.Fatal("expected error for unrecognised format")
	}
}

func TestParseRange_Validation(t *testing.T) {
	_, err := ParseRange("2024-03-15T12:00:00Z", "2024-03-15T10:00:00Z")
	if err == nil {
		t.Fatal("expected error when to < from")
	}
}

func TestRange_Contains(t *testing.T) {
	r, _ := ParseRange("2024-03-15T08:00:00Z", "2024-03-15T18:00:00Z")

	mid, _ := ParseTimestamp("2024-03-15T12:00:00Z")
	before, _ := ParseTimestamp("2024-03-15T07:00:00Z")
	after, _ := ParseTimestamp("2024-03-15T19:00:00Z")

	if !r.Contains(mid) {
		t.Error("expected mid to be contained")
	}
	if r.Contains(before) {
		t.Error("expected before to be excluded")
	}
	if r.Contains(after) {
		t.Error("expected after to be excluded")
	}
}

func TestRange_OpenEnded(t *testing.T) {
	r, _ := ParseRange("2024-03-15T08:00:00Z", "")

	far, _ := time.Parse(time.RFC3339, "2099-01-01T00:00:00Z")
	if !r.Contains(far) {
		t.Error("open-ended range should contain far future")
	}
}
