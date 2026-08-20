---
session_ended: 2026-08-19T20:52:00Z
last_completed: REQ-282
queue_state: 26 pending, 2 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 4
session_depth: moderate
---

# Session Checkpoint

## In Progress (interrupted)


## Completed This Session

`do-work run REQ-268 REQ-270 REQ-276 REQ-282` — four independent correctness REQs, serial mode, each with its own adversarial review:

- REQ-268: never report clean for a state that was never verified (Route B, 91%) — `ec5e550`, **0.215.1**
- REQ-270: carry a worktree builder's hand-back sections into Step 8 (Route B, 58% Partial) — `9fe63fa`, **0.215.2**
- REQ-276: give record-commit-hash's readers the writer's fence guard (Route A, 94%) — `fa19f69`, **0.215.3**
- REQ-282: make the release probes run in a suite checkout (Route B, 98%) — `bc809fd`, **0.216.0**

Every hash confirmed with `record-commit-hash.sh --verify`; `maintainer-verify.sh` exited 0 at every commit boundary. Serial throughout — no worktrees created, none to clean up. One pre-session bookkeeping commit (`fd86625`) recorded REQ-296's abandonment, which was complete in the tree but uncommitted.

## Still Queued

**Twenty-eight** — 26 `pending`, 2 `pending-answers`. Net −2 (four shipped, two follow-ups created).

**Needs you (`pending-answers`, both generation-≥2 follow-ups):**
- **REQ-298** — the unchecked-exit-status primitive REQ-268 fixed was copied from `tools/checks/record-commit-hash.sh`, which still carries it. **Sequence after REQ-276's file work, not before.**
- **REQ-299** — `## Decisions` has REQ-270's exact defect at two read sites outside Step 8 (review-work's traceability check, and the Decision Brief's HANDLED block), so under fan-out a builder's judgment calls never reach the user.

## Session Notes

**Every one of the four reviews found something the builder's own sweep missed, and three of the four findings were in the same file the REQ had just edited.**

- REQ-268 swept for the primitive and stopped at the three instances the REQ listed; three more of the same shape sat below them, including the post-rename guard that runs *after* the file is replaced and whose `0/0` fallback made `[ 0 -gt threshold ]` false for every threshold.
- REQ-270 swept Step 8's substeps correctly and never grepped outside Step 8. An always-loaded crew file still told every builder to do the thing the fix forbids — that one had to be fixed inside the REQ, because shipping the reader-side fix alone would have shipped a fix that does not work.
- REQ-276's helper was defined below the `--verify` branch that calls it. Caught by running the suite, not by reading it.

**The pattern across all three: the sweep was for the instance across one scope, when the condition needed sweeping across the repo.** REQ-270's lesson names it directly.

**A near-miss worth keeping.** REQ-282's obvious one-line implementation — widen the shared version-file resolver — satisfies the feature and silently breaks the constraint the REQ spent a paragraph on, because that resolver is also `next-version`'s writer. The REQ named the constraint; only a test enforced it.

**Estimator ran under on three of four:** 20→33, 15→19, 10→19, 20→12 minutes. The overage is serial review latency, not build time. Four rows appended to `do-work/calibration-log.tsv`; do not recalibrate off four.

## Context Summary

**Two of tonight's four REQs were about the same thing wearing different clothes: a check that answers "fine" for work it never did.** REQ-268 was that condition in the timestamp scripts, REQ-282 was three release probes that had been silently off since the four-skill split — including the duplicate-version-number check that CLAUDE.md says has caught real failures before. Both are now loud. REQ-298 is the remaining reach of the first, and it is the one queued item that touches a file the pipeline depends on every single run.

**Prove a check bites, never that it went quiet.** REQ-282's acceptance needed a deliberate `9.9.9` mismatch on the real repo, because a disabled probe and a working one print the same thing on a clean tree. That test shape is now in `release_test.go` and is worth copying wherever a probe is added.
