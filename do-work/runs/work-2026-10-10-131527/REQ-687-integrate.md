# Integrator brief: REQ-687 (ai-report --kind proposal and root-cause)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-687.

- REQ id: REQ-687. Working REQ path: `do-work/working/REQ-687-ai-report-kind-proposal-root-cause.md`. UR: UR-153 (`do-work/user-requests/UR-153/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-687-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:25:22Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-687-ai-report-proposal-root-cause-kinds`, converted to UTC Z).
- Route B, tdd: false, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-687-probe.sh`.
- Integration order: ninth integrator. REQ-692 released 0.305.109 just before you. Re-read VERSION before the payloads (expect 0.305.110).
- Not the wave's last successful integration. Do not pass the wave-end fact. Your diff redefines what ai-report accepts (unfinished work via --kind): tell your reviewer to run the Restatement Sweep over every place that says ai-report presents only completed work (the builder named architecture-report.md:129 and two intro lines).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-687-ai-report-kind-proposal-root-cause.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Released before the UR-153 members (plus REQ-688 0.305.108 and REQ-692 0.305.109, archived flat): REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed, archive/UR-145/), REQ-655 0.305.105, REQ-657 0.305.106, REQ-656 0.305.107 (UR-144 closed, archive/UR-144/). UR-153 order: REQ-688, REQ-692, REQ-687, REQ-690, REQ-691, then REQ-689 last (closes UR-153, carries the wave-end sweep). UR-153 stays open until REQ-689: archive flat and link any lesson bullet to the flat path; REQ-689's integrator re-points them.
- New tools on main you may use: `worktree merge` (REQ-660), `frontmatter set` and `req append-section` (REQ-659), `finalize --auto-manifest --emit` (REQ-658). The helper scripts also still work.
- Disk: the coordinator cleared the Go build cache earlier; check `df -h /System/Volumes/Data` before the gate and stop with a report if under 5 GB free.
- Recurring trap: pre-dispatch Scope lists without backticked paths made scope-drift refuse (REQ-658, REQ-692); add backticks only, change nothing else.
- Your builder branched from bd56c4b0. Since then REQ-655 (index/find), REQ-657 (judge render check) and REQ-656 (revise) all landed edits to the same toolbox ai-report files. Expect conflicts in ai-report.md (the Do-NOT-use bullet, the checklist end, the $ARGUMENTS line), ai-report-reference.md, ai-report-guide.md, toolbox SKILL.md and help.md. Keep every REQ's text; ai-report must stay routed in ONE SKILL.md row (staged-skills fails on two). REQ-656 made revise skip the completed-work check: make sure your --kind rule and revise's rule do not contradict.
- Stale path in your working REQ: `do-work/working/REQ-687-ai-report-kind-proposal-root-cause.md` (about line 81) cites REQ-654's flat archive path; it now lives under `do-work/archive/UR-144/`. Fix it before qualify.
