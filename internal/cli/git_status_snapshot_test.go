package cli

import (
	"strings"
	"testing"
)

func TestParseGitRepoSnapshot_CollectsBranchAndFileSummary(t *testing.T) {
	raw := strings.Join([]string{
		"# branch.oid 8f1a55fe322bc4f1086dcf6ef9f3a0f755fdd124",
		"# branch.head feature/cache",
		"# branch.upstream origin/feature/cache",
		"# branch.ab +3 -1",
		"1 M. N... 100644 100644 100644 aaaaaaa bbbbbbb README.md",
		"1 .M N... 100644 100644 100644 aaaaaaa bbbbbbb docs/spec.md",
		"? tmp.txt",
		"2 R. N... 100644 100644 100644 aaaaaaa bbbbbbb R100 new/name.md\told/name.md",
	}, "\n")

	snapshot, err := parseGitRepoSnapshot(raw)
	if err != nil {
		t.Fatalf("parseGitRepoSnapshot() error = %v, want nil", err)
	}

	if !snapshot.Status.Dirty {
		t.Fatalf("Dirty = false, want true")
	}
	if snapshot.Branch != "feature/cache" {
		t.Fatalf("Branch = %q, want %q", snapshot.Branch, "feature/cache")
	}
	if snapshot.Status.Upstream != "origin/feature/cache" {
		t.Fatalf("Upstream = %q, want %q", snapshot.Status.Upstream, "origin/feature/cache")
	}
	if snapshot.Status.AheadCount != 3 || snapshot.Status.BehindCount != 1 {
		t.Fatalf("ahead/behind = %d/%d, want 3/1", snapshot.Status.AheadCount, snapshot.Status.BehindCount)
	}
	if snapshot.Staged != 2 || snapshot.Unstaged != 1 || snapshot.Untracked != 1 {
		t.Fatalf("staged/unstaged/untracked = %d/%d/%d, want 2/1/1", snapshot.Staged, snapshot.Unstaged, snapshot.Untracked)
	}
	if len(snapshot.Files) != 4 {
		t.Fatalf("files length = %d, want 4", len(snapshot.Files))
	}
	if !strings.Contains(snapshot.Files[3], "old/name.md -> new/name.md") {
		t.Fatalf("rename file format mismatch: %q", snapshot.Files[3])
	}
}

func TestParseGitRepoSnapshot_InitialDetachedBranch(t *testing.T) {
	raw := strings.Join([]string{
		"# branch.oid (initial)",
		"# branch.head (detached)",
		"? new.txt",
	}, "\n")

	snapshot, err := parseGitRepoSnapshot(raw)
	if err != nil {
		t.Fatalf("parseGitRepoSnapshot() error = %v, want nil", err)
	}
	if !snapshot.Status.HeadMissing {
		t.Fatalf("HeadMissing = false, want true")
	}
	if !snapshot.Status.Detached {
		t.Fatalf("Detached = false, want true")
	}
	if snapshot.Branch != "" {
		t.Fatalf("Branch = %q, want empty", snapshot.Branch)
	}
	if !snapshot.Status.Dirty {
		t.Fatalf("Dirty = false, want true")
	}
	if snapshot.Untracked != 1 {
		t.Fatalf("Untracked = %d, want 1", snapshot.Untracked)
	}
}

func TestParseGitRepoSnapshot_InvalidBranchAB(t *testing.T) {
	raw := strings.Join([]string{
		"# branch.oid 8f1a55fe322bc4f1086dcf6ef9f3a0f755fdd124",
		"# branch.head feature/cache",
		"# branch.ab +x -1",
	}, "\n")

	if _, err := parseGitRepoSnapshot(raw); err == nil {
		t.Fatalf("parseGitRepoSnapshot() error = nil, want non-nil")
	}
}
