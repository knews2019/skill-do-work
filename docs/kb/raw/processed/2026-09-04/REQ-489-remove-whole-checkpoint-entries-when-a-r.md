---
source_type: req_lesson
req_id: REQ-489
req_path: do-work/archive/REQ-489-remove-whole-checkpoint-entries-on-departure.md
date: 2026-09-02
domain: backend
module: skills/do-work/tools/do-work-cli
tags: [backend, remove, whole, checkpoint]
---

# Lessons from REQ-489: Remove whole checkpoint entries when a REQ leaves working

## What the REQ was about

When the canonical `complete` (and, by the same code path, `fail` and the other departures from `do-work/working/`) removes this checkout's entry from `do-work/CHECKPOINT.md`'s `## In Progress (interrupted)` list, it deletes only the `- REQ-NNN: ...` header line. The indented `Last known state:`, `Key files being modified:`, and `Known issues:` lines that Step 10 enrichment adds beneath it stay behind as an unattributed orphan block.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)

## What worked

- Exact-section RED cases exposed both the orphaned continuation block and the inline-heading attraction bug with a small two-file diff.

## What didn't work

- Treating the canonical request-state helper as the only departure writer missed cleanup's independent checkpoint-removal implementation.

## Worth knowing

- A stored-format contract is not closed until every alternate writer of that format is swept; package-local green tests can still leave a repository-wide lifecycle path stale.

## Back-reference

See `do-work/archive/REQ-489-remove-whole-checkpoint-entries-on-departure.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6e92e536`.
