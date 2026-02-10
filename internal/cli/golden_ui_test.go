package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tasuku43/gion-core/workspacerisk"
)

func TestGolden_WSActionSelectorSingle(t *testing.T) {
	lines := renderWorkspaceSelectorLinesWithOptions(
		"active",
		"Action:",
		"run",
		[]workspaceSelectorCandidate{
			{ID: "add-repo", Description: "add repositories", Risk: workspacerisk.WorkspaceRiskClean},
			{ID: "close", Description: "archive this workspace", Risk: workspacerisk.WorkspaceRiskClean},
		},
		map[int]bool{},
		0,
		"",
		selectorMessageLevelMuted,
		"",
		true,
		true,
		true,
		false,
		false,
		120,
	)
	assertGolden(t, "ws_action_selector_single.golden", strings.Join(lines, "\n")+"\n")
}

func TestGolden_RepoPoolSelectorMulti(t *testing.T) {
	lines := renderWorkspaceSelectorLinesWithOptions(
		"active",
		"Repos(pool):",
		"add",
		[]workspaceSelectorCandidate{
			{ID: "example-org/helmfiles", Risk: workspacerisk.WorkspaceRiskClean},
			{ID: "example-org/sre-apps", Risk: workspacerisk.WorkspaceRiskClean},
		},
		map[int]bool{1: true},
		0,
		"",
		selectorMessageLevelMuted,
		"",
		false,
		true,
		false,
		false,
		false,
		120,
	)
	assertGolden(t, "repo_pool_selector_multi.golden", strings.Join(lines, "\n")+"\n")
}

func TestGolden_WSAddRepoPlan(t *testing.T) {
	var out bytes.Buffer
	plan := []addRepoPlanItem{
		{Candidate: addRepoPoolCandidate{RepoKey: "example-org/terraforms"}},
		{Candidate: addRepoPoolCandidate{RepoKey: "tasuku43/gionx"}},
	}
	printAddRepoPlan(&out, "DEMO-0000", plan, false)
	assertGolden(t, "ws_add_repo_plan.golden", out.String())
}

func assertGolden(t *testing.T, name string, got string) {
	t.Helper()
	path := filepath.Join("testdata", "golden", name)
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden file: %v", err)
		}
	}
	wantBytes, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden file %s: %v", path, err)
	}
	if got != string(wantBytes) {
		t.Fatalf("golden mismatch for %s\n--- want ---\n%s\n--- got ---\n%s", name, string(wantBytes), got)
	}
}
