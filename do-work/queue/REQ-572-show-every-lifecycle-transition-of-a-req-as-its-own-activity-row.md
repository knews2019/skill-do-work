---
id: REQ-572
title: 'Show every lifecycle transition of a REQ as its own Activity row'
status: pending
created_at: 2026-09-04T23:16:00Z
user_request: UR-115
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
related: [REQ-573]
batch: activity-history
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/activity.go
  - skills/do-work-board/tools/queue-kanban/activity_test.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
route: B
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-09-04T23:23:52Z
  basis:
    - Route B
    - 6-file write set
    - 3 subsystems involved
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
dispatch_at: 2026-09-04T23:29:19Z
builder_handback_at: 2026-09-04T23:41:02Z
integration_at: 2026-09-04T23:41:59Z
gate_deferred: 'true'
depends_on: [REQ-574]
deferred_implementation_base: 7ad53bff1d867f1453e1e7765e988dedb308e7e1
deferred_implementation_merge: fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
---

# Show Every Lifecycle Transition of a REQ as Its Own Activity Row

## What

The Activity view lists one row per REQ: its newest lifecycle stamp and nothing else. Change the aggregation so every parseable lifecycle stamp inside the window becomes a row, so a REQ that was captured, claimed, dispatched, merged, reviewed, released and completed in the last 24 hours shows all of those transitions, newest first, on the same surface. The Board's detail drawer only prints Created, Claimed and Completed, and the Timeline only draws two spans, so today the full path of a REQ is readable only in its frontmatter or with `git log`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Builder read both primes, both lesson satellites and the crew rules, then settled a four-step approach (RED tests first; aggregation loop over `lifecycleTimestampFields` with a third sort key; client counts; comment restatements). Recorded under `## P-A-U` in `do-work/runs/work-2026-09-04-232225/REQ-572-handback.md`.
- [x] **[APPLY]:** Two commits on the builder branch (6ed61142, 2d8beb40) touching exactly the six Scope files; no file outside Scope written.
- [x] **[UNIFY]:** `git diff --stat main...HEAD` reviewed (6 files, +347/-130); `gofmt -l .` empty; `go vet ./...` clean; debug-artifact scan over added lines empty; each file reviewed line by line (activity.go, activity_test.go, generate.go, board-activity.js, template.html, javascript_behavior_c_test.go); worktree `git status --porcelain` empty.

## Detailed Requirements

- `buildActivityRows` (`activity.go`) emits one `ActivityRow` per parseable stamp on each ticket, not one per ticket. `newestLifecycleStamp` goes away or becomes the single-row special case; do not add a second stamp list, `lifecycleTimestampFields` in `model.go` stays the one enumeration.
- Ordering stays newest first with the same deterministic tie-break, so several rows of one REQ interleave with other REQs by time.
- The client (`board-activity.js`) keeps windowing against the wall clock and keeps the shared filter chips. The summary line must say what it now counts: rows are transitions, REQs are distinct ids, and both numbers are useful ("38 transitions across 21 REQs in the last 24 hours").
- A ticket with no parseable stamp is still skipped, never dated from the zero time.
- Prefer showing all transitions by default over adding a "latest only" toggle. If a toggle is kept, "every transition" is the default state and the toggle is one extra button in the existing Activity window group, not a new control family.
- Click behavior and the Board's `data-detail-kind` attribute are REQ-573's concern (opening the drawer and highlighting sibling rows); this REQ only changes the row set and the counts.

## Constraints

- Go decides which stamps exist and what each records; the client draws and filters. No second definition of a stamp's meaning in JavaScript.
- The existing window, filter chips, and empty-state messages keep working; only the row set and the counts change.
- REQ-571 (removing the board's `pending-heavy-testing` reader case) may touch `model.go` in the same period; this REQ does not edit that file.

## Dependencies

REQ-573 (click a row to open the drawer and highlight sibling rows) depends on this REQ.

## Red-Green Proof
**RED prompt/case:** A ticket with `created_at: 2026-09-04T22:52:00Z` and `claimed_at: 2026-09-04T23:00:17Z` passed to `buildActivityRows` returns one row ("claimed"). On the running board, REQ-570 (deleting the pending-heavy-testing status) in the 24h Activity window shows a single "claimed" row; its capture eight minutes earlier is not visible.
**Why RED now:** `newestLifecycleStamp` keeps only the latest stamp per ticket by design (REQ-568, showing recently touched REQs regardless of status).
**GREEN when:** The same ticket returns two rows, "claimed" at 23:00:17 then "captured" at 22:52:00, in that order; on the board REQ-570 appears twice in the 24h window and the summary line reports both the transition count and the distinct REQ count.
**Validation:** User confirmed (verify-requests, 2026-09-04)

## Builder Guidance

The user is certain about the outcome (see every state a REQ went through) and left the shape to the builder. The Go change is small; most of the judgment is in how the table reads when one REQ has six rows among others. Keep the row shape (`id`, `stampField`, `stampAt`, `transition`) so the payload contract and REQ-573 stay stable.

## Assets

- `do-work/user-requests/UR-115/assets/REQ-572-screenshot-1-activity-view-one-row-per-req.png`: the Activity view at 24h, "38 REQs touched in the last 24 hours", columns REQ / Title / Status / What happened / When / Stamp. REQ-570 has one row: status claimed, "claimed", Sep 4 23:00 UTC, stamp `claimed_at`. REQ-505 above it shows "completed" with stamp `completed_at`; REQ-571 shows "captured" with stamp `created_at`; REQ-507 and REQ-506 show "status changed to pending" with stamp `status_changed_at`. A verify finding card sits above the table. The browser find box has "570" highlighted, 1 of 4 matches, all in the finding card and the one row.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban model, parser, UI, timeline, testing, or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban parsing, views, static output". Over the budget on its own.

## Full Context
See `do-work/user-requests/UR-115/input.md` for complete verbatim input.

*Source: "is this only showing the last status of a REQ? how about if I want to see when it went through all the states of it?" / "ok, do-work capture-request"*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (one Activity row per lifecycle stamp, both counts in the summary line) and the Go entry point is named, but how the table reads when one REQ owns several rows and how the client counts distinct REQs need the existing patterns in the board's assembled client. Exploration, no plan.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings: `do-work/runs/work-2026-09-04-232225/REQ-572-exploration.md`.

- **Go change is confined to `activity.go`.** `buildActivityRows` (lines 56-78) calls `newestLifecycleStamp` once per ticket; looping `lifecycleTimestampFields(ticket)` directly and appending one row per parseable stamp is the whole aggregation change. `newestLifecycleStamp` has exactly one caller and goes away. `model.go` needs no edit.
- **Tie-break defect the change introduces:** the comparator breaks ties on `RequestId` only. One REQ often carries two stamps at the same instant (`completed_at` and `status_changed_at` are written together), so `StampField` must become the third sort key or `sort.Slice` decides.
- **Payload shape stays.** `generatedActivityEntry` (`generate.go:308-318`, tags `id`/`stampField`/`stampAt`/`transition`) already accepts several rows per id; only its comments (generate.go:82-86) restate "newest stamp". Nobody but `web/board-activity.js:21` reads `boardData.activity`.
- **Client:** `renderActivity` (`board-activity.js:43-112`) already windows and filters rows, not REQs. The summary (61-65) and both empty-state strings (108-111) count rows as REQs and need the transition/distinct-REQ wording. `data-activity-request` on every `<tr>` stops being unique; keep it that way (REQ-573 needs sibling rows). Helpers: `requestMatchesFilters` (board-filters.js:65), `makeInstantWithRelativeNode` (board-core.js:81), `setActiveButton` (board-controls.js:10).
- **Tests to rewrite** in `activity_test.go`: `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` (line 19, multi-stamp fixtures now yield more rows), `TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField` (92, becomes "orders one ticket's stamps newest first"), the tail of `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` (237-250, `len(rows)==1` becomes `len(rows)==len(declared)`). Boundary, skip and tie tests survive; extend the tie test with the same-REQ/same-instant case.
- **Node behavior lane pattern:** `runJavaScriptBehaviorProbe` + `sliceBalancedBlockAfter(t, indexHtml, "function renderActivity(")` with a hand-stubbed `document.getElementById` (see `javascript_behavior_a_test.go:420-450`). Probes run only with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on` and node on PATH.
- **Test commands:** fast Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`; Node lane: same wrapper with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 -run '^TestJavaScriptBehavior'`. No prime test map exists for queue-kanban.
- **Restatements to update:** activity.go header/doc comments, generate.go:82-86, board-activity.js:2-4 and 53-55, template.html:460-467 comment. CHANGELOG.md's 0.276.0 entry is history and stays. board-guide.md and both lesson satellites never mention the Activity view.
- **Concerns:** near-duplicate rows for a terminal REQ ("completed" and "status changed to completed" at one instant) are shipped as-is, no suppression rule; `board-data.js` activity array grows several-fold, acceptable at this queue size; the read-only prep another session left at `do-work/runs/work-2026-09-05-005615/activity-prep.md` reaches the same file set and recommends the Node test file over `board.css`.

*Generated by Explore agent; condensed by the work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/activity.go` (modify) — one row per parseable stamp, the stamp field as third sort key, the newest-only helper removed, comments restated
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (modify) — RED first: captured two-row case, rewritten newest-only assertions, same-REQ tie case, all-stamps pin
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — comment restatement only; payload shape unchanged
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (modify) — summary with transition and distinct-REQ counts, empty-state wording, header comment
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — Activity section comment only
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modify) — Node lane test that the rendered summary reports both counts from a synthetic multi-row payload

**Files I will NOT touch:** `model.go` (REQ-571 owns it; `lifecycleTimestampFields` already serves), `web/board.css` and `web/board-detail.js` (row highlighting and drawer opening are REQ-573), `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` (release paths, written by finalization), everything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] A ticket with `created_at` and `claimed_at` yields two rows, "claimed" then "captured", newest first (captured RED/GREEN)
- [ ] Every parseable stamp on every ticket becomes a row; a ticket with no parseable stamp is still skipped
- [ ] Ordering is newest first with a deterministic tie-break that also covers two stamps of one REQ at one instant
- [ ] The Activity summary line reports the transition count and the distinct REQ count; the window, filter chips and both empty states keep working

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/activity.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified, comments only)
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified, comment only)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** `buildActivityRows` now appends one row for every parseable stamp `lifecycleTimestampFields` returns, sorted newest first with `RequestId` then `StampField` as tie-breaks; `newestLifecycleStamp` was deleted. The Activity summary reports transitions and distinct REQs ("3 transitions across 2 REQs in the last 24 hours"), both empty states count transitions, and every comment that said "one row per REQ" was restated. Payload shape unchanged. Merge range `7ad53bff..fbdcd35e`, builder branch head `2d8beb40`. Builder-authored `## Decisions` (D-01 to D-07) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-04-232225/REQ-572-handback.md`.

## Qualification

**Passed.** Merge range `7ad53bff..fbdcd35e` holds exactly the six Scope files; `qualify` and `scope-drift` both returned satisfied on the second run (the first run reported the three P-A-U boxes unchecked, because the builder cannot write the REQ file, and two Scope prose tokens the checker read as paths; both were orchestrator bookkeeping, fixed and re-run).

Read directly from the diff, not from the summary:
- `activity.go`: the per-ticket `newestLifecycleStamp` call is replaced by a loop over `lifecycleTimestampFields(ticket)` that appends one `ActivityRow` per parseable stamp; the helper is deleted; the comparator gains `StampField` ascending as the third key with a comment saying why the direction is arbitrary. No new stamp list; `model.go` untouched.
- `board-activity.js`: the summary builds a distinct-id set over the filtered rows and reports "N transitions across M REQs in the <window>" with real singular forms; the before-filters clause and both empty states count transitions. Row construction and windowing are unchanged, so the wall-clock window and the shared filter chips keep working; `data-activity-request` stays on every row with a comment that it is deliberately non-unique.
- `generate.go` and `template.html`: comment restatements only; `git diff` shows no code line changed in either.
- Tests: one new Go test for the captured two-row case, three rewritten assertions that keep their fixtures, one extended tie test with a mutation check recorded in the hand-back, and a new Node lane test pinning the summary and both empty states against the shipped assembled client.

Requirements traced: one row per stamp (loop), newest first with deterministic ties (comparator), skip unparseable stamps (`parsedOk` guard), both counts in the summary (distinct-id set), window and filters unchanged, payload shape unchanged (`generatedActivityEntry` untouched). Live data flow verified by the builder's generate run against this repository: REQ-570 carries five activity rows in `board-data.js` against one before.

*Checked by work action*

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** go-test-file-budget:do-work-cli:publication-defer-gate-test
- **Repair dependency:** REQ-574
- **Diagnostic evidence:** "post-merge run at 4adcff4e (fbdcd35e merged): FAIL: internal/corehelpers/inventory_test.go accumulated 38.92s; internal/publication/defer_gate_test.go 37.01s; internal/finalization/finalization_recovery_test.go 35.65s; internal/finalization/finalization_req499_test.go 30.85s; each test file must finish under 30s; every test passed"
- **Diagnostic evidence:** "pre-build run at f6c43d22: the same four files over budget (38.61s, 37.93s, 35.73s, 33.72s)"
- **Diagnostic evidence:** "detached diagnostic worktree at base 7ad53bff (clean tree): FAIL: internal/publication/defer_gate_test.go accumulated 32.52s; each test file must finish under 30s; every test passed; queue-kanban 24s"
- **Implementation base:** 7ad53bff1d867f1453e1e7765e988dedb308e7e1
- **Implementation merge:** fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
