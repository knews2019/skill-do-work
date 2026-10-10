# Integrator brief: REQ-690 (do-work status action and run-status command)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-690.

- REQ id: REQ-690. Working REQ path: `do-work/working/REQ-690-run-status-action.md`. UR: UR-153 (`do-work/user-requests/UR-153/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-690-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-690-run-status-action`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-690-run-status-action`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:31:31Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-690-run-status-action`, converted to UTC Z).
- Route C, tdd: true, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-690-probe.sh`.
- Integration order: tenth integrator. REQ-687 released 0.305.110 just before you. Re-read VERSION before the payloads (expect 0.305.111). REQ-691 and then REQ-689 follow you.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (claim staleness classes and remedies wherever work.md, run-with-recovery or the board docs restate them).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-690-run-status-action.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Released before the UR-153 members (plus REQ-688 0.305.108 and REQ-692 0.305.109 and REQ-687 0.305.110, archived flat): REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed, archive/UR-145/), REQ-655 0.305.105, REQ-657 0.305.106, REQ-656 0.305.107 (UR-144 closed, archive/UR-144/). UR-153 order: REQ-688, REQ-692, REQ-687, REQ-690, REQ-691, then REQ-689 last (closes UR-153, carries the wave-end sweep). UR-153 stays open until REQ-689: archive flat and link any lesson bullet to the flat path; REQ-689's integrator re-points them.
- New tools on main you may use: `worktree merge` (REQ-660), `frontmatter set` and `req append-section` (REQ-659), `finalize --auto-manifest --emit` (REQ-658). The helper scripts also still work.
- Disk: the coordinator cleared the Go build cache earlier; check `df -h /System/Volumes/Data` before the gate and stop with a report if under 5 GB free.
- Recurring trap: pre-dispatch Scope lists without backticked paths made scope-drift refuse (REQ-658, REQ-692); add backticks only, change nothing else.
- A merge with no conflict can still put a paragraph in the wrong section (REQ-687 found its --kind text inside REQ-656's Revise steps): read every merged hunk of a shared prose file, not only the conflicted ones.
- Coordinator ruling on builder D-11 (YAGNI, the maintainer's standing preference): delete the two refusals the builder added that the brief did not ask for and no incident earned (board facts with no stale-claim threshold; a --run path that is not a directory). Do it as one commit on the builder branch before the first merge, record it as a decision in the REQ, and keep the tests that pin the REQ's named failures.
- Seams: your builder branched from bd56c4b0. Since then REQ-660 added a typed Worktree block and REQ-658 fields in resultmodel/result_model.go, REQ-659/658/660 touched cmd/do-work-cli/main.go and the corehelpers handler list (the handler-count test may need your command added again), and REQ-692 added a core SKILL.md routing row above verify. Keep everything; your routing row goes above clarify and the argument-hint line is a union. docs/work-guide.md is shared with REQ-689 (later).
- You also change the board (queue-kanban open-work). The board has no version file of its own; expect board heavy lanes from plan-heavy-verification (export QUEUE_KANBAN_BROWSER for browser lanes; a skip is not a pass).
