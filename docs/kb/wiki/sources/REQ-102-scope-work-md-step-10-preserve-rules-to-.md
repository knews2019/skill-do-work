---
title: "Lessons from REQ-102: Scope work.md Step 10 preserve rules to every non-own entry, and pin both label-destruction paths"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-102-scope-work-md-step-10-preserve-rules-to-.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-102: Scope work.md Step 10 preserve rules to every non-own entry, and pin both label-destruction paths

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

REQ-094's review (Important finding) found that `actions/work.md`'s two Step 10 preserve rules — the wholesale-rewrite clause (~line 637) and the session-start delete clause (~line 647) — are scoped to entries "carrying another checkout's `writer:` label", which silently excludes the **label-less report-only** case. A label-less entry in a clean, committed checkpoint survives Crash Recovery's report-only branch (`actions/work-reference.md`, label-less legacy bullet) but then satisfies "no entry carrying another checkout's label remains", so the session-start delete removes it — and the next run classifies that `working/` REQ as "not named there" and ages it into the three-hour takeover ladder, which is exactly what the report-only branch refused.

## Solution summary

Step 10's checkpoint rewrite and session-start delete now preserve every in-progress entry this checkout did not write — label-less legacy entries included, closing the takeover-ladder re-entry hole from REQ-094's review. Contract suite pins both destruction paths.

## What worked

Proving assertion non-vacuity by running the full suite against `git show HEAD:<file>` and diffing the FAIL set — exactly two new lines, nothing else moved.

## Worth knowing

This closed the *second and third* copies of a scoping condition whose first copy was canonical — and turned up a fourth (the template comment in work-reference.md, routed to REQ-096). Echo sites that paraphrase a canonical condition drift; echoes should quote it.

## Back-reference

See `do-work/archive/UR-018/REQ-102-step10-preserve-rule-scoping.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `44d4563`.
