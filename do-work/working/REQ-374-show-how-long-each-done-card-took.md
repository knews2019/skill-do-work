---
id: REQ-374
title: 'Show how long each done card took'
status: claimed
claimed_at: 2026-08-26T13:17:00Z
created_at: 2026-08-26T13:02:22Z
user_request: UR-074
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
route: B
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-26T13:22:00Z
  basis:
    - Route B
    - 4-file write set
    - 1 new files
    - 2 subsystems involved
    - 6 acceptance criteria
    - browser evidence
write_set:
  - skills/do-work-board/tools/queue-kanban/durations.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Show How Long Each Done Card Took

## What

A card in the Kanban board's Recently Done column states when the work finished (`done Aug 26, 12:47 UTC · 9min ago`) but never states how long it took. Add the implementation span to that card: the time from when the builder started the REQ to when it landed in Done with a completed status.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The board already measures this span — `durations.go` computes `completed_at − claimed_at` for the Durations view — but a reader scanning the Done column has to open another view to learn whether a card took twenty minutes or half a day.

## Context

- The span is `claimed_at` → the completion instant the card already renders. `claimedAt` is already in the client payload (`generate.go`), so this needs no new frontmatter field and no second walk of the archive.
- The done line is built in `web/board-cards.js` under `options.showCompleted`, using `makeInstantWithRelativeNode` from `web/board-core.js`.
- The four-hour ceiling is the board's existing read-time rule (`durations.go`'s `analysisOutlierCeiling`, stated once in `skills/do-work/actions/estimate-reference.md` § Calibration). Reuse that rule; do not restate the number as a second definition.

## Detailed Requirements

- Every Recently Done card whose status is `completed` or `completed-with-issues` and whose `claimed_at` parses shows the implementation span alongside the existing done stamp.
- A span at or under four hours renders plainly, in the same relative-time vocabulary the board already uses (`34min`, `2h`).
- A span over four hours renders with a marker saying the session was likely paused, so an overnight pause is never read as hours of work.
- A negative span renders as a broken-stamp warning rather than a number. The board already treats reversed stamps as an anomaly worth surfacing.
- A card with no parseable `claimed_at` shows no duration at all — the rest of the card is unchanged.
- Cancelled cards, which share the Recently Done column, show no duration. The user's request is scoped to completed status.

## Constraints

- Display only. No new frontmatter field, no new write surface, no change to `model.go`'s parsed status vocabulary.
- The four-hour ceiling and the reversed-stamp rule have one definition each already. Read them; do not copy the constant into the client.
- Board changes get a normal skill version bump and a root `CHANGELOG.md` entry; the tool has no independent changelog.

## Red-Green Proof

**RED prompt/case:** Run the board (`do-work-board board`) against an archive holding a REQ with `claimed_at: 2026-08-24T10:05:00Z` and `completed_at: 2026-08-24T12:45:00Z`. Its Recently Done card reads `done Aug 24, 12:45 UTC · <relative>` and says nothing about the 2h40m the work took.
**Why RED now:** `web/board-cards.js` renders only the completion instant on a done card. The span exists in the payload (`claimedAt`) and is computed for the Durations view, but no card reads it.
**GREEN when:** That same card also reads its span — `2h40m` — and three neighbours prove the marked cases: a REQ spanning 18 hours renders its span with a paused marker, a REQ whose `completed_at` precedes its `claimed_at` renders a broken-stamp warning instead of a number, and a REQ with no `claimed_at` renders the done line exactly as it does today with no duration text.
**Validation:** User confirmed (the odd-span handling was put to the user during capture and the mark-them option was chosen).

## Assets

Screenshot of the board's Recently Done column, supplied with the request. It could not be persisted as a file in the capture session; this is the record.

The column header reads `RECENTLY DONE` with a count badge of `18` and a `Copy all` button. Seven cards are visible, each a white rounded panel with a green left edge:

- `REQ-1684` · `completed` — "[impact-user-visible] Convert the two remaining apps/admin/ native confirms to CmsConfirmDialog" · chips `frontend`, `UR-388`, `ROUTE B` · footer `done Aug 26, 12:47 UTC · 9min ago`
- `REQ-1685` · `completed-with-issues` — "[impact-rule-change] Close the preview-to-apply drift window across the four destructive card-mutation flows" · chips `cms`, `UR-389`, `ROUTE C`, `impact-rule-change` · footer `done Aug 26, 12:23 UTC · 33min ago`
- `REQ-1683` · `completed` — "Align the admin Import/Export done feedback with the extended feedback contract" · chips `frontend`, `UR-388`, `ROUTE B` · a `NEEDS REQ-1682` dependency row · footer `done Aug 26, 00:16 UTC · 12h ago`
- `REQ-1682` · `completed` — "[impact-rule-change] Extend the prime-cms-ux feedback contract to admin browser surfaces" · footer `done Aug 25, 23:58 UTC · 12h ago`
- `REQ-1681` · `completed` — "Audit e2e/prime-bowser.md — its instructions predate the layout migration" · footer `done Aug 25, 23:47 UTC · 13h ago`
- `REQ-1680` · `completed` — "Addendum: Import card confirm becomes a CMS dialog with an Import-vs-Restore fork" · footer `done Aug 25, 23:39 UTC · 13h ago`
- `REQ-1679` · `completed` — "Admin can delete a card — any card, admin-only, mapped level assets deleted too" · footer `done Aug 25, 23:27 UTC · 13h ago`

The footer line on every card is the surface this REQ extends: it carries the completion instant and its relative companion, and nothing about elapsed implementation time.

## Full Context
See `do-work/user-requests/UR-074/input.md` for complete verbatim input.

---
*Source: "So when we show the recently done cards, please also show the duration, how long it took since it was started until it is finished to implement that card to make it delivered. By making it delivered, I mean it was moved to the Done column, completed status."*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fully specified and the target files are named, but the existing conventions still need discovery — where the four-hour ceiling is read from on the client side, which test file holds the board's card-rendering probes, and how the done line is styled.

**Planning:** Not required

## Exploration

**F1 — the card's rendered completion instant is not always `completed_at`.** `resolveCompletionTime` (`model.go:1375`) falls back to the `commit:` hash's git committer date when `completed_at` is absent or unparseable, and reports which it used in `completionTimeSource` (`generate.go:190`). The Durations view refuses that fallback: `buildDurationAggregate` requires both `claimed_at` and `completed_at` to parse and skips the ticket otherwise (`durations.go:92-96`). Computing the card's span from `completionTime` would therefore print a duration for REQs the Durations view excludes, and the number would be a git commit delta rather than an implementation span.

**F2 — the reversed span is already detected Go-side and already badged on the card.** `detectCompletionAnomaly` (`model.go:1407-1421`) flags `completed_at` earlier than a parseable `claimed_at`, and `board-cards.js:246-256` renders the `anomaly` badge with that reason. No second broken-stamp channel is needed.

**F3 — the four-hour ceiling is Go-only, by convention.** `analysisOutlierCeiling` (`durations.go:30`) has three references, all Go, none in `web/`. The client never sees 4h or 240: it receives the already-applied verdict as `excludedReason` (`generate.go:264`, `:730`) and branches on the strings `"paused"` / `"reversed"` (`board-durations.js:1267`, `:1434`). Recomputing the ceiling in JavaScript would be a new mirror of a constant that has exactly one definition. The two existing Go-to-client lock-steps are source-text assertions with named guards: `TestFutureStampCauseClauseMatchesTheShippedClient` (`timestamp_test.go:332`) and `TestCalendarDayKeySentinelsMatchTheShippedClient` (`board_synthetic_test.go:442`).

**F4 — a span formatter already exists.** `formatDurationMinutes` (`board-durations.js:443`) renders minutes as `7.5 min` / `1h 59m` with correct-carry rounding, and is what the Durations table already draws. `board-core.js`'s `formatElapsedDuration` (`:128`) is now-anchored and its node variant carries `data-instant-ms`, which the 1s ticker (`board-core.js:236`) rewrites — wrong for a fixed historical span.

**F5 — the client is one concatenated closure, not modules.** `assembleBoardJavaScript` (`generate.go:826-871`) joins the `web/*.js` fragments into a single IIFE, so every helper is a plain sibling binding. `board-durations.js` loads after `board-cards.js` (`generate.go:44-53`), so `formatDurationMinutes` is not bound at definition time but is bound at render time.

**F6 — no test renders a real card DOM.** The Node behavior lane (`TestJavaScriptBehavior*`, roughly 40 tests) slices real functions out of the generated `index.html` via `sliceBalancedBlockAfter` (`generate_test.go:1362`) and runs them under `node` through `runJavaScriptBehaviorProbe` (`generate_test.go:275`). Two existing probes stub `makeRequestCard` rather than building one (`generate_test.go:1718`, `:5473`), so the done line's text is currently unasserted. Closest precedent for a span assertion: `TestJavaScriptBehaviorSpanFormattersCarryRoundedRemainders` (`generate_test.go:6651`).

**F7 — CSS classes to reuse.** `.req-card-completed` (`board.css:1150`) is the done line. `.elapsed-duration` (`board.css:1169`) already carries `tabular-nums` plus a separator `::before`. `.status-invalid-flag` (`board.css:934`) is the established inline flag inside a value span (`board-cards.js:65`). `.instant-with-relative` has no CSS rule at all.

**F8 — fixtures.** `writeVerifyFixture` (`verify_test.go:22`) writes literal frontmatter into a temp repo; `durationTicket` (`durations_test.go:16`) builds an in-memory ticket from a claimed/completed pair.

*Generated by Explore agent*

## Decisions

- [DECIDE & STATE] **D-01**: The span is `completed_at` minus `claimed_at`, read from frontmatter only — never from `completionTime`. Reasoning: F1. A git-dated completion instant measures when a commit landed, not how long the work took, and using it would make the card disagree with the Durations view for exactly the REQs that view excludes. Consequence, accepted: a done card whose completion instant came from `git log` shows no duration. The user's request scopes this to REQs that reached Done with a completed status, and `completed_at` is mandatory on every terminal flip (`actions/work.md` Step 8 substep 1), so the silent case is a bookkeeping defect the `anomaly` badge already reports.
- [DECIDE & STATE] **D-02**: The span's verdict (plain / over-ceiling / reversed) is computed in Go and shipped in the per-request payload; the client renders the verdict and never restates the four-hour ceiling. Reasoning: F3, and the board's own lesson that a rule's second reader must not become a second definition (REQ-219). This reuses the `excludedReason` convention rather than inventing one.
- [DECIDE & STATE] **D-03**: The reversed case reuses the existing completion-anomaly verdict rather than adding a second broken-stamp signal. Reasoning: F2. The done line renders an inline `.status-invalid-flag` reading `reversed stamps` in place of a number, so the line is never silently blank while the card's `anomaly` badge carries the full explanation.

- [DECIDE & STATE] **D-04**: The card draws its span with `formatElapsedDuration` (`board-core.js:128`), the card's own stopwatch vocabulary — `34m 00s`, `2h 40m`, `3d 04h` — not with the Durations view's `formatDurationMinutes`. Reasoning: raised by a PR reviewer against my own exploration note, and correct. `formatDurationMinutes` renders 34 minutes as `34.0 min`; a tenth of a minute is meaningful in the Durations table and is noise on a card footer. `formatElapsedDuration` is already what the state-timer line one row above draws in, so the two time lines on a card read alike, and it is a pure function of two instants, so passing `(claimed, completed)` needs no new formatter. Consequence: the `2h40m` / `34min` strings in `## Red-Green Proof` are illustrative of the reading, not of the spelling; the acceptance criteria below carry the drawn vocabulary. The formatter's clock-skew branch is unreachable here because the client branches on the Go verdict first (D-02/D-03).

<!-- D-XX counter: last used D-04. Next decision: D-05. -->

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/durations.go` — export the per-ticket span + verdict helper the aggregate already computes inline, so the payload and the Durations view read one definition of the ceiling
- `skills/do-work-board/tools/queue-kanban/generate.go` — add the span minutes and its verdict to the per-request client payload
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` — render the duration on the done line
- `skills/do-work-board/tools/queue-kanban/web/board.css` — style the duration reading and its markers
- `skills/do-work-board/tools/queue-kanban/generate_test.go` — payload-shape test plus the Node behavior probe asserting the rendered line
- `skills/do-work-board/tools/queue-kanban/durations_test.go` — unit tests for the exported span/verdict helper

**Acceptance criteria (restated from the REQ):**

1. A Recently-Done card whose status is `completed` or `completed-with-issues`, with both stamps parseable, renders its implementation span on the done line.
2. A span at or under four hours renders plainly, in the card's stopwatch vocabulary (`formatElapsedDuration`: `34m 00s`, `2h 40m`) — the same vocabulary the state-timer line uses, and no new formatter.
3. A span over four hours renders with a paused-session marker.
4. A reversed span renders a broken-stamp flag instead of a number.
5. A card with no parseable `claimed_at` renders the done line exactly as it does today.
6. Cancelled cards render no duration.
