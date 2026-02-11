package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	JiraTypeSprint = "sprint"
	JiraTypeJQL    = "jql"
)

type Config struct {
	Workspace   WorkspaceConfig   `yaml:"workspace"`
	Integration IntegrationConfig `yaml:"integration"`
}

type WorkspaceConfig struct {
	DefaultTemplate string `yaml:"default_template"`
}

type IntegrationConfig struct {
	Jira JiraConfig `yaml:"jira"`
}

type JiraConfig struct {
	DefaultSpace   string `yaml:"default_space"`
	DefaultProject string `yaml:"default_project"`
	DefaultType    string `yaml:"default_type"`
}

func LoadFile(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("read config file: %w", err)
	}
	if strings.TrimSpace(string(b)) == "" {
		return Config{}, nil
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse yaml: %w", err)
	}
	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) Normalize() {
	c.Workspace.DefaultTemplate = strings.TrimSpace(c.Workspace.DefaultTemplate)
	c.Integration.Jira.DefaultSpace = strings.ToUpper(strings.TrimSpace(c.Integration.Jira.DefaultSpace))
	c.Integration.Jira.DefaultProject = strings.ToUpper(strings.TrimSpace(c.Integration.Jira.DefaultProject))
	c.Integration.Jira.DefaultType = strings.ToLower(strings.TrimSpace(c.Integration.Jira.DefaultType))
}

func (c Config) Validate() error {
	issues := make([]string, 0, 2)
	if c.Integration.Jira.DefaultType != "" &&
		c.Integration.Jira.DefaultType != JiraTypeSprint &&
		c.Integration.Jira.DefaultType != JiraTypeJQL {
		issues = append(issues, "integration.jira.default_type must be one of: sprint, jql")
	}
	if c.Integration.Jira.DefaultSpace != "" && c.Integration.Jira.DefaultProject != "" {
		issues = append(issues, "integration.jira.default_space and integration.jira.default_project cannot be combined")
	}
	if len(issues) == 0 {
		return nil
	}
	return fmt.Errorf("invalid config: %s", strings.Join(issues, "; "))
}

func Merge(global Config, root Config) Config {
	global.Normalize()
	root.Normalize()

	out := global
	if root.Workspace.DefaultTemplate != "" {
		out.Workspace.DefaultTemplate = root.Workspace.DefaultTemplate
	}
	if root.Integration.Jira.DefaultSpace != "" {
		out.Integration.Jira.DefaultSpace = root.Integration.Jira.DefaultSpace
	}
	if root.Integration.Jira.DefaultProject != "" {
		out.Integration.Jira.DefaultProject = root.Integration.Jira.DefaultProject
	}
	if root.Integration.Jira.DefaultType != "" {
		out.Integration.Jira.DefaultType = root.Integration.Jira.DefaultType
	}
	out.Normalize()
	return out
}
