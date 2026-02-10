package cli

import (
	"context"
	"strings"

	"github.com/tasuku43/gionx/internal/infra/gitutil"
)

func detectOriginHeadBaseRef(ctx context.Context, worktreePath string) string {
	out, err := gitutil.Run(ctx, worktreePath, "symbolic-ref", "--quiet", "refs/remotes/origin/HEAD")
	if err != nil {
		return ""
	}
	ref := strings.TrimSpace(out)
	const pfx = "refs/remotes/origin/"
	if !strings.HasPrefix(ref, pfx) {
		return ""
	}
	branch := strings.TrimSpace(strings.TrimPrefix(ref, pfx))
	if branch == "" {
		return ""
	}
	return "origin/" + branch
}
