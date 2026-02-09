package archguard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectorReuseAcrossCommands(t *testing.T) {
	files := []string{
		filepath.Join("..", "cli", "context.go"),
		filepath.Join("..", "cli", "repo_discover.go"),
		filepath.Join("..", "cli", "repo_remove.go"),
		filepath.Join("..", "cli", "repo_gc.go"),
		filepath.Join("..", "cli", "ws_launcher.go"),
	}

	for _, file := range files {
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		src := string(b)
		if strings.Contains(src, "tea.NewProgram(") {
			t.Fatalf("selector runtime must be reused via shared prompt, found direct tea.NewProgram in %s", file)
		}
		if !strings.Contains(src, "promptWorkspaceSelector") {
			t.Fatalf("selector-capable command should use shared selector prompt in %s", file)
		}
	}
}
