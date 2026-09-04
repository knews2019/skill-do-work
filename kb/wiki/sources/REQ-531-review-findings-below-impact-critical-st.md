---
title: "Lessons from REQ-531: Review findings below impact-critical stay in the report"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-04/REQ-531-review-findings-below-impact-critical-st.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-531: Review findings below impact-critical stay in the report

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

A review or build finding whose impact is anything other than `impact-critical` is recorded where it was found and never creates a REQ file. The maintainer reads the record and runs `do-work capture` by hand for the ones worth building. Only `impact-critical` still auto-queues.

## Solution summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified) — replaced stale automatic-follow-up pins with a named critical-only/report-only contract while preserving concurrent fast-test-budget coverage and the 8,479-line ceiling.
- `skills/do-work/actions/review-work.md` (modified) — records every finding, auto-queues only critical findings, adds report-only suffix/summary wording, and documents manual capture promotion.
- `skills/do-work/actions/work-reference.md` (modified) — retains noncritical builder discoveries in the current REQ and removes consent/test-hygiene creation paths.
- `skills/do-work/actions/work.md` (modified) — aligns review-result and discovered-task orchestration with critical-only automatic follow-ups.
- `skills/do-work/actions/capture-reference.md` (modified) — gates automatic findings before Fold-First and makes destination 4 report-only, with explicit capture and critical exceptions.
- `skills/do-work/actions/capture.md` (modified) — distinguishes explicit user capture from automatic finding routing.
- `skills/do-work/crew-members/general.md` (modified) — teaches report-only handling for noncritical discoveries.
- `skills/do-work-toolbox/crew-members/general.md` (modified) — keeps the toolbox crew mirror aligned.
- `skills/do-work/docs/review-work-guide.md` (modified) — updates user-facing review follow-up behavior.
- `skills/do-work/docs/work-guide.md` (modified) — updates user-facing builder discovery behavior.
- `skills/do-work/docs/standing-preferences.md` (modified) — removes the obsolete all-findings queue preference.
- `skills/do-work-toolbox/actions/code-review.md` (modified) — aligns the alternate review writer and summary template.
- `skills/do-work/next-steps.md` (modified) — makes post-review queue guidance conditional on critical work.

## What worked

- A named RED contract plus an isolated builder branch made the default reversal testable before prose changed, and the merge-time conflict exposed concurrent contract edits instead of losing either family.

## What didn't work

- The implementation followed the headline's broad “report only” wording past the narrower captured decision, and added phrase pins despite the UR's explicit no-new-pin constraint; independent review caught both after the green tests.

## Worth knowing

- The review contract now keeps its own noncritical findings here. To pursue one, run `do-work capture` with the complete finding line quoted as the source rather than expecting an automatic follow-up.

## Back-reference

See `do-work/archive/UR-102/REQ-531-review-findings-below-impact-critical-stay-in-the-report.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `455462be46b3170a22d331136ec5aa7f7e5a1c60`.
