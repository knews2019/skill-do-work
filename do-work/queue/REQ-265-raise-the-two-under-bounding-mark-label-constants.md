---
id: REQ-265
title: Raise the two under-bounding mark-label constants to the current build
status: pending
created_at: 2026-08-18T20:07:08Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-051
addendum_to: REQ-252
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# Raise the Two Under-Bounding Mark-Label Constants to the Current Build

## What

Chromium 141.0.7390.37 measures the 11px mark-label line box at **12.9631** (constant `durationsMeasuredLabelBoxHeightUnits` records 12.84) and its descent at **2.7778** (constant records 2.41). Per the larger-wins convention both constants should rise (≥12.97 / ≥2.78). Nothing is live-wrong today — the pitch-floor consumer clears the real box at pitch 13, and the ceiling consumer's paired title-ascent constant over-bounds by 0.99 — but the compensation is a coincidence of one consumer, not a guarantee. When raising, re-verify `TestDurationsLastLabelRowClearsPanelBTitle`'s margin (0.12 model units at the larger descent) and expect `TestDurationsLabelRowPitchClearsTheLabelTextBox` to still pass at pitch 13.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-252's builder measured both exceedances and captured the raise as a Discovered Task per the REQ's no-value-change rule; its review (F1a, gate: trivial) verified no assertion flips on any recorded build and routed the raise here as a durable artifact. Created `pending-answers` per the generation-≥2 depth stop.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-252: two measured constants no longer bound the face on a current build. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — wait until an assertion actually flips.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.
