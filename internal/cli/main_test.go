package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	testBase, err := os.MkdirTemp("", "kra-test-home-*")
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
	_ = os.Setenv("KRA_HOME", filepath.Join(testBase, ".kra"))
	_ = os.Setenv("GIT_AUTHOR_NAME", "kra-test")
	_ = os.Setenv("GIT_AUTHOR_EMAIL", "kra-test@example.com")
	_ = os.Setenv("GIT_COMMITTER_NAME", "kra-test")
	_ = os.Setenv("GIT_COMMITTER_EMAIL", "kra-test@example.com")

	code := m.Run()
	_ = os.RemoveAll(testBase)
	os.Exit(code)
}
