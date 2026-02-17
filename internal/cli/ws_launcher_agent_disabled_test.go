//go:build !experimental

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_WS_Launcher_RunAgent_UnsupportedWithoutExperimental(t *testing.T) {
	root := prepareCurrentRootForTest(t)
	if err := os.MkdirAll(filepath.Join(root, "workspaces", "WS1"), 0o755); err != nil {
		t.Fatalf("mkdir workspace: %v", err)
	}

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"ws", "--id", "WS1", "--act", "run-agent"})
	if code != exitUsage {
		t.Fatalf("exit code = %d, want %d", code, exitUsage)
	}
	if !strings.Contains(err.String(), `unsupported --act: "run-agent"`) {
		t.Fatalf("stderr should include unsupported action error: %q", err.String())
	}
}
