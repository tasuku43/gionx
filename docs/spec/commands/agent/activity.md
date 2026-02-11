---
title: "`gionx agent` activity tracking"
status: planned
---

# `gionx agent` activity tracking (baseline)

## Purpose

Provide a baseline command and data contract to observe agent execution across workspaces:
- which agent is running in which workspace
- current lifecycle state (`running` / `succeeded` / `failed` / `unknown`)
- where to inspect logs

This spec is planning-only. No implementation commitment is implied until backlog completion.

## Scope (baseline)

- Define command boundary under `gionx agent ...`.
- Define minimal tracked fields:
  - `workspace_id`
  - `agent_kind`
  - `started_at`
  - `last_heartbeat_at`
  - `status`
  - `log_path`
- Define read model for "activity list" view (CLI/TUI integration is implementation detail).

## Out of scope (baseline)

- Strong process supervision guarantees (PID ownership, hard crash recovery)
- External process discovery for agents launched outside `gionx`
- Long-term retention policy and log redaction policy

## Open questions to resolve before implementation

- Execution model (`gionx` launcher only vs partial external detection)
- Heartbeat mechanism and stale/unknown threshold
- Persistence location and recovery strategy under partial failures
- Security/privacy constraints for log handling
