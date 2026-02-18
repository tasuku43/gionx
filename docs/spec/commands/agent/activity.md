---
title: "`kra agent` activity tracking"
status: implemented
---

# `kra agent` activity tracking

## Purpose

Provide runtime visibility for agent sessions across workspaces with state files under `KRA_HOME` (outside `KRA_ROOT` Git tree).

## Scope (implemented)

- Command boundary: `kra agent ...`
- Command surface:
  - `kra agent run`
  - `kra agent attach`
  - `kra agent stop`
  - `kra agent list` (`ls` alias)
  - `kra agent board`
- Discoverability policy:
  - `kra agent` is executable directly
  - root help intentionally does not list `agent`
- Runtime state location:
  - `KRA_HOME/state/agents/<root-hash>/<session-id>.json`
  - `KRA_HOME` default: `~/.kra`

## Runtime record schema (session file)

- Required fields:
  - `session_id`
  - `root_path`
  - `workspace_id`
  - `execution_scope` (`workspace` | `repo`)
  - `repo_key` (empty for `workspace` scope)
  - `kind`
  - `pid`
  - `started_at`
  - `updated_at`
  - `seq`
  - `runtime_state` (`running` | `idle` | `exited` | `unknown`)
  - `exit_code` (nullable; set when `runtime_state=exited`)

## State write rules (implemented)

- one session = one file
- writes are atomic (`tmp -> rename`)
- `seq` is monotonic per session

## Runtime state model (operator-facing)

- `running`: process alive
- `idle`: process alive but no recent activity signal
- `exited`: process ended
- `unknown`: runtime could not determine state reliably

## `kra agent list` / `kra agent board`

- Data source:
  - directory scan of `KRA_HOME/state/agents/<root-hash>/`
  - missing directory means empty list
- `list` output contract:
  - `tsv` is machine-friendly flat rows
  - `human` is workspace-first summary + per-session tree rows
  - child order is deterministic: `workspace` first, then `repo:<repo_key>`
- `board` output contract:
  - workspace-grouped human view
  - deterministic ordering
- Filtering:
  - workspace id
  - runtime state
  - execution location (`workspace` or `repo:<repo_key>`)
  - kind

## Deferred to AGENT-100

- snapshot fields for attach/input ownership (`attached_clients`, `writer_owner`, `lease_expires_at`)
- append-only runtime events (`events/<session-id>.jsonl`)
- lease/takeover event model
- launch mode metadata (`launch_mode`)
