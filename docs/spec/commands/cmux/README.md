---
title: "`kra cmux` command group"
status: implemented
---

# `kra cmux`

## Purpose

Provide a dedicated command group for cmux integration flows without changing
existing `kra ws --act go` behavior.

## Subcommands (skeleton phase)

- `kra cmux open`
- `kra cmux switch`
- `kra cmux list`
- `kra cmux status`

## Behavior

- `kra cmux --help` prints command-group usage.
- `kra cmux <subcommand> --help` prints subcommand usage.
- Unknown subcommands fail with usage (`exitUsage`).
- Non-help subcommand execution is intentionally unimplemented in this phase and
  returns `not implemented` (`exitNotImplemented`).

## Notes

- This spec covers command routing and usage contracts only.
- Functional cmux integration semantics are specified in follow-up specs.
