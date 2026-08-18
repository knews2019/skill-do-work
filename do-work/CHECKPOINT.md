---
session_ended: 2026-08-18T21:06:10Z
last_completed: REQ-255
queue_state: 12 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 12
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

- REQ-257: Decide whether the timestamp repairer learns offset and fractional stamps — claimed 2026-08-18T21:16:24Z — writer: vm:/home/user/skill-do-work

## Completed This Session

Twelve REQs shipped, 0.212.8 through 0.212.19, each with an independent adversarial review:

- REQ-246: repair wrong queue/working timestamps from the session hook (Route C, 78%) — commit `270a2d0`, **0.212.9**
- REQ-249: decide the cross-package citation path form and sweep (Route B, 97%) — `cc1083c`, **0.212.10**
- REQ-248: anchor the Durations day buckets to UTC midnight (Route B, 96%) — `1cb897f`, **0.212.11**
- REQ-251: retire the stale copies of the future-stamp message (Route A, 95%) — `96bb593`, **0.212.12**
- REQ-250: close the remaining markdown link checker gaps (Route B, 97%) — `330797b`, **0.212.13**
- REQ-247: archive timestamp audit tool driven by git commit times (Route C, 94%) — `4035ddc`, **0.212.14**
- REQ-253: decide the Timestamp rule's two uncovered stamp shapes (Route A, 95%) — `0d8d629`, **0.212.15**
- REQ-254: let qualify tell a check's own output from instrumentation (Route B, 91%) — `116eec6`, **0.212.16**
- REQ-252: record the browser with every measured-face number (Route B, 97%) — `c752529`, **0.212.17**
- REQ-256: disclose the session hook's queue write surface (Route A, 96%) — `fbc14e8`, **0.212.18**
- REQ-255: give the timestamp repairer shape parity with the read-side detectors (Route B, 95%) — `84add20`, **0.212.19**

Plus **0.212.8**, an environment fix committed before the queue run: the board's JavaScript behavior probes now reach Node on stdin instead of as an argv entry, because a probe embedding the assembled client exceeds Linux's 128 KiB per-argument limit. That failure made `maintainer-verify.sh` red on this checkout for a reason unrelated to any REQ.

Every hash was confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exited 0 at every commit boundary. All eleven worktrees and branches were removed with `git worktree remove` / `git branch -d` (never `-D`); the worktrees parent directory is gone and `git worktree list` shows only the main tree.

## Still Queued

Twelve, all `pending`, none claimed. Every one was answered by the user in a single `do-work clarify` session — there are no open questions anywhere in the queue.

- **REQ-257** — teach the repairer offset and fractional-second stamps (gated behind REQ-255, now satisfied)
- **REQ-258** — split the shell behaviour suite per script (user chose the non-recommended branch)
- **REQ-259** — retire the skill-root citation reading at its three unbackticked sites *(externally corroborated on PR #145)*
- **REQ-260** — run the Go formatter as part of the canonical verify *(rescoped by the user from a one-character fix to the rule behind it)*
- **REQ-261** — delete the date-only tripwire and keep the rule *(user settled the underlying question; see Decisions below)*
- **REQ-262** — govern the prompt-kit templates' date headers
- **REQ-263** — tighten qualify's ownership probe and make its WARN legible *(externally corroborated)*
- **REQ-264** — make a disarmed P-A-U audit visible
- **REQ-265** — raise the two under-bounding mark-label constants
- **REQ-266** — name builds beside the JS renderer's measured face numbers (sweep key `durations-measured-face-constants-lack-provenance`)
- **REQ-267** — close the two remaining repairer shape divergences (sweep key `repairer-detector-shape-parity`)
- **REQ-268** — make the archive audit's clean answer trustworthy *(one instance externally reported, both orchestrator-reproduced)*

## Session Notes

- **The instance-versus-class failure recurred in nine of twelve REQs, and the two that beat it did the same thing.** REQ-250 and REQ-255 both grepped or fuzzed the *primitive* before declaring a class closed, and both found the real hole somewhere the instance list never pointed — REQ-250's in a different checker's copy of the same path logic, REQ-255's in a value shape (unterminated frontmatter fence) that no reviewer had contemplated. Every other REQ shipped a mechanism that looked complete and wasn't. **An instance list is a sample; the primitive is the class.**
- **Three independent review channels found different things, and none was redundant.** The pipeline's own adversarial reviews found what execution reveals (mutation-falsifiability, fuzz corruption, class holes); the external automated reviewer on PR #145 found what fresh eyes reveal (calendar-impossible erasure, first-versus-last duplicate keys, a failing file walk reported as clean); the orchestrator's own re-reading found what bookkeeping reveals (a commit message claiming a P-A-U transcription that silently no-op'd). All three fired on work the other two had already passed.
- **A gate can be silently disarmed by the shape of the file it audits.** Every review-generated REQ this session lacked the P-A-U block that qualify's debug-artifact check keys on, so half that gate was off for the REQs most likely to need it — and the false "transcription" claim in two commit messages went unnoticed until a reviewer armed the check by hand. REQ-264 makes that visible; the record was corrected in REQ-254's archived trail rather than by rewriting history.
- **The environment needed real repair before any REQ could be judged.** Go 1.26.1, ShellCheck 0.11.0 and `just` were installed, and one pre-existing Linux-only test failure fixed, before the canonical gate could exit 0 at all. A red baseline is not a REQ's fault, but nothing downstream means anything until it's green.
- **Estimator calibration:** twelve rows appended. Estimates ran 5–45 minutes against wall spans of 12–135, but the wall figures are inflated by server-side API failures (several builders were killed mid-run by 529s and resumed) and by serial integration queuing behind other REQs' reviews. Route A held closest; Route B and C overshot in the estimate and undershot in wall time. **Do not recalibrate from this session's spans** — they measure the harness, not the work.
- **Concurrency held.** Four waves of two to three builders, disjoint write sets each time, zero merge conflicts and zero invisible collisions — the failure mode from the previous session (two REQs declaring the same constant in different files of one Go package) did not recur, because REQ-252 was gated behind REQ-248 rather than run beside it.

## Context Summary (heavy sessions only)

**Read these fresh before starting; twelve REQs of carried-over assumptions are not reliable.**

- `_dev/primes/prime-shell-commands.md` — gained REQ-246, 247, 250, 254 and 255's lessons. The load-bearing one: fix at the primitive and fuzz the value space.
- `_dev/primes/prime-kanban-board.md` — gained REQ-248, 251 and 252's. Measured constants now name their browser build, held by a test.
- `_dev/primes/prime-action-files.md` — gained REQ-249 and 253's. The citation convention changed this session: literal relative paths from the citing file's own directory, checked mechanically.
- `_dev/tests/prescribed-shell-scripts-behavior.sh` — 42 named cases at session start, 64 now.
- `_dev/tests/shipped-package-reference-contract.sh` — now also checks backticked cross-package citations in both topologies, and validates bare `#anchor` links.

**Decisions with reach beyond their own REQ:**

- **The Timestamp rule's date-only tripwire is a builder's prose, not a maintainer decision** — established during clarify by tracing it to the repository import. The user settled it: keep the rule, delete the "revisit if a second consumer appears" clause, add no date-only subcommand. Consumer count has no bearing on whether a shell one-liner suffices. That is REQ-261, and it is the model for any tripwire keyed on a count rather than a condition.
- **REQ-255 D-04: duplicate stamp keys are repaired last-occurrence, not refused** — refusal would have made the SessionStart hook exit nonzero every session on a file the board renders fine. When a fix's two options are "match the read side" and "refuse loudly", an unattended path should not be the one that fails forever.
- **REQ-254 D-02: printed output belongs to whoever owns the process exit** — the condition that separates a check's own reporting from leftover instrumentation. Its implementation (a whole-file text grep) is weaker than its statement, which is REQ-263.
- **REQ-247 D-01: share by sourcing, not duplication** — the archive auditor holds zero predicate code of its own, which is why REQ-255's six-shape fix reached both tools in one edit. Verified by grep at review time, not assumed.
- **REQ-252 D-01: hold a documentation convention with a vacuity-guarded test** — comments do not stay honest on their own; the test walks every Go file and fails a measured constant whose comment names no build.
