---
session_ended: 2026-08-04T09:15:44Z
last_completed: REQ-093
queue_state: 0 pending, 0 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 in-progress
reqs_processed_this_session: 2
session_depth: moderate
---

# Session Checkpoint

## In Progress (interrupted)

Nothing. Both REQs were carried through commit and `do-work/working/` is empty — the stale
`baseline.json` a previous session left there was swept this run. A REQ found in `working/` by a
future session is therefore a **foreign claim** and must be left byte-identical
(`actions/work-reference.md` → Crash Recovery (Step 1)).

## Completed This Session

- REQ-088: fix memory-reference.md's citation of the export-ignored CLAUDE.md (Route A, 96%) —
  v0.169.3, commit `bb8cf3b`
- REQ-093: six shipped Go-source sites cite the export-ignored CLAUDE.md; suite guard caught none
  (Route B, 95%) — v0.169.4, commit `2e29b36`

**REQ-090 was cancelled**, not built — its premise did not reproduce (see Session Notes).

## Queue State — Empty

**The queue is empty and both open URs closed this session.** UR-015 (8 members) and UR-016
(8 members) are both consolidated into `do-work/archive/`. `do-work/user-requests/` is empty.
Nothing is pending, pending-answers, blocked, or in-progress. This is the first fully-clear
boundary since the UR-015/UR-016 audit-remediation batch opened.

`do-work/HANDDOWN-UR-015-016.md` now describes a finished batch — every REQ in its table is
resolved. It is history, not a live handdown.

## Session Notes

- **The recorded "8 FAIL" suite baseline does not reproduce, and REQ-090 was cancelled on that
  basis.** The previous checkpoint recorded `_dev/tests/contract-regressions.sh` as red with 7
  update-script probe failures, confirmed twice by stash-and-compare. This session the suite exits
  **0 with zero FAIL lines**, the standalone probe passes, and `tools/do-work-update.sh`
  demonstrably contains all five strings REQ-090 says are absent (lines 166, 194, 202, 204, 218).
  Ruled out: lingering worktrees (none), cwd sensitivity (green from a subdirectory), a vacuous skip
  path (the probe's only skip is `git` unavailable; `git` is present), and an intervening fix
  (neither file committed to since `b583c78`). **The original observation is unexplained** — that is
  recorded rather than resolved, in the cancelled REQ and in the KB doc for REQ-083 whose lesson
  asserted the red state.
- **The site inventory was a floor for the sixth consecutive batch, and the *guard* was the deeper
  defect.** REQ-088 named one dangling citation; a shipped-path sweep found 8 mention-lines across
  4 files. Running the suite's own `self_citation_pattern` verbatim scored **0 of 8** — it had never
  caught anything, including the line it was written to prevent.
- **Inverting a check beat widening it, and the user's correction was right on the evidence.** The
  builder recommended adding `§` and `:` to the idiom list. Measured: that would have caught 4 of 6
  Go sites (missing two bare-prose ones) and false-positived on `actions/prime.md`, where the colon
  is sentence punctuation. The inverted check — flag any mention, exempt a 14-file per-file
  allowlist — scores **8 of 8**. The general lesson: when a guard enumerates *how* a defect is
  phrased, the enumeration is the bug (`CLAUDE.md` → Closed Enumerations Go Stale).
- **Two handdown traps combined into one silent failure, and it was walked into anyway.** `git mv`
  stages content at move time, and `git add` aborts the whole invocation on one bad pathspec. Both
  are written down in `do-work/HANDDOWN-UR-015-016.md`. Passing a path that `git mv` had already
  moved aborted the staging, and the commit then *succeeded* carrying REQ-090's pre-cancellation
  content and none of REQ-093's answers. **`git status --porcelain -uall` between staging and
  committing is the check** — it was reading the commit's own `--name-status` afterwards that caught
  it, one step too late. The same trap was then anticipated and avoided for REQ-093's archive.
- **The KB backlog is cleared: 34 REQs promoted** into `kb/raw/inbox/` (37 documents there in
  total, including 3 from a prior drop). Prose sections are extracted from each REQ's own words
  rather than paraphrased. One document carries an explicit `[Superseded]` note: REQ-083's
  "suite is red" lesson would otherwise have planted a falsified fact.
  **Next step is `do-work bkb triage`, then `do-work bkb ingest`** — nothing has been compiled into
  the wiki yet, and `kb/raw/_inbox_queue.md` is still an empty table (triage populates it).
- **A near-miss worth remembering:** clearing regenerated output with
  `rm -f kb/raw/inbox/REQ-0[0-9][0-9]-*.md` also matched three pre-existing inbox files from an
  earlier drop (REQ-074/075/076). They were tracked, so `git checkout --` restored them, and the
  `kb_entry`-resolves check is what surfaced it. A glob written to match *your* outputs will match
  someone else's too.

## Context Summary

**Prime files to re-read fresh rather than trust from memory:**
`tools/queue-kanban/prime-do-kanban.md` — checked for staleness this session (its referenced paths
all exist) but it gained no entry, because neither REQ touched board behavior, only comments.

**Contract files amended this session — read the hunks:** `_dev/tests/contract-regressions.sh`
(the citation check was replaced, not tweaked — its enforcement point moved from a pattern list to
`maintainer_doc_mention_allowlist`), `CLAUDE.md` (line ~122 now describes the inverted check),
`actions/memory-reference.md`, `tools/queue-kanban/verify.go` and `verify_test.go` (comments only).

**Two open threads a future session may want:**

1. **The allowlist can rot in two directions the inverted check cannot see** — an entry whose file
   no longer mentions the maintainer doc (the exemption outlives its reason and could hide a real
   citation), and an entry naming a deleted file. Recorded as REQ-093's `[low]` Discovered Task and
   deliberately not built (D-02) — no failing case exists yet.
2. **REQ-088's fixed line is now mildly redundant** — *"nothing carries between command blocks
   (shell state does not survive between prescribed command blocks)"*. The user's authorized
   wording; the declined option (a) would have read better. Recorded as a Minor review finding on
   REQ-088, not queued.
