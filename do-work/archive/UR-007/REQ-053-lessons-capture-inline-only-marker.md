---
id: REQ-053
title: Teach the pipeline's Lessons-capture step to honor a prime file's inline-only marker
status: completed
commit: 8b4dd30
status_changed_at: 2026-07-29T14:02:00Z
domain: general
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
depends_on: []
write_set: ["actions/work.md", "actions/review-work.md"]
---

# Teach the pipeline's Lessons-capture step to honor a prime file's inline-only marker

## What

Make the Lessons-capture prime-link write conditional on the target prime file's inline-only marker: when a prime declares its lessons must be inlined (e.g. `tools/queue-kanban/prime-do-kanban.md`), the pipeline inlines instead of appending an archive link. Touch both write sites — `actions/work.md` Step 8's prime-link write and `actions/review-work.md`'s twin.

## Why

The pipeline currently re-introduces a dead link the next time any REQ lists an inline-only prime: archive links die in consumer installs (the archive isn't shipped), so consumer copies accumulate dead lesson links. Fix the flow, not the instance — the prime's own header contract becomes machine-honored. Approved by the user via `do-work clarify` on 2026-07-29 (follow-up discovered during REQ-034, surfaced in REQ-041).

## Constraints

- Wording must not disturb the normal (non-marked) link path.
- Per the Closed Enumerations rule, key off the marker condition in the prime's header — never a hand-list of marked primes.

## Acceptance

- A REQ listing an inline-only prime gets its lesson inlined, not linked; a REQ listing a normal prime behaves exactly as today.

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified) — Step 8 substep 7 (execute deferred prime-link writes) now checks the target prime file's `## Lessons` section for an inline-only marker before writing; branches to a plain inline bullet or falls through to the unchanged link path.
- `actions/review-work.md` (modified) — Step 9.5 item 4 (standalone-mode prime file update) gets the identical marker check and branch, ahead of its unchanged link path.

**What was done:**
Both write sites now open with a check: if the prime file's `## Lessons` section already opens with an HTML comment containing the phrase "inlined, not linked" (the exact contract text lives in `tools/queue-kanban/prime-do-kanban.md`'s `## Lessons` header), the prime has declared itself inline-only, so the step appends a plain `- REQ-NNN: 1-line summary` bullet with no link instead of computing/verifying an archive path. When the marker is absent, each site falls through to an "Otherwise" branch containing the original link-writing prose verbatim (path computation, existence-verify, and the `[REQ-NNN: ...](...)` append), so the normal path's wording and behavior are unchanged. The check reads a condition on the prime file itself rather than a hand-maintained list of inline-only primes, per the Closed Enumerations rule — any future prime that copies the same marker phrase is honored automatically.
