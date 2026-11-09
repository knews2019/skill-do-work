---
title: "Lessons from REQ-121: One uncommitted-changes inventory, one REQ-association pass"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-121-one-uncommitted-changes-inventory-one-re.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-121: One uncommitted-changes inventory, one REQ-association pass

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Candidate B of REQ-114, split out and approved. Two primitives are currently copy-pasted as prose across several action files; each becomes one shipped script under `tools/checks/`, and the action prose calls it with its manual procedure documented as the fallback.

## Solution summary

` file list, path-match against a candidate set, tie-break on the latest `completed_at`.

## What worked

- **`$(...)` silently destroys NUL-delimited output in bash.** Any prescribed command that reaches for `git ... -z` for path safety must read it through process substitution or a pipe; capturing it in a variable removes the very delimiter that made it safe. Worth checking wherever the skill prescribes `-z`.
- **Porcelain rename records carry a second field.** `R  new\0old\0` — a consumer that reads one field per record parses the origin path as the next record's status bytes and shifts every subsequent row. Any hand-rolled `-z` parse in prose has this bug unless it explicitly discards the origin.
- **The drift was already visible before the extraction.** `commit.md` had learned that a `status: done` REQ must still associate; `inspect.md` had not. Two copies of one primitive do not stay two copies of the same primitive — that is the argument for extraction, and it was observable in the diff rather than hypothetical.

## Back-reference

See `do-work/archive/UR-022/REQ-121-uncommitted-inventory-and-req-association-scripts.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `167b0ae`.
