---
source_type: req_lesson
req_id: REQ-477
req_path: do-work/archive/UR-088/REQ-477-family-keyed-lessons-index-and-trap-promotion.md
date: 2026-09-01
domain: general
module: _dev/primes
tags: [general, family, keyed, lessons, intelligent]
---

# Lessons from REQ-477: Family-keyed lessons, intelligent index, and mandatory Trap promotion

## What the REQ was about

Make the Lessons-Capture Phase produce transferable context: every appended lesson bullet carries a failure-family slug, an intelligent index over the lesson satellites is created/refreshed in the same edit, and a second same-family occurrence makes Trap promotion mandatory instead of a judgment call.

## Solution summary

Lesson writers now emit a literal family discriminator, refresh one reproducible index row in the same edit, and promote a generalized trap by the second occurrence. The initial index is complete for tracked source satellites and honest about partial legacy coverage.

## What worked

**What worked:** Enumerating live writers before editing prevented the standalone review path from silently bypassing the new contract; testing the patch on a clean detached worktree separated REQ-477 evidence from paused REQ-420 changes.

**What didn't:** Running the repository contract suite directly in the shared dirty tree produced many unrelated REQ-420 failures and could not qualify this change.

**Worth knowing:** An output-format rule is incomplete until every writer is swept, even when the REQ names only the primary writer.

## Back-reference

See `do-work/archive/UR-088/REQ-477-family-keyed-lessons-index-and-trap-promotion.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `74b1fd41`.
