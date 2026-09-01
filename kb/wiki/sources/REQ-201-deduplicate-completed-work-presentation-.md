---
title: "Lessons from REQ-201: Deduplicate completed-work presentation publication mechanics"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-201-deduplicate-completed-work-presentation-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-201: Deduplicate completed-work presentation publication mechanics

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Keep completed-work presentation consumers limited to their consumer-specific preferred output and delegate the collision/no-overwrite algorithm and verification to `actions/completed-work-presentation-reference.md`. Sweep every current consumer of that reference so this duplicated contract class cannot recur as parallel wording.

## Solution summary

Made Collision-Safe Publication consumer-neutral and complete, moved generic immutable-output verification into that section, reduced both live consumers to preferred-path/content/result concerns, and added narrow source-seam/restatement contracts.

## What worked

- A consumer sweep found duplication beyond the originally reported video suffix branch and centralized generic verification for both item-level actions.
- Narrow restatement patterns preserved legitimate preferred names, result summaries, and the separate generated-image transaction.

## What didn't work

- The first ratchet accepted a passive checklist heading as an active application pointer and missed a paraphrased whole-artifact path algorithm.

## Worth knowing

Canonicalization tests need to prove an active directive at the execution boundary, not the mere presence of a section name anywhere. Include paraphrase mutations for the rule's meaning, not only copied keywords.

## Back-reference

See `do-work/archive/UR-042/REQ-201-deduplicate-completed-work-publication-mechanics.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `54da281`.
