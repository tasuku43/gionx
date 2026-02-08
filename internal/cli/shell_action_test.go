package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEmitShellActionCD_WritesActionFile(t *testing.T) {
	actionFile := filepath.Join(t.TempDir(), "action.sh")
	t.Setenv(shellActionFileEnv, actionFile)

	if err := emitShellActionCD("/tmp/example dir"); err != nil {
		t.Fatalf("emitShellActionCD() error: %v", err)
	}

	b, err := os.ReadFile(actionFile)
	if err != nil {
		t.Fatalf("ReadFile(action file) error: %v", err)
	}
	got := string(b)
	if !strings.HasPrefix(got, "cd ") {
		t.Fatalf("action should start with cd: %q", got)
	}
	if !strings.Contains(got, "/tmp/example dir") {
		t.Fatalf("action should contain target path: %q", got)
	}
	if !strings.HasSuffix(got, "\n") {
		t.Fatalf("action should end with newline: %q", got)
	}
}

func TestEmitShellActionCD_NoEnv_NoWrite(t *testing.T) {
	t.Setenv(shellActionFileEnv, "")

	if err := emitShellActionCD("/tmp/ignored"); err != nil {
		t.Fatalf("emitShellActionCD() error: %v", err)
	}
}
