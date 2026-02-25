package repodiscovery

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Repo struct {
	RepoUID   string
	RepoKey   string
	RemoteURL string
}

type Provider interface {
	Name() string
	CheckAuth(ctx context.Context) error
	ListOrgRepos(ctx context.Context, org string) ([]Repo, error)
}

type ProviderFactory func() Provider

var (
	providerRegistryMu sync.RWMutex
	providerRegistry   = map[string]ProviderFactory{
		"github": func() Provider { return NewGitHubGHProvider(nil) },
	}
)

func RegisterProvider(name string, factory ProviderFactory) error {
	normalized := normalizeProviderName(name)
	if normalized == "" {
		return fmt.Errorf("provider name is required")
	}
	if factory == nil {
		return fmt.Errorf("provider factory is required")
	}
	providerRegistryMu.Lock()
	defer providerRegistryMu.Unlock()
	if _, exists := providerRegistry[normalized]; exists {
		return fmt.Errorf("provider already registered: %q", normalized)
	}
	providerRegistry[normalized] = factory
	return nil
}

func SupportedProviders() []string {
	providerRegistryMu.RLock()
	defer providerRegistryMu.RUnlock()
	names := make([]string, 0, len(providerRegistry))
	for name := range providerRegistry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func normalizeProviderName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func NewProvider(name string) (Provider, error) {
	normalized := normalizeProviderName(name)
	if normalized == "" {
		normalized = "github"
	}
	providerRegistryMu.RLock()
	factory, ok := providerRegistry[normalized]
	providerRegistryMu.RUnlock()
	if !ok {
		supported := SupportedProviders()
		return nil, fmt.Errorf("unsupported provider: %q (supported: %s)", normalized, strings.Join(supported, ", "))
	}
	return factory(), nil
}
