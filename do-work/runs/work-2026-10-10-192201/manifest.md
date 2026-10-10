# Run Manifest — work-2026-10-10-192201

Run dir: do-work/runs/work-2026-10-10-192201/
Concurrency: four builders at once (UR-154 review follow-ups REQ-693 to REQ-696, claimed with `advance --fan-out 4` at fd9a2378). Coordinator shape: one pre-dispatch agent per REQ (PREDISPATCH-GUIDE.md), one builder per REQ in its own worktree, one integrator per REQ in series (INTEGRATOR-GUIDE.md), in hand-back order; the last closes UR-154 and carries the wave-end sweep.
Status: building
Set aside: none
Hand-back emphasis note: none
Gate record: `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` run by the coordinator in a detached checkout at fd9a2378; pre-flight for each Route B/C member recorded with `advance REQ-NNN --request-path <P> --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-10-10-192201/REQ-NNN-preflight-probe.sh`.

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-693 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-693-worktree-missing-folder-owned-links | do-work/runs/work-2026-10-10-192201/REQ-693-handback.md | landed | 2026-10-10T19:29:01Z |
| REQ-694 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-694-append-section-and-set-guards | do-work/runs/work-2026-10-10-192201/REQ-694-handback.md | landed, integrating | 2026-10-10T19:26:34Z |
| REQ-695 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-695-run-status-missing-run-integrator | do-work/runs/work-2026-10-10-192201/REQ-695-handback.md | landed | 2026-10-10T19:27:42Z |
| REQ-696 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-696-stale-wording-sweep | do-work/runs/work-2026-10-10-192201/REQ-696-handback.md | released 0.305.114 (5bf7dbc5) | 2026-10-10T19:26:48Z |
