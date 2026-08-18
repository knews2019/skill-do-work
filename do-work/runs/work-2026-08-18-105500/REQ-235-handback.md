# REQ-235 hand-back — Timeline period navigation and jump to now

**Branch:** `worktree-agent-REQ-235-timeline-period-navigation`
**Implementation commit:** `c42f823` (`[REQ-235] Give the Timeline day/week/month periods and a real jump to now`)
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-235-timeline-period-navigation`

## File manifest

| Action | File | What changed |
|---|---|---|
| Modified | `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` | Added `TIMELINE_DAY_MS`, `TIMELINE_NOW_JUMP_MARGIN_FRACTION`, `TIMELINE_PERIOD_LEVEL_NAMES`; six new pure functions (`timelinePeriodStart`, `timelineSteppedPeriodStart`, `timelinePeriodWindow`, `timelinePeriodLevelOfWindow`, `timelineNearestPeriodLevel`, `timelineFirstOpenRowIndex`, `timelineNowJump`); rewrote the `timeline-zoom-now` handler to move the window *and* the row list; wired the level buttons and prev/next; added `renderPeriodControls()` to `renderAll()`; renamed the local `wireZoomButton` → `wireToolbarButton` (it now wires period steps too). |
| Modified | `skills/do-work-board/tools/queue-kanban/web/template.html` | New `.timeline-periods` control group in the timeline toolbar (‹, Day, Week, Month, ›, plus a `#timeline-period-state` live-region span), placed **before** the existing zoom group and leaving it untouched; extended the panel `aria-label` and the `.timeline-hint` line. |
| Modified | `skills/do-work-board/tools/queue-kanban/web/board.css` | One rule, `.timeline-period-state`, inserted between `.timeline-toolbar` and `.timeline-legend`. |
| Modified | `skills/do-work-board/tools/queue-kanban/generate_test.go` | Added `regexp` import, helper `rendererDeclarationLine`, and `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow` at the end of the file. |

No file outside the write set was touched. `timeline.go` was not touched and needs no new payload field — `projection.queueEnd` was already shipped and is what the Now jump reads. `durations.go` / `board-durations.js` untouched.

## P-A-U

Recorded here because `## AI Execution State (P-A-U Loop)` lives in the REQ file, which is queue state I must not write. This is what I actually did, for the queue owner to tick from.

### [PLAN]

Read before writing code: `CLAUDE.md` (including § **Kanban Board Write Surfaces** — this REQ adds none, so that sentence needed no amendment); `_dev/primes/prime-kanban-board.md` in full including its `## Lessons` list; `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, `testing.md` (`tdd: true`). Then read the REQ body, and — because the brief made it structural — REQ-233's `timelinePannedWindow` and `timelineKeyboardWindow` in `web/board-timeline.js`, plus its two lock-in tests, to copy their shape before writing mine.

The approach I settled on before writing any code:

- **Period stepping is a candidate-plus-settle, not a new transform.** `timelinePeriodWindow` computes a calendar-aligned candidate with pure UTC arithmetic, then calls `timelineZoomedWindow(candidateStart, candidateEnd, 1, 0, boundStartMs, boundEndMs)`. Factor 1 at anchor fraction 0 reduces that function to "keep this window, apply the model's floor (`TIMELINE_MIN_SPAN_MS`), ceiling (the bound span) and edge clamp" — so the third mover inherits the same limits as the wheel, the zoom buttons and REQ-233's keys, and cannot acquire its own. This mirrors REQ-233's precedent of routing zoom through `timelineZoomedWindow` while pan got its own pure function.
- **The level is derived, not stored.** Rather than add a `periodLevel` field to `timelineViewState` — a second thing that can disagree with the window, i.e. exactly the divergence the REQ forbids — `timelinePeriodLevelOfWindow(start, end)` reads the level back off the window. That makes "a free zoom marks the level inexact" true by construction instead of by remembering to clear a flag, and it means the REQ adds **no state at all**.
- **The anchor is the window midpoint**, one rule for both level selection and stepping, because the midpoint of an exact period is always inside that period.
- **The Now jump returns two values** (`window` and `scrollTop`) from one module-level function, so the button is two assignments and the row-list half of the requirement is testable in Node instead of only in a browser.
- **Test plan:** one Node behaviour probe in the `TestJavaScriptBehavior*` family over a five-month fixture, using `timelineProbePreamble` so it drives the shipped constants; I anticipated needing a sibling helper for the non-numeric `TIMELINE_PERIOD_LEVEL_NAMES` declaration and wrote `rendererDeclarationLine` for it.

One plan item did not survive contact: I planned an explicit period-index clamp inside `timelinePeriodWindow` and it turned out to be dead code (see `## Decisions` #1). Removing it is the only material deviation from the plan above.

### [APPLY]

Code stayed inside the declared write set. `git diff --stat` below lists four files and they are the four declared ones. Specifically **not** touched: `timeline.go` (no payload change was needed — `projection.queueEnd` already ships and is what the Now jump reads), `generate.go`, `durations.go`, `board-durations.js`, `board-controls.js`, and everything under `do-work/` in either tree.

Three places pushed at the boundary, all resolved inside the write set:

1. **`setActiveButton` lives in `board-controls.js`**, which is outside my write set. I *called* it rather than editing it or writing a local copy — all fragments are concatenated into one closure and function declarations hoist, so it was already in scope. No edit to that file.
2. **`rendererNumericConstant` lives in `durations_test.go`**, outside my write set. Its sibling `rendererDeclarationLine` therefore went into `generate_test.go` (in the write set) rather than beside it.
3. **Renaming `wireZoomButton` → `wireToolbarButton`** is a change to pre-existing code that the REQ did not ask for. It is a function local to `renderTimelineView` with five call sites and zero references anywhere else in the repo (verified by `grep -rn wireZoomButton skills/ _dev/`), and I made its name wrong by reusing it for period steps — so it is cleaning my own mess, not improving adjacent code.

Main-tree writes: only the hand-back file. Playwright, whose cwd was the main tree, dropped `req-235-timeline-period-toolbar.png` in the main tree root during the browser drive; I read it for evidence and deleted it before the first hand-back, and confirmed `git status` in the main tree is clean apart from the hand-back and the orchestrator's `.req-reservations/` markers. All other scratch — the `qk-235` binary, the generated board, verify logs — went to `/tmp`.

Sibling-collision handling on the three shared files, per the brief: my `template.html` edit is one new `<div>` inside `#view-timeline`'s toolbar plus two text-only edits on that section's own `aria-label` and `.timeline-hint`; my `board.css` edit is one new rule inserted between `.timeline-toolbar` and `.timeline-legend`; my `generate_test.go` edit is a contiguous append at the very end plus one added import line. Nothing above my additions was reflowed or reordered — confirmed by the hunk headers below.

### [UNIFY]

```
$ git diff --stat HEAD~1 HEAD
 .../tools/queue-kanban/generate_test.go            | 276 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-timeline.js       | 214 +++++++++++++++-
 .../do-work-board/tools/queue-kanban/web/board.css |   9 +
 .../tools/queue-kanban/web/template.html           |  16 +-
 4 files changed, 502 insertions(+), 13 deletions(-)
```

Reviewed every changed file:

- **`generate_test.go`** — exactly two hunks, confirmed by `git diff -U0 ... | grep '^@@'`:
  ```
  @@ -12,0 +13 @@ import (
  @@ -2507,0 +2509,275 @@ func TestTimelinePanelStatesItsKeyboardInteraction(t *testing.T) {
  ```
  One added line (`"regexp"`, in alphabetical position in the existing import block — required by `rendererDeclarationLine`, and the only edit above my block), and one contiguous 275-line append starting at the old end of file (2507 lines → 2783). Checked: the new probe reads its constants through `timelineProbePreamble` / `rendererDeclarationLine` rather than hard-coding them, every assertion message names the failure it pins, and no assertion is a truthiness catch-all.
- **`web/board-timeline.js`** — checked that the five new pure functions are module-level (so `sliceBalancedBlockAfter` can reach them) and brace-balanced with no braces inside comments or string literals, which that helper cannot parse around; that every new listener path goes through `renderAll()` so `renderPeriodControls()` cannot be skipped; that no new listener was added to a node outliving a render outside `addTimelineListener` (the period buttons use `.onclick` assignment, matching the existing toolbar helper, which overwrites rather than stacks); and that the rename hit all five call sites (`grep -rn wireZoomButton` now returns nothing).
- **`web/template.html`** — checked the new group is inside the timeline toolbar and before the untouched zoom group; that the panel `aria-label` still contains the three phrases `TestTimelinePanelStatesItsKeyboardInteraction` requires (`Timeline`, `arrow keys`, `plus and minus`); and that `#timeline-scroll` still carries `tabindex="0"`.
- **`web/board.css`** — one rule, using existing tokens (`--ink-faint`); no change to `.control-button`, `.control-group`, or any shared selector.

Linters and checks run, all from the worktree root or the tool directory:

| Check | Result |
|---|---|
| `gofmt -l .` | no output (exit 0) — nothing unformatted |
| `go vet ./...` | exit 0 |
| `go test ./...` (queue-kanban) | `ok ... 12.718s` |
| `bash _dev/tests/maintainer-verify.sh` | **exit 0**, re-run after the rename; includes ShellCheck warning-level lint, the aggregate contract suite, shipped-package reference contract, `queue-kanban go vet`, uncached ordinary tests, the strict JavaScript behavior lane (`TestMaintainerStrictJavaScriptBehaviorLane` PASS), and the audit-metrics vet + tests |
| Browser console during the live drive | 0 errors after the post-rename reload |

Debug artifacts: none. `git diff HEAD~1 HEAD | grep -nE '^\+.*(console\.(log|debug|warn)|debugger|TODO|FIXME|XXX|fmt\.Print|t\.Log\()'` returns no matches across all four files. The RED-2 defect (the Sunday-origin week) was reverted from a saved copy and the GREEN run re-confirmed after the revert; the working tree at commit time contained no experimental code.

## RED / GREEN, verbatim

### RED 1 — the functions do not exist (proves nothing on its own)

```
=== RUN   TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow
    generate_test.go:2542: web/board-timeline.js declares no numeric constant TIMELINE_DAY_MS
--- FAIL: TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow (0.23s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.542s
FAIL
EXIT=1
```

### RED 2 — code present, one behaviour wrong

The defect: `timelinePeriodStart` computing the week origin as `dayStartMs - instant.getUTCDay() * TIMELINE_DAY_MS` — the classic Sunday-origin week — instead of `- ((instant.getUTCDay() + 6) % 7) * TIMELINE_DAY_MS`. The assertion that exists for exactly that fails:

```
=== RUN   TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow
    generate_test.go:2698: the week window starts at 2026-06-14T00:00:00.000Z (weekday 0); want a Monday at 00:00 UTC
--- FAIL: TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow (0.29s)
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.603s
FAIL
EXIT=1
```

> Note on why the RED-2 defect is the week origin and not the clamp the brief suggested: I wrote `timelinePeriodWindow` with an explicit period-index clamp, and the test **passed without it**. Because the anchor is the window's own midpoint and `timelineZoomedWindow` already clamps to the bounds, the "stop at the last period" behaviour is *emergent* from the one window model — a second clamp would have been dead code. I deleted it and used a defect that the calendar assertions actually catch. See `## Decisions`.

### GREEN

```
=== RUN   TestJavaScriptBehaviorAssembledClientSyntax
--- PASS: TestJavaScriptBehaviorAssembledClientSyntax (0.04s)
=== RUN   TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer
--- PASS: TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer (0.25s)
=== RUN   TestTimelinePanelStatesItsKeyboardInteraction
--- PASS: TestTimelinePanelStatesItsKeyboardInteraction (0.20s)
=== RUN   TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow
--- PASS: TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow (0.25s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.930s
EXIT=0
```

**REQ-233's two tests both still pass** — see the middle two lines above.

## maintainer-verify

Run from the worktree root after every edit including the rename:

```
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (3.69s)
...
Maintainer verification passed.
EXIT=0
```

**Exit code: 0.** Not piped through `tail`/`head` — output was redirected to a file and the exit code echoed on its own line.

## What I rendered and drove

Built `/tmp/qk-235` from the worktree and generated the live repo as the fixture:
`queue-kanban generate --repo-root <worktree> --out /tmp/board-235` → **234 REQs, 19 still open**, range **28 May → ~20 Aug 2026** (three months). Served on `127.0.0.1:8235` and driven in Chromium at 1440×960 through Playwright. Console: no errors (only a `favicon.ico` 404 from the static server on the first load; zero errors on the post-rename reload).

### Axis label strings, before and after each control

Measured by reading `#timeline-axis text` out of the live DOM.

| Step | Axis labels | `#timeline-period-state` | pressed |
|---|---|---|---|
| On arrival (Fit all) | `28 May · 11 Jun · 25 Jun · 9 Jul · 23 Jul · 6 Aug · 20 Aug` | `custom span` | none |
| click **Week** | `6 Jul · 7 Jul · 8 Jul · 9 Jul · 10 Jul · 11 Jul · 13 Jul` | `one week` | `week` |
| click **›** | `13 Jul · 14 Jul · 15 Jul · 16 Jul · 17 Jul · 18 Jul · 20 Jul` | `one week` | `week` |
| click **‹** | `6 Jul · …· 13 Jul` (exact round-trip) | `one week` | `week` |
| click **‹** again | `29 Jun · 30 Jun · 1 Jul · 2 Jul · 3 Jul · 4 Jul · 6 Jul` | `one week` | `week` |
| click **Month** | `1 Jul · 6 Jul · 11 Jul · 16 Jul · 21 Jul · 26 Jul · 1 Aug` | `one month` | `month` |
| click **›** (last month of the range) | `20 Jul · 25 Jul · 30 Jul · 4 Aug · 9 Aug · 14 Aug · 20 Aug` | `custom span` | none |
| click **Day** | `4 Aug 00:00 · 04:00 · 08:00 · 12:00 · 16:00 · 20:00 · 5 Aug 00:00` | `one day` | `day` |
| click **›** | `5 Aug 00:00 · … · 6 Aug 00:00` | `one day` | `day` |
| **Week, then › ×40** | `13 Aug · … · 20 Aug` (stops at the range end) | `custom span` | none |
| **one more ›** | identical — clamped, does not move | `custom span` | none |
| Week, then **+ (zoom in)** | `11 Aug · 12 Aug · 12 Aug · 13 Aug · 14 Aug · 14 Aug · 15 Aug` | `custom span` | none |
| **Fit all** | `28 May · … · 20 Aug` | `custom span` | none |

`6 Jul`, `13 Jul` and `29 Jun` 2026 are all Mondays (confirmed in-page via `Date.prototype.toUTCString`). The month windows start on the 1st and end on the 1st. The day windows run `00:00` → `00:00`.

### `Now`, driven as a real click

Starting from `Fit all` with `scrollTop` forced to `0` (top of the archive, first visible row `REQ-001`), one real click on `#timeline-zoom-now`:

- `scrollTop`: **0 → 720**, which is exactly `firstOpenRowIndex (40) × TIMELINE_ROW_HEIGHT (18)` — the first still-open REQ, `REQ-041`. Verified against the payload, not against a constant.
- `REQ-041`'s row measured with `getBoundingClientRect()`: top **265** against a scroll-host top of **264** — it lands at the top of the viewport.
- `REQ-041`'s open wait bar measures **left 156 → right 551**, i.e. from the plot's left edge to the now-line; the now-line itself measures **x = 551** inside a plot running **156 → 1361**. Both the line and the open bar are on screen, measured from the live DOM.
- Window covers `now` (`2026-08-18T11:25:03Z`) and `projection.queueEnd` (`2026-08-18T11:55:31Z`); axis reads `18 Aug 10:00 … 18 Aug 12:00`.

### REQ-233's keyboard path, driven with real keypresses

Tabbed to `#timeline-scroll` with two real `Tab` presses (`document.activeElement.id === "timeline-scroll"`) — not a scripted `.focus()`. A real `ArrowRight` moved the window forward and the period readout followed it to `custom span`. Zoom `+`/`−` and `Fit all` unchanged.

## How period stepping goes through the one window model

`timelinePeriodWindow(anchorMs, levelName, stepCount, boundStartMs, boundEndMs)` does two things:

1. Computes a **candidate** window by pure UTC calendar arithmetic — `timelinePeriodStart` floors the anchor to the level's boundary, `timelineSteppedPeriodStart` moves it `stepCount` periods and produces the end.
2. Hands that candidate to **`timelineZoomedWindow(candidateStart, candidateEnd, 1, 0, boundStartMs, boundEndMs)`**. Factor 1 at anchor 0 is "keep this window, apply the model's floor, ceiling and edge clamp". It is the same function the wheel, the zoom buttons and REQ-233's `timelineKeyboardWindow` call.

So there is exactly one floor (`TIMELINE_MIN_SPAN_MS`), one ceiling (the bound span) and one edge clamp, and the period path cannot acquire its own. The clamp at the ends of the range is a *consequence* of that single call rather than a second rule: the anchor is the current window's midpoint, which the bounds clamp already keeps inside the range, so `›` held down converges and stops (proved above: 40 clicks then one more, identical window).

**No new state.** There is no stored "current level". `timelinePeriodLevelOfWindow(windowStart, windowEnd)` reads the level back off `timelineViewState` — a window is a level exactly when it starts on that level's boundary and ends on the next one. `timelineViewState` remains the single window model.

## How a free zoom marks the level inexact

Because the level is derived, not stored, it falls out for free: `renderPeriodControls()` is called from `renderAll()`, which every mover already calls (buttons, keydown, wheel, pointermove drag, resize). It sets `aria-pressed` via the existing `setActiveButton` and writes `one day` / `one week` / `one month` or **`custom span`** into `#timeline-period-state` (a `role="status" aria-live="polite"` span). After a `+`, a drag, or `Fit all`, no level button is highlighted and the readout says `custom span` — proven in the table above.

The same honesty applies at the ends of the range: a period the bounds cut short (e.g. the last month, clamped to `20 Jul → 20 Aug`) is genuinely not a whole month, so it reads `custom span` rather than claiming `one month`.

## Integration seams

**None.** Everything landed inside the four write-set files. No payload field, no `timeline.go` change, no new write surface — `CLAUDE.md` § Kanban Board Write Surfaces still counts three and needs no amendment. `setActiveButton` is reused from `board-controls.js` without modification (all fragments share one closure and function declarations hoist).

One coordination note for the integrator rather than a seam: **`web/template.html`, `web/board.css` and `generate_test.go` are shared with the concurrent URs-only-lens builder.** My template edit is a new `<div class="control-group timeline-periods">` inside `#view-timeline`'s toolbar plus two text-only edits on that section's `aria-label` and `.timeline-hint`; my CSS edit is one new rule between `.timeline-toolbar` and `.timeline-legend`; my Go edits are the `regexp` import line and an append at the end of the file. Nothing was restructured.

## Discovered Tasks

- **The axis tick label lies about minutes at sub-hour tick spacing.** `timelineFormatAxisTick` prints `"18 Aug 10:00"` for an instant that is actually `10:55` — the hour is real but the `":00"` is a literal. It is only visible once the window is small enough that six ticks fall inside a couple of hours, which the new Day level and the tightened `Now` window both make routine. Pre-existing REQ-227 behaviour; not fixed here (out of scope, and it touches the axis renderer the REQ told me to leave alone).
- **A week window's six evenly spaced ticks do not land on day boundaries.** Seven days over six intervals gives 1.167-day ticks, so the interior labels read `7 Jul · 8 Jul · 9 Jul · 10 Jul · 11 Jul` and skip 12 Jul. The two ends are exact, so calendar alignment is still visible, but a period-aware tick generator (7 ticks for a week, one per day) would read much better. Deliberately not done — it is a redesign of the axis, and the REQ said keep it small.
- **`Fit all` is what makes the board unreadable in the first place.** With the projection extending the range, every measured bar on the live board sits in the right-hand fraction of the plot. The period controls make that navigable, but a "fit to the last N days" default would address the original complaint more directly. Capture-worthy, not fixable inside this write set.

## Decisions

1. **The period-index clamp was deleted rather than written.** I implemented an explicit "clamp the stepped period to the first/last period that overlaps the range" guard and the test passed without it: anchoring on the window midpoint plus `timelineZoomedWindow`'s existing bounds clamp already makes stepping converge and stop. Per CLAUDE.md's *delete before you add* and coding-guardrails § Earned Defense, an unearned second clamp is exactly the "second definition of where the window goes" this REQ exists to prevent. Cost: RED 2 had to be a different defect (the Monday/Sunday week origin), which is the failure the calendar assertions were written for anyway.

2. **The level is derived, never stored.** The alternative was a `timelineViewState.periodLevel` field. Storing it would have been a second thing that can disagree with the window — precisely the divergence the REQ forbids — and would have needed explicit invalidation on every zoom, drag and keypress. Deriving it makes "a free zoom marks the level inexact" true by construction instead of by remembering to clear a flag.

3. **The anchor is the window's midpoint, for both level selection and stepping.** One rule rather than two. The midpoint of an exact period is always inside that period, so `›` from an exact week is always the next week, and choosing a level from a wide view keeps the reader near what they were looking at rather than teleporting to the start of the range.
   *Known wrinkle:* at the extreme end of the range the window is clamped and its midpoint can fall in the previous period, so a `‹` immediately after hitting the end can skip one period. Fixing it needs remembered state, which decision 2 rules out. It affects one press at one edge; I judged the trade the right way round and am flagging it rather than hiding it.

4. **At the ends of the range the shared clamp wins over grid alignment.** A calendar period that overhangs the padded bounds is shifted inward by `timelineZoomedWindow` and therefore stops being grid-aligned. The alternative — widening the bounds to the period grid for the period path only — would have given the period controls a different range from the pan and drag paths, and a subsequent `←` would then visibly snap the window. Made honest by decision 2: such a window reads `custom span`.

5. **`Now` covers `[min(now, queueEnd), max(now, queueEnd)]` plus a 10% margin** (`TIMELINE_NOW_JUMP_MARGIN_FRACTION`, floored at half the minimum span) rather than preserving the current span. The REQ requires both the now-line and the forecast to be inside the window; preserving the span cannot guarantee that when the reader is zoomed in.

6. **`Now` sets `scrollTop` to exactly `firstOpenRowIndex × TIMELINE_ROW_HEIGHT`**, putting the first still-open row at the top of the viewport, and leaves the scroll alone entirely (`scrollTop: null`) when nothing is open. Context rows above it would have been prettier and less exactly assertable; landing the row at a known position is the behaviour the test can pin.

7. **`wireZoomButton` renamed to `wireToolbarButton`.** It now wires the period steps as well, so the old name had become wrong — cleaning up my own mess rather than improving adjacent code. It is a local function inside `renderTimelineView` with five call sites and no references anywhere else in the repo (checked).

8. **Reused `setActiveButton` from `board-controls.js`** instead of writing a local pressed-state toggler, so the period buttons get the same `is-active` + `aria-pressed` treatment as every other pill group on the board.

9. **Added `rendererDeclarationLine` to `generate_test.go`** as `rendererNumericConstant`'s sibling, so the probe drives the shipped `TIMELINE_PERIOD_LEVEL_NAMES` list rather than a copy in the test that would go stale beside it — same reason `timelineProbePreamble` exists.
