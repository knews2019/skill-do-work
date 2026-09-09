# Approved rev-v2 capture scope

Source report: `ai-reports/2026-09-09_0013_do-work-improvement-review-rev-v2/index.html`
Report SHA-256: `9dcd17027123c5400ab80c29a60cb8370bd7af99706ec203813f38d7b45cb25a`
Repository revision at capture: `f0e0e224f194cb5542b19e4a459ad8131f6e8aaf`
Captured at: 2026-09-08T21:49:05Z

These are the report's exact batch instruction and three retained prompts, adopted by the user's capture instruction. Report labels map to the new IDs in UR-130 (Reconcile prose, simplify instructions, and trim shipped history).

## Ordered batch

> do-work capture-request: Capture the following three requests as one ordered batch: R2, then R1, then R3. These are report labels; resolve them to the created or folded request IDs and encode the order in depends_on. Capture only.

## R2 — REQ-622 (Correct verified prose drift)

> do-work capture-request: Reconcile the 19 open entries in do-work/prose-backlog.md against current code, canonical contracts, and recorded intent. Correct surviving prose discrepancies and record which entries are resolved, obsolete, or still need an intent decision. Prefer removing duplicate statements; preserve runtime behavior.

## R1 — REQ-623 (Reduce orchestration instruction duplication)

> do-work capture-request: Remove proven duplication from the core SKILL.md, work.md, work-reference.md, and review-work.md. Preserve command contracts, judgment, authored evidence, recovery guidance, and behavior. Verify representative normal and recovery paths and report the text reduction; make no unmeasured runtime-savings claim.

## R3 — REQ-624 (Addendum: Trim the shipped changelog to 50 entries)

> do-work capture-request: Keep the newest 50 release entries in CHANGELOG.md and its identical installed mirror. Preserve all older history in dated export-ignored archives with working links. Verify history preservation and mirror equality; retain existing release cadence and version semantics.

## Adopted decision context

The report retains R2 → R1 → R3. R1 uses R2's reconciled wording; R3's edge preserves the chosen serial order rather than a technical prerequisite. R4 is held pending a decision on the gate's responsibilities. R5 and R6 are withdrawn and remain report history, without new backlog entries. Gate and release-policy decisions do not block this batch. Read tracing is optional for duplication cleanup and is needed only to substantiate a runtime-savings claim. This batch measures text reduction and checks preserved behavior. The report does not select release batching, CLI deletion, in-flight reconciliation, or recipe-card imports.
