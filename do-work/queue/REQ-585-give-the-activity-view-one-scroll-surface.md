---
id: REQ-585
title: 'Give the Activity view one scroll surface instead of a scroll box inside the scrolling board'
status: pending
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
---

# Give the Activity View One Scroll Surface Instead of a Scroll Box Inside the Scrolling Board

## What

The Activity view scrolls twice. The transitions table sits in `.activity-table-scroll`, a box capped at `max-height: 70vh` with `overflow: auto`, and that box sits inside `.board-main`, which is the board's own scroll container by design. Two nested scroll regions give two scrollbars, and the mouse wheel moves whichever one the pointer is over. Leave exactly one scroll surface on this view.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

- [ ] Which layout: M1 (one page scroll, sticky header; recommended, fewest rules, keeps the board's one-scroll-container rule), M2 (table fills the viewport, board stops scrolling on this view), or M3 (M1 plus the pinned summary line). **Recommended:** M1. The builder proceeds with M1 if this is still open at claim time.

## Required Lessons — Dropped for Budget

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (5744 tokens, `slugged: partial`, so bare only): matches on "Changing queue-kanban UI or browser behavior". Over the 2000-token budget on its own.
- `_dev/primes/lessons-kanban-board.md` (4820 tokens, `slugged: partial`): matches on "Changing queue-kanban views". Over the budget on its own.

## Assets

- `do-work/user-requests/UR-120/assets/REQ-585-screenshot-1-activity-double-scroll.png`: the live board at 127.0.0.1:8090, Activity view, 24h window, generated 12:11 UTC on 2026-09-05. Above the table, the verify-findings strip with five worktree cards. Below it, "234 transitions across 49 REQs in the last 24 hours" and the table (REQ, Title, Status, What happened, When, Stamp). Two scrollbars are visible on the right: a thick one on the table box starting at the column header, and a thin one on the board area behind it starting under the top bar.

*Source: "<- this double scrolling behavior is not good, check REQs that would address it, if none make one"*
