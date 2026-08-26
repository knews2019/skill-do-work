---
id: REQ-376
title: 'Raise the done line''s faint text to readable contrast'
status: pending-answers
created_at: 2026-08-26T14:40:00Z
user_request: UR-074
addendum_to: REQ-374
domain: ui-design
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
---

# Raise the Done Line's Faint Text to Readable Contrast

## What

The done line's faint companions — `.relative-time` and `.elapsed-duration` — measure roughly 3.3:1 against `<body>` in both themes at 11px, under the 4.5:1 needed for text that size.

## Context

Discovered while working on REQ-374. Pre-existing for the whole line, not introduced by that REQ: measured light `rgb(108,116,128)` at 0.85 opacity on `rgb(245,247,250)`, dark `rgb(107,116,128)` at 0.85 on `rgb(12,14,18)`. REQ-374's new span reading deliberately inherits the same treatment rather than diverging mid-line, which is why the fix belongs to the line as a whole.

Per the prime, measure against `getComputedStyle(document.body).backgroundColor` — the board's SVG and card surfaces are transparent, so a tone judged against a `--surface-*` token is measured against something the reader never sees.

A related but separate gap, reported by REQ-374's review and not covered here: both new markers convey their meaning only through `title`, which a screen reader does not reliably announce and which is unreachable by keyboard on a non-focusable span.

## Red-Green Proof

**RED prompt/case:** compute the contrast ratio of `.relative-time` and `.elapsed-duration` against `getComputedStyle(document.body).backgroundColor` in a rendered board, in both themes. Both read about 3.3:1.
**Why RED now:** the tokens and the 0.85 opacity together land under 4.5:1 for 11px text.
**GREEN when:** both read at least 4.5:1 in both themes, with the done line still visibly quieter than the card title — the point of the treatment is hierarchy, and a fix that flattens it has traded one defect for another.
**Validation:** Inferred during capture.

## Open Questions
- [ ] I discovered this out-of-scope task while working on REQ-374: the done line's faint text measures about 3.3:1 against the page background in both themes, under 4.5:1 for 11px text. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.
