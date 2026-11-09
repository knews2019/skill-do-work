---
title: "Lessons from REQ-197: Normalize completed-work presentation target IDs"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-197-normalize-completed-work-presentation-ta.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-197: Normalize completed-work presentation target IDs

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Make every completed-work presentation ID path touched by UR-042 inherit `actions/work-reference.md` → **Target ID Resolution** before dispatch or archive lookup, so case-insensitive prefixes and numeric-value matching resolve canonical stored IDs such as `REQ-042` and `UR-011`.

This is a standalone user-visible input contract and cannot fold into a sweep: its fix is unrelated to output-directory publication and has one canonical resolver surface.

## Solution summary

Made both completed-work presentation ID entry paths inherit the canonical Target ID Resolution grammar before lookup or migration dispatch, preserved caller-specific search/write behavior, and added replayable assertions for equivalent input spellings and supplied-token output.

## What worked

- Exploration found the single canonical grammar and kept the product change limited to two callers plus one contract seam.
- Mutation testing exposed semantic blindness that ordinary GREEN runs could not reveal.

## What didn't work

- Attempt 1 copied examples into callers and tested for their presence, creating the drift the reference contract forbids.
- The remediation removed the duplication but still used a substring-positive assertion; “read without applying” survived, so the single remediation attempt did not fully close the review finding.

## Worth knowing

When an instruction caller inherits a canonical contract, test the source definition once and test each caller's active, ordered application with adversarial negations. Positive keyword presence alone is not evidence that the directive remains operative.

## Back-reference

See `do-work/archive/UR-042/REQ-197-normalize-completed-work-presentation-target-ids.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `89d068e`.
