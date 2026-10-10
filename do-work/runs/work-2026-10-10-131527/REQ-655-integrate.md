# Integrator brief: REQ-655 (ai-report index and find catalog)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-655.

- REQ id: REQ-655. Working REQ path: `do-work/working/REQ-655-ai-report-index-and-find-catalog.md`. UR: UR-144 (`do-work/user-requests/UR-144/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-655-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-655-ai-report-index-and-find`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-655-ai-report-index-and-find`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:25:37Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-655-ai-report-index-and-find`, converted to UTC Z).
- Route C, tdd: true, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-655-probe.sh`.
- Integration order: fourth integrator. Released before you: REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed; those three now live in do-work/archive/UR-145/). Re-read VERSION before the payloads (expect 0.305.105). REQ-657 (ai-report judge) integrates right after you and shares your toolbox files.
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (the ai-report bundle naming and entry-file pick order).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-655-ai-report-index-and-find-catalog.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- UR-144 stays open: REQ-656 (ai-report revise, depends_on REQ-655) is still queued, so archive flat and link your lesson bullet to the flat path.
- New tool you may use: REQ-658 shipped `finalize --auto-manifest --emit` (see work.md Step 8/9). REQ-658's integrator finalized with it and its emitted commit paths matched the helper's. Either route is fine; the helper scripts still work.
- Builder deviations to judge in Qualification: D-11 (the action runs the verb with `--format text`, because JSON output drops printed lines) and D-12 (an explicit non-alphanumeric boundary instead of \\b for linked ids). Both look sound; confirm against the code.
- Seams: toolboxcommands/commands_test.go handler count goes 7 to 8 with you; REQ-657 will set 9. The toolbox SKILL.md routing row and help.md: append your phrases, do not reflow.
