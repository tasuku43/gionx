package cli

import (
	"bytes"
	"testing"

	"github.com/tasuku43/gionx/internal/core/workspacerisk"
)

func TestGolden_WSSection_FlowResult(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)

	c.printWorkspaceFlowResult("Purged", "✔", []string{"TEST-001", "TEST-002"}, 2, false)
	assertGolden(t, "ws_flow_result.golden", out.String())
}

func TestGolden_WSSection_FlowAbortedResult(t *testing.T) {
	var out bytes.Buffer
	var err bytes.Buffer
	c := New(&out, &err)

	c.printWorkspaceFlowAbortedResult("canceled at Risk", false)
	assertGolden(t, "ws_flow_aborted_result.golden", out.String())
}

func TestGolden_WSSection_CloseRisk(t *testing.T) {
	var out bytes.Buffer
	items := []workspaceRiskDetail{
		{
			id:   "WS1",
			risk: workspacerisk.WorkspaceRiskDirty,
			perRepo: []repoRiskItem{
				{alias: "repo-a", state: workspacerisk.RepoStateDirty},
				{alias: "repo-b", state: workspacerisk.RepoStateUnpushed},
			},
		},
		{
			id:      "WS2",
			risk:    workspacerisk.WorkspaceRiskClean,
			perRepo: []repoRiskItem{{alias: "repo-c", state: workspacerisk.RepoStateClean}},
		},
	}

	printRiskSection(&out, items, false)
	assertGolden(t, "ws_close_risk_section.golden", out.String())
}

func TestGolden_WSSection_PurgeRisk(t *testing.T) {
	var out bytes.Buffer
	selectedIDs := []string{"WS1"}
	riskMeta := map[string]purgeWorkspaceMeta{
		"WS1": {
			status: "active",
			risk:   workspacerisk.WorkspaceRiskDirty,
			perRepo: []repoRiskItem{
				{alias: "repo1", state: workspacerisk.RepoStateDirty},
			},
		},
	}

	printPurgeRiskSection(&out, selectedIDs, riskMeta, false)
	assertGolden(t, "ws_purge_risk_section.golden", out.String())
}
