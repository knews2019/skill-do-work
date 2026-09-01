---
title: "Lessons from REQ-191: Extract an explicit standalone present-video action"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-191-extract-an-explicit-standalone-present-v.md]
related:
  - page: concept-completed-work-presentation
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-191: Extract an explicit standalone present-video action

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Move the existing Remotion specification into `do-work-toolbox present-video [UR|REQ|most recent]` and a dedicated guide. Generate a valid, source-only animated walkthrough only after an explicit video request, with no automatic invocation and no MP4 render path.

## Solution summary

Added the explicit source-only walkthrough action and guide with shared completed-work evidence delegation, proportional verified scenes, successful skip behavior, immutable collision-safe output, and manual foreground-only preview guidance. No tests were touched because REQ-192 owns durable routing and contract-test integration.

## What worked

- Recovering the removed Remotion section from Git history separated its useful source/package contract from the unsafe preview wrapper without restoring obsolete `present-work` modes.
- Treating direct action dispatch as explicit authority kept `$ARGUMENTS` target-only and made the never-automatic boundary auditable.

## What didn't work

- The first action draft delegated publication canonically and then restated the collision branch in later steps and checks; REQ-201 now sweeps that duplicated-rule class.
- Qualification's tracked-reference heuristic warned on the new guide because both linked files were untracked; full-file inspection and the shipped-reference suite were needed to judge the actual link.

## Back-reference

See `do-work/archive/UR-042/REQ-191-extract-explicit-present-video-action.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c5d040a`.
