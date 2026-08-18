---
session_ended: 2026-08-18T12:29:36Z
last_completed: REQ-239
queue_state: 2 pending, 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 10
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

- REQ-242: Stop Panel B's slowest-day annotation colliding with its own title — claimed 2026-08-18T13:05:12Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2
- REQ-244: Cite the Timestamp rule at every timestamp write site — claimed 2026-08-18T13:05:12Z — writer: t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2

## Completed This Session

- REQ-231: Keep Panel A's direct labels clear of the mark band (Route B, 97%) — commit `720f23c`, shipped as **0.208.0**
- REQ-230: Point caller docs at the canonical publication rationale (Route B, 98%) — commit `19669fc`, shipped as **0.208.1**
- REQ-234: Stop the shell behavior suite counting its own cases (Route A, 97%) — commit `48ed251`, shipped as **0.208.2**
- REQ-233: Give the Timeline a keyboard path to zoom and pan (Route B, 96%) — commit `9b2578b`, shipped as **0.209.0**
- REQ-236: Add a URs-only lens to the Board view (Route C, 97%) — commit `456ee9d`, shipped as **0.210.0**
- REQ-235: Give the Timeline period navigation and a jump to now (Route C, 95%) — commit `7cae7a4`, shipped as **0.211.0**
- REQ-240: Stop the Timeline axis printing a fake minute (Route B, 97%) — commit `664b269`, shipped as **0.211.1**
- REQ-237: Backfill the Durations label rows when the longest spans cluster (Route B, 96%) — commit `3720ab9`, shipped as **0.212.0**
- REQ-238: Point present-work at the canonical independent-bytes rationale (Route B, 98%) — commit `d783ec9`, shipped as **0.212.1**
- REQ-239: Give the Timeline's rows a real focus ring (Route B, 97%) — commit `1d76ad1`, shipped as **0.212.2**

Every hash was confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exits 0 at every commit boundary. All ten ran as worktree builders; every worktree and `worktree-agent-*` branch was removed with `git worktree remove` / `git branch -d` (never `-D`), and `git worktree list` shows only the main tree.

## Still Queued

- **REQ-241** (pending): Reconcile the Durations label metrics with the rendered face. Sweep, two instances — `durationsLabelCharacterWidthUnits = 6.2` against a measured 6.61 units/char, and `DURATIONS_LABEL_ROW_HEIGHT = 12` against a declared 13-unit text box. Approved via clarify.
- **REQ-242** (pending): Stop Panel B's slowest-day annotation colliding with its own title. Pre-existing; `209 min` renders through the title on a fixture whose slowest day falls under it. Approved via clarify.
- **REQ-243** (pending-answers): Check that shipped markdown pointers actually resolve. New machinery rather than a repair, which is why it is asking.

**REQ-241 and REQ-242 must not run in parallel with each other.** Both write `web/board-durations.js`, and REQ-241 may move `DURATIONS_LABEL_ROW_HEIGHT` — which shifts every panel below it, including the Panel B title REQ-242 is separating an annotation from. Run REQ-241 first, then re-measure REQ-242's collision against the result: it may change size or disappear. REQ-243 is fully disjoint (`_dev/tests/prescribed-shell-canonicalization.sh`) and can run alongside either.

## Session Notes

- **Every visual REQ this session had a defect that the suite passed and a render caught.** REQ-231's dot-on-label overprinting (55 overlapping pairs → 0), REQ-240's axis reading `18 Aug 11:00` five times, REQ-237's under-filled label rows, REQ-239's open question about whether a ring even suited an 18px row. In each case the assertions were green before and after the discovery. The rule from REQ-226 held for the fourth, fifth and sixth time: generate a board and look at it.
- **A shared browser instance silently invalidates cross-builder DOM evidence.** REQ-237's builder measured a sibling's page and got confident, well-formed, wrong numbers; only unfamiliar REQ ids gave it away. REQ-239's builder was warned mid-build and re-took every reading in an isolated context with `location.href` and the page's own `cssRuleText` returned from the same `evaluate` — all identical. Both facts are now conventions in `_dev/primes/prime-kanban-board.md`.
- **Two builders pushed back on their briefs, and both were right.** REQ-237's brief said `TestOverflowLabelsGoToTheLongestSpans` must pass unchanged; that test's assertions *were* the six-label cap the REQ removes, so no correct implementation could satisfy both. REQ-235's builder wrote a period-index clamp, found the tests passed without it, and deleted it as unearned. A builder that had quietly edited the test, or kept the dead clamp, would have looked identical in the diff.
- **Merge conflicts in `generate_test.go` are not all alike.** Three occurred; two were clean appends where stripping markers works, and REQ-239's was not — both sides ended mid-function with their shared closing braces folded into the common tail, so the naive strip made one function swallow the other. Nothing in the marker positions distinguishes the cases. Compile, then run *both sides'* tests.
- **Three self-inflicted errors, all caught and corrected:** three future-dated `completed_at` stamps written by extrapolating the clock instead of reading it (caught by the board's own future-timestamp check, commit `818ea17`); a dispatch brief that omitted the P-A-U instruction, so two REQs' boxes were filled by the orchestrator from hand-back evidence with an HTML comment saying so; and a `git mv` that ran after its frontmatter edit had already failed, archiving REQ-237 still reading `status: claimed`.
- **Estimator calibration this session:** ten REQs, estimate vs actual wall minutes — 35/37, 20/6, 5/11, 30/70, 60/20, 75/24, 25/21, 30/28, 15/23, 20/30. The two Route C estimates (60, 75) overshot badly against 20 and 24 actual; Route B clustered close except REQ-233 (30 est / 70 actual, the one REQ that needed an integration seam). Worth a recalibration pass.
- A builder left `timeline-focus-ring.png` in the main tree root — a write-set violation. Removed, and every subsequent brief carried a "scratch goes in /tmp" line.

## Context Summary (heavy sessions only)

**Read these fresh before starting; ten REQs of carried-over assumptions are not reliable.**

- `_dev/primes/prime-kanban-board.md` — gained five lesson links (REQ-233, 235, 236, 237, 239, 240) and **two new conventions** this session: render-and-look, and assert page identity inside the measuring call. Entry point for anything under `skills/do-work-board/tools/queue-kanban/`.
- `_dev/primes/prime-shell-commands.md` — gained REQ-230, REQ-234 and REQ-238.
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` — changed by three REQs today (233, 235, 240) and is the file REQ-241/242's neighbours sit beside.

**Decisions with reach beyond their own REQ:**

- **REQ-233 D-01 through D-06 established that every way of moving the Timeline's window goes through `timelineZoomedWindow`.** REQ-235 then added period navigation as a *third* driver without adding a fourth rule — it computes a candidate window and hands it to that same function. Anything that later moves the window must do the same, or the guarantee stops holding.
- **REQ-235 D-02: the period level is derived from the window, never stored.** That is why "a free zoom marks the level inexact" needed no code. Do not add a `periodLevel` field.
- **REQ-236 D-01: a lens and a lens button are no longer one-to-one.** `viewState.lens` holds two values while the Lens group offers three; the third is a fold modifier. Four sites read `viewState.lens === "user-request"` and are correct for both readings *because* of that. A future lens must check this assumption.
- **REQ-237 [MAP CHANGED]: Durations label selection and placement are one descending-magnitude pass.** There is no candidate list and no `durationsLabelTopCount`. Row occupancy is an interval list because a magnitude-ordered walk does not visit x monotonically. Anything influencing *which* spans get labelled has one place to do it and must not reintroduce a pre-placement filter.
- **REQ-240 D-02: the axis label format keys on the gap between ticks, not the window span.** `TIMELINE_AXIS_TICK_COUNT` exists so the formatter's threshold and `renderAxis`'s loop read one number. A third reader must join the same constant.

**Architectural note:** the board now has three chart-ish views with different disciplines — Durations rebuilds its nodes each render and needs no teardown; Timeline binds to the scroll host and `window` and keeps an explicit teardown registry; the Board lenses share one renderer with a fold flag. Copying any one's habit into a fourth view without asking which case it is will be wrong two times in three.
