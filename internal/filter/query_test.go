package filter

import (
	"testing"
)

func TestParseQuery_Valid(t *testing.T) {
	q, err := ParseQuery([]string{"level=error", "svc~auth", "host^web-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(q.Conditions) != 3 {
		t.Fatalf("expected 3 conditions, got %d", len(q.Conditions))
	}
}

func TestParseQuery_Invalid(t *testing.T) {
	_, err := ParseQuery([]string{"nodoperator"})
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestParseQuery_EmptyField(t *testing.T) {
	_, err := ParseQuery([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestQuery_Matches_Equals(t *testing.T) {
	q, _ := ParseQuery([]string{"level=error"})
	if !q.Matches(map[string]string{"level": "error"}) {
		t.Error("expected match")
	}
	if q.Matches(map[string]string{"level": "info"}) {
		t.Error("expected no match")
	}
}

func TestQuery_Matches_Contains(t *testing.T) {
	q, _ := ParseQuery([]string{"msg~timeout"})
	if !q.Matches(map[string]string{"msg": "connection timeout reached"}) {
		t.Error("expected match")
	}
	if q.Matches(map[string]string{"msg": "all good"}) {
		t.Error("expected no match")
	}
}

func TestQuery_Matches_Prefix(t *testing.T) {
	q, _ := ParseQuery([]string{"host^web-"})
	if !q.Matches(map[string]string{"host": "web-01"}) {
		t.Error("expected match")
	}
	if q.Matches(map[string]string{"host": "db-01"}) {
		t.Error("expected no match")
	}
}

func TestQuery_Matches_MissingField(t *testing.T) {
	q, _ := ParseQuery([]string{"level=error"})
	if q.Matches(map[string]string{"msg": "something"}) {
		t.Error("expected no match when field is absent")
	}
}

func TestParseFields_Basic(t *testing.T) {
	fields := ParseFields(`time=2024-01-01T00:00:00Z level=info msg="user logged in" svc=auth`)
	if fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", fields["level"])
	}
	if fields["msg"] != "user logged in" {
		t.Errorf("expected msg unquoted, got %q", fields["msg"])
	}
	if fields["svc"] != "auth" {
		t.Errorf("expected svc=auth, got %q", fields["svc"])
	}
}

func TestParseFields_NoFields(t *testing.T) {
	fields := ParseFields("plain log line with no kv pairs")
	if len(fields) != 0 {
		t.Errorf("expected empty fields, got %v", fields)
	}
}
