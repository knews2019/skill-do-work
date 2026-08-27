---
id: REQ-374
title: 'Show how long each done card took'
status: completed
completed_at: 2026-08-26T14:42:00Z
commit: 5ad1d3d
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
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` and the crew rules; approach recorded as D-01..D-04 in `## Decisions` and the file list in `## Scope`. Ticked by the orchestrator — the dispatch instructed the builder not to write this file, so it could not tick these itself.
- [x] **[APPLY]:** Six files changed, all inside the `## Scope` list; `git diff --stat` shows no undeclared path.
- [x] **[UNIFY]:** Audited by the orchestrator, not self-reported. `git diff --stat` = 6 files, +528/−10, all declared. `gofmt -l .` prints nothing; `go vet ./...` clean; `go test -count=1 ./...` green. Debug-artifact grep over the added lines (`console.log|debugger|fmt.Print|TODO|FIXME`) returns 0 hits. Per file: `durations.go` — helper extracted, aggregate rewired to it, ceiling still has one definition; `generate.go` — three payload fields, gated on `isCompletedStatus`; `board-cards.js` — new node builder carries no `data-instant-ms`; `board.css` — comment-only; `durations_test.go` / `generate_test.go` — new assertions only, no existing test weakened.

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

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** Factored the claim-to-completion span and the read-time rule's verdict out of `buildDurationAggregate` into `measureImplementationSpan`, so the Durations view and the new card reading share one definition of the four-hour ceiling; shipped the measured minutes plus the verdict as three per-request payload fields gated on terminal success; and rendered them on the done line as a plain non-ticking `span.elapsed-duration` reading `took 2h 40m`, with an inline flag for the paused and reversed verdicts.

## Decisions (builder, renumbered from the report's D-04..D-10)

- [DECIDE & STATE] **D-05**: The terminal-success gate lives in `generate.go`, not inside `measureImplementationSpan`. The helper answers what a ticket's span is; the caller decides which tickets to ask. `buildDurationAggregate` already had its own `isCompletedStatus` check, so pushing the gate down would have made one of the two redundant and coupled the helper to a status vocabulary it does not otherwise mention.
- [DECIDE & STATE] **D-06**: Shipped a presence flag rather than inferring absence from a zero. `implementationSpanMinutes` carries `omitempty` and a genuine zero-minute span is possible, so without the flag the client could not tell "took no measurable time" from "unmeasured". Mirrors the existing `hasMedian` convention.
- [DECIDE & STATE] **D-07**: The client formats from the Go-measured minutes (`formatElapsedDuration(0, spanMs)`), not from `Date.parse(claimedAt)` / `Date.parse(completedAt)`. Go's `parseTimestamp` accepts `2006-01-02 15:04:05` and reads it as UTC; V8's `Date.parse` reads that same space-separated form as local time. Re-parsing client-side would silently shift the span on any such REQ and make the card disagree with the Durations view.
- [DECIDE & STATE] **D-08**: Added fixture `REQ-907` — a 34-minute span on a `completed-with-issues` REQ. It is the case that separates the card's vocabulary from the chart's (`34m 00s` vs `34.0 min`), so a revert to `formatDurationMinutes` fails loudly, and it covers `completed-with-issues`, which no other fixture did.
- [DECIDE & STATE] **D-09**: Both markers reuse `.status-invalid-flag` rather than earning a new class and colour tokens — it is the established inline flag inside a value span, and neither marker introduces a new semantic colour. A `" "` text node precedes each flag so the line reads `took 18h 00m likely paused` rather than running together for a screen reader.
- [DECIDE & STATE] **D-10**: No source-text lock-step guard between Go's `"paused"`/`"reversed"` strings and the client's branch, despite `TestFutureStampCauseClauseMatchesTheShippedClient` being the local precedent. The Node probe drives the real payload through the real `makeRequestCard`, so a rename on either side fails it end to end; a grep-style guard would be weaker and a second thing to maintain.
- [DECIDE & STATE] **D-11**: Dropped a planned "the client must not contain the four-hour number" assertion. `strings.Contains` for an absent token is the guard the prime warns about (REQ-245: it passes when the whole string is replaced) and it risked false positives on any `240`. The real protection is that the client has no arithmetic to derive a ceiling from.

<!-- D-XX counter: last used D-11. Next decision: D-12. -->

## Discovered Tasks

- `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails at HEAD on Chromium 141.0.7390.37 — its own vacuity guard fires ("the isolator was not exercised and the mutation pair is vacuous"). Reproduced on an unmodified `git archive HEAD` tree, so it predates this REQ.
- `_dev/tests/maintainer-verify.sh` cannot run in this environment: it exits 1 at `required command is unavailable: shellcheck` before any test runs.
- The done line's faint companions (`.relative-time`, `.elapsed-duration`) measure roughly 3.3:1 against `<body>` in both themes, under 4.5:1 for 11px text. Pre-existing for the whole line; the new reading inherits it rather than diverging mid-line.
- `tools/checks/preflight.sh` writes `do-work/working/baseline.json` and `baseline-failures.txt`, which nothing in this repo gitignores, so every work run leaves an untracked file behind. The suite's changelog says these are meant to be locally excluded; this repo has no such exclusion because the installer never runs against it.

## Testing

**Tests run:**
- `cd skills/do-work-board/tools/queue-kanban && go test -count=1 ./...` — `ok … 85.024s`, exit 0
- `gofmt -l .` — no output; `go vet ./...` — clean
- `bash _dev/tests/maintainer-verify.sh` (the canonical repository gate) — `Maintainer verification passed.`, exit 0
- Strict browser behavior lane with `QUEUE_KANBAN_BROWSER` set — one failure, established as not this REQ's (below)

**Result:** ✓ All passing. The canonical gate needed two tools installed in this container to run at all (ShellCheck 0.11.0 and `just`); neither is a repository change.

**Red-green validation:**
- `TestImplementationSpanVerdictBoundaryReadsTheOutlierCeiling`: ✗ `both stamps parse, yet the span measured nothing` → ✓
- `TestImplementationSpanMarksReversedStampsAndRefusesUnparseableOnes`: ✗ `a reversed pair of parseable stamps measured nothing` → ✓
- `TestImplementationSpanAgreesWithTheDurationsAggregate`: ✗ `REQ-411: card span = 0 min, Durations sample = 1080 min` → ✓
- `TestGeneratedRequestCarriesTheDoneCardImplementationSpan`: ✗ `REQ-901 hasImplementationSpan = false, want true` → ✓
- `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan`: ✗ `REQ-901 done line said "" about its span, want "took 2h 40m"` → ✓

These trace to the REQ's `## Red-Green Proof`: the probe drives the real payload through the real `makeRequestCard` and reads the done line's text, which is the GREEN the proof names. The captured proof's `2h40m` / `34min` spellings are illustrative (D-04); the drawn strings are `2h 40m` and `34m 00s`.

**Orchestrator-run mutations** — the builder supplied a table; these two I ran myself rather than accepting it:
- `analysisOutlierCeiling` comparison `>` → `>=`: `TestImplementationSpanVerdictBoundaryReadsTheOutlierCeiling` fails with `a span exactly at the ceiling read "paused", want the plain verdict`. The test reads the constant rather than restating four hours, so moving the ceiling cannot leave it passing (the REQ-322 lesson).
- Card formatter swapped to the Durations view's `formatDurationMinutes`: `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan` fails. Both mutations reverted; targeted re-run green.

**New tests added:**
- `durations_test.go`: `TestImplementationSpanVerdictBoundaryReadsTheOutlierCeiling`, `TestImplementationSpanMarksReversedStampsAndRefusesUnparseableOnes`, `TestImplementationSpanAgreesWithTheDurationsAggregate`
- `generate_test.go`: `TestGeneratedRequestCarriesTheDoneCardImplementationSpan`, `TestJavaScriptBehaviorDoneCardStatesItsImplementationSpan`, plus the shared six-REQ fixture builder

**Pre-existing failure, not this REQ's:** `TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails on Chromium 141.0.7390.37 with its own vacuity guard — *"the isolator was not exercised and the mutation pair is vacuous."* Verified by checking out HEAD into a separate worktree and running the probe there: byte-identical failure without this REQ's diff. Recorded as a discovered task.

**Rendered evidence:** the builder generated a board from the real archive and read the live DOM at `file:///…/scratchpad/board/render-check.html` on Chromium 141.0.7390.37 — 106 done cards, 94 carrying a span (`REQ-373 … · took 9m 30s`), the paused marker rendering on real archived cards (`REQ-345 → · took 11h 37m likely paused`), and `tickerTouchesSpan: 0` confirming no span node carries `data-instant-ms`.

## Review

**Overall: 77% at review time; remediated to Pass before archive** | 2026-08-26T14:21:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 82% |
| Test Adequacy | 72% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial → Pass after remediation |

**Important findings — all three fixed in this REQ rather than deferred, because each was a declared criterion this REQ had not actually met:**

- **I1 — a zero-minute span rendered `took NaNs`** (`impact-user-visible`). `generate.go`'s `omitempty` on `ImplementationSpanMinutes float64` dropped a genuine `0` from the JSON while `hasImplementationSpan: true` still shipped, so `board-cards.js` computed `Math.round(undefined * 60000)` = `NaN` and `formatElapsedDuration` fell through every tier to `"NaNs"`. Reachable two ways: identical stamps, or date-only stamps on both fields. Directly contrary to D-06, which added the presence flag *because* a zero span is possible and then left the renderer unable to consume the value the flag preserved. **Fixed:** dropped `omitempty` from the float, with the reason recorded at the field. **Pinned by:** fixture `REQ-908` (claimed and completed at the same instant), asserted in the marshalled JSON rather than on the struct — a struct-field assertion would pass with the tag restored — plus the render expectation `took 0s`. **Mutation:** restoring `omitempty` reproduces `done line said "took NaNs" about its span, want "took 0s"`.
- **I2 — `TestImplementationSpanAgreesWithTheDurationsAggregate` did not hold the property its comment claimed** (`impact-rule-change`). The review mutation-proved it: replacing the aggregate's verdict with a genuinely second definition using `6 * time.Hour` left the whole suite green, because all three fixture spans (40 min, 18 h, −3 h) sat far from `analysisOutlierCeiling` and any second ceiling in that gap agreed with the first everywhere the test looked. **Fixed:** added a straddling pair, `REQ-413` at exactly `analysisOutlierCeiling` and `REQ-414` one minute past it, both derived from the constant so moving the real ceiling moves the pair; plus a second vacuity guard asserting the pair still brackets it. **Mutation:** the review's own 6-hour second ceiling now fails with `REQ-414: card verdict = "paused", Durations verdict = ""`.
- **I3 — "a finished span never ticks" was pinned by nothing** (`impact-rule-change`). Adding `dataset.instantMs` to the span node passed the entire suite, though that property is the single reason `makeElapsedDurationNode` was not reused; had it regressed, every done card's span would have been rewritten one second after load as elapsed-since-epoch (roughly `20696d 07h`). Its only evidence was a one-off browser reading. **Fixed:** the probe now reports the span node's `dataset` keys and asserts the set is empty. **Mutation:** adding `dataset.instantMs` fails with `span node carries dataset keys [instantMs]`.

**Minor findings:** 3, report only — two Restatement Sweep hits routed to `do-work/prose-backlog.md` (`web/board-core.js:124-127`'s `formatElapsedDuration` contract still says "for a ticket sitting in a state" and justifies its second-resolution tier by "claim spans are short", now one of two callers; `skills/do-work-board/docs/board-guide.md:35` names the `anomaly` badge as the only reversed-span surface and documents none of the new vocabulary), plus the markers conveying their meaning only through `title`, which a screen reader does not reliably announce.

**Nits:** 4, report only — `durations.go`'s file header still says "Duration aggregation for the board's Durations view"; `dayMedianExclusionReason` names one consumer where there are now two; the span reaches the By-UR lens and Testing view too (correct, since the gate is in the payload, but untested there); the implementation commit bundled the REQ living log.

**Acceptance:** Pass after remediation. At review time: Partial, all six criteria demonstrated on the live archive (337/355 completed REQs spanned, 4 paused, 10 cancelled with none, 18 legitimately unmeasurable) and 13 of 14 adversarial mutations caught, with AC2 failing at its zero boundary. After remediation: the zero boundary renders `took 0s`, and the two unpinned properties are pinned — each fix mutation-verified by the orchestrator, not accepted on report.

**Verification after remediation:** `go test -count=1 ./...` → `ok … 83.745s`, exit 0. `bash _dev/tests/maintainer-verify.sh` run unpiped → `Maintainer verification passed.`, exit 0.

**Follow-ups created:** none. All three Important findings were remediated in this REQ.

**Deferred to `do-work/prose-backlog.md`:** the two Restatement Sweep hits.

*Reviewed by review-work action; remediated and re-verified by the orchestrator*

## Lessons Learned

**What worked:**
- Extracting the span and its verdict into one Go helper before adding the second reader. The card and the Durations chart cannot disagree about the four-hour ceiling because there is nothing for them to disagree with — the ceiling has one definition and both are readers of it. Shipping the applied verdict rather than the threshold is the same move the board already made for `excludedReason`.
- Checking what the card's completion instant actually *is* before measuring against it. `resolveCompletionTime` falls back to the commit's git date, so the obvious implementation — span from `claimed_at` to the instant already on the line — would have printed durations for exactly the REQs the Durations view excludes, and nobody would have noticed until the two numbers were compared.

**What didn't:**
- Pointing the builder at `formatDurationMinutes` because it was "the existing span formatter." It is the *chart's* formatter: `34.0 min`, where the card wants `34m 00s`. Two formatters existed and the vocabulary question — which surface is this? — was the one that decided it, not "is there already a function." A PR reviewer caught the mismatch against the REQ's own acceptance strings while the builder was still running.
- `omitempty` on a `float64` that a presence flag already guards. D-06 added `hasImplementationSpan` precisely because a zero-minute span is real, then left `omitempty` on the value — so the flag said "present" and the value vanished, and the client multiplied `undefined` into `took NaNs`. A presence flag and `omitempty` on the same field are contradictory by construction: one exists to preserve the zero the other deletes.
- Telling the builder "do not edit the REQ file" while the pipeline requires the builder to tick the P-A-U boxes in that file. Qualification failed on unticked boxes that could not have been ticked. The dispatch has to leave the builder the one write the pipeline demands of it, or the orchestrator has to own the audit — which is what happened, and is arguably better, but it should be chosen rather than discovered.

**Worth knowing:**
- **A test that spans a threshold widely does not test the threshold.** The agreement test used 40 min, 18 h and −3 h against a four-hour ceiling and looked thorough; a second ceiling of six hours passed it silently, because any threshold in the gap classifies all three samples identically. Only a pair straddling the real boundary — derived from the constant, so it moves with it — can catch a second definition. This is the sharper form of the prime's REQ-322 lesson: reading the constant is necessary but not sufficient, the fixture has to land where the constant decides something.
- **A property argued at length in a comment is not a property under test.** Two of this REQ's three review findings were exactly that shape: the code explained why a finished span must not tick and why a zero span must survive, and nothing asserted either. The comment is where a reader looks for intent; it is also where an unpinned invariant hides most comfortably.
- Go's `parseTimestamp` accepts `2006-01-02 15:04:05` and reads it as UTC; V8's `Date.parse` reads the same string as local time — a nine-hour divergence under `TZ=Asia/Tokyo`. Any board feature that re-parses a frontmatter stamp client-side inherits that bug. Format from Go-measured values instead.

## Orientation

Every finished card on the board now says how long it took, not just when it landed — `done Aug 25, 08:47 UTC · 1d ago · took 9m 30s` — with a `likely paused` marker where a span crossed the read-time ceiling and a `reversed stamps` flag where the bookkeeping is broken.

The measurement itself lives in the board tool's durations subsystem (`_dev/primes/prime-kanban-board.md`), which already owned the claim-to-completion span for the Durations view; this REQ gave that span a name and a second reader instead of a second definition. `[MAP CHANGED]` is not warranted — no new module, data flow, or contract; the per-request payload gained three additive fields and `.elapsed-duration` now dresses two kinds of node.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md`'s referenced paths all still resolve, and its write-surface count is untouched (this REQ adds no write surface). Two prose restatements elsewhere did go stale and are recorded in `do-work/prose-backlog.md` rather than fixed here.

**Discovered-task dispositions (Step 8):**
- The strict browser lane failure → REQ-375 (`pending-answers`, `impact-rule-change`). Its prose half — the stale "Chromium suppresses the boundary events" reason — was already on `do-work/prose-backlog.md` from REQ-336; this REQ is the failing test, which is not prose-only.
- The done line's contrast → REQ-376 (`pending-answers`, `impact-user-visible`).
- Preflight's unexcluded scratch files → REQ-377 (`pending-answers`, `impact-negligible`).
- The maintainer gate's missing tooling was resolved in-session by installing ShellCheck 0.11.0, `just`, and the Go 1.26.1 toolchain in this container. No repository change is implied, so no REQ was minted.

Fold-first scan before minting: `do-work/queue/` held no REQs at all and `grep -rl "^sweep: true" do-work/queue/` returned nothing, so no pending REQ, sweep or otherwise, shared a root cause with any of the three.
