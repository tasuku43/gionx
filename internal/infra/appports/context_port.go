package appports

import (
	"fmt"
	"time"

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

func (p *ContextPort) ResolveCurrentName(root string) (string, bool, error) {
	return stateregistry.ResolveContextNameByRoot(root)
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
			ContextName: e.ContextName,
			RootPath:    e.RootPath,
			LastUsedAt:  e.LastUsedAt,
		})
	}
	return out, nil
}

func (p *ContextPort) ResolveUseRootByName(name string) (string, bool, error) {
	root, ok, err := stateregistry.ResolveRootByContextName(name)
	if err != nil {
		return "", false, fmt.Errorf("resolve context by name: %w", err)
	}
	return root, ok, nil
}

func (p *ContextPort) RenameContext(oldName string, newName string) (string, error) {
	root, err := stateregistry.RenameContextName(oldName, newName, time.Now())
	if err != nil {
		return "", fmt.Errorf("rename context: %w", err)
	}
	return root, nil
}

func (p *ContextPort) RemoveContext(name string) (string, error) {
	root, err := stateregistry.RemoveContextName(name)
	if err != nil {
		return "", fmt.Errorf("remove context: %w", err)
	}
	return root, nil
}

func (p *ContextPort) CreateContext(name string, rawPath string) (string, error) {
	if p.ResolveUseRootFn == nil {
		return "", fmt.Errorf("validate root: resolver callback is required")
	}
	root, err := p.ResolveUseRootFn(rawPath)
	if err != nil {
		return "", fmt.Errorf("validate root: %w", err)
	}
	if err := stateregistry.SetContextName(root, name, time.Now()); err != nil {
		return "", fmt.Errorf("write context registry: %w", err)
	}
	return root, nil
}

func (p *ContextPort) WriteCurrent(root string) error {
	if err := paths.WriteCurrentContext(root); err != nil {
		return fmt.Errorf("write current context: %w", err)
	}
	return nil
}
