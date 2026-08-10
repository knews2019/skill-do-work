---
id: UR-035
title: Status-colored board cards
created_at: 2026-08-10T21:30:31Z
requests: [REQ-164]
word_count: 21
---

# Status-Colored Board Cards

## Summary

Make mixed-status board cards easier to scan, especially in the **By UR** lens, by giving each workflow state a restrained but distinct visual identity. The user confirmed a colored left rail plus a softly tinted status pill, applied across all states while keeping the card body neutral.

## Extracted Requests

| Request | Intent |
| --- | --- |
| REQ-164 | Distinguish board-card workflow states with an accessible rail-and-pill color treatment. |

## Batch Constraints

- Use a 3px status-colored rail and a softly tinted status pill; do not tint the whole card.
- Distinguish all workflow states: pending amber, claimed blue, blocked/pending-answers/failed red, completed variants green, cancelled gray, and invalid statuses red.
- Keep written status labels so the interface does not rely on color alone.
- Apply the treatment consistently to card-based views, especially the mixed-status By UR lens, without adding a separate legend.
- Preserve hover, keyboard focus, responsive layout, and light/dark theme behavior.
- Add generated-site regression coverage and visually verify the result.

## Screenshot Context

The supplied light-mode screenshot shows the local queue board at `127.0.0.1:8090` with the **By UR** lens active. The `UR-355` group contains eight cards with a mixture of `claimed`, `pending`, and `completed` statuses, but every card uses the same pale gray body, gray left edge, and tiny gray status marker. Completed work is therefore not visually distinguishable from pending work at a glance. A detail drawer for `UR-355` is open on the right half of the screen.

Permanent asset: `do-work/user-requests/UR-035/assets/REQ-164-current-by-ur-board.png`

## Confirmed Choices

- **Emphasis:** colored rail plus softly tinted status pill.
- **Status scope:** visually distinguish every workflow state, not only pending versus completed.

## Full Verbatim Input

```text
think about a visual aid, completed should have a different color then pending, as it is now it is very bland
```

---
*Captured: 2026-08-10T21:30:31Z*
