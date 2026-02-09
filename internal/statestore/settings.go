package statestore

import (
	"context"
	"database/sql"
)

// EnsureSettings is retained as a compatibility shim during SQLite retirement.
func EnsureSettings(ctx context.Context, db *sql.DB, rootPath string, repoPoolPath string) error {
	_ = ctx
	_ = db
	_ = rootPath
	_ = repoPoolPath
	return retiredStateStoreError()
}
