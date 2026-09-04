---
id: REQ-568
title: 'Show recently touched REQs on the board regardless of status'
status: completed
created_at: 2026-09-04T17:54:05Z
user_request: UR-112
domain: ui-design
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-04T18:13:19Z
route: C
write_set:
  - skills/do-work-board/tools/queue-kanban/activity.go
  - skills/do-work-board/tools/queue-kanban/activity_test.go
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/board-filters.js
  - skills/do-work-board/tools/queue-kanban/web/board.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - CHANGELOG.md
  - skills/do-work/CHANGELOG.md
  - VERSION
estimate:
  p50_active_minutes: 50
  confidence: low
  calculated_at: 2026-09-04T18:14:58Z
  basis:
    - Route C
    - 8-file write set
    - 1 new files
    - 3 subsystems involved
    - 4 acceptance criteria
    - browser evidence
    - cross-route regression gates
completed_at: 2026-09-04T20:05:02Z
commit: 6037cd5
release_at: 2026-09-04T20:05:02Z
---

# Show Recently Touched REQs on the Board Regardless of Status

## What

Give the Kanban board one surface that answers "what changed on the queue in the last N hours, and why", listing every REQ whose newest lifecycle stamp falls inside the selected window, newest first, with the stamp and the transition it records. Status must not filter it: a REQ that was claimed, held for heavy testing, deferred, blocked, completed, cancelled, or failed inside the window all belong on it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, the shipped `prime-do-kanban.md`, `_dev/primes/lessons-kanban-board.md`, and `crew-members/ui-design.md`. Approach recorded in `## Plan` above: extract the lifecycle-stamp list to one definition, aggregate the newest stamp per REQ in Go, ship rows unwindowed, window in the client, host it as a sixth view tab.
- [x] **[APPLY]:** Ten files, all inside the declared scope; `tools/checks/scope-drift.sh` reports no undeclared touch.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file. `go vet ./...` clean; `gofmt -l` empty; `node --check` clean on all three touched JS files. No debug artifacts: no `console.log`, `debugger`, `fmt.Print`, `t.Skip` or `TODO` in the diff. Files verified — activity.go (aggregation + comparator + skip rule), activity_test.go (six lock-ins), model.go (extraction, both readers still compile against it), generate.go (payload struct, row build, fragment manifest entry), board-activity.js (window filter, table render, two empty states), board-controls.js (panel map, window group visibility, click wiring, render guard), board-filters.js (invalidate + re-render on filter change), board.js (view state comment, window default, render guard), board.css (view styles only, appended), template.html (view button, window group, panel markup).

## Why

On the board generated 2026-09-04 16:56 UTC, the only claimed card (REQ-505, moving selection and claim behind advance) was 20 minutes old and the newest Recently done card (REQ-485, canonicalizing reservation marker filenames) was two hours old. The maintainer read that as a gap and asked where the work went. Git history showed the loop never paused: REQ-567 (repairing shipped lesson links to archived UR paths), REQ-503 (adding the read-only advance lifecycle command), and REQ-504 (collapsing Step 10 and crash recovery prose into recovery) were each claimed, built, merged, and held as `pending-heavy-testing` between 14:57 and 16:39 UTC. Those three sit under Pending, Waiting, so no "recent" surface on the board showed them. Answering the question needed `git log`.

## Context

- The board has three time surfaces today, none of which answers the question: the Recently done column (terminal states only, 24h / 48h / 7d window), the Timeline view (wait and work spans per REQ, window buttons for last day / 7 / 30 / 90 / all days), and the Calendar view (one entry per REQ on its claim or resolve day). The `open-work` terminal digest excludes finished REQs on purpose.
- Tickets already carry several lifecycle stamps the board parses: `created_at`, `claimed_at`, `completed_at`, `blocked_at`, the REQ-448 phase milestones (`planning_at`, `exploration_at`, and the rest), and `testing_updated_at`. The `pending-heavy-testing` hold itself currently writes no stamp: REQ-503's frontmatter after the hold carries `claimed_at` and the phase milestones but nothing recording when it was held. The builder decides whether the newest-stamp rule is enough or whether the hold transition must also stamp a time; the REQ intent is that the hold event is visible on this surface either way.
- "Newest stamp on a ticket" is the proposed key, and the transition it belongs to is what the row should name in words (claimed, held for heavy testing, blocked, completed, and so on). Which control hosts the surface (a new toolbar tab, a strip like Verify findings, or reuse of the existing Recently done window buttons) is the builder's call; the maintainer's words were "a 'recently touched' window".
- Release commits (for example 0.275.3 at 15:54 UTC in the same gap) are not REQs and are out of scope for this surface unless the builder finds them free to include.
- Board versioning, parser lock-step, and build outputs follow `_dev/primes/prime-kanban-board.md`.

## Red-Green Proof
**RED prompt/case:** Generate the board against this repository at a state where REQ-503, REQ-504, and REQ-567 are `pending-heavy-testing` with claims inside the last two hours and REQ-505 is `claimed`, then look for any surface listing what changed in the last 2 hours.
**Why RED now:** Recently done shows only REQ-485 and older terminal REQs. The three held REQs appear only under Pending, Waiting, with no time ordering and no transition name. The Timeline draws their open bars but not the hold. Nothing on the board says "REQ-504 was held for heavy testing at 16:38".
**GREEN when:** One board surface lists REQ-505, REQ-504, REQ-503, REQ-567, and REQ-485 newest first for a "last 24h" window, each row naming its newest stamp and the transition it records, and a Go test on the aggregation pins that a `pending-heavy-testing` REQ with a hold inside the window is included and ordered by that stamp.
**Validation:** Inferred during capture from the maintainer's accepted proposal.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-kanban-board.md` — 4820 tokens; matches a queue-kanban view change, but the satellite is `slugged: partial`, so no targeted form is legal and the bare entry exceeds the 2000-token capture budget.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5744 tokens; matches queue-kanban model, UI, and timeline behavior, but the satellite is `slugged: partial`, so no targeted form is legal and the bare entry exceeds the 2000-token capture budget.

## Assets

None.

---
*Source: "capture that as a REQ" (accepting the proposal: a "recently touched" window keyed on the newest stamp on a ticket, `updated_at` or the hold time)*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new board surface needs an aggregation over every lifecycle stamp in the model, a new UI control and view in the embedded client, a Go test pinning the ordering, and possibly a new hold stamp written by the work loop. That spans the Go model, the embedded frontend, and the skill's own action prose, so it is multi-component rather than a located edit.

**Planning:** Required

## Plan

### The two builder decisions the REQ hands over

**D-01 — the hold already has a stamp; no new frontmatter field.** `actions/work.md` Step 6.5's heavy-testing hold is specified to write `status_changed_at` alongside the status flip, and `status_changed_at` is already parsed onto the ticket (`model.go` `StatusChangedAt`) and already listed in `detectFutureTimestampFields`. REQ-567 carries it. REQ-503 and REQ-504 do not, and REQ-504 carries no lifecycle stamp at all past `created_at` — that is missing data on two tickets, not a gap in the schema. Adding a second hold stamp would be a second way to say what `status_changed_at` already says, and the first thing to drift. The two unstamped tickets go to `## Discovered Tasks`.

Consequence, stated plainly: REQ-504 cannot appear on a stamp-driven surface, so the captured GREEN's five-ticket list is met for REQ-505, REQ-503, REQ-567 and REQ-485 but not REQ-504. Deriving "touched" from git history instead would invent a second definition of the word and read a different source than every other time surface on this board. Not done.

**D-02 — a new Activity view, not a strip.** The two existing strips (Completion anomalies, Verify findings) are for breakage, and their defining property is that no window and no filter may hide them. This surface is windowed and should honour the shared filters, so putting it in a strip would need an exception to the one rule that makes strips trustworthy. It is a table of rows, and tables on this board live in view panels (Timeline, Durations). So: a sixth entry in the View control group, with its own window control group.

### The aggregation

New `activity.go` + `activity_test.go`, following the `timeline.go` / `durations.go` shape: pure functions over already-parsed tickets, no second walk of the tree and no new frontmatter field.

The set of lifecycle stamps is **not** re-enumerated. `model.go`'s `detectFutureTimestampFields` already holds the one list of every timestamp the frontmatter can carry; a second hand-maintained copy would go stale the first time a stamp is added (CLAUDE.md, *State conditions, not lists*). Extract that table into one `lifecycleTimestampFields(ticket)` returning `{FieldName, RawValue, Transition}`, and give it two readers: the future-stamp check (which reads two fields) and this aggregation (which reads all three).

`buildActivityRows(tickets)` returns one row per ticket that has at least one parseable stamp: the newest parseable stamp, the field it came from, and the transition that stamp records in words. Ordered newest first with an explicit `RequestId` tiebreak, so a tie is decided by the comparator rather than by the sort's stability.

Two transition names depend on the ticket's status rather than the field alone — `status_changed_at` on a `pending-heavy-testing` ticket reads "held for heavy testing" (matching the calendar's existing wording), and `completed_at` reads "completed", "cancelled" or "failed" by status. Everything else is fixed per field.

The window is not applied in Go. The payload ships every row with its stamp, and the client filters against the wall clock at render time — the same shape `recentlyDoneIds` already uses in `board-cards.js`, and for the same reason: a tab left open past the snapshot must not keep counting "last 24 hours" from page-generation time.

### The client

- `web/board-activity.js` (new fragment; register it in `boardJavaScriptFragmentPaths` in `generate.go`, which `TestBoardJavaScriptAssemblyStructure` locks).
- `web/template.html` — a sixth View button, a `#view-activity` panel with its own window control group (6h / 24h / 48h / 7d) and a table: REQ, Title, Status, When, Transition.
- `web/board-controls.js` — panel registration, window-group visibility, first-render guard.
- `web/board.js` — `viewState.activityWindowHours`, `renderedOnce.activity`.
- `web/board.css` — reuse the existing table classes; add only what the new panel needs.

### Test-first sequence (`tdd: true`)

1. `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` — RED first. A `pending-heavy-testing` ticket whose `status_changed_at` is inside the window is present, carries the "held for heavy testing" transition, and sorts above an older `completed_at` ticket and below a newer `claimed_at` one. This is the captured GREEN's Go half.
2. `TestBuildActivityRowsPicksTheNewestStampNotTheFirstField` — a ticket whose newest stamp is not the first field in the table.
3. `TestLifecycleTimestampFieldsFeedsBothReaders` — the extracted list is what `detectFutureTimestampFields` reads, so the two cannot drift apart.
4. Boundary pair straddling the window edge, derived from the constant rather than restated beside it (lesson REQ-374: a fixture that spans a threshold widely does not test it).
5. Render the board and look at it (prime: a chart's correctness is partly a claim about pixels; lesson REQ-285: for a rendering change the screenshot is the test).

### Release

New user-visible view, so a minor bump: 0.275.3 → 0.276.0, root `CHANGELOG.md` entry plus the byte-identical `skills/do-work/CHANGELOG.md` mirror (`_dev/primes/prime-releases.md`).

*Plan written by the work action after the dispatched Plan agent was lost to a session restart.*

## Plan Validation

- **Requirement coverage:** every requirement in What maps to a task — "one surface" → the Activity view; "newest lifecycle stamp" → `buildActivityRows`; "newest first" → the comparator; "the stamp and the transition it records" → the row's two columns; "status must not filter it" → the aggregation reads every ticket regardless of status, which test 1 pins for the held case.
- **No orphan tasks:** the `lifecycleTimestampFields` extraction is the only task not named by the REQ. It is not scope creep: it is what stops this REQ from creating a second stamp list, and it deletes a duplicate rather than adding one.
- **Scope sanity:** 4 tasks (extract the list, aggregate, payload, client). Under the 5-task flag.
- **Consumer field contract:** the payload's `activity` rows carry id, stamp field, stamp instant and transition — the exact fields the client's table and window filter consume.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/activity.go` (new) — the aggregation: newest lifecycle stamp per REQ, its transition, newest first
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (new) — the Go lock-ins, including the captured GREEN's held-REQ case
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — extract the lifecycle-stamp list so it has one definition and two readers
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — the activity payload rows plus the new client fragment in the assembly manifest
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (new) — windowing and drawing
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modify) — panel registration, window group, first-render guard
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modify) — refresh the view when the shared filters change
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify) — the activity window state and the render guard
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the view's styles
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the View button, the window group, the panel and its table
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` (modify) — the release

**Files I will NOT touch:** the work action and its reference under skills/do-work/actions/. The heavy-testing hold already writes the status_changed_at stamp (D-01), so no action prose changes; adding a second hold stamp would have pulled them in.

**Acceptance criteria (restated from REQ):**
- [x] One board surface lists every REQ whose newest lifecycle stamp falls inside the selected window
- [x] Newest first
- [x] Each row names the stamp and the transition it records
- [x] Status does not filter it — claimed, held, blocked, completed, cancelled and failed all belong
- [x] A Go test pins that a `pending-heavy-testing` REQ with a hold inside the window is included and ordered by that stamp

**Ordering note, recorded rather than hidden:** this section was written after implementation, not before it. Step 5.5 wants it first and this run did not do that. The declared list and the Implementation Summary list agree exactly, so there is no drift to find, but the declaration did not do the job it exists for on this REQ.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/activity.go` (new)
- `skills/do-work-board/tools/queue-kanban/activity_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-activity.js` (new)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)

**What was done:** Added an Activity view to the board. `buildActivityRows` picks each REQ's newest parseable lifecycle stamp and the transition it records, ordered newest first with an explicit id tiebreak; a REQ with no parseable stamp is skipped rather than dated from zero. The stamp enumeration that lived inline in `detectFutureTimestampFields` became `lifecycleTimestampFields`, now carrying a transition phrase per field and read by both that check and the new aggregation, so the list has one definition. The rows ship unwindowed in the `activity` payload and `board-activity.js` filters them against the wall clock at render time, matching how `recentlyDoneIds` already behaves. The view is a sixth tab with its own 6h/24h/48h/7d window, honours the shared filter chips, and distinguishes "nothing moved" from "your filters hide what moved".

## Decisions

- **D-01**: The heavy-testing hold needs no new frontmatter field. `status_changed_at` is already written by the hold, already parsed onto the ticket, and already in the future-stamp list. Reasoning: a second stamp saying the same thing is the first thing to drift. DECIDE & STATE — reversible, and it removes a schema change rather than adding one.
- **D-02**: The surface is a view tab, not a strip. Reasoning: the two existing strips are for breakage and their defining property is that no window and no filter can hide them; this surface is windowed and filterable, so hosting it there would need an exception to the one rule that makes strips trustworthy. DECIDE & STATE — reversible, and it adds no rule exception.
- **D-03**: The Go aggregation ships every row and the client applies the window, rather than Go pre-windowing. Reasoning: it matches `recentlyDoneIds`'s existing shape, and it is what keeps "last 24 hours" true in a tab left open past the snapshot. DECIDE & STATE.
- **D-04**: The captured GREEN lists its five REQs as 505, 504, 503, 567, 485, but the stamps the same REQ cites put REQ-504's 16:38 hold above REQ-505's 16:36 claim, and REQ-567's 15:19 hold above REQ-503's 15:04. The test encodes the stamp order. Reasoning: the capture's listing was approximate and the surface's whole contract is that stamps decide the order. DECIDE & STATE.
- **D-05**: Repaired `skills/do-work-board/justfile.template` and the repository's own `justfile` outside this REQ's original scope, on the user's explicit instruction, because the canonical repository gate could not go green without it. Two recipes used `name default="x" *args:`, which no `just` version parses. Recorded here because it is a scope expansion, committed separately from this REQ's work.

## Discovered Tasks

- `[impact-negligible]` `_dev/tests/session-start-hook-behavior.sh` runs 25s standalone and 51s under a cold full-gate run, against the gate's own 30-second per-test-file budget. It passed at 29s on a warm re-run. The budget has no headroom on a slower machine, so this file will flake the gate rather than fail it honestly. → report only
- `[impact-negligible]` No `just` version is pinned anywhere in the repository, and the managed template needs 1.29 or newer for `[positional-arguments]`. A pin (or a version check in `maintainer-verify.sh` beside the existing Go and ShellCheck checks) would have caught D-05's defect at the point it was introduced rather than at the next gate run that happened to have `just` installed. → report only

## Qualification

Passed — 10 files verified against the merge range `76d74a0..dc6dc52`, 5 acceptance criteria traced, P-A-U confirmed against the diff.

Two qualification findings were fixed rather than waved through:
- `gofmt -l` flagged `generate.go` after the payload field was added. The UNIFY note had claimed gofmt was clean before it was run; the file was reformatted and the claim now matches the check.
- `TestBoardJavaScriptAssemblyStructure` failed on the new client fragment. That is the lock working as designed (`_dev/primes/prime-kanban-board.md`: the manifest in `generate.go` is the source and the test pins it), so both its authored-inventory and execution-order lists were updated in the same change.

Substantive-changes check: `activity.go` is 104 lines of aggregation with no placeholder returns; `board-activity.js` renders real rows from the payload and was verified doing so in a browser, not merely by reading it.

## Testing

**Tests run:** `go test .` in `skills/do-work-board/tools/queue-kanban` (whole package, 59s), `go vet ./...`, `gofmt -l .`, `node --check` on each touched fragment, and the canonical repository gate `bash _dev/tests/maintainer-verify.sh`.
**Result:** ✓ All passing.

**Red-green validation:** traces to the captured `## Red-Green Proof`.
- `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition`: ✗ before implementation (compile failure — `buildActivityRows` and `ActivityRow` undefined) → ✓ after. This is the captured GREEN's Go half: a `pending-heavy-testing` REQ with a hold inside the window is included, carries "held for heavy testing", and orders by that stamp rather than by its older claim.
- The captured RED case was also reproduced against the real repository rather than only in fixtures: the board generated before this change had no surface listing the last N hours, and the board generated after it lists 44 REQs for the last 24 hours, with the GREEN's five named REQs present in stamp order — REQ-505 17:31, REQ-504 16:37, REQ-503 15:43, REQ-567 15:19, REQ-485 14:56.

**New tests added:**
- `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` — the captured GREEN.
- `TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField` — pins "newest", not "first present field"; `created_at` is first in the list and nearly universal, so a first-match reader would pass every other assertion here.
- `TestBuildActivityRowsStraddlesTheWindowBoundary` — both samples derived from the window constant, one minute either side, so a second wrong cutoff cannot classify the pair alike (lesson REQ-374).
- `TestBuildActivityRowsSkipsTicketsWithNoParseableStamp` — absence stays absence rather than becoming a zero-time row.
- `TestBuildActivityRowsBreaksStampTiesDeterministically` — runs the same tie in both input orders, so the sort's stability cannot answer for the comparator (lesson REQ-318).
- `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` — the anti-drift pin across both readers of the extracted list.

**Existing tests updated (cross-REQ impact):**
- `generate_test.go` (`TestBoardJavaScriptAssemblyStructure` and the authored-inventory assertion): the new client fragment is a deliberate manifest change, and this lock exists to make it deliberate.

**Render evidence:** static board generated from this repository at `2026-09-04 19:57 UTC` and measured in Chromium via Playwright, `location.href` returned alongside every measurement. Light and dark both read against the real painted `body` background (light `rgb(245,247,250)`, dark `rgb(12,14,18)`); no horizontal page overflow; the window buttons move the row count as expected (6h → 11 rows, 24h → 44, 7d → 164).

**Thirty-second file budget:** the queue-kanban package is one Go test binary at 59s, which is the package's existing shape and not a per-file measure; the gate's own budget reporter accepts it. `activity_test.go`'s own six tests run in under 0.02s combined.

## Review

**Requirements check — did we build what was asked?** Yes, and verified against the real queue rather than only fixtures. Every clause of the What section is met: one surface, every REQ whose newest lifecycle stamp is inside the window, newest first, each row naming the stamp and the transition, and status not filtering it (the rendered 24h view carries `claimed`, `pending-heavy-testing` and `completed` rows side by side).

The captured GREEN named five REQs. All five appear, in stamp order. One deviation from the capture's literal wording, recorded as D-04: the capture listed them 505, 504, 503, 567, 485, and the stamps the same REQ cites put 504 above 505 and 567 above 503. The surface's contract is that stamps decide the order, so the stamps won.

**One capture-time premise turned out to be wrong, in the REQ's favour.** The Context section says the `pending-heavy-testing` hold "currently writes no stamp" and cites REQ-503 as carrying none. That was true of the tree this REQ was captured against. After merging `origin/main`, REQ-503 and REQ-504 both carry `status_changed_at` from their holds, so the builder decision the REQ hands over ("whether the hold transition must also stamp a time") resolves to *no change needed* rather than to a schema addition. The surface therefore meets the full GREEN instead of the partial one an earlier reading of this run predicted.

**Code review.** The aggregation is a pure function over already-parsed tickets, matching `timeline.go` and `durations.go`, with no new frontmatter field and no second tree walk. Three judgment calls carry their reasoning in comments where the next reader will need it: why absence is skipped rather than zero-dated, why the tie is broken in the comparator, and why the window is applied client-side. The one structural change outside the new files — extracting `lifecycleTimestampFields` — removes a duplicate rather than adding an abstraction, and is pinned by a test that fails if either reader stops using it.

**Restatement sweep.** This REQ redefines nothing other text restates. It adds a payload field and a view; it does not change a contract token, a schema field's semantics, a gate's wording, or a prescribed command's output shape. `prime-do-kanban.md`'s "three write surfaces" count is untouched — the Activity view is read-only. No stale restatement found. → report only.

**Acceptance testing — does it actually work?** Yes, measured in a browser against a board generated from this repository, in both themes, with `location.href` returned alongside every measurement.

**Findings:**
- `impact-negligible` — the Activity table sits below the data-warnings banner and the completion-anomalies strip, both of which render above every view panel by design, so on this repository (9 anomalies) the view needs a scroll before its first row. Pre-existing layout behaviour shared with Timeline and Durations, not introduced here. → report only.

**Acceptance: Pass. Overall: 92%.**

## Lessons Learned

**What worked:** Extracting the stamp list before writing the aggregation, rather than copying it. The copy would have been three lines shorter and would have gone stale at the next schema addition; the extraction is pinned by a test that fails if either reader drifts off it. Also: generating the real board and reading it. The 44-row rendered view is what proved the surface answers the maintainer's actual question, and no assertion in the package could have told me that.

**What didn't:** The first fixture for the captured GREEN used round offsets (-2h, -20min) instead of the instants the REQ itself cites, and its expected ordering was simply wrong — REQ-485's completion landed newer than two of the holds. The test failed on my arithmetic, not on the code. Rebuilding the fixture from the REQ's own quoted times (16:38, 16:36, 15:19, 15:04, 14:56) made it both correct and checkable against the source. A fixture whose numbers come from the REQ can be audited; one whose numbers come from convenience cannot.

**Worth knowing:** `status_changed_at` is the hold's stamp and was already contracted — REQ-568's Context says otherwise because it was captured against a tree where two tickets had lost theirs. Check the current tree before treating a capture-time observation about missing data as a schema gap. Separately: the shipped `justfile.template` had two recipes no `just` version could parse, and the gate test that catches it only bites where `just` is installed, so a missing optional tool hid a real defect on `main` for a full release.

## Orientation

The board grows a sixth view, **Activity**, answering "what changed on the queue in the last N hours, and why" — every REQ's newest lifecycle stamp and the transition it records, newest first, with status deliberately not filtering it. Lives in the queue-kanban board subsystem (`_dev/primes/prime-kanban-board.md`), as a Go aggregation beside `timeline.go` and `durations.go` plus one new client fragment.

[MAP CHANGED] The lifecycle-stamp enumeration moved from being an inline detail of the future-stamp check to a named shared concept, `lifecycleTimestampFields`, that now carries a transition phrase per stamp. Anything added to the REQ timestamp schema from here reaches both the future-stamp warning and the Activity view from one edit, and a new stamp must now say what it records rather than only when.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` and the shipped `prime-do-kanban.md` were re-read against this change. Every path each names still exists. `prime-do-kanban.md`'s "three write surfaces" count and its Read-first list remain accurate — the Activity view reads only, and `activity.go` is a peer of the aggregations already described rather than a new routing concept. Neither prime is stale.
