---
session_ended: 2026-08-18T01:54:45Z
last_completed: REQ-232
queue_state: 0 pending, 4 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 6
session_depth: heavy
---

# Session Checkpoint

## Completed This Session

- REQ-225: State verified-exact-publication once as a condition in the shipped shell guide (Route B, 95%) — commit `a54d5c4`, shipped as **0.205.1**
- REQ-226: Stop the Durations chart from silently overprinting and clipping (Route C, 94%) — commit `787c846`, shipped as **0.205.2**
- REQ-227: Add the Timeline view with two-segment REQ bars (Route C, 92%) — commit `17b9422`, shipped as **0.206.0**
- REQ-228: Project the remaining queue onto the Timeline (Route C, 94%) — commit `2daefd1`, shipped as **0.207.0**
- REQ-229: Verify the published path in the download and screenshot helpers (Route B, 96%) — commit `2f1cde5`, shipped as **0.207.1**
- REQ-232: Stop shipped prose from counting the board's views (Route B, 96%) — commit `ea6b1c3`, shipped as **0.207.2**

Every hash was confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exits 0 at hand-back.

## In Progress (interrupted)
- REQ-238: Point present-work at the canonical independent-bytes rationale — claimed 2026-08-18T11:57:10Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
- REQ-239: Give the Timeline's rows a real focus ring — claimed 2026-08-18T11:57:10Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
- REQ-237: Backfill the Durations label rows when the longest spans cluster — claimed 2026-08-18T11:42:03Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2



## Still Queued

- REQ-230: Point caller docs at the canonical publication rationale (pending-answers — 1 question). Sweep, one instance: `present-work.md:137` restates the container-not-a-collision rationale that REQ-225 gave a single home.

All four are `pending-answers` by design, not by accident: three are cascade-depth or scope-boundary decisions, and the fourth needs a definition of "one case" that only the maintainer should pick. `do-work clarify` is the next step.

## Session Notes

- **The environment had no working baseline on arrival.** `maintainer-verify.sh` could not run at all: no `shellcheck`, no `just`, and Go 1.24.7 against the required exactly-`go1.26.1`. ShellCheck 0.11.0, just 1.43.0 and Go 1.26.1 were installed into a session-local directory and prepended to `PATH`; nothing in the repository was changed to accommodate the gap. A future session in a fresh container will hit the same wall — the three-line install is worth having written down somewhere runnable.
- **Writing a rule down worked, immediately and on its own author.** REQ-225 moved the verified-exact-publication rule out of one script's section into a shared condition. Reading the guide's three publication sections against that single statement surfaced two more helpers that never verified — one of them (`capture-screenshot.sh --staged`) destroying the dispatch's only copy of a screenshot on a false success. That is the fifth and sixth closure of this defect class here, and the first found by reading rather than by a review sweep.
- **Three of the four defects fixed this session were invisible to the test suite and visible in a render.** The Durations label blob, the clipped Panel B bar, and the remainder sentence overprinted by the marks it described all passed every assertion before and after. Generate a board and look at it; a passing suite is not evidence about pixels.
- **A restore from a mid-session copy silently reverted a fix** (REQ-226, D-03) and the whole suite still passed, because nothing pinned which row the remainder used. The near-miss is why `durationsRemainderBaselineY` exists as a named function with its own probe rather than as an inline expression. Prefer re-applying an edit over restoring a file whose copy predates other edits.
- **The board's two chart views now have opposite listener disciplines.** Durations binds to nodes it rebuilds each render and needs no teardown; Timeline binds to the scroll container and `window`, which outlive a render, and keeps an explicit teardown registry. Copying either one's habit into a third view without asking which case it is will be wrong half the time (REQ-227, D-07).
- **Pre-existing anomaly, not touched:** `do-work/user-requests/UR-049/` still holds an `input.md` byte-identical to the one already archived at `do-work/archive/UR-049/input.md`. Commit `c0331e8` consolidated the UR by copying rather than moving. Cleanup never deletes durable artifacts, so it was reported rather than removed; `do-work forensics` is the right place to decide.
- Reservation markers under `do-work/.req-reservations/` are left untracked as usual; `cleanup-req-reservations.sh` reaps them on SessionStart.

## Context Summary (heavy sessions only)

**Read these fresh before starting; six REQs of carried-over assumptions are not reliable.**

- `_dev/primes/prime-kanban-board.md` — gained three lesson links this session (REQ-226, REQ-227, REQ-228, REQ-232) and is the entry point for anything touching `skills/do-work-board/tools/queue-kanban/`.
- `_dev/primes/prime-shell-commands.md` — gained REQ-225 and REQ-229. Its § *Closed Enumerations Go Stale* was applied three separate times this session (REQ-225, REQ-232, and REQ-234's finding), which is worth noticing as a pattern rather than as three coincidences.
- `skills/do-work/docs/prescribed-shell-primitives.md` § **Verified exact publication** — new this session, and now the canonical statement six shipped helpers are measured against.

**Decisions with reach beyond their own REQ:**

- REQ-226 D-02 — placement anchors labels *before* the mark by preference, because a left-to-right greedy walk reuses space it has already passed. Anchoring after cost a label on the real board.
- REQ-227 D-07 — the listener teardown registry, above.
- REQ-228 D-02 — the forecast uses bucket medians rather than the REQ's own `estimate:` block, because the board parses no nested frontmatter blocks. If nested parsing ever lands, that is the first consumer.
- REQ-228's own lesson is the constraint to preserve: the projection can be wrong about *timing* without being wrong about *order*, because it borrows work.md's ordering rule rather than inventing a scheduler. Anything that makes the duration model cleverer must not change the ordering source.

**Architectural note:** `timeline.go` is the first thing on the board that produces a future instant. Its projection is deliberately crude and its honesty machinery — the stated assumptions, the exclusion list, the decline-on-thin-history gate — is what makes it publishable. Adding sophistication there should mean adding caveats, not adding cleverness.
