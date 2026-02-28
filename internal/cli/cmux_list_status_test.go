package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/tasuku43/kra/internal/cmuxmap"
	"github.com/tasuku43/kra/internal/infra/cmuxctl"
)

type fakeCMUXStatusClient struct {
	workspaces []cmuxctl.Workspace
	err        error
}

func (f *fakeCMUXStatusClient) ListWorkspaces(context.Context) ([]cmuxctl.Workspace, error) {
	return f.workspaces, f.err
}

func TestCLI_CMUX_List_JSON_Success(t *testing.T) {
	root := prepareCurrentRootForTest(t)
	store := cmuxmap.NewStore(root)
	if err := store.Save(cmuxmap.File{
		Version: cmuxmap.CurrentVersion,
		Workspaces: map[string]cmuxmap.WorkspaceMapping{
			"WS1": {
				NextOrdinal: 3,
				Entries: []cmuxmap.Entry{
					{CMUXWorkspaceID: "CMUX-1", Ordinal: 1, TitleSnapshot: "WS1 | one [1]"},
					{CMUXWorkspaceID: "CMUX-2", Ordinal: 2, TitleSnapshot: "WS1 | one [2]"},
				},
			},
		},
	}); err != nil {
		t.Fatalf("save mapping: %v", err)
	}
	prev := newCMUXListClient
	newCMUXListClient = func() cmuxListClient {
		return &fakeCMUXStatusClient{
			workspaces: []cmuxctl.Workspace{
				{ID: "CMUX-1"},
				{ID: "CMUX-2"},
			},
		}
	}
	t.Cleanup(func() { newCMUXListClient = prev })

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"cmux", "list", "--format", "json", "--workspace", "WS1"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q out=%q)", code, exitOK, err.String(), out.String())
	}
	if err.Len() != 0 {
		t.Fatalf("stderr should be empty: %q", err.String())
	}

	var resp struct {
		OK     bool `json:"ok"`
		Result struct {
			Items          []map[string]any `json:"items"`
			RuntimeChecked bool             `json:"runtime_checked"`
			PrunedCount    int              `json:"pruned_count"`
		} `json:"result"`
	}
	if uerr := json.Unmarshal(out.Bytes(), &resp); uerr != nil {
		t.Fatalf("json unmarshal error: %v", uerr)
	}
	if !resp.OK || len(resp.Result.Items) != 2 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	if !resp.Result.RuntimeChecked || resp.Result.PrunedCount != 0 {
		t.Fatalf("unexpected runtime/prune summary: %+v", resp.Result)
	}
}

func TestCLI_CMUX_List_JSON_PrunesMissingEntries(t *testing.T) {
	root := prepareCurrentRootForTest(t)
	store := cmuxmap.NewStore(root)
	if err := store.Save(cmuxmap.File{
		Version: cmuxmap.CurrentVersion,
		Workspaces: map[string]cmuxmap.WorkspaceMapping{
			"WS1": {
				NextOrdinal: 3,
				Entries: []cmuxmap.Entry{
					{CMUXWorkspaceID: "CMUX-1", Ordinal: 1, TitleSnapshot: "WS1 | one [1]"},
					{CMUXWorkspaceID: "CMUX-2", Ordinal: 2, TitleSnapshot: "WS1 | one [2]"},
				},
			},
		},
	}); err != nil {
		t.Fatalf("save mapping: %v", err)
	}
	prev := newCMUXListClient
	newCMUXListClient = func() cmuxListClient {
		return &fakeCMUXStatusClient{
			workspaces: []cmuxctl.Workspace{
				{ID: "CMUX-1"},
			},
		}
	}
	t.Cleanup(func() { newCMUXListClient = prev })

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"cmux", "list", "--format", "json", "--workspace", "WS1"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q out=%q)", code, exitOK, err.String(), out.String())
	}
	if err.Len() != 0 {
		t.Fatalf("stderr should be empty: %q", err.String())
	}

	var resp struct {
		OK     bool `json:"ok"`
		Result struct {
			Items []struct {
				CMUXID string `json:"cmux_workspace_id"`
			} `json:"items"`
			PrunedCount int `json:"pruned_count"`
		} `json:"result"`
	}
	if uerr := json.Unmarshal(out.Bytes(), &resp); uerr != nil {
		t.Fatalf("json unmarshal error: %v", uerr)
	}
	if !resp.OK || len(resp.Result.Items) != 1 || resp.Result.Items[0].CMUXID != "CMUX-1" || resp.Result.PrunedCount != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}

	updated, lerr := store.Load()
	if lerr != nil {
		t.Fatalf("load mapping: %v", lerr)
	}
	if len(updated.Workspaces["WS1"].Entries) != 1 || updated.Workspaces["WS1"].Entries[0].CMUXWorkspaceID != "CMUX-1" {
		t.Fatalf("mapping should be pruned to CMUX-1 only: %+v", updated.Workspaces["WS1"].Entries)
	}
}

func TestCLI_CMUX_Status_JSON_ReportsExists(t *testing.T) {
	root := prepareCurrentRootForTest(t)
	store := cmuxmap.NewStore(root)
	if err := store.Save(cmuxmap.File{
		Version: cmuxmap.CurrentVersion,
		Workspaces: map[string]cmuxmap.WorkspaceMapping{
			"WS1": {
				NextOrdinal: 3,
				Entries: []cmuxmap.Entry{
					{CMUXWorkspaceID: "CMUX-1", Ordinal: 1, TitleSnapshot: "WS1 | one [1]"},
					{CMUXWorkspaceID: "CMUX-2", Ordinal: 2, TitleSnapshot: "WS1 | one [2]"},
				},
			},
		},
	}); err != nil {
		t.Fatalf("save mapping: %v", err)
	}

	prev := newCMUXStatusClient
	newCMUXStatusClient = func() cmuxStatusClient {
		return &fakeCMUXStatusClient{
			workspaces: []cmuxctl.Workspace{
				{ID: "CMUX-1"},
			},
		}
	}
	t.Cleanup(func() { newCMUXStatusClient = prev })

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"cmux", "status", "--format", "json", "--workspace", "WS1"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q out=%q)", code, exitOK, err.String(), out.String())
	}
	if err.Len() != 0 {
		t.Fatalf("stderr should be empty: %q", err.String())
	}

	var resp struct {
		OK     bool `json:"ok"`
		Result struct {
			Items []struct {
				CMUXID string `json:"cmux_workspace_id"`
				Exists bool   `json:"exists"`
			} `json:"items"`
		} `json:"result"`
	}
	if uerr := json.Unmarshal(out.Bytes(), &resp); uerr != nil {
		t.Fatalf("json unmarshal error: %v", uerr)
	}
	if !resp.OK || len(resp.Result.Items) != 2 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	gotExists := map[string]bool{}
	for _, it := range resp.Result.Items {
		gotExists[it.CMUXID] = it.Exists
	}
	if !gotExists["CMUX-1"] || gotExists["CMUX-2"] {
		t.Fatalf("exists flags = %+v, want CMUX-1=true CMUX-2=false", gotExists)
	}
}

func TestCLI_CMUX_List_JSON_DoesNotPruneWhenRuntimeIsEmpty(t *testing.T) {
	root := prepareCurrentRootForTest(t)
	store := cmuxmap.NewStore(root)
	if err := store.Save(cmuxmap.File{
		Version: cmuxmap.CurrentVersion,
		Workspaces: map[string]cmuxmap.WorkspaceMapping{
			"WS1": {
				NextOrdinal: 2,
				Entries: []cmuxmap.Entry{
					{CMUXWorkspaceID: "CMUX-1", Ordinal: 1, TitleSnapshot: "WS1 | one [1]"},
				},
			},
		},
	}); err != nil {
		t.Fatalf("save mapping: %v", err)
	}
	prev := newCMUXListClient
	newCMUXListClient = func() cmuxListClient {
		return &fakeCMUXStatusClient{workspaces: []cmuxctl.Workspace{}}
	}
	t.Cleanup(func() { newCMUXListClient = prev })

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"cmux", "list", "--format", "json", "--workspace", "WS1"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q out=%q)", code, exitOK, err.String(), out.String())
	}

	var resp struct {
		OK     bool `json:"ok"`
		Result struct {
			Items []struct {
				CMUXID string `json:"cmux_workspace_id"`
			} `json:"items"`
			PrunedCount int `json:"pruned_count"`
		} `json:"result"`
	}
	if uerr := json.Unmarshal(out.Bytes(), &resp); uerr != nil {
		t.Fatalf("json unmarshal error: %v", uerr)
	}
	if !resp.OK || len(resp.Result.Items) != 1 || resp.Result.PrunedCount != 0 {
		t.Fatalf("unexpected response: %+v", resp)
	}

	after, lerr := store.Load()
	if lerr != nil {
		t.Fatalf("load mapping: %v", lerr)
	}
	if len(after.Workspaces["WS1"].Entries) != 1 || after.Workspaces["WS1"].Entries[0].CMUXWorkspaceID != "CMUX-1" {
		t.Fatalf("mapping should remain unchanged when runtime empty: %+v", after.Workspaces["WS1"].Entries)
	}
}
