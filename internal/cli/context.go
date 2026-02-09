package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tasuku43/gionx/internal/app/contextcmd"
	"github.com/tasuku43/gionx/internal/infra/appports"
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
	case "create":
		return c.runContextCreate(args[1:])
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
	svc := contextcmd.NewService(appports.NewContextPort(resolveContextUseRoot))
	current, err := svc.Current(wd)
	if err != nil {
		fmt.Fprintf(c.Err, "resolve current context: %v\n", err)
		return exitError
	}
	fmt.Fprintln(c.Out, current)
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

	svc := contextcmd.NewService(appports.NewContextPort(resolveContextUseRoot))
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
		name := strings.TrimSpace(e.ContextName)
		if name == "" {
			name = "(unnamed)"
		}
		fmt.Fprintf(c.Out, "%s%s  path=%s  last_used_at=%s\n", uiIndent, name, e.RootPath, last)
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

	svc := contextcmd.NewService(appports.NewContextPort(resolveContextUseRoot))
	root, err := svc.Use(args[0])
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "resolve context by name:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		case strings.HasPrefix(err.Error(), "context not found:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		case strings.HasPrefix(err.Error(), "write current context:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		default:
			fmt.Fprintf(c.Err, "run context use usecase: %v\n", err)
		}
		return exitError
	}
	useColorOut := writerSupportsColor(c.Out)
	printResultSection(c.Out, useColorOut, styleSuccess(fmt.Sprintf("Context selected: %s", args[0]), useColorOut), styleMuted(fmt.Sprintf("path: %s", root), useColorOut))
	return exitOK
}

func (c *CLI) runContextCreate(args []string) int {
	if len(args) == 0 {
		c.printContextUsage(c.Err)
		return exitUsage
	}
	if args[0] == "-h" || args[0] == "--help" || args[0] == "help" {
		c.printContextUsage(c.Out)
		return exitOK
	}
	name := strings.TrimSpace(args[0])
	if name == "" {
		fmt.Fprintln(c.Err, "context name is required")
		c.printContextUsage(c.Err)
		return exitUsage
	}
	rawPath := ""
	useNow := false
	rest := args[1:]
	for len(rest) > 0 {
		switch rest[0] {
		case "--path":
			if len(rest) < 2 {
				fmt.Fprintln(c.Err, "--path requires a value")
				c.printContextUsage(c.Err)
				return exitUsage
			}
			rawPath = strings.TrimSpace(rest[1])
			rest = rest[2:]
		case "--use":
			useNow = true
			rest = rest[1:]
		default:
			if strings.HasPrefix(rest[0], "--path=") {
				rawPath = strings.TrimSpace(strings.TrimPrefix(rest[0], "--path="))
				rest = rest[1:]
				continue
			}
			fmt.Fprintf(c.Err, "unexpected args for context create: %q\n", strings.Join(rest, " "))
			c.printContextUsage(c.Err)
			return exitUsage
		}
	}
	if rawPath == "" {
		fmt.Fprintln(c.Err, "context create requires --path <path>")
		c.printContextUsage(c.Err)
		return exitUsage
	}

	svc := contextcmd.NewService(appports.NewContextPort(resolveContextUseRoot))
	root, err := svc.Create(name, rawPath, useNow)
	if err != nil {
		switch {
		case strings.HasPrefix(err.Error(), "validate root:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		case strings.HasPrefix(err.Error(), "write context registry:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		case strings.HasPrefix(err.Error(), "write current context:"):
			fmt.Fprintf(c.Err, "%v\n", err)
		default:
			fmt.Fprintf(c.Err, "run context create usecase: %v\n", err)
		}
		return exitError
	}

	useColorOut := writerSupportsColor(c.Out)
	if useNow {
		printResultSection(c.Out, useColorOut, styleSuccess(fmt.Sprintf("Context created and selected: %s", name), useColorOut), styleMuted(fmt.Sprintf("path: %s", root), useColorOut))
		return exitOK
	}
	printResultSection(c.Out, useColorOut, styleSuccess(fmt.Sprintf("Context created: %s", name), useColorOut), styleMuted(fmt.Sprintf("path: %s", root), useColorOut))
	return exitOK
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
