---
id: UR-130
title: 'Reconcile prose, simplify instructions, and trim shipped history'
created_at: 2026-09-08T21:49:05Z
requests: [REQ-622, REQ-623, REQ-624]
word_count: 15
---
# Reconcile prose, simplify instructions, and trim shipped history

## Summary
Capture the three retained requests from the rev-v2 report as one ordered batch. This invocation captures only; it does not start implementation.

## Extracted Requests
- R2 → REQ-622 (Correct verified prose drift): reconcile the 19 open prose-backlog entries.
- R1 → REQ-623 (Reduce orchestration instruction duplication): reduce proven duplication in the four core instruction files; depends on REQ-622.
- R3 → REQ-624 (Addendum: Trim the shipped changelog to 50 entries): retain the newest 50 releases in identical live changelogs and preserve older history; depends on REQ-623. Its `addendum_to` links the earlier changelog trim, REQ-022 (CHANGELOG: keep newest 20 entries live, archive the rest with a no-git pointer).

## Batch Constraints
The approved order is R2 → R1 → R3, encoded in `depends_on`; the R3 edge preserves serial ordering rather than a technical prerequisite. No overall queue priority was specified. Preserve existing behavior and release semantics. Read tracing is optional for duplication cleanup; text reduction alone is not measured runtime savings.

## Decision History
The latest report retains these three proposals. R4 remains held pending the gate decision. R5/R6 were withdrawn and remain in report history, with no new backlog entries. Release batching, CLI deletion, in-flight reconciliation, and recipe-card imports remain outside this batch. The gate and release-policy choices do not block these three requests.

## Source and Assets
- Adopted report: `ai-reports/2026-09-09_0013_do-work-improvement-review-rev-v2/index.html` (SHA-256 `9dcd17027123c5400ab80c29a60cb8370bd7af99706ec203813f38d7b45cb25a`).
- Exact approved prompts and decision context: `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.
- Backlog snapshot identifying the 19 captured items: `do-work/user-requests/UR-130/assets/prose-backlog-at-capture.md`.
- Repository revision at capture: `f0e0e224f194cb5542b19e4a459ad8131f6e8aaf`.

The short instruction below refers to that earlier report and planning. The adopted material above is preserved as context rather than presented as words in the latest message.

## Capture Decisions
The current queue contains seven unrelated requests and no eligible same-root-cause destination. Each retained report item therefore becomes a new REQ. The changelog request links its completed predecessor without changing archived records. Requests are judged user-visible to maintainers/agents through corrected guidance or a smaller shipped history. The two prose investigations require substantive judgment; the bounded history split is mechanical. All preserve behavior and use focused completion proof rather than test-first requirements.

## Full Verbatim Input
> ```
> ok, given the report and the planning so far go ahead and capture the requests
> ```
