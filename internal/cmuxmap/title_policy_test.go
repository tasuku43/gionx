package cmuxmap

import "testing"

func TestFormatWorkspaceTitle(t *testing.T) {
	got, err := FormatWorkspaceTitle("MVP-020", "implement auth", 2)
	if err != nil {
		t.Fatalf("FormatWorkspaceTitle() error: %v", err)
	}
	want := "MVP-020 | implement auth"
	if got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
}

func TestFormatWorkspaceTitle_EmptyTitleUsesFallback(t *testing.T) {
	got, err := FormatWorkspaceTitle("MVP-020", "   ", 1)
	if err != nil {
		t.Fatalf("FormatWorkspaceTitle() error: %v", err)
	}
	want := "MVP-020 | (untitled)"
	if got != want {
		t.Fatalf("title = %q, want %q", got, want)
	}
}

func TestAllocateOrdinal_InitializesAndIncrements(t *testing.T) {
	f := File{Version: CurrentVersion, Workspaces: map[string]WorkspaceMapping{}}

	first, err := AllocateOrdinal(&f, "WS1")
	if err != nil {
		t.Fatalf("AllocateOrdinal(first) error: %v", err)
	}
	second, err := AllocateOrdinal(&f, "WS1")
	if err != nil {
		t.Fatalf("AllocateOrdinal(second) error: %v", err)
	}
	if first != 1 || second != 2 {
		t.Fatalf("ordinals = (%d, %d), want (1, 2)", first, second)
	}
	if f.Workspaces["WS1"].NextOrdinal != 3 {
		t.Fatalf("next_ordinal = %d, want 3", f.Workspaces["WS1"].NextOrdinal)
	}
}

func TestAllocateOrdinal_RequiresWorkspaceID(t *testing.T) {
	f := defaultFile()
	if _, err := AllocateOrdinal(&f, ""); err == nil {
		t.Fatalf("AllocateOrdinal() error = nil, want non-nil")
	}
}
