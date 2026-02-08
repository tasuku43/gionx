package appports

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tasuku43/gionx/internal/paths"
	"github.com/tasuku43/gionx/internal/stateregistry"
	"github.com/tasuku43/gionx/internal/statestore"
)

type InitPort struct {
	EnsureLayoutFn  func(root string) error
	TouchRegistryFn func(root string) error
}

func NewInitPort(ensureLayoutFn func(root string) error, touchRegistryFn func(root string) error) *InitPort {
	return &InitPort{
		EnsureLayoutFn:  ensureLayoutFn,
		TouchRegistryFn: touchRegistryFn,
	}
}

func (p *InitPort) EnsureLayout(root string) error {
	if p.EnsureLayoutFn == nil {
		return fmt.Errorf("init layout: ensure layout callback is required")
	}
	if err := p.EnsureLayoutFn(root); err != nil {
		return fmt.Errorf("init layout: %w", err)
	}
	return nil
}

func (p *InitPort) EnsureState(ctx context.Context, root string) error {
	dbPath, err := paths.StateDBPathForRoot(root)
	if err != nil {
		return fmt.Errorf("resolve state db path: %w", err)
	}
	repoPoolPath, err := paths.DefaultRepoPoolPath()
	if err != nil {
		return fmt.Errorf("resolve repo pool path: %w", err)
	}
	db, err := statestore.Open(ctx, dbPath)
	if err != nil {
		return fmt.Errorf("open state store: %w", err)
	}
	defer func(db *sql.DB) { _ = db.Close() }(db)
	if err := statestore.EnsureSettings(ctx, db, root, repoPoolPath); err != nil {
		return fmt.Errorf("initialize settings: %w", err)
	}
	return nil
}

func (p *InitPort) TouchRegistry(root string) error {
	if p.TouchRegistryFn != nil {
		if err := p.TouchRegistryFn(root); err != nil {
			return fmt.Errorf("update root registry: %w", err)
		}
		return nil
	}
	if err := stateregistry.Touch(root, time.Now()); err != nil {
		return fmt.Errorf("update root registry: %w", err)
	}
	return nil
}
