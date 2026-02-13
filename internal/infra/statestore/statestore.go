package statestore

import (
	"context"
	"database/sql"

	base "github.com/tasuku43/kra/internal/statestore"
)

type EnsureRepoInput = base.EnsureRepoInput
type RootRepoCandidate = base.RootRepoCandidate
type RepoPoolCandidate = base.RepoPoolCandidate
type WorkspaceListItem = base.WorkspaceListItem
type AddWorkspaceRepoInput = base.AddWorkspaceRepoInput
type ArchiveWorkspaceInput = base.ArchiveWorkspaceInput
type WorkspaceRepo = base.WorkspaceRepo
type CreateWorkspaceInput = base.CreateWorkspaceInput
type WorkspaceAlreadyExistsError = base.WorkspaceAlreadyExistsError
type PurgeWorkspaceInput = base.PurgeWorkspaceInput
type ReopenWorkspaceInput = base.ReopenWorkspaceInput

func Open(ctx context.Context, dbPath string) (*sql.DB, error) {
	return base.Open(ctx, dbPath)
}

func EnsureSettings(ctx context.Context, db *sql.DB, root string, repoPoolPath string) error {
	return base.EnsureSettings(ctx, db, root, repoPoolPath)
}

func ListRepoUIDs(ctx context.Context, db *sql.DB) ([]string, error) {
	return base.ListRepoUIDs(ctx, db)
}

func LookupRepoRemoteURL(ctx context.Context, db *sql.DB, repoUID string) (string, bool, error) {
	return base.LookupRepoRemoteURL(ctx, db, repoUID)
}

func EnsureRepo(ctx context.Context, db *sql.DB, in EnsureRepoInput) error {
	return base.EnsureRepo(ctx, db, in)
}

func DeleteReposByUIDs(ctx context.Context, db *sql.DB, repoUIDs []string) error {
	return base.DeleteReposByUIDs(ctx, db, repoUIDs)
}

func ListRootRepoCandidates(ctx context.Context, db *sql.DB, startDay int) ([]RootRepoCandidate, error) {
	return base.ListRootRepoCandidates(ctx, db, startDay)
}

func AddWorkspaceRepo(ctx context.Context, db *sql.DB, in AddWorkspaceRepoInput) error {
	return base.AddWorkspaceRepo(ctx, db, in)
}

func DeleteWorkspaceRepoBinding(ctx context.Context, db *sql.DB, workspaceID string, repoUID string) error {
	return base.DeleteWorkspaceRepoBinding(ctx, db, workspaceID, repoUID)
}

func IncrementRepoUsageDaily(ctx context.Context, db *sql.DB, repoUID string, day int, nowUnix int64) error {
	return base.IncrementRepoUsageDaily(ctx, db, repoUID, day, nowUnix)
}

func TouchRepoUpdatedAt(ctx context.Context, db *sql.DB, repoUID string, nowUnix int64) error {
	return base.TouchRepoUpdatedAt(ctx, db, repoUID, nowUnix)
}

func ListRepoPoolCandidates(ctx context.Context, db *sql.DB, startDay int) ([]RepoPoolCandidate, error) {
	return base.ListRepoPoolCandidates(ctx, db, startDay)
}

func ListWorkspaceRepos(ctx context.Context, db *sql.DB, workspaceID string) ([]WorkspaceRepo, error) {
	return base.ListWorkspaceRepos(ctx, db, workspaceID)
}

func LookupWorkspaceStatus(ctx context.Context, db *sql.DB, workspaceID string) (string, bool, error) {
	return base.LookupWorkspaceStatus(ctx, db, workspaceID)
}

func ListWorkspaces(ctx context.Context, db *sql.DB) ([]WorkspaceListItem, error) {
	return base.ListWorkspaces(ctx, db)
}

func ArchiveWorkspace(ctx context.Context, db *sql.DB, in ArchiveWorkspaceInput) error {
	return base.ArchiveWorkspace(ctx, db, in)
}

func CreateWorkspace(ctx context.Context, db *sql.DB, in CreateWorkspaceInput) (int, error) {
	return base.CreateWorkspace(ctx, db, in)
}

func MarkWorkspaceRepoMissing(ctx context.Context, db *sql.DB, workspaceID string, repoUID string, nowUnix int64) (bool, error) {
	return base.MarkWorkspaceRepoMissing(ctx, db, workspaceID, repoUID, nowUnix)
}

func PurgeWorkspace(ctx context.Context, db *sql.DB, in PurgeWorkspaceInput) error {
	return base.PurgeWorkspace(ctx, db, in)
}

func ReopenWorkspace(ctx context.Context, db *sql.DB, in ReopenWorkspaceInput) error {
	return base.ReopenWorkspace(ctx, db, in)
}
