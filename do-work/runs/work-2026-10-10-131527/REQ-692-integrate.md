# Integrator brief: REQ-692 (validate-feedback --capture and --run chain)

Read `INTEGRATOR-GUIDE.md` in this directory first; it holds every argv. This file holds the facts for REQ-692.

- REQ id: REQ-692. Working REQ path: `do-work/working/REQ-692-validate-feedback-capture-run.md`. UR: UR-153 (`do-work/user-requests/UR-153/input.md`).
- Run directory: `do-work/runs/work-2026-10-10-131527/` (run id `work-2026-10-10-131527`).
- Hand-back (landed, builder done): `do-work/runs/work-2026-10-10-131527/REQ-692-handback.md`. Read it in full: file manifest, P-A-U text, proof record, Decisions, Discovered Tasks, proposed CHANGELOG entry, proposed lesson bullet.
- Operative name (branch = worktree basename): `worktree-agent-REQ-692-validate-feedback-capture-chain`; worktree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-692-validate-feedback-capture-chain`. Base `bd56c4b0`.
- Dispatch instant (for the builder-work timing event): `2026-10-10T13:26:22Z`. The hand-back landed a while ago, so per fan-out-reference "Landed hand-back" skip the builder-work event and say so in `## Timing` notes; stamp `builder_handback_at` from the builder commit's committer date (`git log -1 --format=%cI worktree-agent-REQ-692-validate-feedback-capture-chain`, converted to UTC Z).
- Route B, tdd: false, impact-user-visible, effort-substantive. Test-gate probe: `do-work/runs/work-2026-10-10-131527/REQ-692-probe.sh`.
- Integration order: eighth integrator. REQ-688 (capture-files --example) released 0.305.108 just before you. Re-read VERSION before the payloads (expect 0.305.109).
- Not the wave's last successful integration. Do not pass the wave-end fact. Your reviewer runs the Restatement Sweep only if your diff redefines something other text restates (validate-feedback's description in toolbox help.md and the core routing row).
- Mid-run addendum text: none.
- Do not run Step 10 (`advance --checkpoint`, the loop, worktree cleanup of other REQs); the coordinator does after you report.
- Finalization: `UR_CLOSES=0`; archive path `do-work/archive/REQ-692-validate-feedback-capture-run.md`.
- Changelog: the builder's proposed entry is in the hand-back; rewrite its title so it says what shipped in plain words, verify it is unused in `CHANGELOG.md`.
- Released before the UR-153 members (plus REQ-688 0.305.108, archived flat): REQ-660 0.305.102, REQ-659 0.305.103, REQ-658 0.305.104 (UR-145 closed, archive/UR-145/), REQ-655 0.305.105, REQ-657 0.305.106, REQ-656 0.305.107 (UR-144 closed, archive/UR-144/). UR-153 order: REQ-688, REQ-692, REQ-687, REQ-690, REQ-691, then REQ-689 last (closes UR-153, carries the wave-end sweep). UR-153 stays open until REQ-689: archive flat and link any lesson bullet to the flat path; REQ-689's integrator re-points them.
- New tools on main you may use: `worktree merge` (REQ-660), `frontmatter set` and `req append-section` (REQ-659), `finalize --auto-manifest --emit` (REQ-658). The helper scripts also still work.
- Disk: the coordinator cleared the Go build cache earlier; check `df -h /System/Volumes/Data` before the gate and stop with a report if under 5 GB free.
- Seams: capture.md now has REQ-688's Step 5 changes; your one clause near line 119 is separate. The core skills/do-work/SKILL.md routing rows: REQ-689, 690 and 691 add rows later; you are first, insert where your builder put it (above the verify row).
- Contract trap: the heavy staged-skills contract fails if any shipped file contains the text "do-work validate-feedback"; the drain proves it.
