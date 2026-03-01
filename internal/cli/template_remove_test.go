package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tasuku43/kra/internal/testutil"
)

func TestCLI_TemplateRemove_NameFlag_RemovesTemplate(t *testing.T) {
	env := testutil.NewEnv(t)
	env.EnsureRootLayout(t)
	target := filepath.Join(env.Root, "templates", "custom")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir custom template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write custom file: %v", err)
	}
	if err := gitCommitPaths(env.Root, "seed template custom", "templates/custom"); err != nil {
		t.Fatalf("seed template commit: %v", err)
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	c := New(&out, &errBuf)

	code := c.Run([]string{"template", "remove", "--name", "custom"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q)", code, exitOK, errBuf.String())
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("template should be removed, stat err=%v", err)
	}
	if !strings.Contains(out.String(), "Removed 1 / 1") || !strings.Contains(out.String(), "✔ custom") {
		t.Fatalf("stdout missing remove summary: %q", out.String())
	}
	if !strings.Contains(out.String(), "commit: ") {
		t.Fatalf("stdout missing commit line: %q", out.String())
	}
	msg, err := gitLogHeadSubject(env.Root)
	if err != nil {
		t.Fatalf("git log head subject: %v", err)
	}
	if msg != "template-remove: custom" {
		t.Fatalf("commit subject=%q, want %q", msg, "template-remove: custom")
	}
}

func TestCLI_TemplateRemove_RMAlias_Works(t *testing.T) {
	env := testutil.NewEnv(t)
	env.EnsureRootLayout(t)
	target := filepath.Join(env.Root, "templates", "custom")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir custom template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write custom file: %v", err)
	}
	if err := gitCommitPaths(env.Root, "seed template custom", "templates/custom"); err != nil {
		t.Fatalf("seed template commit: %v", err)
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	c := New(&out, &errBuf)
	code := c.Run([]string{"template", "rm", "custom"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q)", code, exitOK, errBuf.String())
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("template should be removed, stat err=%v", err)
	}
}

func TestCLI_TemplateRemove_PromptWhenNameOmitted(t *testing.T) {
	env := testutil.NewEnv(t)
	env.EnsureRootLayout(t)
	target := filepath.Join(env.Root, "templates", "custom")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir custom template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write custom file: %v", err)
	}
	if err := gitCommitPaths(env.Root, "seed template custom", "templates/custom"); err != nil {
		t.Fatalf("seed template commit: %v", err)
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	c := New(&out, &errBuf)
	c.In = strings.NewReader("custom\n")
	code := c.Run([]string{"template", "remove"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q)", code, exitOK, errBuf.String())
	}
	if !strings.Contains(errBuf.String(), "Inputs:") || !strings.Contains(errBuf.String(), "name: custom") {
		t.Fatalf("stderr missing inputs echo: %q", errBuf.String())
	}
}

func TestCLI_TemplateRemove_Missing_ReturnsStructuredError(t *testing.T) {
	env := testutil.NewEnv(t)
	env.EnsureRootLayout(t)

	var out bytes.Buffer
	var errBuf bytes.Buffer
	c := New(&out, &errBuf)
	code := c.Run([]string{"template", "remove", "--name", "missing"})
	if code != exitError {
		t.Fatalf("exit code = %d, want %d", code, exitError)
	}
	if !strings.Contains(errBuf.String(), "Error:") || !strings.Contains(errBuf.String(), "reason:") || !strings.Contains(errBuf.String(), `template "missing" not found`) {
		t.Fatalf("stderr missing structured error: %q", errBuf.String())
	}
}

func TestCLI_TemplateRemove_AbortWhenStagedOutsideAllowlist(t *testing.T) {
	env := testutil.NewEnv(t)
	env.EnsureRootLayout(t)
	target := filepath.Join(env.Root, "templates", "custom")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir custom template: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("x\n"), 0o644); err != nil {
		t.Fatalf("write custom file: %v", err)
	}
	if err := gitCommitPaths(env.Root, "seed template custom", "templates/custom"); err != nil {
		t.Fatalf("seed template commit: %v", err)
	}
	outside := filepath.Join(env.Root, "README.md")
	if err := os.WriteFile(outside, []byte("outside\n"), 0o644); err != nil {
		t.Fatalf("write outside file: %v", err)
	}
	cmd := exec.Command("git", "add", "--", "README.md")
	cmd.Dir = env.Root
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git add README.md: %v (out=%s)", err, string(b))
	}

	var out bytes.Buffer
	var errBuf bytes.Buffer
	c := New(&out, &errBuf)
	code := c.Run([]string{"template", "remove", "--name", "custom"})
	if code != exitError {
		t.Fatalf("exit code = %d, want %d", code, exitError)
	}
	if !strings.Contains(errBuf.String(), "outside allowlist") {
		t.Fatalf("stderr missing allowlist error: %q", errBuf.String())
	}
}

func gitCommitPaths(root string, message string, pathsToAdd ...string) error {
	addArgs := append([]string{"add", "--"}, pathsToAdd...)
	add := exec.Command("git", addArgs...)
	add.Dir = root
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git %v: %v (%s)", addArgs, err, string(out))
	}
	commit := exec.Command("git", "commit", "-m", message)
	commit.Dir = root
	if out, err := commit.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit: %v (%s)", err, string(out))
	}
	return nil
}
