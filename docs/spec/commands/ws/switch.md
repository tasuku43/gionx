---
title: "`kra ws switch`"
status: implemented
---

# `kra ws switch [--id <id> | --current | --select] [--cmux <id|ref>] [--format human|json]`

## Purpose

Switch focus to mapped cmux workspace from workspace action entrypoint.

## Inputs

- target mode:
  - `--id <id>`
  - `--current`
  - `--select`
- cmux target hint:
  - `--cmux <id|ref>`

## Behavior

- Uses cmux mapping-aware switch flow from workspace command entrypoint.
- `--id` scopes resolution to one `kra` workspace.
- `--current` resolves workspace from current path context.
- `--select` enables interactive fallback for workspace/cmux selection.
- `--cmux` accepts mapped cmux workspace id or `workspace:<n>` ordinal handle.
- `last_used_at` is updated after successful switch.

## Notes

- Workspace-level command shape is `kra ws switch ...` (not `kra cmux switch ...`).
