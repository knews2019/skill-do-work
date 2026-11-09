---
source_type: req_lesson
req_id: REQ-001
req_path: do-work/archive/legacy/REQ-001-split-work-md-orchestrator.md
date: 2026-06-01
module: actions
tags: [actions, orchestrator, reference, companion]
---

# Lessons from REQ-001: Code review: split actions/work.md into orchestrator + reference companion

## What the REQ was about

`actions/work.md` has grown to 1,074 lines — the longest file in the repo, on the hottest path (read on every `do-work run` and every `pipeline` build step). The body has accreted multiple sub-systems that have natural seams (Schema Read Contract, wave execution, archive-collision rules, plan validation, scope drift, qualification, TDD verify, prime-link deferral, follow-up REQ template, commit-format prose).

## Solution summary

A faithful *move-not-rewrite* split, following the `bkb.md`/`bkb-reference.md` pattern. A single deterministic Python transform excised blocks by original line range and substituted pointers; the companion was assembled from the same ranges. Verified line-by-line that every substantive original line survives in one of the two files.

## What worked

- Modeling the split as two range-lists over the *original* line numbers — "work.md excisions (→ pointer)" and "companion blocks (→ header)" — let one deterministic transform do the move with no line-shift bugs. A line-by-line content-preservation diff (every substantive original line must reappear in new+companion) is the right "test" for a move-not-rewrite refactor: it proved nothing was silently dropped and surfaced exactly which deltas were intentional.

## What didn't work

- Taking the REQ's enumerated move-list literally would have shipped a ~910-line `work.md` that *fails its own `<700` acceptance* — the Requirements and Acceptance were inconsistent. Resolved by D-01 (move the enumerated sections plus enough additional reference-grade content to hit the target).

## Worth knowing

- `work.md`'s document-level Red Flags + Verification Checklist had drifted mid-document inside Step 6, orphaning the Step-6.3 qualify-fail/pass logic beneath them; relocating the headings to the tail reconnected it. The Schema Read Contract is referenced by-name (no path) throughout `work.md` — those internal refs stay stable after the move; only cross-file references needed a path update.

## Back-reference

See `do-work/archive/legacy/REQ-001-split-work-md-orchestrator.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `3d600f6`.
