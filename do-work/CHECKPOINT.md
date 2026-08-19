---
session_ended: 2026-08-19T09:57:20Z
last_completed: REQ-265
queue_state: 12 pending, 4 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 6
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

Six REQs shipped across two waves, 0.212.20 through 0.212.25, each with an independent adversarial review:

- REQ-259: retire the skill-root citation reading at its three unbackticked sites (Route B, 91%) — merge `4081c50`, **0.212.20**
- REQ-260: run the Go formatter as part of the canonical verify (Route A, 96%) — `307e146`, **0.212.21**
- REQ-257: decide the offset/fractional refusal — permanent, and pinned (Route B, 90%) — `6afcbd5`, **0.212.22**
- REQ-261: delete the date-only tripwire and keep the rule (Route A, 96%) — `210abba`, **0.212.23**
- REQ-267: close the two remaining repairer shape divergences (Route B, 96%) — `f7441d7`, **0.212.24**
- REQ-265: raise the mark-label line-box bound and drop its duplicate descent (Route A, 92%) — `1227678`, **0.212.25**

Every hash confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exited 0 at every commit boundary. All six worktrees and branches removed with `git worktree remove` / `git branch -d` (never `-D`); `git worktree list` shows only the main tree.

## Still Queued

**Sixteen** — twelve `pending`, four `pending-answers`. The queue **grew** this session: six shipped, ten created. Every one of the ten came from an adversarial review or an external one finding a real defect, so this is the process working rather than failing — but it means "drain the queue" is not a reachable end state. Stop when the findings stop being worth fixing.

**Needs you (`pending-answers`):**
- **REQ-274** — the "SessionStart hook exits nonzero" mechanism is false and still stated where it carries a decision's rationale
- **REQ-275** — a third repairer/board divergence on the field-name axis: the repairer keys on the `_at` suffix, the board on six hand-kept names. Latent today
- **REQ-276** — `record-commit-hash.sh` guards its writer against an unterminated fence but not its readers, on the last check every REQ passes through
- **REQ-278** — nothing bounds the Durations label face off Linux; scoped to *measuring*, with a closed font stack as the cheap option

**Ready to build (`pending`):** REQ-258, 262, 263, 264, 266, 268, 269, 270, 271, 272, 273, 277.

## Session Notes

**The scheduling bottleneck is one file.** `_dev/tests/prescribed-shell-scripts-behavior.sh` is written by REQ-258, 263, 264, 268 and 271, so exactly one of those can run per wave. It forces roughly four more waves regardless of builder count. REQ-258 restructures that file wholesale, so it should probably run alone or first.

**Corrections made on the record this session, each caught by someone other than its author:**
- REQ-267's severity framing was wrong — the hook runs the repairer under `|| true`, so a failure prints a banner line and never breaks a session. The orchestrator repeated the false claim twice before a builder read the hook. Residue tracked as REQ-274.
- REQ-265's integration seam, applied by the orchestrator, shipped a false safety claim ("the declared parts round both up" — the declared descent rounds *down*). Its reviewer caught it; fixed before shipping rather than deferred.
- REQ-261's record said "same sentence" where it was "same paragraph, three sentences earlier". Corrected in place.
- Two builders reached right conclusions via provably false arguments (REQ-257 on why offsets are refused, REQ-265 on what triggers the pitch alarm). Both were caught only because the reviewers checked the *reasoning* and not just the verdict.

**The measurement lesson, twice in one REQ.** REQ-267's fuzz first missed the wedge because it compared only whether a file was *mutated*, and the wedge *refuses*. Then both its fuzz and its reviewer's held the field name constant, so neither could see the `_at`-suffix divergence. A fuzz's blind spots are the axes it holds constant, and its oracle decides which failures are visible at all.

**Do not recalibrate the estimator from this session's rows.** Six rows appended to `do-work/calibration-log.tsv`; the wall spans measure serial integration queuing and review latency far more than they measure work — REQ-265's 15-minute estimate carries a 53-minute span almost entirely spent waiting for the integration slot.

## Context Summary

**Re-read prime files fresh; do not trust carried-over assumptions.** Both `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-kanban-board.md` gained lessons this session, and `prime-shell-commands.md` also gained a new trap entry (a checking tool that reports on stdout while exiting zero makes an exit-status-shaped gate lane decorative).

**Decisions with reach made this session:**
- The offset/fractional refusal in the timestamp repairer is **permanent and pinned**, and the reason is now "the arithmetic is the risk, not the obstacle" — not the repudiated "we cannot decide it".
- A fence-broken REQ file is **refused** by the repairer, matching the board's reader exactly. The refusal list gained an honest entry for a residual the fix itself creates (non-ASCII whitespace padding).
- The canonical gate now fails on unformatted Go, reading the formatter's *output* rather than its exit status.
- The Durations mark-label face has **one** measured descent bound, not two disagreeing ones.
- A conditional in standing prose is clutter when it keys on a count that does not bear on the argument, and a boundary when it keys on a condition that does. The directory it lives in is not the discriminator.
