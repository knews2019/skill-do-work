# Compound-Engineering Integration Guide

do-work is compatible with the [compound-engineering plugin](https://github.com/EveryInc/compound-engineering-plugin) (CE). This guide explains what that means today, what it will mean later, and how to set it up.

## The short version

- **do-work orchestrates, CE augments.** do-work runs the full REQ cycle (capture → work → review → archive) on its own. Specific seams optionally hand off to CE skills when CE is installed.
- **No hard dependency.** If CE is not installed, do-work behaves exactly as without the integration. Handoffs degrade to a saved prompt the user can run later.
- **Current integration:** one seam — after a REQ's review passes and Lessons Learned are captured, do-work offers to promote those lessons into CE's `docs/solutions/` knowledge base via the `ce-compound` skill.

## Why CE

CE is a plugin (for Claude Code, Cursor, Codex, Gemini, Kiro, OpenCode, Pi) built around the compound-engineering thesis: **make each unit of engineering work easier than the last by capturing solved problems so future work can build on them.** It ships with specialized reviewer agents (Rails, TypeScript, Python, security, data integrity), a knowledge base (`ce-compound`), brainstorming and planning skills (`ce-brainstorm`, `ce-plan`), session research (`ce-sessions`), and more.

do-work already produces the raw material CE compounds — every REQ's `## Lessons Learned` section is a small captured learning. Without the integration, those lessons live in the archived REQ and nowhere else. With the integration, relevant lessons get promoted into a durable, queryable knowledge base that CE maintains.

## Current integration point: the ce-compound handoff

### What it does

After a REQ's review completes:

1. do-work captures the `## Lessons Learned` section (existing behavior).
2. do-work assembles a structured payload from the REQ: title, inferred category from `domain`, symptoms from `## What`, solution summary from `## Implementation Summary`, and prevention guidance from Lessons Learned.
3. do-work asks the user: run `/ce-compound` now, save for later, or skip.
4. If the user runs it now, CE's compound skill writes a file under `docs/solutions/<category>/<slug>.md`.
5. do-work records the outcome on the REQ as `ce_compound_status` (and `ce_solution_path` when promoted).

The handoff **never auto-promotes**, **never blocks if CE is not installed**, and **never retries on failure**. In unattended pipeline runs it defaults to `ce_compound_status: pending` — the user can batch-dispatch later.

### Where it runs

- **Standalone mode** (`do-work review REQ-NNN`): review-work action, Step 9.5, after lesson capture.
- **Pipeline mode** (`do-work` full cycle): work action, Step 7.5, after lesson capture.

Both call the same reference file: `actions/ce-compound-handoff.md`.

### REQ frontmatter fields

Two optional fields, both set by the handoff, both absent on REQs that predate CE integration:

```yaml
ce_compound_status: promoted   # promoted | pending | declined | skipped
ce_solution_path: docs/solutions/frontend/user-avatar-default-state.md
```

| Status | Meaning |
|---|---|
| `promoted` | CE wrote a solution file. `ce_solution_path` points to it. |
| `pending` | User chose "save for later" or the pipeline ran unattended. Safe to revisit. |
| `declined` | User actively refused. Do not re-offer. |
| `skipped` | Auto-skipped because handoff conditions weren't met (empty Lessons, Route A with no gotchas). |

### Sample handoff payload

When the handoff fires, it prints something like this before asking the user:

```yaml
title: User avatar component — fallback to default state on missing image
date: 2026-04-23
category: frontend
module: src/components
problem_type: best_practice
symptoms:
  - Avatar crashed when image prop was undefined
what_didnt_work:
  - Initial CSS-modules approach conflicted with the project's styled-components convention
solution: >
  Wrapped the existing Avatar.tsx with a UserAvatar component that handles the missing-image case
  and keeps sizing constraints from AppShell.tsx's grid.
prevention:
  - Avatar sizes above 48px require checking the sidebar grid in AppShell.tsx
  - Match the project's existing styling convention before introducing a new one
tags: [frontend, react, components]
```

The user sees this, picks one of run / save / skip, and control returns to the calling action.

## Roadmap (not yet implemented)

These integration points are designed for but not yet wired up. Tracked here so users know what is coming and can lobby for ordering.

- **Reviewer agents in review-work / code-review** — dispatch `ce-kieran-rails-reviewer`, `ce-kieran-typescript-reviewer`, `ce-kieran-python-reviewer`, `ce-dhh-rails-reviewer`, and `ce-data-integrity-guardian` from do-work's review actions based on REQ `domain` or changed-file patterns.
- **ce-plan → Route C planning** — for complex REQs, delegate the planning step to `ce-plan`. The plan artifact would live at `docs/plans/YYYY-MM-DD-<slug>-plan.md` and work.md would implement unit-by-unit from there.
- **ce-brainstorm → capture** — let users run `ce-brainstorm` first, then have do-work's capture action ingest the resulting `docs/brainstorms/*.md` into a REQ with full context preserved.
- **ce-demo-reel → present-work** — chain demo GIF/video generation into do-work's client-facing deliverables.
- **ce-polish-beta → review-work UI mode** — hand UI-feature reviews off to CE's polish workflow for dev-server + browser-reachability verification.

These phases will land as separate versions. Each one will reuse the same two primitives the ce-compound handoff established: **do-work pilots, CE specializes** and **degrade gracefully when CE is absent**.

## Installing compound-engineering

The handoff is harmless without CE — it prints a saved prompt and records `ce_compound_status: pending`. To actually promote solutions, install CE.

Installation depends on your harness. The canonical install pointers live at the [CE plugin repo](https://github.com/EveryInc/compound-engineering-plugin):

- **Claude Code:** `/plugin install compound-engineering@compound-engineering-plugin`
- **Codex, Cursor, Gemini, Kiro, OpenCode, Pi:** use the `@every-env/compound-plugin` Bun converter documented in the CE repo's README. It installs skills and agents into the target harness's native plugin/skill directories.

After install, confirm CE is available by running `/ce-compound` in your harness (it will prompt for a solution payload). If that works, the next do-work REQ you complete will surface the handoff with the "run now" option enabled.

## Troubleshooting

- **"The handoff asked but nothing happened when I chose 'run now'."** Your harness may not have CE installed, or the slash-command dispatcher isn't routing `/ce-compound`. Check that the CE plugin is enabled and that `.claude/skills/ce-compound/SKILL.md` (or the equivalent path for your harness) exists.
- **"My REQ's `ce_compound_status` is missing."** The handoff only runs on REQs with a non-empty `## Lessons Learned` section. Route A REQs without exploration or gotchas skip the handoff entirely and leave the fields unset — this is intentional.
- **"I declined the handoff but want to promote now."** Invoke CE directly: `/ce-compound` with the REQ path, or copy the saved payload from when the handoff first ran. The `declined` status prevents do-work from re-offering on its own, but you can always run CE yourself.
- **"The payload has `[TBD]` in important fields."** do-work couldn't extract that field from the REQ structure (often because the REQ is a legacy shape without `## Implementation Summary` or `## What`). Either fill the field in manually before dispatching, or let CE prompt for it interactively.

## Design principles (for contributors adding new integration points)

When wiring up the next integration point (e.g., reviewer agents, ce-plan), follow the pattern the ce-compound handoff established:

1. **do-work pilots.** do-work assembles context and decides *whether* to hand off. CE specializes and executes. do-work never embeds CE logic; it only offers the dispatch.
2. **User consent required.** Never auto-dispatch a CE skill without asking. Unattended runs default to `pending`, not `promoted`.
3. **Degrade gracefully.** If CE is not installed, the action must still complete. The handoff becomes a saved prompt, not an error.
4. **Record the outcome on the REQ.** Use a `ce_*` frontmatter field namespaced to the integration (e.g., `ce_compound_status`, `ce_plan_path`, `ce_review_findings`). Legacy REQs without the field are fine; never backfill retroactively.
5. **One reusable reference file per integration.** Like `actions/ce-compound-handoff.md`, keep the handoff logic in one place and let calling actions point to it. Avoid duplicating the instructions in both review-work and work.
