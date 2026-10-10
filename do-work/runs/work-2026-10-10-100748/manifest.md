# Run Manifest — work-2026-10-10-100748

Run dir: do-work/runs/work-2026-10-10-100748/
Concurrency: 4 (UR-152 batch upstream-report-accepts, eight REQs claimed in two claims of four; wave A = REQ-679, REQ-680, REQ-681, REQ-682, wave B = REQ-683, REQ-684, REQ-685, REQ-686; coordinator shape: a pre-dispatch agent writes the pre-dispatch sections, briefs and probes, one builder per REQ in its own worktree, one integrator per REQ in series, order REQ-679, REQ-680, REQ-681, REQ-682, REQ-683, REQ-684, REQ-685, REQ-686)
Status: consumed   # UR-152 closed: REQ-681 0.305.94, REQ-680 0.305.95, REQ-683 0.305.96, REQ-684 0.305.97, REQ-679 0.305.98, REQ-685 0.305.99, REQ-682 0.305.100, REQ-686 0.305.101 (1fc1e3e6); all eight archived under do-work/archive/UR-152/
Set aside: none
Hand-back emphasis note: none
Gate record (coordinator, once at the dispatch revision, then per REQ): run `bash _dev/tests/maintainer-verify.sh` on a quiet machine; for each Route B REQ (679, 682, 685) record the green gate with `advance REQ-NNN --request-path <working REQ path> --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-10-10-100748/REQ-NNN-preflight-probe.sh`. Route A REQs (680, 681, 683, 684, 686) have no pre-flight phase.

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-679 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-679-timeline-probes-fixed-data | do-work/runs/work-2026-10-10-100748/REQ-679-handback.md | dispatched | 2026-10-10T10:26:48Z |
| REQ-680 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-680-archive-fetch-absent-target-test | do-work/runs/work-2026-10-10-100748/REQ-680-handback.md | dispatched | 2026-10-10T10:26:48Z |
| REQ-681 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-681-append-section-entry-visible-sections | do-work/runs/work-2026-10-10-100748/REQ-681-handback.md | dispatched | 2026-10-10T10:26:48Z |
| REQ-682 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-682-hidden-timeline-scroll-guard | do-work/runs/work-2026-10-10-100748/REQ-682-handback.md | dispatched | 2026-10-10T10:26:48Z |
| REQ-683 | builder agent in worktree (dispatched by the coordinator, wave B) | worktree-agent-REQ-683-interview-cadence-am-pm | do-work/runs/work-2026-10-10-100748/REQ-683-handback.md | dispatched | 2026-10-10T10:28:06Z |
| REQ-684 | builder agent in worktree (dispatched by the coordinator, wave B) | worktree-agent-REQ-684-qualify-summary-messages | do-work/runs/work-2026-10-10-100748/REQ-684-handback.md | dispatched | 2026-10-10T10:28:41Z |
| REQ-685 | builder agent in worktree (dispatched by the coordinator, wave B) | worktree-agent-REQ-685-rollback-uses-open-root | do-work/runs/work-2026-10-10-100748/REQ-685-handback.md | dispatched | 2026-10-10T10:29:20Z |
| REQ-686 | builder agent in worktree (dispatched by the coordinator, wave B) | worktree-agent-REQ-686-pushed-back-finding-doc-lines | do-work/runs/work-2026-10-10-100748/REQ-686-handback.md | dispatched | 2026-10-10T10:29:55Z |
