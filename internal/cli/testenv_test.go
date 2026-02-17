package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tasuku43/kra/internal/paths"
)

func setKraHomeForTest(t *testing.T) string {
	t.Helper()
	base := t.TempDir()
	home := filepath.Join(base, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("mkdir test HOME dir: %v", err)
	}
	t.Setenv("HOME", home)
	kraHome := filepath.Join(base, ".kra")
	t.Setenv("KRA_HOME", kraHome)
	return kraHome
}

func prepareCurrentRootForTest(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	setKraHomeForTest(t)
	if err := os.MkdirAll(filepath.Join(root, "workspaces"), 0o755); err != nil {
		t.Fatalf("create workspaces/: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "archive"), 0o755); err != nil {
		t.Fatalf("create archive/: %v", err)
	}
	if err := paths.WriteCurrentContext(root); err != nil {
		t.Fatalf("WriteCurrentContext() error: %v", err)
	}
	runGitCmd := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v failed: %v (output=%s)", args, err, string(out))
		}
	}
	runGitCmd("init", "-b", "main")
	runGitCmd("config", "user.email", "test@example.com")
	runGitCmd("config", "user.name", "test")
	return root
}
