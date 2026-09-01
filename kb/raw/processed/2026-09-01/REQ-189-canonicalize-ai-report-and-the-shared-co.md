---
source_type: req_lesson
req_id: REQ-189
req_path: do-work/archive/UR-042/REQ-189-canonical-ai-report-and-shared-evidence-contract.md
date: 2026-08-15
domain: general
module: _dev/primes
tags: [general, canonicalize, report, shared, completed]
---

# Lessons from REQ-189: Canonicalize ai-report and the shared completed-work evidence contract

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

## Back-reference

See `do-work/archive/UR-042/REQ-189-canonical-ai-report-and-shared-evidence-contract.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bb7ae54`.
