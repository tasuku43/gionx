package statestore

import (
	"context"
	"database/sql"
	"fmt"
)

type WorkspaceAlreadyExistsError struct {
	ID     string
	Status string
}

func (e *WorkspaceAlreadyExistsError) Error() string {
	return fmt.Sprintf("workspace already exists (id=%q status=%q)", e.ID, e.Status)
}

type CreateWorkspaceInput struct {
	ID        string
	Title     string
	SourceURL string
	Now       int64
}

type ArchiveWorkspaceInput struct {
	ID                string
	ArchivedCommitSHA string
	Now               int64
}

type ReopenWorkspaceInput struct {
	ID                string
	ReopenedCommitSHA string
	Now               int64
}

type PurgeWorkspaceInput struct {
	ID  string
	Now int64
}

func LookupWorkspaceStatus(ctx context.Context, db *sql.DB, id string) (string, bool, error) {
	_ = ctx
	_ = db
	_ = id
	return "", false, retiredStateStoreError()
}

func CreateWorkspace(ctx context.Context, db *sql.DB, in CreateWorkspaceInput) (int, error) {
	_ = ctx
	_ = db
	_ = in
	return 0, retiredStateStoreError()
}

func ArchiveWorkspace(ctx context.Context, db *sql.DB, in ArchiveWorkspaceInput) error {
	_ = ctx
	_ = db
	_ = in
	return retiredStateStoreError()
}

func ReopenWorkspace(ctx context.Context, db *sql.DB, in ReopenWorkspaceInput) error {
	_ = ctx
	_ = db
	_ = in
	return retiredStateStoreError()
}

func PurgeWorkspace(ctx context.Context, db *sql.DB, in PurgeWorkspaceInput) error {
	_ = ctx
	_ = db
	_ = in
	return retiredStateStoreError()
}
