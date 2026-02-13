package appports

import (
	"fmt"
	"time"

	"github.com/tasuku43/kra/internal/stateregistry"
)

type InitPort struct {
	EnsureLayoutFn   func(root string) error
	TouchRegistryFn  func(root string) error
	SetContextNameFn func(root string, contextName string) error
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

func (p *InitPort) SetContextName(root string, contextName string) error {
	if p.SetContextNameFn != nil {
		if err := p.SetContextNameFn(root, contextName); err != nil {
			return fmt.Errorf("update context registry: %w", err)
		}
		return nil
	}
	if err := stateregistry.SetContextName(root, contextName, time.Now()); err != nil {
		return fmt.Errorf("update context registry: %w", err)
	}
	return nil
}
