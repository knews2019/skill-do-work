---
id: REQ-347
title: "Colour Durations marks by a chosen channel"
status: pending
created_at: 2026-08-23T22:37:52Z
user_request: UR-069
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-346]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-346, REQ-348, REQ-349, REQ-350, REQ-351, REQ-352, REQ-353, REQ-354]
batch: durations-panel-improvement
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-durations.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - skills/do-work-board/tools/queue-kanban/web/board.css
---

# Colour Durations Marks by a Chosen Channel

## What

REQ-346's lane answers "which dots are one UR" positionally. Add the second half the user asked for:
a colour-by control on the Durations view offering route, UR and domain, so identity is also
readable inside the plot.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The proposal offered the lane alone as the cheaper, reversible half and left colour as an open
question; the user chose both. The lane scales past 60 URs where a colour channel does not, so the
two are complementary rather than redundant: the lane always answers grouping, colour answers it
without a second glance in the cases a palette can carry.

## Detailed Requirements

- A colour-by control on the Durations view with three channels: **route** (today's behaviour, the
  default), **UR**, **domain**.
- **The active channel is always named on screen.** A mark's fill means something different under
  each channel, so a legend or caption that states which channel is live is part of the change, not
  polish. Two identity channels competing silently is the failure mode this REQ has to avoid.
- **UR colouring degrades honestly.** 66 URs exceed any usable categorical palette: decide and state
  the rule for what happens past the palette's capacity (a shared "other" bucket, colouring only the
  URs visible in the window, or another rule you can defend), and make the legend say it.
- The control follows the same visibility discipline REQ-353 applies to the topbar: it appears on
  Durations and nowhere it does nothing.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first.
- Generate a board and look at it, in both light and dark. A categorical palette that reads in one
  theme and not the other is not done.
- Do not change what a mark *is* — position, radius and the read-time exclusion rule stay as they
  are. This REQ changes fill and nothing else.

## Dependencies

`depends_on: REQ-346` — the UR join and its lane land first.

## Builder Guidance

**Certainty: firm that all three channels are wanted, open on the palette and the overflow rule.**
Route colour is the default so a reader who never touches the control sees today's board. All three
channels ship: if the domain channel turns out to carry almost no signal on real data (most REQs
share a domain), that is an honest reading of the archive and worth noting in the hand-back, not a
reason to omit the channel. Dropping one is a decision for the maintainer, not the builder.

## Red-Green Proof

**RED prompt/case:** Generate a board, open Durations. There is no way to colour the marks by
anything but route, so a dense column of interleaved URs looks identical whichever request its marks
belong to.

**Why RED now:** Mark fill is bound to route with no channel selector anywhere in the view.

**GREEN when:** a colour-by control offers route, UR and domain; switching it recolours panel A's
marks; the on-screen legend names the active channel; and the UR channel's behaviour past the
palette's capacity is both implemented and stated in the legend.

**Validation:** User confirmed (capture answer Q4: lane plus colour-by toggle).

---
*Source: capture answer to the report's open question Q3, `ai-reports/2026-08-23_2200_durations-panel-improvement-proposal/index.html`.*
