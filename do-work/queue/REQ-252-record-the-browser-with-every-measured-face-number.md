---
id: REQ-252
title: Record the browser with every measured-face number in the Durations tests
status: pending
created_at: 2026-08-18T13:56:12Z
status_changed_at: 2026-08-18T13:56:12Z
user_request: UR-051
addendum_to: REQ-241
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: durations-measured-face-constants-lack-provenance
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/durations.go
---

# Record the Browser With Every Measured-Face Number in the Durations Tests

## What

The Durations suite now holds several constants that are a **browser's answer**, not arithmetic — and two independent measurements of the same face, taken on different Chromium builds, disagreed by enough to matter. Every such constant needs the build it came from recorded next to it, and the numbers that are now known to be per-browser need saying so.

## Context

Found when REQ-241 and REQ-242 both measured the 12px axis-title face and landed on different values, then declared the same constant in different files of one package. The collision surfaced at integration as a Go redeclaration — git could not see it, because the two edits never touched adjacent lines. It was resolved to the larger value (12.1) on the reasoning that a title box reaching higher makes every clearance test demand more room than the render needs.

That resolution is sound and stays. What it exposed is that these numbers carry no provenance, so the next disagreement looks like a mistake rather than a difference of build.

## Instances

- [ ] **The Panel B clearance budget is per-browser and roughly 7× thinner on a current build than recorded.** REQ-241's D-03 measured 1.364 units of headroom above Panel B's title and named a row pitch of 15 as the point where it breaks. REQ-242's reviewer measured **0.185 units** on Chromium 146 — still positive, still zero intersections, and demonstrably not caused by REQ-242 (the SVG for that region is byte-identical across its range). Anyone about to spend that budget must re-measure first, and the prose should say so.
- [ ] **No measured-face constant records the Chromium build it came from.** `generate_test.go`'s block documents the procedure and viewport, which is most of the way there; `durations_test.go`'s do not. A number that differs between builds and does not name its build cannot be re-derived or argued with.
- [ ] **`durations.go`'s `formatDurationLabelMinutes` emits an ASCII hyphen where the renderer draws U+2212.** The two glyphs differ by 1.73 units in this face, so the Go side models a narrower string than the browser draws. **Currently width-neutral** — both are one character and the width model counts characters — but it becomes a real under-estimate the moment anyone replaces the flat constant with a per-glyph model, which REQ-241 attempted and abandoned. Worth a comment at minimum.
- [ ] **`durationsLabelRemainderReserveUnits` over-reserves more at the new width.** It scales as `24 × durationsLabelCharacterWidthUnits`, so raising that constant to 7.15 widened the reservation to 171.6 against a widest composable remainder sentence of 123.18 units. Over-reserving is the safe direction — it drops a label, which the remainder counts — but the gap is now large enough to be a deliberate choice rather than an accident.

## Requirements

- Every constant in these files that is a browser measurement names the browser and build it was taken on, in the same comment as the number.
- The Panel B clearance budget's prose states that it is per-browser and gives both measurements, so nobody spends a budget that does not exist on their machine.
- No behaviour change is required by this REQ. If a measurement turns out to demand one, capture that separately rather than folding it in.

## Red-Green Proof

**RED prompt/case:** this is largely a documentation-fidelity REQ; the mechanical half is whether any assertion currently reads a measured constant without a recorded provenance. A check that every `durationsMeasured*` constant's comment names a browser would fail today.
**Why RED now:** `durations_test.go`'s measured constants name no build, and the two files disagreed by 0.86 units on the same face.
**GREEN when:** each constant carries its provenance and the suite still exits 0.
**Validation:** Integration seam on REQ-241/REQ-242 plus REQ-242's review finding F3.
