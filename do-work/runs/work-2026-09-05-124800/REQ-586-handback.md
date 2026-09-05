# Hand-back — REQ-586 (keep the board top bar to one line)

## Branch

- Branch: `worktree-agent-REQ-586-top-bar-one-line`
- Head commit: `cc56324f` — `[REQ-586] Keep the board top bar to one line`
- Base: `59f169d0` (contains REQ-585's merge `c08ac2b4`)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1248/worktree-agent-REQ-586-top-bar-one-line`
- Working tree clean apart from this hand-back file, which is in the main tree and was never staged.

## File manifest

All paths are inside `skills/do-work-board/tools/queue-kanban/`. Every one of them
is in the REQ's declared `write_set`; nothing outside it was touched.
`web/board-activity.js` was in the write set but needed no change — the chips keep
their ids and their `setActiveButton` call, so only their home moved.

| Verb | Path | What changed |
|---|---|---|
| modified | `web/template.html` | identity block became one row (wordmark, dot, project, dot, clock, clipped full stamp); view buttons reordered; `#activity-window-group` moved out of `.board-controls` into `#view-activity` inside a new `.activity-summary-row`, losing its `hidden` attribute and its "Touched in" label |
| modified | `web/board.css` | `.board-identity` is a baseline flex row with `white-space: nowrap`; new `.board-identity-dot` and `.board-generated-clock`; `.board-generated` is now clip-hidden; new `.activity-summary-row`; `.activity-summary` keeps its 40px pad and gives up its bottom margin to the row |
| modified | `web/board-controls.js` | deleted the `#activity-window-group` hand toggle in `applyView`; added `generatedStampTooltipText` and `renderTopBarIdentity`, called from `wireControls` |
| modified (tests) | `javascript_behavior_c_test.go` | two new tests plus two helpers, `viewTargetsInDocumentOrder` and `sliceMarkupElementAfter` |

No file was created or deleted.

## P-A-U

### [PLAN]

Read before coding: the five crew-member files the brief names plus `testing.md`,
`_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`,
and both lesson satellites (see Lesson evidence).

1. **Identity.** `GENERATED_AT_DISPLAY` is substituted **once** (`generate.go`
   uses `strings.Replace(..., 1)`, and `generate_test.go` fails on any leftover
   placeholder), so the placeholder can appear in exactly one place in the
   template. `board.js` line 60 appends a ticking relative-time node into
   `#board-generated`, and `board.js` is not in the write set. That settles the
   split the brief asked me to check: `#board-generated` keeps the server's full
   stamp and the ticking age, and is clipped out of the visible line; a new
   `#board-generated-clock` carries the visible short time, derived in the client
   from `boardData.generatedAt` (the same instant, same zone, cut to the minute).
   The tooltip is assembled from `#board-generated`'s live child text, so the age
   reaches it without a second one-second ticker.
2. **Chips.** Move the markup into `#view-activity` beside `#activity-summary`
   inside a `.activity-summary-row` flex row, delete the `hidden` attribute and
   the `applyView` toggle. Keep the 40px top pad on `.activity-summary` itself
   rather than moving it to the row, and align the row on baselines: the summary
   is then never pushed down, so REQ-585's probe still reads exactly 40px. The
   wiring in `wireControls` selects `[data-activity-window]` document-wide, so it
   survives the move untouched.
3. **View order.** Reorder the six `data-view-target` buttons in the template.

Test plan: one Node-lane test for the chips (markup ancestry, chip inventory,
the deleted toggle, the view order, and the shipped window writer driven for
real), one for the identity renderer. The one-line layout is measured in Chrome.

### [APPLY]

Coded exactly as planned, inside the write set. The only planned file left
untouched was `web/board-activity.js`, which needed no edit.

### [UNIFY]

`git diff --stat` for the commit:

```
 .../queue-kanban/javascript_behavior_c_test.go     | 294 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-controls.js       |  48 +++-
 .../do-work-board/tools/queue-kanban/web/board.css |  54 +++-
 .../tools/queue-kanban/web/template.html           |  50 ++--
 4 files changed, 420 insertions(+), 26 deletions(-)
```

Linters and per-file review:

- `gofmt -l .` in the package — no output, exit 0.
- `go vet ./...` in the package — clean, exit 0.
- Debug artifacts: `git diff -U0 | grep -E '^\+' | grep -E 'console\.log|debugger|TODO|FIXME|XXX'` — no matches.
- `web/template.html` — read the whole diff. The identity keeps `id="board-project"`
  and `id="board-generated"`; the placeholder still appears exactly once; the
  separators are `aria-hidden`; the moved chip group keeps its four
  `data-activity-window` values, its `aria-pressed` states, its `role="group"`
  and its `aria-label="Activity window"` (which is what still names the group for
  assistive tech after the visible "Touched in" label was dropped).
- `web/board.css` — read the whole diff. The REQ-585 Activity block is extended,
  not undone: `.board-main:has(> #view-activity.is-active)` and the 40px summary
  pad are byte-identical to what REQ-585 shipped, including its comment.
- `web/board-controls.js` — read the whole diff. One deletion (the no-op toggle),
  one addition (two functions and their call). The comment above the remaining
  toggles still describes them correctly.
- `javascript_behavior_c_test.go` — read the whole diff; both new tests fail for
  the right reason before the change (see Test evidence) and neither asserts on
  anything a refactor would move.
- Frontend checklist: no console errors in the headless render; keyboard reach
  unchanged (the chips are still `<button>`s in tab order, now inside the view);
  the clipped full stamp keeps the date readable for screen readers; no
  dependency added.

## Test evidence

Every command was run from the worktree. Exit statuses are the runner's own.

| # | Command | Exit |
|---|---|---|
| T1 | `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...` | 0, 61 tests, wall 7s |
| T2 | `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...` | 0, 397 tests, wall 46s |
| T3 | `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test ./... -run '^TestBrowserBehaviorActivityViewHasOneScrollSurface$' -count=1` | 0, PASS in 1.32s |

T2's wall time is 46s for the whole package; no single test file exceeded the
30-second budget (`slowest-file=generate_test.go:9.25s`).

### RED

Run before any implementation edit, same command as T1 narrowed to the two new
tests. Both failed on their assertions, not on compilation:

```
--- FAIL: TestJavaScriptBehaviorActivityWindowChipsRenderInsideTheActivityView
    the top bar still carries the Activity window chips; they belong on the Activity summary line
    the Activity view is missing the window chip "data-activity-window=\"6\""
    the Activity view is missing the window chip "data-activity-window=\"24\" aria-pressed=\"true\""
    the Activity view is missing the window chip "data-activity-window=\"48\""
    the Activity view is missing the window chip "data-activity-window=\"168\""
    board-controls.js still toggles #activity-window-group by hand; #view-activity hides it now
    the view pill reads [board calendar durations timeline activity testing], want [board activity calendar timeline durations testing]
--- FAIL: TestJavaScriptBehaviorTopBarIdentityIsOneLineWithAFullStampTooltip
    the identity block is missing one of its four parts: <div class="board-identity"> …
```

One captured assertion was wrong about shipped behavior and was corrected in
RED, before the implementation: the REQ asks the summary to read "in the last 48
hours" after the 48h chip, but `activityWindowPhrase` has always spelled whole-day
windows in days, so the sentence reads "in the last 2 days". Changing that
phrasing is not this REQ's business (its Constraints say the window values and
behavior stay exactly as today), so the test asserts the words the board actually
uses and carries a comment saying why. See D-05.

### GREEN

Same command after the implementation: both tests pass, and the full Node lane
(T1) and full package (T2) are green. The 48h and 6h cases drive the shipped
`applyActivityWindowSelection`, so the summary sentence and all four
`aria-pressed` values are read back from the real writer.

### Identity-line measurement (the layout fact the Node lane cannot see)

Measured by hand, since neither the Node lane nor a shipped browser probe covers
the top bar. Method: built the pre-change package from `HEAD` web assets into a
scratch copy, generated two static boards from the same `--repo-root`
(`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`, so the project name is
the real `skill-do-work2`), appended a measurement script to each `index.html`,
and dumped the DOM from headless Chrome after clicking through to the Activity
view. Browser: **Chrome 152.0.7977.76 headless** (same build the REQ-585 probe
reports). Each measurement carries its own `location.href`, and before/after came
from two different files.

| Measurement at 1400x900 | before | after |
|---|---|---|
| `.board-topbar` height | 126.0 px | **68.0 px** |
| `.board-identity` height | 85.0 px | 27.0 px |
| visible top-bar pills | Filters, View, Activity window | Filters, View |
| identity visible text | `do-work/ queue skill-do-work2 Generated 2026-09-05 13:42 UTC 12s ago` | `do-work/ queue · skill-do-work2 · 13:42 UTC` |
| identity `title` | (none) | `Generated 2026-09-05 13:42 UTC · 14s ago` |
| chips inside `#view-activity` | false | true |
| chips inside `.board-topbar` | true | false |
| view order | board, calendar, durations, timeline, activity, testing | board, activity, calendar, timeline, durations, testing |

68.0 px is the REQ's GREEN number exactly, and the identity's 27 px is one line
of an 18 px wordmark. A screenshot of the after board on the Activity view shows
the summary line reading `251 transitions across 46 REQs in the last 24 hours
[6h][24h][48h][7d]` with 24h pressed:
`/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/ab1d5ff6-3b27-4868-aaaf-f40f63a49993/scratchpad/activity-after.png`
(session scratch, not committed).

REQ-585's probe (T3) reports `summaryTextTop=40.0` with the chips in place, which
is the number the brief required it to keep.

Narrow widths: headless Chrome clamps its window to a 500 px minimum, so 320 px
and 375 px could not be measured. At 500 px the top bar drops from 264.5 px to
165.0 px, and the page's 26 px of horizontal overflow (`scrollWidth` 526 against
`clientWidth` 500) is identical before and after, so it is pre-existing and not
something this change introduced. See DT-2.

## Lesson evidence

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — read (the
  REQ lists it as dropped for budget; the brief restores it). The 0.295.1 entry
  applies directly: fragment order in the manifest is real, and a probe that
  assigns state by hand cannot see the bug — so both new tests drive shipped
  functions (`applyActivityWindowSelection`, `renderTopBarIdentity`) rather than
  setting `aria-pressed` or the tooltip themselves.
- `_dev/primes/lessons-kanban-board.md` — read the entries touching the web
  assets and the behavior lane (REQ-185 on the counted lane, REQ-376 on measuring
  text on the real element).
- `_dev/primes/prime-kanban-board.md` and
  `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — read in full.
  Two traps shaped the work: `web/` is embedded, so every measurement above came
  from a freshly built binary, not from files on disk; and a claim about pixels
  needs a render, which is why the identity line was measured before and after
  rather than argued from the CSS.
- No listed lesson path was missing.

## Decisions

**D-01 — `#board-generated` carries the full stamp; a new `#board-generated-clock`
carries the visible time. DECIDE & STATE.** The brief asked me to check what
writes into `#board-generated` before choosing. Two facts decide it: the
`GENERATED_AT_DISPLAY` placeholder is substituted once and `generate_test.go`
fails on a leftover, so it can only appear in one place; and `board.js` (outside
the write set) appends the ticking age into `#board-generated`, so that element
cannot be the short visible clock without showing the age beside it. Reversible.

**D-02 — the full stamp is clip-hidden, not `display: none`. DECIDE & STATE.**
A `title` tooltip is poorly announced by screen readers, so removing the date
from the accessibility tree would have made the board less readable, not more.
Clipping keeps it announced and keeps board.js's ticking node alive.

**D-03 — the tooltip is rebuilt on `pointerenter`/`focusin`, not on a timer.
DECIDE & STATE.** A `title` set once at boot would say "0s ago" forever. A second
one-second interval would duplicate the ticker board.js already owns. A tooltip
is only read while open, so refreshing it when the pointer arrives is both
cheaper and correct. Reversible; the test asserts only that something refreshes it.

**D-04 — the "Touched in" label was dropped when the chips moved. DECIDE & STATE.**
The REQ's target rendering is "…in the last 24 hours [6h][24h][48h][7d]": the
sentence already says what the chips window, so the label would repeat it two
words later. The group keeps `role="group"` and `aria-label="Activity window"`,
so its accessible name is unchanged. Reversible in one line.

**D-05 — the captured RED assertion "in the last 48 hours" was adapted to the
shipped phrase "in the last 2 days". DECIDE & STATE.** `activityWindowPhrase`
renders any whole-day window in days and always has. The REQ's Constraints
forbid changing window behavior, so asserting words the board does not use would
have pinned nothing. The test carries the reason inline. See DT-1 for the
underlying wording mismatch, which I did not fix.

**D-06 — the 40 px pad stays on `.activity-summary`; the row aligns on
baselines. DECIDE & STATE.** The brief allowed moving the pad to the new row
provided the probe still measured 40 px. Leaving it where REQ-585 put it is the
smaller change and makes the invariant structural rather than arithmetic: the
summary is the tallest item in the row, so nothing can push its text down.
Measured at exactly 40.0 px (T3).

**D-07 — `<span>` separators carrying a literal `·`, marked `aria-hidden`.
DECIDE & STATE.** CSS `::before` content is not reliably reachable and would put
punctuation in the stylesheet; a hidden span keeps the mark out of the
accessibility tree while leaving it selectable and greppable in the template.

## Discovered Tasks

- **DT-1** — The 48h chip's sentence reads "in the last 2 days" while the chip
  itself says "48h", because `activityWindowPhrase` (`web/board-activity.js`)
  spells any whole-day window in days. The 7d chip has the same shape and reads
  correctly; only 48h reads as a different unit from the control the reader just
  clicked. Pre-existing, untouched here. impact-user-visible → report only.
- **DT-2** — The board's page overflows horizontally by 26 px at a 500 px
  viewport (`documentScrollWidth` 526 against `documentClientWidth` 500),
  measured identically on the pre-change and post-change builds, so it is not
  caused by this REQ. impact-user-visible → report only.
- **DT-3** — Headless Chrome clamps its window to 500 px, so the frontend
  crew-member's 320 px check cannot be performed with `--window-size` alone; it
  needs a CDP device-metrics override. The existing browser probes have the same
  blind spot. impact-low → report only.

## Integration seams

The release paths are owned by finalization, so I wrote none of them. One entry
belongs in the root `CHANGELOG.md` (and the board rides the suite version per
`_dev/primes/prime-kanban-board.md`, which gives the tool no changelog of its
own). Suggested text, in the repo's descriptive-title house style:

```markdown
### Board top bar on one line, with the Activity window chips inside the view

The identity block is now a single line — `do-work/ queue · skill-do-work2 ·
13:42 UTC` — with the full generation stamp and its age in the line's tooltip,
and the Touched-in chips moved from the top bar onto the Activity view's summary
line. The bar keeps two control pills on every view and measures 68px at 1400px,
where it used to grow to 126px. The view buttons now read Board, Activity,
Calendar, Timeline, Durations, Testing.
```

Nothing else outside my write set needs a line: no test elsewhere in the package
referenced `#activity-window-group`, the identity classes, or the view order, and
the full package suite (T2, 397 tests) is green without edits to any other file.
