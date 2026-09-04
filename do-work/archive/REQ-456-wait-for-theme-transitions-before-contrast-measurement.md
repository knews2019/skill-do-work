---
id: REQ-456
title: 'Wait for theme transitions before contrast measurement'
status: completed
completed_at: 2026-09-02T20:14:45Z
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
claimed_at: 2026-09-02T19:35:52Z
dispatch_at: 2026-09-02T19:44:23Z
builder_handback_at: 2026-09-02T19:59:29Z
integration_at: 2026-09-02T19:59:29Z
review_at: 2026-09-02T20:13:29Z
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-02T19:36:04Z
  basis:
    - trivial short-circuit
release_at: 2026-09-02T20:14:46Z
commit: 0ac93b69399a23c1940b1ab62277b5799da488ad
kb_status: promoted
kb_entry: REQ-456-wait-for-theme-transitions-before-contra.md
---

# Wait For Theme Transitions Before Contrast Measurement

## What

Wait for the completion card's browser-reported theme animations to finish after changing emulated color scheme, then sample computed colors for contrast. Use a browser condition rather than a fixed sleep.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this transition-time contrast-sampling root cause.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the required board primes, lessons, crew rules, and bug-fix spec. Reproduce the immediate-sampling defect with the existing generated-board/DevTools harness and a temporary transition-duration mutation, then make the browser test await the completion card's reported animations before reading computed styles. Replay the same mutation and an intentionally poor settled palette, run the focused browser test repeatedly plus the module and strict browser lanes, and leave production CSS unchanged.
- [x] **[APPLY]:** The contrast probe now awaits every animation reported by the completed card through each animation's `finished` promise after setting emulated media and viewport metrics, before any computed-style read. Root cause: the former path synchronized neither the theme change nor the card's CSS transition timeline with sampling, so dark text could be compared with an intermediate light card. The production palette and transitions are unchanged.
- [x] **[UNIFY]:** Reviewed `completion_contrast_browser_test.go` for exact post-emulation ordering, bounded browser-promise behavior, retained light/dark and genuine-contrast assertions, and absence of debug artifacts; reviewed this REQ for the P-A-U/root-cause record. `web/board.css` was restored byte-for-byte after temporary mutation proofs. `git diff --check`, `gofmt`, `go vet ./...`, the focused Chrome 152 test at `-count=10`, `go test ./...`, and both maintainer strict JavaScript/browser behavior lanes passed.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the exact browser regression test and requires one focused observable wait without production changes.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-kanban-board.md` — 4707 tokens, over the 2000-token budget; `slugged: partial` so no targeted family form. Matched on queue-kanban browser behavior. Read separately under the prime's touch-conditional Lessons Discipline.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/completion_contrast_browser_test.go` (modified)

**What was done:** The real generated-board contrast probe now awaits every animation reported by the completed card after each emulated theme and viewport change, before reading computed foreground and background colors. Production CSS remains unchanged.

## Qualification

Passed — 1 file verified, 5 requirements traced, P-A-U confirmed. The wait is after both emulation calls and before every computed-style read; the existing light/dark, body/card, hierarchy, and insufficient-contrast assertions remain intact.

## Testing

**Tests run:**
- `QUEUE_KANBAN_BROWSER='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' go test -run '^TestBrowserBehaviorCompletionCompanionsKeepReadableContrast$' -count=10`
- `go test ./...`
- strict JavaScript behavior lane
- strict browser behavior lane with Google Chrome 152.0.7977.66

**Result:** ✓ All passing

**Red-green validation:**
- Completion-card transition mutation: ✗ immediate path sampled dark text against the intermediate light card at 2.36:1 → ✓ awaited path waited on the 2-second browser animation and passed.
- Settled low-contrast palette mutation: ✗ remained detected at 1.45:1 against the body and 1.28:1 against the card.

**New tests added:**
- Existing real-browser contrast regression strengthened at its theme-switch seam; no separate test file added.

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-02T20:13:29Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings:** None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — Chrome 152 passed the focused real-browser probe across both themes and all three viewport widths on three independent runs.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Waiting on the card's own Web Animations `finished` promises makes the contrast probe synchronize with browser state and fail closed if an animation is cancelled, without guessing at elapsed time.

**What didn't:** Stretching the CSS transition duration alone did not reliably expose the immediate-sampling path; a temporary harness-level animation made the RED timing window deterministic and was restored after the proof.

**Worth knowing:** Preserve a settled low-contrast mutation alongside the flake reproduction. A wait can eliminate false failures while accidentally weakening the assertion, so the same probe must still reject a genuinely poor final palette.

## Orientation

The queue-kanban real-browser contrast probe now measures completed cards only after their theme animation state settles; the production board palette and transitions are unchanged.
