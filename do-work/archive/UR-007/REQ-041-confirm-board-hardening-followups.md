---
id: REQ-041
title: "Confirm: three small board/pipeline hardening follow-ups from REQ-034"
status: completed
created_at: 2026-07-28T23:18:24Z
completed_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-034
builder_decided: true
domain: general
prime_files: []
review_generated: false
maintenance: false
---

# Confirm: three small board/pipeline hardening follow-ups from REQ-034

## What

While building REQ-034 (board overlap badges), the builder and reviewer surfaced four small hardening ideas that were outside the declared scope. None is a defect — the feature works and is tested — so each is a "should we also do this" question rather than something already built. Approve any subset; approved items become one small pending REQ each (or one combined REQ).

## What the Builder Chose

Nothing — all were recorded as discovered tasks instead of being fixed inline, per the out-of-scope contract.

## What Would Change

Each approved item is a small, independent change; declining any leaves today's working behavior untouched.

## Open Questions

- [x] Add a `generate_test.go` substring assertion covering the badge render path (`badge-write-overlap` / `writeSetOverlaps` in the inlined board JS)? → Confirmed: Yes — queued as REQ-051
  Recommended: Yes — cheapest real coverage; the badge's entire render path currently has zero automated tests (the Go tests cover only the annotation).
  Value: A refactor that drops the badge renderer fails a test instead of shipping a silent regression.
  Risk: Near zero — one assertion in an existing test file, same style as its neighbors.
  Also: No — accept that frontend rendering stays covered only by manual board checks.
- [x] Add a contract-regression ratchet pinning the display-only invariant (`annotateWriteSetOverlap` called after `bucketColumns`, plus the display-only wording in `actions/board.md`)? → Confirmed: Yes — queued as REQ-052
  Recommended: Yes — the invariant is currently protected by one Go test and prose; a ratchet also guards the instruction-side wording.
  Value: Prevents a future edit from quietly turning the display annotation into column logic or deleting the doc claim.
  Risk: Low — ratchets are string anchors; REQ-033's review showed weak anchors can be gutted, so anchor on the call-site line and heading.
  Also: No — the existing `TestWriteSetOverlapNeverAffectsColumnPlacement` is deemed sufficient.
- [x] Teach the pipeline's Lessons-capture step to honor a prime file's inline-only marker (so it stops appending archive links that die in consumer installs, e.g. `tools/queue-kanban/prime-do-kanban.md`)? → Confirmed: Yes — queued as REQ-053
  Recommended: Yes — the pipeline re-introduces a dead link the next time any REQ lists that prime; fix the flow, not just the instance.
  Value: Consumer installs stop accumulating dead lesson links; the prime's own header contract becomes machine-honored.
  Risk: Low — a conditional in `actions/work.md` Step 8's prime-link write (and `actions/review-work.md`'s twin); wording must not disturb the normal link path.
  Also: No — treat it as this one prime's quirk and keep inlining by hand when noticed.
- [x] Create a user-facing `docs/board-guide.md` (board features are currently documented only in agent-facing `actions/board.md`)? → Yes, write it now (user overrode the builder's "No" recommendation) — queued as REQ-054
  Recommended: No — a new user guide is real ongoing doc surface; the board UI is largely self-explanatory and `actions/board.md` already serves agents.
  Value: A linkable feature tour for humans (columns, badges, testing view).
  Risk: Ongoing drift surface — every board change would need a third doc updated.
  Also: Yes — write a short guide now while the feature set is fresh.

## Implementation

**No changes needed in this REQ.** All four questions resolved by the user via `do-work clarify` on 2026-07-29. Items 1–3 confirmed as recommended; item 4 (board guide) approved against the builder's "No" recommendation — `docs/` already carries a per-action guide convention it slots into. Approved items queued as REQ-051 (badge render test), REQ-052 (display-only ratchet), REQ-053 (inline-only lessons marker), REQ-054 (board guide).

*Resolved via clarify questions*
