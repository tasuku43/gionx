package cli

import (
	"testing"
	"time"
)

func TestAgentBrokerSessionOutput_IdleOutputBecomesRunning(t *testing.T) {
	now := time.Unix(200, 0)
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{
			SessionID:    "s-1",
			RuntimeState: "idle",
			UpdatedAt:    100,
			Seq:          1,
		},
		seqParser: newAgentTerminalSequenceParser(),
	}

	_, snapshot, _ := session.appendOutputAndSnapshotWritable([]byte("working..."), now)
	if snapshot == nil || snapshot.RuntimeState != "running" {
		got := "<nil>"
		if snapshot != nil {
			got = snapshot.RuntimeState
		}
		t.Fatalf("runtime should become running after output activity, got=%s", got)
	}
}

func TestAgentBrokerSessionOutput_DoesNotDropToIdleOnOutput(t *testing.T) {
	now := time.Unix(200, 0)
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{
			SessionID:    "s-1",
			RuntimeState: "running",
			UpdatedAt:    100,
			Seq:          1,
		},
		seqParser: newAgentTerminalSequenceParser(),
	}

	_, snapshot, _ := session.appendOutputAndSnapshotWritable([]byte("\n› prompt\r\n"), now.Add(500*time.Millisecond))
	if snapshot != nil && snapshot.RuntimeState == "idle" {
		t.Fatalf("runtime should not turn idle while output is still arriving")
	}
}

func TestAgentBrokerSessionMarkIdleOnSilence_AfterOutput(t *testing.T) {
	now := time.Unix(300, 0)
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{
			SessionID:    "s-1",
			RuntimeState: "running",
			UpdatedAt:    now.Unix() - int64(agentRuntimeIdleAfterSilence/time.Second) - 1,
			Seq:          2,
		},
		lastOutputAt: now.Add(-agentRuntimeIdleAfterSilence - time.Second),
	}

	snapshot := session.markIdleOnSilence(now)
	if snapshot == nil {
		t.Fatalf("snapshot should be emitted when running output is silent")
	}
	if snapshot.RuntimeState != "idle" {
		t.Fatalf("runtime_state=%q, want=idle", snapshot.RuntimeState)
	}
}

func TestAgentBrokerSessionOutput_OSC133D_BecomesWaitingInput(t *testing.T) {
	now := time.Unix(400, 0)
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{
			SessionID:    "s-1",
			RuntimeState: "running",
			UpdatedAt:    100,
			Seq:          1,
		},
		seqParser: newAgentTerminalSequenceParser(),
	}

	_, snapshot, _ := session.appendOutputAndSnapshotWritable([]byte("\x1b]133;D;0\x07"), now)
	if snapshot == nil {
		t.Fatalf("snapshot should be emitted for state transition")
	}
	if snapshot.RuntimeState != "waiting_input" {
		t.Fatalf("runtime_state=%q, want=waiting_input", snapshot.RuntimeState)
	}
}

func TestAgentBrokerSessionMarkIdleOnSilence_DoesNotChangeWaitingInput(t *testing.T) {
	now := time.Unix(500, 0)
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{
			SessionID:    "s-1",
			RuntimeState: "waiting_input",
			UpdatedAt:    100,
			Seq:          1,
		},
		lastOutputAt: now.Add(-agentRuntimeIdleAfterSilence - time.Second),
	}

	snapshot := session.markIdleOnSilence(now)
	if snapshot != nil {
		t.Fatalf("waiting_input should not be downgraded to idle by silence timeout")
	}
}

func TestAgentBrokerSessionAcquireControl_FirstOwnerThenSpectator(t *testing.T) {
	session := &agentBrokerSession{
		attachments: map[string]*agentBrokerAttachment{
			"a-1": {clientID: "c-1"},
			"a-2": {clientID: "c-2"},
		},
	}

	inputOK1, resizeOK1 := session.acquireControl("c-1", "interactive")
	if !inputOK1 || !resizeOK1 {
		t.Fatalf("first interactive client should own input/resize lease")
	}
	inputOK2, resizeOK2 := session.acquireControl("c-2", "interactive")
	if inputOK2 || resizeOK2 {
		t.Fatalf("second interactive client should be spectator while first owner is attached")
	}
}

func TestAgentBrokerSessionAcquireControl_SpectatorModeNeverOwns(t *testing.T) {
	session := &agentBrokerSession{
		attachments: map[string]*agentBrokerAttachment{
			"a-1": {clientID: "c-1"},
		},
	}
	inputOK, resizeOK := session.acquireControl("c-1", "spectator")
	if inputOK || resizeOK {
		t.Fatalf("spectator mode should never acquire input/resize lease")
	}
}

func TestAgentBrokerHandleResizeRequest_DeniesNonOwner(t *testing.T) {
	session := &agentBrokerSession{
		record: agentRuntimeSessionRecord{SessionID: "s-1"},
		attachments: map[string]*agentBrokerAttachment{
			"a-1": {clientID: "c-owner"},
			"a-2": {clientID: "c-other"},
		},
		inputOwnerClientID:  "c-owner",
		resizeOwnerClientID: "c-owner",
	}
	server := &agentBrokerServer{
		sessions: map[string]*agentBrokerSession{"s-1": session},
	}

	resp := server.handleResizeRequest(agentBrokerRequest{
		SessionID: "s-1",
		ClientID:  "c-other",
		Cols:      120,
		Rows:      40,
	})
	if resp.OK {
		t.Fatalf("non-owner resize should be denied")
	}
	if resp.Error != "resize lease denied" {
		t.Fatalf("error=%q, want=resize lease denied", resp.Error)
	}
}
