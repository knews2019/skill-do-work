---
id: REQ-239
title: Give the Timeline's rows a real focus ring
status: pending-answers
created_at: 2026-08-18T11:09:44Z
user_request: UR-051
addendum_to: REQ-233
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Give the Timeline's Rows a Real Focus Ring

## What

`web/board.css` sets `.timeline-row { outline: none; }`, so a focused Timeline row falls back to `.timeline-row:focus .timeline-row-hit { fill: var(--surface-2); }` — a one-step background tint. Every other focusable thing on the board gets a 2px `--accent-claimed` ring. Give the rows the same ring, and key it on `:focus-visible` so a pointer click does not draw one.

## Context

Found by REQ-233's review. REQ-233 added a visible focus ring to the *chart container* because its requirement said "focus is visible on whatever element takes the keyboard interaction — a focus ring that exists only in the default user-agent style is not enough on a dark surface". The rows next door fail the same test for a different reason: not a user-agent default, but an explicit `outline: none` with a weak substitute.

This is a sibling of REQ-233's requirement rather than part of it — rows are not what REQ-233 added, and rows were already keyboard-activatable before it (pinned by `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard`). It is recorded separately so the fix gets its own before/after rather than riding in on another REQ's merge.

The chart container's ring is the model to copy, including its reasoning: it uses `--accent-claimed`, the token every other ring on the board uses, and `outline-offset: -2px` because the container is flush under the axis. A row's correct offset may differ — a row is not clipped the same way — so this is a judgment, not a copy.

## Requirements

- A keyboard-focused Timeline row draws a focus indicator of the same weight as the rest of the board's rings, using the same token.
- It keys on `:focus-visible`, so a pointer click does not draw one — matching `.control-button:focus-visible`, `.req-card:focus-visible`, and `.calendar-chip:focus-visible`.
- The ring is not clipped by the row's own geometry or by the scroll container.
- The existing `.timeline-row:focus .timeline-row-hit` tint either stays as a complement or is removed deliberately, with the choice stated — two overlapping focus signals is a decision, not an accident.
- No change to row activation behaviour; `TestJavaScriptBehaviorTimelineRowsActivateFromTheKeyboard` still passes.

## Red-Green Proof

**RED prompt/case:** an assertion over the generated stylesheet that `.timeline-row` carries a `:focus-visible` rule with a non-`none` outline, in the same shape the check would use for `.control-button:focus-visible`.
**Why RED now:** `.timeline-row` sets `outline: none` and has no `:focus-visible` rule at all.
**GREEN when:** the assertion passes and a real Tab press onto a row draws a visible, unclipped ring — verified in a browser, not by a programmatic `.focus()`, which does not trigger `:focus-visible` and will report a false negative.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] While reviewing REQ-233 I found that the Timeline's rows have their focus outline explicitly switched off, and get only a faint background tint instead — one shade different from the row's normal colour. Everything else on the board that can be focused gets a clear 2px coloured ring, which is what REQ-233 just added to the chart itself. So a keyboard user moving down the rows has a much weaker sense of where they are than anywhere else on the board. Nothing is broken and the rows still work; this is about how visible the current position is. The fix is a few lines of CSS plus one assertion. I am asking rather than doing it because the tint was written deliberately — someone chose to turn the outline off — and I cannot tell from the code whether that was to avoid a clipped or ugly ring on a dense chart, which is a real concern on rows that are only a few pixels tall. If it was, the answer might be a better tint rather than a ring. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the tint is the deliberate choice for dense rows and should stay.
