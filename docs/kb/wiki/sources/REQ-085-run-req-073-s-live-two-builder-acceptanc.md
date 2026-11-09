---
title: "Lessons from REQ-085: Run REQ-073's live two-builder acceptance test and record what it found"
type: source-summary
topic_cluster: worktree-and-parallel-dispatch
sources: [raw/processed/2026-09-01/REQ-085-run-req-073-s-live-two-builder-acceptanc.md]
related:
  - page: REQ-073-fan-out-dispatch-n-concurrent-builders-u
    rel: complements
  - page: REQ-082-the-fan-out-hand-back-file-has-no-legal
    rel: complements
  - page: REQ-091-the-hand-back-merge-fails-while-the-owne
    rel: complements
  - page: REQ-092-actions-work-md-has-no-wave-selection-or
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-085: Run REQ-073's live two-builder acceptance test and record what it found

Part of the [[concept-worktree-isolation-and-parallelism]] cluster.

## What the REQ was about

REQ-073 raised worktree dispatch from one builder to N and shipped as `completed` at v0.166.0. Its
`## Red-Green Proof` GREEN condition includes a live run of two concurrent builders; that run has never
happened. Two consecutive session checkpoints now carry it as deferred. Everything built since has been
serial, so grep proves the prose and nothing proves two builders compose.

## Solution summary

Ran REQ-073's live two-builder acceptance test for the first time. Two real queue REQs were built concurrently in two git worktrees on two branches, integrated serially by one owner, and archived; then a deliberately overlapping pair was run to confirm the merge refuses. Two defects were found and filed rather than fixed.

## What worked

- **Running it found things reading it could not.** F-01 is a hard `exit 2` on the first merge of the run, in a procedure that had been reviewed, contract-asserted, and shipped. It was invisible to grep because neither half is wrong on its own — Step 2's claim is correct, Step 6's merge is correct, and only their ordering against a tracked `do-work/` fails. This is the second time in three REQs that a documented-but-unexecuted path turned out to be broken (REQ-082 was the first).
- **Using real queue REQs rather than synthetic ones made check 2 meaningful.** Two throwaway builders would have exercised the merges just as well, but the changelog/version check only tests anything when there are two real deliverables to write entries for — and writing them is where the serial-only rule actually bites.
- **REQ-084's probe, shipped an hour earlier, was the instrument for check 3.** `git diff --name-only <integration>...<branch> -- do-work/` answered "did this builder write queue state" directly, including the committed case. Check 3 would otherwise have been a porcelain glance that could not see a committed violation.

## What didn't work

- **The `-D` slip on the throwaway branches.** After `reset --hard`, `-d` would have refused and the correct move was to report that refusal. Reaching for `-D` on branches that felt disposable is exactly the reflex the constraint exists to prevent — and "these were only throwaways" is the rationalization that makes it feel fine. Recorded in the run above rather than omitted.

## Worth knowing

- **`<pre>` is per REQ, not per run, and this run shows why.** The two ranges are `306c1f4..a17e6af` and `3ccbf36..5cfe1b5` — different lower bounds, because REQ-086's changelog and archive commits landed between the merges. A single run-level `<pre>` would have swept REQ-086's bookkeeping into REQ-087's range and misattributed it.
- **A mid-build scope extension can silently un-non-overlap a pair.** Builder B added two files to its own `write_set` during the build (REQ-087 D-02/D-03). It stayed disjoint from Builder A here, but nothing checked that — the pair was validated at pick time and never re-validated. Under real concurrency the merge would catch a textual collision, but the joint-wrongness case it cannot see is precisely REQ-073's unexercised line-proximity scenario.
- **The run directory earned its keep for an unexpected reason.** F-01's workaround was to commit the claim bookkeeping, and because `do-work/runs/` is a committable path by design, the manifest and both hand-backs went into that same commit and survived as durable state rather than as untracked scratch.

## Back-reference

See `do-work/archive/UR-016/REQ-085-run-the-live-two-builder-acceptance-test.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `b224e8a`.
