# Builder brief — REQ-586

## Where you work

- **Your worktree (cd here first):** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905-1248/worktree-agent-REQ-586-top-bar-one-line`
- **Your branch (already checked out there):** `worktree-agent-REQ-586-top-bar-one-line`
- **Route:** A
- **Base commit:** 59f169d0 (main; contains REQ-585's merge c08ac2b4, the Activity view's single-scroll change)

You are the builder. The orchestrator runs in the main checkout at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` and is the only writer of `do-work/`. Commit your work on your own branch in your own worktree and hand back a manifest; the orchestrator merges.

## Never touch

- Anything under `do-work/` — with exactly one exception, the hand-back file named below, which you write by its absolute main-tree path and never stage or commit.
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work-board/tools/queue-kanban/VERSION` — release paths owned by finalization.
- Any file outside the write set declared in the REQ below. If you need one, stop and report it to the orchestrator in the hand-back instead of writing it, unless the REQ's own requirements already demand that file class (then flag the contradiction and proceed).
- Do not run `bash _dev/tests/maintainer-verify.sh` (the repository gate). The orchestrator owns it and concurrent runs corrupt each other's timing budgets. Run only the focused tests named below.
- Do not build or serve the board on port 8090: a live board owned by the user is running there. Use another port if you serve one at all.

## Rules to load and follow (read these first, from your worktree)

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/shared-principles.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/frontend.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/testing.md` (the REQ is `tdd: true`)

Also read every path in the REQ's `prime_files` (`_dev/primes/prime-kanban-board.md`), the shipped prime it points at (`skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`), and the `lessons-<name>.md` satellite beside each whose Read-first or Traps entries your change touches. The REQ's `## Required Lessons — Dropped for Budget` names both satellites as over budget; that is a record, not a prohibition — read the parts that touch the web assets and the Node behaviour lane.

## What to build

Three changes, all in the REQ below (read its Detailed Requirements and the Addendum):

1. **One-line identity** in `web/template.html` lines 16–19: the wordmark, project and time on a single `nowrap` line, `do-work/ queue · skill-do-work2 · 12:17 UTC`, with the full "Generated 2026-09-05 12:17 UTC · 37s ago" text in a `title` tooltip on that line. Keep `id="board-project"` and `id="board-generated"`: `web/board.js` line 60 reads `board-generated` and the serve-mode refresh writes into it. Check what it writes before deciding which element carries the visible short time and which carries the full text.
2. **Touched-in chips move into the Activity view.** Delete `#activity-window-group` from the top bar and render the same four buttons (same `data-activity-window` values, `aria-pressed`, `setActiveButton` call in `board-activity.js`) beside the Activity summary line, so the line reads "236 transitions across 49 REQs in the last [6h] [24h] [48h] [7d]". `board-controls.js` line 46 (`document.getElementById("activity-window-group").hidden = viewState.view !== "activity"`) becomes unnecessary once the group lives inside `#view-activity`; delete it.
   Two existing readers constrain the markup: `javascript_behavior_c_test.go` line 2466 reads `#activity-summary`'s `textContent` as the sentence, and REQ-585's `activity_scroll_browser_probe_test.go` measures where `#activity-summary`'s text starts (40 px below the board's top edge). So keep `#activity-summary` as the sentence-only element and put the chips in a sibling inside one row container (a `<p>` cannot hold the pill `div`); if the 40 px top padding REQ-585 put on `.activity-summary` moves to that row container, the probe's measurement must still come out at 40 px. Run that probe.
3. **View button order** in `web/template.html` lines 71–88: Board, Activity, Calendar, Timeline, Durations, Testing. The order is declared only there; `board-controls.js` reads the buttons from the DOM.

REQ-585's Activity block in `web/board.css` (the `:has()` padding rule and the 40 px summary padding) is the base you build on; extend it, do not undo it.

## P-A-U phasing (mandatory, reported in the hand-back)

The REQ file is the orchestrator's, so report your P-A-U record under a `## P-A-U` heading in the hand-back instead of ticking boxes in the REQ:
- **[PLAN]** — brief technical approach, written before code.
- **[APPLY]** — code exactly as planned, strictly inside the declared write set.
- **[UNIFY]** — run `git diff --stat`, run the native linters (`gofmt -l .`, `go vet ./...` for Go changes), verify no debug artifacts (`console.log`, `debugger`, stray `TODO`) in added lines, and list each file you checked and what you checked.

## Focused tests

Every test-file invocation must finish in under 30 seconds. From the repo root of your worktree:
- Node lane (the RED/GREEN lane): `QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' ./...`
- Go: `bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`
- REQ-585's browser probe, which your markup change must keep green: `QUEUE_KANBAN_BROWSER_PROBES=on QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" go test ./... -run '^TestBrowserBehaviorActivityViewHasOneScrollSurface$' -count=1` inside `skills/do-work-board/tools/queue-kanban`.

`tdd: true`. The `## Red-Green Proof` in the REQ is the captured RED/GREEN pair: write the Node-lane test first (the 24h chip is a descendant of `#view-activity`, not of `.board-topbar`; clicking 48h updates the summary text and `aria-pressed`; the button order in the view pill), observe it fail, then make it pass. The one-line identity is a layout fact the Node lane cannot see: verify it in a real engine (extend the browser probe or measure it by hand) and record the numbers.

## Hand-back (write this file, then stop)

Write **`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-124800/REQ-586-handback.md`** using that absolute path — it is the one main-tree path you may write, and you must never stage or commit it.

It must contain, each under its own `##` heading:
- `## Branch` — the branch name and the head commit you left on it.
- `## File manifest` — every source file created/modified/deleted with the verb, plus tests touched.
- `## P-A-U` — the three phases above.
- `## Test evidence` — every command you ran, its exit status, the RED observation (test name + failure) and the GREEN observation, and the identity-line measurement.
- `## Lesson evidence` — each lesson satellite or family you read, and any listed path that was missing.
- `## Decisions` — significant choices as `D-NN` starting at **D-01**, each with reasoning. Mark a reversible low-reach choice DECIDE & STATE; mark an irreversible, taste-dependent or contestable one ESCALATE and add `Value:` and `Risk:` lines.
- `## Discovered Tasks` — out-of-scope findings. Do not fix them inline.
- `## Integration seams` — any exact line that belongs in a file outside your write set, with where it goes. The orchestrator applies it.

---

# The request
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
claimed_at: 2026-09-05T13:24:34Z
---

# Keep the Board Top Bar to One Line: Single-Line Identity and Touched-In Chips Inside the Activity View

## What

The top bar grows whenever its three control groups (filters, views, Touched-in) no longer fit beside the identity block: the controls wrap to two rows and the identity block (`do-work/ queue`, project, "Generated", date) wraps to four lines, so a 68 px bar becomes about 150 px. That space matters most on the Activity view, where the reader wants rows on screen to click a REQ and see every row of it highlighted (REQ-573). Two changes, both chosen by the user:

1. **O1, one-line identity.** Render the identity as one `nowrap` line, `do-work/ queue · skill-do-work2 · 12:17 UTC`, with the full "Generated 2026-09-05 12:17 UTC · 37s ago" text kept in a `title` tooltip on that line. The bar no longer grows when the controls beside it wrap.
2. **O2, Touched-in chips move into the Activity view.** Delete the `#activity-window-group` pill from the top bar and render the same four chips (6h, 24h, 48h, 7d) on the Activity view's summary line, so it reads "236 transitions across 49 REQs in the last [6h] [24h] [48h] [7d]". The top bar keeps two groups (filters, views) and fits on one row at far narrower widths. The Timeline already keeps its period controls inside its view, so this follows the existing pattern.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
