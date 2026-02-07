---
title: "`gionx ws go`"
status: implemented
---

# `gionx ws go [--archived] [--emit-cd] [<id>]`

## Purpose

Jump to a workspace directory as a "start work" action.

## Inputs

- `<id>` (optional): workspace id for direct mode
- `--archived` (optional): target archived workspaces instead of active workspaces
- `--emit-cd` (optional): emit shell snippet for `cd` integration

## Behavior

### Mode selection

- If `<id>` is provided:
  - resolve the target directly
- If `<id>` is omitted:
  - launch shared selector UI (`commands/ws/selector.md`)
  - default scope is `active`; use `--archived` to switch scope
  - selection cardinality is single-only
  - UI is single-select mode (no checkbox markers / no `selected: n/m` footer)

### Target path

- active target: `GIONX_ROOT/workspaces/<id>/`
- archived target (with `--archived`): `GIONX_ROOT/archive/<id>/`

### UX detail

- In standard mode, print `Result:` and destination path.
- In `--emit-cd` mode, print only a shell snippet (`cd '<path>'`) to stdout.
- In selector mode, non-TTY invocation must fail (no fallback).

### Shell integration

- `gionx` cannot mutate the parent shell cwd directly.
- For practical navigation, support shell-evaluable output in `--emit-cd` mode.
- Expected usage style: `eval "$(gionx ws go --emit-cd)"`.

## Errors

- no matching workspace in selected scope
- invalid mixed selection (more than one selected in selector mode)
- target directory does not exist
- non-TTY invocation in selector mode
