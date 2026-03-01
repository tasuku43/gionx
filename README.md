# kra

`kra` is a local CLI for AI-agent-oriented, ticket-driven workspace operations.
It standardizes task workspaces on the filesystem through template-driven workspace scaffolding and per-workspace Git worktrees for the repositories each task actually needs.
The default template includes `notes/` and `artifacts/`, and you can define your own structures (for example `AGENTS.md`, `CLOUD.md`, and custom directories) to provide a predictable workspace scaffold for multi-repo execution and continuous accumulation of intermediate outputs.

`kra` works standalone for workspace lifecycle operations, and becomes even more valuable when paired with `cmux` for agent runtime/workspace orchestration.
Ticket providers are designed to be extensible; currently, Jira is supported.

## cmux integration

`kra` provides an operational framework for using `cmux` in ticket-driven, filesystem-based workflows.

In this model, `kra` acts as the glue across:

- ticket
- filesystem task workspace (`workspaces/<id>/`)
- `cmux` runtime workspace

This gives you an operational 1:1:1 mapping across ticket, task workspace, and runtime workspace, so agents can run in a consistent task-scoped workspace and continuously write intermediate outputs as work progresses.

`kra` works standalone for workspace lifecycle operations, but this `kra` + `cmux` operating model is where the overall system value becomes strongest.

## worktree management

`kra` lets you manage Git worktrees per task workspace with `kra ws add-repo` / `kra ws remove-repo`.
You can attach only the repositories needed for the current task, exactly when they become necessary.

Repositories are attached under:

- `workspaces/<id>/repos/<alias>/`

This keeps execution context task-scoped and avoids mixing temporary task outputs into a single long-lived repository.

When work spans multiple repositories, this model lets you build a task-scoped, monorepo-like execution surface quickly, so agents can use the right repository context at the right time without long setup cycles.

## Why kra exists

AI-agent-driven work now produces large volumes of intermediate outputs, not only code changes but also investigations, analyses, logs, and notes.
Those artifacts need a temporary but structured home before they are distilled into final deliverables.
In practice, people often dump them into ad hoc locations first, such as random local directories or inside an active repository, and later struggle with cleanup and traceability.

In AI-agent-driven workflows that span multiple repositories and task contexts, common problems include:

- Non-code outputs end up in inconsistent locations, making task-scoped recovery harder.
- Task context is difficult to keep organized in a way that supports fast resume.
- After task switching, the exact repo/branch combination for a task is easy to lose.
- Manual workspace rebuilds drift over time (missing repos, extra repos, wrong branch context).

`kra` addresses this by making the filesystem workspace the unit of execution and traceability.
External ticket systems remain the source of truth for task management, while `kra` provides a repeatable local operating model.

## What it gives you

- A template-driven workspace scaffold for each task:
  - start with the default template (`notes/`, `artifacts/`)
  - replace or extend it with your own structure (`AGENTS.md`, `CLOUD.md`, custom directories/files)
- Per-task, per-workspace multi-repo worktree management:
  - add only repositories required for that task
  - keep them under `workspaces/<id>/repos/<alias>/` for a predictable local execution surface
- An operational bridge between planning and runtime:
  - anchor ticket context to a filesystem workspace
  - pair that workspace with `cmux` runtime workspace(s) as an operational model
- State-first lifecycle operations with explicit transitions:
  - create, open, close, reopen, purge with clear state rules
  - archive completed workspaces by default instead of destructive deletion
- Guardrails for risky operations:
  - evaluate workspace risk (`dirty`, `unpushed`, `diverged`, `unknown`)
  - apply confirmation gates for destructive flows
- Automation-ready output contracts:
  - use a shared JSON envelope (`ok`, `action`, `workspace_id`, `result`, `error`) across supported commands

## Quickstart (5 minutes)

```sh
# 1) initialize a root (interactive)
kra init

# 2) create a workspace from your ticket id (interactive)
kra ws create TASK-1234

# 3) register repositories into the repo pool
kra repo add git@github.com:org/backend.git git@github.com:org/frontend.git

# 4) attach needed repositories to the workspace (interactive selector + prompts)
kra ws add-repo --id TASK-1234

# 5) open workspace context
kra ws open --id TASK-1234

# 6) inspect current state
kra ws dashboard

# 7) close when done (archives notes/artifacts, removes worktrees)
kra ws close --id TASK-1234
```

`kra init`, `kra ws create`, and `kra ws add-repo` will guide you with prompts in this quickstart flow.
You can also use interactive selection for workspace-targeted commands (for example: `kra ws open --select`, `kra ws close --select`).

After this flow, your task context and artifacts remain reviewable under `archive/<id>/`, while active workspace area stays clean.

## Boundaries

`kra` is intentionally focused.

- It is not a replacement for Jira or other external ticket systems.
- It is not an agent runtime/session manager.
  - Agent runtime orchestration is expected to be handled by tools such as `cmux`.
- It is not a GUI planning tool.

## Installation

### Homebrew (stable releases)

```sh
brew tap tasuku43/kra
brew install kra
```

### GitHub Releases (manual)

1. Download an archive for your OS/arch from GitHub Releases.
2. Extract and place `kra` on your `PATH`.
3. Verify with `kra version`.

### Build from source

Requirements:

- Go 1.24+
- Git

```sh
go build -o kra ./cmd/kra
./kra version
```

## Further reading

- Start here: `docs/START_HERE.md`
- Install guide: `docs/guides/INSTALL.md`
- Command reference: `docs/guides/COMMANDS.md`
- Product concept: `docs/concepts/product-concept.md`
- Specs: `docs/spec/README.md`
- Releasing: `docs/ops/RELEASING.md`

## Contributing

See `CONTRIBUTING.md`.

## Support

See `SUPPORT.md`.

## Security

See `SECURITY.md`.

## Code of Conduct

See `CODE_OF_CONDUCT.md`.

## License

See `LICENSE`.

## Maintainer

- @tasuku43
