package cmuxctl

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
)

type fakeRunner struct {
	lastName string
	lastArgs []string
	stdout   []byte
	stderr   []byte
	err      error
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) ([]byte, []byte, error) {
	f.lastName = name
	f.lastArgs = append([]string{}, args...)
	return f.stdout, f.stderr, f.err
}

func TestClientCapabilities_ParsesMethods(t *testing.T) {
	f := &fakeRunner{stdout: []byte(`{"methods":["workspace.create","workspace.select"]}`)}
	c := &Client{Runner: f, SocketPath: "/tmp/cmux.sock"}

	got, err := c.Capabilities(context.Background())
	if err != nil {
		t.Fatalf("Capabilities() error: %v", err)
	}
	wantArgs := []string{"--socket", "/tmp/cmux.sock", "--json", "capabilities"}
	if !reflect.DeepEqual(f.lastArgs, wantArgs) {
		t.Fatalf("args = %v, want %v", f.lastArgs, wantArgs)
	}
	if _, ok := got.Methods["workspace.create"]; !ok {
		t.Fatalf("workspace.create capability missing: %v", got.Methods)
	}
}

func TestClientCreateWorkspace_ParsesOKResponse(t *testing.T) {
	f := &fakeRunner{stdout: []byte("OK 1234-ABCD\n")}
	c := &Client{Runner: f}

	got, err := c.CreateWorkspace(context.Background())
	if err != nil {
		t.Fatalf("CreateWorkspace() error: %v", err)
	}
	if got != "1234-ABCD" {
		t.Fatalf("workspace id = %q, want %q", got, "1234-ABCD")
	}
	wantArgs := []string{"new-workspace"}
	if !reflect.DeepEqual(f.lastArgs, wantArgs) {
		t.Fatalf("args = %v, want %v", f.lastArgs, wantArgs)
	}
}

func TestClientCreateWorkspaceWithCommand_BuildsArgs(t *testing.T) {
	f := &fakeRunner{stdout: []byte("OK ws-1\n")}
	c := &Client{Runner: f}

	got, err := c.CreateWorkspaceWithCommand(context.Background(), "cd '/tmp/ws-1'")
	if err != nil {
		t.Fatalf("CreateWorkspaceWithCommand() error: %v", err)
	}
	if got != "ws-1" {
		t.Fatalf("workspace id = %q, want %q", got, "ws-1")
	}
	wantArgs := []string{"new-workspace", "--command", "cd '/tmp/ws-1'"}
	if !reflect.DeepEqual(f.lastArgs, wantArgs) {
		t.Fatalf("args = %v, want %v", f.lastArgs, wantArgs)
	}
}

func TestClientListWorkspaces_JSONMode(t *testing.T) {
	f := &fakeRunner{stdout: []byte(`{"workspaces":[{"id":"id1","ref":"workspace:1","index":0,"title":"t","selected":true}]}`)}
	c := &Client{Runner: f}

	got, err := c.ListWorkspaces(context.Background())
	if err != nil {
		t.Fatalf("ListWorkspaces() error: %v", err)
	}
	if len(got) != 1 || got[0].ID != "id1" || got[0].Ref != "workspace:1" || !got[0].Selected {
		t.Fatalf("unexpected workspaces: %+v", got)
	}
	wantArgs := []string{"--json", "list-workspaces"}
	if !reflect.DeepEqual(f.lastArgs, wantArgs) {
		t.Fatalf("args = %v, want %v", f.lastArgs, wantArgs)
	}
}

func TestClientRenameWorkspace_RequiresInputs(t *testing.T) {
	c := &Client{}
	if err := c.RenameWorkspace(context.Background(), "", "title"); err == nil {
		t.Fatalf("RenameWorkspace() with empty workspace should fail")
	}
	if err := c.RenameWorkspace(context.Background(), "ws", ""); err == nil {
		t.Fatalf("RenameWorkspace() with empty title should fail")
	}
}

func TestClientSelectWorkspace_RequiresInput(t *testing.T) {
	c := &Client{}
	if err := c.SelectWorkspace(context.Background(), ""); err == nil {
		t.Fatalf("SelectWorkspace() with empty workspace should fail")
	}
}

func TestClientSendText_BuildsCommandArgs(t *testing.T) {
	f := &fakeRunner{stdout: []byte("OK\n")}
	c := &Client{Runner: f}

	if err := c.SendText(context.Background(), "ws-1", "surface:1", "cd /tmp\n"); err != nil {
		t.Fatalf("SendText() error: %v", err)
	}
	wantArgs := []string{"send", "--workspace", "ws-1", "--surface", "surface:1", "cd /tmp\n"}
	if !reflect.DeepEqual(f.lastArgs, wantArgs) {
		t.Fatalf("args = %v, want %v", f.lastArgs, wantArgs)
	}
}

func TestClientCommandError_UsesStderr(t *testing.T) {
	f := &fakeRunner{
		stderr: []byte("boom"),
		err:    errors.New("exit status 1"),
	}
	c := &Client{Runner: f}
	err := c.SelectWorkspace(context.Background(), "ws-1")
	if err == nil {
		t.Fatalf("SelectWorkspace() error = nil, want non-nil")
	}
	if got := err.Error(); got == "" || !containsAll(got, []string{"cmux select-workspace", "boom"}) {
		t.Fatalf("error message = %q, want command + stderr", got)
	}
}

func containsAll(s string, parts []string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}
