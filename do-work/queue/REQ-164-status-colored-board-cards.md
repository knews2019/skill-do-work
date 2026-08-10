---
id: REQ-164
title: "Status-colored board cards"
status: pending
created_at: 2026-08-10T21:30:31Z
user_request: UR-035
domain: ui-design
prime_files: []
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
