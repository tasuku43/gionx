---
title: "gionx ideas"
status: idea
updated: 2026-02-08
---

# Ideas (Not Scheduled Yet)

This file is a place to capture ideas that are *not* in the backlog yet, but are part of the current direction.
Nothing here is a spec (`docs/spec/**`) or an implementation commitment.

## 1) Jira Integration: Ticket URL → Workspace Creation

### Goal

- Accept a Jira ticket URL (e.g. `https://<your-domain>/browse/PROJ-123`) and create a `gionx` workspace for it.
- Also support “create workspaces for all tickets in a sprint” as a bulk operation.

### Minimal baseline flow (proposal)

- Input: ticket URL
- Fetch from Jira API: minimum metadata such as `issueKey`, `summary`, `status`, `assignee`
- Create a workspace equivalent to `gionx ws create` (naming can be provisional, but must handle collisions/duplicates)
- Persist an “external reference (Jira)” link on the workspace for future sync and UI display

### Sprint bulk creation (extension)

- Input: sprint identifier (board/sprint ID, or URL)
- Enumerate all tickets in the sprint and create only missing workspaces (idempotency is key)
- Failure behavior: stop-on-first-error vs allow partial success; safe re-run and recovery story

### Open questions (intentionally deferred)

- Workspace shape (worktrees, templates, bootstrapped files, etc.)
- Auth strategy (PAT / OAuth / keychain / env) and secure storage
- Naming rules (e.g. `PROJ-123-<slug>`) and collision resolution
- What Jira metadata to store in the state DB (minimal vs future-proof)
- Rate limiting, incremental sync, and tracking changes (moved issues, sprint changes, deletions)

## 2) Workspace List TUI: “Check Off” Work Like a ToDo

### Goal

- Provide a polished TUI for the workspace list (think: checking off items in a ToDo list).
- By default, checked items are archived (initially: equivalent to `ws close`).

### Experience sketch

- List view (filter, search, sort)
- Toggle selection with `Space` (or similar)
- Default action: checked → archive (close)
- Safety: undo, confirmations, dry-run, and a preview for batch operations

### Open questions

- UI framework choice (defer; specify the UX and state transitions first)
- Whether to support actions beyond “check=close” (done / snooze / defer)
- Ensuring the interaction model is safe given `close` can involve Git commits in the root

## 3) View: Agent Activity Across Workspaces

Examples:
- Claude Code / Codex CLI / Cursor CLI / Gemini CLI (assume more in the future)

### Goal

- See “which agent is running in which workspace” and “what state it is in” from a single view.
- Make long-running work visible and reduce conflicts (e.g. multiple agents operating in the same workspace).

### Minimal approach (proposal)

- Encourage running agents via `gionx`, recording metadata at start
  - Example: `gionx agent run --ws <id> -- <agent-command...>`
- Candidate fields: `workspace_id`, `agent_kind`, `started_at`, `last_heartbeat_at`,
  `status (running/succeeded/failed/unknown)`, `log_path`
- Display via TUI/CLI (may be integrated into the workspace list)

### Open questions

- Execution model (PID tracking, tmux integration, log capture, heartbeats, crash detection)
- How to treat agents started outside of `gionx` (accept as unobservable vs attempt discovery)
- Security/privacy (logs may contain sensitive data)

## 4) Logical Work State in Active Workspaces (`todo` vs `in-progress`)

### Goal

- Keep physical workspace lifecycle as-is: `active` and `archive`.
- Within `active`, derive a logical work state (`todo` or `in-progress`) at read time.
- Make status visible in `gionx ws list` and actionable commands (e.g. `go`) so users can quickly tell
  whether work has started.

### Constraints / principles

- Do not persist this logical work state in the DB.
- Compute it from observable signals at runtime (derived state only).
- Keep the classification deterministic and explainable (avoid opaque heuristics).

### Candidate signals (proposal)

- Workspace-side Git activity:
  - local commit history since workspace creation or since last lifecycle event
  - modified/untracked files in the workspace repos
- Repo-level activity hints:
  - whether branches/worktrees under `repos/` show active development movement
  - optional comparison against base/default branch for "work has diverged"

### Open questions

- Exact decision rule and precedence when signals disagree (e.g. clean tree but local commits exist)
- Performance budget for list rendering when many workspaces/repos are present
- UX wording in list/go flows (`todo`, `in-progress`, or symbols/colors) and fallback when unknown

## How to turn these into specs/backlog (proposal)

- Start by specifying “what we store in state (DB)” and the “command boundaries” in `docs/spec/**`.
- Only move items into `docs/BACKLOG.md` once ambiguity is reduced enough to define “done” precisely.

## 5) Workspace Unified Entry (Idea Competition 2026-02-08)

### Context

- Current `ws` operations are split across commands (`ws close`, `ws go`, `ws add-repo`).
- Goal: make human workflow feel like operating one workspace-management tool, while keeping automation strong.

### Compared ideas

- A: one screen with mixed `close/go` shortcuts
- B: keep command split + add unified UI entry
- C: select one workspace first, then choose action (`close/go`)
- D: plan/apply two-phase execution for destructive operations

### Agreed direction (for backlog/spec)

- Keep **dual entry**:
  - human: unified interactive launcher
  - agent: operation-fixed non-interactive commands
- Human canonical entrypoint: `gionx ws select`
- `gionx ws` becomes context-aware launcher:
  - outside workspace: same as `ws select`
  - inside workspace: skip workspace selection, show action menu for current workspace
- In-workspace action menu:
  - order: `add-repo` first, `close` second
  - `go` is excluded
- Selection model:
  - start with single-select only
  - multi-select is postponed
- Listing:
  - keep `ws list` read-only
  - add `ws ls` alias
- Shell integration:
  - move toward post-exec action protocol (`action file`) for parent-shell effects
- Scope note:
  - `ws current` is out of initial scope (possible later)

### Current status

- This idea set is now scheduled into backlog as:
  - `UX-WS-020` through `UX-WS-026` in `docs/BACKLOG.md`
- Next step is spec-first refinement for those tickets.
