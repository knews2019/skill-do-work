# Builder brief — REQ-585

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1248/worktree-agent-REQ-585-one-scroll-surface`
- **Your branch (already checked out there):** `worktree-agent-REQ-585-one-scroll-surface`
- **Route:** A
- **Base commit:** 1f43608b (main)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work-board/tools/queue-kanban/VERSION` — release paths owned by finalization.
- Any file outside the write boundary below. If you need one, stop and report it to the orchestrator in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.
- Do not build or serve the board on port 8090: a live board owned by the user is running there. Use another port if you serve one at all.

## Write boundary (derived from the REQ; no `write_set` was declared)

- `skills/do-work-board/tools/queue-kanban/web/board.css` — the Activity-view rules (the `/* ---- activity view` block near line 3148) and, if you scope the padding move through a class, the `.board-main` rule near line 417.
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` — only if the Activity-only scoping needs a class toggled where the view changes (`viewState.view`), instead of a CSS `:has()` selector. Prefer whichever is smaller and clearest; say why in `## Decisions`.
- One new or extended Go probe test in `skills/do-work-board/tools/queue-kanban/` (`*_browser_probe_test.go` shape) pinning the GREEN condition. Read `browser_probe_test.go` lines 30–110 for the harness: probes run only with `QUEUE_KANBAN_BROWSER_PROBES=on`, find a browser on PATH or via `QUEUE_KANBAN_BROWSER=<path>`, render a page in a real engine and read back one JSON result node (`queue-kanban-probe-result`). Follow `timeline_browser_probe_test.go` (`TestBrowserBehaviorTimelineBarsCarryTheirStatusColour`) as the pattern. On this machine try `QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"`. If no engine runs, the test must still compile and skip cleanly, and you say so in the hand-back.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/frontend.md`

Also read every path in the REQ's `prime_files` (`_dev/primes/prime-kanban-board.md`), the shipped prime it points at (`skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`), and the `lessons-<name>.md` satellite beside each whose Read-first or Traps entries your change touches. The REQ's `## Required Lessons — Dropped for Budget` names both satellites as over budget; that is a record, not a prohibition — read the parts that touch `web/board.css` and the browser probe lane.

## What to build (the decided layout: M1)

The orchestrator resolved the REQ's open question as **D-01: M1, one page scroll with the column header stuck to the board's top edge**. The mockup of exactly this is `ai-reports/2026-09-05_1520_activity-view-double-scroll-mockups/mockups/m1-one-page-scroll.html` (serve the repo root over HTTP to open it); its `<style>` block is the whole proposed diff:

```css
.activity-table-scroll { max-height: none; overflow: visible; }
/* the board's 24px top padding is scrollable content, so rows would show
   through above the stuck header; move it onto the summary line instead */
.board-main       { padding-top: 0; }   /* Activity view only */
#view-activity    { padding-top: 0; }
.activity-summary { padding-top: 40px; }
```

Scope the padding move to the Activity view only (other views keep their 24 px top padding). The `thead th { position: sticky; top: 0 }` rule already exists; keep it. Update the block comment above `.activity-table-scroll` so it no longer claims the table "borrows the timeline table's metrics" for scrolling. Note that `REQ-578` already removed the verify-findings strip from this view (merged), so nothing but the summary line sits above the table.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the write boundary.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes), verify no debug artifacts (`console.log`, `debugger`, stray `TODO`) in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. Use, from the repo root of your worktree:
- Go (includes the assembly and CSS-shape tests): `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`
- Node lane: `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`
- Your browser probe: `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test ./... -run '<your test name>' -count=1` inside `skills/do-work-board/tools/queue-kanban`.

`tdd: false`, but the REQ's `## Red-Green Proof` is the acceptance target: measure it. RED today is `[true, true]` for `[.board-main scrolls, .activity-table-scroll scrolls]` on a populated Activity view; GREEN is `[true, false]` with the column header still visible after scrolling the board 700 px and no row showing above it. Record the actual numbers (clientHeight/scrollHeight of both, header top after scroll) in `## Test evidence`, from a real engine if one runs.

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-124800/REQ-585-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, and the RED/GREEN measurements above.
- `## Lesson evidence` — each lesson satellite or family you read, and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN` starting at **D-02** (D-01 is taken), each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings. Do not fix them inline. (The Timeline view has the same nested-scroll shape in `.timeline-scroll`; the orchestrator is capturing that separately as its own REQ, so do not touch it.)
- `## Integration seams` — any exact line that belongs in a file outside your write boundary, with where it goes. The orchestrator applies it.

---

# The request
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
claimed_at: 2026-09-05T12:46:41Z
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
