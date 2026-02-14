---
title: "Agent Skillpack (project-local, flow-oriented)"
status: implemented
---

# Agent Skillpack

## Purpose

Define a tool-provided, project-local skillpack model that improves execution flow quality
without prescribing the user's business domain details.

## Design principle

- Skills should optimize flow (how to work), not domain content (what to conclude).
- Users should not need to author all base skills manually.
- `kra` should provide maintainable default skillpacks for effective usage.

## Scope

This concept covers:

- project-local skill source of truth:
  - `<KRA_ROOT>/.agent/skills/`
- tool-managed base skillpack contents
- guidance synchronization in `KRA_ROOT/AGENTS.md`

## Baseline skillpack (v1)

Bootstrap installs a default pack under:

- `KRA_ROOT/.agent/skills/.kra-skillpack.yaml`
- `KRA_ROOT/.agent/skills/flow-investigation/SKILL.md`
- `KRA_ROOT/.agent/skills/flow-execution/SKILL.md`
- `KRA_ROOT/.agent/skills/flow-insight-capture/SKILL.md`

Included flow patterns:

- investigation flow
- execution/change flow
- evidence/summarization flow
- insight capture proposal flow

These are process templates, not domain-specific rule bundles.

## AGENTS.md relation

`KRA_ROOT/AGENTS.md` should explicitly guide agents to:

- use project-local skills under `.agent/skills`
- prefer flow templates for consistency and traceability
- store reusable insights into workspace-local worklog paths
- propose insight capture in conversation and persist only after explicit approval

## Versioning

Skillpacks include explicit version metadata (`.kra-skillpack.yaml`) to allow:

- upgrade visibility
- compatibility checks
- controlled migration in existing roots

## Non-goals (v1)

- enforcing one fixed skill implementation for all model providers
- opinionated domain playbooks
- remote registry dependency for baseline operation
- forced overwrite of existing skill files
