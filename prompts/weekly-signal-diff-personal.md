# Weekly Signal Diff — Personal Sidecar (placeholder)

> Placeholder template for the `weekly-signal-diff` personal sidecar. Declares no real lanes on its own. Create a project-local copy anywhere in your repo to add user-specific lanes to the weekly scan.

**Aliases:** (none — this is not invoked directly)

**When to use:**
- Auto-discovered by `do work prompts run weekly-signal-diff`; not run on its own
- The skill ships this file as a placeholder. To personalize the scan, create a file named `weekly-signal-diff-personal.md` anywhere in your project (project root, `.claude/`, `do-work/`, wherever fits) and declare your real lanes there.
- Your project-local copy overrides this placeholder at Phase 3 of the main prompt.

**Inputs / flags:**
- None. Data sidecar only.

---

## What this file is

This file is a placeholder shipped with the `do-work` skill. It intentionally contains no real lanes — only an illustrative template row. The library `weekly-signal-diff` prompt searches your project for a file named `weekly-signal-diff-personal.md` at Phase 3. If it finds one, it loads the personal lanes declared there. If no project-local copy exists, the weekly scan runs with the 10 core lanes only and this placeholder contributes nothing.

## How to activate

1. Create or copy a `weekly-signal-diff-personal.md` anywhere in your project, for example:
   - Project root: `./weekly-signal-diff-personal.md`
   - Claude config: `./.claude/weekly-signal-diff-personal.md`
   - Do-work directory: `./do-work/weekly-signal-diff-personal.md`
2. Fill the table below with lanes tied to your active projects, toolchains, and concerns. Replace the placeholder row.
3. Run `do work prompts run weekly-signal-diff` — the prompt will auto-discover your copy and load the lanes.

## Lane table shape

Your personal lanes table must have four columns (Number, Category, Suggested entities, Why this lane matters). One placeholder row is shown below; replace it with real lanes or delete it entirely.

| # | Category | Suggested entities | Why this lane matters |
|---|---|---|---|
| 11 | [Your vertical — e.g., billing platforms, commerce tooling, supply-chain security] | [Vendors and tools you work with in this vertical] | [Why a shift in this lane would affect your active work or toolchain] |

## Rules for personal lanes

- Each row declares one lane that will be scanned every week alongside the 10 core lanes.
- Numbering (`#` column) is cosmetic — contiguous numbering from 11 is conventional, but the main prompt does not enforce it.
- Personal lanes get the same full-coverage treatment as core lanes: a per-lane scan-notes paragraph every week, weighted at least as heavily as any core lane.
- Swap entities or remove lanes as your active projects shift.

## Adding, removing, or editing lanes in your project-local copy

- **Add a lane:** append a new row to the table.
- **Remove a lane:** delete the row. The main prompt recalculates the total (`N / N lanes covered`) automatically.
- **Swap entities within a lane:** edit the Suggested entities column. The main prompt uses BKB context to promote whichever entities are most relevant this week.
