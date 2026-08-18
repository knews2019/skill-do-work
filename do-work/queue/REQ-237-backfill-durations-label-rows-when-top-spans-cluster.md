---
id: REQ-237
title: Backfill the Durations label rows when the longest spans cluster
status: pending
status_changed_at: 2026-08-18T10:46:58Z
created_at: 2026-08-18T10:52:00Z
user_request: UR-051
addendum_to: REQ-231
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
effort_estimate: normal
write_set:
- skills/do-work-board/tools/queue-kanban/durations.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Backfill the Durations Label Rows When the Longest Spans Cluster

## What

In the Durations view's overflow lane, the board picks the six longest spans in a band as label candidates, then walks them left to right and gives each the first text row where it does not touch a label already placed. A candidate that fits nowhere is simply dropped and counted. Nothing then offers the freed space to the seventh-longest span, so where several of the longest spans finish close together in time, the lane's two text rows end up mostly empty while the remainder count carries almost everything.

## Context

Found while reviewing REQ-231, which introduced the six-longest selection rule (the maintainer's chosen "Alternative 2"). Measured on two synthetic 60-sample boards:

- **Magnitude correlated with completion time** (every long span crowded at the right edge): **2 labels drawn out of 6 candidates**, 58 in the remainder.
- **Magnitude scattered across the window**: **5 of 6 place**, which is the healthy case.

So this is the tail, not the norm — but the tail is exactly the shape a burst of long REQs produces, which is also when a reader most wants the lane to talk. Before REQ-231 the lane filled both rows on the same dense fixture (27 labels), though half of them were unreadable under the dots, which is the defect REQ-231 fixed. The two label rows are a fixed budget either way; the question is only whether an unusable candidate's slot should pass to the next-longest span.

The change would be local: `selectDurationLabelCandidates` currently returns a fixed set of six before placement runs, so placement has no way to ask for a replacement. Making the two cooperate — placement pulling the next candidate when one is dropped — is a real change in shape, not a constant tweak, which is why it is a question rather than a fix.

## Requirements

- On a band where selected candidates collide, the label rows carry as many of the band's longest spans as physically fit, rather than stopping at the first six by magnitude.
- Every drawn label is still one of the band's longer spans — backfill may not reintroduce the left-edge first-fit sampling REQ-231 removed.
- Labelled + hidden still equals the band's sample count, so nothing is silently dropped.
- No change to the payload's shape (`labelRow` / `labelAnchor` / per-band hidden counts).

## Red-Green Proof

**RED prompt/case:** a test on a magnitude-gradient fixture (long spans crowded at one edge) asserting that the number of drawn labels equals what the two rows can physically hold, not the number of top-N candidates that happened to fit.
**Why RED now:** measured at 2 labels of a possible ~13 row-slots on exactly that fixture.
**GREEN when:** the same test passes, `TestOverflowLabelsGoToTheLongestSpans` still passes unchanged, and re-rendering the gradient fixture shows both rows populated with long spans.
**Validation:** apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-231: the Durations chart's top strip picks the six longest jobs to label, but if several of those six finished at nearly the same time they cannot all fit side by side, and the ones that do not fit are just dropped — nothing offers their space to the seventh-longest job instead. On a test board where the longest jobs all clustered, that left 2 labels drawn where the strip had room for about 13; on a board where the long jobs were spread out, 5 of 6 placed fine. So it only bites when a burst of slow work lands together, which is arguably when you most want to read the strip. Fixing it means letting the placement step ask the selection step for a replacement whenever it drops one, which is a genuine change to how the two cooperate rather than a tuning knob — and you may reasonably prefer the current simpler rule, where the remainder count carries the overflow and the strip stays predictable. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the remainder count already states what is not shown, and a two-step selection is more machinery than the lane is worth.
