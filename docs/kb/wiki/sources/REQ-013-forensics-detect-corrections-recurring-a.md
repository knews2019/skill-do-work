---
title: "Lessons from REQ-013: forensics: detect corrections recurring across archived REQ Lessons Learned"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-013-forensics-detect-corrections-recurring-a.md]
related:
  - page: REQ-014-add-crew-members-maintenance-md-codifyin
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-013: forensics: detect corrections recurring across archived REQ Lessons Learned

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

Extend `actions/forensics.md` (read-only pipeline diagnostics) with a check that scans
the `## Lessons Learned` sections of archived REQs in `do-work/archive/` and flags any
correction or lesson theme that recurs across multiple REQs as a single harness-level
finding. Today nothing aggregates lessons across REQs — `actions/kb-lessons-handoff.md`
only pulls one REQ's lessons at a time into the KB inbox. This imports the Agent
Maintenance Loop's "the same correction across multiple runs is signal that the harness
is teaching the wrong thing" into do-work's own diagnostics.

## Solution summary

Added a tenth read-only forensics check that aggregates `## Lessons Learned` across all archived REQs (enumerated with `find do-work/archive -name 'REQ-*.md'`, which surfaces both loose and UR-nested files), groups lessons by heuristic theme, and reports themes recurring across 2 distinct REQs as Info/watch and 3+ as Warning/strong-signal, each with REQ IDs and a "fix the harness, not the next run" pointer. Updated the two closed/teaser enumerations of the check list so the new check isn't orphaned.

## What worked

- Validating the Red-Green Proof by *actually running* the prescribed `find` against the live archive — not just reasoning about it — proved both the loose-vs-nested handling and that the documented glob trap is real (`ls` drops `UR-001/REQ-012`). Empirical proof beat prose assertion.

## What didn't work

- The second proof pair (REQ-008+REQ-010) is a looser theme match than the first; "read the whole list before editing" and "read complementary source files before editing" only group together under a broad "read full context first" theme. The heuristic grouping is doing real work here — worth keeping the degenerate-grouping Red Flag.

## Worth knowing

- `do-work/` is in `.git/info/exclude` (local), so queue/working files are git-ignored and archived REQs must be force-added (`git add -f`) to enter the commit — same pattern the prior archived REQs followed. Future forensics checks that read across REQs should reuse `find do-work/archive -name 'REQ-*.md'` rather than a glob.

## Back-reference

See `do-work/archive/UR-002/REQ-013-recurring-correction-detector.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `e2e506a`.
