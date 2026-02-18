---
title: "`kra agent attach` session reattach"
status: implemented
---

# `kra agent attach`

## Purpose

Attach current terminal to an existing broker-managed agent session.

This command is for returning to an already running session from the current
workspace/repo context, not for global manager discovery.

## Scope (implemented)

- Command:
  - `kra agent attach [--session <id>]`
- Behavior:
  - resolve current `KRA_ROOT`
  - when `--session` is omitted, resolve current context scope from `cwd`
  - connect broker socket for the root hash
  - attach terminal stream to selected session PTY stream
- Session selection:
  - if `--session` is set:
    - attach directly, fail if not found
  - if `--session` is omitted:
    - non-interactive: fail (`--session is required`)
    - interactive: select from sessions in current scope

## Context Resolution Rules

- Inside `workspaces/<id>/repos/<repo-key>/...`:
  - candidate scope = same `workspace + repo`
- Inside `workspaces/<id>/...`:
  - candidate scope = same workspace
- At `KRA_ROOT` root:
  - fail (context too broad for `attach`)
- Outside `KRA_ROOT`:
  - fail

## Output Contract

- Success:
  - terminal enters attached stream until session exits or connection closes
- Errors:
  - clear reason + next action (missing broker, session not found, invalid context)
  - non-zero exit code

## Deferred to AGENT-100

- writer lease / takeover control
- dangerous key confirmation (`Ctrl-C`, `Ctrl-D`, `Ctrl-Z`)
