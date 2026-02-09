package statestore

import (
	"context"
	"database/sql"
	"fmt"
)

const retiredStateStoreMessage = "state store is retired (sqlite removed)"

func retiredStateStoreError() error {
	return fmt.Errorf(retiredStateStoreMessage)
}

// Open is kept for backward compatibility during the SQLite retirement.
// State store is no longer available and always returns an error.
func Open(ctx context.Context, dbPath string) (*sql.DB, error) {
	_ = ctx
	_ = dbPath
	return nil, retiredStateStoreError()
}
