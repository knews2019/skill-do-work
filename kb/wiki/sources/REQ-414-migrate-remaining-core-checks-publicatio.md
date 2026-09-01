---
title: "Lessons from REQ-414: Migrate remaining core checks, publication helpers, Git helpers, and surveys"
type: source-summary
topic_cluster: shell-and-automation
sources: [raw/processed/2026-09-01/REQ-414-migrate-remaining-core-checks-publicatio.md]
related:
  - page: concept-prescribed-shell-commands
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-414: Migrate remaining core checks, publication helpers, Git helpers, and surveys

Part of the [[concept-prescribed-shell-commands]] cluster.

## What the REQ was about

Move all remaining core utility domain logic into `do-work-cli` subcommands.

## Solution summary

- **Implementation commit:** `0689970c`
- **Lifecycle result:** `completed-with-issues` after one remediation and one fresh re-review.
- **Residuals:** REQ-462 through REQ-465 retain inventory, reservation-authority, structured-projection, and differential-parity closure work.

## What worked

- Registration, renderer agreement, and broad green gates do not prove a migration preserved behavior; differential fixtures must compare exact retained statuses, ordered facts, paths, actions, and side effects across combined-state and authority-boundary cases.
- Destructive cleanup needs fresh authorization as well as object identity: revalidate the evidence that permits deletion and every eligibility predicate immediately before the mutation.

## Back-reference

See `do-work/archive/REQ-414-migrate-core-helpers.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0689970c`.
