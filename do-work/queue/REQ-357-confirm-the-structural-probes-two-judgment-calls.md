---
id: REQ-357
title: "[impact-negligible] Confirm: the structural probe's two judgment calls"
status: pending-answers
created_at: 2026-08-24T09:50:00Z
user_request: UR-068
addendum_to: REQ-343
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Confirm: The Structural Probe's Two Judgment Calls

## What

REQ-343 shipped two decisions the builder made on its own judgment because no answer was available
mid-build. Both are live in `verify.go` now. Confirm or overturn them.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Neither is a defect and neither blocks anything — REQ-343 is complete and green. They are recorded
here because a builder's judgment call that nobody ever confirms becomes a rule by default, and both
of these are cheap to reverse now and awkward to reverse once operators are used to the output.

## Open Questions

- [ ] Should verify report an empty `id:` field at all?
  `deriveRequestIdFromFilename` means the board never loses the REQ over a missing `id:`, so the real
  exposure is narrower than the other damage classes: a file rename silently renumbers the REQ.
  REQ-343 flags it as damage and says exactly that in the finding's detail text.
  Recommended: Keep it — the rename hazard is real and the detail text already scopes the claim.
  Also: Drop the `id` branch (the other three classes are unaffected); or keep it but downgrade the
  wording so it reads as a caution rather than damage.

- [ ] Should the `archive/legacy/` carve-out key on the directory, or on a `created_at` cutoff?
  REQ-343 uses the directory, because that is what its `## Context` names and it needs no date
  arithmetic. The consequence: a REQ written today and dropped into `archive/legacy/` is exempt from
  the `user_request` check.
  Recommended: Keep the directory — a structural probe with a clock dependency is worse, and the
  exemption is visible the moment anyone looks at the directory.
  Also: Add a `created_at` cutoff alongside the directory; or drop the carve-out and backfill
  `user_request` into the 11 legacy REQs instead.

## Context

Both decisions are recorded as D-07 and D-08 in the archived REQ-343, with the builder's Value and
Risk lines. The code they describe is live: `appendStructuralDamageFindings` and
`isLegacyArchiveRequestPath` in `skills/do-work-board/tools/queue-kanban/verify.go`.

`status: pending-answers`, so the work loop walks past it until `do-work clarify` resolves the two
questions. Answering may need no code change at all — "keep both" closes this REQ with a note.

---
*Source: builder-decided questions from REQ-343 (UR-068).*
