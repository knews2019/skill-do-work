---
source_type: req_lesson
req_id: REQ-192
req_path: do-work/archive/UR-042/REQ-192-migrate-presentation-routing-docs-and-contracts.md
date: 2026-08-15
domain: testing
module: _dev/primes
tags: [testing, migrate, completed, work, presentation]
---

# Lessons from REQ-192: Migrate completed-work presentation routing documentation and contracts

## What the REQ was about

Update toolbox routing, discovery, completion-flow recommendations, cross-references, tests, and release notes so the suite presents one unambiguous completed-work choice: detailed report through `ai-report`, portfolio through `present-work`, and animated walkthrough through `present-video`.

## Solution summary

Migrated live routing, discovery, full-cycle guidance, caller descriptions, and portfolio retention documentation to one explicit three-command ownership model. Added RED-first durable contracts for routing, inventories, report evidence modes, portfolio branches, source-only video, archive safety/status, and retired-workflow rejection while preserving the four focused follow-up defects for their owning REQs.

## What worked

- A condition-based surface sweep reduced a broad migration to 20 live files while leaving accurate history, generated artifacts, presentation mechanics, and focused review follow-ups untouched.
- Adding the routing and inventory tests before product edits produced a useful RED signal for every stale command family and then verified the exact three-owner model.

## What didn't work

- The first GREEN contract block under-specified load order, retired portfolio workflows, and unsafe guide commands. Qualification caught those gaps and sent the test seam back for focused correction.
- The widened unsafe-preview matcher still encoded literal examples too narrowly; review mutation probes exposed fixed-port flag and nonliteral platform-opener escapes, now owned by REQ-202.

## Worth knowing

Presentation routing aliases must remain exact and mutually exclusive, while guardrail applicability should be condition-based with action names treated only as examples. Test executable command segments separately from negative explanatory prose so safety assertions can be both broad and precise.

## Back-reference

See `do-work/archive/UR-042/REQ-192-migrate-presentation-routing-docs-and-contracts.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a00ee67`.
