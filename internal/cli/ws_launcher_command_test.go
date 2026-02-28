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
	code := c.Run([]string{"ws", "close", "--id", "WS1"})
	if code != exitOK {
		t.Fatalf("ws close --id exit code = %d, want %d (stderr=%q)", code, exitOK, err.String())
	}
	if _, statErr := os.Stat(filepath.Join(env.Root, "archive", "WS1")); statErr != nil {
		t.Fatalf("archive/WS1 should exist after close: %v", statErr)
	}
}
