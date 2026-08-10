---
id: REQ-164
title: "Status-colored board cards"
status: completed
claimed_at: 2026-08-10T21:41:17Z
completed_at: 2026-08-10T21:58:15Z
commit: 2905360
route: A
created_at: 2026-08-10T21:30:31Z
user_request: UR-035
domain: ui-design
prime_files: [skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: false
suggested_spec: ui-component
depends_on: []
maintenance: false
related: []
batch: status-colored-board-cards
---

# Status-Colored Board Cards

## What

Make board cards easier to scan by giving workflow states distinct, restrained visual treatments across card-based views, particularly the **By UR** lens where mixed statuses currently appear uniformly gray.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reuse the existing semantic accent/tint tokens, bind card-level CSS variables to each recognized `data-status`, add a renderer class for unrecognized statuses, style a neutral-text status pill plus 3px rail, remove the By-UR neutral override, and add generated-site source-contract coverage before running the queue-kanban Go tests.
- [x] **[APPLY]:** Updated the shared request-card CSS/renderer and generated-site regression test exactly as planned; the focused regression failed on the old By-UR override before implementation and passes after it.
- [x] **[UNIFY]:** Reviewed `web/board.css` for semantic mappings, neutral bodies, focus/hover preservation, and removal of the By-UR neutral override; `web/board.js` for card-level invalid-state propagation without schema changes; and `generate_test.go` for complete generated-site source-contract coverage. `gofmt`, `go test ./...`, `go vet ./...`, and `git diff --check` are clean with no debug artifacts.

## Why

Completed work should look different from pending work. The current mixed-status presentation is too visually bland to scan quickly.

## Context

The standard column view already carries semantic status accents, but cards in the By UR lens reset to the same neutral accent. The user confirmed a **rail + pill** treatment and asked that every workflow state receive a distinct treatment.

## Detailed Requirements

- Use a 3px status-colored left rail and a compact, softly tinted status pill while keeping the card body neutral.
- Map `pending` to amber and `claimed` to blue.
- Map `pending-answers`, recognized blocked variants, and `failed` to red.
- Map `completed` and `completed-with-issues` to green; their written labels continue to distinguish the exact outcome.
- Map `cancelled` to neutral gray and retain its existing struck-through title treatment.
- Map unrecognized statuses to the red invalid treatment, including the card rail and pill.
- Apply the status treatment consistently anywhere standard REQ cards render, especially the mixed-status By UR lens.
- Do not add a separate legend; keep the written status inside every pill so status remains understandable without color.
- Preserve hover behavior, keyboard focus visibility, responsive wrapping, and both light and dark themes.
- Keep data schemas, commands, generated board data, and other public interfaces unchanged.
- Add generated-site regression coverage for the status mappings and the removal of the By UR gray override.
- Visually verify representative pending, claimed, blocked, completed, completed-with-issues, failed, cancelled, and invalid cards in light and dark modes.

## Constraints

- Reuse the board's existing semantic accent and tint palette rather than introducing a second color system.
- Color must remain a redundant cue: readable status text and existing cancellation/invalid markers stay present.
- Avoid broad card-body tinting; the agreed emphasis is the rail and status pill.
- Preserve existing filtering, bucketing, drawer, and card interaction behavior.

## Dependencies

None.

## Builder Guidance

**Certainty level: Firm.** The visual direction, state mapping, accessibility fallback, and validation cases were confirmed during capture. Implementation details may follow the board's existing CSS/JavaScript structure.

## Red-Green Proof

**RED prompt/case:** Open the board in light mode, select **Board → By UR**, and inspect a UR containing at least one pending card and one completed card. Both currently use the same gray left edge and tiny gray status treatment, so they are not distinguishable at a glance.

**Why RED now:** The By UR card grid neutralizes the status accent, and the current status marker is too visually quiet to separate mixed workflow states.

**GREEN when:** The pending card has an amber 3px rail and amber-tinted status pill, the completed card has a green rail and green-tinted status pill, every other state follows the confirmed mapping, the card bodies remain neutral, and the labels remain readable in light and dark modes.

**Validation:** User confirmed the rail-and-pill emphasis, the all-states mapping, and this capture after reviewing the implementation plan.

## Assets

`do-work/user-requests/UR-035/assets/REQ-164-current-by-ur-board.png` — light-mode By UR board showing eight mixed-status `UR-355` cards rendered with nearly identical gray visual treatment; the `UR-355` detail drawer is open on the right.

## Full Context

See `do-work/user-requests/UR-035/input.md` for the complete verbatim request, confirmed choices, constraints, and screenshot description.

---
*Source: "think about a visual aid, completed should have a different color then pending, as it is now it is very bland" — rail + pill and all-state mappings confirmed during capture.*

---

## Triage

**Route: A** - Simple

**Reasoning:** This is a focused styling change in the existing board card renderer and stylesheet, with a narrow generated-site regression seam and no new component or public interface.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** Bound the shared request-card rail and pill to semantic status tokens, widened the rail to 3px, removed the By-UR gray override, and propagated invalid status state to the card. Added generated-site regression coverage for every recognized mapping, invalid statuses, the compact pill contract, and By-UR consistency.

## Qualification

Passed — all 3 implementation files are present in the diff and substantive. Every detailed requirement traces to either the card-level status variables, invalid-state renderer class, retained focus/hover rules, or generated-site regression case; the existing board data flow remains unchanged and still supplies `status` plus `statusUnrecognized` to the shared card renderer. Mechanical qualification, `git diff --check`, formatter, test, and vet checks are clean.

## Testing

**Tests run:** `go test ./...`, `go vet ./...`, `gofmt -w generate_test.go`, `git diff --check`, and local-browser visual/computed-style checks in the By-UR lens at 1440×1000 light, 1280×900 dark, 768×900, and 320×760.

**Result:** ✓ All passing. Light and dark checks covered pending, claimed, pending-answers/blocked, failed, both completed variants, cancelled, and invalid cards; every rail measured 3px, card bodies stayed neutral, pills retained written labels, cancellation remained struck through, and invalid styling reached the whole card. Keyboard focus retained its 2px outline, hover retained lift/background feedback, and status labels remained visible at both responsive widths.

**Red-green validation:**
- `TestGenerateInlinesSemanticStatusCardStyles`: ✗ before implementation on the By-UR neutral override → ✓ after implementation.

**New tests added:**
- `TestGenerateInlinesSemanticStatusCardStyles` in `skills/do-work-board/tools/queue-kanban/generate_test.go`.

**Existing tests updated (cross-REQ impact):** None.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-10T21:57:12Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — all status mappings, shared-card behavior, accessibility cues, themes, and responsive widths passed source, automated, and browser checks.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Request-card status identity now lives at card level in `web/board.css`; `makeRequestCard` supplies `data-status` plus the invalid-state class, and `TestGenerateInlinesSemanticStatusCardStyles` guards the shared Columns/By-UR contract. Architecture and public interfaces are unchanged.
