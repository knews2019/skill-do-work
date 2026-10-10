# Integrator brief: REQ-657 (ai-report judge render-check command)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-657.

- REQ id: REQ-657. Working REQ path: `do-work/working/REQ-657-ai-report-judge-render-check-command.md`. UR: UR-144 (`do-work/user-requests/UR-144/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-657-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-657-ai-report-judge-render-check`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-657-ai-report-judge-render-check`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:25:39Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-657-ai-report-judge-render-check`, converted to UTC Z).
- Route B, tdd: true, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-657-probe.sh`.
- Integration order: fifth integrator. Released before you: REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed), REQ-655 0.305.105 (ai-report index and find, archived flat). Re-read VERSION before the payloads (expect 0.305.106).
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (the render check step in ai-report.md Step 7 and the architecture/stakeholder report steps).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-657-ai-report-judge-render-check-command.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Seams with REQ-655 (now on main): toolboxcommands/commands.go keeps both registrations; commands_test.go handler count becomes 9 (base 7, REQ-655 made it 8, your builder made it 8 from base). The toolbox SKILL.md routing row (about line 24) and help.md: union of both REQs' phrases. `skills/do-work/docs/command-line-guide.md:26`: your builder added a sentence that ai-report-judge has no recipe; REQ-655 did not touch that line.
- REQ-655's review left report-only F2: lines next to its text are stale at `skills/do-work-toolbox/actions/ai-report.md:30`, `:135` and `skills/do-work-toolbox/docs/ai-report-guide.md:21`. If your merge already touches those lines, fix them in the merge commit and say so; otherwise leave them.
- Browser tests: the probe sets DO_WORK_HEAVY_TESTS=1 and needs Chrome; a skip is not a pass.
- UR-144 stays open: REQ-656 (ai-report revise) is building now and closes UR-144 later. Archive flat; link any lesson to the flat path.
- New tools on main you may use: `finalize --auto-manifest --emit` (REQ-658), `frontmatter set` and `req append-section` (REQ-659). The helper scripts also still work.
