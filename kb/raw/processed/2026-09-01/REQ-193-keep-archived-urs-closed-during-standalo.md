---
source_type: req_lesson
req_id: REQ-193
req_path: do-work/archive/UR-043/REQ-193-keep-archived-urs-closed-during-review.md
date: 2026-08-15
domain: general
module: _dev/primes
tags: [general, keep, archived, closed, during]
---

# Lessons from REQ-193: Keep archived URs closed during standalone review

## What the REQ was about

Co-locate the archived-UR invariant with standalone review's archived-input read so a reviewer can create a follow-up without moving or reopening the already-closed User Request folder. Make the downstream archive instruction explicit enough that the completed follow-up returns to that folder in place.

## Solution summary

**[MAP CHANGED]** Standalone review now has a closed-UR lifecycle from context read through completed follow-up placement: same-UR fixes queue normally, but the archived UR never reopens or moves, and the finished review-generated REQ returns to that existing folder in place.

## What worked

- Keying the authority rule to use of the archived fallback, rather than review mode, covers both standalone review and later orchestrated review of its generated fix.
- A literal full-predicate contract locks the marker, conjunction, archived-folder existence, and same-UR relation without adding another scanner.

## What didn't work

- The first wording treated orchestrated review as universally open-UR, overlooking review-generated work whose UR deliberately stays archived.
- The first broad regex required the archive path token but let deletion or negation of `already exists` survive.

## Worth knowing

The in-place archive override is legitimate only when both `review_generated: true` and the matching archived UR folder already exist. Generic live REQs under closed URs remain anomalies for REQ-194's detector path.

## Back-reference

See `do-work/archive/UR-043/REQ-193-keep-archived-urs-closed-during-review.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6fcc433`.
