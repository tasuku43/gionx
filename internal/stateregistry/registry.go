package stateregistry

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/tasuku43/kra/internal/paths"
)

type Entry struct {
	ContextName string `json:"context_name,omitempty"`
	RootPath    string `json:"root_path"`
	FirstSeenAt int64  `json:"first_seen_at"`
	LastUsedAt  int64  `json:"last_used_at"`
}

type filePayload struct {
	Entries []fileEntryPayload `json:"entries"`
}

type fileEntryPayload struct {
	ContextName string `json:"context_name,omitempty"`
	RootPath    string `json:"root_path"`
	StateDBPath string `json:"state_db_path,omitempty"` // legacy compatibility (read-only)
	FirstSeenAt int64  `json:"first_seen_at"`
	LastUsedAt  int64  `json:"last_used_at"`
}

func Path() (string, error) {
	return paths.RegistryPath()
}

func Load(path string) ([]Entry, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("registry path is required")
	}

	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read root registry %s: %w", path, err)
	}
	if len(strings.TrimSpace(string(b))) == 0 {
		return nil, nil
	}

	var p filePayload
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("root registry is malformed: %s (fix or remove this file and retry): %w", path, err)
	}
	out := make([]Entry, 0, len(p.Entries))
	for _, e := range p.Entries {
		out = append(out, Entry{
			ContextName: strings.TrimSpace(e.ContextName),
			RootPath:    e.RootPath,
			FirstSeenAt: e.FirstSeenAt,
			LastUsedAt:  e.LastUsedAt,
		})
	}
	return out, nil
}

func Touch(rootPath string, now time.Time) error {
	rootAbs, err := cleanAbs(rootPath)
	if err != nil {
		return fmt.Errorf("resolve root_path: %w", err)
	}

	registryPath, err := Path()
	if err != nil {
		return fmt.Errorf("resolve root registry path: %w", err)
	}

	entries, err := Load(registryPath)
	if err != nil {
		return err
	}

	nowUnix := now.Unix()
	if nowUnix <= 0 {
		nowUnix = time.Now().Unix()
	}

	updated := false
	for i := range entries {
		if entries[i].RootPath != rootAbs {
			continue
		}
		if entries[i].FirstSeenAt <= 0 {
			entries[i].FirstSeenAt = nowUnix
		}
		if entries[i].LastUsedAt < nowUnix {
			entries[i].LastUsedAt = nowUnix
		}
		updated = true
		break
	}
	if !updated {
		entries = append(entries, Entry{
			RootPath:    rootAbs,
			FirstSeenAt: nowUnix,
			LastUsedAt:  nowUnix,
		})
	}

	slices.SortFunc(entries, func(a, b Entry) int {
		return strings.Compare(a.RootPath, b.RootPath)
	})

	if err := writeAtomic(registryPath, entries); err != nil {
		return err
	}
	return nil
}

func writeAtomic(path string, entries []Entry) error {
	entryPayloads := make([]fileEntryPayload, 0, len(entries))
	for _, e := range entries {
		entryPayloads = append(entryPayloads, fileEntryPayload{
			ContextName: strings.TrimSpace(e.ContextName),
			RootPath:    e.RootPath,
			FirstSeenAt: e.FirstSeenAt,
			LastUsedAt:  e.LastUsedAt,
		})
	}
	payload := filePayload{Entries: entryPayloads}
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal root registry: %w", err)
	}
	b = append(b, '\n')

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create root registry dir %s: %w", dir, err)
	}

	tmp, err := os.CreateTemp(dir, ".registry-*.tmp")
	if err != nil {
		return fmt.Errorf("create root registry temp file in %s: %w", dir, err)
	}
	tmpPath := tmp.Name()
	cleanup := func() {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
	}

	if _, err := tmp.Write(b); err != nil {
		cleanup()
		return fmt.Errorf("write root registry temp file %s: %w", tmpPath, err)
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("close root registry temp file %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("replace root registry %s: %w", path, err)
	}
	return nil
}

func cleanAbs(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("empty path")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

func SetContextName(rootPath string, contextName string, now time.Time) error {
	rootAbs, err := cleanAbs(rootPath)
	if err != nil {
		return fmt.Errorf("resolve root_path: %w", err)
	}
	contextName = strings.TrimSpace(contextName)
	if contextName == "" {
		return fmt.Errorf("context name is required")
	}

	registryPath, err := Path()
	if err != nil {
		return fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := Load(registryPath)
	if err != nil {
		return err
	}

	nowUnix := now.Unix()
	if nowUnix <= 0 {
		nowUnix = time.Now().Unix()
	}

	rootIdx := -1
	for i := range entries {
		if entries[i].RootPath == rootAbs {
			rootIdx = i
			break
		}
	}
	for i := range entries {
		if strings.TrimSpace(entries[i].ContextName) != contextName {
			continue
		}
		if i != rootIdx {
			return fmt.Errorf("context name already exists: %s", contextName)
		}
	}

	if rootIdx >= 0 {
		entries[rootIdx].ContextName = contextName
		if entries[rootIdx].FirstSeenAt <= 0 {
			entries[rootIdx].FirstSeenAt = nowUnix
		}
		if entries[rootIdx].LastUsedAt < nowUnix {
			entries[rootIdx].LastUsedAt = nowUnix
		}
	} else {
		entries = append(entries, Entry{
			ContextName: contextName,
			RootPath:    rootAbs,
			FirstSeenAt: nowUnix,
			LastUsedAt:  nowUnix,
		})
	}

	slices.SortFunc(entries, func(a, b Entry) int {
		return strings.Compare(a.RootPath, b.RootPath)
	})
	return writeAtomic(registryPath, entries)
}

func ResolveRootByContextName(contextName string) (string, bool, error) {
	contextName = strings.TrimSpace(contextName)
	if contextName == "" {
		return "", false, fmt.Errorf("context name is required")
	}
	registryPath, err := Path()
	if err != nil {
		return "", false, fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := Load(registryPath)
	if err != nil {
		return "", false, err
	}
	for _, e := range entries {
		if strings.TrimSpace(e.ContextName) == contextName {
			return e.RootPath, true, nil
		}
	}
	return "", false, nil
}

func ResolveContextNameByRoot(rootPath string) (string, bool, error) {
	rootAbs, err := cleanAbs(rootPath)
	if err != nil {
		return "", false, fmt.Errorf("resolve root_path: %w", err)
	}
	registryPath, err := Path()
	if err != nil {
		return "", false, fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := Load(registryPath)
	if err != nil {
		return "", false, err
	}
	for _, e := range entries {
		if e.RootPath == rootAbs {
			name := strings.TrimSpace(e.ContextName)
			if name == "" {
				return "", false, nil
			}
			return name, true, nil
		}
	}
	return "", false, nil
}

func RenameContextName(oldName string, newName string, now time.Time) (string, error) {
	oldName = strings.TrimSpace(oldName)
	newName = strings.TrimSpace(newName)
	if oldName == "" || newName == "" {
		return "", fmt.Errorf("old and new context names are required")
	}
	if oldName == newName {
		return "", fmt.Errorf("new context name must be different")
	}

	registryPath, err := Path()
	if err != nil {
		return "", fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := Load(registryPath)
	if err != nil {
		return "", err
	}

	oldIdx := -1
	for i := range entries {
		if strings.TrimSpace(entries[i].ContextName) == oldName {
			oldIdx = i
			break
		}
	}
	if oldIdx < 0 {
		return "", fmt.Errorf("context not found: %s", oldName)
	}
	for i := range entries {
		if i == oldIdx {
			continue
		}
		if strings.TrimSpace(entries[i].ContextName) == newName {
			return "", fmt.Errorf("context name already exists: %s", newName)
		}
	}

	nowUnix := now.Unix()
	if nowUnix <= 0 {
		nowUnix = time.Now().Unix()
	}
	entries[oldIdx].ContextName = newName
	if entries[oldIdx].LastUsedAt < nowUnix {
		entries[oldIdx].LastUsedAt = nowUnix
	}

	slices.SortFunc(entries, func(a, b Entry) int {
		return strings.Compare(a.RootPath, b.RootPath)
	})
	if err := writeAtomic(registryPath, entries); err != nil {
		return "", err
	}
	return entries[oldIdx].RootPath, nil
}

func RemoveContextName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("context name is required")
	}
	registryPath, err := Path()
	if err != nil {
		return "", fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := Load(registryPath)
	if err != nil {
		return "", err
	}

	idx := -1
	for i := range entries {
		if strings.TrimSpace(entries[i].ContextName) == name {
			idx = i
			break
		}
	}
	if idx < 0 {
		return "", fmt.Errorf("context not found: %s", name)
	}
	root := entries[idx].RootPath
	entries = append(entries[:idx], entries[idx+1:]...)
	if err := writeAtomic(registryPath, entries); err != nil {
		return "", err
	}
	return root, nil
}
