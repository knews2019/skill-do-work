# Run Manifest — work-2026-10-10-131527

Run dir: do-work/runs/work-2026-10-10-131527/
Concurrency: eleven builders at once (UR-145 tool REQs REQ-658/659/660 and UR-144 ai-report REQs REQ-655/657 handed over by the UR-144 handoff at c644ee47, plus the UR-153 recaptures REQ-687 to REQ-692 claimed with two `advance --fan-out 4` calls at bd56c4b0). Coordinator shape: one pre-dispatch agent per REQ (PREDISPATCH-GUIDE.md), one builder per REQ in its own worktree, one integrator per REQ in series (INTEGRATOR-GUIDE.md). Integration order: hand-back order, with REQ-660, REQ-659, REQ-658 (tool REQs) taken first whenever their hand-back has landed, and REQ-655 before REQ-657. Coordinator ruling: the UR-144 handoff's fixed order (660, 659, 658, 655, 657) existed because the tool REQs fed the recaptured coordinator REQs; those recaptures (REQ-687 to REQ-692) are now built in parallel from the same base with no depends_on edge, so a fixed order would only leave the serial integrator idle while landed hand-backs wait.
Status: building
Set aside: none
Hand-back emphasis note: none
Queued after this wave: REQ-656 (ai-report revise, depends_on REQ-655).
Gate record: `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` run by the coordinator in a detached checkout at bd56c4b0 (the claim tip; later commits before dispatch change only do-work/ files); pre-flight for each Route B/C member recorded with `advance REQ-NNN --request-path <P> --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- bash do-work/runs/work-2026-10-10-131527/REQ-NNN-preflight-probe.sh`.

| REQ | Builder | Operative name | Handback file | Status | Dispatch instant |
|-----|---------|----------------|---------------|--------|------------------|
| REQ-660 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-660-worktree-lifecycle-command | do-work/runs/work-2026-10-10-131527/REQ-660-handback.md | released 0.305.102 (500a1fc3) | 2026-10-10T13:23:49Z |
| REQ-659 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-659-frontmatter-set-append-section | do-work/runs/work-2026-10-10-131527/REQ-659-handback.md | released 0.305.103 (4f533248) | 2026-10-10T13:24:46Z |
| REQ-658 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-658-finalize-auto-manifest | do-work/runs/work-2026-10-10-131527/REQ-658-handback.md | released 0.305.104 (90a3680b), closed UR-145 | 2026-10-10T13:25:57Z |
| REQ-655 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-655-ai-report-index-and-find | do-work/runs/work-2026-10-10-131527/REQ-655-handback.md | released 0.305.105 (424cac54) | 2026-10-10T13:25:37Z |
| REQ-657 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-657-ai-report-judge-render-check | do-work/runs/work-2026-10-10-131527/REQ-657-handback.md | released 0.305.106 (bbc03236) | 2026-10-10T13:25:39Z |
| REQ-687 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds | do-work/runs/work-2026-10-10-131527/REQ-687-handback.md | released 0.305.110 (18167d00) | 2026-10-10T13:25:22Z |
| REQ-688 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-688-capture-files-example-fence-fix | do-work/runs/work-2026-10-10-131527/REQ-688-handback.md | released 0.305.108 (196690b1) | 2026-10-10T13:23:59Z |
| REQ-689 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-689-run-coordinate-mode | do-work/runs/work-2026-10-10-131527/REQ-689-handback.md | landed | 2026-10-10T13:23:28Z |
| REQ-691 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-691-trace-coverage-action | do-work/runs/work-2026-10-10-131527/REQ-691-handback.md | landed, integrating | 2026-10-10T13:26:03Z |
| REQ-692 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-692-validate-feedback-capture-chain | do-work/runs/work-2026-10-10-131527/REQ-692-handback.md | released 0.305.109 (08759290) | 2026-10-10T13:26:22Z |
| REQ-690 | builder agent in worktree (dispatched by the coordinator) | worktree-agent-REQ-690-run-status-action | do-work/runs/work-2026-10-10-131527/REQ-690-handback.md | released 0.305.111 (8ab3d327) | 2026-10-10T13:31:31Z |
| REQ-656 | builder agent in worktree (dispatched by the coordinator after REQ-655 archived; base 7fbeff30) | worktree-agent-REQ-656-ai-report-revise-sibling-bundle | do-work/runs/work-2026-10-10-131527/REQ-656-handback.md | released 0.305.107 (a6c0af0d), closed UR-144 | 2026-10-10T16:01:00Z |
