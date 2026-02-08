package appports

import (
	"fmt"

	"github.com/tasuku43/gionx/internal/app/contextcmd"
	"github.com/tasuku43/gionx/internal/paths"
	"github.com/tasuku43/gionx/internal/stateregistry"
)

type ContextPort struct {
	ResolveUseRootFn func(raw string) (string, error)
}

func NewContextPort(resolveUseRootFn func(raw string) (string, error)) *ContextPort {
	return &ContextPort{ResolveUseRootFn: resolveUseRootFn}
}

func (p *ContextPort) ResolveCurrentRoot(cwd string) (string, error) {
	return paths.ResolveExistingRoot(cwd)
}

func (p *ContextPort) ListEntries() ([]contextcmd.Entry, error) {
	registryPath, err := stateregistry.Path()
	if err != nil {
		return nil, fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := stateregistry.Load(registryPath)
	if err != nil {
		return nil, err
	}
	out := make([]contextcmd.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, contextcmd.Entry{
			RootPath:   e.RootPath,
			LastUsedAt: e.LastUsedAt,
		})
	}
	return out, nil
}

func (p *ContextPort) ResolveUseRoot(raw string) (string, error) {
	if p.ResolveUseRootFn == nil {
		return "", fmt.Errorf("validate root: resolver callback is required")
	}
	root, err := p.ResolveUseRootFn(raw)
	if err != nil {
		return "", fmt.Errorf("validate root: %w", err)
	}
	return root, nil
}

func (p *ContextPort) WriteCurrent(root string) error {
	if err := paths.WriteCurrentContext(root); err != nil {
		return fmt.Errorf("write current context: %w", err)
	}
	return nil
}
