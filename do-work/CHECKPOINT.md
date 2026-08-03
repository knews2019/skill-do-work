---
session_ended: 2026-08-03T22:18:00Z
last_completed: REQ-082
queue_state: 5 pending, 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 6
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

Nothing. Stopped at a REQ boundary on user request — REQ-082 was carried through commit rather than
left claimed, so `do-work/working/` holds only `baseline.json` (a pre-flight record, not a claim).

A REQ found in `do-work/working/` by a future session is therefore a **foreign claim** and must be
left byte-identical (`actions/work-reference.md` → Crash Recovery (Step 1)).

## Completed This Session

- REQ-077: Crash recovery's own-crash branch is unreachable (Route C, 96%) — v0.167.1, commit `598ef35`
- REQ-078: The Windows timestamp fallback cannot run on stock Windows (Route C, 94%) — v0.167.2, commit `7998740`
- REQ-079: Two guards pin the weaker fingerprint (Route B, 97%) — v0.167.3, commit `8fdce3c`
- REQ-080: The capture template emits a stray instruction line (Route A, 98%) — v0.167.4, commit `8ee717b`
- REQ-081: `next-version` ignores flags after the bump size (Route B, 95%) — v0.168.0, commit `84d79c1`
- REQ-082: The fan-out hand-back file has no legal write location (Route B, 96%) — v0.168.1, commit `1cff0a7`

**UR-015 is fully built** (REQ-077…080). **UR-016 is 2 of 5** (REQ-081, REQ-082 done; REQ-083, 084,
085 remain). Both UR folders stay in `do-work/user-requests/` because each has unresolved members —
UR-015 now holds three review-generated follow-ups.

## Still Queued

- REQ-083: verify reports every builder worktree as a fixable orphan (pending) — UR-016, original batch
- REQ-084: verify's queue-state probe misses a committed builder (pending) — UR-016, original batch
- REQ-085: run REQ-073's live two-builder acceptance test (pending) — UR-016, original batch
- REQ-086: in-progress record unstated at three consumer sites (pending, `depends_on: [REQ-077]` — satisfied) — review-generated
- REQ-087: board and verify hand users the POSIX-only timestamp command (pending, `depends_on: [REQ-078]` — satisfied) — review-generated
- REQ-088: confirm the fix for memory-reference.md's export-ignored CLAUDE.md citation (**pending-answers**) — discovered task, needs `do-work clarify`

## Session Notes

- **Rigor is calibrated per REQ, by user decision this session.** Full Route C pipeline for the
  genuinely complex ones; lean Route A/B for well-scoped ones. The Restatement Sweep runs on every
  REQ regardless — it produced the only two defects the builds themselves missed (REQ-077's third
  premise fingerprint; REQ-078's three tool-side prescriptions).
- **KB handoff is batched, by user decision.** Every completed REQ carries `kb_status: pending`,
  including REQ-074/075/076 from the previous session. Nine REQs' lessons are now waiting. Offer one
  consolidated drop into `kb/raw/inbox/` at the end of the run.
- **A sweep REQ's site inventory is a floor — third batch running.** REQ-077's audit named two premise
  fingerprints; the tree had three. REQ-078's named eleven inline sites; a third grep shape found
  three more outside `actions/`. REQ-081's requirement 5 predicted "none" and the audit found the
  defect live on five subcommands. REQ-079 and REQ-082 were the exceptions — their inventories held.
  **Run the second and third grep shapes before believing any REQ's site list.**
- **Watch the staging trap.** `git mv` stages the file's content *at move time*. REQ-077's first
  commit captured the pre-edit REQ because the later appends were never re-added; it had to be
  amended. Always `git add` the archived REQ explicitly after every append, before committing.
- **Observing the RED is cheap and repeatedly worth it.** Running `HEAD`'s copy of the suite against
  an injected regression took under a minute in REQ-077, REQ-079 and REQ-081, and in each case turned
  an argument about coverage into an observation. `git show HEAD:_dev/tests/contract-regressions.sh >
  _dev/tests/.old-suite.sh` works because `repo_root` derives from `BASH_SOURCE`; `git archive` does
  **not** work, since `_dev/`, `CLAUDE.md` and `.gitattributes` are all export-ignored.
- **Two REQs deliberately extended their own scope, both declared.** REQ-077 D-01 (a stale gloss in a
  declared file) and REQ-081 D-02 (leftover-argument rejection across all six subcommands, pulling in
  `serve.go`). Both amended `## Scope` and `write_set` when the drift was taken, not after, and both
  flagged it for the reviewer rather than presenting it as inevitable.
- **REQ-084 must confirm REQ-082's reconciliation.** REQ-082 landed first; its requirement 7 says
  whichever REQ lands second confirms in its Qualification that the hand-back file cannot trip the
  builder-wrote-`do-work/` probe. The structural argument is already written in REQ-082's
  Qualification — REQ-084 needs to restate it against the probe it actually builds.

## Context Summary (heavy session)

**Key decisions worth carrying:**

- **REQ-077 D-02 / REQ-082 r5 — closed enumerations, twice in one session.** Both REQs shipped a rule
  whose first draft was a list and whose final form is a trigger condition with the list marked
  illustrative. The pull toward enumeration is strong enough that it happened *inside* a REQ about
  enumerations going stale.
- **REQ-079 D-02 — a guard's blast radius is the pattern *and* the filter.** The pattern decides what
  the premise looks like; the filter decides where it is wrong. Widening the first without
  understanding the second either misses the class or flags true statements.
- **REQ-081 D-02 — `os.Exit` in a command function is an untestability boundary, and bugs hide there.**
  The whole argument-handling path had zero coverage while the surrounding release logic had plenty.
- **REQ-082 D-01 — put a carve-out inside the sentence it modifies**, so a maintenance pass cannot
  read the prohibition without also reading the exception.

**Prime files to re-read fresh, not trusted from memory:** `tools/queue-kanban/prime-do-kanban.md`
(REQ-083 and REQ-084 both list it; its `## Lessons` are inlined-not-linked and gained a REQ-081 entry
this session).

**Contract files this batch has been amending repeatedly** — expect overlap and read the hunks rather
than assuming contamination: `actions/work-reference.md` (REQ-077, 078, 082),
`_dev/tests/contract-regressions.sh` (all six), `actions/work.md` (REQ-077, 078, 081).
