package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tasuku43/kra/internal/app/repocmd"
	"github.com/tasuku43/kra/internal/infra/appports"
)

func (c *CLI) runRepoAdd(args []string) int {
	outputFormat := "human"
	repoSpecs := make([]string, 0, len(args))
	for i := 0; i < len(args); i++ {
		arg := strings.TrimSpace(args[i])
		switch arg {
		case "-h", "--help", "help":
			c.printRepoAddUsage(c.Out)
			return exitOK
		case "--format":
			if i+1 >= len(args) {
				fmt.Fprintln(c.Err, "--format requires a value")
				c.printRepoAddUsage(c.Err)
				return exitUsage
			}
			outputFormat = strings.TrimSpace(args[i+1])
			i++
		default:
			if strings.HasPrefix(arg, "--format=") {
				outputFormat = strings.TrimSpace(strings.TrimPrefix(arg, "--format="))
				continue
			}
			if strings.HasPrefix(arg, "-") {
				fmt.Fprintf(c.Err, "unknown flag for repo add: %q\n", arg)
				c.printRepoAddUsage(c.Err)
				return exitUsage
			}
			repoSpecs = append(repoSpecs, arg)
		}
	}
	switch outputFormat {
	case "human", "json":
	default:
		fmt.Fprintf(c.Err, "unsupported --format: %q (supported: human, json)\n", outputFormat)
		c.printRepoAddUsage(c.Err)
		return exitUsage
	}
	if len(repoSpecs) == 0 {
		if outputFormat == "json" {
			_ = writeCLIJSON(c.Out, cliJSONResponse{
				OK:     false,
				Action: "repo.add",
				Error: &cliJSONError{
					Code:    "invalid_argument",
					Message: "repo add requires at least one <repo-spec>",
				},
			})
			return exitUsage
		}
		c.printRepoAddUsage(c.Err)
		return exitUsage
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.Err, "get working dir: %v\n", err)
		return exitError
	}
	ctx := context.Background()
	repoUC := repocmd.NewService(appports.NewRepoPort(c.ensureDebugLog, c.touchStateRegistry))
	session, err := repoUC.Run(ctx, repocmd.Request{
		CWD:           wd,
		DebugTag:      "repo-add",
		RequireGit:    true,
		TouchRegistry: true,
	})
	if err != nil {
		if outputFormat == "json" {
			_ = writeCLIJSON(c.Out, cliJSONResponse{
				OK:     false,
				Action: "repo.add",
				Error: &cliJSONError{
					Code:    "internal_error",
					Message: err.Error(),
				},
			})
			return exitError
		}
		fmt.Fprintf(c.Err, "%v\n", err)
		return exitError
	}
	c.debugf("run repo add count=%d", len(repoSpecs))

	requests := make([]repoPoolAddRequest, 0, len(repoSpecs))
	for _, arg := range repoSpecs {
		requests = append(requests, repoPoolAddRequest{RepoSpecInput: strings.TrimSpace(arg)})
	}
	if outputFormat == "json" {
		outcomes := applyRepoPoolAdds(ctx, session.RepoPoolPath, requests, repoPoolAddDefaultWorkers, c.debugf, nil)
		items := make([]map[string]any, 0, len(outcomes))
		success := 0
		for _, o := range outcomes {
			if o.Success {
				success++
			}
			items = append(items, map[string]any{
				"repo_key": o.RepoKey,
				"success":  o.Success,
				"reason":   strings.TrimSpace(o.Reason),
			})
		}
		if success == len(outcomes) {
			_ = writeCLIJSON(c.Out, cliJSONResponse{
				OK:     true,
				Action: "repo.add",
				Result: map[string]any{
					"added": success,
					"total": len(outcomes),
					"items": items,
				},
			})
			return exitOK
		}
		_ = writeCLIJSON(c.Out, cliJSONResponse{
			OK:     false,
			Action: "repo.add",
			Result: map[string]any{
				"added": success,
				"total": len(outcomes),
				"items": items,
			},
			Error: &cliJSONError{
				Code:    "conflict",
				Message: fmt.Sprintf("failed to add %d repo(s)", len(outcomes)-success),
			},
		})
		return exitError
	}

	useColorOut := writerSupportsColor(c.Out)
	printRepoPoolSection(c.Out, requests, useColorOut)
	outcomes := applyRepoPoolAddsWithProgress(ctx, session.RepoPoolPath, requests, repoPoolAddDefaultWorkers, c.debugf, c.Out, useColorOut)
	printRepoPoolAddResult(c.Out, outcomes, useColorOut)
	if repoPoolAddHadFailure(outcomes) {
		return exitError
	}
	return exitOK
}
