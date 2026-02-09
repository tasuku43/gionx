---
title: "UI Selector"
status: implemented
---

# UI Selector (shared `single` / `multi`)

## Purpose

Define one shared interaction contract for interactive list selection used by `ws`, `repo`, and `context`.

## Interaction

- Cursor: `Up` / `Down`
- Confirm: `Enter`
- Cancel: `Esc` or `Ctrl+C`
- Filter input: direct typing
- Filter delete: `Backspace` / `Delete`

Multi-select:
- Toggle selection: `Space` (or full-width space)
- Footer includes `selected: n/m`

Single-select:
- `Space` does nothing
- Footer does not show `selected: n/m`
- `Enter` starts confirm transition:
  - mark focused row as selected (`●`)
  - keep frame visible briefly (`0.2s`)
  - lock cursor/input during the transition
  - then proceed to next stage
  - when `GIONX_REDUCED_MOTION=1`, skip delay and proceed immediately

## Visual model

- Selection marker is shared:
  - unselected: `○`
  - selected: `●`
- Focus marker is shared:
  - focused row starts with `>`
- Focus row may use subtle background highlight on color-capable terminals.
- Section/body indentation follows shared UI indentation constants.

## Color semantics

- Selected marker uses `accent`.
- Non-selected marker uses primary text color.
- During single-select confirm delay, non-selected rows are dimmed with `text.muted`.
- Errors/messages use semantic tokens from `docs/spec/concepts/ui-color.md`.
- No-color terminals must preserve all meaning via markers/text only.

## Reuse policy

- Selector rendering/runtime must be implemented as shared component(s), not per-command bespoke logic.
- `ws`, `repo`, and `context` interactive selections must use the same component and contract.
