package workspacerisk

import "strings"

type RepoStatus struct {
	Upstream    string
	AheadCount  int
	BehindCount int
	Dirty       bool
	Detached    bool
	HeadMissing bool
	Error       error
}

type RepoState string
type WorkspaceRisk string

const (
	RepoStateUnknown  RepoState = "unknown"
	RepoStateDirty    RepoState = "dirty"
	RepoStateDiverged RepoState = "diverged"
	RepoStateUnpushed RepoState = "unpushed"
	RepoStateClean    RepoState = "clean"
)

const (
	WorkspaceRiskUnknown  WorkspaceRisk = "unknown"
	WorkspaceRiskDirty    WorkspaceRisk = "dirty"
	WorkspaceRiskDiverged WorkspaceRisk = "diverged"
	WorkspaceRiskUnpushed WorkspaceRisk = "unpushed"
	WorkspaceRiskClean    WorkspaceRisk = "clean"
)

func ClassifyRepoStatus(status RepoStatus) RepoState {
	if status.Error != nil {
		return RepoStateUnknown
	}
	if status.Dirty {
		return RepoStateDirty
	}
	if status.Detached || status.HeadMissing {
		return RepoStateUnknown
	}
	if strings.TrimSpace(status.Upstream) == "" {
		return RepoStateUnknown
	}
	if status.AheadCount > 0 && status.BehindCount > 0 {
		return RepoStateDiverged
	}
	if status.AheadCount > 0 {
		return RepoStateUnpushed
	}
	return RepoStateClean
}

func Aggregate(repos []RepoState) WorkspaceRisk {
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range repos {
		switch repo {
		case RepoStateUnknown:
			hasUnknown = true
		case RepoStateDirty:
			hasDirty = true
		case RepoStateDiverged:
			hasDiverged = true
		case RepoStateUnpushed:
			hasUnpushed = true
		}
	}
	switch {
	case hasUnknown:
		return WorkspaceRiskUnknown
	case hasDirty:
		return WorkspaceRiskDirty
	case hasDiverged:
		return WorkspaceRiskDiverged
	case hasUnpushed:
		return WorkspaceRiskUnpushed
	default:
		return WorkspaceRiskClean
	}
}

func AggregateForState(repos []RepoState) WorkspaceRisk {
	hasDirty := false
	hasUnknown := false
	hasDiverged := false
	hasUnpushed := false
	for _, repo := range repos {
		switch repo {
		case RepoStateDirty:
			hasDirty = true
		case RepoStateUnknown:
			hasUnknown = true
		case RepoStateDiverged:
			hasDiverged = true
		case RepoStateUnpushed:
			hasUnpushed = true
		}
	}
	switch {
	case hasDirty:
		return WorkspaceRiskDirty
	case hasUnknown:
		return WorkspaceRiskUnknown
	case hasDiverged:
		return WorkspaceRiskDiverged
	case hasUnpushed:
		return WorkspaceRiskUnpushed
	default:
		return WorkspaceRiskClean
	}
}
