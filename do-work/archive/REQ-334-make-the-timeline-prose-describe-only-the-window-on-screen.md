---
id: REQ-334
title: "Make the timeline prose describe only the window on screen"
status: pending
created_at: 2026-08-23T12:26:00Z
user_request: UR-066
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-327, REQ-328]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Make the Timeline Prose Describe Only the Window on Screen

## What

Four sentences the view emits assert things the window they describe does not carry. Each is a small edit;
they are one REQ because they share one rule — a sentence about the chart may only claim what the current
window actually draws.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

These are the sentences a reader quotes and screenshots. The forecast paragraph in particular contradicts
itself inside one sentence pair, which reads as a broken board rather than as a subtlety.

## Context

- **P1 — "measured to the now-line at &lt;t&gt;" in windows with no now-line.** `renderSummary` (`:1157`)
  always appends it, while `drawNowRule` (`:1506`) draws nothing when `nowMs` is outside the window. Press
  Week on any past week and the summary names a now-line that is not on screen, beside open bars clipped
  flush at the frame. Say where the measurement came from without implying a rule is visible.
- **P2 — the forecast contradicts itself under a filter.** `renderTimelineForecast` (`:913`) prefixes "This
  covers the whole queue, not the rows shown." and then, when the chain is empty, adds "Nothing left to
  schedule — every remaining REQ is listed below." Both appear together: measured with `domain=ui-design`
  (1 row), the paragraph reads "This covers the whole queue, not the rows shown. Nothing left to schedule —
  every remaining REQ is listed below." while one row is drawn and the excluded paragraph immediately under it
  names REQ-325, which is not listed below. The "listed below" half is only true when the rows on screen are
  the whole queue — which is exactly the case the prefix denies.
- **P3 — the rightmost axis label names a day the window excludes.** `2026-07-06 → 2026-07-13` prints a
  `13 Jul` tick while the `to` field calls 12 Jul the last day included. REQ-327 changes the tick instants;
  this is the wording half of that decision and lands after it.
- **P4 — the legend's "Vertical rules" group omits the gridlines.** `template.html:296-298` names the now
  rule and the queue-end rule. The plot draws a third kind of vertical line, one per axis tick
  (`drawGridlines`, `:1445`), and the reader has nothing to look it up under.

## Detailed Requirements

1. The summary names the now-line only when the now-line is drawn. When it is not, it still has to say what
   an open span was measured against — the instant, without claiming a visible rule.
2. The forecast's "every remaining REQ is listed below" is emitted only when the rows on screen really are the
   whole queue. Under any subset it says what is true of a subset. Keep the "whole queue, not the rows shown"
   prefix, which is right and load-bearing for the numbers under it.
3. The excluded paragraph is consistent with whichever branch the forecast took: it must not name a REQ as
   excluded directly under a sentence claiming everything remaining is listed.
4. The legend gains an entry for the gridlines, so every vertical line the plot draws has a key.
5. Every count in every sentence matches the number of rows the chart draws in that state. This already
   holds at Fit all (317 REQs / svg height 5706 / 317 table rows all agree) and must keep holding.

## Constraints

- `depends_on: [REQ-327, REQ-328]`. P3 is the wording half of REQ-327's tick decision, and the summary's open
  count is wrong for a different reason until REQ-328 lands — fixing the sentence first would pin the wrong
  number.
- Do not blur the two empty-state messages: "No REQ matches the current filters" and "Nothing was drawn
  between …" have different remedies (`:1026-1028`).
- The forecast's DOM is rebuilt only when the subset answer flips (`:1166-1176`); keep that, and do not add a
  per-frame rebuild for the sake of a sentence.

## Red-Green Proof

**RED prompt/case:** (P1) Set From 2026-07-06 / to 2026-07-12 and read `#timeline-summary`. (P2) Filter
`domain=ui-design` and read `#timeline-forecast` and `#timeline-excluded` together. (P4) Read the legend and
count the kinds of vertical line in the plot.

**Why RED now:** (P1) the summary says "measured to the now-line at 2026-08-23 11:13 UTC" for a window ending
2026-07-13, where no now-rule is drawn. (P2) the forecast says "not the rows shown" and "every remaining REQ
is listed below" in the same breath, above one row, with REQ-325 named as excluded below it. (P4) the legend
lists 2 vertical rules; the plot draws 3 kinds.

**GREEN when:**
- No window that lacks a now-rule produces a sentence naming one.
- No state produces "every remaining REQ is listed below" while the drawn rows are a subset of the queue.
- The legend accounts for every vertical line the plot draws.
- A Node behaviour probe drives `renderSummary`'s and `renderTimelineForecast`'s branches over the windows the
  view can reach and asserts each sentence against that window's own facts.

**Validation:** Inferred during capture; each sentence above was read out of a live render.

## Full Context

See `do-work/user-requests/UR-066/input.md`.
