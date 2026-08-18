---
id: REQ-241
title: Reconcile the Durations label metrics with the face actually rendered
status: pending
status_changed_at: 2026-08-18T12:17:03Z
created_at: 2026-08-18T12:09:46Z
user_request: UR-051
addendum_to: REQ-237
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: durations-label-metric-constants-disagree-with-rendered-face
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
---

# Reconcile the Durations Label Metrics With the Face Actually Rendered

## What

Two constants that describe the Durations label face disagree with what the browser actually draws. Neither causes a visible collision today, but both are now load-bearing in a way they were not before, because REQ-237 made the label rows actually fill up.

## Context

Found by REQ-237's build and confirmed independently against the merged tree. Both are pre-existing; REQ-237 is what made them reachable, not what caused them.

Until REQ-237, the overflow lane drew two or three labels on any real board, so the slack in these constants never mattered. It now draws as many as fit — measured 21 on a clustered 60-sample fixture — and the slack is what stands between "packed" and "overlapping".

## Instances

- [ ] **`durationsLabelCharacterWidthUnits = 6.2` under-estimates the 11px sans face by ~7%.** Its comment calls the value "deliberately generous"; measurement says otherwise — a 14-character label renders **92.52 user units**, i.e. **6.61 units/char**, not the 86.8 the constant predicts. The 6-unit separation rule absorbs the difference, so nothing collides: the tightest same-row gap in a full 21-label render is 3.08 units, and a direct DOM measurement finds **0 same-row overlaps**. But the real margin is about half what the code claims, and the comment is actively misleading about which direction the error runs.

- [ ] **`DURATIONS_LABEL_ROW_HEIGHT = 12` is smaller than the text box the same file declares.** `DURATIONS_LABEL_TEXT_ASCENT = 11` plus a 2-unit descent is a 13-unit box on a 12-unit pitch; the rendered font box measures 12.83 units. Measured on a densely-populated lane: **20 cross-row bounding-box intersections, each 1.6px deep.** This is line-box padding rather than ink — the render shows two cleanly separated rows, and a screenshot confirms it — but it means no test can honestly assert row-against-row separation the way `TestDurationsLabelRowsClearTheMarkBands` asserts row-against-mark.

## Requirements

- Each constant either matches the face the browser renders, or its comment states the measured value and why the code deliberately differs. A constant whose comment claims a safety margin in the wrong direction is the specific defect here.
- Whatever changes, the same-row separation guarantee holds: **0 same-row label overlaps** at full density, measured from the live DOM rather than computed from the constants under test.
- REQ-231's guarantee holds unchanged: **0 label/mark overlaps** in either band, at any density.
- If the row pitch changes, Panels B and C shift with it and `describeAtPointer` still resolves the same panel for the same pointer position — the same constraint REQ-231 worked under.
- Changing label counts across the view is an accepted consequence of retuning the width model; say so in the REQ trail with before/after counts on a real board rather than only on a fixture.

## Red-Green Proof

**RED prompt/case:** a test asserting each constant against the measured face — that `durationsLabelCharacterWidthUnits` is not below the rendered units-per-character, and that `DURATIONS_LABEL_ROW_HEIGHT` is not below the declared ascent-plus-descent box.
**Why RED now:** 6.2 < 6.61, and 12 < 13.
**GREEN when:** the assertions pass, and a rendered clustered fixture still shows 0 same-row label overlaps and 0 label/mark overlaps measured in the live DOM.
**Validation:** Review finding on REQ-237; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] While reviewing REQ-237 I found two numbers in the Durations chart that describe the label text and are both slightly wrong about it. One says each character is 6.2 units wide when the browser draws 6.61 — and its comment claims the estimate is deliberately generous, which is backwards. The other says a row of label text is 12 units tall when the same file elsewhere describes that text as a 13-unit box. Nothing overlaps on screen today: I measured a fully packed lane and found no two labels on the same row touching, and no label touching a dot. But the spare room absorbing both errors is about half what the code claims it is, and until this session the chart only ever drew two or three labels so it never mattered. It does now — REQ-237 made the rows fill up. Fixing the width number will change how many labels the chart draws on every board, including yours, which is why it is your call and not a quiet tidy-up: it is a visible change to a view you look at, made for correctness reasons rather than because anything is broken. The alternative is to leave the numbers and fix the comments so they stop claiming a margin that is not there. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No — just correct the misleading comments and leave the numbers alone, since nothing overlaps.
  → Confirmed: Yes, add to queue — full fix, retune the constants to match the measured face (user accepted that label counts across the view will visibly change). [2026-08-18, via do-work clarify]
