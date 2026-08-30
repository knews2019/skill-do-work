# REQ-390 Hand-back — Replace the Timeline's Day/Week/Month Periods with Trailing Windows

## Branch

`worktree-agent-REQ-390-timeline-trailing-window-periods`

| Commit | Subject |
|---|---|
| `ebd62c280bf070b1827b3b8b801cf069bc43f46b` | REQ-390: add the failing trailing-window tests (RED) |
| `10472cbc23c755adf1530df126fcc0d211a3128e` | REQ-390: replace the timeline's Day/Week/Month with trailing windows |

Base: `18c19f1`. Working tree clean. Not pushed, not merged.

## Implementation Summary

| File | Action | What actually landed |
|---|---|---|
| `skills/do-work-board/tools/queue-kanban/web/template.html` | modify | The `.timeline-periods` group's three `data-timeline-period="day\|week\|month"` buttons became five `="1\|7\|30\|90\|all"` buttons labelled Last day / Last 7 days / Last 30 days / Last 90 days / All days. Arrow `title`/`aria-label` became "Step back one window" / "Step forward one window". The toolbar comment, the `#view-timeline` panel `aria-label` and the `.timeline-hint` period sentence were rewritten for trailing windows. `.timeline-periods`, `#timeline-period-state`, `#timeline-period-prev/next` and `aria-label="Timeline period"` are unchanged. |
| `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` | modify | Added `timelineTrailingWindow(trailingWindowValue, nowMs, boundStartMs, boundEndMs)` and `timelineUtcDayStart(epochMs)`. Renamed `applyPeriodWindow`→`applyTrailingWindow`, `applyPeriodStep`→`applyTrailingWindowStep`, `renderPeriodControls`→`renderTrailingWindowControls` (now derived from the DOM control set). `steppedWindowFor` calls `timelinePannedWindow` directly. Deleted `TIMELINE_PERIOD_LEVEL_NAMES`, `timelinePeriodStart`, `timelineSteppedPeriodStart`, `timelinePeriodWindow`, `timelinePeriodAnchor`, `timelinePeriodLevelOfWindow`, `timelinePeriodGridOfWindow`, `timelineSteppedWindow` and their comment blocks. Re-pointed six prose sites that named the deleted machinery. |
| `skills/do-work-board/tools/queue-kanban/web/board.css` | modify — **outside the declared Scope, see D-03** | Added `.timeline-periods { flex-wrap: wrap; }` with the measurement that earned it, and rewrote the now-inaccurate `.timeline-period-state` comment. |
| `skills/do-work-board/tools/queue-kanban/generate_test.go` | modify | Added `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow`. Replaced `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow` with `TestJavaScriptBehaviorTimelineNowJumpLandsOnTheOpenWork` (its Now-jump half kept, calendar half dropped). Deleted `TestJavaScriptBehaviorTimelinePeriodChipsLandOnNowAndStepsStayAligned`. Re-pointed the typed-dates and axis-labels tests off the deleted helpers. Fixed two stale comments naming `timelinePeriodWindow`. |
| `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` | modify | Added `TestBrowserBehaviorTimelineTrailingWindowsEndAtNow`. Migrated three existing probes off `day`/`week`/`month`, including rebuilding the landing probe's clause-(6) refusal case. |

Net: 714 insertions, 1018 deletions against `18c19f1`.

`timelineTrailingWindow` in full:

```js
function timelineTrailingWindow(trailingWindowValue, nowMs, boundStartMs, boundEndMs) {
  var trailingDayCount = Number(trailingWindowValue);
  var candidateStartMs = boundStartMs;
  var candidateEndMs = boundEndMs;
  if (isFinite(trailingDayCount) && trailingDayCount > 0) {
    candidateStartMs = nowMs - trailingDayCount * TIMELINE_DAY_MS;
    candidateEndMs = nowMs;
  }
  candidateStartMs = Math.min(Math.max(candidateStartMs, boundStartMs), boundEndMs);
  candidateEndMs = Math.max(Math.min(candidateEndMs, boundEndMs), boundStartMs);
  return timelineZoomedWindow(candidateStartMs, candidateEndMs, 1, 0, boundStartMs, boundEndMs);
}
```

## P-A-U

### [PLAN]

Read `CLAUDE.md`, `_dev/primes/prime-kanban-board.md` in full, `_dev/primes/lessons-kanban-board.md`, and the `general.md` / `coding-guardrails.md` / `communication-style.md` / `testing.md` crew members before any code. Approach, as committed:

1. One new pure function, `timelineTrailingWindow`, clamping each endpoint into the bounds before `timelineZoomedWindow(…, 1, 0, …)` — the pattern `timelineTypedWindow` and the deleted `timelinePeriodWindow` both use, so the chips inherit the shared floor/ceiling/edge rules and acquire none of their own.
2. Keep a day-only `timelineUtcDayStart` so `timelineTypedWindow`'s surviving clamp does not break when `timelinePeriodStart` goes — the exploration's BLOCKER, confirmed real at `board-timeline.js:565-566`.
3. `‹`/`›` collapse to `timelinePannedWindow` (D1); the whole grid-vs-exactness split is deleted, not ported.
4. Lit chip and state text derived by re-asking each `#view-timeline [data-timeline-period]` in the DOM (D4).
5. Split the Now-jump coverage out of the doomed calendar test rather than losing it — the exploration was right that it is the only Node-lane coverage of `timelineNowJump`'s NaN/degenerate cases and `timelineFirstOpenRowIndex` ordering.

No lesson in `lessons-kanban-board.md` contradicted the approach. REQ-235 (derive, don't store), REQ-237 (say loudly which tests a rule-removal puts in scope), REQ-305 (drive the shipped function, not a copy), REQ-322 (read the constant, don't restate it) and REQ-345 (this probe going vacuous) all shaped what landed and are cited at the sites that honour them.

### [APPLY]

Tests first (`ebd62c2`), then implementation and migration (`10472cb`). Everything inside the Scope declaration except `web/board.css`, which the render forced — see D-03.

### [UNIFY]

`git diff --stat 18c19f1`:

```
 .../tools/queue-kanban/generate_test.go            | 896 ++++++---------------
 .../queue-kanban/timeline_browser_probe_test.go    | 404 +++++++---
 .../tools/queue-kanban/web/board-timeline.js       | 387 +++------
 .../do-work-board/tools/queue-kanban/web/board.css |  17 +-
 .../tools/queue-kanban/web/template.html           |  28 +-
 5 files changed, 714 insertions(+), 1018 deletions(-)
```

File-by-file review, and what was checked in each:

| File | Checked |
|---|---|
| `web/template.html` | Read the whole diff. Five buttons in DOM order with the REQ's exact labels; no `day`/`week`/`month` value remains; `aria-pressed="false"` and `class="control-button"` on every chip; the group's own id/class/`aria-label` untouched; three prose surfaces (panel `aria-label`, toolbar comment, `.timeline-hint`) rewritten and no longer say "calendar". |
| `web/board-timeline.js` | Read the whole diff. `node --check` on the file passes. Every deleted symbol has zero remaining references anywhere in the module (grep for all eleven names returns nothing). The pre-clamp is present and carries its earned reason. `stepLandsOffTheData`, `renderControlAvailability`, `timelineOwnedControls`, `retireTimelineControls` untouched as planned. Six prose sites that named the deleted machinery re-pointed; two reflowed comments re-checked for line shape. |
| `web/board.css` | Read the whole diff (17 lines). One new rule, three declarations of context in its comment including the measured number and the browser build; the `.timeline-period-state` comment now describes what the span actually says. No other selector touched. |
| `generate_test.go` | Read every changed hunk. `gofmt -l` clean, `go vet` clean. Confirmed the Now-jump half kept every assertion it had (both degenerate cases, both `holdsNow` checks, the scroll assertion and its anti-vacuity guard). Confirmed the typed-dates test now slices `timelineUtcDayStart` and still reads `lastDayStartMs` off the shipped helper rather than restating it. Confirmed the axis oracle reads both `TIMELINE_WEEK_ALIGNMENT_MS` and `TIMELINE_DAY_MS` off the renderer. |
| `timeline_browser_probe_test.go` | Read every changed hunk. `gofmt -l` clean, `go vet` clean. Confirmed every snapshot in the new probe is `href`-checked, that the new probe's three premises fatal rather than pass vacuously, and that the rebuilt clause (6) carries a premise the live queue cannot silently satisfy. |

Debug artifacts: `git diff 18c19f1 -U0 | grep` for `console.log/debug/warn`, `debugger`, `TODO`, `FIXME`, `XXX`, `HACK`, `.only(`, `.skip(`, `t.Skip(`, `fmt.Print` over all added lines — **no matches**.

Serial-only integrator files: `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` — untouched, per D9. Nothing under `do-work/` written in either tree except this hand-back file.

## Testing

Module directory unless stated: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-390-timeline-trailing-window-periods/skills/do-work-board/tools/queue-kanban`.

Browser lane used `QUEUE_KANBAN_BROWSER=/tmp/claude-0/-home-user-skill-do-work/cee71dc2-6250-5a0b-9c51-9822eef12052/scratchpad/chrome151/chrome-linux64/chrome`, which self-reports **Google Chrome for Testing 151.0.7922.174**, per the brief's correction 1. The plan's `/opt/pw-browsers/chromium-1194` (Chromium 141) was not used anywhere.

### RED, before any change to `template.html` or `board-timeline.js`

```
$ go test -count=1 -run '^TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow$' -v .
=== RUN   TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow
    generate_test.go:6976: anchor "function timelineTrailingWindow(" not found in the generated page
--- FAIL: TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow (2.07s)
FAIL
exit status 1
```

```
$ QUEUE_KANBAN_BROWSER=…/chrome151/chrome-linux64/chrome \
    go test -count=1 -run '^TestBrowserBehaviorTimelineTrailingWindowsEndAtNow$' -v .
=== RUN   TestBrowserBehaviorTimelineTrailingWindowsEndAtNow
    timeline_browser_probe_test.go:1899: the period toolbar offers 3 controls
      ([{Value:day Label:Day Pressed:false} {Value:week Label:Week Pressed:false}
        {Value:month Label:Month Pressed:false}]), want the five trailing windows
      [{Value:1 Label:Last day Pressed:false} {Value:7 Label:Last 7 days Pressed:false}
       {Value:30 Label:Last 30 days Pressed:false} {Value:90 Label:Last 90 days Pressed:false}
       {Value:all Label:All days Pressed:false}]
--- FAIL: TestBrowserBehaviorTimelineTrailingWindowsEndAtNow (3.47s)
FAIL
exit status 1
```

Both are the REQ's `## Red-Green Proof` RED case. The browser probe fails at the assertion (the shipped control set is still day/week/month); the Node probe fails in `sliceBalancedBlockAfter` because the function under test does not exist yet — for a Node-lane probe of an absent shipped function that is the only available RED, and the browser half supplies the assertion-level one.

These two commands and their outputs are the RED half of commit `ebd62c2`, which contains only the tests.

### GREEN, in the plan's order

| # | Command | Exit | Result |
|---|---|---|---|
| 1a | `go test -count=1 -run '^TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow$' -v .` | 0 | `--- PASS (2.18s)` |
| 1b | `QUEUE_KANBAN_BROWSER=… go test -count=1 -run '^TestBrowserBehaviorTimelineTrailingWindowsEndAtNow$' -v .` | 0 | `--- PASS (3.50s)` |
| 2 | `go test -count=1 -run '^TestJavaScriptBehaviorTimeline' -v .` | 0 | 19 tests, all PASS, 32.65s — includes the two re-pointed oracles and the split-out `TestJavaScriptBehaviorTimelineNowJumpLandsOnTheOpenWork` |
| 3 | `QUEUE_KANBAN_BROWSER=… go test -count=1 -run '^TestBrowserBehaviorTimeline' -v .` | 0 | 14 tests, all PASS, 55.55s — includes all three migrated probes |
| 4a | `go vet ./...` | 0 | clean |
| 4b | `go test -count=1 ./...` | 0 | `ok … 277.69s` |
| 5 | `go test -count=1 -run '^TestMaintainerStrictJavaScriptBehaviorLane$' -v .` | 0 | `--- PASS (80.70s)` |
| 6 | `bash _dev/tests/maintainer-verify.sh` (repo root, not piped, `QUEUE_KANBAN_BROWSER` **unset**) | **0** | `Maintainer verification passed.` |

Step 6 was run twice — once before and once after the final comment reflows — exit 0 both times. It prints `SKIP: no browser is available; strict browser behavior lane was not run.` because no chrome name is on this machine's PATH. `QUEUE_KANBAN_BROWSER` was deliberately not exported into it (see D-05).

Named RED→GREEN transitions:

| Test | Before | After |
|---|---|---|
| `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow` | FAIL — `anchor "function timelineTrailingWindow(" not found in the generated page` | PASS |
| `TestBrowserBehaviorTimelineTrailingWindowsEndAtNow` | FAIL — `the period toolbar offers 3 controls ([day Day] [week Week] [month Month]), want the five trailing windows` | PASS |

`TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` — the pre-existing failure the exploration recorded on Chromium 141 — **passes** on Chrome 151, as the brief said it would. `TestBrowserBehaviorCompletionCompanionsKeepReadableContrast` was not run (brief correction 2; the whole strict browser lane was never invoked).

### Render evidence

Live board built from this worktree with `queue-kanban generate`, opened in **Google Chrome for Testing 151.0.7922.174**, `location.href` returned in the same `evaluate` as every measurement below.

Toolbar geometry, `.timeline-periods` against `.timeline-toolbar`, HEAD (`18c19f1`) versus this branch, same browser build, same board:

| Window | HEAD group width / right / fits | REQ-390 group width / right / fits | chip rows |
|---|---|---|---|
| 1400×900 | 366.78 / 418.78 / **true** | 685.22 / 737.22 / **true** | 1 |
| 900×900 | 366.78 / 418.78 / **true** | 685.22 / 737.22 / **true** | 1 |
| 500×900 *(before the CSS rule)* | 366.78 / 406.78 / **true** | 465.63 / 505.63 / **false** — toolbar right edge is 469, and the chip labels wrapped to 2–3 text lines (heights 51 and 70.5 against 31.5) | 2 |
| 500×900 *(after `.timeline-periods { flex-wrap: wrap }`)* | — | 405 / 445 / **true**, all chips back to 31.5px | 2 |

500px is Chrome headless's minimum window width on this build; it stands in for the plan's narrow case and for `frontend.md`'s 320px, which the browser will not honour. The page never scrolled sideways at any width, before or after (`documentElement.scrollWidth === clientWidth` at 500 and 900, both HEAD and this branch) — the pre-fix symptom was the `›` arrow and the state span clipped at the panel edge, not a page-wide overflow.

Chip behaviour, this repo's live board (`now = 2026-08-29T22:29:05Z`), measured at 1400px, 900px and 500px — identical at all three:

```
press 1   -> lit ['1']   | state 'last day'     | 2026-08-28 22:29 UTC → 2026-08-29 22:29 UTC
press 7   -> lit ['7']   | state 'last 7 days'  | 2026-08-22 22:29 UTC → 2026-08-29 22:29 UTC
press 30  -> lit ['30']  | state 'last 30 days' | 2026-07-30 22:29 UTC → 2026-08-29 22:29 UTC
press 90  -> lit ['90']  | state 'last 90 days' | 2026-05-31 22:29 UTC → 2026-08-29 22:29 UTC
press all -> lit ['all'] | state 'all days'     | 2026-05-27 20:14 UTC → 2026-09-01 05:28 UTC
```

Every chip lights alone, names itself in the state span, and ends the window at 22:29 — the board's own now, minute-truncated by the readout. All days spans the full recorded range including forecast and padding, as D5 requires.

The `part of` branch cannot be reached on this repo's board (its archive is older than 90 days), so it was exercised on a **synthetic young board** — the same generated page with `timeline.rangeStart` narrowed to three days before its own now, mutated before the Timeline view first renders:

```
narrowed range: 2026-08-26T22:29:05Z → 2026-08-29T22:29:05Z
press 1   -> lit ['1']   | state 'last day'            | 2026-08-28 22:29 → 2026-08-29 22:29
press 7   -> lit ['7']   | state 'part of last 7 days' | 2026-08-26 20:50 → 2026-08-29 22:29
press 30  -> lit ['7']   | state 'part of last 7 days' | 2026-08-26 20:50 → 2026-08-29 22:29
press 90  -> lit ['7']   | state 'part of last 7 days' | 2026-08-26 20:50 → 2026-08-29 22:29
press all -> lit ['all'] | state 'all days'            | 2026-08-26 20:50 → 2026-08-30 10:38
after a free zoom -> lit chips: 0 | state 'custom span' | 2026-08-27 12:55 → 2026-08-29 18:32
```

Last day is unclipped and says so; 7/30/90 all produce the identical clipped window and report it as the tightest honest description; All days is a distinct window here and is not mislabelled "part of last day"; a free zoom lights nothing. The clipped window still ends at now, which is the defect the pre-clamp exists to prevent.

## Integration Seams

**None.** There is no line for the orchestrator to add to a file outside my write set. `web/board.css` was the one file outside the Scope declaration that needed a change, and rather than hand it back as a seam I applied it — see D-03 for why, and reverse it there if you disagree.

Two notes that are *not* seams but need your eyes at Step 8/9:

- **The release ritual is entirely undone**, per D9 and brief correction 4: `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` are all untouched. This is a user-visible control-set change, so it reads as a **minor** bump to me; the changelog title should say what was delivered — something like `Timeline Trailing Windows` rather than a codename.
- **`do-work/prose-backlog.md:10`** is an open entry against `timeline_browser_probe_test.go`'s "three of the four properties" doc comment (the landing probe, now at roughly line 1959). My rebuild changed clause (6)'s *shape* but not the clause *count* — still six — so per the exploration's instruction I left the comment and the backlog entry alone. A worktree builder may not write any `do-work/` path. The entry is still open and still accurate.

## Decisions

**D-01 — Split the Now-jump coverage out of the doomed calendar test instead of deleting it. DECIDE & STATE.**
The plan's task 4 and D8 both say to delete `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow` outright. The exploration flagged this as understated, and it was right: that test carried the only Node-lane coverage of `timelineNowJump`'s two degenerate cases (`queueEnd === now`, and `NaN` forecast — both REQ-235 defects) and of `timelineFirstOpenRowIndex`'s decide-the-scroll-after-the-row-refresh ordering with its own anti-vacuity guard (REQ-319). None of those functions is being deleted and none is covered elsewhere in the Node lane; the browser probe's Now clauses (1) and (2) do not reach the NaN forecast or the row-list ordering. I kept that half verbatim under `TestJavaScriptBehaviorTimelineNowJumpLandsOnTheOpenWork` and deleted only the calendar half. D8 is honoured in substance — no calendar-period assertion survives — and no coverage left with it. The other doomed test, `TestJavaScriptBehaviorTimelinePeriodChipsLandOnNowAndStepsStayAligned`, was deleted whole as D8 says, because it is purely calendar-period; the behaviours worth keeping from it (settle-clamp discipline, inverse steps, chip idempotence, no off-calendar sliver at the range edge) are all carried by the new trailing-window test's four cases.

**D-02 — Keep `timelineUtcDayStart` rather than deleting `timelinePeriodStart` outright. DECIDE & STATE.**
Task 3 lists `timelinePeriodStart` for deletion, but `timelineTypedWindow` calls it twice at its surviving day-granularity clamp (`board-timeline.js:565-566` at HEAD). Deleting it as listed would have shipped a `ReferenceError` in the client's typed-date path. `timelineUtcDayStart(epochMs)` is the day-only remainder, the typed-dates test slices it in place of the old one, and the day clamp is unchanged in behaviour. The exploration named this as a BLOCKER and it was real.

**D-03 — Touched `web/board.css`, which is outside the Scope declaration's "Files I will touch". DECIDE & STATE, but flagged loudly because it crosses my write boundary.**
The plan's own risk list anticipates this: *"the fix if the render is bad is either shorter chip labels with a `Window` control-label prefix (the Durations precedent) or a `.timeline-periods` rule in board.css — either way, do not decide it from the markup, decide it from the render."* The render was bad: at 500px the group measured 465.63px inside a 453px toolbar on Chrome for Testing 151.0.7922.174 and clipped the `›` arrow and the state span, where HEAD's three short chips fit. Of the two sanctioned fixes, shortening the labels would contradict the REQ's own GREEN condition, which names those exact five labels. So I added `.timeline-periods { flex-wrap: wrap; }` — one declaration, the exact mirror of `.timeline-range { flex-wrap: wrap; }` three rules below it, with the measured number and the browser build recorded in its comment. I also rewrote the adjacent `.timeline-period-state` comment, which restated the deleted "which level the window is exactly showing" vocabulary. Reversing this is one `git revert -n` of the two CSS hunks; the cost of reversing is a narrow-width clip, not a broken test.

**D-04 — "part of" is decided by coverage, not by comparing spans. DECIDE & STATE.**
A chip is reported as clipped when the settled window does not *cover* what the chip asked for — `settledStart > now - N days || settledEnd < now`. The obvious alternative, comparing the settled span against the requested span, misreports in both directions: `timelineZoomedWindow` floors any span at one hour and caps it at the bound span, so on a very young board a naive `settledSpan < requestedSpan` says "part of last day" for a window that is in fact the whole archive, and a `!==` comparison says "custom span" for the chip the reader just pressed. This is the risk the plan names; the coverage test is immune to both the floor and the 2% padding.

**D-05 — Ran the canonical gate without `QUEUE_KANBAN_BROWSER`, so the strict browser lane was skipped rather than run. DECIDE & STATE.**
`maintainer-verify.sh` runs `TestMaintainerStrictBrowserBehaviorLane` only when `QUEUE_KANBAN_BROWSER` is set or a chrome name is on PATH; neither holds here, so it prints its SKIP line and exits 0. Exporting the Chrome 151 path into that invocation would have pulled in `TestBrowserBehaviorCompletionCompanionsKeepReadableContrast`, which the brief records as failing at HEAD in this container on both Chrome 141 and 151 with no source changes — turning a pre-existing unrelated failure into this REQ's gate failure. Instead I ran the whole timeline browser probe family individually under Chrome 151 (14 tests, all pass) and left the gate as the repo defines it. **The strict browser lane was therefore not exercised as a lane on this machine.**

**D-06 — Rebuilt the landing probe's clause (6) around Fit-all-plus-one-step, with a new premise. DECIDE & STATE.**
The old case relied on the `week` chip navigating to the filtered REQ's own week, which trailing windows remove. Rebuilt as the plan directs: under the REQ-164 filter, press Fit all, then `timeline-period-next`, and assert the arrow reports itself disabled and the press does not move the window. The premise is read rather than restated, and it is a *different* premise from the old one: the fitted window is an outer bound on everything drawn under the filter, so a full screenful forward necessarily lands past all of it — *unless the pan clamps at the right-hand bound*. So the guard checks that the payload's whole range (which is exactly what the new "All days" chip spans, captured unfiltered in the same probe) has more than one screenful of room to the right of the fitted window. The chart-not-empty premise is clause (4)'s existing `t.Fatalf` on the same `filteredFit` snapshot, which runs first; restating it in clause (6) would be a guard no mutation can break. This is REQ-345's lesson applied — the case cannot pass vacuously on an empty chart or on a clamped pan.

**D-07 — `1` rather than `7` for the grouping probe's second window. DECIDE & STATE.**
`TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders` asserts `fitAllIds != lastDayIds`. The exploration flagged that trailing `1` is not a superset of the old calendar `day` and the assertion is data-dependent. The plan's `1` passed on this repo's live queue with a large margin, so I kept it; the risk note's fallback to `7` was not needed. The JSON keys and Go fields were renamed `day`→`lastDay` so a plain-text search for "day" no longer lands on a calendar concept.

**D-08 — The chip values are `1|7|30|90|all` and the DOM identifiers are unchanged, as D2 and D3 direct.** No decision of mine; recorded so the integrator can see the wart was accepted deliberately: `data-timeline-period` keeps a noun that now means "trailing window", because renaming it would churn `board.css`, four probe call sites and the no-match toolbar stub for no user-visible gain, and the REQ's GREEN condition names the attribute literally.

## Discovered Tasks

- **`.timeline-toolbar`'s children wrap at 900px and below on this board.** `toolbarChildRows` measured 4 at 1400px, 900px and 500px — the legend, the window group, the range fields and the zoom group each on their own line even at 1400px, giving a 166px-tall toolbar at 1400 and 243px at 900. This is pre-existing (HEAD measures the same) and unrelated to this REQ, but a four-row toolbar above the chart is a lot of vertical space before any data.
- **`timeline_browser_probe_test.go`'s landing-probe doc comment still says "three of the four properties"** while the test has six numbered clauses. This is the open `do-work/prose-backlog.md:10` entry; I did not drain it because a worktree builder may not write any `do-work/` path and the clause count did not change.
- **No test covers `renderTrailingWindowControls`'s "part of" branch.** The clipping *arithmetic* is pinned by `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow` case (2), and I verified the rendered text by hand on a synthetic young board (evidence above), but nothing in the suite would catch a regression in the branch that turns a clipped match into the words "part of last 7 days". It needs a fixture board younger than 90 days, which neither the live-board browser lane nor the current Node fixtures supply.
- **`@media (max-width: 480px)` in `board.css` covers only `.durations-chart-list`.** The timeline toolbar has no narrow-width rules of its own; the fix I added works by wrapping rather than by restyling. If the board ever claims 320px support properly, this toolbar is where to start.

## Cross-REQ Test Changes

Every change below is intentional: this REQ removes a rule (calendar-period windows), and per `general.md` § Cross-REQ Test-Break Rules and `lessons-kanban-board.md` REQ-237, every test asserting that rule's shape is in scope and is named here rather than quietly edited.

| Prior REQ | Test | What changed | Why the behaviour change is intentional |
|---|---|---|---|
| **REQ-235** (Add period zoom and jump to now on the timeline) | `TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow` | **Split.** The calendar half — Monday-aligned weeks, month-as-calendar-arithmetic, midnight-to-midnight days, period-index clamping at the range end, exact-level readback — deleted. The Now-jump half kept verbatim as `TestJavaScriptBehaviorTimelineNowJumpLandsOnTheOpenWork`. | The functions the calendar half drove (`timelinePeriodStart`, `timelineSteppedPeriodStart`, `timelinePeriodWindow`, `timelinePeriodAnchor`, `timelinePeriodLevelOfWindow`) no longer exist: the REQ replaces calendar periods with trailing windows. The Now button is untouched by this REQ, so its coverage stays — see D-01. |
| **REQ-235 / REQ-326** | `TestJavaScriptBehaviorTimelinePeriodChipsLandOnNowAndStepsStayAligned` | **Deleted whole**, with its comment block. | It tests only the deleted calendar machinery (chip anchoring on now vs the reader's position, calendar-grid stepping, next-then-previous as a calendar identity, edge-period alignment). Adapting it would mean writing new tests inside old scaffolding. The behaviours worth keeping — clamp-before-settle discipline, chip idempotence, inverse steps, no sliver at the range edge — are carried by `TestJavaScriptBehaviorTimelineTrailingWindowsEndAtNow` cases (1)–(4). The strict Node lane selects by `^TestJavaScriptBehavior` with no minimum count beyond non-zero, so the deletion cannot silently empty it; the lane still passes with 44 `TestJavaScriptBehavior*` functions. |
| **REQ-320** (Show and set the timeline window's start and end) | `TestJavaScriptBehaviorTimelineTypedDatesMoveTheWindow` | Slices `timelineUtcDayStart` in place of `timelinePeriodStart`/`timelineSteppedPeriodStart`/`timelinePeriodLevelOfWindow`; `lastDayStartMs` now reads `timelineUtcDayStart(boundEnd - 1)`; `periodRoundTrips` (day/week/month, with a `reparsedLevel` field) became `windowRoundTrips` over two literal `Date.UTC` pairs — one day and seven days — and the `reparsedLevel` assertion was dropped. | REQ-320's actual property — that rendering a window into the two date fields and parsing it back are inverses — is unchanged and still asserted, on the same shapes of window. What was dropped is the *secondary* claim that a typed pair can light a period chip, which is no longer true and should not be: trailing windows are not midnight-aligned, so a typed whole-day pair cannot equal a chip's window. That is a deliberate consequence of the REQ, not a weakening. The `SameDaySpanMs` assertion kept its value (exactly one day) and lost only its "so the Day chip can light for it" clause. |
| **REQ-240 / REQ-327** (Stop the timeline axis printing a fake minute) | `TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant` | The week-boundary oracle `timelinePeriodStart(tickMs, "week") === tickMs` became `(tickMs - TIMELINE_WEEK_ALIGNMENT_MS) % (7 * TIMELINE_DAY_MS) === 0`; the `timelinePeriodStart` slice dropped; the failure message re-pointed. | The property is unchanged — a week-long axis gap lands on the same Monday the rest of the view uses. Only the oracle moved, from a deleted function onto the constant the shipped axis code (`timelineTickAtOrBefore`) actually aligns to. Both constants are still *read off the renderer* rather than restated beside the test, which is REQ-322's rule; if Monday moves, `TIMELINE_WEEK_ALIGNMENT_MS` moves and this moves with it. |
| **REQ-320** | `TestBrowserBehaviorTimelineRangeFieldsShowTheWindowTheChartIsDrawnAt` | Two `[data-timeline-period="month"]` clicks became `="30"`; `monthFrom`/`MonthFrom` renamed `trailingFrom`/`TrailingFrom`; one message reworded. | The probe only needs *a window that is not the bounds* so a clamp is observable, and a 30-day trailing window is one. All four of its properties (cleared field restored, clamped commit shown honestly, in-range date applied exactly, mid-edit respected then released) are unchanged. |
| **REQ-235 / REQ-345** | `TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable` | Clause (3): `week` click → `7`, `currentWeek`/`afterStepFromCurrentWeek` renamed `trailingSevenDays`/`afterStepFromTrailingSevenDays`, comment rewritten for a now-anchored window. New `allDays` snapshot captured unfiltered. Clause (6): **rebuilt** — the two `filteredWeek` states replaced by one forward step from `filteredFit`, with a new room-to-the-right premise. | Clause (3)'s contract is unchanged: read the arrow's own verdict, check the press against it, never assert which regime the live queue is in — REQ-345's lesson, still honoured. Clause (6) had to be rebuilt because it was built on the `week` chip navigating to the filtered REQ's week, which trailing windows remove; the *property* it pins — the only deterministic arrow-refusal case in the suite — is preserved, with a premise the live queue cannot silently satisfy. See D-06. |
| **REQ-313 / REQ-338** | `TestBrowserBehaviorTimelineListsRowsBeneathUserRequestHeaders` | `day` click → `1`; `day`/`dayRebuild` JSON keys and `Day`/`DayRebuild` Go fields renamed `lastDay`/`lastDayRebuild`; two messages reworded. | The probe needs a *narrower window than Fit all* that lists a different REQ set; trailing `1` supplies one. All of its grouping, header-id, cell-headers, tab-stop and virtualization assertions are unchanged. The `fitAllIds != lastDayIds` assertion is data-dependent and was checked live — see D-07. |

## Remediation

**Commit:** `403d78e16784f220583d19b66effa5ed02ec8b03` on `worktree-agent-REQ-390-timeline-trailing-window-periods`, one file, `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (+88/-35). Not pushed, not merged.

### What the new premise is, and why the live queue cannot move it

The old clause took the window Fit all chose under the `REQ-164` filter and then *measured* whether the board had a screenful of room to its right. Both halves of the premise came from that one window, and its **width** is live queue data.

Root cause, exactly: the shared search filter matches by citation as well as by id and title (`web/board-filters.js` → `citationMatchedTicketId`). `do-work/working/REQ-390-timeline-trailing-window-periods.md` acquired `REQ-164` in its `citedTicketIds` when the queue owner recorded the merge notes, so the needle now matches **two** rows — REQ-164 (28 minutes long, 2026-08-10) and REQ-390 itself, still open, whose bar runs to the now-line. The filtered fit stretched from 20 days wide, against 1.91 days of bound padding. The test's own REQ file changed the fixture the test reads.

The new premise is built, not found, and it separates the two facts a refusal needs:

1. **Right edge past everything drawn** — still supplied by Fit all, which is an outer bound on the drawn extent by construction (`timelineFitWindow` fits `drawnExtent` plus breathing room, and `drawnExtent` is computed from `filterMatchedSegments`, so it does not move when the window does).
2. **More than one of the window's own screenfuls of room to its right** — now supplied by making the window **one hour wide** rather than by hoping the filtered set is narrow. The probe zooms in with the anchor pinned to the window's right edge: `handleTimelineWheel` clamps `anchorFraction` into `[0, 1]`, so a `WheelEvent` with `ctrlKey` and `clientX: 1e6` anchors the zoom exactly on the window's end, and `timelineZoomedWindow` at anchor 1.0 leaves `windowEndMs` untouched. Repeated until a press stops narrowing, that lands on `TIMELINE_MIN_SPAN_MS` (one hour) with the edge unmoved.

Two setup assertions pin the construction before the guard runs, so neither half can silently rot: the narrowed window's end must equal the fit window's end, and its span must be at most the one-hour floor.

Why the guard can no longer fire on live data: the room past any filtered extent is at least the bound padding, which is a fixed 2% of the board's range (`boundPaddingMs` in `web/board-timeline.js`), less the fit's breathing room. Clause (4) — already a hard assertion in the same test — caps the filtered fit at half the unfiltered span, which caps that breathing room at half the padding. So the room is at least ~0.96% of the board's range against a fixed one-hour window. **Residual condition, stated rather than hidden:** that clears one hour on any board whose data spans more than about 4.4 days. It is no longer a fact about *where* the queue's data sits, only about the board being more than a few days wide.

### RED evidence

The worktree's own `do-work/` is at the branch tip and still passes, so the failure was reproduced on the merged data by copying the main tree's `do-work/` beside a copy of the package (`resolveRepoRoot` walks up for `do-work/`; no git needed). Pre-fix file, merged queue:

```
--- FAIL: TestBrowserBehaviorTimelineNowAndFitAllLandSomewhereReadable (6.86s)
    timeline_browser_probe_test.go:2235: the board's whole range (2026-05-27 20:04 UTC → 2026-09-01 14:30 UTC)
    leaves 1.91 days to the right of the window Fit all chose under the one-REQ filter
    (2026-08-10 12:11 UTC → 2026-08-30 16:39 UTC), which is not more than that window's own 20.19-day span;
    ... (summary "2 REQs in the window ... 1 still open, measured to the now-line at 2026-08-30 07:20 UTC.")
FAIL   exit 1
```

### Pass counts

| Where | Runs | Result |
|---|---|---|
| Worktree data (branch tip queue), fixed clause | 6 | 6 pass, 0 fail |
| Merged data (main tree's `do-work/`), fixed clause | 5 | 5 pass, 0 fail |
| Adversarial fixture (below), fixed clause | 3 | 3 pass, 0 fail |
| Timeline family, `-run '^TestBrowserBehaviorTimeline\|^TestTimeline'` | 1 | exit 0, 57.5s |

### Step 3 — proving it is not luck on today's data

Two independent checks.

**A fixture board built to trip the old premise.** A synthetic `do-work/` tree (scratch only, not committed) with the shape that broke the live board, made worse: 18 filler REQs spread over 12 days so the whole board spans 12.06 days and the bound padding is 5.79 hours rather than the live board's 1.9 days; `REQ-164` archived 5 days back; and `REQ-777`, still `claimed`, whose body cites REQ-164 — so the needle's drawn extent runs to the now-line exactly as REQ-390 made it. Pre-fix file on that fixture:

```
timeline_browser_probe_test.go:2235: the board's whole range (2026-08-18 01:45 UTC → 2026-08-30 14:46 UTC)
leaves 0.20 days to the right of the window Fit all chose under the one-REQ filter
(2026-08-25 05:08 UTC → 2026-08-30 09:57 UTC), which is not more than that window's own 5.20-day span
FAIL   exit 1
```

Fixed file on the same fixture: 3 runs, all exit 0. Generator kept at `…/scratchpad/req390/make-fixture.py`.

**Mutation, to show the guard was not made unreachable by weakening the clause.** With `stepLandsOffTheData` forced to `return false` in `web/board-timeline.js` (scratch copies only, restored after), the clause fails on both boards — the assertions, not the guard:

```
merged data:  the step-forward arrow is enabled on 2026-08-30 15:49 UTC → 2026-08-30 16:49 UTC ...
              pressing the step-forward arrow moved the window ... and drew 0 segments there
fixture:      the step-forward arrow is enabled on 2026-08-30 08:57 UTC → 2026-08-30 09:57 UTC ...
              pressing the step-forward arrow moved the window ... and drew 0 segments there
```

Both readouts are exactly one hour wide and end on the fit's own right edge, which is the construction the two setup assertions pin.

### Gate

`bash _dev/tests/maintainer-verify.sh` from the worktree root, no `QUEUE_KANBAN_BROWSER` set: **exit 0**, direct status, unpiped. Its browser lane stayed in the default skipped state (`SKIP: no browser is available; strict browser behavior lane was not run.`).
