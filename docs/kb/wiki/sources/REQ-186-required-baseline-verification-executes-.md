---
title: "Lessons from REQ-186: Required baseline verification executes two child suites twice"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-186-required-baseline-verification-executes-.md]
related:
  - page: REQ-182-public-work-and-schema-vocabularies-drif
    rel: complements
  - page: REQ-184-live-board-origin-checks-have-no-trusted
    rel: complements
  - page: REQ-185-javascript-behavior-probes-can-all-skip-
    rel: complements
  - page: REQ-187-no-single-local-maintainer-command-prove
    rel: complements
  - page: REQ-188-hotspot-output-silently-drops-unavailabl
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-186: Required baseline verification executes two child suites twice

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Give each duplicated baseline child suite one owner: remove the aggregate's redundant direct prescribed-shell edge and stop requiring a second standalone shipped-reference invocation when it has no distinct mode or fixture.

## Solution summary

**Verified unchanged:** `_dev/tests/staged-skills-contract.sh` still invokes `prescribed-shell-scripts-behavior.sh` for standalone-semantics coverage.

## Worth knowing

- A required aggregate should give each identical child invocation one owner. If a nested suite already preserves the needed standalone semantics and failure propagation, a second direct edge adds runtime without adding evidence.
- Maintainer hand-back instructions should name the aggregate baseline once and reserve standalone child runs for genuinely distinct modes, fixtures, or touched-area focus.

## Back-reference

See `do-work/archive/UR-041/REQ-186-baseline-suite-single-ownership.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0ab2b79`.
