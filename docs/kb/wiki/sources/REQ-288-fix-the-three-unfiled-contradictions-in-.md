---
title: "Lessons from REQ-288: Fix the three unfiled contradictions in clarify's Step 4"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-288-fix-the-three-unfiled-contradictions-in-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-288: Fix the three unfiled contradictions in clarify's Step 4

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

`skills/do-work/actions/clarify.md` Step 4 carries three shipped contradictions. Each is a pair of
statements that cannot both be followed, and none is filed anywhere. Fix all three in one pass.

## Solution summary

Fixed all three contradictions in one pass. The durable record now includes its dated reasoning note and cites the Timestamp rule for the date; that rule's date-only paragraph was rewritten to key on the condition instead of two named consumers. Per-question verbs no longer set REQ-level status — Step 5 computes it once from every question's outcome, with any remaining open question holding the REQ in the queue. The `completed` fast path routes on the existing `builder_decided: true` marker, which gained the Schema Read Contract row it never had, and the literal-string defense that existed only to protect the old prose-keyed routing was deleted.

## What worked

Checking whether the marker K4 asked for already existed before adding one. It did — `builder_decided: true`, stamped by exactly the right emitter and absent from exactly the REQs that must not take that branch — so the fix collapsed from "new field, second emitter, contract row, migration" to "route on this, add the missing row". `maintenance.md`'s deletion questions are what prompted the check; the REQ itself said "add a marker" and would have been followed literally otherwise.

## What didn't work

The mutation runs kept hitting the command timeout, because `contract-regressions.sh` runs the whole prescribed-shell suite before its own assertions and takes minutes. Two mutations were lost to timeouts before switching to extracting the block and running the two `grep` patterns directly against an original and a mutated copy — same assertion, seconds instead of minutes. For a prose-assertion check, testing the pattern against two copies of the file is equivalent to running the suite and is the right loop to iterate in.

## Worth knowing

K4's failure mode is worth remembering as a shape, not just as a bug. `review-work.md` had *documented* it — it knew the routing was prose-keyed, knew a rewording would misroute an approved follow-up into being archived unbuilt, and defended it by requiring an exact sentence in every consent question it wrote. A rule that requires every future author to reproduce a magic string is a defense with a per-use failure probability. The marker moves the decision to something the emitter sets once and nobody retypes.

## Back-reference

See `do-work/archive/UR-059/REQ-288-fix-the-three-unfiled-contradictions-in-clarifys-step-4.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c25ee71`.
