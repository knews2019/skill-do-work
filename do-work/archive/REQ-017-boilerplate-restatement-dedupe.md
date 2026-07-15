---
id: REQ-017
title: "Dedupe intra-file guard restatements in note, scan-ideas, commit, quick-wins"
status: completed
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
claimed_at: 2026-07-15T18:05:00Z
completed_at: 2026-07-15T18:12:00Z
route: A
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
- [x] Every deleted row provably maps to a surviving Rule (record the mapping in
      the Implementation Summary).
- [x] No unique constraint deleted; word deltas recorded per file.

## Open Questions
(none)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read all four files' guard sections in full; built the deletion→surviving-rule mapping before cutting anything (maintenance.md: subtraction is not vandalism).
- [x] **[APPLY]:** note.md CR+RF+VC removed (each row maps to Rules: not-a-task, append-only, bullets-only, never-commits, empty-input); scan-ideas.md RF+VC removed (map to Philosophy/Rules/Output: grounded, focus, product-breadth, no-dupes, read-only, 8-15 cap); commit.md Checklist + Common-mistakes blocks removed (.env was stated 7x; now in What-NOT-Do + Red Flags + Verification only) and 3 generic CR rows dropped; quick-wins.md 2 generic CR rows dropped (long-file row = generic wisdom; padding row = restates Rules "Be honest about impact").
- [x] **[UNIFY]:** Word deltas: note 1,027→800, scan-ideas 1,027→856, commit 1,904→1,527, quick-wins 1,844→1,801 (net −818). Kept intact: commit.md terminal-status Red Flag + REQ-traceability CR rows; quick-wins scan-breadth + dynamic-refs rows and the maintenance-marker Rules paragraph. Out-of-scope files untouched.

## Triage

Route A — surgical deletions with a pre-built mapping; no exploration.

## Implementation Summary

**What was done:** Removed intra-file guard restatements in 4 small actions; every deleted row maps to a surviving rule.

Files changed:
- `actions/note.md` (modified) — CR+RF+VC sections → one-line pointer at Rules. Mapping: auto-promote→"A note is not a task"; tidy/sort→"Append-only"; empty→"Empty input is a no-op"; UR/REQ-created→"never let a note kick off capture"; frontmatter/header→"Write bullets, never frontmatter"; ran-commit→"The action never commits".
- `actions/scan-ideas.md` (modified) — RF+VC sections → one-line pointer. Mapping: no-evidence→Rules "Grounded"; off-focus→"Respect the focus"; all-refactors→Philosophy "Product thinking"; dupes→"No duplicates"; created-REQ→"Read-only"; 30+ ideas→Output "More than 15 dilutes".
- `actions/commit.md` (modified) — Checklist (restated Steps 1-6 headings) and Common-mistakes (all 6 rows restated: add -A→Steps/VC; --no-verify→Error Handling; .env→What-NOT-Do/RF/VC; giant commit→Steps; unrelated grouping→CR; exclusion check→What-NOT-Do) removed; CR rows 1-2-4 (generic git advice; .env kept in 3 other homes) dropped.
- `actions/quick-wins.md` (modified) — CR rows 1 (generic) and 4 (restates Rules "Be honest about impact") dropped.
- `actions/version.md`, `CHANGELOG.md` — version 0.123.2 + entry.

## Testing

- Non-behavioral (instruction dedupe): regression evidence = the mapping above; `bash _dev/tests/contract-regressions.sh` re-run to confirm no asserted pattern was in the deleted blocks.

## Lessons Learned

**What worked:** Building the deletion→survivor mapping BEFORE editing made "subtraction, not vandalism" checkable.
**What didn't:** The original hypothesis (strip boilerplate sections from all small actions) was wrong — the audit found the content hard-won; only the restatements were bloat.
**Worth knowing:** commit.md's .env rule appeared in SEVEN places; when a rule feels over-restated, count its homes before assuming one deletion suffices.
