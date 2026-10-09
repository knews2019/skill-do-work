# Run Manifest — work-2026-10-09-160815

Run dir: do-work/runs/work-2026-10-09-160815/
Concurrency: 3 (UR-143 wave REQ-650 + REQ-651 + REQ-652; REQ-653 waits on REQ-652 and runs as a second wave; coordinator shape: coordinator pre-dispatch, one builder per REQ in its own worktree, one integrator per REQ in series)
Status: integrating (REQ-652 0.305.85, REQ-650 0.305.86 done; REQ-651 integrating; REQ-653 built, integrates last and closes UR-143)
Set aside: none
Hand-back emphasis note: none

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-650 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-650-the-disk-probe-measures-the-repo-root-only | do-work/runs/work-2026-10-09-160815/REQ-650-handback.md | handback received b504b4a3; merged 674fa943 (review 98%, two report-only nits); gate + three heavy lanes incl. browser green; finalized 1c4b5f8c as 0.305.86, archived flat | 2026-10-09T16:11:57Z |
| REQ-651 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-651-ai-report-on-board-activity-correlation-determinism | do-work/runs/work-2026-10-09-160815/REQ-651-handback.md | handback received 5ac41434 (one new ai-reports bundle, 4 files); integrator dispatched third (no release) | 2026-10-09T16:11:57Z |
| REQ-652 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated | do-work/runs/work-2026-10-09-160815/REQ-652-handback.md | handback received 8f6f6e89; merged 92a6951c, review fix 0a1a8da0, final merge e8b3687c (re-review 98%); gate + staged-skills green; finalized 1cef71b7 as 0.305.85, archived flat | 2026-10-09T16:11:57Z |
| REQ-653 | builder agent in worktree (dispatched by the coordinator, wave 2) | worktree-agent-REQ-653-fan-out-orchestration-prose-companion-reference | do-work/runs/work-2026-10-09-160815/REQ-653-handback.md | handback received d1bf4855 (18 files, +225/-144, four test scripts green; work.md 12,706 / work-reference.md 17,046 words); integrator fourth and last, closes UR-143 | 2026-10-09T16:34:28Z |
