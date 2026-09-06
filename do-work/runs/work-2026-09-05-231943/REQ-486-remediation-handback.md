# REQ-486 remediation hand-back

REQ-486 is the addendum that made the board's By UR groups collapsible and put a
progress summary strip on both UR surfaces. An independent three-lens review
scored it 79% / Partial with three demonstrated findings. All three are closed,
plus the two minor items named with them.

- Branch: `worktree-agent-REQ-486-collapsible-ur-progress-summaries`
- Head: `cb6cfd5`
- Worktree: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-486-collapsible-ur-progress-summaries`
- Not pushed. Nothing staged or committed in `/home/user/skill-do-work`.

## Files changed (all inside the declared write set)

- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js`
- `skills/do-work-board/tools/queue-kanban/web/board-core.js`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go`
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go`
- `skills/do-work-board/docs/board-guide.md`

No release change. `VERSION` is already 0.304.0 from this request's own second
increment and the changelog entry already covers this work.

## F1 — `Remaining ~0 min` with no qualifier

**What was wrong.** Each member's remaining contribution was floored at zero, as
the request asks, but nothing counted the floor. `userRequestSummaryRemainingText`
then took the clean path (`knownForecastCount > 0`, `unknownForecastCount === 0`)
and printed a bare `~0 min`. A reader could not tell "nearly done" from "every
member has already blown its estimate". On the real board that was four of the
five user requests with a live claim.

**Fix.** `summarizeUserRequestProgress` carries a new `overrunForecastCount`, and
the formatter renders it.

**Wording, and why.** `~0 min (2 over estimate)`. Two reasons for that phrasing
over "overrun" or a sentence: the Active figure already prints `2 excluded`,
`1 unmeasured`, `1 unknown`, so `2 over estimate` is the same N-word grammar a
reader has already learned on the line above; and "over estimate" is plainer
English than "overrun" for a non-native reader. Where both qualifiers apply the
figure reads `at least ~2h 30m (1 unknown, 1 over estimate)`.

**Evidence.** New assertions in
`TestJavaScriptBehaviorUserRequestProgressForecastsRemainingTime` cover a user
request (UR-803) whose two claimed members have each run past their saved
estimate — 600 minutes against 20, 240 against 45 — so every known contribution
floors and the sum is a true zero.

- RED (counter removed): `forecast rollup = {RemainingMinutes:150 KnownForecastCount:4 UnknownForecastCount:1 OverrunForecastCount:0 ...}` — FAIL
- RED (qualifier removed from the text): `remaining metric = "at least ~2h 30m (1 unknown)", want it to carry "1 over estimate"` — FAIL
- GREEN: `--- PASS: TestJavaScriptBehaviorUserRequestProgressForecastsRemainingTime (3.11s)`

## F1b — a rejected claim stamp counted as zero elapsed inside Remaining

**What was wrong.** `liveClaimElapsedMinutes` returns null for an unparseable or
future-skewed stamp. Active disclosed it with the clock-skew marker; the forecast
did `Math.max(0, estimateMinutes - (liveMinutes || 0))` and charged the member
its full estimate as a KNOWN forecast, with no qualifier at all.

**Fix, and why this mechanism.** A claimed member whose stamp the board refuses
is now an UNKNOWN forecast member, exactly like a member with no estimate and no
confident fallback. That is the disclosure channel the request already
requires ("disclosed as unknown rather than treated as zero", "preserve the
known forecast as explicitly partial"), so it needs no third qualifier and no new
field. How much of the estimate has been spent is precisely what the bad stamp
hides, so the member's remainder is genuinely unknown, not zero.

**Evidence.** UR-804 in the same test: one claimed member stamped 2099-01-01
carrying a saved 130-minute estimate.

- RED (branch disabled): `skewed-claim rollup = {RemainingMinutes:130 KnownForecastCount:1 UnknownForecastCount:0 ...}` — FAIL, which is exactly the pre-fix behaviour
- GREEN: the member is unknown, `remainingIsPartial` is true, and the metric reads `unknown`

## F2 — the browser probe measured a layout two of its three widths never had

**What was actually wrong — two leaks, not one.** The review named the docked
drawer, and that is real: `.detail-drawer` is a grid column above the 760px
breakpoint, and `measureProbePage` never closed it between widths. Closing it
alone did NOT restore the wide layout — it made things worse (2px). Instrumenting
the body grid found a second leak the review had not seen: the probe's own result
node, a plain `<pre>` appended to `<body>`, lands in the auto-sized second grid
column. Once it holds a measurement that is one ~32,000px-wide line of JSON, the
auto column takes the whole viewport and the board column collapses to nothing.
With the drawer open the two partly cancelled, which is why the corrupted run
still produced plausible-looking 259px and 579px numbers.

**Fix.** Both page-state leaks are reset: the result node is taken out of the
layout before it is attached (`position: fixed`, off-screen, 1px wide), and any
open drawer is closed at the start of every measurement.

**Re-logged `.ur-group` box, both themes, full run:**

| Case | Before this fix | After | Pre-rewrite reference |
|------|-----------------|-------|-----------------------|
| 320px | 16.0-289.0 (273.0) | 16.0-289.0 (273.0) | 273 |
| 768px | 28.0-287.0 (259.0) | 28.0-725.0 (697.0) | 697 |
| 1280px | 28.0-607.0 (579.0) | 28.0-1237.0 (1209.0) | 1209 |

The isolated subtest now agrees with the full run: `-run '.../light/1280px$'`
alone reports `1209.0 wide`, the same as inside the sequence.

**The lock-in, which is the real deliverable here.** `checkUserRequestProgressStrip`
now asserts the measured group box is the width the case's label claims — the
viewport less the board's own horizontal chrome (47 CSS px at 320, 71 at 768 and
1280), with a 96px allowance for future padding work. Its absence is why the
rewrite passed review the first time.

- RED (drawer close removed): `measured 259.0 CSS px wide (28.0-287.0) at the case labelled 768 px` and `measured 579.0 ... at the case labelled 1280 px` — reproduces the review's numbers exactly
- RED (result-node styling removed): `measured 2.0 CSS px wide (28.0-30.0)` at 768 and 1280
- GREEN: all six subtests pass, both themes

## F3 — the tick ordering was documented as asserted and was not

**Fix.** New probe
`TestJavaScriptBehaviorTickRefreshesExistingSurfacesBeforeTheSummaryPass` makes
the order observable. It replaces `document.querySelectorAll` so that only the
summary pass's own selector, `[data-ur-summary-id]`, throws, then drives one
`refreshTickingSurfaces` and checks two things: the injected failure really fired
once and escaped the tick (otherwise the probe proves nothing), and the existing
claim stopwatch still advanced 1m 30s to 2m 30s. It can only have advanced if the
older pass had already run.

**Ablation — swap alone, nothing else:**

- RED: swapping the two lines in `board-core.js` gives `claim stopwatch went "1m 30s" -> "1m 30s"` — FAIL, and it is the only test in the lane that reds
- GREEN at HEAD

The two comments that made the false claim are corrected: the freeze-guard
probe's doc comment now says it holds the narrowing and points at the new probe
for the ordering, and `board-core.js` says the order is asserted by a probe that
throws in the summary pass.

## F4 — two one-line fixes

`TestJavaScriptBehaviorUserRequestSummaryPathCarriesNoCatchAndNoCompletedAt`
slices the summary path out of the assembled page (`generateLiveSite`) and
asserts it carries no `try {` / `try{` / `catch (` / `catch(` and no
`completedAt`. Matched on syntax rather than on the bare words on purpose: the
module's own comments state both rules in prose, and the word "registry" in
`refreshUserRequestSummaryNodes` contains "try" — which is why the original
manual greps returned 1 instead of 0.

- RED: wrapping the summary pass's selector read in try/catch and adding a
  `request.completedAt` read fails all three tokens
- GREEN at HEAD

`HeadPrecedesDetil` is spelled `HeadPrecedesDetail`, so a plain-text search for
the correct spelling now finds the field and both of its reads, not only the JSON
tag.

## Docs

`skills/do-work-board/docs/board-guide.md` now documents the `N over estimate`
qualifier and states that a claim stamp the board refuses makes that member's
remaining time unknown.

## Lanes

Run individually, not through the canonical gate. Each reports its own exit line.

| Lane | Command | Exit | Result |
|------|---------|------|--------|
| Fast stage | `bash _dev/tests/maintainer-verify.sh` | 0 | `Maintainer verification passed.`, gate wall 78s; queue-kanban tests=402, do-work-cli tests=782 |
| Strict JavaScript | `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` | 0 | `go-test budget: ... wall=6s tests=72`; verbose run: `ok ... 6.248s`, 72 pass, 0 skip, 0 fail |
| Strict browser | `bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` | 0 | `go-test budget: ... wall=97s tests=35`; verbose run: `ok ... 96.083s`, 35 pass, 0 skip, 0 fail, Chromium 141 (HeadlessChrome/141.0.0.0) |

No lane skipped. The fast stage prints one SKIP line for the staged-skills,
updater and installer probes, which are heavy-only and unrelated to this work.

## Findings from the review NOT addressed

These were reported-only and outside the three the task named: the contrast
measured against `<body>` rather than `.ur-group`'s own surface; `Math.round`
printing 100% for an unfinished user request; the untested missing-request
narrowing; the duplicate `user_request_clipboard_browser_probe_test.go` entry in
the write set; and the nits.
