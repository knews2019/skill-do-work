---
session_ended: 2026-08-20T08:41:00Z
last_completed: REQ-258
queue_state: 29 pending, 2 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 1
session_depth: light
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

`do-work run` — stopped by the user after one REQ, at the commit boundary, for a fresh-session handoff.

- **REQ-258**: Split the prescribed shell behavior suite per script (Route B, 94% / Acceptance Pass) — `1cc1836`, **0.216.2**

Hash confirmed with `record-commit-hash.sh --verify`; `maintainer-verify.sh` exited 0 immediately before the commit. Serial mode — no worktrees created, none to clean up.

**A second session was writing this repo concurrently.** Two commits landed mid-run from outside this session: `031c546` (clarify approving REQ-298 and REQ-299 to `pending`) and `1311300` (duration-label rounding, version 0.216.0 → 0.216.1). Neither touched any file REQ-258 wrote, verified by `git show --stat` on both. `031c546` did commit this session's in-flight `git mv` of REQ-258 into `working/` as an R100 rename, which was harmless — content was identical at that instant — and this session's archive move landed on top of it correctly.

## Still Queued

**Thirty-one** — 29 `pending`, 2 `pending-answers`. Net +3 (one shipped, three follow-ups created, two flipped from `pending-answers` to `pending` by the parallel clarify).

**Needs you (`pending-answers`, both REQ-258 discoveries awaiting the consent flow):**
- **REQ-301** — `tools/checks/qualify.sh` has no rename/copy detection, so every code-relocation REQ gets a false `[UNIFY]` FAIL on pre-existing debug markers inside the moved text. REQ-258 hit it on four fixture `TODO` strings. The risk is habituation: a gate that cries wolf on a whole category of change trains builders to wave it away.
- **REQ-302** — `effort_estimate: trivial` produced a 5-minute P50 for REQ-258's 19-file restructure. One data point; the REQ asks the question before proposing a fix.

**Queued at `pending` and worth sequencing deliberately:**
- **REQ-300** — the text that still plans around the pre-split shell suite: `RESTART-PROMPT.md:33` and the `write_set` of REQ-263, REQ-264, REQ-271. **Run it before those three**, or their board display and any wave planning read the dissolved monolith.

## Session Notes

**REQ-258's only real risk was unverifiable success.** A 1882-line file split 17 ways can claim "no case changes" and be wrong in a way no test catches, because the suite tests the *scripts*, not itself. The thing that made the claim checkable was a line-multiset comparison of every non-blank case line, original against split: 1756 in, 1756 out. That, plus a deliberate mutation proving the runner still exits 1, is the whole evidence base — the green run alone proves nothing, because a suite that silently ran zero cases prints the same thing.

**The split surfaced two defects that were invisible while it was one file:** `generate-report-image`'s cases were interleaved around its `-batch` sibling's, `repair-req-timestamps`' around `audit-archive-timestamps`', and a `publish-portfolio-summary` fixture setup block sat under the `install-memory-hooks` header entirely. None of that was findable by reading 1882 lines top to bottom.

**Estimator ran badly under: 5 → 51 minutes.** Not an estimator fault — `effort_estimate: trivial` short-circuits signal extraction by design, and the field was misjudged at capture. That is REQ-302. One row appended to `do-work/calibration-log.tsv`; do not recalibrate off it.
