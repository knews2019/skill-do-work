---
title: "Lessons from REQ-453: Keep targeted UR dependency closures in the run"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-04/REQ-453-keep-targeted-ur-dependency-closures-in.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-453: Keep targeted UR dependency closures in the run

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

Keep every pending member of a targeted user-request dependency closure in the authoritative run set, then re-evaluate downstream members after their prerequisites integrate. Do not falsely declare a dependent concurrently runnable during fan-out, and do not silently leave it behind when the targeted workflow stops after the returned set.

## Solution summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/nextselection/next_selection_test.go` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

## What worked

- A post-evaluation fixed point retained dependency chains without widening the result schema, and mutation-based action tests pinned the selector/action ownership boundary.

## What didn't work

- Replaying the initial fan-out bound before projecting frozen membership let a newly captured out-of-ledger root consume the only slot and falsely strand retained work.

## Worth knowing

- In targeted replay, observe the complete canonical ready set first, project it onto frozen membership second, and apply the saved scheduling bound last. Bounding before projection changes semantics, not just throughput.

## Back-reference

See `do-work/archive/REQ-453-keep-targeted-ur-dependency-closures-in-run.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `62ef510d`.
