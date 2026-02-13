package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCLI_Shell_Help(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)

	code := c.Run([]string{"shell", "--help"})
	if code != exitOK {
		t.Fatalf("exit code=%d, want=%d", code, exitOK)
	}
	if !strings.Contains(out.String(), "kra shell") {
		t.Fatalf("stdout missing usage: %q", out.String())
	}
	if err.Len() != 0 {
		t.Fatalf("stderr not empty: %q", err.String())
	}
}

func TestCLI_Shell_Init_Zsh_PrintsEvalReadyFunction(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)

	code := c.Run([]string{"shell", "init", "zsh"})
	if code != exitOK {
		t.Fatalf("exit code=%d, want=%d (stderr=%q)", code, exitOK, err.String())
	}
	text := out.String()
	if !strings.Contains(text, `eval "$(kra shell init zsh)"`) {
		t.Fatalf("missing bootstrap hint: %q", text)
	}
	if !strings.Contains(text, "kra() {") {
		t.Fatalf("missing function definition: %q", text)
	}
	if !strings.Contains(text, `KRA_SHELL_ACTION_FILE="$__kra_action_file"`) {
		t.Fatalf("missing action file env wiring: %q", text)
	}
	if !strings.Contains(text, `command kra "$@"`) {
		t.Fatalf("missing unified command passthrough: %q", text)
	}
	if !strings.Contains(text, `if [ -s "$__kra_action_file" ]; then`) {
		t.Fatalf("missing action file apply block: %q", text)
	}
}

func TestRenderShellInitScript_UnsupportedShell(t *testing.T) {
	_, err := renderShellInitScript("nushell")
	if err == nil {
		t.Fatal("expected unsupported shell error")
	}
}
