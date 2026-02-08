package appports

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tasuku43/gionx/internal/gitutil"
	"github.com/tasuku43/gionx/internal/paths"
	"github.com/tasuku43/gionx/internal/statestore"
)

type RepoPort struct {
	EnsureDebugLogFn func(root string, tag string) error
	TouchRegistryFn  func(root string) error
}

func NewRepoPort(
	ensureDebugLogFn func(root string, tag string) error,
	touchRegistryFn func(root string) error,
) *RepoPort {
	return &RepoPort{
		EnsureDebugLogFn: ensureDebugLogFn,
		TouchRegistryFn:  touchRegistryFn,
	}
}

func (p *RepoPort) EnsureGitInPath() error {
	if err := gitutil.EnsureGitInPath(); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (p *RepoPort) ResolveRoot(cwd string) (string, error) {
	root, err := paths.ResolveExistingRoot(cwd)
	if err != nil {
		return "", fmt.Errorf("resolve GIONX_ROOT: %w", err)
	}
	return root, nil
}

func (p *RepoPort) EnsureDebugLog(root string, tag string) error {
	if p.EnsureDebugLogFn == nil {
		return nil
	}
	if err := p.EnsureDebugLogFn(root, tag); err != nil {
		return fmt.Errorf("enable debug logging: %w", err)
	}
	return nil
}

func (p *RepoPort) ResolveStateDBPath(root string) (string, error) {
	dbPath, err := paths.StateDBPathForRoot(root)
	if err != nil {
		return "", fmt.Errorf("resolve state db path: %w", err)
	}
	return dbPath, nil
}

func (p *RepoPort) ResolveRepoPoolPath() (string, error) {
	repoPoolPath, err := paths.DefaultRepoPoolPath()
	if err != nil {
		return "", fmt.Errorf("resolve repo pool path: %w", err)
	}
	return repoPoolPath, nil
}

func (p *RepoPort) OpenState(ctx context.Context, dbPath string) (*sql.DB, error) {
	db, err := statestore.Open(ctx, dbPath)
	if err != nil {
		return nil, fmt.Errorf("open state store: %w", err)
	}
	return db, nil
}

func (p *RepoPort) EnsureSettings(ctx context.Context, db *sql.DB, root string, repoPoolPath string) error {
	if err := statestore.EnsureSettings(ctx, db, root, repoPoolPath); err != nil {
		return fmt.Errorf("initialize settings: %w", err)
	}
	return nil
}

func (p *RepoPort) TouchRegistry(root string) error {
	if p.TouchRegistryFn != nil {
		if err := p.TouchRegistryFn(root); err != nil {
			return fmt.Errorf("update root registry: %w", err)
		}
		return nil
	}
	if err := TouchStateRegistry(root); err != nil {
		return fmt.Errorf("update root registry: %w", err)
	}
	return nil
}
