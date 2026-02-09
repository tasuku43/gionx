package statestore

import (
	"context"
	"database/sql"
)

type WorkspaceListItem struct {
	ID        string
	Status    string
	UpdatedAt int64
	Title     string
	SourceURL string
}

type WorkspaceRepo struct {
	RepoUID   string
	RepoKey   string
	Alias     string
	Branch    string
	BaseRef   string
	MissingAt sql.NullInt64
}

func ListWorkspaces(ctx context.Context, db *sql.DB) ([]WorkspaceListItem, error) {
	_ = ctx
	_ = db
	return nil, retiredStateStoreError()
}

func ListWorkspaceRepos(ctx context.Context, db *sql.DB, workspaceID string) ([]WorkspaceRepo, error) {
	_ = ctx
	_ = db
	_ = workspaceID
	return nil, retiredStateStoreError()
}

func MarkWorkspaceRepoMissing(ctx context.Context, db *sql.DB, workspaceID string, repoUID string, now int64) (bool, error) {
	_ = ctx
	_ = db
	_ = workspaceID
	_ = repoUID
	_ = now
	return false, retiredStateStoreError()
}
