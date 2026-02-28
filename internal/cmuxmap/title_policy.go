package cmuxmap

import (
	"fmt"
	"strings"
)

const UntitledWorkspace = "(untitled)"

func FormatWorkspaceTitle(kraWorkspaceID string, kraTitle string, ordinal int) (string, error) {
	id := strings.TrimSpace(kraWorkspaceID)
	if id == "" {
		return "", fmt.Errorf("workspace id is required")
	}
	if ordinal < 1 {
		return "", fmt.Errorf("ordinal must be >= 1")
	}
	title := strings.TrimSpace(kraTitle)
	if title == "" {
		title = UntitledWorkspace
	}
	return fmt.Sprintf("%s | %s [%d]", id, title, ordinal), nil
}

func AllocateOrdinal(f *File, kraWorkspaceID string) (int, error) {
	if f == nil {
		return 0, fmt.Errorf("file is required")
	}
	if err := normalize(f); err != nil {
		return 0, err
	}
	id := strings.TrimSpace(kraWorkspaceID)
	if id == "" {
		return 0, fmt.Errorf("workspace id is required")
	}
	ws := f.Workspaces[id]
	if ws.NextOrdinal < 1 {
		ws.NextOrdinal = 1
	}
	ordinal := ws.NextOrdinal
	ws.NextOrdinal = ordinal + 1
	f.Workspaces[id] = ws
	return ordinal, nil
}
