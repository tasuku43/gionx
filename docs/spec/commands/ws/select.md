---
title: "`kra ws` selection entrypoint policy"
status: implemented
---

# `kra ws` selection entrypoint

## Purpose

Unify workspace targeting into explicit modes and remove implicit cwd-based resolution.

## Dual-entry contract

- Human facade (interactive):
  - `kra ws --select`
- Non-interactive execution path (operation-fixed):
  - `kra ws <go|add-repo|remove-repo|close|reopen|purge> ...`
- Both facades must converge to the same operation core behavior for each action.
- Parent-shell side effects (for example `cd`) are applied only through action-file protocol.

## Entry policy

- `kra ws` must use explicit target mode:
  - `--id <id>`: resolve by workspace id
  - `--current`: resolve from current path context (`workspaces/<id>/...` or `archive/<id>/...`)
  - `--select`: interactive selection-first flow
- `kra ws` without any target mode must fail with usage error.
- `kra ws --id <id>` resolves launcher target explicitly by id.
- `kra ws --current` resolves launcher target from current path only when explicitly set.
- `kra ws --select` always starts from workspace selection.
- `kra ws --select <go|close|add-repo|remove-repo|reopen|unlock|purge>` skips action menu and executes fixed action.
- `kra ws --select reopen|unlock|purge` implicitly switches to archived scope.
- `kra ws --select --archived go|add-repo|remove-repo|close` must fail with usage error.
- `kra ws --select --multi` requires action.
- `kra ws --select --multi <close|reopen|purge>` enables multi-selection and executes the fixed action for each
  selected workspace.
- `kra ws --select --multi close` is active-scope only (`--archived` is invalid).
- `kra ws --select --multi reopen|purge` implicitly switches to archived scope.
- `kra ws --select --multi` runs lifecycle commits by default; `--no-commit` disables commits for selected action.
- `kra ws --select --multi --commit` is accepted for backward compatibility and keeps default behavior.
- `go|add-repo|remove-repo` are not supported in `--multi` mode.
- `kra ws` must not auto-resolve workspace from current path unless `--current` is explicitly set.
- unresolved invocation should fail and instruct users to use one of `--id`, `--current`, or `--select`.
- Backward compatibility: `kra ws select ...` may remain as alias to `kra ws --select ...`.

## Selection flow

- Stage 1: select exactly one workspace from list scope.
- Stage 1 (multi): select one or more workspaces from list scope when `--multi` is set.
- Stage 2: select action for selected workspace.
  - active scope: `go`, `add-repo`, `remove-repo`, `close`
  - archived scope: `reopen`, `purge`
- Stage 3: dispatch to operation command with explicit `<id>`.
- Stage 3 (multi): dispatch selected fixed action for each selected workspace id.

## Action routing policy

- Edit operations for existing workspace resources are routed by `ws <action>` subcommands.
- Read-only commands remain subcommands (`ws list`, `ws ls`) and resource creation remains `ws create`.
- Non-interactive usage should prefer explicit `--id` and must not rely on interactive selectors.
