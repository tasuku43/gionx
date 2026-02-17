---
title: "`kra agent logs` retirement"
status: planned
---

# `kra agent logs` retirement (v3 draft)

## Purpose

Retire `kra agent logs` from the default agent surface in v3.

## Rationale

- v2 `logs` depended on externally managed log files and did not provide reliable runtime truth.
- v3 shifts to PTY-managed runtime session state under `KRA_HOME/state/agents/`.
- MVP focus is current runtime visibility (`run/list/board/stop`), not file-log inspection.

## Scope (v3 draft)

- `kra agent logs` command is removed from CLI usage.
- `log_path` is removed from runtime activity schema.
- `kra agent run` no longer accepts `--log-path`.

## Out of scope (v3 draft)

- Reintroducing log inspection in MVP.
- Structured trace/event querying.
