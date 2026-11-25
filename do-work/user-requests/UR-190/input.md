---
id: UR-190
title: 'Simplify capture and verification and scan for further instruction cuts'
created_at: 2026-11-25T22:57:02Z
requests: [REQ-758, REQ-759, REQ-760]
word_count: 17
---
## Summary
Queue the first two validated instruction changes and a separate read-only scan for additional improvements.

## Extracted Requests

| Request | Intended result |
|---|---|
| REQ-758 | Preserve the force of the request during capture |
| REQ-759 | Drop routine capture verification and self-scoring |
| REQ-760 | Scan skill instructions for new maintenance opportunities |

## Batch Constraints
Preserve the user’s words, intended result, constraints, and builder discretion. The first two requests implement the explicitly selected audit findings. The third investigates fresh opportunities beyond the original five and reports them for a maintainer decision.

## Referenced findings

1. Delete the passing-mention rule. Trim the closing capture checklists where they repeat the comparison already performed before writing the queue.
2. Remove the suggestion to run `verify-requests` after every new capture. When verification is requested, report the original clause, what the capture changed or missed, and the correction. Delete the percentage table and re-scoring step.

## Full Verbatim Input
> ```
> ok, let's capture 1 and 2, but also capture in 3 a new scan for new improvements
> ```
