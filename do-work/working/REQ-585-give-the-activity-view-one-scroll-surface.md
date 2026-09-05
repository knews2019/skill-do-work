---
id: REQ-585
title: 'Give the Activity view one scroll surface instead of a scroll box inside the scrolling board'
status: claimed
created_at: 2026-09-05T12:26:00Z
user_request: UR-120
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
related: [REQ-578, REQ-573]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T12:48:56Z
  basis:
    - trivial short-circuit
dispatch_at: 2026-09-05T12:51:06Z
builder_handback_at: 2026-09-05T13:06:44Z
integration_at: 2026-09-05T13:07:15Z
review_at: 2026-09-05T13:24:00Z
kb_status: pending
commit: c08ac2b4b4b60cd9b4713078771e7b46bc10dd73
claimed_at: 2026-09-05T12:46:41Z
---

# Give the Activity View One Scroll Surface Instead of a Scroll Box Inside the Scrolling Board

## What

The Activity view scrolls twice. The transitions table sits in `.activity-table-scroll`, a box capped at `max-height: 70vh` with `overflow: auto`, and that box sits inside `.board-main`, which is the board's own scroll container by design. Two nested scroll regions give two scrollbars, and the mouse wheel moves whichever one the pointer is over. Leave exactly one scroll surface on this view.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reconciled from the builder hand-back: probe first against unmodified CSS to measure RED, then apply the M1 diff (drop the table's max-height and overflow, zero the board's and view's top padding for this view only, 40 px on the summary line), scope with `:has()` in CSS, re-run the probe for GREEN plus the Go, Node and browser lanes.
- [x] **[APPLY]:** Reconciled from the builder hand-back: `web/board.css` activity block and one new probe test; `board-controls.js` left untouched because the `:has()` route needed no toggle (D-03).
- [x] **[UNIFY]:** Reconciled from the builder hand-back: `git diff --stat` shows two files (29+/8- in board.css, 265 new probe lines); `gofmt -l .` empty and `go vet ./...` exit 0; grep of added lines for console.log, debugger, TODO, FIXME, XXX found nothing; both files read in full, sticky header rule untouched, new selector specificity checked against the 760 px media query.

## Red-Green Proof
**RED prompt/case:** Serve the board (`queue-kanban serve`), open the Activity view with the 24h window on a queue with more transitions than fit on one screen, and run in the console: `const m=document.querySelector('.board-main'), t=document.querySelector('.activity-table-scroll'); [m.scrollHeight>m.clientHeight, t.scrollHeight>t.clientHeight]`.
**Why RED now:** It returns `[true, true]` (measured 2026-09-05 12:15 UTC at 2466×1323 CSS px: board area 1231 px tall holding 1442 px, table box 926 px tall holding 6815 px) and two scrollbars are visible in the screenshot under Assets.
**GREEN when:** Exactly one of the two is `true`. With the recommended layout (M1 below) `.board-main` scrolls and the table reports `scrollHeight === clientHeight`; the column header is still visible after scrolling the board 700 px, and no row shows through above it. A browser probe in the board's existing probe lane (`*_browser_probe_test.go`) pins this, since the fact is layout and the Node behavior lane cannot see it.
**Validation:** Inferred during capture

## Context

Three mockups, each a scrollable page built from the board's stylesheet and the real Activity rows, plus the cause and a side-by-side table, are in `ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/index.html` (serve the repo over HTTP to open it; the iframes do not load from `file://`). The CSS under each mockup is the whole change for that variant.

- M1, one page scroll: `.activity-table-scroll { max-height: none; overflow: visible; }`; the `thead` keeps its sticky rule and now sticks to the top of `.board-main`. The board's 24 px top padding is scrollable content, so rows would show through above the stuck header; move that padding onto the summary line for this view (scoped, so other views keep it).
- M2, table fills the viewport: `.board-main` stops scrolling on this view and the table takes the remaining height as the only scroller. Gives Activity a scroll model no other view has.
- M3, M1 plus the summary line pinned above the header at a fixed height.

REQ-578 (hide the verify-findings strip on the Activity view) has merged, so nothing but the summary line sits above the table any more; that is why M1 needs no other change. REQ-573 (open the detail drawer from an Activity row) touches the same rows and stylesheet section; neither depends on the other.

## Open Questions

- [~] Which layout: M1 (one page scroll, sticky header; recommended, fewest rules, keeps the board's one-scroll-container rule), M2 (table fills the viewport, board stops scrolling on this view), or M3 (M1 plus the pinned summary line). **Recommended:** M1. The builder proceeds with M1 if this is still open at claim time. → **D-01**: Builder chose: M1 (one page scroll, sticky header). Reasoning: the recommended answer at capture, the smallest diff (two declarations removed plus one scoped padding move), and the only variant that keeps the board's one-scroll-container rule intact. Value: the double scroll is gone with no new scroll model for readers to learn, and REQ-586 (top bar on one line, chips into the Activity view) builds on the same summary line either way. Risk: if the user wanted the summary count pinned (M3), that is a follow-up of a few CSS lines, fully reversible.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

## Assets

- `do-work/user-requests/UR-120/assets/REQ-585-screenshot-1-activity-double-scroll.png`: the live board at 127.0.0.1:8090, Activity view, 24h window, generated 12:11 UTC on 2026-09-05. Above the table, the verify-findings strip with five worktree cards. Below it, "234 transitions across 49 REQs in the last 24 hours" and the table (REQ, Title, Status, What happened, When, Stamp). Two scrollbars are visible on the right: a thick one on the table box starting at the column header, and a thin one on the board area behind it starting under the top bar.

*Source: "<- this double scrolling behavior is not good, check REQs that would address it, if none make one"*

---

## Triage

**Route: A** - Simple

**Reasoning:** A styling change in one named stylesheet section with a measured RED and a mockup of the target layout; the request names the file, the rules, and the proof.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go` (new)

**What was done:** The Activity table's own scroll box is gone: `.activity-table-scroll` lost its `max-height: 70vh` and `overflow: auto`, so the board area is the only scroll surface and the existing sticky column header now pins to the board's top edge. The board's 24 px top padding is zeroed on the Activity view only, through `.board-main:has(> #view-activity.is-active)`, and the summary line carries 40 px instead, so the reader sees the same layout at rest and no row paints above the stuck header while scrolling. A browser probe renders a synthetic 120-row Activity view in headless Chrome and asserts that exactly one element scrolls, that the header sits at 0 px after a 700 px scroll, and that the padding rule does not leak to the Kanban view.

**Implementation range:** `db36c8ca..c08ac2b4`. Builder commit `94532c45`.

## Qualification

Passed the request-bound advance qualify gate for the range `db36c8ca..c08ac2b4` (`state: satisfied`). Two files: the Activity block of `web/board.css` and one new browser probe test. The one warning, QUALIFY-NEW-FILE-UNWIRED on `activity_scroll_browser_probe_test.go`, is the expected shape of a Go test file: nothing references a `_test.go` file statically, the Go toolchain discovers it, and the hand-back's T3 and T4 runs show it executing in the browser lane. Requirements traced against the diff: the inner scroll box is gone (the two declarations deleted), exactly one scroll surface remains (measured `[true, false]`), the sticky header rule is untouched and pins at 0 px after a 700 px scroll, and the padding move is scoped to the Activity view (`:has()` rule, Kanban view measured back at 24 px). The P-A-U boxes were reconciled from the builder hand-back. No debug artifacts in the added lines.

## Testing

**Tests run:** `do-work/runs/work-2026-09-05-124800/REQ-585-probe.sh` through the advance test gate (`run-blocked-check`, raw status 0): the new browser probe plus `TestGenerateInlines*`, `TestBoardJavaScriptAssemblyStructure` and `TestJavaScriptBehaviorActivity*`, run in the detached worktree at the merge revision `c08ac2b4` because the shared main tree carried a sibling session's uncommitted edits to the same package.
**Result:** ✓ All passing. Builder lanes at the branch head `94532c45`: Go lane 396 tests, wall 45 s, slowest file 9.66 s; strict Node lane 59 tests; whole browser lane `^TestBrowserBehavior` green in 103 s on Chrome 152.0.7977.76.

**Repository gate:** `bash _dev/tests/maintainer-verify.sh` from a detached worktree at `c08ac2b4`. First run exited 1 at 13:10 UTC on the per-file budget alone: `internal/finalization/finalization_recovery_test.go` 32.85 s and `finalization_req499_test.go` 30.67 s in the do-work-cli module, under a load average above 20 from the sibling session's heavy lanes; nothing in the queue-kanban module or this REQ's files was involved. **Repository gate retry:** first run exited 1, rerun exited 0 (13:18 to 13:21 UTC, load 12; the same two files at 28.61 s and under). Recorded through the advance green gate (`satisfied`). HEAD moved to `ff796db5` during the gate through the sibling's REQ-578 release and REQ-582 merge; those commits touch do-work-cli sources and shipped prose, not the board package this REQ changed.

**Red-green validation:** traced to `## Red-Green Proof`, measured by the probe in a real engine (1600×900, 120 transition rows):
- `TestBrowserBehaviorActivityViewHasOneScrollSurface`: ✗ before implementation (`[board scrolls, table scrolls]` = `[true, true]`, table 530/3509 px, header 55.5 px below the board's top edge after scrolling) → ✓ after (`[true, false]`, table 3509/3509 px, header at 0.0 px after a 700 px scroll, 0 rows above it, summary text at 40.0 px before and after, board padding 0 px on Activity and 24 px back on the Kanban view).

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/activity_scroll_browser_probe_test.go` (browser lane; skips cleanly without an engine)

**Heavy verification plan:** planned with `plan-heavy-verification` over `db36c8ca..c08ac2b4`, uncovered paths none.
- Range: `db36c8ca`..`c08ac2b4`
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — both changed paths match subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree match; this lane is the one that executes the new probe
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both paths match subtree `skills`

*Verified by work action*

## Decisions

D-01 (M1, one page scroll) is recorded under Open Questions. The builder's decisions, copied from its hand-back:

- **D-02 — Verified the mockup's padding claim instead of taking it on faith (DECIDE & STATE):** a third CSS variant with only the scroll box removed measured the header pinned 24 px down with two rows visible above it, so Chrome pins a sticky header to the container's content box, not its padding box. The padding move stays and the stylesheet comment records the measurement with the browser build (Chrome 152.0.7977.76).
- **D-03 — Scoped the padding move with `:has()` in CSS rather than a class toggled in `board-controls.js` (DECIDE & STATE):** one rule in the file that owns the layout, no second place deciding "this is the Activity view"; `board.css` already ships a `:has()` rule and the browser lane targets current Chromium. `board-controls.js` is unmodified.
- **D-04 — One `40px` summary padding instead of tracking the board's padding per breakpoint (DECIDE & STATE):** below the 760 px breakpoint the gap is 6 px taller than before; exact tracking would cost three rules for six pixels on a phone-width board. Recorded in the stylesheet comment; one declaration in the existing media query undoes it.
- **D-05 — Kept the class name `.activity-table-scroll` (DECIDE & STATE):** renaming touches `web/template.html`, outside the write boundary; a comment marks the name as historical.
- **D-06 — Built the probe on a synthetic 40-REQ fixture rather than the live queue (DECIDE & STATE):** the live queue's last 24 hours might not overflow anything on a quiet day and the probe would pass while measuring nothing; the fixture yields exactly 120 rows and the probe asserts the count before comparing heights.

## Discovered Tasks

- The Timeline view has the same nested-scroll shape in `.timeline-scroll`. Captured separately as REQ-587 (give the Timeline view one scroll surface) before this build finished → report only, `impact-user-visible`
- `.activity-table-scroll` is now named for a scroll box that does not scroll; a rename touches `web/template.html`, `web/board.css` and the probe. Cosmetic → report only, `impact-negligible`
- The Activity view lays out every transition row in the window at full height instead of inside a 70vh clip. Row count is unchanged and 120 rows measured fine; worth a look only if the 7-day window on a busy queue feels slow → report only, `impact-negligible`

## Review

**Overall: 94%** | 2026-09-05T13:24:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

Reviewed in orchestrated mode against the merge range `db36c8ca..c08ac2b4` (two files), the REQ, and UR-120's verbatim input. Requirements: the inner scroll box is gone (both declarations deleted), exactly one scroll surface remains (measured `[true, false]` in Chrome 152), the sticky header pins at 0 px after a 700 px scroll with no row painting above it, and the padding move is scoped to the Activity view (`:has()` rule; Kanban view measured back at 24 px). Code: four declarations plus one new rule inside the Activity block, the block comment now states the measured sticky behaviour with the browser build beside it, and the `:has()` selector's specificity beats the 760 px media query shorthand. Tests: the new probe follows the existing browser-lane harness, refuses to compare heights until the fixture's 120 rows are laid out (REQ-291's lesson), reads its two constants through Go rather than restating them, and asserts the view-scoping in both directions. Scope: `board-controls.js` was in the boundary and left alone because the CSS route made it unnecessary; the only debt is the class name `.activity-table-scroll`, which now names a box that does not scroll.

Restatement sweep: the change redefines nothing other text restates. The mockup report under `ai-reports/` describes the same M1 diff and is a proposal record, not a contract.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- `.activity-table-scroll` keeps a name that describes the old scroll box; renaming needs `web/template.html`, `web/board.css` and the probe in one edit — impact-negligible → report only
- Below the 760 px breakpoint the summary line carries 40 px instead of the 34 px it had, a 6 px difference the builder chose over a per-breakpoint rule (D-04) — impact-negligible → report only
- `:has()` is now relied on for layout, not only decoration; the strict browser lane targets current Chromium and every current engine supports it, so this is a note, not a defect — impact-negligible → report only

**Acceptance:** Pass — the REQ's RED/GREEN pair was measured in a real engine before and after the change (`[true, true]` → `[true, false]`, header offset 55.5 px → 0.0 px), the whole browser lane passed, and the rebuilt board served on port 8090 shows one scroll surface on the Activity view.
**Suggested testing:** 2 items — check the Activity view once REQ-573's row-click highlight merges on top (same block of `board.css`, different rules), and glance at a 7-day window on a busy queue for layout cost now that every row is laid out at full height.
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Measuring the mockup's claim in a third CSS variant before trusting it (D-02). The mockup and the REQ said the board's top padding would let rows show through above the stuck header; the builder's reading of the specification said the opposite, and the measurement settled it in one probe run.
**What didn't:** Reasoning from the specification about where a sticky element pins. Chrome pins a sticky header inside a scroll container to the container's content box, not its padding box, so any padding on the scroll container is a band the rows scroll through above the header.
**Worth knowing:** When a view removes an inner scroll box in favour of the board's own scroll, the board's top padding has to move onto the first element of that view, scoped to that view. `.board-main:has(> #view-activity.is-active)` is the pattern; REQ-587 (the Timeline's same fix) should reuse it rather than invent a second scoping mechanism.

## Orientation

Now the Activity view scrolls as one surface: the board area scrolls, the transition rows are ordinary content, and the column header pins to the top edge while rows pass under it. Lives in the queue-kanban board's stylesheet (`web/board.css`, Activity block), with a browser-lane probe pinning the one-scroll-surface condition. No map change.

## Heavy Verification Plan

Held at Step 7.7 after a passing review. Base `db36c8ca0aa3da54edc149f73ce678e8c42c209e`, target `c08ac2b4b4b60cd9b4713078771e7b46bc10dd73` (the landed merge, recorded in `commit:`). Planned with `plan-heavy-verification --manifest _dev/tests/heavy-lanes.json`; uncovered paths none.

- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `activity_scroll_browser_probe_test.go` and `web/board.css` match subtree `skills/do-work-board/tools/queue-kanban`
- `queue-kanban-browser`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same subtree match
- `staged-skills`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — both paths match subtree `skills`
