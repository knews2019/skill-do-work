---
id: REQ-005
title: "Code review: fix SKILL.md routing-table priority off-by-one cross-references"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:40:21Z
completed_at: 2026-05-29T18:43:00Z
commit: c58c367
route: A
review_generated: true
source: code-review
scope: SKILL.md
---

# Code Review Fix: Fix SKILL.md Routing-Table Priority Cross-References

## What

SKILL.md's routing table at L66-98 lists Verify keywords at **priority 5** (line 72). But the Verb Reference rows for `ui-review` (line 128) and `review-work` (line 129) both refer to verify as **"priority 4"**. The `slop-check` row (line 149) correctly says priority 5.

A reader cross-checking the carve-out ("Do NOT use 'check ui' — consumed by verify at priority 4") and then trying to find priority 4 ends up at "Action verbs (± REQ IDs ± flags)" → work, which is the wrong route. The actual dispatch logic is unaffected, but this is a real abstraction-health concern in the project's central dispatch model.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). Nothing automated checks routing-row references; the typo survives because no diff-time validator exists.

## Requirements

- Change "priority 4" to "priority 5" in `SKILL.md:128` (ui-review row).
- Change "priority 4" to "priority 5" in `SKILL.md:129` (review-work row).
- Grep the entire `SKILL.md` for any other `priority \d+` references and audit each against the actual routing-table priorities.

## Acceptance

- All numeric priority references in SKILL.md resolve to the correct row in the routing table.
- `grep -nE 'priority [0-9]+' SKILL.md` is hand-audited and clean.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `architecture.md` F3

---

## Triage

**Route: A** — Simple

**Reasoning:** Names a specific file (`SKILL.md`), with two specific line numbers and exact-text changes. The grep audit is a one-command check. No exploration or planning needed.

**Planning:** Not required.

## Implementation Summary

**Files changed:**
- `SKILL.md` (modified) — L128 and L129: "priority 4" → "priority 5" in ui-review and review-work Verb Reference rows.

**What was done:** Edited the two off-by-one cross-references in the Verb Reference table so they point at the correct row in the routing table (verify keywords are priority 5, not 4). Ran the prescribed `grep -nE 'priority [0-9]+' SKILL.md` audit before and after; the remaining six priority references (2, 7, 11×3, 28×2, 29×2, and the corrected 5 on L149) all resolve correctly against the routing-table rows.

## Qualification

Passed — 1 file verified via `git diff`, 1 requirement traced (Acceptance criterion: grep is clean). No P-A-U state on Route A. Edit is substantive (corrects two distinct semantic errors), no debug artifacts, no scope drift.

## Testing

**Tests run:** `grep -nE 'priority [0-9]+' SKILL.md` (no automated test suite for routing-prose; the grep audit is the prescribed acceptance check)
**Result:** ✓ All eight matches hand-audited against the routing table (L66-98) — every reference now points at the correct row.

**Red-green validation:** Non-behavioral text edit; no test-first evidence required. Regression evidence is the grep audit above.

*Verified by work action*

## Review

**Overall: 95%** | 2026-05-29T18:42Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor
**Acceptance:** Pass — both off-by-ones corrected; full grep audit clean.
**Suggested testing:** None — there's no diff-time validator that could catch this class of error; that's a separate, larger problem the REQ correctly punts on.
**Follow-ups created:** None

*Reviewed by work action (Route A self-review)*

