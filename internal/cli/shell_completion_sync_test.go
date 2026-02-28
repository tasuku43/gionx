package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	krapaths "github.com/tasuku43/kra/internal/paths"
)

var completionHelpFlagPattern = regexp.MustCompile(`--[a-z0-9][a-z0-9-]*|-h`)

func TestCompletionFlags_AreDocumentedInHelp(t *testing.T) {
	t.Setenv(experimentsEnvKey, experimentInsightCapture)
	prepareCompletionHelpRoot(t)

	checkPath := func(path string, flags []string) {
		t.Helper()

		needsValidation := false
		for _, flag := range flags {
			if flag != "--help" && flag != "-h" {
				needsValidation = true
				break
			}
		}
		if !needsValidation {
			return
		}

		helpText := runHelpTextForCompletionPath(t, path)
		helpFlags := extractHelpFlags(helpText)
		for _, flag := range flags {
			if flag == "--help" || flag == "-h" {
				continue
			}
			if _, ok := helpFlags[flag]; ok {
				continue
			}
			t.Fatalf("completion flag %q for %q is not documented in help output", flag, path)
		}
	}

	for _, cmd := range kraCompletionCommandFlagOrder {
		checkPath(cmd, kraCompletionCommandFlags[cmd])
	}
	for _, path := range kraCompletionPathFlagOrder {
		checkPath(path, kraCompletionPathFlags[path])
	}
}

func prepareCompletionHelpRoot(t *testing.T) {
	t.Helper()

	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "workspaces"), 0o755); err != nil {
		t.Fatalf("create workspaces dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "archive"), 0o755); err != nil {
		t.Fatalf("create archive dir: %v", err)
	}

	kraHome := filepath.Join(t.TempDir(), ".kra")
	t.Setenv("KRA_HOME", kraHome)
	if err := krapaths.WriteCurrentContext(root); err != nil {
		t.Fatalf("write current context: %v", err)
	}

	prevWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir root: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(prevWD)
	})
}

func runHelpTextForCompletionPath(t *testing.T, path string) string {
	t.Helper()

	args := append(strings.Fields(path), "--help")
	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)

	code := c.Run(args)
	if code != exitOK {
		t.Fatalf("help for %q failed: code=%d stderr=%q", path, code, err.String())
	}
	return out.String()
}

func extractHelpFlags(helpText string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, token := range completionHelpFlagPattern.FindAllString(helpText, -1) {
		out[token] = struct{}{}
	}
	return out
}
