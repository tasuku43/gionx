package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const workspaceMetaFilename = ".gionx.meta.json"

type workspaceMetaFile struct {
	SchemaVersion int                        `json:"schema_version"`
	Workspace     workspaceMetaWorkspace     `json:"workspace"`
	ReposRestore  []workspaceMetaRepoRestore `json:"repos_restore"`
}

type workspaceMetaWorkspace struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	SourceURL   string `json:"source_url"`
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type workspaceMetaRepoRestore struct {
	RepoUID   string `json:"repo_uid"`
	RepoKey   string `json:"repo_key"`
	RemoteURL string `json:"remote_url"`
	Alias     string `json:"alias"`
	Branch    string `json:"branch"`
	BaseRef   string `json:"base_ref"`
}

func newWorkspaceMetaFileForCreate(id string, description string, now int64) workspaceMetaFile {
	return workspaceMetaFile{
		SchemaVersion: 1,
		Workspace: workspaceMetaWorkspace{
			ID:          id,
			Description: description,
			SourceURL:   "",
			Status:      "active",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		ReposRestore: make([]workspaceMetaRepoRestore, 0),
	}
}

func writeWorkspaceMetaFile(wsPath string, meta workspaceMetaFile) error {
	if strings.TrimSpace(wsPath) == "" {
		return fmt.Errorf("workspace path is required")
	}
	metaPath := filepath.Join(wsPath, workspaceMetaFilename)
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal workspace meta: %w", err)
	}
	b = append(b, '\n')

	tmp, err := os.CreateTemp(wsPath, ".gionx-meta-*.tmp")
	if err != nil {
		return fmt.Errorf("create workspace meta temp file: %w", err)
	}
	tmpPath := tmp.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write workspace meta temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close workspace meta temp file: %w", err)
	}
	if err := os.Rename(tmpPath, metaPath); err != nil {
		return fmt.Errorf("replace workspace meta file %s: %w", metaPath, err)
	}
	return nil
}
