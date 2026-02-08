package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tasuku43/gionx/internal/app/contextcmd"
	"github.com/tasuku43/gionx/internal/paths"
	"github.com/tasuku43/gionx/internal/stateregistry"
)

func (c *CLI) runContext(args []string) int {
	if len(args) == 0 {
		c.printContextUsage(c.Err)
		return exitUsage
	}

	switch args[0] {
	case "-h", "--help", "help":
		c.printContextUsage(c.Out)
		return exitOK
	case "current":
		return c.runContextCurrent(args[1:])
	case "list":
		return c.runContextList(args[1:])
	case "use":
		return c.runContextUse(args[1:])
	default:
		fmt.Fprintf(c.Err, "unknown command: %q\n", strings.Join(append([]string{"context"}, args[0]), " "))
		c.printContextUsage(c.Err)
		return exitUsage
	}
}

func (c *CLI) runContextCurrent(args []string) int {
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help":
			c.printContextUsage(c.Out)
			return exitOK
		}
		fmt.Fprintf(c.Err, "unexpected args for context current: %q\n", strings.Join(args, " "))
		c.printContextUsage(c.Err)
		return exitUsage
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.Err, "get working dir: %v\n", err)
		return exitError
	}
	svc := contextcmd.NewService(&contextAdapter{})
	root, err := svc.Current(wd)
	if err != nil {
		fmt.Fprintf(c.Err, "resolve GIONX_ROOT: %v\n", err)
		return exitError
	}
	fmt.Fprintln(c.Out, root)
	return exitOK
}

func (c *CLI) runContextList(args []string) int {
	if len(args) > 0 {
		switch args[0] {
		case "-h", "--help", "help":
			c.printContextUsage(c.Out)
			return exitOK
		}
		fmt.Fprintf(c.Err, "unexpected args for context list: %q\n", strings.Join(args, " "))
		c.printContextUsage(c.Err)
		return exitUsage
	}

	svc := contextcmd.NewService(&contextAdapter{})
	entries, err := svc.List()
	if err != nil {
		fmt.Fprintf(c.Err, "%v\n", err)
		return exitError
	}
	if len(entries) == 0 {
		fmt.Fprintln(c.Out, "(none)")
		return exitOK
	}

	fmt.Fprintln(c.Out, "Contexts:")
	for _, e := range entries {
		last := time.Unix(e.LastUsedAt, 0).UTC().Format(time.RFC3339)
		fmt.Fprintf(c.Out, "%s%s  last_used_at=%s\n", uiIndent, e.RootPath, last)
	}
	return exitOK
}

func (c *CLI) runContextUse(args []string) int {
	if len(args) == 0 {
		c.printContextUsage(c.Err)
		return exitUsage
	}
	if len(args) > 1 {
		fmt.Fprintf(c.Err, "unexpected args for context use: %q\n", strings.Join(args[1:], " "))
		c.printContextUsage(c.Err)
		return exitUsage
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		c.printContextUsage(c.Out)
		return exitOK
	}

	svc := contextcmd.NewService(&contextAdapter{})
	root, err := svc.Use(args[0])
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "validate root:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		case strings.HasPrefix(err.Error(), "write current context:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		default:
			fmt.Fprintf(c.Err, "run context use usecase: %v\n", err)
		}
		return exitError
	}
	useColorOut := writerSupportsColor(c.Out)
	printResultSection(c.Out, useColorOut, styleSuccess(fmt.Sprintf("Context set: %s", root), useColorOut))
	return exitOK
}

type contextAdapter struct{}

func (a *contextAdapter) ResolveCurrentRoot(cwd string) (string, error) {
	return paths.ResolveExistingRoot(cwd)
}

func (a *contextAdapter) ListEntries() ([]contextcmd.Entry, error) {
	registryPath, err := stateregistry.Path()
	if err != nil {
		return nil, fmt.Errorf("resolve root registry path: %w", err)
	}
	entries, err := stateregistry.Load(registryPath)
	if err != nil {
		return nil, err
	}
	out := make([]contextcmd.Entry, 0, len(entries))
	for _, e := range entries {
		out = append(out, contextcmd.Entry{
			RootPath:   e.RootPath,
			LastUsedAt: e.LastUsedAt,
		})
	}
	return out, nil
}

func (a *contextAdapter) ResolveUseRoot(raw string) (string, error) {
	root, err := resolveContextUseRoot(raw)
	if err != nil {
		return "", fmt.Errorf("validate root: %w", err)
	}
	return root, nil
}

func (a *contextAdapter) WriteCurrent(root string) error {
	if err := paths.WriteCurrentContext(root); err != nil {
		return fmt.Errorf("write current context: %w", err)
	}
	return nil
}

func resolveContextUseRoot(raw string) (string, error) {
	root, err := filepath.Abs(strings.TrimSpace(raw))
	if err != nil {
		return "", err
	}
	root = filepath.Clean(root)

	fi, err := os.Stat(root)
	if err == nil {
		if !fi.IsDir() {
			return "", fmt.Errorf("not a directory: %s", root)
		}
		return root, nil
	}
	if !os.IsNotExist(err) {
		return "", err
	}

	parent := filepath.Dir(root)
	pfi, err := os.Stat(parent)
	if err != nil || !pfi.IsDir() {
		if err == nil {
			return "", fmt.Errorf("parent is not a directory: %s", parent)
		}
		return "", fmt.Errorf("parent directory missing: %s", parent)
	}
	return root, nil
}
