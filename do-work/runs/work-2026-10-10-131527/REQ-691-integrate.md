# Integrator brief: REQ-691 (do-work trace coverage action)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-691.

- REQ id: REQ-691. Working REQ path: `do-work/working/REQ-691-trace-coverage-action.md`. UR: UR-153 (`do-work/user-requests/UR-153/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-691-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-691-trace-coverage-action`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-691-trace-coverage-action`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:26:03Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-691-trace-coverage-action`, converted to UTC Z).
- Route C, tdd: false, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-691-probe.sh`.
- Integration order: eleventh integrator. REQ-690 (do-work status, run-status, queue-kanban open-work) released 0.305.111 just before you. Re-read VERSION before the payloads (expect 0.305.112). REQ-689 integrates last, after you.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (builder D-10: an ask covering part of its REQ counts only its own criteria).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-691-trace-coverage-action.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Released before the UR-153 members (plus REQ-688 0.305.108 and REQ-692 0.305.109 and REQ-687 0.305.110 and REQ-690 0.305.111, archived flat): REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed, archive/UR-145/), REQ-655 0.305.105, REQ-657 0.305.106, REQ-656 0.305.107 (UR-144 closed, archive/UR-144/). UR-153 order: REQ-688, REQ-692, REQ-687, REQ-690, REQ-691, then REQ-689 last (closes UR-153, carries the wave-end sweep). UR-153 stays open until REQ-689: archive flat and link any lesson bullet to the flat path; REQ-689's integrator re-points them.
- New tools on main you may use: `worktree merge` (REQ-660), `frontmatter set` and `req append-section` (REQ-659), `finalize --auto-manifest --emit` (REQ-658). The helper scripts also still work.
- Disk: the coordinator cleared the Go build cache earlier; check `df -h /System/Volumes/Data` before the gate and stop with a report if under 5 GB free.
- Recurring trap: pre-dispatch Scope lists without backticked paths made scope-drift refuse (REQ-658, REQ-692); add backticks only, change nothing else.
- A merge with no conflict can still put a paragraph in the wrong section (REQ-687 found its --kind text inside REQ-656's Revise steps): read every merged hunk of a shared prose file, not only the conflicted ones.
- Seams: your builder branched from bd56c4b0. REQ-690 added a queue-kanban open-work subcommand: skills/do-work-board/tools/queue-kanban/main.go will conflict with your request-commits subcommand; keep both. Core skills/do-work/SKILL.md: REQ-692 added a row above verify and REQ-690 one above clarify; your row goes above the capture fallback; the argument-hint line is a union of every REQ's tokens. capture.md and capture-reference.md now carry REQ-688 and REQ-692 edits; your line citations into them may have moved: re-check every file:line your diff cites.
- Board change: expect board heavy lanes (export QUEUE_KANBAN_BROWSER for browser lanes; a skip is not a pass).
