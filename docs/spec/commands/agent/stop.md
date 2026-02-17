---
title: "`kra agent stop` baseline"
status: planned
---

# `kra agent stop` v3 draft

## Purpose

Stop one running agent session managed by `kra agent run`.

## Scope (v3 draft)

- Command:
  - `kra agent stop (--session <id> | --workspace <id> [--repo <repo-key>] [--kind <agent-kind>])`
- Required options:
  - either `--session`, or `--workspace` selector set
- Optional options:
  - `--repo`
  - `--kind`
- Runtime data source:
  - `KRA_HOME/state/agents/<root-hash>/<session-id>.json`
- Behavior:
  - resolve current `KRA_ROOT`
  - locate target session
  - if session is already `exited`, return idempotent success
  - send termination signal to target process
  - wait bounded grace period, then force kill if still alive
  - persist final runtime state (`exited`) with `updated_at` and `seq` increment
  - print final status line with `session_id`

## Out of scope (v3 draft)

- Distributed stop across remote hosts.
- Multi-step approval flow for stop.
