---
id: REQ-439
title: 'Anchor trailing timeline windows before display padding'
status: pending
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
related: [REQ-437, REQ-438, REQ-440, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Anchor Trailing Timeline Windows Before Display Padding

## What

Make trailing-window controls anchor to the latest semantic recorded, drawn, or projected endpoint before cosmetic range padding is added. Keep padded bounds for rendering and clamping without allowing “Last day” or similar chips to select only empty display space.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
