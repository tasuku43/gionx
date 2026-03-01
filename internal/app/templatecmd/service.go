package templatecmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/tasuku43/kra/internal/infra/gitutil"
)

func CommitCreate(ctx context.Context, root string, templateName string) (string, error) {
	return commitTemplateChange(ctx, root, templateName, "template-create")
}

func CommitRemove(ctx context.Context, root string, templateName string) (string, error) {
	return commitTemplateChange(ctx, root, templateName, "template-remove")
}

func commitTemplateChange(ctx context.Context, root string, templateName string, commitPrefix string) (string, error) {
	if err := ensureRootGitWorktree(ctx, root); err != nil {
		return "", err
	}

	templateArg := filepath.ToSlash(filepath.Join("templates", templateName))
	templatePrefix, err := toGitTopLevelPath(ctx, root, filepath.Join("templates", templateName))
	if err != nil {
		return "", err
	}
	templatePrefix += string(filepath.Separator)
	resetArgs := []string{templateArg}

	if _, err := gitutil.Run(ctx, root, "add", "-A", "--", templateArg); err != nil {
		if strings.Contains(err.Error(), "did not match any files") || strings.Contains(err.Error(), "did not match any file") {
			return "", nil
		}
		resetTemplateStaging(ctx, root, resetArgs)
		return "", err
	}
	out, err := gitutil.Run(ctx, root, "diff", "--cached", "--name-only")
	if err != nil {
		resetTemplateStaging(ctx, root, resetArgs)
		return "", err
	}
	staged := 0
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		p := strings.TrimSpace(line)
		if p == "" {
			continue
		}
		staged++
		if !strings.HasPrefix(filepath.FromSlash(p), templatePrefix) {
			resetTemplateStaging(ctx, root, resetArgs)
			return "", fmt.Errorf("unexpected staged path outside allowlist: %s", p)
		}
	}
	if staged == 0 {
		return "", nil
	}
	if _, err := gitutil.Run(ctx, root, "commit", "--only", "-m", fmt.Sprintf("%s: %s", commitPrefix, templateName), "--", templateArg); err != nil {
		resetTemplateStaging(ctx, root, resetArgs)
		return "", err
	}
	sha, err := gitutil.Run(ctx, root, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(sha), nil
}

func resetTemplateStaging(ctx context.Context, root string, args []string) {
	for _, arg := range args {
		if strings.TrimSpace(arg) == "" {
			continue
		}
		_, _ = gitutil.Run(ctx, root, "reset", "-q", "--", arg)
	}
}

func ensureRootGitWorktree(ctx context.Context, root string) error {
	out, err := gitutil.Run(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return fmt.Errorf("KRA_ROOT must be a git working tree: %w", err)
	}

	got := filepath.Clean(strings.TrimSpace(out))
	want := filepath.Clean(root)

	if gotEval, err := filepath.EvalSymlinks(got); err == nil {
		got = gotEval
	}
	if wantEval, err := filepath.EvalSymlinks(want); err == nil {
		want = wantEval
	}

	rel, relErr := filepath.Rel(got, want)
	if relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("KRA_ROOT must be inside the git worktree: toplevel=%s root=%s", strings.TrimSpace(out), root)
	}
	return nil
}

// toGitTopLevelPath converts a path relative to KRA_ROOT into a path
// relative to the enclosing git toplevel (the path domain of `git diff --name-only`).
func toGitTopLevelPath(ctx context.Context, root string, rootRelativePath string) (string, error) {
	topRaw, err := gitutil.Run(ctx, root, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
	}
	top := filepath.Clean(strings.TrimSpace(topRaw))
	root = filepath.Clean(root)
	if resolved, evalErr := filepath.EvalSymlinks(top); evalErr == nil {
		top = filepath.Clean(resolved)
	}
	if resolved, evalErr := filepath.EvalSymlinks(root); evalErr == nil {
		root = filepath.Clean(resolved)
	}

	relRoot := "."
	if r, relErr := filepath.Rel(top, root); relErr == nil {
		r = filepath.Clean(r)
		if r == ".." || strings.HasPrefix(r, ".."+string(filepath.Separator)) {
			relRoot = "."
		} else {
			relRoot = r
		}
	}

	p := filepath.Clean(filepath.FromSlash(rootRelativePath))
	if relRoot == "." {
		return p, nil
	}
	return filepath.Clean(filepath.Join(relRoot, p)), nil
}
