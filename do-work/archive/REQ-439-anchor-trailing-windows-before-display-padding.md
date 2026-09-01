---
id: REQ-439
title: 'Anchor trailing timeline windows before display padding'
status: completed
route: A
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T21:50:05Z
  basis:
    - trivial short-circuit
related: [REQ-437, REQ-438, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T21:48:36Z
completed_at: 2026-09-01T22:01:29Z
---

# Anchor Trailing Timeline Windows Before Display Padding

## What

Make trailing-window controls anchor to the latest semantic recorded, drawn, or projected endpoint before cosmetic range padding is added. Keep padded bounds for rendering and clamping without allowing “Last day” or similar chips to select only empty display space.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Preserve the latest semantic endpoint immediately before display padding, pass it through both production trailing-window call sites, and add a test-first JavaScript probe that executes the shipped bounds block and nested apply caller for the 95-day drained case.
- [x] **[APPLY]:** Preserved the semantic endpoint before display padding and threaded it through both production trailing-window computations; added the test-first production-caller regression.
- [x] **[UNIFY]:** Reviewed `board-timeline.js` and `generate_test.go` in full diff context; verified both chip call sites share the endpoint, padded bounds still own clamping/rendering, existing helper callers retain fallback behavior, and no debug artifacts were introduced. Focused trailing-window tests pass.

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Anchor trailing windows before adding display padding.` The feedback repeats this as Finding 21 with the same 95-day replay.
- **Evidence:** `board-timeline.js:1684-1715` reads and then pads the semantic range; `timelineTrailingWindow` at lines 401-410 clamps its anchor to that already-padded `boundEndMs`.
- **Origin / earned by:** REQ-425 (Stop the Timeline's Trailing-Window Controls Assuming Now and a Full Screenful Are Inside the Bounds)/`cbe9de07` fixed chips collapsing when `now` was out of bounds but conflated semantic and display endpoints. A 95-day drained board receives 1.9 days of right padding, making “Last day” contain no recorded activity; the existing test calls the helper without production padding.
- **Surface-cost:** N/A. Separating the semantic anchor from cosmetic bounds is a direct correction; add one production-path regression rather than a new control abstraction.

## Detailed Requirements

- Preserve the unpadded latest meaningful endpoint before adding display padding.
- Use that endpoint to anchor finite trailing-window chips when `now` lies beyond recorded activity.
- Retain padded bounds for drawing, panning, zooming, and final window clamping.
- Cover a drained history longer than 50 days through the production caller, not only a pure helper supplied with unpadded bounds.

## Constraints

- Do not remove the visual breathing-room padding.
- Do not add a second timeline state model; carry only the semantic endpoint required by the controls.

## Red-Green Proof

**RED prompt/case:** Render a drained 95-day history whose current time is beyond the range, apply production's 2% padding, and select “Last day.”
**Why RED now:** The chip anchors to the padded end, 1.9 days after the latest activity, so its one-day window is entirely empty.
**GREEN when:** “Last day” ends at the unpadded meaningful endpoint and intersects recorded activity, while rendered bounds still retain their padding.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Preserve the visual padding and correct only the control's semantic anchor.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Findings 18 and 21 from the validated external feedback.*

## Triage

**Route: A** - Simple

**Reasoning:** The regression is isolated to one renderer-local value being overwritten before two known call sites, with an existing JavaScript behavior harness for the exact 95-day case.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

The renderer stored semantic and display ends in the same `boundEndMs` variable. Adding 2% breathing-room padding overwrote the last meaningful endpoint before the production chip handlers called `timelineTrailingWindow`, so finite windows anchored to empty display space even though the pure helper tests supplied unpadded bounds.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** The renderer now preserves the latest recorded, drawn, or projected endpoint before adding display padding and supplies it to both finite trailing-window call sites. A production-path JavaScript regression executes the shipped padding block and nested chip caller over a drained 95-day history, while confirming padded display bounds remain intact.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed. The new test follows the shipped bounds block and production apply caller, both consumer call sites use the preserved endpoint, and the diff contains no debug artifacts or unrelated changes.

## Testing

**Tests run:** focused trailing-window JavaScript behavior tests; `go vet ./... && go test ./... -count=1` in the queue-kanban module; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing, including the strict JavaScript behavior lane and canonical maintainer verification. The optional browser lane was unavailable and skipped by the canonical gate; no browser-only evidence is required for this arithmetic-only correction.

**Red-green validation:**
- `TestJavaScriptBehaviorTimelineTrailingWindowAnchorsBeforeDisplayPadding`: ✗ before implementation (“Last day” ended at the padded boundary 1.9 days after activity) → ✓ after (ends at the meaningful endpoint while padded display bounds remain)

**New tests added:**
- Production bounds/caller regression for a drained 95-day history with 2% display padding

*Verified by work action*

## Review

**Overall: 98%** | 2026-09-01T22:01:05Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — finite chips anchor before cosmetic padding, all display-bound consumers keep their padded clamps, and the production 95-day replay is green.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Timeline trailing-window controls now distinguish meaningful queue time from cosmetic chart padding inside the board renderer. The Kanban board prime remains current.
