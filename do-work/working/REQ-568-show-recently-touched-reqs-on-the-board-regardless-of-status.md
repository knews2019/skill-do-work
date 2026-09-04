---
id: REQ-568
title: 'Show recently touched REQs on the board regardless of status'
status: claimed
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
---

# Show Recently Touched REQs on the Board Regardless of Status

## What

Give the Kanban board one surface that answers "what changed on the queue in the last N hours, and why", listing every REQ whose newest lifecycle stamp falls inside the selected window, newest first, with the stamp and the transition it records. Status must not filter it: a REQ that was claimed, held for heavy testing, deferred, blocked, completed, cancelled, or failed inside the window all belong on it.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
