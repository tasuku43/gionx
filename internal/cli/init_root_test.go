package cli

import (
	"path/filepath"
	"testing"
)

func TestNormalizeInitRoot_ExpandsTildeToHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := normalizeInitRoot("~")
	if err != nil {
		t.Fatalf("normalizeInitRoot(\"~\") error = %v", err)
	}
	if got != home {
		t.Fatalf("normalizeInitRoot(\"~\") = %q, want %q", got, home)
	}
}

func TestNormalizeInitRoot_ExpandsTildePrefixToHomeSubdir(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := normalizeInitRoot("~/testroot")
	if err != nil {
		t.Fatalf("normalizeInitRoot(\"~/testroot\") error = %v", err)
	}
	want := filepath.Join(home, "testroot")
	if got != want {
		t.Fatalf("normalizeInitRoot(\"~/testroot\") = %q, want %q", got, want)
	}
}
