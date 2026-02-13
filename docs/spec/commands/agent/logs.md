---
title: "`kra agent logs` baseline"
status: implemented
---

# `kra agent logs` baseline

## Purpose

Provide a direct command to inspect logs of the tracked agent activity for a workspace.

## Scope (baseline)

- Command:
  - `kra agent logs --workspace <id> [--tail <n>] [--follow]`
- Required options:
  - `--workspace`
- Optional options:
  - `--tail` (default: `100`)
  - `--follow`
- Data source:
  - `KRA_ROOT/.kra/state/agents.json` (resolve record by `workspace_id`)
  - read `log_path` from resolved record
- Behavior:
  - if workspace record is missing: fail
  - if `log_path` is empty: fail
  - resolve relative `log_path` against `KRA_ROOT`
  - print last `N` lines
  - when `--follow` is enabled, continue streaming appended lines

## Out of scope (baseline)

- Session selection by historical run id
- Log redaction policy and retention policy
- Structured log query language
