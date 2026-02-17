---
title: "`kra agent run` baseline"
status: planned
---

# `kra agent run` v3 draft

## Purpose

Start one agent session under PTY and register runtime activity in `KRA_HOME` state.

## Scope (v3 draft)

- Command:
  - `kra agent run [--workspace <id>] [--repo <repo-key>] [--kind <agent-kind>]`
- Interactive behavior:
  - with no args, command enters interactive selector flow
  - workspace selector must include active workspaces only
  - execution target must always be selected when `--repo` is not given:
    - run at workspace scope
    - run at repo scope (pick repo key)
  - if `--kind` is omitted, prompt for kind selection
- Flags removed in v3:
  - `--task`
  - `--instruction`
  - `--status`
  - `--log-path`
- Runtime write target:
  - `KRA_HOME/state/agents/<root-hash>/<session-id>.json`
- Behavior:
  - resolve current `KRA_ROOT`
  - resolve run target (workspace scope or repo scope)
  - start child process on a PTY
  - create new `session_id` and write initial runtime state
  - while process is alive, update `updated_at` and `seq` based on PTY I/O and lifecycle events
  - on process exit, persist final record (`runtime_state=exited`, `exit_code=<code>`)
  - print a human confirmation line including `session_id`

## Out of scope (v3 draft)

- Rich instruction/task metadata capture in `run`.
- Cross-process command injection channel from `agent run` MVP.
