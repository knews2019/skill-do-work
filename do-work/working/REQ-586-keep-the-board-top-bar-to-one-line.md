---
id: REQ-586
title: 'Keep the board top bar to one line: single-line identity and Touched-in chips inside the Activity view'
status: claimed
created_at: 2026-09-05T12:40:00Z
user_request: UR-121
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-585]
related: [REQ-585, REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board-activity.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T13:24:56Z
  basis:
    - trivial short-circuit
dispatch_at: 2026-09-05T13:26:18Z
builder_handback_at: 2026-09-05T13:46:18Z
integration_at: 2026-09-05T13:46:18Z
review_at: 2026-09-05T13:50:28Z
kb_status: pending
commit: 2ea0b1508b579df9fb96b8d27c2c71c7f8180017
claimed_at: 2026-09-05T13:24:34Z
---

# Keep the Board Top Bar to One Line: Single-Line Identity and Touched-In Chips Inside the Activity View

## What

The top bar grows whenever its three control groups (filters, views, Touched-in) no longer fit beside the identity block: the controls wrap to two rows and the identity block (`do-work/ queue`, project, "Generated", date) wraps to four lines, so a 68 px bar becomes about 150 px. That space matters most on the Activity view, where the reader wants rows on screen to click a REQ and see every row of it highlighted (REQ-573). Two changes, both chosen by the user:

1. **O1, one-line identity.** Render the identity as one `nowrap` line, `do-work/ queue · skill-do-work2 · 12:17 UTC`, with the full "Generated 2026-09-05 12:17 UTC · 37s ago" text kept in a `title` tooltip on that line. The bar no longer grows when the controls beside it wrap.
2. **O2, Touched-in chips move into the Activity view.** Delete the `#activity-window-group` pill from the top bar and render the same four chips (6h, 24h, 48h, 7d) on the Activity view's summary line, so it reads "236 transitions across 49 REQs in the last [6h] [24h] [48h] [7d]". The top bar keeps two groups (filters, views) and fits on one row at far narrower widths. The Timeline already keeps its period controls inside its view, so this follows the existing pattern.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reconciled from the builder hand-back: the generated-at placeholder is substituted once and board.js appends the ticking age into `#board-generated`, so that element keeps the full stamp (clipped) and a new `#board-generated-clock` shows the short time; the chips move into a `.activity-summary-row` beside the summary with the 40 px pad left on the summary; the six view buttons reorder in the template; two Node-lane tests written first.
- [x] **[APPLY]:** Reconciled from the builder hand-back: four files inside the declared write set (`template.html`, `board.css`, `board-controls.js`, `javascript_behavior_c_test.go`); `board-activity.js` needed no change because the chips keep their ids and their `setActiveButton` call.
- [x] **[UNIFY]:** Reconciled from the builder hand-back: `git diff --stat` shows four files, 420 insertions and 26 deletions; `gofmt -l .` empty and `go vet ./...` clean; grep of added lines for console.log, debugger, TODO, FIXME, XXX found nothing; each file read in full, REQ-585's `:has()` rule and 40 px pad byte-identical, placeholder still present exactly once, the chip group keeps `role`, `aria-label` and all four `data-activity-window` values.

## Detailed Requirements

- The chips keep their ids, `data-activity-window` values, `aria-pressed` handling, and the `setActiveButton` call in `board-activity.js`; only their home moves. `board-controls.js` line 46 (`document.getElementById("activity-window-group").hidden = viewState.view !== "activity"`) becomes unnecessary once the group lives inside `#view-activity`, which is hidden with the view; delete it rather than keeping a no-op.
- The summary line and the chips stay on one line at desktop widths; the chips sit after the count text, styled as the same pill group, not as a second bar.
- With REQ-585 (give the Activity view one scroll surface) landed, the summary line is the natural place for the chips whether or not it is pinned (M3); if the pinned variant was chosen, the chips are pinned with it.
- The one-line identity keeps the wordmark, project name and time in that order, separated by a middle dot, and keeps the existing `id="board-project"` and `id="board-generated"` hooks that `board.js` (line 60 reads `board-generated`) and the serve-mode refresh write into; check for readers before renaming anything.
- The `@media (max-width: 760px)` rule that stacks the top bar vertically stays as it is; this REQ is about the widths above it.

## Constraints

- No new control, no new state: the window values, the default of 24h, and the persistence behavior stay exactly as today.
- Board version and changelog per `_dev/primes/prime-kanban-board.md`.

## Dependencies

Depends on REQ-585 (give the Activity view one scroll surface): both edit the Activity summary line and the same block of `board.css`, so this builds after that lands rather than beside it.

## Red-Green Proof
**RED prompt/case:** In the Node behavior lane (`javascript_behavior_*_test.go`), load the board with an Activity payload, switch to the Activity view, and assert that the element carrying the `data-activity-window="24"` button is a descendant of `#view-activity` and not of `.board-topbar`; then click the `48h` button and assert the summary text says "in the last 48 hours" and the `48h` button has `aria-pressed="true"`.
**Why RED now:** `template.html` line 94 places `#activity-window-group` inside `.board-topbar`'s `.board-controls`; the descendant assertion fails.
**GREEN when:** The assertion passes, the top bar contains exactly two `.control-group` pills on every view, and on the served board at a width where the old bar wrapped (about 1400 px), the top bar is one row of 68 px with the identity on a single line and the chips visible on the Activity summary line. The one-line identity is a layout fact the Node lane cannot see; verify it with a browser probe or a captured screenshot in the hand-back.
**Validation:** User confirmed (chose O1 and O2 from three options, 2026-09-05)

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

## Assets

- `do-work/user-requests/UR-121/assets/REQ-586-screenshot-1-top-bar-wrapped.png`: the top bar of the M1 mockup (`ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/mockups/m1-one-page-scroll.html`) in a narrow frame. Left, the identity block on four lines: "do-work/", "queue", "skill-do-work2", "Generated 2026-09-05" and "12:17 UTC" on a fifth. Right, the controls on two rows: the filter pill (Filter id or title, All domains, All statuses) above, the view pill (Board, Calendar, Durations, Timeline, Activity selected, Testing) and the Touched-in pill (6h, 24h selected, 48h, 7d) below. The mockup copies the shipped top bar rules, so the real board wraps the same way at the width where its three groups stop fitting on one row.

## Full Context
See `do-work/user-requests/UR-121/input.md` for complete verbatim input.

*Source: "this part of the header is still taking up too much vertical space, and that is precious when I want to click a req and I want to highlight all of it's occurances" / "ok, do o1 and o2 capture it first"*

## Addendum (2026-09-05)

User added:

> ````text
> while we are at it the order should be Board, Activity, Calendar, Timeline, Durations
> 
> testing can remain last
> ````

- Reorder the view buttons in the top bar to: Board, Activity, Calendar, Timeline, Durations, Testing. The order is declared once, in `template.html` lines 71 to 88 (the `data-view-target` buttons); `board-controls.js` reads the buttons from the DOM, so nothing else keeps the list.
- Asset: `do-work/user-requests/UR-122/assets/REQ-586-screenshot-2-view-tab-order.png`, the view pill today: Board, Calendar, Durations, Timeline (selected, with focus ring), Activity, Testing.
- Same proof lane as the rest of this REQ: the behavior test can assert the `data-view-target` values of the pill's buttons in document order.

---

## Triage

**Route: A** - Simple

**Reasoning:** Three named changes in five named files with a runnable RED in the existing Node behaviour lane: a one-line identity block, the Touched-in chips moving into the Activity summary line, and a fixed button order. No discovery needed; the request names the lines.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified)
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modified)

**What was done:** The top bar's identity block is one nowrap line, `do-work/ queue · skill-do-work2 · 13:42 UTC`, with the full "Generated … · Ns ago" stamp clip-hidden in `#board-generated` (still written by board.js) and surfaced as the line's tooltip on hover or focus; the Touched-in chips left the top bar for a new `.activity-summary-row` beside the Activity summary line, keeping their ids, values and `aria-pressed` handling, and the hand toggle in `applyView` was deleted; the six view buttons now read Board, Activity, Calendar, Timeline, Durations, Testing. Two Node-lane tests pin the chip ancestry, the deleted toggle, the view order, the shipped window writer and the identity renderer. Measured in headless Chrome 152 at 1400 px: the top bar is 68 px where it was 126 px.

**Implementation range:** `0edce090..2ea0b150`. Builder commit `cc56324f`.

## Qualification

Passed the request-bound advance qualify gate for the range `0edce090..2ea0b150` (`state: satisfied`, no findings). Four files, all inside the declared five-file write set; `board-activity.js` was declared and left untouched because the chips keep their ids and their `setActiveButton` call. Requirements traced against the diff: the identity is one nowrap line with the full stamp in the tooltip and both `board-project` and `board-generated` ids kept (board.js still writes into the latter); `#activity-window-group` is gone from the top bar and sits in a new row beside `#activity-summary` inside `#view-activity`, with the same values, `aria-pressed` handling and group label; the `applyView` hand toggle is deleted; the view buttons read Board, Activity, Calendar, Timeline, Durations, Testing; the 760 px media rule is untouched; REQ-585's `:has()` rule and 40 px pad are byte-identical. The P-A-U boxes were reconciled from the builder hand-back. No debug artifacts in the added lines.

## Testing

**Tests run:** `do-work/runs/work-2026-09-05-124800/REQ-586-probe.sh` through the advance test gate (`run-blocked-check`, raw status 0): the two new Node-lane tests, the other Activity behaviour tests, REQ-585's browser probe, and the generate and assembly tests, run in the detached worktree at the merge revision `2ea0b150` because the shared main tree carried a sibling session's uncommitted edits to the same package.
**Result:** ✓ All passing. Builder lanes at the branch head `cc56324f`: strict Node lane 61 tests, wall 7 s; Go lane 397 tests, wall 46 s, slowest file 9.25 s; REQ-585's browser probe green in 1.32 s (summary text still at 40.0 px).

**Repository gate:** `bash _dev/tests/maintainer-verify.sh` from a detached worktree at `2ea0b150`, first run exited 0 (13:46 to 13:49 UTC, load 7; slowest do-work-cli test file 24.97 s). Recorded through the advance green gate (`satisfied`).

**Red-green validation:** traced to `## Red-Green Proof`, run before any implementation edit and failing on assertions, not compilation:
- `TestJavaScriptBehaviorActivityWindowChipsRenderInsideTheActivityView`: ✗ before ("the top bar still carries the Activity window chips"; the four chips missing from the Activity view; "board-controls.js still toggles #activity-window-group by hand"; view pill read board calendar durations timeline activity testing) → ✓ after, with the 48h and 6h clicks driving the shipped `applyActivityWindowSelection` and all four `aria-pressed` values read back.
- `TestJavaScriptBehaviorTopBarIdentityIsOneLineWithAFullStampTooltip`: ✗ before (identity block missing its parts) → ✓ after.
- The captured phrase "in the last 48 hours" was adapted to the board's own "in the last 2 days": `activityWindowPhrase` has always spelled whole-day windows in days and this REQ's constraints forbid changing window behaviour (D-05, DT-1).
- The one-line identity is layout, so it was measured by hand in headless Chrome 152.0.7977.76 at 1400×900 on two static boards generated from this repository: top bar 126.0 px → 68.0 px, identity block 85.0 px → 27.0 px, visible pills three → two, chips inside `#view-activity` false → true, view order as requested, identity `title` "Generated 2026-09-05 13:42 UTC · 14s ago".

**New tests added:**
- `TestJavaScriptBehaviorActivityWindowChipsRenderInsideTheActivityView` and `TestJavaScriptBehaviorTopBarIdentityIsOneLineWithAFullStampTooltip` in `javascript_behavior_c_test.go`, with helpers `viewTargetsInDocumentOrder` and `sliceMarkupElementAfter`

**Heavy verification plan:** planned with `plan-heavy-verification` over `0edce090..2ea0b150`, uncovered paths none.
- Range: `0edce090`..`2ea0b150`
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — all four changed paths match subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree match
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four paths match subtree `skills`

*Verified by work action*

## Decisions

Copied from the builder hand-back:

- **D-01 — `#board-generated` carries the full stamp; a new `#board-generated-clock` carries the visible time (DECIDE & STATE):** the placeholder is substituted once and `generate_test.go` fails on a leftover, and board.js (outside the write set) appends the ticking age into `#board-generated`, so that element cannot double as the short clock.
- **D-02 — the full stamp is clip-hidden, not `display: none` (DECIDE & STATE):** a `title` tooltip is poorly announced, so the date stays in the accessibility tree and board.js's ticking node stays alive.
- **D-03 — the tooltip is rebuilt on `pointerenter` and `focusin`, not on a timer (DECIDE & STATE):** a title set once would say "0s ago" forever and a second interval would duplicate board.js's ticker.
- **D-04 — the "Touched in" label was dropped when the chips moved (DECIDE & STATE):** the sentence already says what the chips window; the group keeps `role="group"` and `aria-label="Activity window"`.
- **D-05 — the captured RED phrase "in the last 48 hours" was adapted to the shipped "in the last 2 days" (DECIDE & STATE):** `activityWindowPhrase` spells whole-day windows in days and the REQ's constraints forbid changing window behaviour; the test carries the reason inline. The wording mismatch is DT-1 below.
- **D-06 — the 40 px pad stays on `.activity-summary`; the row aligns on baselines (DECIDE & STATE):** the summary is the tallest item in the row, so nothing can push its text down; REQ-585's probe measured 40.0 px.
- **D-07 — `<span>` separators carrying a literal `·`, marked `aria-hidden` (DECIDE & STATE):** keeps the mark out of the accessibility tree and greppable in the template.

## Discovered Tasks

- **DT-1** — the 48h chip's sentence reads "in the last 2 days" while the chip says "48h", because `activityWindowPhrase` spells any whole-day window in days; pre-existing → report only, `impact-user-visible`
- **DT-2** — the page overflows horizontally by 26 px at a 500 px viewport, identically before and after this change → report only, `impact-user-visible`
- **DT-3** — headless Chrome clamps its window to 500 px, so the frontend crew's 320 px check needs a device-metrics override; the existing browser probes share the blind spot → report only, `impact-negligible`

## Review

**Overall: 93%** | 2026-09-05T13:50:28Z

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Reviewed in orchestrated mode against the merge range `0edce090..2ea0b150` (four files), the REQ with its addendum, and UR-121 and UR-122's verbatim input. Requirements: the identity is one nowrap line (`do-work/ queue · skill-do-work2 · 13:42 UTC`) with the full stamp clip-hidden and surfaced as the line's tooltip, both ids kept and board.js still writing the ticking age into `#board-generated`; the chips moved into a row beside the Activity summary with values, `aria-pressed`, group role and label intact, and the hand toggle in `applyView` is gone; the view pill reads Board, Activity, Calendar, Timeline, Durations, Testing; the 760 px media rule is untouched; REQ-585's rules are byte-identical. The one requirement not met to the letter is the summary text after clicking 48h: the REQ said "in the last 48 hours" and the board has always said "in the last 2 days" (D-05, DT-1); the builder was right not to change window phrasing under a constraint that forbids window behaviour changes. Code: the tooltip is rebuilt on hover and focus from `#board-generated`'s live child text, so it stays correct without a second ticker; the clip-hide keeps the date in the accessibility tree; the placeholder still appears once. Tests: two Node-lane tests written first and observed failing for the right reasons, both driving the shipped writers rather than assigning state by hand; REQ-585's browser probe still measures the summary text at 40.0 px; the identity line was measured in headless Chrome before and after (top bar 126 px to 68 px at 1400 px). Scope: four of the five declared files, nothing outside.

Restatement sweep: the template comment that explained why the chips lived in the top bar moved with them; no other text restates the identity layout or the button order. The board guide (`skills/do-work-board/docs/board-guide.md`) was checked for a view-order or "Touched in" description: none.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- The 48h chip's sentence says "in the last 2 days" while the chip says "48h" (`activityWindowPhrase` spells whole-day windows in days); pre-existing and outside this REQ's constraints — impact-user-visible → report only
- The "Touched in" label was dropped with the move (D-04); the sentence carries the meaning and the group keeps its accessible name, but a reader used to the label loses a cue — impact-negligible → report only
- `renderTopBarIdentity` is called once from `wireControls` and returns silently when any of its three nodes is missing, so a template rename would drop the clock without a test failure; the new identity test pins the ids today — impact-negligible → report only

**Acceptance:** Pass — both RED tests went green, the whole Node lane and the whole package are green, REQ-585's probe still reads 40.0 px, and the before/after render measured the bar at 68 px with the chips inside the Activity view and the buttons in the requested order.
**Suggested testing:** 2 items — hover the identity line on the served board and confirm the tooltip shows the full stamp with a live age; check the top bar at a width where the two remaining pills wrap (roughly 1000 px) to confirm the identity stays one line.
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Checking what already writes into an element before moving it. The `GENERATED_AT_DISPLAY` placeholder is substituted exactly once and board.js appends into `#board-generated`, so the visible clock had to be a new element and the full stamp had to stay; finding that before coding avoided a second ticker and a failing generate test.
**What didn't:** The captured RED phrase "in the last 48 hours" did not match the board's own wording for whole-day windows. The test now asserts the words the board uses and says why; the mismatch is a report-only task, not a silent fix.
**Worth knowing:** Two readers constrain the Activity summary markup: the Node-lane test reads `#activity-summary`'s textContent as the sentence, and REQ-585's browser probe measures where that text starts. Anything placed beside the sentence has to be a sibling, aligned on the baseline, never a child.

## Orientation

Now the top bar stays one row: the identity reads `do-work/ queue · project · HH:MM UTC` with the full stamp as its tooltip, the Touched-in chips live on the Activity view's summary line, and the view buttons read Board, Activity, Calendar, Timeline, Durations, Testing. Lives in the queue-kanban board's template, controls script and stylesheet, pinned by two Node-lane tests. No map change.

## Heavy Verification Plan

Held at Step 7.7 after a passing review. Base `0edce090`, target `2ea0b150` (the landed merge, recorded in `commit:`). Planned with `plan-heavy-verification --manifest _dev/tests/heavy-lanes.json`; uncovered paths none.

- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `template.html`, `board.css`, `board-controls.js` and `javascript_behavior_c_test.go` match subtree `skills/do-work-board/tools/queue-kanban`
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree match
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all four paths match subtree `skills`
