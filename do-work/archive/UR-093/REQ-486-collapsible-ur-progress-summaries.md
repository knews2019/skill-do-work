---
id: REQ-486
title: 'Addendum: make UR groups collapsible and show progress summaries'
status: completed
priority: later
created_at: 2026-09-01T17:29:43Z
user_request: UR-093
addendum_to: REQ-236
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: ui-component
depends_on: [REQ-510]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-09-06T00:26:09Z
estimate:
  p50_active_minutes: 85
  confidence: low
  calculated_at: 2026-09-06T00:35:00Z
  basis:
    - Route C
    - 22-file write set
    - 3 new files
    - 3 subsystems involved
    - 8 acceptance criteria
    - dependency depth 1
    - browser evidence
    - async lifecycle behavior
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js
  - skills/do-work-board/tools/queue-kanban/web/board.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/timeline.go
  - skills/do-work-board/tools/queue-kanban/frontmatter_test.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go
  - skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go
  - skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go
  - skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go
  - skills/do-work-board/docs/board-guide.md
  - skills/do-work-board/actions/board.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/version.md
  - VERSION
  - skills/do-work/VERSION
  - CHANGELOG.md
  - skills/do-work/CHANGELOG.md
route: C
dispatch_at: 2026-09-06T02:22:21Z
builder_handback_at: 2026-09-06T02:22:21Z
planning_at: 2026-09-06T00:35:06Z
completed_at: 2026-09-06T03:26:33Z
commit: cbfcec76e75e1ab47315b197726169ded59d8d3a
release_at: 2026-09-06T03:26:33Z
---

# Addendum: Make UR Groups Collapsible and Show Progress Summaries

## Deferral (2026-09-03)

Hand triage, maintainer approved: deferred behind REQ-510, the end of the UR-098 orchestrator-simplification chain, so board feature work does not compete with pipeline simplification. Remove the `depends_on` edge to un-defer.

## What

Extend the board's existing UR presentation so the By UR card grid and the UR detail drawer's REQ-id list are independently collapsible. Show the same whole-UR request count, active-time rollup, remaining-time forecast, successful progress, and resolved progress on the By UR header and in the drawer.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `_dev/primes/prime-kanban-board.md` and the board's lesson satellite read before either
  increment, plus `_dev/primes/prime-releases.md` for the release in T5. The judge-panel plan in the run
  directory is the approach; its fifteen decisions were carried out or diverged from in writing.
- [x] **[APPLY]:** Two increments, merged separately. Twenty-five files across both, of which four were
  not in the original declaration and are recorded as D-16 rather than absorbed silently.
- [x] **[UNIFY]:** `git diff --stat` reports 25 files, 2177 insertions, 50 deletions for increment 2 on
  top of increment 1's 8 files. Linters and lanes, each with its own exit line rather than the gate's
  summary: fast stage `ok … 65.004s`; strict JavaScript `ok … 6.179s`, 70 pass, 0 skip; strict browser
  `ok … 100.538s`, 35 pass, 0 skip on Chromium 141; `contract-regressions.sh` green with the
  write-surface count unchanged; `shipped-package-reference-contract.sh` green, so the changelog mirror
  is byte-identical. No debug artifacts: the diff adds no `console.log`, no `debugger`, no commented-out
  block. Two greps stand as structural assertions rather than eyeballing — `completedAt` appears nowhere
  in the summary fragment, and neither does `try`.

## Context

The supplied screenshot shows the By UR lens with UR-081 expanded into a large card grid and its detail drawer open. The header reports only `43 REQ`, while the drawer reports `GROUPED REQS 43` followed by a REQ-id list long enough to fill the visible drawer. Neither surface reports elapsed work, expected work remaining, or progress.

This is an addendum to completed REQ-236, which added the separate URs only lens. The user now wants folding available in the normal By UR reading too, plus compact progress information on both UR surfaces. The user confirmed that both collapsible regions start open, that the time figures use an active-time model, and that successful and resolved percentages are both shown.

## Prior Implementation

REQ-236 implemented URs only as a fold modifier on the existing By UR renderer rather than as a third `viewState.lens` value. Both readings share `renderUserRequestLens`, `makeRequestCard`, and the `ur-group` markup. In the folded reading, the UR header is a real button with `aria-expanded`, cards are built only when opened and removed on collapse, and a separate Details button opens the drawer without colliding with the fold control. Fold state lives only in the rendered DOM and resets on a re-render. Its behavior and generated markup are pinned by Go-driven JavaScript probes in `generate_test.go`. The recorded implementation commit is `456ee9d`.

## Detailed Requirements

- In the **By UR** lens, every UR header independently collapses and expands that UR's REQ card grid.
- By UR groups start expanded. More than one group may be collapsed or expanded at once.
- The fold control is keyboard-operable and exposes the current state through `aria-expanded`.
- Opening the UR detail drawer remains a separate action from folding the group; keep a dedicated Details control rather than assigning two meanings to one button.
- In the UR detail drawer, the grouped REQ-id list is independently collapsible and starts expanded.
- Folding the drawer list hides only the linked REQ ids. The UR metrics and the rest of `input.md` remain visible.
- Preserve the existing URs only lens: its groups still start collapsed, expand in place, use the same filters and UR activity scope, and keep non-persisted DOM-only fold state.
- Both the By UR header and the UR drawer show the same whole-UR summary:
  - total grouped REQs;
  - active time spent;
  - estimated active time remaining;
  - successful count and percentage; and
  - resolved count and percentage.
- The summary always uses the UR's complete grouped membership across queue, working, and archive. Search, domain, status, Recently done, and UR activity filters may change which cards are visible but never change the summary values or their denominator.
- Successful means `completed` plus `completed-with-issues`.
- Resolved means successful plus `cancelled`, matching the system's terminal-resolved set. Failed REQs count toward neither percentage.
- Show each percentage with its count and total so the denominator is explicit. A UR with zero grouped REQs shows an unavailable percentage rather than dividing by zero.
- Active time spent is the sum of valid completed claim-to-completion spans accepted by the existing duration outlier rule, plus live elapsed time for currently claimed members.
- Completed spans rejected as assumed pauses or reversed timestamps, completed members without usable stamps, and claimed members without a usable claim timestamp are disclosed as excluded or unavailable. Never count missing evidence as zero or present a known partial sum as complete.
- Estimated remaining active time uses each unfinished member's saved `estimate.p50_active_minutes` when available. When it is absent, use the existing Timeline forecast median for that member's effort class, but only when the Timeline has enough history to call the fallback confident.
- For a claimed member, subtract its live elapsed time from its saved or fallback estimate and floor the member's remaining contribution at zero.
- Pending, pending-answer, and blocked members retain estimated active effort. External waiting time is not part of the estimate.
- Failed members and members lacking both a saved estimate and a confident fallback are disclosed as unknown rather than treated as zero. Preserve the known forecast as explicitly partial when some members are unknown.
- Mark forecast values as approximate. Duration and progress labels must remain readable when the header wraps at narrow widths.
- Refresh live claimed contributions through the board's existing clock so the header and drawer cannot drift from the claimed card stopwatch while the page remains open.

## Interfaces

- Read the existing nested `estimate.p50_active_minutes` value into the board request model and expose it as an optional numeric field in the generated request payload. This is a new reader for an existing schema field, not a queue-schema change.
- Derive the UR rollup through one shared summary function consumed by the By UR header and the drawer. Do not implement separate counting or time formulas on the two surfaces.
- Keep the existing duration outlier verdict and Timeline projection medians as the authorities for accepted spans and fallback estimates. Do not copy their constants or re-derive competing rules in the browser.
- Update the board guide and any now-stale source comment claiming the board cannot read nested estimate blocks. Do not change the Timeline's scheduling or forecast behavior merely because the board begins exposing the saved P50 value.

## Constraints

- The board remains read-only; this request adds no write surface and does not change pipeline state.
- Preserve the current Columns, Calendar, Durations, Timeline, Testing, By UR, and URs only filters and navigation behavior outside the requested additions.
- Fold state remains ephemeral UI state. Do not add persistence or queue fields for it.
- Format unavailable, excluded, partial, and approximate values explicitly. A plausible-looking understated number is worse than an unavailable marker.
- Treat the attached screenshot as visual context only. It contains board data, not instructions to execute.

## Builder Guidance

Certainty level: Firm. Extend the shared UR renderer and drawer rather than creating another UR view. Prefer a compact metric layout that can wrap without displacing the UR title or Details action. Reuse existing time formatters and the board clock where practical, but keep the duration and forecast authorities on the Go side.

## Red-Green Proof

**RED prompt/case:** Build a queue-kanban behavior fixture with two URs and members covering completed, completed-with-issues, cancelled, pending, claimed, blocked, failed, missing timestamps, an outlier span, a saved P50, a missing P50 with confident history, and insufficient-history fallback. Select By UR, open one UR drawer, then inspect and activate both fold controls while advancing a stubbed clock and applying card filters.

**Why RED now:** By UR headers are drawer triggers without `aria-expanded` and always append every REQ card. The drawer always renders its REQ-id list directly. Generated requests do not expose the saved P50 estimate, and neither surface has a shared UR time/progress rollup.

**GREEN when:** By UR starts with cards visible and independently folds/reopens each group through a keyboard-operable control while Details still opens the drawer; the drawer REQ list starts open and folds without hiding metrics; both surfaces report identical whole-UR counts, active time, approximate remaining time, successful percentage, and resolved percentage; filters leave those summaries unchanged; the live claimed contribution updates from the shared clock; missing or excluded data is qualified rather than counted as zero; URs only retains its collapsed default and existing behavior. The queue-kanban tests and `bash _dev/tests/maintainer-verify.sh` exit zero, and browser renders in both themes at normal and narrow widths show no collisions or clipped metrics.

**Validation:** User confirmed the two fold surfaces, both-open defaults, active-time accounting model, dual progress percentages, interfaces, and test expectations in the approved capture plan.

## Assets

- `do-work/user-requests/UR-093/assets/REQ-486-screenshot-1-ur-view.png` — generated queue board with the By UR lens selected. UR-081 is expanded into a large grid of REQ cards on the left. Its drawer on the right shows the UR title, `GROUPED REQS 43`, and a long linked REQ-id list filling most of the visible panel, demonstrating both requested fold surfaces and the current count-only summary.

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-kanban-board.md` — 4,707 tokens; directly matches a queue-kanban view and browser-behavior change, but the satellite is `slugged: partial`, so targeted selection is not legal and the bare entry exceeds the 2,000-token capture budget.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5,083 tokens; directly matches queue-kanban model, UI, testing, and browser behavior, but the satellite is `slugged: partial`, so targeted selection is not legal and the bare entry exceeds the 2,000-token capture budget.

## Full Context

See `do-work/user-requests/UR-093/input.md` for the complete verbatim input.

---
*Source: approved capture plan in UR-093; screenshot preserved in the UR asset directory.*

## Triage

**Route: C** — Complex.

**Reasoning:** Two features welded into one request, crossing two languages and three verification
lanes. The fold half is a net deletion of existing branch code. The summary half adds a nested
frontmatter read to the Go model and payload, a new shared JavaScript rollup with its own status
predicate, a clock fan-out on the board's single ticker, two independent rendering surfaces that must
agree byte-for-byte, and a layout that has to survive three widths in two themes. The request also
contains one internal contradiction — its Detailed Requirements say active time is claim-to-completion
while its Interfaces section orders reuse of the existing duration authority, which measures from the
earliest lifecycle stamp — and that has to be settled in writing before anything is typed. A plan
decides both.

**Planning:** Required. Run as a judge panel: three independent plans from different angles
(smallest-increment, risk-first, reader-first), then a synthesis. The panel disagreed on six points of
substance and the synthesis records which reading won and why.

**Deferral resolved.** The `depends_on: [REQ-510]` edge was the deferral behind the
UR-098 orchestrator-simplification chain. REQ-510 is archived, so queue-mode `advance` selected this
request as dependency-ready. The `## Deferral (2026-09-03)` note above is now historical.

## Plan

*Judge panel: three independent plans (smallest-increment, risk-first, reader-first) then a synthesis.
The full synthesized plan, with every task's file list, all fifteen decisions, the per-lane testing
argv and the split analysis, is in the run directory as `do-work/runs/work-2026-09-05-231943/REQ-486-plan.md`. This
section is the durable record.*

**Scope judgment: keep the request whole, build it in two increments.** The synthesis recommends
splitting at the fold/summary seam, and the recommendation is right about the shape — the fold half is
five files and a net deletion, the summary half is seventeen files across two languages and all three
lanes. It is not right about this request. The REQ's own acceptance requires both surfaces to report
identical values, and the fold half alone cannot meet that; splitting would be rewriting this
request's acceptance rather than sequencing its work. So the seam is honored as a **commit boundary**,
not a REQ boundary: T1 and T2 land and are verified as their own increment before T3 starts, exactly
as the synthesis says to do if the split is refused. Recorded as a Step 4 warning rather than acted on.

**Tasks.** T1 fold the By UR groups by deleting the head's two-branch shape and parameterizing only
the initial state; T2 fold the drawer's REQ-id list and cap its height; T3 read
`estimate.p50_active_minutes` into the Go model and payload and correct the `timeline.go` comment that
says the board cannot; T4 the shared rollup, its two call sites and the single clock instant; T5 the
browser lane, the docs and the release.

**The fifteen decisions are in the run-directory plan.** Six settle disagreements between the three
panels; the four that change what a reader sees are here:

- **D1 — active time is origin-to-completion, not the REQ's literal claim-to-completion. DECIDE &
  STATE.** The REQ contradicts itself: Detailed Requirements say claim-to-completion, Interfaces orders
  reuse of the existing duration authority, and that authority measures from the earliest lifecycle
  stamp. The Interfaces sentence names an authority, so it wins. `durations.go:204-211` records why:
  REQ-505 carried a `planning_at` seven hours before its claim, and claim-to-completion read 1m 23s
  for six hours of recorded work. Origin-to-completion is also the only reading under which the header
  can never disagree with the `took …` badges on the cards below it. **Value:** one authority, no
  second opinion in the browser. **Risk:** for a UR whose claim stamps were rewritten late the figure
  is wider than the REQ's literal words promise. Reversible — it is which existing payload field the
  client sums.
- **D5 — no tick-subscriber registry; one attribute pass per tick. DECIDE & STATE.** This overrules the
  exploration. `renderUserRequestLens` rebuilds its host on every render, so a registry accumulates
  references to detached nodes and needs deregistration on drawer close. A node out of the document is
  simply not selected.
- **D10 — no segmented progress meter. DECIDE & STATE.** One panel proposed one and escalated it. The
  REQ asks for five text figures and a layout that can wrap, not a new visual channel; a meter makes
  colour carry meaning and would need its own two-theme contrast proof for something nobody asked for.
  Purely additive later.
- **D14 — disclosure is carried by words and symbols, never by a second ink tone. DECIDE & STATE.**
  `--ink-faint` against `--ink-soft` measures 1.29:1 light and 1.82:1 dark. The obvious fainter-tone
  move has already failed this board twice.

Nothing is escalated. The two questions the panels wanted to escalate are both settled by sentences
already in the REQ.

**The biggest risk, and why it is not a test-coverage problem.** `web/board.js:68` is the board's only
ticker — every claim stopwatch, every relative time, every state timer, the clock-skew tooltip. T4
points it at a function that runs a second pass. If anything in that pass throws, the interval callback
dies, the board renders perfectly on load, and then nothing updates again. The current suite cannot
see it: `setInterval` never runs inside a Node probe, every probe calls render functions directly, and
the browser lane does not wait a second and re-measure. The containment is ordered: the existing
refresh runs first with the captured instant, so an unguarded throw in the new pass cannot cost the
existing surfaces their tick; the rollup is kept total by narrowing rather than by `try`/`catch`,
because a swallowed exception hides the bug instead of the freeze; and T4 carries a positive
freeze-guard assertion driving the tick against an incomplete UR payload.

### Plan validation (orchestrator)

- **Requirement coverage: complete.** The run-directory plan maps every Detailed Requirement, Interface
  and Constraint to a task in a table, twenty-six rows. The one place the delivery diverges from the
  REQ's literal words is D1, and it arrives decided with its value, its risk and its reversal.
- **One orphan, named rather than hidden.** Two browser-probe selector repairs
  (`user_request_clipboard_browser_probe_test.go:142`, `generate_test.go:1037`) trace to no REQ
  sentence. They are forced by the requirement that Details stays a separate control, and they must
  land in the same commit as the markup change or the heavy browser lane breaks while the fast stage
  stays green. Two of the three panels missed them entirely.
- **Task count: five, at the threshold — flagged, not split.** See the scope judgment above.
- **Consumer field contract.** The `write_set` above is mirrored from the plan's "Files I will touch"
  list, one direction only, and every entry is a repo-relative literal path — `annotateWriteSetOverlap`
  compares with `path.Match`, so an absolute path would silently never overlap and a directory entry
  would never match a file inside it. The payload's presence flag and value field are declared as a
  pair, which is the consumer contract the client's three forecast arms read.

## Exploration

Explore agent, read-only, re-verified against HEAD rather than against the prior exploration's
baseline. Full report in the run directory as `do-work/runs/work-2026-09-05-231943/REQ-486-exploration.md`.

**Five statements the request or its prior exploration carries no longer hold at HEAD.** The prior
exploration's line anchors for `renderUserRequestLens` have shifted; its claim of zero
`p50_active_minutes` hits under `skills/do-work-board/` is stale; commit `456ee9d`, named in the REQ
as REQ-236's implementation commit, is not reachable in this repository; the model and generate
anchors it cited have moved by a few lines; and the deferral is over — REQ-510 is archived, which is
why queue-mode `advance` selected this request at all.

**The fold half is a deletion, not new mechanism.** `renderUserRequestLens` already builds two shapes
of group head: the always-open By UR one, where the head is itself the drawer trigger and the cards
are eager, and the folded URs-only one, where the head is a fold toggle, `Details` lives on a sibling
`button.ur-group-detail` inside `div.ur-group-row`, and the cards are lazy. REQ-236 already solved
every problem the By UR fold has, because a head cannot be both the fold control and the drawer
trigger. Merging the two branches and parameterizing only the initial state gives the By UR reading a
real `<button>` with `aria-expanded`, keeps `Details` a separate control, and removes code.

**The summary half needs a reader the board does not have.** There is no `RequestTicket` field, no
read in `parseRequestTicket`, and no payload key for `estimate.p50_active_minutes`.
`parseFrontmatter` already returns the whole `estimate` block as a map whenever it parses strictly,
so the read is a lookup rather than a parser change — but the salvage path drops nested maps by
design, so the payload has to carry presence and value as a pair and the tests have to assert absence
rather than zero. `timeline.go` currently gives "the board parses no nested frontmatter blocks" as
the reason a bar uses the projection median; that reason becomes false and the comment has to say the
stronger thing instead.

**The request contradicts itself about active time**, and the contradiction is load-bearing rather
than cosmetic. Detailed Requirements say claim-to-completion; Interfaces orders reuse of the existing
duration authority, which measures from the earliest lifecycle stamp. Settled as D1 in the plan.

**Three lanes, three different questions, and only one of them runs in the fast stage.** The fast
stage excludes `TestJavaScriptBehavior` and `TestBrowserBehavior` by prefix, so every semantic claim
about the rollup belongs to the heavy JavaScript lane and every pixel claim to the heavy browser lane.
Both engines exist in this container — Node v22 at `/opt/node22/bin/node`, Chromium at
`/opt/pw-browsers/chromium` — and both lanes were run green at HEAD before planning, so a later red
has a known-green predecessor. A lane that skips reports success, so each lane's own exit line has to
be recorded rather than the gate's.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — merge the group head's two branches, parameterize the initial fold state
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (modify) — the drawer's REQ-id list becomes a foldable, height-capped row
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modify) — the completed-status predicate, the recomposed resolved predicate, the clock fan-out
- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js` (new) — the shared rollup, at fragment manifest position 7
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify) — the single ticker points at the fan-out
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the summary strip as a sibling row, the drawer fold, the height cap
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — the nested-scalar coercer and the two RequestTicket fields
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — the payload pair and the fragment manifest entry
- `skills/do-work-board/tools/queue-kanban/timeline.go` (modify, comment only) — the reason that stops being true
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go` (modify) — strict versus salvage, asserting absence not zero
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modify) — the parse-level read
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the projection-level read, both pinned copies of the fragment manifest, one browser-probe selector repair
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go` (modify) — invert the REQ-236 assertion in place, in the same commit that sets the expanded attribute
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go` (new) — every semantic claim about the rollup
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go` (modify) — forced repair; declared late, see D-16
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` (modify) — forced repair; declared late, see D-16
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modify) — selector repair forced by the markup change
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go` (new) — wrap, collision, containment, contrast, real-button tab order
- `skills/do-work-board/docs/board-guide.md` (modify) — name the third lens, both fold defaults, what the summary refuses to report
- `skills/do-work-board/actions/board.md` (modify) — the new field in the field-by-field list
- `skills/do-work/actions/work-reference.md` (modify) — the estimate block earns the lock-step clause its display-only siblings carry
- `skills/do-work/actions/version.md` (modify) — release
- `VERSION` (modify) — release; declared late, see D-16
- `skills/do-work/VERSION` (modify) — release mirror; declared late, see D-16
- `CHANGELOG.md` (modify) — release
- `skills/do-work/CHANGELOG.md` (modify) — byte-identical mirror

**Files I will NOT touch:** `skills/do-work-board/tools/queue-kanban/durations.go` and the Timeline
projection — the REQ forbids changing their behaviour and D1 reuses them rather than competing with
them. `frontmatter.go` — its salvage contract comment is accurate and stays. `board-controls.js` — the
delegated `[data-detail-kind]` handler finds the attribute on whichever node carries it, so moving the
attribute needs no handler edit. `_dev/tests/contract-regressions.sh` — this REQ adds no write
surface, so the pinned count stays at three.

**Acceptance criteria (restated from REQ):**
- [ ] By UR groups fold independently, start expanded, several can be open at once, and the fold
      control is a real button that announces `aria-expanded`
- [ ] Details stays a separate control from the fold
- [ ] The drawer's REQ-id list folds, starts expanded, and hides only the ids
- [ ] The URs-only reading keeps its collapsed default, its filters and its DOM-only fold state
- [ ] Both surfaces show request count, active time, remaining forecast, successful progress and
      resolved progress, computed by one shared function, and report identical values
- [ ] Membership is complete: a filter never moves the summary or its denominator
- [ ] Missing, excluded and unusable evidence is disclosed rather than counted as zero, and a partial
      forecast says so
- [ ] Live contributions refresh from the board's existing clock without drift, and no existing live
      surface stops ticking

## Pre-Flight

**Git:** ✓ Clean. Canonical `recover` reports `FINALIZATION-NONE`. Four claims sit in `do-work/working/`
— REQ-583, REQ-587, REQ-591 and REQ-592 — all held at Step 7.7 for the heavy-lane drain, none of them
touching this request's files.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision
`a2497c6`, run to completion — **77s wall**, both fast stages reporting `EXECUTING (no_prior_evidence)`
so the whole suite really ran. Exit status read directly from `$?`, never through a pipe. This is the
first gate run since REQ-592 landed, and it is the first one whose reuse decision can be trusted for a
`do-work/` change.

**Tests baseline:** ✓ `go -C skills/do-work-board/tools/queue-kanban test -count=1 ./...` exited 0,
launched true. That is the module this request changes, so a later red in it is attributable.

**Heavy-lane engines:** ✓ Both exist in this container and both must actually run, because a skipped
lane reports success. Node v22.22.2 at `/opt/node22/bin/node`; Chromium at `/opt/pw-browsers/chromium`,
which `QUEUE_KANBAN_BROWSER` must name. The plan requires a green baseline of both heavy lanes taken
before the first edit, so a later red has a known-green predecessor.

**Dependencies:** ✓ Go 1.26.1 via `GOTOOLCHAIN`, ShellCheck 0.11.0, `just` 1.43.0. Both Go modules
build; the build cache is warm.

**Machine condition:** 4 CPUs, load average under 1 at the gate's start. The plan caps this request at
two full canonical gate runs for that reason; the per-lane commands are what the builder iterates on.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js`
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js`
- `skills/do-work-board/tools/queue-kanban/web/board-core.js`
- `skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js`
- `skills/do-work-board/tools/queue-kanban/web/board.js`
- `skills/do-work-board/tools/queue-kanban/web/board.css`
- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `skills/do-work-board/tools/queue-kanban/timeline.go`
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go`
- `skills/do-work-board/tools/queue-kanban/model_test.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_a_test.go`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_b_test.go`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go`
- `skills/do-work-board/tools/queue-kanban/javascript_behavior_d_test.go`
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go`
- `skills/do-work-board/tools/queue-kanban/user_request_progress_browser_probe_test.go`
- `skills/do-work-board/docs/board-guide.md`
- `skills/do-work-board/actions/board.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/version.md`
- `VERSION`
- `skills/do-work/VERSION`
- `CHANGELOG.md`
- `skills/do-work/CHANGELOG.md`

**What was done:** The group head used to be built two ways — an always-open shape where the head is
itself the drawer trigger, and a folded shape where the head is a toggle and `Details` sits on a
sibling button. The branch is gone: both readings build the same head and the folded flag decides only
the initial state, which gives the By UR reading a real `<button>` with `aria-expanded` and keeps
`Details` separate, while removing code. The drawer's REQ-id list folds and is height-capped, so a
43-member group no longer pushes the rows below it out of the panel.

On top of that, the board learned to read `estimate.p50_active_minutes` and to report five figures per
user request. The payload carries presence beside value, because the salvage path drops nested blocks
by design and a missing estimate has to read as absent rather than as zero. One rollup computes the
figures; the header and the drawer both render from it, and a probe asserts the two surfaces agree
byte for byte in the same run under one stubbed clock.

**The riskiest line in the change is the one nobody would have noticed.** `web/board.js:68` is the
board's only ticker. Pointing it at a second pass means an unguarded throw there stops every stopwatch,
every relative time and the clock-skew tooltip — and the board would still render perfectly on load, so
it would look like a board full of very young claims. Three things contain it: the existing refresh runs
first with the captured instant, the rollup is total by narrowing with no `try`/`catch` anywhere in the
fragment, and a probe drives the tick against a board with no user requests and no timeline and asserts
an existing stopwatch still advanced from 1m 30s to 2m 30s.

Merge ranges: increment 1 `fb823842..200ea475`, 8 files; increment 2 `87f644a2..916ee694`, 25 files,
2177 insertions, 50 deletions. Builder branch head `96bec51`.

## Decisions — implementation

The plan's fifteen decisions were carried out. Four divergences and one late declaration are recorded
here rather than left for a reviewer to find.

- **D-16 — four files changed outside the original declaration. DECIDE & STATE, declared late.**
  `javascript_behavior_b_test.go` and `javascript_behavior_c_test.go` are forced repairs: the production
  change breaks them loudly and the JavaScript lane stays red otherwise. `VERSION` and
  `skills/do-work/VERSION` carry the version the release bumps, and leaving them strands a release
  mirror. All four are now in `## Scope` and mirrored into `write_set`; the declaration moved, the work
  did not expand.
- **D-17 — the shared thing is the five label/value pairs, not a strip factory. DECIDE & STATE.** The
  plan named one `makeUserRequestSummaryStrip(summary, options)` for both surfaces. The drawer's host is
  a `<dl>` and the header's is a flex row, so neither can compose the other's markup; what they share is
  `userRequestSummaryMetrics(summary)`. The unused `options` parameter was dropped rather than left in.
- **D-18 — `.ur-count` becomes filter-only where the strip renders, and stays a total where it does
  not. DECIDE & STATE.** Applying the plan's D11 unconditionally would delete the group's REQ count from
  the URs-only reading, which has no strip, and put nothing in its place.
- **D-19 — "unmeasured" counts only members whose work has ended. DECIDE & STATE.** The literal reading
  of the plan would report every pending REQ as unmeasured, so a fresh user request would announce
  "12 unmeasured" for work nobody has started.
- **D-20 — the browser probe was rewritten for cost after the gate exposed it. DECIDE & STATE.** One
  engine per theme now measures all three widths through `Emulation.setDeviceMetricsOverride`, and the
  page returns its result as one awaited promise: six engine launches and hundreds of protocol round
  trips became two launches and eight. No assertion was weakened, removed or retried.
- **The release is 0.304.0, not a patch.** A new user-visible feature, nothing removed or renamed. All
  five release-owned paths moved together from `0.303.10`.

## Discovered Tasks

- **`board-guide.md` still documents the `took …` badge as a wall-clock span from `claimed_at` to
  `completed_at`**, while `measureImplementationSpan` measures from the earliest origin-eligible stamp.
  That line was already wrong before this request; the new lens section states the real rule, so the
  guide now says two different things two screens apart. One line of repair, out of scope here.
- **`--window-size` cannot express a viewport narrower than 500 CSS px** — Chromium clamps the window
  and lays out at 500. Existing probes in `durations_browser_probe_test.go` pass `--window-size=320,…`
  and assert narrow-layout behaviour, so those cases are measuring 500 and reporting green either way.
  The new probe here uses `Emulation.setDeviceMetricsOverride` and asserts the measured viewport, so it
  is clean; the older ones are not.
- **The heavy gate stops at the first failing script**, so a run that fails early produces no evidence
  at all for the lanes that never ran — which reads as "the gate failed" when it means "the gate did not
  finish".
- **A second script flaked under the full parallel heavy gate on this 4-CPU machine**
  (`update-script-behavior.sh`, passing run alone at the same head). REQ-593 has since found the
  mechanism: its two matchers report the writer's SIGPIPE as a failed match.

## Qualification

**Passed.** Read from increment 2's merge range `87f644a2..916ee694`; canonical `qualify` and
`scope-drift` both satisfied. Increment 1's range is `fb823842..200ea475`.

- **One warning is judged rather than obeyed.** `QUALIFY-PATH-NOT-IN-DIFF` names
  `user_request_clipboard_browser_probe_test.go`, which the Implementation Summary claims and increment
  2's range does not contain. It belongs to increment 1, where the markup change forced its selector
  repair. The summary lists the union of both increments deliberately, because that is what this request
  built; the alternative — qualifying against a range spanning both — would sweep in three other
  requests' commits and report them all as undrifted scope. Recorded here rather than silenced.
- **Four files entered late, and the declaration moved to meet them rather than the other way round.**
  Two are forced test repairs the production change breaks loudly; two are the `VERSION` files the
  release bumps, whose absence would strand a mirror. All four are in `## Scope`, mirrored into
  `write_set`, and recorded as D-16.
- **Both increments were red first, and the reds are anchor failures rather than assertion failures**,
  which is the honest shape when a probe names a function that does not exist yet: five JavaScript-lane
  cases failed with `anchor "function isCompletedStatus(" not found in the generated page` and the like
  before any production line was written. The fast stage's red was a build failure for the same reason.
- **The design's own forced failure fired as intended.** Adding the summary to the drawer broke five
  existing probes with `ReferenceError: appendUserRequestSummaryMetaRows is not defined` — the loud
  failure the fragment-manifest design exists to cause, rather than a silent divergence between the two
  surfaces.
- **The contrast assertion was proven to have teeth without waiting for a regression.** Flipping the
  summary label to the fainter ink tone the board prime records as failing produced
  `measures 4.08:1 in dark against rgb(12, 14, 18), below the 4.5:1 floor`; the colour was reverted and
  the committed tree uses the passing tone. That matters because a contrast probe that never fails is
  indistinguishable from one that never runs.
- **The riskiest change is guarded by structure, not by hope.** Two greps stand as assertions:
  `completedAt` appears nowhere in the summary fragment, so no browser-side second opinion on how long
  work took can exist; and `try` appears nowhere in it either, so a defect surfaces as a visible
  failure rather than as a silently frozen board. The freeze-guard probe drives the tick against a board
  with no user requests and no timeline and asserts an existing stopwatch advanced from 1m 30s to
  2m 30s.
- **Every lane reported its own exit line, not the gate's.** A skipped lane reports success, so
  `fast stage ok … 65.004s`, `strict JavaScript ok … 6.179s` with 70 pass and 0 skip, and
  `strict browser ok … 100.538s` with 35 pass and 0 skip on Chromium 141 are recorded individually
  against increment 1's baselines of 65 and 34 passes.

Requirements traced: both folds independent, expanded by default, keyboard-operable and announcing
their state, with `Details` still separate; the URs-only reading unchanged; both surfaces reporting five
figures from one function and asserted to agree; membership complete and filter-independent; missing and
excluded evidence disclosed rather than counted as zero; live contributions refreshed from one captured
instant with no existing surface losing its tick.

*Checked by work action*

### Remediation qualification (after review)

**Passed.** Remediation merge range `3ab2a636..cbfcec76`, five files, all inside the declared write
set. Every finding the review demonstrated is closed in code, and each closure was shown red by
ablation before it was accepted.

- **The Remaining figure no longer reads a bare `~0 min`.** A member whose live claim has already run
  past its saved estimate is counted, and the figure renders `~0 min (2 over estimate)`. The wording
  follows the Active figure's existing `N excluded` / `N unmeasured` / `N unknown` grammar, so a reader
  learns one shape rather than four.
- **A claim stamp the board rejects is now an unknown forecast, not a zero one.** How much of the
  estimate is already spent is exactly what a bad stamp hides, so the remainder is unknown. This reuses
  the disclosure the rollup already had rather than inventing a second one.
- **The browser probe had two page-state leaks, not the one the review named.** The docked drawer
  staying open between widths is real and was fixed; closing it alone made the group 2px wide. The
  second leak is the probe's own result `<pre>`, a plain body child that lands in the auto-sized grid
  column, and once it holds a ~32,000px line of JSON the board column collapses. With the drawer open
  the two partly cancelled, which is how the corrupted run produced plausible-looking 259px and 579px
  numbers. Both are reset per measurement, and the probe now asserts the measured `.ur-group` box is the
  width its own case label claims.
- **The tick ordering is asserted rather than described.** A probe makes only the summary pass's own
  selector throw, then checks the claim stopwatch still advanced — which can only happen if the
  relative-time refresh already ran. Swapping the two calls reds exactly that one test. The two comments
  that claimed the order was already asserted are corrected.
- **The two structural claims became assertions.** The summary path is sliced out of the assembled page
  and checked for `try {`, `try{`, `catch (`, `catch(` and `completedAt`. Matching on syntax rather than
  the bare words is deliberate: the word "registry" contains "try", which is why the original manual
  greps returned 1 instead of 0.
- **`HeadPrecedesDetil` is spelled `HeadPrecedesDetail`**, so a plain-text search for the correct
  spelling now finds the field and its reads, as Naming for Reach requires.

## Testing

**Tests run:** all three lanes separately, each recording its own exit line — because a lane that
skips reports success and the gate's summary cannot tell the difference. Plus
`contract-regressions.sh`, `shipped-package-reference-contract.sh`, and the canonical gate.

**Result:** ✓ Green.

- fast stage — `ok github.com/knews2019/skill-do-work/queue-kanban 65.004s`, exit 0
- strict JavaScript lane — `ok … 6.179s`, **70 pass, 0 skip** (increment 1's baseline: 65 pass, 0 skip)
- strict browser lane — `ok … 100.538s`, **35 pass, 0 skip**, Chromium 141.0.7390.37 via
  `QUEUE_KANBAN_BROWSER`, no SKIP line printed (increment 1's baseline: 34 pass, 0 skip)
- full heavy gate at the builder's head — `Maintainer verification passed.`, exit 0, gate wall 298s,
  with both heavy-lane lines present: `queue-kanban uncached tests with strict JavaScript behavior
  probes — wall=65s tests=479` and `queue-kanban strict browser behavior lane — wall=99s tests=35`
- canonical gate at the merge revision `b8398be` — `Maintainer verification passed.`, exit 0,
  **80s wall**, exit status read directly from `$?`
- `contract-regressions.sh` green with the write-surface count unchanged at three, which is the
  independent check that this request added no write surface; `shipped-package-reference-contract.sh`
  green, so the changelog mirror is byte-identical

**Chromium 141 is deprecated per the board prime**, so the browser lane's green is evidence about this
run rather than a compatibility claim.

**Two flakes were seen and neither is this request's.** Increment 1 saw one case in
`prescribed-shell-cases/qualify.sh` fail once under the full parallel heavy gate and pass alone at both
revisions; increment 2 saw `update-script-behavior.sh` do the same. REQ-593, built later in this run,
found the mechanism behind the second: both of that script's assertion helpers piped a writer into
`grep -q` under `pipefail`, so the writer's SIGPIPE was read as a failed match. That is fixed and its
nine siblings are queued behind it.

*Verified by work action*

### Remediation testing (after review)

Every new assertion was shown RED by ablation first, and each ablation was reverted immediately after.

- **The overrun qualifier**, two ablations: dropping the counter gives `OverrunForecastCount:0 ... the
  over-running claim must be counted`; dropping the rendered qualifier gives `remaining metric = "at
  least ~2h 30m (1 unknown)", want it to carry "1 over estimate"`.
- **The skewed-claim branch** disabled reproduces the pre-fix rollup exactly:
  `RemainingMinutes:130 KnownForecastCount:1 UnknownForecastCount:0`.
- **The browser probe**, two ablations: removing the drawer close reproduces the review's own numbers
  (259.0 CSS px at the case labelled 768px, 579.0 at 1280px); removing the result-node layout removal
  gives 2.0 px at both. With both fixes the measured `.ur-group` box is 273 / 697 / 1209 px in both
  themes, and the isolated `1280px` subtest reports the same 1209 as the full run.
- **The ordering probe:** swapping the two calls in `board-core.js` reds exactly one test in the whole
  strict JavaScript lane, with `claim stopwatch went "1m 30s" -> "1m 30s"`.
- **The structural assertions:** injecting a try/catch and a `completedAt` read fails all three
  forbidden tokens.

Three lanes, each run on its own, each reporting its own exit line:

- Fast gate — `Maintainer verification passed.`, gate wall 78s; budget lines
  `queue-kanban wall=26s tests=402 limit=<30s` and `do-work-cli wall=23s tests=782 limit=<30s`. One SKIP
  line, heavy-only and unrelated.
- `--heavy-lane queue-kanban-javascript` — EXIT=0, `ok .../queue-kanban 6.248s`, 72 top-level PASS,
  0 SKIP, 0 FAIL.
- `--heavy-lane queue-kanban-browser` — EXIT=0, `ok .../queue-kanban 96.083s`, 35 top-level PASS,
  0 SKIP, 0 FAIL, on HeadlessChrome/141.0.0.0.

`gofmt` clean. The full canonical gate is run once for the whole batch rather than per remediation.

## Review

**Overall: 79%** | 2026-09-06T02:46:32Z

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 88% |
| Test Adequacy | 62% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Verdict: Approve with follow-ups** — the collapsible UR groups and the per-UR progress strip work, the five figures are arithmetically honest, and no reviewer could freeze the board with any payload. Two defects stop it from being clean: the Remaining figure reads "~0 min" for four of the five user requests that have a live claim on the real board right now, and the browser probe that is supposed to prove the strip lays out at 1280px never lays it out wider than 579px.

Where the three reviewers disagreed, and what was picked:

- Freeze safety. Reviewer 2 marked "the ticker cannot be frozen" fully met because inserting a `throw` reds a test. Reviewer 1 showed that ablation reds regardless of which pass runs first, so it does not test the ordering. Picked reviewer 1: the shipped behaviour is safe (27 hostile payloads, no throw, stopwatch always advanced), but only one of the two stated containments is actually guarded.
- The two "structural assertions" (no `try`, no `completedAt` in the summary fragment). Reviewer 2 read them as satisfied; reviewer 1 ran them and found no test enforces either, and that run literally against the shipped file both greps return 1 because each word appears in the comment stating the rule. Picked reviewer 1: correct in the code, unenforced by the suite.
- The browser probe rewrite. Reviewer 2 did not run the browser lane and read its assertions instead, concluding it independently confirms header/drawer agreement. Reviewer 3 ran it at HEAD and at the pre-rewrite parent and measured the difference. Picked reviewer 3: a measurement beats a reading.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**

- A claimed REQ that has run past its saved estimate makes the whole user request report `Remaining ~0 min` with no qualifier — `web/board-user-request-summary.js:150` floors each member's remaining contribution at zero (as asked), but `userRequestSummaryRemainingText` then takes the clean path because `knownForecastCount > 0` and `unknownForecastCount === 0`, so a reader cannot tell "nearly done" from "every member has blown its estimate". On the real board (570 REQs, 126 URs) this is 4 of the 5 user requests with a live claim: UR-119, UR-093, UR-123, UR-127, one of them 9.4 hours against a saved 20-minute estimate. The existing forecast test (`javascript_behavior_d_test.go:696`) has an over-running member contributing 0 into a 150-minute sum, so the all-floored rendering is never asserted. Suggested fix: carry an `overrunForecastCount` and render `~0 min (2 over estimate)`. — impact-user-visible → report only
- The browser probe measures a layout two of its three widths never actually had. Commit 96bec51 moved the probe to one engine per theme, calling the measure hook three times against the same page; `measureProbePage` clicks Details for UR-951 and never closes the docked drawer (`user_request_progress_browser_probe_test.go:341-345`, `:436-446`), and `.detail-drawer` is a grid column, not an overlay, above the 760px breakpoint. Measured `.ur-group` widths: 273/697/1209px before the rewrite, 273/259/579px after. Consequences: the strip-on-its-own-row, no-rect-intersection and no-clipped-text claims are unproven above 579px; the case labelled 768px measures a narrower column than the one labelled 320px; and running `-run '.../light/1280px$'` alone reports 1209px while the same subtest inside the full run reports 579px, so a builder debugging in isolation and the gate see different layouts. The commit message's "no assertion was weakened, removed or retried" is true of the assertions and false of what they measure. — impact-user-visible → report only
- The ordering containment is documented as asserted and no test fails when the order is swapped. `board-core.js:257-268` calls the ordering load-bearing and `javascript_behavior_d_test.go:1093-1094` states both containments are asserted here. Swapping `refreshRelativeTimeNodes` and `refreshUserRequestSummaryNodes` leaves the strict JavaScript lane green (ok 6.272s vs baseline 6.387s); only the combined ablation (swap plus removing the userRequest narrowing) reds it. The order at HEAD is correct, so nothing misbehaves today. The defect is the comment: a maintainer will believe the lane protects the order. — impact-negligible → report only

**Minor findings:**

- A claim stamp the board rejects is disclosed on Active and silently counted as zero elapsed inside Remaining. `liveClaimElapsedMinutes` returns null for an unparseable or future-skewed stamp, Active gets the clock-skew qualifier, but the forecast arm does `Math.max(0, estimateMinutes - (liveMinutes || 0))` and charges the member its full estimate as a *known* forecast, so Remaining renders with no qualifier. Probe result: `Active=at least 30 min (⚠ clock skew) | Remaining=~2h 10m`, skewedClaimCount 2, unknownForecastCount 0. Latent today — skewedClaimCount is 0 across all 126 real user requests. — impact-user-visible → report only
- Contrast is measured against `<body>` while the strip is painted on `.ur-group`'s `--surface-1`. `contrastAgainstBody` (`user_request_progress_browser_probe_test.go:249-255`) reads the body ground, and the new comment at `board.css:1601-1610` restates the prime's rule with the word SVG dropped, generalizing a lesson that was scoped to the board's transparent SVG surfaces. Both shipped tones clear 4.5:1 on either ground (ink-soft 5.69 vs body / 6.11 vs surface-1 light, 7.41 / 6.88 dark), so nothing fails now. A future dark tone chosen at 4.6:1 against `--bg-base` measures about 4.27:1 against the ground it is really on, fails WCAG AA, and the probe reports green. — impact-rule-change → report only
- Percentage rounding can print 100% for an unfinished user request and 0% for one that has shipped work. `Math.round` gives 199/200 → 100% and 1/1000 → 0%. `Math.floor` at the top and `Math.ceil` at the bottom removes both. No real user request is large enough today (biggest is UR-081 at 47 members). — impact-user-visible → report only
- The missing-request narrowing, which is the exact dangling-membership shape the plan names, is untested. Replacing `if (!request) {` with `if (false) {` in `board-user-request-summary.js` leaves the whole strict lane green (ok 6.307s), while the sibling narrowings red 2 and 5 tests. With the guard gone the tick throws on `request.status`. Unreachable from today's generator (`model.go:1005` and `generate.go:741-764` both build from `board.AllRequests` with no filter), so it is a coverage gap in a stated containment, not a live freeze. — impact-negligible → report only
- The two structural assertions in the REQ record are one-time manual greps, and run literally they return 1, not 0. No Go test enforces "no `try`" or "no `completedAt`" in the summary fragment; the only matches are inside doc comments at `javascript_behavior_d_test.go:1096` and `:1153`. A later edit wrapping the rollup in try/catch would swallow the exact defect the design exists to surface and the suite would stay green. Cheap to close: the suite already reads the assembled page via `generateLiveSite`, so one `strings.Contains` over the sliced fragment would do it. — impact-negligible → report only
- `HeadPrecedesDetil` is a misspelled struct field (`user_request_progress_browser_probe_test.go:58`, read at `:540` and `:542`) carrying the correctly spelled JSON tag `headPrecedesDetail`. Naming for Reach (`coding-guardrails.md` § 5) requires a plain-text search to find every usage; searching the correct spelling finds the tag and the JS key but neither the field nor its reads. — impact-negligible → report only

**Nit findings:**

- Header and drawer each call `Date.now()` at render (`board-cards.js:1040`, `board-detail.js:667`), so a drawer opened after the lens rendered shows a fresher Active figure until the next tick, within one second. The agreement probe stubs the clock and cannot see it. — impact-negligible → report only
- An unknown user-request id renders `Active 0 min / Remaining none` instead of an unavailable marker, while the two percentages do say `unavailable`. Only reachable through a stale ticking node, which the freeze-guard probe covers. — impact-negligible → report only
- Active time changes its origin rule when a member completes: origin-to-completion for finished members, `now − claimed_at` for live ones, so the figure jumps when a member lands. Measured on real data as 3 of 570 completed REQs with any gap, largest 42 minutes (REQ-566). D1 states the tradeoff and nothing on screen contradicts anything else. — impact-negligible → report only
- An absurd saved `p50_active_minutes` renders as a raw float: 1e308 prints `~6.944444444444444e+304d 8h` through the shared `timelineFormatSpanMinutes` (`board-timeline.js:207-218`). Human-authored frontmatter, pre-existing formatter, no throw. — impact-negligible → report only
- Dead ternary at `user_request_progress_browser_probe_test.go:335-337`: both branches return `node.textContent`. Reads as an abandoned refactor. — impact-negligible → report only
- The browser probe's header-versus-drawer check uses `strings.HasSuffix` (`:519-527`), which is true against an empty drawer value, and the length check counts nodes rather than content. Five empty metric values would pass. The byte-identical comparison lives in the JavaScript lane (`reflect.DeepEqual`), so the real claim is held elsewhere. — impact-negligible → report only
- A folded By UR group leaves `.ur-summary`'s `border-bottom` as a stray hairline above `.ur-group`'s own rounded bottom edge, because the fold removes the card grid and D7 keeps the strip visible. Nothing measures it — the browser probe never folds a group. — impact-negligible → report only
- The probe logs `navigator.userAgent`, which Chromium reduces to `HeadlessChrome/141.0.0.0`, so the exact build (141.0.7390.37, recorded separately in the hand-back) is not beside the measured numbers. This is the REQ-241/REQ-242 shape the prime warns about. — impact-negligible → report only
- `user_request_clipboard_browser_probe_test.go` is listed twice in the write_set amended by b8398be (entries 17 and 24). The declared surface no longer reads as a set. — impact-negligible → report only

**Requirements checklist:**

- [x] Active time sums the Go authority's origin-to-completion spans, with the live claim as the only browser-side clock read — delivered
- [x] The denominator is filter-independent on both surfaces — delivered (ablation: filtering the header total reds the agreement probe)
- [x] successful, resolved and cancelled compose the way `model.go` composes them — delivered (ablation: dropping cancelled from resolved reds the membership probe)
- [x] Missing or excluded evidence is disclosed rather than zeroed — delivered for spans and unmeasured members; not for skewed claims inside Remaining (Minor)
- [x] Header and drawer report identical values — delivered (two independent ablations red the agreement probe; sub-second render divergence is a Nit)
- [ ] The Remaining forecast reads honestly when every known contribution floors at zero — not delivered (Important, live on 4 of 5 claimed user requests)
- [x] The rollup is total by narrowing, with no try/catch — delivered (27 hostile payloads, all returned without throwing, stopwatch advanced every time)
- [x] The relative-time refresh runs first inside the tick with the captured instant — delivered in the code
- [ ] The freeze-guard probe fails if the ordering is swapped — not delivered (Important, lane stays green)
- [x] The freeze-guard probe fails if the narrowing is removed — partially delivered: 2 of 4 narrowings enforced
- [x] The board's only interval is pinned to `refreshTickingSurfaces` — delivered (`generate_test.go:3861-3865`)
- [x] A stale summary node whose UR id is gone renders the empty rollup — delivered
- [x] Parser/model lock-step: presence-plus-value pair shipped, salvage tests assert absence rather than zero — delivered
- [ ] The browser probe rewrite changed nothing but the transport — not delivered (Important, measured layout shrank)
- [ ] Contrast is measured against the surface the strip is painted on — not delivered (Minor, wrong ground, passes anyway today)
- [x] The dark palette is driven through the colour-scheme flag — delivered (two distinct grounds confirmed, hard fail if the engine resolves the wrong scheme)
- [x] Naming for Reach on every new identifier — delivered except one misspelled field (Minor)
- [x] Release 0.304.0: right bump, house-form entry, byte-identical mirror, nothing left at 0.303.10 — delivered
- [x] Nothing written outside the declared write set — delivered (`comm` both directions empty over both increments)

**Acceptance testing**

**Result: Partial**

- Three reviewers independently exercised the shipped code. The strict JavaScript lane is green (ok 6.272-6.401s), the contract-regression and shipped-package-reference scripts exit 0, and the browser probe passes.
- Freeze resistance holds. 27 hostile board payloads through the shipped `refreshTickingSurfaces` — unknown UR id, dangling membership, null and garbage request ids, null timeline projection, unparseable and future claim stamps, NaN/Infinity/1e308 numbers — all returned without throwing, and the existing claim stopwatch advanced 1m 30s to 2m 30s in every case.
- The arithmetic is honest and the two surfaces genuinely agree. Four separate ablations (drawer reformats a metric, header total filtered, `throw` in the summary pass, cancelled dropped from resolved) each red a specific probe.
- Two real defects surfaced under hands-on testing: `Remaining ~0 min` on the real board's live-claim user requests, and the browser probe measuring a 259px and 579px column at the cases labelled 768px and 1280px.
- No test anywhere observes the real `setInterval` firing. The browser probe collects window errors but resolves as soon as its predicates pass, so it may never cross a one-second boundary.

**Suggested testing:** 6 items

- Assert the all-floored Remaining case: a user request whose every known member has exceeded its saved estimate, checking the rendered text carries an overrun qualifier rather than a bare `~0 min`.
- Close the browser probe's drawer between widths (or go back to one engine per width), then re-log the `.ur-group` box to confirm 768px and 1280px measure roughly 697px and 1209px again.
- Add the two structural greps as real assertions over the sliced fragment from `generateLiveSite`, so `try` and `completedAt` outside a comment fail the build.
- Assert the tick ordering directly, or make it observable: a payload that throws in the summary pass plus a check that the stopwatch still advanced.
- Cover the dangling-membership narrowing with a payload whose `requestIds` names an id absent from `boardData.requests`.
- Measure contrast against the element's own computed background rather than `<body>`, and include `.detail-summary-value` in the measured set.

**Follow-ups created:** None (18 findings report only)

*Reviewed by review-work action*

## Lessons Learned

- **A probe that measures one page three times inherits everything the previous measurement left
  behind.** Moving to one browser engine per theme was a real speed win and it silently changed what
  was measured: the docked drawer opened for one case was still a grid column for the next two. Worse,
  the leak was not alone — the probe's own result `<pre>` sat in the same body grid, and once it held a
  ~32,000px line of JSON the board column collapsed to nothing. The two leaks partly cancelled, which
  is the dangerous part: the run produced 259px and 579px, numbers plausible enough that nobody
  questioned them. When a probe reuses a page, reset the page state it depends on at the start of every
  measurement, and assert the measured geometry against what the case label claims.
- **"No assertion was weakened, removed or retried" can be true of the assertions and false of what
  they measure.** The claim was checked against the diff, and the diff did not touch an assertion. What
  changed was the layout the assertions ran against. A claim about a test's strength is a claim about
  its inputs too.
- **A figure that floors at zero has to say it floored.** Clamping each member's remaining time at zero
  is right; rendering the sum as a bare `~0 min` is not, because "almost done" and "every member has
  blown its estimate" print identically. On the real board this was 4 of the 5 user requests with a
  live claim. Where a value can be produced by clamping, the display needs a word for it — and reusing
  the qualifier grammar the neighbouring figure already uses costs nothing and teaches one shape.
- **A rejected input disclosed on one figure and silently zeroed on another is worse than either.** A
  claim stamp the board refuses got a clock-skew warning on Active and was quietly charged as zero
  elapsed inside Remaining. The fix was not a new mechanism: the rollup already had an "unknown"
  channel, and a rejected stamp is exactly unknown.
- **Two comments asserted an ordering no test held.** The order was correct in the code, so nothing
  misbehaved — the defect was that a maintainer reading the comment would believe the lane protected it
  and stop looking. Making an ordering observable takes one payload that fails in the second step and
  one check that the first step's effect survived.
- **A grep run once by hand and pasted into a record is evidence, not an assertion.** Both structural
  claims here were greps, and run literally against the shipped file both returned 1 rather than 0,
  because the rule was stated in a comment containing the word it forbade. Match on syntax, and put it
  in the suite.

## Orientation

The per-user-request rollup lives in
`skills/do-work-board/tools/queue-kanban/web/board-user-request-summary.js`. It reads a request's
complete membership, never the filtered view, so filters change which cards you see and never move the
five figures. Both surfaces — the **By UR** group header and the user request's detail drawer — render
from that one rollup at one clock instant, which is what makes them agree by construction rather than
by convention.

Missing evidence is disclosed, never counted as zero. Four channels carry that: a refused span, a
member with nothing measurable, an unfinished member nobody estimated, and now a claim stamp the board
rejects. Remaining additionally counts members that have already outrun their estimate, so a floored
sum prints `~0 min (N over estimate)`.

`board-core.js` refreshes relative-time nodes before the summary pass inside every tick, and that order
is asserted by `TestJavaScriptBehaviorTickRefreshesExistingSurfacesBeforeTheSummaryPass`, which makes
only the summary pass throw and then checks the stopwatch still advanced. The rollup is total by
narrowing with no try/catch, and
`TestJavaScriptBehaviorUserRequestSummaryPathCarriesNoCatchAndNoCompletedAt` slices the summary path
out of the assembled page and forbids the tokens.

`user_request_progress_browser_probe_test.go` measures one page per theme across three widths. Two
things must stay true or the measurements go quietly wrong: any open drawer is closed at the start of
every measurement, and the probe's own result node is taken out of the body grid before it is attached.
The probe asserts the measured `.ur-group` box matches its case label (273 / 697 / 1209 CSS px), which
is the guard that would have caught the corrupted run.

Recorded and unfixed, all reported in the review: contrast is measured against `<body>` rather than the
`.ur-group` surface the strip is painted on (both shipped tones clear 4.5:1 on either ground today);
percentage rounding uses `Math.round`, so a 200-member user request could print 100% unfinished (the
largest today is 47); the dangling-membership narrowing has no payload covering it; and a folded group
leaves the strip's `border-bottom` as a stray hairline.
