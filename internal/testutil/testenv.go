package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tasuku43/kra/internal/paths"
)

type Env struct {
	Root    string
	KraHome string
}

func NewEnv(t *testing.T) Env {
	t.Helper()

	root := t.TempDir()
	base := t.TempDir()
	home := filepath.Join(base, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", home, err)
	}
	kraHome := filepath.Join(base, ".kra")

	t.Setenv("HOME", home)
	t.Setenv("KRA_HOME", kraHome)
	if err := paths.WriteCurrentContext(root); err != nil {
		t.Fatalf("WriteCurrentContext(%q): %v", root, err)
	}

	return Env{
		Root:    root,
		KraHome: kraHome,
	}
}

func (e Env) RepoPoolPath() string {
	return filepath.Join(e.KraHome, "repo-pool")
}

func (e Env) EnsureRootLayout(t *testing.T) {
	t.Helper()
	mustMkdirAll(t, filepath.Join(e.Root, "workspaces"))
	mustMkdirAll(t, filepath.Join(e.Root, "archive"))
	mustMkdirAll(t, filepath.Join(e.Root, "templates", "default", "notes"))
	mustMkdirAll(t, filepath.Join(e.Root, "templates", "default", "artifacts"))
	agentsPath := filepath.Join(e.Root, "templates", "default", "AGENTS.md")
	if err := os.WriteFile(agentsPath, []byte("# test template\n"), 0o644); err != nil {
		t.Fatalf("write %q: %v", agentsPath, err)
	}
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}
