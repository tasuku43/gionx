package cmuxmap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMappingPath(t *testing.T) {
	root := "/tmp/kra-root"
	got := MappingPath(root)
	want := filepath.Join(root, ".kra", "state", "cmux-workspaces.json")
	if got != want {
		t.Fatalf("MappingPath() = %q, want %q", got, want)
	}
}

func TestStoreLoad_NotExist_ReturnsDefault(t *testing.T) {
	root := t.TempDir()
	s := NewStore(root)

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if got.Version != CurrentVersion {
		t.Fatalf("version = %d, want %d", got.Version, CurrentVersion)
	}
	if got.Workspaces == nil {
		t.Fatalf("workspaces should not be nil")
	}
	if len(got.Workspaces) != 0 {
		t.Fatalf("workspaces len = %d, want 0", len(got.Workspaces))
	}
}

func TestStoreSaveLoad_RoundTripAndNormalize(t *testing.T) {
	root := t.TempDir()
	s := NewStore(root)

	in := File{
		Version: CurrentVersion,
		Workspaces: map[string]WorkspaceMapping{
			"WS1": {
				NextOrdinal: 0,
				Entries: []Entry{
					{CMUXWorkspaceID: "B", Ordinal: 2},
					{CMUXWorkspaceID: "A", Ordinal: 0},
				},
			},
		},
	}
	if err := s.Save(in); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	got, err := s.Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	ws, ok := got.Workspaces["WS1"]
	if !ok {
		t.Fatalf("workspace WS1 not found")
	}
	if ws.Entries[0].CMUXWorkspaceID != "A" || ws.Entries[0].Ordinal != 1 {
		t.Fatalf("first entry = %+v, want id=A ordinal=1", ws.Entries[0])
	}
	if ws.NextOrdinal != 3 {
		t.Fatalf("next_ordinal = %d, want 3", ws.NextOrdinal)
	}

	if _, err := os.Stat(MappingPath(root)); err != nil {
		t.Fatalf("mapping file stat error: %v", err)
	}
}

func TestStoreLoad_UnsupportedVersionErrors(t *testing.T) {
	root := t.TempDir()
	s := NewStore(root)
	path := MappingPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error: %v", err)
	}
	if err := os.WriteFile(path, []byte(`{"version":99,"workspaces":{}}`), 0o644); err != nil {
		t.Fatalf("WriteFile() error: %v", err)
	}

	_, err := s.Load()
	if err == nil {
		t.Fatalf("Load() error = nil, want unsupported version error")
	}
}
