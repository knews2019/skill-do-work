---
id: UR-044
title: Modularize the framework-free queue board client
created_at: 2026-08-15T09:12:23Z
requests: [REQ-195]
word_count: 2
---

# Modularize the Framework-Free Queue Board Client

## Summary

Capture the single accepted upstream change from the immediately preceding `do-work-toolbox validate-feedback` report: split the queue-kanban browser client into ordered private closure fragments while preserving its exact framework-free runtime behavior.

## Extracted Requests

| Request | Title | Source verdict |
|---|---|---|
| REQ-195 | Modularize the framework-free queue board client | Accept — feedback item 11 |

## Batch Constraints

- Preserve provenance from consumer commit `d9237469478bd65e3574f2e80d7b57aac9148dfe` and its archived consumer REQ-197 record.
- Capture only feedback item 11. Items 1–10 were already upstream and are not new work.
- Preserve the validated runtime: one private classic-script IIFE, no new browser globals, ES modules, network requests, bundler, or user-visible behavior.
- Use a fresh upstream pre-change baseline for parity evidence; do not install the consumer hash as a permanent golden assertion.
- Do not copy consumer browser-harness paths, consumer maintainability limits, application-specific limits, consumer version/changelog history, or consumer request identifiers.
- Capture does not execute this request.

## Resolved Context From the Validated Handoff

Current upstream still has a 2,524-line, 105,430-byte `web/board.js`. The consumer shell plus eight ordered fragments reconstruct that current upstream file byte-for-byte. The accepted structure is `board.js` plus `board-core.js`, `board-filters.js`, `board-cards.js`, `board-calendar.js`, `board-testing.js`, `board-detail.js`, `board-controls.js`, and `board-clipboard.js`, assembled by an explicit Go-owned manifest through one exact internal placeholder. Static and live modes must consume the same assembled client.

The validated handoff also requires structural assembly tests, retargeting upstream-only contract checks and factual source-owner comments, a package prime update, one-time static/live browser characterization, and integration-time release regeneration. It explicitly omits the consumer Playwright path and line-ceiling shell check.

## Full Verbatim Input

do-work capture-request

---
*Captured: 2026-08-15T09:12:23Z*
