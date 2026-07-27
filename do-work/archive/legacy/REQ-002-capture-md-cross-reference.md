---
id: REQ-002
title: "Code review: fix bare-name cross-reference in actions/capture.md:63"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:44:30Z
completed_at: 2026-05-29T18:44:45Z
commit: 2952f24
route: A
review_generated: true
source: code-review
scope: actions/capture.md
---

# Code Review Fix: Fix Bare-Name Cross-Reference in capture.md

## What

`actions/capture.md:63` says `` See `review-work.md`. `` (bare name) immediately after the same paragraph correctly uses `` `actions/review-work.md` may append … ``. This violates the explicit CLAUDE.md rule: *"Cross-reference other actions by their file path (e.g., `actions/work.md`) — not by name."*

This is the only bare-name reference in the file. 47 path-style references dominate across the codebase.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). Single-line fix.

## Requirements

- Change `` See `review-work.md`. `` to `` See `actions/review-work.md`. `` at `actions/capture.md:63`.

## Acceptance

- `grep -nE '\breview-work\.md\b' actions/capture.md` shows only path-prefixed references.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `consistency.md` F2

---

## Triage

**Route: A** — Simple

**Reasoning:** One-line text replacement at a specific file:line. No exploration or planning needed.

**Planning:** Not required.

## Implementation Summary

**Files changed:**
- `actions/capture.md` (modified) — L63: `` See `review-work.md`. `` → `` See `actions/review-work.md`. ``

**What was done:** Replaced the single bare-name cross-reference with a path-prefixed one, aligning with the CLAUDE.md cross-reference rule and the rest of the file's 47 path-style references.

## Qualification

Passed — 1 file verified via `git diff`, single requirement met. Acceptance grep (`grep -nE '\breview-work\.md\b' actions/capture.md`) shows the lone remaining reference is the path-prefixed one.

## Testing

**Tests run:** `grep -nE '\breview-work\.md\b' actions/capture.md`
**Result:** ✓ Only path-prefixed reference present.

*Verified by work action*

## Review

**Overall: 100%** | 2026-05-29T18:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor
**Acceptance:** Pass — exact-text replacement matches the acceptance grep.
**Follow-ups created:** None

*Reviewed by work action (Route A self-review)*

