---
title: "Lessons from REQ-185: JavaScript behavior probes can all skip while the board suite passes"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-185-javascript-behavior-probes-can-all-skip-.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-185: JavaScript behavior probes can all skip while the board suite passes

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

Add an explicit maintainer-strict JavaScript behavior lane so the board suite cannot report success after all four incident-sensitive Node probes skip, and convert remaining state-transition claims from source-token checks to executable behavior.

## Solution summary

**Behavior:** Maintainers can select one stable strict test entrypoint that fails if no JavaScript behavior probe actually starts. Ordinary package tests still skip when Node is unavailable, while Node-capable runs execute the production predicates, empty-state decisions, recent-window refresh, testing view copy, and confirmed testing transition.

## What worked

- An optional-tool test lane needs two separate contracts: ordinary consumers may skip unavailable probes, while the maintainer entrypoint must count attempted behavior and reject an otherwise green zero-probe run.
- Executing pure helpers is not enough when the regression lives in caller composition or a hidden-state branch. Mutation-resistant coverage must observe the production caller and each cache/render branch whose transition is part of the claim.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Back-reference

See `do-work/archive/UR-041/REQ-185-javascript-behavior-reachability.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `a7beda7`.
