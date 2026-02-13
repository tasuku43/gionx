---
title: "`kra shell init`"
status: implemented
---

# `kra shell init [<shell>]`

## Purpose

Print shell integration script that applies parent-shell side effects via action-file protocol.

## Inputs

- `<shell>` (optional): shell name (`zsh`, `bash`, `sh`, `fish`)
  - when omitted, detect from `$SHELL`
  - when detection fails, fallback to `zsh`

## Behavior

- Print eval-ready script to stdout.
- Script contains:
  - one-time setup hint comment
  - shell function `kra` override
  - for all command paths, set `KRA_SHELL_ACTION_FILE=<tempfile>` and let `kra` emit post-exec action
    (for example `cd '<path>'`) into that file when needed
  - after command success, apply action file content if present
- Unsupported shell names must fail with usage error.

## Output examples

- POSIX shells:
  - `eval "$(kra shell init zsh)"`
- fish:
  - `eval (kra shell init fish)`
