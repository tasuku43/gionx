package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	testBase, err := os.MkdirTemp("", "gionx-test-home-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create isolated test home: %v\n", err)
		os.Exit(1)
	}
	testHome := filepath.Join(testBase, "home")
	if err := os.MkdirAll(testHome, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "failed to create isolated HOME dir: %v\n", err)
		os.Exit(1)
	}
	_ = os.Setenv("HOME", testHome)
	_ = os.Setenv("GIONX_HOME", filepath.Join(testBase, ".gionx"))
	_ = os.Setenv("GIT_AUTHOR_NAME", "gionx-test")
	_ = os.Setenv("GIT_AUTHOR_EMAIL", "gionx-test@example.com")
	_ = os.Setenv("GIT_COMMITTER_NAME", "gionx-test")
	_ = os.Setenv("GIT_COMMITTER_EMAIL", "gionx-test@example.com")

	code := m.Run()
	_ = os.RemoveAll(testBase)
	os.Exit(code)
}
