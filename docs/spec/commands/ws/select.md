---
title: "`gionx ws` selection entrypoint policy"
status: implemented
---

# `gionx ws` / `gionx ws select`

## Purpose

Unify interactive selection into a single entrypoint while keeping operation commands explicit for automation.

## Dual-entry contract

- Human facade (interactive):
  - `gionx ws`
  - `gionx ws select`
- Agent facade (operation-fixed):
  - `gionx ws go`
  - `gionx ws close`
  - `gionx ws add-repo`
  - `gionx ws reopen`
  - `gionx ws purge`
- Both facades must converge to the same operation core behavior for each action.
- Parent-shell side effects (for example `cd`) are applied only through action-file protocol.

## Entry policy

- `gionx ws` is context-aware launcher.
- `gionx ws --id <id>` resolves launcher target explicitly by id.
- `gionx ws select` always starts from workspace selection.
- `gionx ws select --act <go|close|add-repo|reopen|purge>` skips action menu and executes fixed action.
- `gionx ws` must not auto-fallback to workspace list selection when current path cannot resolve workspace.
  unresolved invocation should fail and instruct users to run `gionx ws select`.
- `gionx ws` must resolve target workspace by either:
  - explicit `--id <id>`
  - current workspace context path (`workspaces/<id>/...` or `archive/<id>/...`)

## Selection flow

- Stage 1: select exactly one workspace from list scope.
- Stage 2: select action for selected workspace.
  - active scope: `go`, `add-repo`, `close`
  - archived scope: `reopen`, `purge`
- Stage 3: dispatch to operation command with explicit `<id>`.

## Operation command policy

- `ws go/close/reopen/purge` require explicit `<id>`.
- Operation-level `--select` is not supported.
- `ws add-repo` keeps direct/cwd resolution behavior and supports `--id`.
- `ws close` supports `--id` and cwd resolution fallback.
- Agent/non-interactive usage should prefer explicit `--id` and must not rely on interactive selectors.
