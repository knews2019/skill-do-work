---
title: "Lessons from REQ-479: Enforce required-lessons reads and audit un-promoted families"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-479-enforce-required-lessons-reads-and-audit.md]
related:
  - page: REQ-477-family-keyed-lessons-intelligent-index-a
    rel: complements
  - page: REQ-478-capture-stamps-required-lessons-under-a-
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-479: Enforce required-lessons reads and audit un-promoted families

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

Make the work pipeline read every stamped lessons reference before implementation, and extend `prime audit` to catch missed promotions and index drift.

## Solution summary

Every claimed REQ now re-evaluates the current lessons index before implementation, so pre-existing and serially captured work can inherit lessons created after capture. Builders read captured and claim-time pointers unconditionally while retaining the broader touched-prime rule, and prime audit can identify every accepted promotion/index drift class.

## What worked

Treating claim as a second projection point closes the exact time-order gap capture cannot see, while reusing the same named budget contract avoids a competing rule.

## What didn't work

Nesting the consult under the original “Routes B and C” exploration heading initially excluded Route A; checking the trigger against every route exposed it before release.

## Worth knowing

Context routing that runs at capture only is stale by construction for serial batches; enforce it again at claim, and make that consult independent of whether the route performs exploration.

## Back-reference

See `do-work/archive/UR-088/REQ-479-enforce-required-lessons-reads-and-audit-unpromoted-families.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9a1b7bfb`.
