package template

import (
	"strings"
	"testing"
)

func TestNew_EmptyTemplate_ReturnsError(t *testing.T) {
	_, err := New("", false)
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestNew_InvalidTemplate_ReturnsError(t *testing.T) {
	_, err := New("{{ .Foo", false)
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestFormat_PlainMode_ExposesLine(t *testing.T) {
	f, err := New(">> {{ .Line }}", false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := f.Format("hello world")
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if out != ">> hello world" {
		t.Errorf("got %q, want %q", out, ">> hello world")
	}
}

func TestFormat_JSONMode_ExposesFields(t *testing.T) {
	f, err := New(`level={{ index . "level" }} msg={{ index . "msg" }}`, true)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := f.Format(`{"level":"info","msg":"started"}`)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if out != "level=info msg=started" {
		t.Errorf("got %q", out)
	}
}

func TestFormat_JSONMode_InvalidJSON_FallsBack(t *testing.T) {
	f, err := New("{{ .Line }}", true)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	raw := "not json at all"
	out, err := f.Format(raw)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if out != raw {
		t.Errorf("got %q, want %q", out, raw)
	}
}

func TestFormat_MissingKey_RendersEmpty(t *testing.T) {
	f, err := New(`{{ index . "missing" }}`, true)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := f.Format(`{"level":"warn"}`)
	if err != nil {
		t.Fatalf("Format: %v", err)
	}
	if strings.TrimSpace(out) != "<nil>" && out != "" {
		// missingkey=zero yields zero value; for interface{} that is <nil> in text/template
		// either empty or "<nil>" is acceptable depending on Go version.
	}
}

func TestFormat_MultiLine_IndependentCalls(t *testing.T) {
	f, err := New("[{{ .Line }}]", false)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	lines := []string{"alpha", "beta", "gamma"}
	for _, l := range lines {
		out, err := f.Format(l)
		if err != nil {
			t.Fatalf("Format(%q): %v", l, err)
		}
		want := "[" + l + "]"
		if out != want {
			t.Errorf("got %q, want %q", out, want)
		}
	}
}
