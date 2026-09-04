---
source_type: req_lesson
req_id: REQ-340
req_path: do-work/archive/UR-065/REQ-340-finish-the-report-image-interruption-sweep.md
date: 2026-08-24
domain: testing
module: _dev/primes
tags: [testing, addendum, finish, report]
---

# Lessons from REQ-340: Addendum: finish the report-image interruption sweep

## What the REQ was about

REQ-325 closed two interruption defects in the per-image helper. The prime's rule is to grep the
same primitive across every caller before calling the class closed, and the batch that drives that
helper carries both of them. A third, smaller instance sits in the helper REQ-325 did fix.

## Solution summary

**Files changed:**
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (modified)
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (modified)

## What worked

- ` before starting: it
- records which bash semantics were confirmed by probe (a trap does fire from inside `wait`; a group
- kill does reach the backend; a direct child is always reapable by KILL) and which fix shapes
- worked. Re-deriving those cost most of that REQ's time.
- The launch-window instance could not be pinned by a fixture in REQ-325 — the parent won the race
- 160/160 times. Do not ship a stress case that has never failed; that REQ's D-02 states why.

## Back-reference

See `do-work/archive/UR-065/REQ-340-finish-the-report-image-interruption-sweep.md` for the full REQ — plan, exploration, implementation, review, and lessons.
