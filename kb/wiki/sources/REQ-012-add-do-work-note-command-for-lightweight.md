---
title: "Lessons from REQ-012: Add do-work note command for lightweight roadmap notes"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-012-add-do-work-note-command-for-lightweight.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-012: Add do-work note command for lightweight roadmap notes

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Add a `do-work note <text>` command that appends a lightweight, dated note to `do-work/notes.md`. The `do-work roadmap` action reads this file and renders a **Notes** section at the top of its output, before the REQ queue. Notes have no frontmatter, no RED/GREEN proof, no domain — they are ephemeral next-step hints that a user deletes directly from `notes.md` when no longer relevant.

## Solution summary

Added a lightweight `do-work note <text>` channel that appends a dated hint to `do-work/notes.md`; `do-work roadmap` surfaces those hints at the top of its survey. Notes are ephemeral working-tree-only data (no UR/REQ, no schema, no commit) that the user deletes by hand.

## What worked

- The note action is prose, but its file logic (append / render / delete / skip) is concrete enough to *execute* in bash as a real GREEN test — simulating the four states gave genuine behavioral evidence without a test harness.

## What didn't work

- The REQ's suggested "near roadmap" routing placement collided with the priority-cross-reference fragility REQ-005 fixed (priorities 28/29 are referenced by number in prose). Priority 31 — after the last keyword, before the descriptive-content fallback — was the only safe slot (D-01).

## Worth knowing

- *(corrected post-archival by d528aec — the original lesson here generalized falsely from this source repo's local setup)* In end-user installs, `do-work/` is the **committable Trail of Intent**: `notes.md`, URs, and REQs are committed in the user's normal flow, and only `do-work/pipeline.json` and `do-work/runs/` are git-excluded (via the shipped `.gitignore`). The note action itself never commits, but the file is committable. The `git add -f` workaround was only needed because *this source repo* keeps `do-work/` untracked via a local `.git/info/exclude` entry — that does not generalize and `-f` should not be taught as the normal path (it masks real ignore rules).

## Back-reference

See `do-work/archive/UR-001/REQ-012-note-command-roadmap-notes.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5d048cb`.
