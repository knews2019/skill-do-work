---
id: REQ-322
title: "Name the REQ on its own timeline row"
status: pending
created_at: 2026-08-22T22:08:34Z
user_request: UR-065
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-318, REQ-319, REQ-320, REQ-321, REQ-323, REQ-324]
batch: timeline-ux-audit
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-timeline.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Name the REQ on Its Own Timeline Row

## What

Show each REQ's title in the row's label column, and put its detail in a tooltip at the
pointer instead of only in a readout at the foot of the panel.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The 104px label column holds `REQ-012` and nothing else. The title lives in a one-line
readout at the very bottom of the panel — below a five-sentence hint paragraph, some 700px
from the pointer — and in the collapsed table. Scanning three hundred rows tells the reader
nothing about what any of them are.

## Context

`renderVisibleRows` draws one `<text class="timeline-row-label">` per row holding `row.id`.
`timelineRowDescription` already composes the full sentence (id, title, route, status, both
spans, projection slot, anomaly reason) and is used for both the row's `aria-label` and the
foot readout. The pieces exist; they are just not where the eye is.

`TIMELINE_LABEL_WIDTH = 104` is subtracted from the plot width in `plotWidth()`, so widening
the label column narrows the plot. That trade is the substance of this REQ: enough width to
recognize a title, not so much that the chart loses its span.

## Detailed Requirements

- Each row shows its title beside its id, truncated to fit the label column.
- Truncate against a **measured** face, not an estimated character width. This module has
  been bitten twice: `REQ-292` (a width model returns the same number for every face, so
  slots never move and a wider face draws past them) and `REQ-241`/`REQ-242` (the same 12px
  face measured 12.0372 and 11.2300 on two Chromium builds). Record the browser and build
  beside any measured number, and take the larger where two disagree.
- Hovering a bar shows a tooltip at the pointer carrying what the foot readout carries.
  A native SVG `<title>` on the row group is the cheapest thing that works and needs no
  positioning code; a custom tooltip is fine if it earns the complexity.
- The foot readout stays as the `aria-live` region — it is the non-pointer path and it is
  not this REQ's to remove.
- A REQ with no title (nothing in `requestsById`) still renders its id and does not draw an
  empty tooltip.

## Constraints

- The plot must stay wide enough to be a chart. If the label column grows, say by how much
  and show the rendered result at a narrow viewport, not just at full width.
- Virtualization: the label is rebuilt per visible row on every scroll, so whatever measures
  the text must not measure per row per frame.
- Serial with the rest of the `timeline-ux-audit` batch.

## Builder Guidance

**Certainty: Exploratory.** The user asked for a more useful view, not for a title column —
this is capture's answer to "make it more useful UIUX", and the trade at its centre is a
judgment call: every pixel the label column gains, the chart loses. Pick a width, render it
at a full and a narrow viewport, and expect to move the number. Say what you picked and
what the render showed.

The tooltip and the title column are separable. If measured-face truncation turns out to
cost more than it returns, a native SVG `<title>` on the row group alone is an acceptable
smaller landing — say so rather than spending the REQ's budget on the harder half.

## Red-Green Proof

**RED prompt/case:** Open the Timeline tab and, without moving the pointer or opening the
table, name any three REQs on screen. Impossible — every row reads `REQ-0NN` and nothing
else. Hover one bar and the description appears at the bottom of the panel, below the hint
paragraph.

**Why RED now:** The row label is `row.id` alone, and the only place a title is rendered is
the foot readout and the collapsed table.

**GREEN when:** each row shows a truncated title next to its id; a title too long for the
column is cut with an ellipsis at a measured boundary rather than overrunning into the plot;
hovering a bar raises a tooltip at the pointer carrying id, title, status and both spans;
and the foot readout still announces the same text for a screen reader.

**Validation:** Inferred during capture — this is an audit finding, not one of the user's
four items. The user asked for a full audit and for the view to be more useful.

## Assets

Screenshot described in `do-work/user-requests/UR-065/input.md` — a column of bare ids, and
the hovered row's description stranded at the foot of the page.

---
*Source: audit finding, UR-065 — "audit the timeline view, and make it more useful UIUX."*
