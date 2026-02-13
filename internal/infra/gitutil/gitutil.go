package gitutil

import "context"

import base "github.com/tasuku43/kra/internal/gitutil"

func EnsureGitInPath() error { return base.EnsureGitInPath() }

func Run(ctx context.Context, dir string, args ...string) (string, error) {
	return base.Run(ctx, dir, args...)
}

func RunBare(ctx context.Context, barePath string, args ...string) (string, error) {
	return base.RunBare(ctx, barePath, args...)
}

func ShowRefExistsBare(ctx context.Context, barePath string, ref string) (bool, error) {
	return base.ShowRefExistsBare(ctx, barePath, ref)
}

func EnsureBareRepoFetched(ctx context.Context, repoSpecInput string, barePath string, defaultBranch string) (string, error) {
	return base.EnsureBareRepoFetched(ctx, repoSpecInput, barePath, defaultBranch)
}

func DefaultBranchFromRemote(ctx context.Context, repoSpecInput string) (string, error) {
	return base.DefaultBranchFromRemote(ctx, repoSpecInput)
}

func CheckRefFormat(ctx context.Context, ref string) error {
	return base.CheckRefFormat(ctx, ref)
}

func IsIgnored(ctx context.Context, dir string, path string) (bool, error) {
	return base.IsIgnored(ctx, dir, path)
}
