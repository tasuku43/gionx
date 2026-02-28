---
title: "`kra shell completion`"
status: implemented
---

# `kra shell completion [<shell>]`

## Purpose

Print shell completion script for `kra`.

## Inputs

- `<shell>` (optional): `zsh`, `bash`, `sh`, `fish`
  - when omitted, detect from `$SHELL`
  - when detection fails, fallback to `zsh`

## Behavior

- Print completion script to stdout.
- Top-level suggestions include root commands and global flags.
- Subcommand suggestions are provided for:
  - `bootstrap`
  - `context`
  - `repo`
  - `template`
  - `shell`
  - `ws`
- Unsupported shell names fail with usage error.

## Usage examples

- zsh:
  - `source <(kra shell completion zsh)`
- bash:
  - `source <(kra shell completion bash)`
- fish:
  - `kra shell completion fish | source`

