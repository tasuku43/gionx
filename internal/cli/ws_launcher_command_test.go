package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tasuku43/kra/internal/testutil"
)

func TestCLI_WS_Launcher_RequiresIDOrWorkspaceContext(t *testing.T) {
	env := testutil.NewEnv(t)
	initAndConfigureRootRepo(t, env.Root)

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"ws"})
	if code != exitUsage {
		t.Fatalf("ws exit code = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(err.String(), "ws requires one of --id <id>, --current, or --select") {
		t.Fatalf("stderr missing unresolved launcher message: %q", err.String())
	}
}

func TestCLI_WS_Launcher_WithIDAndFixedAction(t *testing.T) {
	env := testutil.NewEnv(t)
	initAndConfigureRootRepo(t, env.Root)

	{
		var out bytes.Buffer
		var err bytes.Buffer
		c := New(&out, &err)
		code := c.Run([]string{"ws", "create", "--no-prompt", "WS1"})
		if code != exitOK {
			t.Fatalf("ws create exit code = %d, want %d (stderr=%q)", code, exitOK, err.String())
		}
	}

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	actionFile := filepath.Join(t.TempDir(), "action.sh")
	t.Setenv(shellActionFileEnv, actionFile)
	code := c.Run([]string{"ws", "go", "--id", "WS1"})
	if code != exitOK {
		t.Fatalf("ws go --id exit code = %d, want %d (stderr=%q)", code, exitOK, err.String())
	}
	if out.String() != "" {
		t.Fatalf("stdout should be empty for ws go default mode: %q", out.String())
	}
	action, readErr := os.ReadFile(actionFile)
	if readErr != nil {
		t.Fatalf("ReadFile(action) error: %v", readErr)
	}
	want := filepath.Join(env.Root, "workspaces", "WS1")
	if !strings.Contains(string(action), "cd ") || !strings.Contains(string(action), want) {
		t.Fatalf("action file missing destination: %q", string(action))
	}
}
