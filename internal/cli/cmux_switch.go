package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	appcmux "github.com/tasuku43/kra/internal/app/cmux"
	"github.com/tasuku43/kra/internal/infra/cmuxctl"
	"github.com/tasuku43/kra/internal/infra/paths"
)

type cmuxSwitchClient interface {
	ListWorkspaces(ctx context.Context) ([]cmuxctl.Workspace, error)
	SelectWorkspace(ctx context.Context, workspace string) error
}

var newCMUXSwitchClient = func() cmuxSwitchClient { return cmuxctl.NewClient() }

func (c *CLI) runCMUXSwitch(args []string) int {
	outputFormat := "human"
	workspaceID := ""
	cmuxHandle := ""
	for len(args) > 0 && strings.HasPrefix(args[0], "-") {
		switch args[0] {
		case "-h", "--help", "help":
			c.printCMUXSwitchUsage(c.Out)
			return exitOK
		case "--format":
			if len(args) < 2 {
				fmt.Fprintln(c.Err, "--format requires a value")
				c.printCMUXSwitchUsage(c.Err)
				return exitUsage
			}
			outputFormat = strings.TrimSpace(args[1])
			args = args[2:]
		case "--workspace":
			if len(args) < 2 {
				fmt.Fprintln(c.Err, "--workspace requires a value")
				c.printCMUXSwitchUsage(c.Err)
				return exitUsage
			}
			workspaceID = strings.TrimSpace(args[1])
			args = args[2:]
		case "--cmux":
			if len(args) < 2 {
				fmt.Fprintln(c.Err, "--cmux requires a value")
				c.printCMUXSwitchUsage(c.Err)
				return exitUsage
			}
			cmuxHandle = strings.TrimSpace(args[1])
			args = args[2:]
		default:
			if strings.HasPrefix(args[0], "--format=") {
				outputFormat = strings.TrimSpace(strings.TrimPrefix(args[0], "--format="))
				args = args[1:]
				continue
			}
			if strings.HasPrefix(args[0], "--workspace=") {
				workspaceID = strings.TrimSpace(strings.TrimPrefix(args[0], "--workspace="))
				args = args[1:]
				continue
			}
			if strings.HasPrefix(args[0], "--cmux=") {
				cmuxHandle = strings.TrimSpace(strings.TrimPrefix(args[0], "--cmux="))
				args = args[1:]
				continue
			}
			fmt.Fprintf(c.Err, "unknown flag for cmux switch: %q\n", args[0])
			c.printCMUXSwitchUsage(c.Err)
			return exitUsage
		}
	}
	switch outputFormat {
	case "human", "json":
	default:
		fmt.Fprintf(c.Err, "unsupported --format: %q (supported: human, json)\n", outputFormat)
		c.printCMUXSwitchUsage(c.Err)
		return exitUsage
	}
	if len(args) > 0 {
		fmt.Fprintf(c.Err, "unexpected args for cmux switch: %q\n", strings.Join(args, " "))
		c.printCMUXSwitchUsage(c.Err)
		return exitUsage
	}
	if workspaceID != "" {
		if err := validateWorkspaceID(workspaceID); err != nil {
			return c.writeCMUXSwitchError(outputFormat, "invalid_argument", workspaceID, fmt.Sprintf("invalid workspace id: %v", err), exitUsage)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return c.writeCMUXSwitchError(outputFormat, "internal_error", workspaceID, fmt.Sprintf("get working dir: %v", err), exitError)
	}
	root, err := paths.ResolveExistingRoot(wd)
	if err != nil {
		return c.writeCMUXSwitchError(outputFormat, "internal_error", workspaceID, fmt.Sprintf("resolve KRA_ROOT: %v", err), exitError)
	}

	svc := appcmux.NewService(func() appcmux.Client {
		return cmuxSwitchClientAdapter{inner: newCMUXSwitchClient()}
	}, newCMUXMapStore)
	var selector appcmux.SwitchSelector
	if outputFormat != "json" {
		selector = cmuxSwitchSelector{cli: c}
	}
	switchResult, code, msg := svc.Switch(context.Background(), root, workspaceID, cmuxHandle, outputFormat == "json", selector)
	if code != "" {
		exitCode := exitError
		if code == "invalid_argument" {
			exitCode = exitUsage
		}
		return c.writeCMUXSwitchError(outputFormat, code, workspaceID, msg, exitCode)
	}

	if outputFormat == "json" {
		_ = writeCLIJSON(c.Out, cliJSONResponse{
			OK:          true,
			Action:      "cmux.switch",
			WorkspaceID: switchResult.WorkspaceID,
			Result: map[string]any{
				"kra_workspace_id":  switchResult.WorkspaceID,
				"cmux_workspace_id": switchResult.CMUXWorkspaceID,
				"ordinal":           switchResult.Ordinal,
				"title":             switchResult.Title,
			},
		})
		return exitOK
	}

	fmt.Fprintln(c.Out, "switched cmux workspace")
	fmt.Fprintf(c.Out, "  kra: %s\n", switchResult.WorkspaceID)
	fmt.Fprintf(c.Out, "  cmux: %s\n", switchResult.CMUXWorkspaceID)
	fmt.Fprintf(c.Out, "  title: %s\n", switchResult.Title)
	return exitOK
}

type cmuxSwitchClientAdapter struct {
	inner cmuxSwitchClient
}

func (a cmuxSwitchClientAdapter) Capabilities(context.Context) (cmuxctl.Capabilities, error) {
	return cmuxctl.Capabilities{}, fmt.Errorf("unsupported")
}
func (a cmuxSwitchClientAdapter) CreateWorkspaceWithCommand(context.Context, string) (string, error) {
	return "", fmt.Errorf("unsupported")
}
func (a cmuxSwitchClientAdapter) RenameWorkspace(context.Context, string, string) error {
	return fmt.Errorf("unsupported")
}
func (a cmuxSwitchClientAdapter) SelectWorkspace(ctx context.Context, workspace string) error {
	return a.inner.SelectWorkspace(ctx, workspace)
}
func (a cmuxSwitchClientAdapter) SetStatus(context.Context, string, string, string, string, string) error {
	return fmt.Errorf("unsupported")
}
func (a cmuxSwitchClientAdapter) ListWorkspaces(ctx context.Context) ([]cmuxctl.Workspace, error) {
	return a.inner.ListWorkspaces(ctx)
}
func (a cmuxSwitchClientAdapter) Identify(context.Context, string, string) (map[string]any, error) {
	return nil, fmt.Errorf("unsupported")
}

type cmuxSwitchSelector struct {
	cli *CLI
}

func (s cmuxSwitchSelector) SelectWorkspace(candidates []appcmux.SwitchWorkspaceCandidate) (string, error) {
	items := make([]workspaceSelectorCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		items = append(items, workspaceSelectorCandidate{
			ID:    candidate.WorkspaceID,
			Title: fmt.Sprintf("%d mapped", candidate.MappedCount),
		})
	}
	selected, err := s.cli.promptWorkspaceSelectorSingle("active", "switch", items)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tty") {
			return "", fmt.Errorf("interactive workspace selection requires a TTY")
		}
		return "", err
	}
	if len(selected) != 1 {
		return "", fmt.Errorf("cmux switch requires exactly one workspace selected")
	}
	return strings.TrimSpace(selected[0]), nil
}

func (s cmuxSwitchSelector) SelectEntry(_ string, candidates []appcmux.SwitchEntryCandidate) (string, error) {
	items := make([]workspaceSelectorCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		items = append(items, workspaceSelectorCandidate{
			ID:    candidate.CMUXWorkspaceID,
			Title: candidate.Title,
		})
	}
	selected, err := s.cli.promptWorkspaceSelectorWithOptionsAndMode("active", "switch", "cmux:", "cmux", items, true)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "tty") {
			return "", fmt.Errorf("interactive cmux selection requires a TTY")
		}
		return "", err
	}
	if len(selected) != 1 {
		return "", fmt.Errorf("cmux switch requires exactly one target selected")
	}
	return strings.TrimSpace(selected[0]), nil
}

func (c *CLI) writeCMUXSwitchError(format string, code string, workspaceID string, message string, exitCode int) int {
	if format == "json" {
		_ = writeCLIJSON(c.Out, cliJSONResponse{
			OK:          false,
			Action:      "cmux.switch",
			WorkspaceID: workspaceID,
			Error: &cliJSONError{
				Code:    code,
				Message: message,
			},
		})
		return exitCode
	}
	if workspaceID != "" {
		fmt.Fprintf(c.Err, "cmux switch (%s): %s\n", workspaceID, message)
	} else {
		fmt.Fprintf(c.Err, "cmux switch: %s\n", message)
	}
	return exitCode
}
