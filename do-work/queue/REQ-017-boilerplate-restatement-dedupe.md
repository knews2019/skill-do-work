---
id: REQ-017
title: "Dedupe intra-file guard restatements in note, scan-ideas, commit, quick-wins"
status: pending
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: true
---

# Boilerplate restatement dedupe (4 files)

## What
Remove only guard content that restates the same file's Rules 1:1; keep every
unique hard-won row:

- `actions/note.md` — Common Rationalizations + Red Flags rows that map onto
  Rules (append-only, not-a-task, empty-input, never-commits). Keep the Rules.
- `actions/scan-ideas.md` — Red Flags + Verification rows restating Rules
  (grounded/no-generic/read-only/focus/8-15 count). Keep the Rules.
- `actions/commit.md` — collapse the FIVE overlapping guard sections
  (Checklist + Common-mistakes + CR + RF + VC) to the standard triad; drop
  near-generic git-advice rows; keep file-specific rows (terminal-status
  filtering, completed-with-issues association).
- `actions/quick-wins.md` — drop the generic CR rows (generic refactoring
  wisdom); keep the section and the maintenance-marker contract.

**Explicitly out of scope:** abandon, board, clarify, slop-check, forensics,
verify-requests, validate-feedback — their guard sections were verified hard-won
(audit §1d) and stay untouched.

## Why
The audit falsified the "boilerplate is filler" hypothesis: content is real, the
defect is triple-stating within a file. Audit DELETE bucket (revised scope).

## Acceptance criteria
- [ ] Every deleted row provably maps to a surviving Rule (record the mapping in
      the Implementation Summary).
- [ ] No unique constraint deleted; word deltas recorded per file.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:**
- [ ] **[APPLY]:**
- [ ] **[UNIFY]:**
