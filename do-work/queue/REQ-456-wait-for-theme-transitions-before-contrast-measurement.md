---
id: REQ-456
title: 'Wait for theme transitions before contrast measurement'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-454, REQ-455, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Wait For Theme Transitions Before Contrast Measurement

## What

Wait for the completion card's browser-reported theme animations to finish after changing emulated color scheme, then sample computed colors for contrast. Use a browser condition rather than a fixed sleep.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this transition-time contrast-sampling root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #15 — P2 — source:** `completion_contrast_browser_test.go:52-55`

> ````text
> [P2] Wait for theme transitions before measuring contrast — [prj].claude/skills/do-work-board/
> tools/queue-kanban/completion_contrast_browser_test.go:52-55
> After Emulation.setEmulatedMedia changes the color scheme, the card has 120 ms background/color transitions, but the test now
> samples computed styles immediately. The dark-theme case can therefore measure an intermediate light-to-dark palette and
> intermittently report incorrect contrast; wait for the card's browser-reported animations before reading colors.
> ````

- **Evidence:** `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go:46-62` changes emulated media and samples immediately, while `web/board.css:824-843` applies a 120 ms background transition. Replaying Chrome 151 ten times produced a failure where dark text was sampled against still-light `rgb(241,244,248)`, yielding 2.36:1. The browser harness contract requires waiting on browser-observable conditions.
- **Surface-cost result:** Earned — the flake is reproduced. One browser-reported animation wait is cheaper than intermittent false contrast failures; do not add a fixed-duration sleep.

## Detailed Requirements

- After each emulated theme change, wait until the completion card's relevant browser-reported animations have finished.
- Read computed foreground and background colors only after the wait condition succeeds.
- Preserve both light- and dark-theme contrast assertions.
- Keep the test sensitive to genuine insufficient contrast.
- Use browser-observable animation state rather than elapsed-time guessing.

## Constraints

- Do not add a fixed 120 ms or padded sleep.
- Keep production transitions unchanged unless an independent product defect is proven.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm. Prefer the browser's animation API or an equivalent state predicate already supported by the harness.

## Red-Green Proof

**RED prompt/case:** Repeatedly switch the contrast browser fixture from light to dark emulated media and sample through the current immediate path.
**Why RED now:** Computed styles can be captured during the 120 ms transition, reproducing a false 2.36:1 dark-theme result.
**GREEN when:** Repeated runs wait for reported animations, consistently sample the settled palette, and still fail for an intentionally low-contrast settled fixture.
**Validation:** User confirmed after validate-feedback accepted Finding #15.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #15, captured by UR-085.*
