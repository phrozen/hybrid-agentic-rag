package logger

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/pterm/pterm"
)

func TestMain(m *testing.M) {
	// Deterministic, ANSI-free output for substring assertions.
	pterm.DisableStyling()
	os.Exit(m.Run())
}

func TestDefaultLogger_PackageFuncs(t *testing.T) {
	var buf bytes.Buffer
	old := Default()
	SetDefault(New(Config{Writer: &buf}))
	t.Cleanup(func() { SetDefault(old) })

	Info("via package func", "tool", "search")

	out := buf.String()
	t.Logf("rendered:\n%s", out)
	if !strings.Contains(out, "via package func") || !strings.Contains(out, "tool") {
		t.Errorf("package-level Info did not route to default logger:\n%s", out)
	}
}

func TestNew_LevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	l := New(Config{Writer: &buf, Level: "warn"})

	l.Info("should be filtered out")
	l.Warn("should appear")

	out := buf.String()
	if strings.Contains(out, "filtered out") {
		t.Errorf("info record leaked at warn level:\n%s", out)
	}
	if !strings.Contains(out, "should appear") {
		t.Errorf("warn record missing:\n%s", out)
	}
}
