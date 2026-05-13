package transform

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoOps_PassThrough(t *testing.T) {
	tr := New(nil)
	out, err := tr.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"level":"info"}` {
		t.Fatalf("unexpected: %s", out)
	}
}

func TestApply_Rename_JSON(t *testing.T) {
	tr := New([]Op{{Kind: "rename", Field: "msg", Value: "message"}})
	out, err := tr.Apply(`{"msg":"hello","level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["msg"]; ok {
		t.Error("old field 'msg' should be removed")
	}
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
}

func TestApply_Drop_JSON(t *testing.T) {
	tr := New([]Op{{Kind: "drop", Field: "secret"}})
	out, err := tr.Apply(`{"secret":"s3cr3t","level":"warn"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["secret"]; ok {
		t.Error("field 'secret' should have been dropped")
	}
}

func TestApply_Set_JSON(t *testing.T) {
	tr := New([]Op{{Kind: "set", Field: "env", Value: "prod"}})
	out, err := tr.Apply(`{"level":"debug"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
}

func TestApply_Set_PlainText(t *testing.T) {
	tr := New([]Op{{Kind: "set", Field: "env", Value: "staging"}})
	out, err := tr.Apply("plain log line")
	if err != nil {
		t.Fatal(err)
	}
	if out != "plain log line env=staging" {
		t.Errorf("unexpected: %q", out)
	}
}

func TestApply_Rename_MissingField_NoError(t *testing.T) {
	tr := New([]Op{{Kind: "rename", Field: "nonexistent", Value: "other"}})
	out, err := tr.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["other"]; ok {
		t.Error("'other' should not appear when source field is absent")
	}
}
