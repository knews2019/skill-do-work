---
title: "Lessons from REQ-189: Canonicalize ai-report and the shared completed-work evidence contract"
type: source-summary
topic_cluster: presentation-and-reporting
sources: [raw/processed/2026-09-01/REQ-189-canonicalize-ai-report-and-the-shared-co.md]
related:
  - page: REQ-190-reduce-present-work-to-portfolio-only-be
    rel: complements
  - page: REQ-191-extract-an-explicit-standalone-present-v
    rel: complements
  - page: REQ-192-migrate-completed-work-presentation-rout
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-189: Canonicalize ai-report and the shared completed-work evidence contract

Part of the [[concept-completed-work-presentation]] cluster.

## What the REQ was about

Make `do-work-toolbox ai-report [UR|REQ|most recent]` the single canonical detailed presentation for one completed UR or REQ, for both visual and non-visual work. Extract the archive-reading, evidence, and output-safety rules shared with `present-video` into one completed-work presentation reference instead of duplicating them across action files.

## Solution summary

Extracted shared completed-work resolution and evidence rules into one toolbox reference, narrowed `ai-report` to its report-specific workflow, preserved its stronger visual verification machinery, and added a first-class non-visual evidence path without introducing alternate briefs, explainers, or video behavior.

## What worked

- Subtracting duplicated target/evidence prose before extracting one shared reference kept the report action focused on presentation-specific judgment.
- Naming visual and non-visual evidence modes explicitly preserved the mature screenshot machinery without forcing backend and refactor work into screenshot-shaped output.

## What didn't work

- The first shared resolver described exact archive lookup without inheriting the suite-wide ID-token normalization contract; a canonical reference must also inherit every upstream input grammar it consumes.
- The retained image-generation block still created its public directory before success was known, contradicting the newly conditional bundle shape.

## Worth knowing

When an instruction says an output directory exists only on success, audit the prescribed shell block itself—not just surrounding prose—for eager `mkdir`. Completed-work ID readers inherit `work-reference.md`'s Target ID Resolution contract even when their search locations differ.

## Back-reference

See `do-work/archive/UR-042/REQ-189-canonical-ai-report-and-shared-evidence-contract.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bb7ae54`.
