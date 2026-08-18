---
id: REQ-252
title: Record the browser with every measured-face number in the Durations tests
status: claimed
created_at: 2026-08-18T13:56:12Z
claimed_at: 2026-08-18T19:12:47Z
route: B
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
depends_on: [REQ-248]
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/durations.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T19:13:28Z
  basis:
    - Route B
    - 3-file write set
    - 3 acceptance criteria
    - browser evidence
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

## Ordering Gate

- [~] **D-01: gated behind REQ-248 rather than left free-running.** Two reasons, and the second is the binding one. Textual: both REQs write `skills/do-work-board/tools/queue-kanban/generate_test.go`, and this session proved that two REQs in one Go package can collide in a way git cannot see — REQ-241 and REQ-242 each declared `durationsMeasuredAxisTitleAscentUnits` in different files, never touched adjacent lines, and the merge failed to compile. Semantic: this REQ records the provenance of measurements taken against Panel B's geometry, and REQ-248 moves that geometry. Recording provenance for a layout about to change would produce numbers that are wrong on arrival.
  **Value:** the one-line `do-work run` resume cannot schedule these two together, because auto-wave reads `depends_on` and deliberately does not read `write_set`.
  **Risk:** if REQ-248 is abandoned, this REQ stays gated behind it and needs the gate cleared by hand. Acceptable — REQ-248 is the highest-value item in the queue.

---

## Triage

**Route: B** - Medium

**Reasoning:** Documentation-fidelity sweep over measured constants plus a fresh per-build measurement of the Panel B budget; the what is clear, the current numbers need deriving in the live DOM on this build.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — every measured constant names its browser build
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — same, and the Panel B budget prose gains both measurements + the per-browser warning
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — hyphen-vs-U+2212 comment; remainder-reserve over-reservation stated as deliberate; the REQ-260 gofmt nit at the `Truncate(24*time.Hour)` line may be fixed in passing (same file, one character)

**Files I will NOT touch:**
- `web/board-durations.js` — REQ-248's geometry landed; no behaviour change belongs to this REQ.
- Any constant's *value* — this REQ records provenance; a measurement demanding a change is captured separately.

**Acceptance criteria (restated from REQ):**
- [ ] Every browser-measured constant names the browser and build in the same comment as the number.
- [ ] The Panel B clearance budget prose states it is per-browser and gives both measurements.
- [ ] No behaviour change; suite still exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.14 tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium present

*Checked by work action*
