---
title: "Lessons from REQ-176: Implement the maintainability-audit action in do-work-toolbox"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-176-implement-the-maintainability-audit-acti.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-176: Implement the maintainability-audit action in do-work-toolbox

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Author a new do-work-toolbox action, `maintainability-audit`, from the validated draft spec in UR-040: a grounded, interactive, read-only codebase maintainability audit (measured metrics with user-calibrated bands → hotspot-scoped judgment → root-cause classes → persistent report with cross-run deltas), whose findings feed `do-work-toolbox validate-feedback`. Ship it as an action + reference companion, routed in SKILL.md and listed in help.

## Solution summary

Authored the maintainability-audit action + reference companion encoding all 22 validated requirements (grounding → calibration gate → measured metrics via audit-metrics with manual fallbacks → hotspot-scoped judgment → root-cause classes → persistent do-work/audits/ report with deltas → loop footer into do-work-toolbox validate-feedback), routed it in the toolbox with the `audit codebase` trigger takeover, and extended the route-count contract.

## What worked

**What worked:** Plan-first with the traceability table — all 22 requirements landed on the first build pass and the reviewer confirmed 22/22 with zero remediation. Having the Plan agent pre-verify the contract suites' exact assertions (route-count array, noun checks, link parser) meant no suite surprises after authoring.
**What didn't:** The capture-seeded write_set missed three files the plan surfaced (`code-review.md`, core `help.md`, `staged-skills-contract.sh`) — routing takeovers always touch the OLD owner's Use-when text and the route-count contract, not just the router. Prescribed blocks with per-metric band flags are easy to leave incomplete (words flags omitted where lines flags were present) — bands-only-from-flags means an omitted placeholder silently loses a whole metric's bands.
**Worth knowing:** The `## Instances` heading in finding templates must ship demoted (`#### Instances`) or it terminates the pasteable `## Findings` section (D-03). The environmental process-tree probe failure (sandbox-only) is the recorded baseline for suite runs in this session — surfaces byte-identical to origin/main.

## Back-reference

See `do-work/archive/UR-040/REQ-176-implement-maintainability-audit-action.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f845e1c`.
