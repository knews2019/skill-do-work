---
title: "Lessons from REQ-172: Make screenshot source cleanup best-effort"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-172-make-screenshot-source-cleanup-best-effo.md]
related:
  - page: REQ-173-handle-first-line-bom-in-just-collision-
    rel: complements
  - page: REQ-174-validate-root-markdown-fence-info
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-172: Make screenshot source cleanup best-effort

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

After a screenshot has been copied, byte-verified, and installed without clobbering, a failure to remove the staged source must warn without invalidating the permanent asset.

## Solution summary

Staged-source removal now warns without invalidating an already verified permanent screenshot. The executable lifecycle regression forces only that removal to fail and proves the source, destination, temporary-copy cleanup, and later no-clobber collision behavior.

## What worked

Replaying only the staged-source `rm` separated post-install cleanup from the strict installation boundary.

## What didn't work

Treating every cleanup failure as transactional failure created a state that the normal retry could not repair.

## Worth knowing

Once the permanent asset is byte-verified and no-clobber installed, cleanup warnings must not revoke its validity.

## Back-reference

See `do-work/archive/UR-039/REQ-172-make-screenshot-source-cleanup-best-effort.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9bf5a19`.
