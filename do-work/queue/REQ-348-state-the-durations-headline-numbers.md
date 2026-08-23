---
id: REQ-348
title: "State the Durations view's own headline numbers"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-068
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-346]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-342, REQ-343, REQ-344, REQ-345, REQ-346, REQ-347, REQ-349, REQ-350]
batch: durations-panel-improvement
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# State the Durations View's Own Headline Numbers

## What

The view makes a reader hover to learn the numbers it exists to report. Add a stat tile row above
panel A, a rolling median line on panel B, and the ticks panel C is missing.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The summary line gives counts and the exclusion rule and nothing else. The median, the p90 and the
completions-per-active-day rate — the numbers a reader opens a durations view for — appear nowhere.
Panel B annotates exactly one day (its slowest) and draws no trend; panel C labels only its peak,
with no zero tick and no gridline. Every one of those figures is derivable from the payload already
shipped.

## Detailed Requirements

- **A small stat tile row above panel A**: median REQ duration, p90, active-day count against the
  axis span, and REQs per active day.
- **Each figure is stated for the window actually plotted**, not for all history.
- **A rolling median line on panel B** over its per-day bars, so "is it getting slower" is answerable
  without reading 29 bars one at a time.
- **Panel C gets its zero and midpoint ticks.** Today only its peak is labelled.
- **Derive every figure from the samples already in the payload.** No new aggregate in `durations.go`.
- **Keep one measure per scale**, and keep the existing summary sentence's exclusion-rule wording
  intact.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- Panels A and B disagree about which samples count on purpose, and that rule lives in one place.
  Whichever panel a tile describes, it uses that panel's rule and says so if the two could be
  confused.
- Generate a board and look at it.

## Dependencies

`depends_on: REQ-346` — the tiles are stated for the plotted window, so the window state must exist
first. This resolves the report's own "whether the tiles follow A3's window or all history": they
follow the window.

## Builder Guidance

**Certainty: firm.** Four tiles, one line, two ticks, all from existing data. Keep it small — this is
the cheapest legibility gain in the batch and it should not grow a new aggregation layer.

## Red-Green Proof

**RED prompt/case:** Generate a board for this repository and open Durations. Nowhere on the view does
it say the median is 13.2 minutes or the p90 is 42.7; panel B draws 29 bars with one annotation and no
trend; panel C's only labelled value is its peak.

**Why RED now:** The summary composes counts and the exclusion rule only, and neither panel B nor
panel C draws an aggregate or a full tick set.

**GREEN when:** a tile row above panel A states median, p90, active days against axis span and REQs
per active day for the plotted window; panel B carries a rolling median line; panel C has zero and
midpoint ticks; and the summary sentence's exclusion-rule wording is unchanged.

**Validation:** User confirmed (bundled invocation).

---
*Source: prompt A7, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html` (finding F7).*
