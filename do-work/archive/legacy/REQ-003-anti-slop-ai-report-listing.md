---
id: REQ-003
title: "Code review: add ai-report to anti-slop loading lists in CLAUDE.md and anti-slop.md"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:47:00Z
completed_at: 2026-05-29T18:47:30Z
commit: bad903e
route: A
review_generated: true
source: code-review
scope: CLAUDE.md, crew-members/anti-slop.md
---

# Code Review Fix: Add ai-report to anti-slop Loading Lists

## What

CLAUDE.md line 159 enumerates exactly 5 places `crew-members/anti-slop.md` is loaded: present-work (Step 4), review-work (Step 9), pipeline (Step 5), kb-lessons-handoff (Step 2), and slop-check. But `actions/ai-report.md` lines 10 and 41 **also** explicitly load anti-slop ("Anti-slop or it doesn't ship" + "Read `crew-members/anti-slop.md`. Keep all seven principles active").

The JIT_CONTEXT comment at `crew-members/anti-slop.md:3` has the same gap.

This is the closest thing the skill has to a single-source-of-truth for "where does this crew-member load." Any future addition of an anti-slop consumer will look at one of these two lists, see the gap, and have to guess whether it's intentional.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). Documentation drift between the listed callers and the actual callers.

## Requirements

- Add `ai-report (Step 1 principle loading; applied inline to every section per Step 6)` to the anti-slop bullet in CLAUDE.md's Agent Rules section (line 159 area).
- Add the same entry to the JIT_CONTEXT comment at `crew-members/anti-slop.md:3`.

## Acceptance

- Both lists enumerate the same 6 callers: present-work, review-work, pipeline, kb-lessons-handoff, slop-check, ai-report.
- `grep -l 'anti-slop' actions/*.md` shows no caller missing from either list.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `architecture.md` F1

---

## Triage

**Route: A** — Simple

**Reasoning:** Two enumerated lists, exact addition known, two files named. No exploration needed beyond the acceptance grep.

**Planning:** Not required.

## Implementation Summary

**Files changed:**
- `CLAUDE.md` (modified) — L159: added `ai-report (Step 1 principle loading; applied inline to every section per Step 6)` to the anti-slop callers list.
- `crew-members/anti-slop.md` (modified) — L3 JIT_CONTEXT: added `ai-report's section drafting (Step 1 principle load + applied inline through Step 6)`.

**What was done:** Closed the documentation drift between the anti-slop loading list and the actual callers. Both source-of-truth lists now enumerate the same six callers: present-work, review-work, pipeline, kb-lessons-handoff, ai-report, slop-check.

## Qualification

Passed — 2 files modified per scope. `grep -l 'anti-slop' actions/*.md` returns exactly the six caller files now enumerated in both lists. No other anti-slop callers exist.

## Testing

**Tests run:** `grep -l 'anti-slop' actions/*.md`
**Result:** ✓ 6 callers: ai-report, kb-lessons-handoff, pipeline, present-work, review-work, slop-check. All six appear in both lists.

*Verified by work action*

## Review

**Overall: 100%** | 2026-05-29T18:47Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor
**Acceptance:** Pass — both lists enumerate the same six callers; grep is clean.
**Follow-ups created:** None

*Reviewed by work action (Route A self-review)*

