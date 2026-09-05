---
id: REQ-572
title: 'Show every lifecycle transition of a REQ as its own Activity row'
status: completed
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
dispatch_at: 2026-09-05T12:06:43Z
builder_handback_at: 2026-09-04T23:41:02Z
integration_at: 2026-09-04T23:41:59Z
gate_deferred: 'true'
depends_on: [REQ-574]
claimed_at: 2026-09-05T12:00:55Z
review_at: 2026-09-05T12:21:17Z
commit: fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
heavy_verified_at: 2026-09-05T13:16:11Z
heavy_verified_revision: 9233921395b509d06d440f955e0bdb0c289958bf
completed_at: 2026-09-05T13:16:28Z
release_at: 2026-09-05T13:16:28Z
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

## Saved-Range Resume Proof — Rejected (2026-09-05)

The saved pair `7ad53bff..fbdcd35e` resolves and is an ancestor of HEAD, but three later commits touch protected implementation paths (`a55f24ce` and `4c76c332` for REQ-571, `ae184a7b` for REQ-576). Proven path drift, so reuse is rejected: the two saved pointers were deleted from the frontmatter and every prior qualification, testing and review claim below is superseded. This REQ returns to Step 6 implementation against current main.

## Prior Implementation Summary (superseded — saved-range drift)

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/activity.go` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified, comments only)
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified, comment only)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** `buildActivityRows` now appends one row for every parseable stamp `lifecycleTimestampFields` returns, sorted newest first with `RequestId` then `StampField` as tie-breaks; `newestLifecycleStamp` was deleted. The Activity summary reports transitions and distinct REQs ("3 transitions across 2 REQs in the last 24 hours"), both empty states count transitions, and every comment that said "one row per REQ" was restated. Payload shape unchanged. Merge range `7ad53bff..fbdcd35e`, builder branch head `2d8beb40`. Builder-authored `## Decisions` (D-01 to D-07) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-04-232225/REQ-572-handback.md`.

## Prior Qualification (superseded — saved-range drift)

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

## Implementation Summary

**Files changed:** None in this dispatch. Every file this REQ declares already carries its change from the earlier merge `fbdcd35e`, which is still an ancestor of `HEAD`:

- `skills/do-work-board/tools/queue-kanban/activity.go` (unchanged — per-stamp loop, three-key comparator, `newestLifecycleStamp` absent)
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (unchanged — seven Activity tests including the captured RED/GREEN case)
- `skills/do-work-board/tools/queue-kanban/generate.go` (unchanged — comments only; payload shape untouched)
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (unchanged — distinct-id set, both counts, both empty states)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (unchanged — Activity section comment)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (unchanged — Node-lane summary test)

**What was done:** The saved-range resume proof was rejected on path drift, so the REQ returned to implementation. The builder re-derived the state of all six files against current main instead of trusting the prior record, found every acceptance criterion already satisfied, and committed nothing. The builder branch `worktree-agent-REQ-572-activity-rows` sits at its base commit with zero commits of its own; there is no hand-back range to merge. The delivered behavior lives in the merge range `7ad53bff..fbdcd35e`, which every evidence step below reads. Builder-authored `## Decisions` (D-01 to D-05) and `## Discovered Tasks` live in `do-work/runs/work-2026-09-05-120117/REQ-572-handback.md`.

## Qualification

**Passed.** Read from the tree and the range `7ad53bff..fbdcd35e`, not from the prior summary.

- `activity.go:66-102` — one `ActivityRow` appended per parseable stamp from `lifecycleTimestampFields(ticket)`, the `parsedOk` guard skipping unparseable values with no zero-time fallback, and a comparator whose keys are stamp time (descending), `RequestId` (descending), then `StampField` (ascending) with a comment saying the last direction is arbitrary and only determinism matters. `grep -rn newestLifecycleStamp` returns nothing across the module, so the helper is gone rather than shadowed.
- `web/board-activity.js:62-80` — a distinct-id set over the filtered rows and a summary reading "N transitions across M REQs in the <window>" with real singular forms and a "(K before filters)" clause. Windowing still runs off the wall clock and the shared chip filter still applies, so neither behavior regressed.
- `generate.go` and `template.html` carry comment restatements only; `generatedActivityEntry` keeps `id`/`stampField`/`stampAt`/`transition`, so REQ-573 (opening the drawer from an Activity row) still has the payload contract it was promised.
- The three commits that voided the saved range were read rather than assumed harmless: `a55f24ce` and `4c76c332` (REQ-571) change ten lines of doc comment and test-fixture vocabulary in `activity.go`, `activity_test.go` and `javascript_behavior_c_test.go`; `ae184a7b` (REQ-576) touches one comment inside a different struct in `generate.go`. No aggregation logic, comparator or assertion structure moved.

Requirements traced: one row per stamp (the loop), newest first with a tie-break that separates two stamps of one REQ at one instant (the third key), unparseable stamps still skipped (the guard), both counts in the summary (the distinct-id set), window/filters/empty states unchanged, payload shape unchanged. Live data flow verified end to end: a real board generated against this repository's own queue holds 2047 activity rows across 561 REQs with 549 REQs carrying more than one row, REQ-570 shows all nine of its transitions where it previously showed one, every adjacent pair is newest-first, and no row is dated from the zero time.

An empty diff is accepted here only because the evidence proves the work is present and pinned, not because the tests pass. See the RED observation in Testing.

*Checked by work action*

## Testing

**Focused tests (builder, in its worktree at `09a13839`):**
- `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` — exit 0, 392 tests, wall 45s, slowest file `generate_test.go` at 9.06s against the 30s per-file budget.
- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1` on the same wrapper with `-run '^TestJavaScriptBehavior'` — exit 0, 56 tests, wall 7s, slowest file `javascript_behavior_c_test.go` at 2.12s.
- The eight Activity tests re-run verbosely to prove none was skipped: all eight `--- PASS`, with the Node-lane summary test taking 1.80s, which is a real `node` execution rather than a skip.

**Red-green validation** (traced to `## Red-Green Proof`): the captured RED could not be observed in its original order, because the implementation arrived before this dispatch. The builder produced the equivalent by restoring the pre-change behavior in its worktree and re-running: exit 1 with six failures, every one on an assertion rather than a compile error, including `activity_test.go:32: rows = 1, want 2` on the captured REQ-570 case and `javascript_behavior_c_test.go:2492: summary ... = "3 REQs touched in the last 24 hours", want both counts`. The two window and skip tests stayed green under the revert, which is the correct shape: the tests that pin this REQ fail without it, the neighbouring ones do not. The tree was restored with `git checkout --` and `git status --porcelain` is empty. This adaptation is recorded as builder decision D-02.

**Canonical repository gate:** `bash _dev/tests/maintainer-verify.sh` — see the run-level record; the gate was red twice at the claim revision on three stale archive links no claimed REQ owns, fixed in place at `09a13839` (release 0.294.1), and green there.

## Review

**Overall: 97%** | 2026-09-05T12:19:29Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Approve.** The empty diff is correct. Every acceptance criterion is present in current main, and the evidence supports completion rather than merely asserting it — the reviewer reproduced the builder's RED experiment independently in a scratch copy of the module and confirmed the tests fail for the right reason without the change.

### Empty-diff judgment (the first thing the orchestrator asked)

The delivered behavior lives in `7ad53bff..fbdcd35e`, which `git merge-base --is-ancestor fbdcd35e HEAD` confirms is still in main. Nothing in the REQ is unbuilt. Verified from the tree, not the summary:

- `activity.go:65-103` — one `ActivityRow` per parseable stamp from `lifecycleTimestampFields(ticket)`, the `parsedOk` guard skipping unparseable values with no zero-time fallback, and a three-key comparator (stamp time descending, `RequestId` descending, `StampField` ascending). `grep -rn newestLifecycleStamp` over the whole repo outside `do-work/` returns nothing, so the old helper is gone, not shadowed.
- `web/board-activity.js:59-81` — distinct-id set over the filtered rows, summary "N transitions across M REQs in the <window>" with real singular forms and the "(K before filters)" clause. Windowing is still `Date.now()` (line 21); the shared chip filter is still `requestMatchesFilters` (line 61); both empty states count transitions (lines 128-132).
- `generate.go` and `web/template.html` — raw range diff read line by line: comment blocks only, no code line changed. `generatedActivityEntry` still ships `id`/`stampField`/`stampAt`/`transition`, so REQ-573 (clicking an Activity row to open the drawer and highlight sibling rows) keeps the payload contract it was promised.
- No "latest only" toggle was added; `#activity-window-group` still holds only the four window buttons.

**Independent RED reproduction (not taken from the hand-back).** The module was copied to a scratch directory and the per-stamp loop replaced with a newest-stamp-only pick; five Go tests failed on their assertions, not on compilation — `activity_test.go:32 rows = 1, want 2`, `:97 in-window rows = 5, want 8`, `:139 rows = 1, want 4`, `:231 rows = 1, want 4`, `:300 rows = 1, want 14`. Reverting the summary string in the same scratch copy failed `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` as the sixth. `TestBuildActivityRowsStraddlesTheWindowBoundary` and `TestBuildActivityRowsSkipsTicketsWithNoParseableStamp` stayed green under both mutations, which is the correct shape: the tests that pin this REQ fail without it and the neighbouring ones do not. The main checkout was never modified for this experiment.

### Restatement Sweep (Step 6)

The diff redefines what one Activity row means, so the sweep was run rather than skipped. Every consumer and restatement of "one row per REQ" was located and checked:

- Canonical home `activity.go` header and `ActivityRow` doc — restated. `generate.go:82-86` and `generate.go:312-317` — restated. `web/board-activity.js:1-14` and the render comments — restated. `web/template.html:463-472` — restated.
- `boardData.activity` has exactly one reader (`board-activity.js:22`); `board-controls.js` and `board-filters.js` only touch the render-once flag, not the row semantics.
- `data-activity-request` is set on every `<tr>` and is deliberately non-unique, with a comment saying so and a Node-lane assertion pinning the repeated id.
- **One stale restatement found**, recorded below as M1.
- `CHANGELOG.md:197` (the 0.276.0 entry) and the archived REQ-568 record describe what those versions shipped and are history — correctly left alone.
- `skills/do-work-board/docs/board-guide.md` and `prime-do-kanban.md` mention only the queue activity *calendar*, a different surface, so neither restates this contract.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- M1 — `web/template.html:473`: `aria-label="Recently touched REQs"` on `#view-activity` still names the region by the old unit. The comment directly above it was restated in the same diff hunk and the accessible name was missed, so a screen-reader user hears "Recently touched REQs" for a table that now lists transitions. — impact-user-visible → report only
- M2 — `do-work/runs/work-2026-09-05-120117/REQ-572-handback.md` `## Discovered Tasks`: the three lines carry `impact-process` and `impact-cosmetic`, tokens the impact vocabulary does not define (`impact-critical`, `impact-user-visible`, `impact-rule-change`, `impact-negligible`). The orchestrator routes on that token, so an undefined one cannot be routed. — impact-rule-change → report only
- M3 — Traceability of an empty-diff REQ: the `commit:` this REQ receives at finalization will point at a commit containing no implementation, so a later standalone review reading that hash hits the "nothing to review" exit. The delivered range `7ad53bff..fbdcd35e` appears in the REQ body but in no frontmatter field. — impact-rule-change → report only

**Nit findings:**
- N1 — `activity.go:66`: `make([]ActivityRow, 0, len(tickets))` is still sized one row per ticket. Real data is 2051 rows from 561 tickets, so the slice regrows about two times per build. — impact-negligible → report only
- N2 — The REQ's `## AI Execution State` still cites the superseded dispatch's hand-back (`work-2026-09-04-232225`) and its two builder commits; the current dispatch's hand-back is at `work-2026-09-05-120117`. Accurate about how the code arrived, stale about which dispatch certified it. — impact-negligible → report only
- N3 — `skills/do-work-board/docs/board-guide.md` never documents the Activity view at all. Pre-existing from REQ-568 (showing recently touched REQs regardless of status), not introduced by this diff. — impact-negligible → report only

### Requirements Checklist

- [x] `buildActivityRows` emits one row per parseable stamp, not one per ticket — delivered (`activity.go:71-82`)
- [x] `newestLifecycleStamp` gone; no second stamp list; `model.go` untouched — delivered (grep returns nothing; `model.go` absent from the range)
- [x] Newest first with a deterministic tie-break covering two stamps of one REQ at one instant — delivered (three-key comparator; `TestBuildActivityRowsBreaksStampTiesDeterministically` drives both arms)
- [x] A ticket with no parseable stamp is still skipped, never dated from the zero time — delivered (`parsedOk` guard; pinned by test and by zero unparseable rows on live data)
- [x] Client keeps wall-clock windowing and the shared filter chips — delivered (`Date.now()`, `requestMatchesFilters`)
- [x] Summary reports transitions and distinct REQs; empty states keep working — delivered (distinct-id set; Node lane pins five cases including both empty states)
- [x] No "latest only" toggle added; all transitions are the default — delivered
- [x] Payload shape unchanged for REQ-573 — delivered (`id`/`stampField`/`stampAt`/`transition` confirmed on generated output)

### Acceptance Testing

**Result: Pass**

- `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` — exit 0, 447 tests, wall 64s, slowest file `generate_test.go` at 12.52s against the 30s per-file budget.
- `gofmt -l .` empty; `go vet ./...` exit 0; debug-artifact scan over the six declared files returns nothing.
- End-to-end, run by the reviewer rather than read from the hand-back: `go run . generate` against this repository produced 2051 activity rows across 561 distinct REQs, 549 of them carrying more than one row; every adjacent pair is newest-first; zero rows carry an unparseable instant; the payload keys are exactly `id,stampField,stampAt,transition`. REQ-570 (deleting the pending-heavy-testing status) shows its whole path — `completed_at > release_at > review_at > integration_at > builder_handback_at > dispatch_at > planning_at > claimed_at > created_at` — where it previously showed one row. In a live 24h window the surface holds 236 transitions across 49 REQs, which is what the summary line will read.
- Payload cost measured rather than assumed: the activity array is 213 KB inside a 12.3 MB `board-data.js`, so the several-fold row growth the Exploration flagged is 1.7% of the file. Risk: None.
- Finding-Closure Ratchet: not applicable — this REQ is not review- or triage-finding-origin.

### Suggested Additional Testing

- Manual: open the Activity view in a browser at 24h and confirm 236 rows scroll and read sensibly with one REQ owning up to nine of them — the reviewer verified the data, not the visual density.
- Manual: check the near-duplicate pairs the Exploration accepted (`completed_at` "completed" beside `status_changed_at` "status changed to completed" at one instant) read as intended rather than as a bug, now that they are visible on real data.
- Accessibility: confirm the region's accessible name after M1 is decided — a screen reader currently announces "Recently touched REQs" over a transition table.
- Regression: REQ-573 (clicking an Activity row to open the drawer and highlight sibling rows) builds directly on the non-unique `data-activity-request`; exercise it against a REQ with nine rows, not two.

**Acceptance:** Pass — focused Go and Node lanes green (447 tests, exit 0), and a real board generated by the reviewer shows every REQ's full transition path with the payload contract intact.
**Suggested testing:** 4 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Rejecting the saved-range resume on drift and re-deriving the state of all six files cost one builder turn and produced better evidence than the mechanical proof would have: the builder read the three commits that voided the range and found them to be comment sweeps, and both builder and reviewer independently reproduced RED by reverting the behavior rather than trusting a green suite.

**What didn't:** Nothing was rebuilt, because nothing needed rebuilding. The saved-range proof rejects on *any* commit history against a protected path, which is the safe direction but treats a ten-line comment sweep exactly like a logic change — the cost is one full builder dispatch to discover the diff is empty.

**Worth knowing:** A test that passes on arrival proves nothing about whether it can fail. For work that lands before its own verification, the honest substitute for RED-before-GREEN is to revert the behavior, watch the assertions fail for the right reason, and restore — recorded here as builder decision D-02 and reproduced independently at review. Also: the Node behavior lane exits 0 when its probes are skipped for a missing `node` binary, so a verbose run proving `--- PASS` with a real duration is what makes an assertion about the client trustworthy.

## Orientation

The board's Activity view now shows a request's whole path rather than only where it stands: every lifecycle stamp becomes its own row, newest first, and the summary line counts both transitions and distinct requests. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), aggregation in Go and drawing in the assembled client. On this repository's own queue that is 2051 rows across 561 requests where the view previously showed 561. The payload shape is unchanged, so REQ-573 (opening the detail drawer from an Activity row) still has the contract it was promised. No prime was made stale: the kanban prime's referenced paths all still exist.

## Heavy Verification Plan

- **Base revision:** 7ad53bff1d867f1453e1e7765e988dedb308e7e1
- **Target revision:** fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
- **Planned at:** 2026-09-05T12:21:17Z, from `_dev/tests/heavy-lanes.json`

| Lane | Argv | Why it was selected |
| --- | --- | --- |
| `queue-kanban-javascript` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` | all six changed paths matched subtree `skills/do-work-board/tools/queue-kanban` |
| `queue-kanban-browser` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` | same subtree match |
| `staged-skills` | `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` | same subtree match |

No path was left uncovered by the manifest. The request stays `claimed` in `do-work/working/` with its `commit:` landed, so dependent work builds against the delivered source while the lanes wait for the queue-exhaustion drain.

## Heavy Verification Result

- **Target revision:** fbdcd35e0908aca6a01f554cc9b7fd7c85347a49
- **Execution revision:** 9233921395b509d06d440f955e0bdb0c289958bf
- **Run at:** 2026-09-05T13:16:11Z, from a detached worktree (the shared main tree carried another session's uncommitted work, which a lane result must not be attributed to)

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| `queue-kanban-javascript` | 0 | 8s | executed |
| `queue-kanban-browser` | 0 | 119s | executed |
| `staged-skills` | 0 | 43s | executed |

Every lane this request selected was present in the run, exited 0, and none was skipped. The browser lane executed against Google Chrome rather than skipping, which is what makes it evidence: an earlier pass reported it skipped for a missing browser, and a skipped lane is not a pass.

