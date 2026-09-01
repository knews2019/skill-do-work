---
title: "Lessons from REQ-198: Publish generated directory only after image success"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-198-publish-generated-directory-only-after-i.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-198: Publish generated directory only after image success

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Align the ai-report image-generation shell block with the conditional bundle contract: an all-failed generation attempt must not leave an empty published `generated/` directory, while successful images still publish there with status-backed freshness.

This is a standalone user-visible artifact-shape contract and cannot fold into a sweep: its root cause is the image helper's publication timing, not target parsing or a repeated prose class.

## Solution summary

Made raster generation batch-private until status-backed success exists, conditionally published the verified batch with one adjacent rename, cleaned all-failed/interrupted staging, and added exact prescribed-block replay coverage for all-failed and mixed-success outcomes.

## What worked

- Replaying the exact fenced shell block proved the user-visible directory shape rather than relying only on prose patterns.
- Adjacent private staging and status-plus-size filtering cleanly separate current images from stale or failed targets.

## What didn't work

- Normal all-failed/mixed cases were not enough to prove the complete batch lifecycle: caller interruption can orphan children, and plain `mv` has destination-directory nesting semantics under a coordinated collision.

## Back-reference

See `do-work/archive/UR-042/REQ-198-publish-generated-directory-only-after-image-success.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `00db46c`.
