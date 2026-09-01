---
title: "Lessons from REQ-277: State the mark-label face constant's real scope at its canonical home"
type: source-summary
topic_cluster: kanban-board-and-ui
sources: [raw/processed/2026-09-01/REQ-277-state-the-mark-label-face-constant-s-rea.md]
related:
  - page: concept-kanban-board-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-277: State the mark-label face constant's real scope at its canonical home

Part of the [[concept-kanban-board-architecture]] cluster.

## What the REQ was about

REQ-265 consolidated two duplicate descent bounds into one: `durationsMeasuredMarkLabelDescentUnits` in `generate_test.go` is now the package's **single** bound for the `.durations-mark-label` face, read from **two** files. Its documentation did not follow it.

- `generate_test.go:1968-1971` still calls it "The annotation box's descent", and its block header at `:1924` still scopes the block to "panel B's slowest-day annotation, and the faces around it". Neither says the constant is now read cross-file.
- `durations_test.go:616-627` explains the cross-file dependency thoroughly — **from the consumer side only.** The canonical home says nothing.

## Solution summary

Ran the sweep the REQ requires and reported it in full; verified all three named instances are moot post-REQ-292 rather than performing edits that would have made them false; and restored the build-provenance enforcement that REQ-292 removed while four of its subjects were still standing, correcting the comment that still pointed at the deleted test.

## What worked

**What worked:** Doing the sweep mechanically instead of reading the three named lines. Enumerating every `durationsMeasured*` name against its readers took one shell loop, proved all three instances moot in one table, and — because the table made the block's contents concrete — led straight to reading the block header that turned out to be the real defect. The REQ predicted its own instances would be moot and asked for the sweep anyway; that instruction is what earned the finding.

**What didn't:** Nothing failed in this REQ's own work, but it caught a mistake made a few hours earlier in this same run. REQ-292's justification for deleting a guard — "no such constant survives" — was checked against the two files that REQ was clearing and not against the package. The lesson is narrow and repeatable: **when a deletion is justified by "nothing uses this any more", the scope of that claim has to be the scope of the thing being deleted.** The guard was package-wide; the check was file-wide.

**Worth knowing:** A removed guard is more dangerous than an absent one, because the prose that described it usually survives. Here the block header went on asserting that every measured constant's build was enforced for hours after the enforcement was gone — and it read as true, because it had been. When deleting a check, grep for what *claims* it exists.

## Back-reference

See `do-work/archive/UR-051/REQ-277-restate-the-mark-label-face-constant-at-its-home.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `54282b0`.
