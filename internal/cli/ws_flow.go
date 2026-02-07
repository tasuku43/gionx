package cli

import (
	"errors"
	"fmt"
	"strings"
)

var errWorkspaceFlowCanceled = errors.New("workspace flow canceled")

type workspaceSelectRiskResultFlowConfig struct {
	FlowName string

	SelectIDs func() ([]string, error)

	CollectRiskDetails func([]string) ([]workspaceRiskDetail, error)
	PrintRisk          func([]workspaceRiskDetail, bool)
	ConfirmRisk        func() (bool, error)

	ApplyOne func(string) error

	ResultVerb string
	ResultMark string
}

func (c *CLI) runWorkspaceSelectRiskResultFlow(cfg workspaceSelectRiskResultFlowConfig, useColor bool) ([]string, error) {
	if cfg.SelectIDs == nil {
		return nil, fmt.Errorf("workspace flow: SelectIDs is required")
	}
	if cfg.ApplyOne == nil {
		return nil, fmt.Errorf("workspace flow: ApplyOne is required")
	}
	flowName := strings.TrimSpace(cfg.FlowName)
	if flowName == "" {
		flowName = "workspace flow"
	}

	selectedIDs, err := cfg.SelectIDs()
	if err != nil {
		return nil, err
	}
	c.debugf("%s selected=%v", flowName, selectedIDs)

	if cfg.CollectRiskDetails != nil {
		riskItems, err := cfg.CollectRiskDetails(selectedIDs)
		if err != nil {
			return nil, err
		}
		if hasNonCleanRisk(riskItems) {
			c.debugf("%s risk stage entered", flowName)
			if cfg.PrintRisk != nil {
				cfg.PrintRisk(riskItems, useColor)
			}
			if cfg.ConfirmRisk != nil {
				ok, err := cfg.ConfirmRisk()
				if err != nil {
					return nil, err
				}
				if !ok {
					c.printWorkspaceFlowAbortedResult("canceled at Risk", useColor)
					c.debugf("%s canceled at risk stage", flowName)
					return nil, errWorkspaceFlowCanceled
				}
			}
		}
	}

	done := make([]string, 0, len(selectedIDs))
	for _, workspaceID := range selectedIDs {
		if err := cfg.ApplyOne(workspaceID); err != nil {
			return done, fmt.Errorf("apply workspace %s: %w", workspaceID, err)
		}
		done = append(done, workspaceID)
	}

	c.printWorkspaceFlowResult(cfg.ResultVerb, cfg.ResultMark, done, len(selectedIDs), useColor)
	c.debugf("%s result done=%v", flowName, done)
	return done, nil
}

func (c *CLI) printWorkspaceFlowAbortedResult(reason string, useColor bool) {
	fmt.Fprintln(c.Out)
	fmt.Fprintln(c.Out, renderResultTitle(useColor))
	fmt.Fprintf(c.Out, "%saborted: %s\n", uiIndent, reason)
}

func (c *CLI) printWorkspaceFlowResult(verb string, mark string, done []string, total int, useColor bool) {
	if verb == "" {
		verb = "Done"
	}
	if mark == "" {
		mark = "✔"
	}

	fmt.Fprintln(c.Out)
	fmt.Fprintln(c.Out, renderResultTitle(useColor))
	fmt.Fprintf(c.Out, "%s%s %d / %d\n", uiIndent, verb, len(done), total)
	for _, id := range done {
		fmt.Fprintf(c.Out, "%s%s %s\n", uiIndent, mark, id)
	}
}
