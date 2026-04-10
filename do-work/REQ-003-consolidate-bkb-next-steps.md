---
id: REQ-003
title: Consolidate bkb next-steps into next-steps.md
status: pending
created_at: 2026-04-10T19:30:00Z
user_request: UR-001
domain: general
prime_files: []
tdd: false
---

# Consolidate bkb next-steps into next-steps.md

## What
Move the 86-line per-sub-command "Next steps" section from the end of `actions/build-knowledge-base.md` (lines ~1080-1165) into `next-steps.md`, which is the canonical location for post-action next-step suggestions. Remove the section from the action file and replace with a reference to `next-steps.md` if needed. The existing brief bkb entry in `next-steps.md` (around line 135) should be expanded with the per-sub-command detail.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context
`next-steps.md` is referenced by `SKILL.md` as the canonical source for post-action next-step suggestions. `build-knowledge-base.md` embeds its own next-steps section (12 blocks, one per sub-command), creating two sources of truth. The existing `next-steps.md` entry for bkb is a 3-line generic suggestion — the per-sub-command detail from the action file is more useful.

## Red-Green Proof
**RED prompt/case:** `build-knowledge-base.md` contains ~86 lines of "Next steps:" blocks (lines 1080-1165) that duplicate the role of `next-steps.md`.
**Why RED now:** Two sources of truth for bkb next-steps — one in the action file, one in `next-steps.md`.
**GREEN when:** Single source of truth in `next-steps.md` with per-sub-command bkb next-steps; `build-knowledge-base.md` no longer contains its own next-steps section.
**Validation:** Inferred during capture

---
*Source: Address quick-wins report findings*
