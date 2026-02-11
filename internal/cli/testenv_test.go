package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tasuku43/gionx/internal/paths"
)

func setGionxHomeForTest(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	home := filepath.Join(base, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir test HOME dir: %v", err)
	}
	t.Setenv("HOME", home)
	gionxHome := filepath.Join(base, ".gionx")
	t.Setenv("GIONX_HOME", gionxHome)
	return gionxHome
}

func prepareCurrentRootForTest(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	setGionxHomeForTest(t)
	if err := os.MkdirAll(filepath.Join(root, "workspaces"), 0o755); err != nil {
		t.Fatalf("create workspaces/: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "archive"), 0o755); err != nil {
		t.Fatalf("create archive/: %v", err)
	}
	if err := paths.WriteCurrentContext(root); err != nil {
		t.Fatalf("WriteCurrentContext() error: %v", err)
	}
	return root
}
