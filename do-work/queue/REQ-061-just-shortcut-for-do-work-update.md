---
id: REQ-061
title: Add a just shortcut for do-work updates
status: pending
created_at: 2026-07-29T20:43:22Z
user_request: UR-009
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
---

# Add a just shortcut for do-work updates

## What

Extend the target repository's installed `just` recipes so users can run `just run-do-work-update` to update the project-local do-work skill programmatically, without invoking the agent-driven `do-work update` action. Install the recipe alongside the existing `just run-kanban` shortcut.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

- The command belongs in the target repository where do-work is installed, like the existing `run-kanban` recipe installed by `do-work install just-kanban`.
- The update command must update the do-work skill programmatically, avoiding the token cost of asking an agent to run `do-work update`.
- Preserve the update flow's project-local safety and integrity checks rather than creating an unchecked overwrite path.

## Red-Green Proof
**RED prompt/case:** After installing the do-work `just` recipes in a target repository, `just run-do-work-update` is unavailable, so updating requires an agent invocation of `do-work update`.
**Why RED now:** The installed recipe block provides only `run-kanban`, `kanban-static`, and `kanban-summary`; no programmatic update shortcut is installed.
**GREEN when:** Installing the relevant do-work `just` recipes adds a parseable `run-do-work-update` command in the target repository, and running it performs the project-local do-work update with the same safety checks expected of the existing update workflow.
**Validation:** Inferred during capture

---
*Source: UR-009 — "do-work capture-request: add a feature to install in the target repository (where the skill is installed) a just command, like install run-kanban, add also install run-do-work-update which updates programmatically the do-work skill, so we don't need to burn tokens to call `do-work update`"*

Think carefully before answering.
