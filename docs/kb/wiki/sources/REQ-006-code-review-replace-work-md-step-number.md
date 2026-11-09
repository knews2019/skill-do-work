---
title: "Lessons from REQ-006: Code review: replace work.md step-number coupling with named contracts"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-006-code-review-replace-work-md-step-number.md]
related:
  - page: REQ-001-code-review-split-actions-work-md-into-o
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-006: Code review: replace work.md step-number coupling with named contracts

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Multiple action files reference `actions/work.md`'s internal step numbers as if they were a stable interface:

- `actions/kb-lessons-handoff.md:3` — "actions/work.md (Step 7.5) and actions/review-work.md (Step 9.5)"
- `actions/review-work.md:288,300,302` — "actions/work.md's Step 7.5 handles it" (3×)
- `actions/commit.md:7` — "actions/work.md's Step 9 commits..."
- `actions/review-work.md:383` — "actions/work.md's Step 9 handles the commit"

## Solution summary

Replaced brittle "work.md Step N" cross-references with stable named-phase references throughout the repo, and gave the two handoff points explicit named-contract headers so future renumbering can't silently invalidate callers.

## What worked

- Landing REQ-001 first meant the named-contract headers slotted cleanly into the already-restructured work.md.

## What didn't work

- The acceptance grep `work\.md.*Step [0-9]+` is over-broad — it matches any line where `work.md` and a `Step N` co-occur, including a file's *own* step numbers (false positives) and the `work.md` substring inside `review-work.md`. A literal "return 0" required rewording legitimate non-coupling lines, not just the real couplings.

## Worth knowing

- The acceptance was scoped narrower (4 files) than its own grep (`actions/*.md`); honoring the grep meant expanding into capture.md/roadmap.md/CLAUDE.md (D-01). When an acceptance criterion is a repo-wide assertion, the scope list is a floor, not a ceiling.

## Back-reference

See `do-work/archive/legacy/REQ-006-named-contracts.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5410f97`.
