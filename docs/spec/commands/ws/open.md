---
title: "`kra ws open`"
status: implemented
---

# `kra ws open [--id <id> | --current | --select] [--multi] [--concurrency <n>] [--format human|json]`

## Purpose

Open cmux workspace(s) from workspace action entrypoint.

## Inputs

- target mode:
  - `--id <id>`
  - `--current`
  - `--select`
- batch options:
  - `--multi`
  - `--concurrency <n>`

## Behavior

- Uses cmux integration flow from workspace command entrypoint.
- `--id` targets one active workspace.
- `--current` resolves workspace from current path under `workspaces/<id>/...`.
- `--select` opens workspace selector and resolves target workspace(s) interactively.
- `--multi` enables multi-target open flow.
- `--concurrency` is valid only with `--multi`.
- JSON mode remains non-interactive.

## Notes

- Parent shell cwd mutation still follows action-file protocol.
- Workspace-level command shape is `kra ws open ...` (not `kra cmux open ...`).
