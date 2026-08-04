---
session_ended: 2026-08-04T00:29:00Z
last_completed: REQ-092
queue_state: 0 pending, 2 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 8
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

Nothing. The run reached a clean queue boundary — every REQ was carried through commit, and
`do-work/working/` is empty. A REQ found there by a future session is therefore a **foreign claim** and
must be left byte-identical (`actions/work-reference.md` → Crash Recovery (Step 1)).

## Completed This Session

- REQ-083: verify reports every builder worktree as a fixable orphan (Route B, 96%) — v0.168.2, commit `f6c1514`
- REQ-084: verify's queue-state probe misses a committed builder (Route B, 97%) — v0.168.3, commit `0d61054`
- REQ-086: in-progress record unstated at the two out-of-pipeline movers (Route A) — v0.168.4, merge `a17e6af`
- REQ-087: board and verify hand users the POSIX-only timestamp command (Route A) — v0.168.5, merge `5cfe1b5`
- REQ-085: run REQ-073's live two-builder acceptance test (Route B, 92%, **Acceptance: Partial**) — v0.168.6, commit `b224e8a`
- REQ-089: board drawer's Copy omits the ticket's frontmatter (Route B, 96%) — v0.169.0, commit `ba54b5d`
- REQ-091: hand-back merge collides with the owner's staged claim (Route A) — v0.169.1, commit `ecf1966`
- REQ-092: no wave-driving path in actions/work.md (Route A) — v0.169.2, commit `92bebe0`

**REQ-086 and REQ-087 were built by two concurrent builders in git worktrees**, as REQ-085's live
fan-out acceptance test. They carry merge hashes, not plain commits — read their diffs with
`git show --first-parent -m <hash>`.

**UR-017 is closed and archived.** UR-015 and UR-016 each stay in `do-work/user-requests/` on one
`pending-answers` member (REQ-088, REQ-090).

## Still Queued

- REQ-088: confirm the fix for memory-reference.md's export-ignored CLAUDE.md citation (**pending-answers**) — UR-015
- REQ-090: seven update-script behavior probes fail on the base branch (**pending-answers**) — UR-016, raised by REQ-083

Both need `do-work clarify`. Nothing is `pending`; nothing is blocked.

## Session Notes

- **The contract suite is red on the base branch and was red before this session.**
  `_dev/tests/contract-regressions.sh` reports 8 FAIL lines — 7 update-script behavior probes covering
  `tools/do-work-update.sh` plus the summary line. Confirmed pre-existing by stashing REQ-083's changes
  and re-running. Every REQ this session was checked against that count rather than against zero. This
  is REQ-090's subject; until it resolves, "8 FAIL" is the baseline, not a regression.
- **The site inventory is a floor — for the fourth and fifth consecutive batch.** REQ-087's REQ named
  three sites carrying the POSIX-only timestamp command; a grep across `tools/` found a fourth
  (`model.go:321`), and fixing it broke a test that pinned the old literal. REQ-083's requirement 6
  listed five fixture cases, none of which exercised the third state its own split introduced — the
  review caught it. **Re-run the grep before believing any REQ's list**, and expect the extra site to
  have a test pinned to it.
- **Two REQs' own reviews found real gaps and closed them in place** (REQ-083's untested third
  category, REQ-089's untested serve path). Both are recorded in the REQ rather than left silent — a
  self-caught gap that leaves no trace reads later like a gap nobody looked for.
- **KB handoff is still batched and still pending.** Every REQ this session carries `kb_status: pending`,
  on top of the nine already waiting from the previous two sessions. That is now seventeen REQs' lessons
  queued for one consolidated drop into `kb/raw/inbox/`.

## Context Summary (heavy session)

**Key decisions worth carrying:**

- **REQ-085 D-01 / the fan-out run — a documented-but-unexecuted path was broken, again.** The live
  two-builder test found two defects in five minutes, both invisible to grep because neither half was
  wrong alone: the hand-back merge collides with the owner's own staged claim (REQ-091, fixed), and
  nothing in `actions/work.md` drives a wave (REQ-092, decided as documentation). This is the second
  consecutive batch where executing a spec beat reading it — REQ-082 was the first.
- **REQ-092 D-01 — "document the boundary" was the right answer to "should we build this."** The
  capability was never unreachable; it had been performed successfully by hand from the spec. What was
  missing was a sentence saying which document owns it. The tell that a gap is documentation rather
  than machinery: you can already do the thing.
- **REQ-083 D-03 / REQ-084 D-01 — a category is a routing decision, not a label.** Both REQs split one
  finding into several because the *remedies* differed, not because the states did. Where two states
  need the same action, one category is right; where they need different ones, prose in a `Remedy`
  string is not enough, because nothing can act on prose.
- **REQ-089 — test the round-trip from parsed files, never hand-built structs.** A zero-value field
  concatenates to nothing, so a struct-level fixture passes identically before and after. Only a
  fixture read from disk — with a comment line and irregular spacing in it — can tell "original bytes
  survived" from "a plausible fence was produced."

**Prime files to re-read fresh, not trusted from memory:** `tools/queue-kanban/prime-do-kanban.md` —
it gained three entries this session (REQ-083, REQ-084, REQ-089), and its `## Lessons` are
inlined-not-linked, so the entries are the record.

**Contract files amended repeatedly this session** — read the hunks rather than assuming:
`tools/queue-kanban/verify.go` (REQ-083, REQ-084, REQ-087), `actions/forensics.md` (REQ-083, REQ-084,
REQ-086), `actions/work.md` and `actions/work-reference.md` (REQ-091, REQ-092),
`tools/queue-kanban/model.go` (REQ-087, REQ-089).
