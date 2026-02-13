package cli

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/tasuku43/kra/internal/core/repospec"
	"github.com/tasuku43/kra/internal/core/repostore"
	"github.com/tasuku43/kra/internal/gitutil"
	"github.com/tasuku43/kra/internal/testutil"
)

func seedRepoPoolAndState(t *testing.T, env testutil.Env, repoSpecInput string) (repoUID string, repoKey string, alias string) {
	t.Helper()
	ctx := context.Background()

	spec, err := repospec.Normalize(repoSpecInput)
	if err != nil {
		t.Fatalf("Normalize(repoSpec): %v", err)
	}
	repoKey = fmt.Sprintf("%s/%s", spec.Owner, spec.Repo)
	repoUID = fmt.Sprintf("%s/%s", spec.Host, repoKey)
	alias = spec.Repo

	defaultBranch, err := gitutil.DefaultBranchFromRemote(ctx, repoSpecInput)
	if err != nil {
		t.Fatalf("DefaultBranchFromRemote() error: %v", err)
	}
	barePath := repostore.StorePath(env.RepoPoolPath(), spec)
	if _, err := gitutil.EnsureBareRepoFetched(ctx, repoSpecInput, barePath, defaultBranch); err != nil {
		t.Fatalf("EnsureBareRepoFetched() error: %v", err)
	}
	return repoUID, repoKey, alias
}

func addRepoSelectionInput(baseRef string, branch string) string {
	return fmt.Sprintf("1\n%s\n%s\n\n", strings.TrimSpace(baseRef), strings.TrimSpace(branch))
}
