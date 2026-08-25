---
id: REQ-347
title: "Colour Durations marks by a chosen channel"
status: completed
completed_at: 2026-08-24T13:53:03Z
commit: 26f347f
claimed_at: 2026-08-24T13:33:27Z
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
  - skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Colour Durations Marks by a Chosen Channel

## What

REQ-346's lane answers "which dots are one UR" positionally. Add the second half the user asked for:
a colour-by control on the Durations view offering route, UR and domain, so identity is also
readable inside the plot.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the Kanban prime and implementation rules. Keep the existing client-side REQ join and Panel A geometry; add a Durations-only route/UR/domain control, select fills through one renderer seam, and render a dynamic named legend/readout. URs sort by ID, the first twelve use a categorical palette, and later IDs share the named Other URs bucket; missing values use a named visually distinct unknown mark. Lock the behavior through the existing Node renderer stub before implementation, then inspect light and dark output.
- [x] **[APPLY]:** Implemented the three channels, named missing and overflow states, themed categorical fills, and renderer-level behavior coverage within the declared files.
- [x] **[UNIFY]:** Reviewed the five declared files and their diff; ran JavaScript syntax checks, focused renderer behavior coverage, the module suite, and a live Durations preview.

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
- **Every channel has a rule for samples that lack its value, and the legend states it.** This is not
  only the UR channel's problem: nine of this repository's 305 samples carry no domain at all
  (REQ-001 through REQ-007 plus REQ-010 and REQ-011, all pre-dating the field), and the durations
  aggregate includes every one of them. Give them an explicit unknown bucket that is named in the
  legend and visually distinct from a real category — never an arbitrary default colour, and never
  absent from the legend while still drawn in the plot. Apply the same treatment to a missing UR on
  the UR channel, so one rule covers all three channels.
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

## Triage

**Route: B** - Medium

**Reasoning:** The behaviour and acceptance criteria are specific, but the existing Durations control, palette, legend, and browser-probe conventions must be located before implementation.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`web/board-durations.js` owns mark fills, the client-side sample-to-request join, and the Durations legend/readout. `boardData.requests[sample.id]` already provides `userRequestId` and `domain`, so no payload change is needed. `web/board-controls.js`, `web/template.html`, and `web/board.css` contain the established control, markup, and themed-palette patterns. `generate_test.go` has the renderer and DOM-stub seam for a behaviour-level lock-in.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — select and apply route, UR, or domain mark fills and dynamic reader-facing labels
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified) — Durations-only colour channel control and rerender handling
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified) — control group and dynamic legend host
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — light/dark categorical and unknown-state presentation
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — renderer-level TDD coverage for channels, overflow, and unknown values

**Files I will NOT touch:** Go payload and aggregation code; mark geometry, radius, and reversed-stamp semantics.

**Acceptance criteria (restated from REQ):**
- [ ] Durations exposes route, UR, and domain colour channels only while its view is active.
- [ ] The active channel, missing-value treatment, and UR palette overflow rule are stated on screen.
- [ ] UR and domain use the existing client-side request join; no new measurement or payload duplication is introduced.
- [ ] Reversed-stamp marks retain their critical treatment; mark geometry and read-time exclusions do not change.
- [ ] Route remains the default, switching recolours Panel A, and the screen-reader text follows the active channel.
- [ ] Focused renderer tests and a light/dark browser render verify the change.

## Decisions

**D-01:** Use the first 12 sorted UR IDs as the named categorical palette, route subsequent IDs to a shared `Other URs` fill, and render missing UR/domain values as outlines. **Reasoning:** the order is deterministic, the legend makes the limit explicit, and neither overflow nor missing data pretends to be an identifiable category.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — adds route, UR, and domain colour selection; dynamic legend/readout text; deterministic UR overflow; and explicit unknown-value treatment
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (modified) — wires the Durations-only colour control and immediate rerender
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified) — adds the three-channel control and dynamic legend host
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — adds paired light/dark categorical palette tokens and outlined unknown markers
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — exercises channel changes, deterministic overflow, missing values, reversed stamps, and accessible copy through the renderer seam

**What was done:** Panel A now defaults to route colour but can switch to UR or domain. Each channel names its meaning on screen, UR overflow shares a clearly named bucket, missing values remain visibly distinct, and reversed stamps keep their critical fill.

## Qualification

Passed — all five declared project files exist, contain substantive changes, and are wired into the Durations renderer. `qualify.sh` passed; the Implementation Summary and Scope file lists match exactly.

## Testing

**Tests run:**
- `GOTOOLCHAIN=go1.26.1 go test -count=1 -run '^TestJavaScriptBehaviorDurationsColourChannelsNameAndRecolourPanelA$' -v .` — pass
- `node --check web/board-durations.js && node --check web/board-controls.js` — pass
- `GOTOOLCHAIN=go1.26.1 go test -count=1 ./...` — pass (40.010s)
- `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh` — pass; its optional browser lane skipped because no browser binary is configured

**Render evidence:** A freshly served live board was exercised in both themes. Route-to-UR switching changed all 40 sampled mark appearances; the live UR view rendered 313 marks, 12 outlined unknowns, the 12 named palette fills plus `Other URs`, and the required legend/ARIA copy. The inspected dark render had `body` background `rgb(12, 14, 18)` and the same visible categorical/unknown treatment.

**Red-green validation:** The renderer test now proves the captured GREEN behaviour. The interrupted builder did not leave a durable RED transcript, so that pre-change run cannot be independently reproduced from this checkout.

## Review

**Overall: 96%** | 2026-08-24T13:53Z

| Dimension | Score |
|---|---:|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 94% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** No implementation findings. The initial review held acceptance only for a missing dark render; a regenerated dark preview was then inspected and cleared the hold.

**Acceptance:** Pass — route, UR, and domain channels work in the live view; their reader-facing rules, overflow behaviour, unknown states, and reversed-stamp treatment were verified in light and dark renders.

## Lessons Learned

**What worked:** Rendering the live board exposed the actual category mix and confirmed that the legend, outlines, and `Other URs` bucket remain readable at archive scale.

**What didn't:** A light-only render did not meet the REQ's visual evidence requirement; the review caught the gap before archival.

**Worth knowing:** The Durations client consumes `boardData.requests` inside its module closure, so browser inspection must validate rendered marks and copy rather than assume that join data is globally visible.

## Orientation

Durations now has a second identity channel: a reader can keep route colour or switch Panel A to UR or domain without altering the measured data, mark geometry, or anomaly treatment.
