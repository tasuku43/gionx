package cli

import (
	"strings"
	"testing"
)

func TestBuildAddRepoInputsLines_BaseOnlySingleRepo(t *testing.T) {
	rows := []addRepoInputProgress{
		{
			RepoKey: "tasuku43/puml-parser-php",
			BaseRef: "origin/main",
		},
	}

	lines := buildAddRepoInputsLines("TEST-010", rows, 0, false)
	got := strings.Join(lines, "\n")
	want := strings.Join([]string{
		"  • repos:",
		"    └─ tasuku43/puml-parser-php",
		"       ├─ base_ref: origin/main",
	}, "\n")

	if !strings.Contains(got, want) {
		t.Fatalf("unexpected inputs block:\n%s", got)
	}
}

func TestBuildAddRepoInputsLines_BaseOnlyNonActiveRepo(t *testing.T) {
	rows := []addRepoInputProgress{
		{
			RepoKey: "tasuku43/puml-parser-php",
			BaseRef: "origin/main",
		},
	}

	lines := buildAddRepoInputsLines("TEST-010", rows, -1, false)
	got := strings.Join(lines, "\n")
	want := strings.Join([]string{
		"  • repos:",
		"    └─ tasuku43/puml-parser-php",
		"       └─ base_ref: origin/main",
	}, "\n")

	if !strings.Contains(got, want) {
		t.Fatalf("unexpected inputs block:\n%s", got)
	}
}

func TestBuildAddRepoInputsLines_BaseAndBranchSingleRepo(t *testing.T) {
	rows := []addRepoInputProgress{
		{
			RepoKey: "tasuku43/puml-parser-php",
			BaseRef: "origin/main",
			Branch:  "dddd",
		},
	}

	lines := buildAddRepoInputsLines("TEST-010", rows, 0, false)
	got := strings.Join(lines, "\n")
	want := strings.Join([]string{
		"  • repos:",
		"    └─ tasuku43/puml-parser-php",
		"       ├─ base_ref: origin/main",
		"       └─ branch: dddd",
	}, "\n")

	if !strings.Contains(got, want) {
		t.Fatalf("unexpected inputs block:\n%s", got)
	}
}

func TestBuildAddRepoInputsLines_FirstRepoFinalizedThenSecondBaseOnly(t *testing.T) {
	rows := []addRepoInputProgress{
		{
			RepoKey: "tasuku43/puml-parser-php",
			BaseRef: "origin/main",
			Branch:  "dddd",
		},
		{
			RepoKey: "tasuku43/dependency-analyzer",
			BaseRef: "origin/main",
		},
	}

	lines := buildAddRepoInputsLines("TEST-010", rows, 1, false)
	got := strings.Join(lines, "\n")
	want := strings.Join([]string{
		"  • repos:",
		"    ├─ tasuku43/puml-parser-php",
		"    │  ├─ base_ref: origin/main",
		"    │  └─ branch: dddd",
		"    └─ tasuku43/dependency-analyzer",
		"       ├─ base_ref: origin/main",
	}, "\n")

	if !strings.Contains(got, want) {
		t.Fatalf("unexpected inputs block:\n%s", got)
	}
}
