---
title: "Lessons from REQ-452: Refuse ambiguous explicit request IDs"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-452-refuse-ambiguous-explicit-request-ids.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-452: Refuse ambiguous explicit request IDs

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Preserve duplicate queue-record collision evidence when resolving numeric request IDs, and return an ambiguity exclusion when an explicit `REQ-NNN` token cannot identify exactly one file. Explicit targeting may bypass documented dependency, assignment, and impact gates, but never repository identity ambiguity.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_targets_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)

## What worked

- Caller-seam discovery-order replays exposed both arbitrary selection and the missing normalized collision evidence.
- Keeping collision paths as repository-model evidence let the selector return one typed, actionable ambiguity exclusion.

## What didn't work

- The first fixture used filenames that already shared the target number, masking the frontmatter-only gap.
- The remediation reused a suffix-tolerant filename parser for frontmatter, which fixed genuine numeric equivalents but admitted malformed trailing text.

## Worth knowing

- Collision fixtures must decouple filename claims from frontmatter claims. Use unrelated filename numbers when proving frontmatter identity, and include a malformed adjacent value as a negative control.

## Back-reference

See `do-work/archive/REQ-452-refuse-ambiguous-explicit-request-ids.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bbc57391`.
