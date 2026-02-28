package cmuxmap

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	CurrentVersion = 1
	fileName       = "cmux-workspaces.json"
)

type Entry struct {
	CMUXWorkspaceID string `json:"cmux_workspace_id"`
	Ordinal         int    `json:"ordinal"`
	TitleSnapshot   string `json:"title_snapshot"`
	CreatedAt       string `json:"created_at"`
	LastUsedAt      string `json:"last_used_at"`
}

type WorkspaceMapping struct {
	NextOrdinal int     `json:"next_ordinal"`
	Entries     []Entry `json:"entries"`
}

type File struct {
	Version    int                         `json:"version"`
	Workspaces map[string]WorkspaceMapping `json:"workspaces"`
}

type Store struct {
	path string
}

func NewStore(root string) Store {
	return Store{path: MappingPath(root)}
}

func MappingPath(root string) string {
	return filepath.Join(root, ".kra", "state", fileName)
}

func (s Store) Load() (File, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultFile(), nil
		}
		return File{}, fmt.Errorf("read cmux mapping: %w", err)
	}

	var out File
	if err := json.Unmarshal(data, &out); err != nil {
		return File{}, fmt.Errorf("parse cmux mapping: %w", err)
	}
	if err := normalize(&out); err != nil {
		return File{}, err
	}
	return out, nil
}

func (s Store) Save(in File) error {
	if err := normalize(&in); err != nil {
		return err
	}

	data, err := json.MarshalIndent(in, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal cmux mapping: %w", err)
	}
	data = append(data, '\n')

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create cmux mapping dir: %w", err)
	}

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write cmux mapping temp: %w", err)
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("replace cmux mapping: %w", err)
	}
	return nil
}

func defaultFile() File {
	return File{
		Version:    CurrentVersion,
		Workspaces: map[string]WorkspaceMapping{},
	}
}

func normalize(in *File) error {
	if in.Version == 0 {
		in.Version = CurrentVersion
	}
	if in.Version != CurrentVersion {
		return fmt.Errorf("unsupported cmux mapping version: %d", in.Version)
	}
	if in.Workspaces == nil {
		in.Workspaces = map[string]WorkspaceMapping{}
	}

	for wsID, ws := range in.Workspaces {
		if ws.Entries == nil {
			ws.Entries = []Entry{}
		}

		sort.SliceStable(ws.Entries, func(i, j int) bool {
			if ws.Entries[i].Ordinal != ws.Entries[j].Ordinal {
				return ws.Entries[i].Ordinal < ws.Entries[j].Ordinal
			}
			return ws.Entries[i].CMUXWorkspaceID < ws.Entries[j].CMUXWorkspaceID
		})

		maxOrdinal := 0
		for i := range ws.Entries {
			if ws.Entries[i].Ordinal < 1 {
				ws.Entries[i].Ordinal = i + 1
			}
			if ws.Entries[i].Ordinal > maxOrdinal {
				maxOrdinal = ws.Entries[i].Ordinal
			}
		}
		if ws.NextOrdinal < 1 {
			ws.NextOrdinal = 1
		}
		if ws.NextOrdinal <= maxOrdinal {
			ws.NextOrdinal = maxOrdinal + 1
		}
		in.Workspaces[wsID] = ws
	}
	return nil
}
