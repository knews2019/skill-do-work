---
title: "Lessons from REQ-281: Reconcile the calibration log against the frontmatter it was derived from"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-281-reconcile-the-calibration-log-against-th.md]
related:
  - page: REQ-280-probe-timestamp-ordering-and-point-check
    rel: depends-on
  - page: REQ-284-emit-every-verify-finding-from-the-board
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-281: Reconcile the calibration log against the frontmatter it was derived from

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

`do-work/calibration-log.tsv` is an independent third record of every REQ's wall span, written once at archive time by `actions/work.md` Step 8 substep 7.5 as `completed_at − claimed_at`, and read back by `actions/estimate-reference.md:94` as the corpus the scoring table is fit from. Nothing ever compares it against the frontmatter it came from.

Add a `queue-kanban verify` probe that recomputes each row's `wall_minutes` from its REQ's `claimed_at` and `completed_at` and reports rows that disagree by more than a minute.

## Solution summary

Added a read-only probe to `queue-kanban verify` that recomputes every `do-work/calibration-log.tsv` row's `wall_minutes` from its REQ's `claimed_at`/`completed_at` and reports disagreements beyond a one-minute tolerance, naming both values and stating that either record may be the correct one. Rows that cannot be reconciled at all are a separate finding, and an absent log is a reported skipped probe. Nothing is marked fixable and nothing is repaired. Added the probe to forensics Check 14's table, and corrected the single log row this session had written wrong.

## What worked

Writing a throwaway Python recomputation *before* the Go probe, then comparing the two on 72 real rows. It confirmed the REQ's measurement independently, and later it was the acceptance evidence — a fixture proves a probe handles the case you imagined, while two independent implementations agreeing on real data proves it handles the ones you did not. Also: checking whether the new failing condition breaks the canonical gate *before* claiming the REQ was done. It did not, but that was luck, and finding out at review time rather than at commit time was not.

## What didn't work

Nothing failed outright, but the probe immediately flagged a row **this session had written wrong four hours earlier** — REQ-274's, logged as 7 against a true span of 5, because the calibration arithmetic used a hardcoded `claimed_at` string instead of reading back the one actually stamped into the file. The REQ was built to catch exactly that class and caught its own author. The habit worth taking: when a step writes a derived value from a stamp, read the stamp back from the file it was written to rather than reusing the variable.

## Worth knowing

REQ-241, REQ-243 and REQ-245 log three different spans against one identical `claimed_at` of `2026-08-18T12:43:06Z`. That pattern reads like a fan-out wave whose members all recorded the wave's dispatch instant — which would mean the *frontmatter* is the wrong record for those three, not the log. If that is confirmed in REQ-311, REQ-280's ordering probe is currently reading three REQs' spans wrong too. Nobody should batch-rewrite this file before that question is settled.

## Back-reference

See `do-work/archive/UR-057/REQ-281-reconcile-the-calibration-log-against-req-frontmatter.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a868827`.
