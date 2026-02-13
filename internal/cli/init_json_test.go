package cli

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCLI_Init_JSON_Succeeds(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH")
	}
	setGitIdentity(t)

	root := t.TempDir()
	setKraHomeForTest(t)

	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)
	code := c.Run([]string{"init", "--format", "json", "--root", root, "--context", "jsonctx"})
	if code != exitOK {
		t.Fatalf("exit code = %d, want %d (stderr=%q)", code, exitOK, err.String())
	}
	resp := decodeJSONResponse(t, out.String())
	if !resp.OK || resp.Action != "init" {
		t.Fatalf("unexpected json response: %+v", resp)
	}
	if got, _ := resp.Result["root"].(string); got != root {
		t.Fatalf("result.root = %q, want %q", got, root)
	}
	if got, _ := resp.Result["context_name"].(string); got != "jsonctx" {
		t.Fatalf("result.context_name = %q, want %q", got, "jsonctx")
	}
	if err.Len() != 0 {
		t.Fatalf("stderr not empty: %q", err.String())
	}
}

func TestCLI_Init_JSON_RequiresRootAndContext(t *testing.T) {
	setKraHomeForTest(t)

	tests := []struct {
		name string
		args []string
	}{
		{name: "missing root", args: []string{"init", "--format", "json", "--context", "jsonctx"}},
		{name: "missing context", args: []string{"init", "--format", "json", "--root", filepath.Join(t.TempDir(), "r")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer
			var err bytes.Buffer
			c := New(&out, &err)
			code := c.Run(tt.args)
			if code != exitUsage {
				t.Fatalf("exit code = %d, want %d", code, exitUsage)
			}
			resp := decodeJSONResponse(t, out.String())
			if resp.OK || resp.Action != "init" || resp.Error.Code != "invalid_argument" {
				t.Fatalf("unexpected json response: %+v", resp)
			}
			if err.Len() != 0 {
				t.Fatalf("stderr should be empty in json mode: %q", err.String())
			}
		})
	}
}
