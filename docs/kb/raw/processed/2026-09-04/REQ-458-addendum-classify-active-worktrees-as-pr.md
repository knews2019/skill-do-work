---
source_type: req_lesson
req_id: REQ-458
req_path: do-work/archive/UR-086/REQ-458-classify-active-worktrees-as-present.md
date: 2026-09-03
domain: backend
module: _dev/primes
tags: [backend, addendum, classify, active]
---

# Lessons from REQ-458: Addendum: classify active worktrees as present and non-fixable

## What the REQ was about

Correct REQ-083 (Verify reports every builder worktree as a fixable orphan, including active and unmerged ones) so a branch being merged into the integration branch is not, by itself, enough to call its worktree a leftover or mechanically fixable. A dirty worktree or a worktree belonging to an unfinished run must be reported as present and non-fixable; only clean merged residue from finished work may be reported as a fixable leftover.

## Solution summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

## What worked

- Reaching for the tri-state the file already used next door. `worktreeMergeState` had exactly the right shape — merged / unmerged / *git declined to answer* — and the fix for F1 was to stop answering "is this finished?" with a boolean and copy that shape. REQ-457 landed on the same answer independently in the same week, which is what made it worth a family slug rather than a one-off note.
- Running the built binary against a purpose-made fixture instead of trusting the unit tests. Both F1 repros and the M1 index write were found that way; `TestVerifyWritesNothing` passes and cannot catch M1, because its fixture is not a git repo so the probe never runs there.

## What didn't work

- The first pass fixed the two cases the user reported and stopped. A third input class — a REQ id the board never saw — still reached `Fixable: true` while its remedy asserted the REQ had left `working/`, which nothing had established. Closing the reported instances is not closing the defect.
- Declaring `forensics.md` in Scope before reading it. It turned out to have nothing to update by design, so the declaration became drift that had to be explained rather than an edit that had to be made.
- Adding a probe without asking whether the probe itself mutates. `git status` refreshes the index it reads; five of five worktree indices changed hash before `--no-optional-locks` went in, against the tool's own written promise that verify writes nothing at all.

## Worth knowing

- `git status` is not read-only, and the flag that makes it so (`--no-optional-locks`) is a top-level option — it must precede `-C`, not follow the subcommand.
- A REQ file parked outside the scanned sections is recorded as `StrayRequestFiles` and never enters `board.RequestsById`, so its `status:` is unreadable through the board no matter what it says. Any board lookup keyed on presence has to treat absence as unknown, not as a negative.
- A fixture that plants no REQ file is not testing finishedness — it is riding on whatever the absent case happens to do. REQ-083's fixture asserted two fixable leftovers exactly that way, which is why it had to change here.

## Back-reference

See `do-work/archive/UR-086/REQ-458-classify-active-worktrees-as-present.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `ea2bab0`.
