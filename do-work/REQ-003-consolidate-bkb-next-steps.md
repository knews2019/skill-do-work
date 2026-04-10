---
id: REQ-003
title: Consolidate bkb next-steps into next-steps.md
status: done
completed_at: 2026-04-10T19:34:00Z
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
- [x] **[PLAN]:** Move 86 lines of per-sub-command next-steps from build-knowledge-base.md into next-steps.md, replacing the generic 3-line bkb entry. Remove the section from the action file.
- [x] **[APPLY]:** Replaced the generic "After build knowledge base" entry in next-steps.md with 11 per-sub-command blocks (init, triage, ingest, query, lint, resolve, close, rollup, defrag, garden, crew). Removed the "## Next Steps" section from build-knowledge-base.md.
- [x] **[UNIFY]:** Verified both files. next-steps.md gained the detailed entries. build-knowledge-base.md lost 88 lines (the section + separator). No other files affected.

## Context
`next-steps.md` is referenced by `SKILL.md` as the canonical source for post-action next-step suggestions. `build-knowledge-base.md` embeds its own next-steps section (12 blocks, one per sub-command), creating two sources of truth. The existing `next-steps.md` entry for bkb is a 3-line generic suggestion — the per-sub-command detail from the action file is more useful.

## Red-Green Proof
**RED prompt/case:** `build-knowledge-base.md` contains ~86 lines of "Next steps:" blocks (lines 1080-1165) that duplicate the role of `next-steps.md`.
**Why RED now:** Two sources of truth for bkb next-steps — one in the action file, one in `next-steps.md`.
**GREEN when:** Single source of truth in `next-steps.md` with per-sub-command bkb next-steps; `build-knowledge-base.md` no longer contains its own next-steps section.
**Validation:** Inferred during capture

## Implementation Summary

| File | Status | What changed |
|------|--------|-------------|
| `next-steps.md` | (modified) | Replaced generic 3-line bkb entry with 11 per-sub-command next-step blocks (init, triage, ingest, query, lint, resolve, close, rollup, defrag, garden, crew). |
| `actions/build-knowledge-base.md` | (modified) | Removed "## Next Steps (shown after each sub-command)" section (88 lines). |

---
*Source: Address quick-wins report findings*
