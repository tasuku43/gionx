package statestore

import (
	"context"
	"database/sql"
	"fmt"
)

type RepoConflictError struct {
	RepoUID   string
	RepoKey   string
	RemoteURL string
}

func (e *RepoConflictError) Error() string {
	return fmt.Sprintf("repo already exists with different metadata (repo_uid=%q repo_key=%q remote_url=%q)", e.RepoUID, e.RepoKey, e.RemoteURL)
}

type WorkspaceRepoAliasConflictError struct {
	WorkspaceID string
	Alias       string
}

func (e *WorkspaceRepoAliasConflictError) Error() string {
	return fmt.Sprintf("repo alias already exists in workspace (workspace_id=%q alias=%q)", e.WorkspaceID, e.Alias)
}

type WorkspaceRepoAlreadyBoundError struct {
	WorkspaceID string
	RepoUID     string
}

func (e *WorkspaceRepoAlreadyBoundError) Error() string {
	return fmt.Sprintf("repo already bound to workspace (workspace_id=%q repo_uid=%q)", e.WorkspaceID, e.RepoUID)
}

type EnsureRepoInput struct {
	RepoUID   string
	RepoKey   string
	RemoteURL string
	Now       int64
}

type AddWorkspaceRepoInput struct {
	WorkspaceID   string
	RepoUID       string
	RepoKey       string
	Alias         string
	Branch        string
	BaseRef       string
	RepoSpecInput string
	Now           int64
}

type RootRepoCandidate struct {
	RepoUID           string
	RepoKey           string
	RemoteURL         string
	WorkspaceRefCount int
	Score             float64
}

type RepoPoolCandidate struct {
	RepoUID   string
	RepoKey   string
	RemoteURL string
	Score     float64
}

func WorkspaceRepoAliasExists(ctx context.Context, db *sql.DB, workspaceID string, alias string) (bool, error) {
	_ = ctx
	_ = db
	_ = workspaceID
	_ = alias
	return false, retiredStateStoreError()
}

func EnsureRepo(ctx context.Context, db *sql.DB, in EnsureRepoInput) error {
	_ = ctx
	_ = db
	_ = in
	return retiredStateStoreError()
}

func LookupRepoRemoteURL(ctx context.Context, db *sql.DB, repoUID string) (string, bool, error) {
	_ = ctx
	_ = db
	_ = repoUID
	return "", false, retiredStateStoreError()
}

func AddWorkspaceRepo(ctx context.Context, db *sql.DB, in AddWorkspaceRepoInput) error {
	_ = ctx
	_ = db
	_ = in
	return retiredStateStoreError()
}

func DeleteWorkspaceRepoBinding(ctx context.Context, db *sql.DB, workspaceID string, repoUID string) error {
	_ = ctx
	_ = db
	_ = workspaceID
	_ = repoUID
	return retiredStateStoreError()
}

func TouchRepoUpdatedAt(ctx context.Context, db *sql.DB, repoUID string, now int64) error {
	_ = ctx
	_ = db
	_ = repoUID
	_ = now
	return retiredStateStoreError()
}

func ListRepoPoolCandidates(ctx context.Context, db *sql.DB, startDay int) ([]RepoPoolCandidate, error) {
	_ = ctx
	_ = db
	_ = startDay
	return nil, retiredStateStoreError()
}

func IncrementRepoUsageDaily(ctx context.Context, db *sql.DB, repoUID string, day int, now int64) error {
	_ = ctx
	_ = db
	_ = repoUID
	_ = day
	_ = now
	return retiredStateStoreError()
}

func ListRepoUIDs(ctx context.Context, db *sql.DB) ([]string, error) {
	_ = ctx
	_ = db
	return nil, retiredStateStoreError()
}

func ListRootRepoCandidates(ctx context.Context, db *sql.DB, startDay int) ([]RootRepoCandidate, error) {
	_ = ctx
	_ = db
	_ = startDay
	return nil, retiredStateStoreError()
}

func DeleteReposByUIDs(ctx context.Context, db *sql.DB, repoUIDs []string) error {
	_ = ctx
	_ = db
	_ = repoUIDs
	return retiredStateStoreError()
}
