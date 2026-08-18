---
id: REQ-231
title: Keep Panel A's direct labels clear of the mark band
status: pending-answers
domain: general
created_at: 2026-08-18T00:55:10Z
user_request: UR-051
addendum_to: REQ-226
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Discovered Task: Keep Panel A's Direct Labels Clear of the Mark Band

## What

In the Durations view's overflow lane, the first row of direct labels sits at the same height as the marks themselves, so in a dense lane a label can be crossed by a *neighbouring* mark. REQ-226 stopped labels from overprinting each other; this is the remaining overlap, between a label and a dot that is not its own.

## Context

Found while implementing REQ-226. `DURATIONS_LANE_MARK_Y` is 40 with a mark radius of 5, so marks occupy roughly y 35-45; `DURATIONS_LANE_LABEL_ROW_Y` is 44, so a first-row label's text box occupies roughly y 33-46. A label always clears its *own* mark, because it is drawn 9 units to one side of it, so the overlap only appears where the band is dense enough for other marks to crowd the label.

REQ-226's collision rule is deliberately label-against-label — its Requirement 1 says "skip any that would collide with one already placed". Extending it to label-against-mark would drop labels wherever the band is dense, which is the opposite of what the remainder count was added for. The straightforward fix is instead geometric: give the lane roughly 12 more user units so both label rows sit below the marks, and shift the panels beneath it down to match.

REQ-226 could not do that: its constraints state "Panel A's existing scale-break design is correct and stays. This REQ fixes the labelling on top of it; it does not redesign the panel." Changing the lane's height is that redesign, so it is a separate decision.

Visible in REQ-226's synthetic 60-sample fixture; invisible on this repository's own board, whose lane carries three samples.

## Requirements

- No first-row label in the overflow lane may share vertical space with the mark band, at any density.
- Panel A's scale break, overflow lane, `60+` tick, and two label rows all stay — this is a spacing change, not a redesign of the device.
- Panels B and C shift down by whatever Panel A grows by; `DURATIONS_MEDIAN_TITLE_Y` is the panel-A/B boundary `describeAtPointer` keys on, so the hover readout must still resolve the same panel for the same pointer position.
- The reversed band gets the same treatment or an explicit note saying why it does not need it.

## Red-Green Proof

**RED prompt/case:** A test asserting that the mark band (`DURATIONS_LANE_MARK_Y` ± the mark radius) and every label row's text box (baseline minus ascent, baseline plus descent) do not intersect, read from the renderer's own constants the way `TestDurationLabelGeometryMatchesTheRenderer` already reads them.
**Why RED now:** marks span roughly y 35-45 and row 0's text box roughly y 33-46, so the two intersect over about 10 units.
**GREEN when:** the same test passes, and re-rendering REQ-226's synthetic dense fixture shows the first-row label clear of the mark blob.
**Validation:** Discovered during REQ-226; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-226: in the board's Durations chart, the top strip that holds the very long REQs draws its first line of text at the same height as the dots themselves. Each dot's own label is offset sideways so it never covers its own dot, but where many long REQs finish close together, a *neighbouring* dot can sit on top of a label. REQ-226 stopped the labels covering each other; this is the leftover case of a label being covered by a dot. Fixing it means making that strip about 12 units taller and moving the two charts below it down to match — a small layout change, which is exactly what REQ-226 was told not to do ("Panel A's existing scale-break design is correct and stays"), so it is your call rather than mine. It shows up on a board like the one in your screenshot and not on this repository's own board, which has three long REQs in that strip. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
