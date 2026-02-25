package repodiscovery

import (
	"context"
	"strings"
	"testing"
)

type stubProvider struct {
	name string
}

func (s stubProvider) Name() string { return s.name }

func (s stubProvider) CheckAuth(context.Context) error { return nil }

func (s stubProvider) ListOrgRepos(context.Context, string) ([]Repo, error) { return nil, nil }

func TestRegisterProvider_AllowsCustomProvider(t *testing.T) {
	name := "custom-provider-test"
	if err := RegisterProvider(name, func() Provider {
		return stubProvider{name: "custom"}
	}); err != nil {
		t.Fatalf("RegisterProvider() error: %v", err)
	}
	p, err := NewProvider(name)
	if err != nil {
		t.Fatalf("NewProvider() error: %v", err)
	}
	if got := p.Name(); got != "custom" {
		t.Fatalf("provider name = %q, want custom", got)
	}
}

func TestRegisterProvider_RejectsDuplicate(t *testing.T) {
	name := "duplicate-provider-test"
	if err := RegisterProvider(name, func() Provider {
		return stubProvider{name: "dup"}
	}); err != nil {
		t.Fatalf("RegisterProvider() first error: %v", err)
	}
	err := RegisterProvider(name, func() Provider {
		return stubProvider{name: "dup2"}
	})
	if err == nil {
		t.Fatalf("RegisterProvider() duplicate error = nil, want error")
	}
	if !strings.Contains(err.Error(), "already registered") {
		t.Fatalf("duplicate error = %q, want already registered", err.Error())
	}
}

func TestNewProvider_UnsupportedMessageIncludesSupportedList(t *testing.T) {
	_, err := NewProvider("definitely-unsupported-provider")
	if err == nil {
		t.Fatalf("NewProvider() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "supported:") {
		t.Fatalf("error = %q, want supported list", err.Error())
	}
}
