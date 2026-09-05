# Hand-back — REQ-585 (give the Activity view one scroll surface instead of a scroll box inside the scrolling board)

## Branch

- Branch: `worktree-agent-REQ-585-one-scroll-surface`
- Head commit: `94532c45163df7be4624d24407272ba1d37b1ebd` (`[REQ-585] Give the Activity view one scroll surface`)
- Base commit: `1f43608b`
- One commit. Worktree is clean; `git status --porcelain` is empty.

## File manifest

| Verb | Path |
|---|---|
| modified | `skills/do-work-board/tools/queue-kanban/web/board.css` |
| created (test) | `skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go` |

`git diff --numstat 1f43608b HEAD`:

```
265	0	skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go
33	8	skills/do-work-board/tools/queue-kanban/web/board.css
```

No other file was touched. Nothing under `do-work/` was staged or committed; this hand-back is the only main-tree write.

### What changed in `web/board.css`

Inside the `/* ---- activity view */` block only:

- `.activity-table-scroll` loses `max-height: 70vh` and `overflow: auto`. It keeps `margin-top` and `font-size`, so the table's cell metrics are unchanged.
- New rule `.board-main:has(> #view-activity.is-active) { padding-top: 0 }` — the padding move, scoped to this one view.
- `#view-activity` padding goes from `16px 24px 24px` to `0 24px 24px`.
- `.activity-summary` gains `padding-top: 40px`, which is the 24px plus 16px that came off the two rules above.
- The block comment no longer says the table borrows the timeline table's metrics for scrolling. It now says it borrows the timeline's **cell** metrics and explicitly not its scroll box, and it records why the paddings are zeroed with the measurement and the browser build.
- A short comment above `.activity-table-scroll` says the class name is now historical.

The `thead th { position: sticky; top: 0 }` rule was not touched, as instructed.

## P-A-U

### [PLAN]

Written before any code.

1. Add a browser probe first and run it against unmodified CSS, so the RED number is measured rather than quoted from the REQ. The probe builds a synthetic repo of 40 archived REQs, each carrying `created_at`, `claimed_at` and `completed_at` inside the Activity view's default 24-hour window, which is 120 transition rows. It generates the real static site, opens the real `index.html` in headless Chrome at 1600x900, clicks the Activity view, and reads `clientHeight` and `scrollHeight` from both `.board-main` and `.activity-table-scroll`, then scrolls the board 700px and reads the sticky header's offset from the board's top inner edge plus a count of rows painting above it.
2. Apply the M1 diff from the mockup: drop the table's `max-height` and `overflow`, zero the board's and the view's top padding for this view only, and put 40px on the summary line.
3. Scope the padding move with `:has()` rather than a class toggled in `board-controls.js`. `:has()` is one CSS rule in the file that already owns the layout, versus a new class name, a new toggle line in `applyView`, and a second place the rule can drift from. `board.css` already uses `:has()` at the markdown checkbox-list rule, so this is not a new dependency for the client.
4. Re-run the probe for GREEN, then run the Go lane, the Node lane, and the whole browser lane.

### [APPLY]

Coded exactly as planned, inside the write boundary. `web/board-controls.js` was in the boundary but was **not** touched, because the `:has()` route made a JS toggle unnecessary (see D-03).

One deviation from the plan's order, reported rather than hidden: between step 1 and step 2 I ran a third variant of the CSS to test the mockup's own justification. See D-02.

### [UNIFY]

`git diff --stat` (against base `1f43608b`):

```
 .../activity_scroll_browser_probe_test.go          | 265 +++++++++++++++++++++
 .../do-work-board/tools/queue-kanban/web/board.css |  41 +++-
 2 files changed, 298 insertions(+), 8 deletions(-)
```

Linters, from `skills/do-work-board/tools/queue-kanban`:

| Check | Result |
|---|---|
| `gofmt -l .` | no output (exit 0) |
| `go vet ./...` | exit 0 |

Files verified one by one:

- `web/board.css` — read the whole diff. Checked that the four changed declarations are all inside the activity block, that `thead th` sticky is untouched, that the new selector's specificity (one id, two classes) beats the `@media (max-width: 760px)` shorthand at line 2037 that also sets `.board-main` padding, and that the block comment's claims now match the rules under it. No `TODO`, no debug rule, no `!important`.
- `activity_scroll_browser_probe_test.go` — read the whole file. No `console.log`, no `debugger`, no `TODO`. Every name has two words. The probe's two magic numbers are named constants that the JavaScript reads through `strconv.Itoa` rather than restating, so changing the constant changes the probe.
- Grep over the added lines for `console.log|debugger|TODO|FIXME|XXX`: no matches.
- `git status --porcelain` after the commit: empty, so nothing was left uncommitted and nothing under `do-work/` was staged.

## Test evidence

All commands run from the worktree.

| # | Command | Exit |
|---|---|---|
| T1 | `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` | 0 |
| T2 | `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` | 0 |
| T3 | `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test ./... -run '^TestBrowserBehaviorActivityViewHasOneScrollSurface$' -count=1` | 0 |
| T4 | `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="…/Google Chrome" go test ./... -run '^TestBrowserBehavior' -count=1` | 0 |
| T5 | `gofmt -l .` and `go vet ./...` | 0 |

T1 reported `wall=45s tests=396 slowest-file=generate_test.go:9.66s limit=<30s`. T2 reported `wall=6s tests=59 slowest-file=javascript_behavior_c_test.go:2.25s`. The 30-second rule is per test file and the slowest file is 9.66s, so both lanes are inside it. T4 is the whole browser lane at 103s, run as insurance because the change edits a stylesheet other probes measure against; it is not one of the three focused commands.

T1 and T2 were re-run at the committed revision, after the last comment edit, so these exit codes describe exactly what is on the branch.

### RED and GREEN, measured

A real engine ran both. Browser: **Chrome 152.0.7977.76** (`HeadlessChrome/152.0.0.0`), macOS arm64, window 1600x900, 120 transition rows.

RED is the same probe run against the unmodified stylesheet (I restored `web/board.css`, ran it, then restored the change):

| Measurement | RED (before) | GREEN (after) |
|---|---|---|
| `.board-main` clientHeight / scrollHeight | 665 / 677 | 665 / 3657 |
| `.board-main` scrolls | **true** | **true** |
| `.activity-table-scroll` clientHeight / scrollHeight | 530 / 3509 | 3509 / 3509 |
| `.activity-table-scroll` scrolls | **true** | **false** |
| `[board scrolls, table scrolls]` | `[true, true]` | `[true, false]` |
| board `scrollTop` reachable at a 700px request | 12 (that is the maximum) | 700 |
| column header top, relative to the board's top inner edge, after scrolling | 55.5px | **0.0px** |
| column header height | 28.5px | 28.5px |
| rows painting above the stuck header | 0 (nothing scrolled, so nothing could) | **0** |
| summary line's text top, relative to the board's top inner edge | 40.0px | 40.0px |
| `.board-main` computed `padding-top` on the Activity view | 24px | 0px |
| `.board-main` computed `padding-top` back on the Kanban view | 24px | 24px |

RED matches the REQ's captured RED of `[true, true]`. GREEN matches the REQ's target of `[true, false]` with the header still visible after 700px of scroll and no row above it. The board's RED `scrollHeight` of 677 against a `clientHeight` of 665 is the 12px of slack the 70vh cap left, which is why the RED run cannot even reach 700px of scroll.

Two extra checks the REQ did not ask for but that the change makes cheap:

- **The summary line does not move.** Its text starts 40.0px below the board's top inner edge before and after. The padding moved boxes without moving a pixel of what the reader sees, at the default width.
- **The scoping holds.** After clicking back to the Kanban view in the same session, `.board-main` computed `padding-top` is 24px again. The probe asserts both halves, so a rule that leaked to every view fails the test.

### Responsive check

Re-ran T3 at a 700px-wide window, below the stylesheet's `max-width: 760px` breakpoint where `.board-main` padding becomes `18px 16px 48px`. Result: board 478/3649, table 3509/3509, header offset 0.0px, 0 rows above it, Kanban-view padding 18px. The Activity-only rule wins inside the media query as intended. The one cosmetic difference is that the summary line keeps the 40px pad there instead of the 34px it had before; see D-04.

Accessibility and errors: no console errors were raised by any probe run, and the change adds no interactive element, so keyboard reachability is unchanged. Removing the inner scroll box removes a focusable scroll region rather than adding one.

## Lesson evidence

Read in full:

- `_dev/primes/prime-kanban-board.md` — the maintainer prime for this tool.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — the shipped prime, its Read first, Do not edit, Traps and Stakes.
- `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `frontend.md`.

Both lesson satellites the REQ lists as dropped for budget were opened and searched for the entries that touch `web/board.css` and the browser probe lane, as the brief directs:

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — present. Searched for scroll, sticky, overflow, CSS and probe entries. Nothing in it constrains this change.
- `_dev/primes/lessons-kanban-board.md` — present. Two entries changed how I worked:
  - **REQ-291 (the browser behavior probe lane):** an unrendered element measures as zeros, so a probe's default failure is a successful-looking measurement of nothing. Applied directly: this probe refuses to compare heights until the row count matches what the fixture produced and until the board, the table and the header all report a positive height. Without that guard, a table that failed to render would report `0 === 0` and read as the fix working.
  - **REQ-322 (a constant a decision turns on must be read by the test, not restated beside it):** the 700px scroll distance and the 40-request fixture size are named Go constants, and the JavaScript reads the scroll distance through `strconv.Itoa` instead of carrying its own literal.
  - Also noted from the maintainer prime: a measured face is per-browser and a constant that does not name its build cannot be argued with. The stylesheet comment and this hand-back both name Chrome 152.0.7977.76 beside the measurement.

No listed path was missing.

## Decisions

D-01 is the orchestrator's (M1, one page scroll with the column header stuck to the board's top edge). Mine start at D-02.

### D-02 — Verified the mockup's justification for the padding move instead of taking it on faith. DECIDE & STATE

The mockup and the REQ both claim the board's 24px top padding would let rows show through above the stuck header. My own reading of the CSS specification said the opposite: a sticky element in a scroll container is constrained by the scrollport, which is the padding box, so the header should have pinned flush to the visible top edge and covered everything under it.

Rather than guess, I ran a third CSS variant — only `max-height` and `overflow` removed, no padding move — and measured it. The result: header offset 24.0px, **2 rows visible above it**. Chrome pins the sticky header to the container's **content** box, not its padding box. The mockup is right and my reading was wrong, so the padding move stays.

This cost one extra probe run and it is why the stylesheet comment now states the rule as measured behavior with the browser build beside it, instead of repeating a claim nobody had checked. Reversible and low reach: it changed no code, only what I trusted.

### D-03 — Scoped the padding move with `:has()` in CSS rather than a class toggled in `board-controls.js`. DECIDE & STATE

The brief allowed either. The CSS route is one rule, `.board-main:has(> #view-activity.is-active) { padding-top: 0 }`, in the file that already owns the layout. The JavaScript route would need a new class name, a new line in `applyView`, and a second place where "this view is the Activity view" is decided — the exact second-definition shape the board's own comments warn about elsewhere.

`:has()` is not a new dependency: `web/board.css` already ships `.markdown-body ul:has(> li > input[type="checkbox"])`, and the strict browser lane targets current stable Chromium. `web/board-controls.js` is therefore unmodified.

The selector reads `.is-active` rather than `:not([hidden])` because `.view-panel.is-active` is already the rule that decides whether a panel is displayed, so the two cannot disagree.

### D-04 — Wrote the summary's replacement padding as a single `40px` instead of tracking the board's padding per breakpoint. DECIDE & STATE

`.board-main` has two top paddings: 24px normally and 18px below the 760px breakpoint. The summary's 40px restores 24 + 16 exactly at the default width, and comes out 6px taller than the old 18 + 16 = 34px on a narrow window.

Making it exact would mean a custom property on `.board-main`, a second declaration inside the media query, and a `calc()` on the summary — three rules to buy six pixels of leading above one line of text on a phone-width board. I wrote the one constant and recorded the trade in the stylesheet comment so the next reader does not think it was an oversight. Measured and confirmed at 700px: the layout is correct, only the gap differs.

Fully reversible: if the 6px is ever unwanted, one declaration inside the existing media query fixes it.

### D-05 — Kept the class name `.activity-table-scroll` although it no longer scrolls. DECIDE & STATE

Renaming it needs an edit to `web/template.html`, which is outside my write boundary. I left the name and added a two-line comment saying it is historical. See Integration seams for the exact change if the orchestrator wants it; it is cosmetic and safe to skip.

### D-06 — Built the probe on a synthetic 40-REQ fixture rather than on the real `do-work/` tree. DECIDE & STATE

`generateLiveSiteInDir` renders the repository's own queue, which several probes in this lane use. For this measurement it would make the test depend on how busy the queue happened to be in the last 24 hours: a quiet day produces too few transition rows to overflow anything, and the probe would pass while measuring nothing. The fixture produces exactly 120 rows every time, and the probe asserts it got 120 before it compares any height.

## Discovered Tasks

- The Timeline view has the same nested-scroll shape in `.timeline-scroll` (`web/board.css`). Not touched, per the brief; the orchestrator is capturing it separately. → report only, `impact-user-visible`
- `.activity-table-scroll` is now a class named for a scroll box that does not scroll. Renaming it to something like `.activity-table-frame` touches `web/template.html`, `web/board.css` and the new probe. Cosmetic naming debt, no behavior. → report only, `impact-internal`
- The Activity view now renders every transition row in the window as real DOM. At the live board's current volume this is a few hundred rows and measured fine (3657px of content, no perceptible lag in the probe). It was previously capped in height but never in row count, so this change does not add rows — but it does mean a very wide window setting, such as 7 days on a busy queue, now lays out every row at full height instead of inside a 70vh clip. `crew-members/frontend.md` suggests virtualizing lists above roughly 100 visible items. Worth a look only if someone reports the 7-day window feeling slow. → report only, `impact-internal`

## Integration seams

Nothing is required. One optional change, if the orchestrator wants D-05 resolved in the same merge:

- In `skills/do-work-board/tools/queue-kanban/web/template.html`, line 488, `<div class="activity-table-scroll">` would become `<div class="activity-table-frame">`, with the same rename applied to the six `.activity-table-scroll` selectors in `web/board.css` and the four occurrences in `activity_scroll_browser_probe_test.go`. I did not do it because `web/template.html` is outside my write boundary. If it is applied, the browser probe must be re-run, since it addresses the table through that class.
