---
title: "Lessons from REQ-155: Review fix: Correct the manual Stop-hook object path"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-155-review-fix-correct-the-manual-stop-hook-.md]
related:
  - page: concept-queue-task-lifecycle
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-155: Review fix: Correct the manual Stop-hook object path

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Make the no-JSON-tool instruction identify individual nested Stop-hook objects exactly, so following it cannot delete a whole hooks array containing custom neighbors.

## Solution summary

Manual fallback guidance now identifies individual retired hook objects and matches automated wrapper cleanup semantics, while structural fixtures prove custom neighbors survive and guard-only wrappers do not remain empty.

## What worked

- Using the already-correct jq/Python structures as the manual contract made the wording mechanically checkable rather than interpretive.
- A guard-only wrapper fixture caught the difference between removing a wrapper and merely leaving `{"hooks": []}` behind.

## What didn't work

- The earlier exact-output test locked in a JSON path that selected an array instead of its objects; exact-string tests are only as safe as the semantics they encode.

## Worth knowing

- Manual mutation guidance should name both the narrow target and the precise condition for removing its parent, especially when custom siblings can share that parent.

## Back-reference

See `do-work/archive/UR-031/REQ-155-correct-manual-stop-hook-object-path.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `c1f8e21`.
