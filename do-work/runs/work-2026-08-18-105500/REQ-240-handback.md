# REQ-240 Hand-Back — Stop the Timeline axis printing a fake minute

**Branch:** `worktree-agent-REQ-240-timeline-axis-minute`
**Commit:** `4a90f80` — `[REQ-240] Stop the Timeline axis printing a fake minute`
**Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-240-timeline-axis-minute`

## File Manifest

| File | Verb | What |
|---|---|---|
| `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` | modified | Rewrote `timelineFormatAxisTick`; added `TIMELINE_AXIS_TICK_COUNT` and `TIMELINE_YEAR_MS`; `renderAxis` now reads its tick count from the constant instead of a local literal. |
| `skills/do-work-board/tools/queue-kanban/generate_test.go` | modified | Added `TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant` (+130 lines, one new test, no edits to existing ones). |

`git diff --stat` for the commit:

```
 .../tools/queue-kanban/generate_test.go            | 130 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-timeline.js       |  29 +++--
 2 files changed, 152 insertions(+), 7 deletions(-)
```

Nothing outside the declared write set. `durations.go` / `durations_test.go` untouched (I only *read* `rendererNumericConstant` from `durations_test.go`; no edit). `web/board.css` and `web/template.html` not needed — no integration seam there.

## The Fix

```js
function timelineFormatAxisTick(epochMs, spanMs) {
  var instant = new Date(epochMs);
  var calendarDate = instant.getUTCDate() + " " + TIMELINE_MONTHS[instant.getUTCMonth()];
  if (spanMs / TIMELINE_AXIS_TICK_COUNT < TIMELINE_DAY_MS) {
    return (
      calendarDate + " " +
      String(instant.getUTCHours()).padStart(2, "0") + ":" +
      String(instant.getUTCMinutes()).padStart(2, "0")
    );
  }
  if (spanMs >= TIMELINE_YEAR_MS) {
    return calendarDate + " " + instant.getUTCFullYear();
  }
  return calendarDate;
}
```

Format chosen: **`18 Aug 11:36`** — the existing shape with the real minute in place of the literal. It is character-for-character the same width as the `18 Aug 11:00` it replaces, so the axis geometry does not move at all (measured, below).

The branch condition changed from a span threshold (`spanMs <= 3 days`) to a **tick-gap** threshold (`spanMs / TIMELINE_AXIS_TICK_COUNT < TIMELINE_DAY_MS`). Reason in `## Decisions`.

## RED / GREEN

### RED run 1 — the reported defect, code present, assertion failing for its own reason

Probe drives `timelineFormatAxisTick` over the tick instants of the window `Now` produces (span = `TIMELINE_MIN_SPAN_MS`, start 11:26 UTC).

```
--- FAIL: TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant (0.25s)
    generate_test.go:3093: the Now window draws 7 ticks with only 2 distinct labels: ["18 Aug 11:00" "18 Aug 11:00" "18 Aug 11:00" "18 Aug 11:00" "18 Aug 12:00" "18 Aug 12:00" "18 Aug 12:00"]
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.699s
EXIT=1
```

Not a reference error — `timelineFormatAxisTick` already existed, so the assertion failed for exactly the reason it was written, and it reproduced the REQ's measured 7-ticks/2-distinct on the first run.

### RED run 2 — the extended assertions, one rule broken

After the fix was in place I put the *old* span threshold back (`spanMs <= 3 * TIMELINE_DAY_MS`) with everything else unchanged, to prove the newly-added four-day case bites:

```
--- FAIL: TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant (0.24s)
    generate_test.go:3137: the free zoom, four days window draws 7 ticks with only 5 distinct labels: ["15 Aug" "15 Aug" "16 Aug" "17 Aug" "17 Aug" "18 Aug" "19 Aug"]
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.650s
EXIT=1
```

### GREEN — with REQ-235's and REQ-233's tests checked explicitly

```
=== RUN   TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard
--- PASS: TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard (0.24s)
=== RUN   TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer
--- PASS: TestJavaScriptBehaviorTimelineKeyboardMovesTheSameWindowAsThePointer (0.24s)
=== RUN   TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow
--- PASS: TestJavaScriptBehaviorTimelinePeriodStepsOnCalendarBoundariesAndJumpsToNow (0.24s)
=== RUN   TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant
--- PASS: TestJavaScriptBehaviorTimelineAxisLabelsNameTheirOwnInstant (0.24s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	1.302s
EXIT=0
```

REQ-233's two keyboard tests: **PASS**. REQ-235's period/Now test: **PASS**. Neither was modified.

The probe reads `TIMELINE_MIN_SPAN_MS`, `TIMELINE_DAY_MS` and `TIMELINE_AXIS_TICK_COUNT` through `timelineProbePreamble`, and `TIMELINE_YEAR_MS` and `TIMELINE_MONTHS` through `rendererDeclarationLine`, so it cannot pass against numbers the shipped view does not use.

## maintainer-verify

Run from the worktree root at the exact tree state that was committed (re-run after the browser work, because building the comparison board involved a stash/pop):

```
Maintainer verification passed.
VERIFY_EXIT=0
```

Also `gofmt -l .` (empty) and `go vet ./...` (clean) in the tool module, and `go test ./...` for the whole module: `ok ... 13.018s`.

## Browser Evidence — the axis actually read on screen

Two boards generated from this worktree's own repo (235 REQs, 54 URs), served on loopback and driven in a real browser: **before** = the shipped `board-timeline.js` restored via stash (`/tmp/board-240-before`, port 8241), **after** = this commit (`/tmp/board-240`, port 8242). Labels read out of the live DOM (`.timeline-axis-label`), not inferred from source.

| State | BEFORE | AFTER |
|---|---|---|
| **Fit all** | `28 May, 11 Jun, 25 Jun, 9 Jul, 23 Jul, 6 Aug, 20 Aug` — 7/7 distinct | **identical**, 7/7 distinct |
| **Month** | `1 Jul, 6 Jul, 11 Jul, 16 Jul, 21 Jul, 26 Jul, 1 Aug` — 7/7 distinct | **identical**, 7/7 distinct |
| **Week** | `13 Jul, 14 Jul, 15 Jul, 16 Jul, 17 Jul, 18 Jul, 20 Jul` — 7/7 distinct | **identical**, 7/7 distinct |
| **Day** | `16 Jul 00:00, 16 Jul 04:00, 16 Jul 08:00, 16 Jul 12:00, 16 Jul 16:00, 16 Jul 20:00, 17 Jul 00:00` — 7/7 distinct | **identical**, 7/7 distinct |
| **Now** | `18 Aug 11:00, 18 Aug 11:00, 18 Aug 11:00, 18 Aug 12:00, 18 Aug 12:00, 18 Aug 12:00, 18 Aug 12:00` — **7 ticks, 2 distinct** | `18 Aug 11:21, 18 Aug 11:36, 18 Aug 11:51, 18 Aug 12:06, 18 Aug 12:21, 18 Aug 12:36, 18 Aug 12:51` — **7 ticks, 7 distinct** |

The four period/fit states keep byte-identical labels; only `Now` changes, which is the requirement.

One more live state, not in the five but worth recording — a free zoom out of `Day` into the 3-to-6-day band, which is the second half of the defect and was **visibly duplicating on the shipped board**:

| State | BEFORE | AFTER |
|---|---|---|
| Day then 3 × zoom-out (≈4.1 days) | `7 Jul, 8 Jul, 8 Jul, 9 Jul, 10 Jul, 10 Jul, 11 Jul` — **7 ticks, 5 distinct** | `7 Jul 10:50, 8 Jul 03:13, 8 Jul 19:36, 9 Jul 12:00, 10 Jul 04:23, 10 Jul 20:46, 11 Jul 13:09` — **7 ticks, 7 distinct** |

**Legibility, measured rather than assumed** (`getBoundingClientRect` on each label, sorted by x):

- `Now` at a 1754px axis: label width 72–73px, minimum gap between neighbours **162px**, zero overlaps.
- `Now` at a 781px axis (900px viewport, the tight case): gaps `0, 37, 36, 36, 37, 0`. **The BEFORE board measures the identical `0, 37, 36, 36, 37, 0` with identical 72px widths** — `18 Aug 11:00` and `18 Aug 11:36` are the same string length, so this state's geometry did not move by one pixel. The two zeroes are the first/last labels, which are `text-anchor: start` / `end`; they touch their neighbour, they do not overlap.
- The four-day band, the only state where labels got *longer* (30–36px → 66–72px): gaps `9, 43, 42, 39, 37, 0` at 781px. Same width class as the `Day` window that already ships (whose tightest gap is 3px at the same width), so no new overlap risk.

No console errors on either board while driving it.

I did **not** capture a PNG: the browser tool's file root is the main tree, and writing there is out of bounds for a builder. The label strings and bounding boxes above come from the live DOM, which is stronger evidence for a text defect than a screenshot anyway.

## P-A-U

### [PLAN]

Read, in the worktree: `CLAUDE.md`, `_dev/primes/prime-kanban-board.md` and its `## Lessons` links (REQ-227 for this view's own history, REQ-235 for what the Now button changed, REQ-233 for the keyboard path), and the four crew members named in the brief (`general.md`, `coding-guardrails.md`, `communication-style.md`, `testing.md`). Read the REQ body from the main tree read-only.

Then read the code before designing: `timelineFormatAxisTick`, `renderAxis`, and every function that assigns `timelineViewState` — `timelineZoomedWindow`, `timelinePannedWindow`, `timelineKeyboardWindow`, `timelinePeriodWindow`, `timelineNowJump`, and the fit assignment. That last sweep is what let me bound the problem: **every** window setter either routes through `timelineZoomedWindow` (which clamps to `TIMELINE_MIN_SPAN_MS`) or preserves an existing span, and the initial fit's bounds are themselves at least `TIMELINE_MIN_SPAN_MS` wide. So the window span is never below one hour, the tick gap is never below ten minutes, and **minute resolution is sufficient** — seconds would have been a guess.

Approach: test first (`tdd: true`), formatter only, tick positions untouched.

REQ-227's lesson ("draw in pixels so a zoom has no scale to invalidate") did not conflict — this change adds no coordinate arithmetic. REQ-235's lesson about RED runs starting as reference errors shaped how I sequenced the two REDs.

### [APPLY]

Code stayed inside the write set: two files, both declared. No `board.css`, no `template.html`, no `durations*`, no `VERSION`, no `CHANGELOG.md`, nothing under `do-work/`, nothing in the main tree. All scratch (`/tmp/qk-240`, `/tmp/qk-240-before`, `/tmp/board-240`, `/tmp/board-240-before`, log files) is under `/tmp`. The main tree carries no artifact from this build — checked: no stray PNG, and `.playwright-mcp/` there is pre-existing and gitignored (`.gitignore:1`).

### [UNIFY]

- `git diff --stat`: 2 files, +152 / −7 (above).
- `web/board-timeline.js` — checked the diff line by line: three hunks, one constant pair, one function body, one literal replaced by the constant it should always have been. `renderAxis`'s loop arithmetic is unchanged, so tick x-positions are bit-identical.
- `generate_test.go` — one new test appended, no existing test touched (`git diff` shows a pure addition at the end of the file).
- Linters: `gofmt -l .` empty, `go vet ./...` clean, `bash _dev/tests/maintainer-verify.sh` exit **0** (ShellCheck, shipped-package contract, strict JavaScript behaviour lane, both Go modules).
- No debug artifacts: no `console.log`, no commented-out code, no temp files in the worktree (`git status --short` clean after commit).

## Decisions

**D1 — Print the tick's real minute, keeping the existing format shape.** `18 Aug 11:36` rather than anything longer. Rejected: adding seconds (the tick gap can never be under ten minutes, so seconds carry no information and cost width); switching to a locale/ISO form (`2026-08-18 11:36` is 40% wider and the view is UTC-only by design, stated in `timelineFormatStamp`). The chosen form is the same character count as the string it replaces, which is why requirement 3 (legibility at the existing tick count and font size) is satisfied by construction, not by luck.

**D2 — Key the format on the gap between ticks, not on the window span.** The old `spanMs <= 3 days` encoded an assumption about the tick count without saying so: with seven ticks it meant "gaps of 12h or less". Written as `spanMs / TIMELINE_AXIS_TICK_COUNT < TIMELINE_DAY_MS` the rule states what it means — *show the time once two ticks can land on the same day* — and it stays true if the tick count is ever changed. That required hoisting `renderAxis`'s `var tickCount = 6` into `TIMELINE_AXIS_TICK_COUNT`, which is the whole reason the constant exists: the formatter and the generator must read one number, or the threshold silently rots the next time someone edits the other.

This also fixes a **second instance of the same defect that the REQ did not name**: at spans between 3 and 6 days the old code fell into the date-only branch with gaps under a day, so it repeated labels there too. Measured live on the before-board at ≈4.1 days: 7 ticks, 5 distinct (`7 Jul, 8 Jul, 8 Jul, 9 Jul, 10 Jul, 10 Jul, 11 Jul`). Requirement 2 says *at any window span*, so I read this as in scope rather than as widening it — it is the same formatter, the same one-line condition, and no change to tick generation. Flagging it explicitly because it is the one place the diff does more than the REQ's headline.

**D3 — Add the year once the window can repeat a day-and-month.** `Fit all` covers the whole capture history and only grows. Once that reaches back past a year, the date-only branch collides again: a two-year window puts ticks four months apart, and `18 Aug` appears in two different years. Three lines (`spanMs >= TIMELINE_YEAR_MS` → append `getUTCFullYear()`) close that permanently. It changes nothing today — this repo's history spans about three months, and the live `Fit all` labels are byte-identical before and after (table above). Rejected: dropping the day at long spans (`Aug 2026`), which is arguably prettier but changes labels at spans where the current ones are not yet wrong, and requirement 5 says leave those alone.

I considered *not* doing D3 on scope-discipline grounds. I kept it because "two ticks at different instants must not render identical labels" is stated as a property of the formatter at any span, and shipping a fix that is knowingly wrong one zoom level out is the shape of defect this REQ exists to close. It is trivially revertible on its own — the three lines are a self-contained branch.

**D4 — One test, seven windows, two assertions.** The probe drives `Now`, `Day`, a four-day free zoom, `Week`, `Month`, `Fit all`, and a two-year `Fit all`. Each window asserts (a) labels pairwise distinct and (b) the whole label matches a pattern built from *that tick's own* day/hour/minute/year. Assertion (b) is what makes the test unable to pass on a literal: every number in the label has to be one the instant carries. It also pins the format per window, so the three period levels and `Fit all` are held to the labels they already render — requirement 5, as a test rather than as a promise. I did not add a separate test per window; that would be test-per-method symmetry, and the interesting failure is the same one in every case.

**D5 — Not touched: tick generation.** The week window's 1.167-day interval and `Fit all`'s crushed bars are REQ-235's, per the brief. I did not design against them: keying on the tick gap is neutral to how the ticks are placed, so a future period-aware generator can replace `renderAxis`'s loop without touching the formatter.

## Discovered Tasks

- None new. The two adjacent problems I was told about (week ticks at 1.167-day intervals; `Fit all` spanning the whole history) are already recorded on REQ-235 and I left both alone. Worth noting for whoever picks up the tick generator: the date-only label branch is now reachable only at gaps of a day or more, so a period-aware generator that lands ticks on calendar boundaries would let that branch drop the year tier again.

## Integration Seams

**None.** No shared file touched — no `board.css`, no `template.html`, no `durations.go` / `durations_test.go`, no serial-only file (`VERSION`, `skills/do-work/VERSION`, `actions/version.md`, `CHANGELOG.md`). The one shared *helper* I used, `rendererNumericConstant`, lives in `durations_test.go` and I only call it; if the sibling builder changes its signature the merge will show it as a normal Go compile break, not a silent conflict.
