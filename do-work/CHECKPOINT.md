---
session_ended: 2026-08-04T04:28:03Z
last_completed: REQ-088
queue_state: 0 pending, 2 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 1
session_depth: light
---

# Session Checkpoint

## In Progress (interrupted)

Nothing. REQ-088 was carried through commit and `do-work/working/` holds no claim. A REQ found
there by a future session is therefore a **foreign claim** and must be left byte-identical
(`actions/work-reference.md` → Crash Recovery (Step 1)).

One untracked leftover does sit in `do-work/working/`: `baseline.json`, a stale pre-flight artifact
from the 2026-08-03 session (`"exit_status": 0`). It is not a claim and blocks nothing, but it is
available to the next REQ's Step 6.5 baseline comparison, where it would read as that REQ's own
measurement. Recorded as REQ-088's `[low]` Discovered Task; sweeping it is `do-work cleanup`'s job
and is awaiting the user's consent.

## Completed This Session

- REQ-088: fix memory-reference.md's citation of the export-ignored CLAUDE.md (Route A, 96%) —
  v0.169.3, commit `bb8cf3b`

## Still Queued

- REQ-090: seven update-script behavior probes fail on the base branch (**pending-answers**) —
  UR-016. **Its premise did not reproduce this session** — see Session Notes.
- REQ-093: six shipped Go-source sites cite the export-ignored CLAUDE.md, and the suite's guard
  catches none of them (**pending-answers**) — UR-015, raised by REQ-088.

Both need `do-work clarify`. Nothing is `pending`; nothing is blocked.

## Session Notes

- **The recorded "8 FAIL" suite baseline does not reproduce.** The previous checkpoint recorded
  `_dev/tests/contract-regressions.sh` as red on the base branch with 8 FAIL lines — 7 update-script
  behavior probes plus the summary — and said it was confirmed pre-existing by stash-and-compare.
  This session the suite exits **0 with zero FAIL lines**, the standalone probe reports
  `update-script behavior probes passed`, and `tools/do-work-update.sh` demonstrably contains all
  five strings REQ-090 says are absent (lines 166, 194, 202, 204, 218). Ruled out: lingering
  worktrees (`git worktree list` shows only the main checkout), cwd sensitivity (green from a
  subdirectory), and a vacuous skip (the probe has exactly one skip path — `git` unavailable — and
  `git` is present). Neither file has been committed to since `b583c78`. **The cause of the earlier
  observation is unexplained**; what is established is that the failure state is not present now.
  This is the diagnosis REQ-090's Open Question asked for, and it changes the answer set: the
  live options are now "close it as not-reproducible" or "keep it open to find out why it was ever
  red," not "fix the updater."
- **The site inventory was a floor for the sixth consecutive batch.** REQ-088 named one dangling
  citation; a sweep of every shipped path found six more of the same defect class in
  `tools/queue-kanban/verify.go` and `verify_test.go`, plus a seventh ambiguous one in
  `prompts/prompt-kit-step6-constraint-architecture.md:78`. Queued as REQ-093 rather than swept in,
  because REQ-088's `## Answer` explicitly scoped the change to one line.
- **The suite's citation guard has zero coverage of the idioms actually in use.** Running
  `self_citation_pattern` verbatim over the shipped paths it scans returns **0 hits** — it matched
  neither the six sites above nor the line REQ-088 was filed to fix. It recognizes `see CLAUDE.md`,
  `per CLAUDE.md`, and `CLAUDE.md →`; the real occurrences are `CLAUDE.md §` and `CLAUDE.md:`.
- **The Restatement Sweep paid off in the confirming direction this time.** It found no stale
  restatement, but it did establish that the shell-state rule is already restated inline at nine
  shipped sites and cites `CLAUDE.md` at none of them — which turned REQ-088's option (b) premise
  ("no shipped file owns this rule today") from an assumption into a verified fact, and showed the
  chosen fix is the established house pattern rather than a one-off.
- **UR-015 did not close, by one REQ.** REQ-088 was its last unresolved member; REQ-093 is a new
  `pending-answers` member of the same UR, so UR-015 stays in `do-work/user-requests/` and REQ-088
  archived to `do-work/archive/` root rather than consolidating the UR folder.
- **The KB handoff backlog is larger than previously recorded:** 33 archived REQs carry
  `kb_status: pending`, not the 9 or 17 stated earlier — the earlier counts covered only recent
  batches, while the backlog reaches back to REQ-001. `kb/raw/inbox/` exists, so a consolidated
  drop is possible whenever the user wants it. Still batched, still awaiting consent.
