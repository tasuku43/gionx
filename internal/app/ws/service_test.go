package ws

import (
	"context"
	"testing"
)

type stubResolver struct {
	byID map[string]WorkspaceRef
}

func (s stubResolver) ResolveFromPath(_ context.Context, _ string) (WorkspaceRef, bool, error) {
	return WorkspaceRef{}, false, nil
}

func (s stubResolver) ResolveByID(_ context.Context, id string) (WorkspaceRef, bool, error) {
	ref, ok := s.byID[id]
	return ref, ok, nil
}

type stubSelector struct{}

func (stubSelector) SelectWorkspace(_ context.Context, _ Scope, _ string, _ bool) (string, error) {
	return "", nil
}

func (stubSelector) SelectAction(_ context.Context, _ WorkspaceRef, _ bool) (Action, error) {
	return "", nil
}

func TestServiceRun_FixedRunAgent_AllowedForActive(t *testing.T) {
	svc := NewService(
		stubSelector{},
		stubResolver{
			byID: map[string]WorkspaceRef{
				"WS-1": {ID: "WS-1", Status: ScopeActive},
			},
		},
	)

	got, err := svc.Run(context.Background(), LauncherRequest{
		WorkspaceID: "WS-1",
		FixedAction: ActionRunAgent,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got.Action != ActionRunAgent {
		t.Fatalf("result action = %q, want %q", got.Action, ActionRunAgent)
	}
	if got.Workspace.ID != "WS-1" || got.Workspace.Status != ScopeActive {
		t.Fatalf("result workspace = %+v, want id=WS-1 status=active", got.Workspace)
	}
}

func TestServiceRun_FixedRunAgent_RejectedForArchived(t *testing.T) {
	svc := NewService(
		stubSelector{},
		stubResolver{
			byID: map[string]WorkspaceRef{
				"WS-9": {ID: "WS-9", Status: ScopeArchived},
			},
		},
	)

	_, err := svc.Run(context.Background(), LauncherRequest{
		WorkspaceID: "WS-9",
		FixedAction: ActionRunAgent,
	})
	if err != ErrActionNotAllowed {
		t.Fatalf("Run() error = %v, want %v", err, ErrActionNotAllowed)
	}
}
