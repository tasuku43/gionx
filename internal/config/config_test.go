package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadFile_MissingIsEmpty(t *testing.T) {
	cfg, err := LoadFile(filepath.Join(t.TempDir(), "missing.yaml"))
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if cfg != (Config{}) {
		t.Fatalf("LoadFile() = %+v, want zero", cfg)
	}
}

func TestLoadFile_NormalizeAndValidate(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
workspace:
  default_template: "  custom "
integration:
  jira:
    default_space: " abc "
    default_type: " JQL "
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}
	if cfg.Workspace.DefaultTemplate != "custom" {
		t.Fatalf("workspace.default_template = %q, want %q", cfg.Workspace.DefaultTemplate, "custom")
	}
	if cfg.Integration.Jira.DefaultSpace != "ABC" {
		t.Fatalf("integration.jira.default_space = %q, want %q", cfg.Integration.Jira.DefaultSpace, "ABC")
	}
	if cfg.Integration.Jira.DefaultType != JiraTypeJQL {
		t.Fatalf("integration.jira.default_type = %q, want %q", cfg.Integration.Jira.DefaultType, JiraTypeJQL)
	}
}

func TestLoadFile_InvalidTypeFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
integration:
  jira:
    default_type: board
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := LoadFile(path)
	if err == nil {
		t.Fatalf("LoadFile() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "integration.jira.default_type") {
		t.Fatalf("error = %q, want default_type hint", err)
	}
}

func TestLoadFile_SpaceProjectConflictFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(path, []byte(`
integration:
  jira:
    default_space: SRE
    default_project: APP
`), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	_, err := LoadFile(path)
	if err == nil {
		t.Fatalf("LoadFile() error = nil, want non-nil")
	}
	if !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("error = %q, want conflict hint", err)
	}
}

func TestMerge_RootOverridesGlobal(t *testing.T) {
	global := Config{
		Workspace: WorkspaceConfig{DefaultTemplate: "default"},
		Integration: IntegrationConfig{
			Jira: JiraConfig{
				DefaultSpace: "TEAM",
				DefaultType:  JiraTypeSprint,
			},
		},
	}
	root := Config{
		Workspace: WorkspaceConfig{DefaultTemplate: "custom"},
		Integration: IntegrationConfig{
			Jira: JiraConfig{
				DefaultProject: "APP",
				DefaultType:    JiraTypeJQL,
			},
		},
	}

	got := Merge(global, root)
	if got.Workspace.DefaultTemplate != "custom" {
		t.Fatalf("workspace.default_template = %q, want %q", got.Workspace.DefaultTemplate, "custom")
	}
	if got.Integration.Jira.DefaultProject != "APP" {
		t.Fatalf("integration.jira.default_project = %q, want %q", got.Integration.Jira.DefaultProject, "APP")
	}
	if got.Integration.Jira.DefaultType != JiraTypeJQL {
		t.Fatalf("integration.jira.default_type = %q, want %q", got.Integration.Jira.DefaultType, JiraTypeJQL)
	}
}
