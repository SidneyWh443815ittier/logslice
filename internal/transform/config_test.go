package transform

import "testing"

func TestConfig_IsEnabled_False(t *testing.T) {
	c := Config{}
	if c.IsEnabled() {
		t.Error("expected IsEnabled=false for empty config")
	}
}

func TestConfig_IsEnabled_True(t *testing.T) {
	c := Config{Ops: []string{"drop:secret"}}
	if !c.IsEnabled() {
		t.Error("expected IsEnabled=true")
	}
}

func TestConfig_Build_NoOps(t *testing.T) {
	c := Config{}
	tr, err := c.Build()
	if err != nil {
		t.Fatal(err)
	}
	if tr == nil {
		t.Fatal("expected non-nil transformer")
	}
}

func TestConfig_Build_ValidOps(t *testing.T) {
	c := Config{Ops: []string{"rename:msg:message", "drop:token", "set:env:prod"}}
	tr, err := c.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tr.ops) != 3 {
		t.Errorf("expected 3 ops, got %d", len(tr.ops))
	}
}

func TestConfig_Build_InvalidOp(t *testing.T) {
	c := Config{Ops: []string{"badop"}}
	_, err := c.Build()
	if err == nil {
		t.Error("expected error for malformed op")
	}
}

func TestConfig_Build_UnknownKind(t *testing.T) {
	c := Config{Ops: []string{"uppercase:field"}}
	_, err := c.Build()
	if err == nil {
		t.Error("expected error for unknown op kind")
	}
}

func TestConfig_Build_RenameWithoutValue(t *testing.T) {
	c := Config{Ops: []string{"rename:field"}}
	_, err := c.Build()
	if err == nil {
		t.Error("expected error: rename requires a value")
	}
}
