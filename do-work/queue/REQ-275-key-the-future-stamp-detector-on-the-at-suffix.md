---
id: REQ-275
title: Key the board's future-stamp detector on the _at suffix instead of six hand-kept names
status: pending-answers
created_at: 2026-08-18T23:38:35Z
status_changed_at: 2026-08-18T23:38:35Z
user_request: UR-056
addendum_to: REQ-267
domain: general
review_generated: true
sweep: true
sweep_key: future-stamp-detector-field-list-closed
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work-board/tools/queue-kanban/model.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
---

# Key the Board's Future-Stamp Detector on the `_at` Suffix Instead of Six Hand-Kept Names

## What

A **third** divergence family between the timestamp repairer and the board's readers, on an axis neither of the two fuzzes covered — the **field name**.

The repairer repairs any top-level field whose name ends in `_at`, by suffix rule, and its header says so explicitly: *"The rule is the `_at` SUFFIX, not a list of field names: a schema that grows a new stamp field is covered the day it is added, with nothing to remember."* The board's `detectFutureTimestampFields` (`skills/do-work-board/tools/queue-kanban/model.go:1374-1384`) checks a **fixed six-name list**.

Confirmed by execution: a file carrying `reviewed_at: 2093-01-01T00:00:00Z` is rewritten unattended by the repairer while the board never badges it.

**Latent, not live.** No field in today's schema is missed, and no corruption occurs — the derivation is unchanged. But this is the one axis where the unattended writer is provably outside the read side's envelope, and it bites silently the day the schema grows an `_at` field.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Requirements

- `detectFutureTimestampFields` keys on the **condition** (a top-level frontmatter key ending in `_at`) rather than an enumeration, matching the repairer's own stated rule — `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale.
- Consider deliberately whether any currently-listed name must stay special-cased, and say so if it must; a silent widening that changes which cards get badged is a user-visible change, not a refactor.
- A lock-in that fails if the detector goes back to an enumeration — for example, a probe asserting a novel `*_at` field is detected.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** a REQ file carrying `reviewed_at` set beyond the skew horizon is rewritten by `skills/do-work/scripts/repair-req-timestamps.sh` and produces **no** future-stamp warning from the board. Observed today.

**GREEN:** the board warns on it, and the repairer's and the board's field sets agree by construction rather than by coincidence.

## Context

REQ-267's independent review, Important finding 2 (gate: rule-change). Created `pending-answers` per the generation-≥2 cascade stop.

Deliberately **not** folded into REQ-274 even though both came from the same review: different root cause (a stale enumeration versus a false stated mechanism), different file, different fix. Two sweeps, not one.

Worth knowing: **both fuzzes missed this**, the builder's and the reviewer's, because both varied one field name. A fuzz's blind spots are the axes it holds constant, and that is the second time in this REQ that the *oracle* rather than the fix was the limiting factor.

## Open Questions

- [ ] REQ-267's review found a third repairer/board divergence on an axis neither fuzz covered: the repairer repairs any `_at` field by suffix while the board checks six hand-kept names, so a future `_at` field would be rewritten unattended and never badged. Latent today. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — nothing in today's schema is missed, so accept the divergence until a new `_at` field is added.
